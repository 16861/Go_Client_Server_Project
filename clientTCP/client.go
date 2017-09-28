//client go

package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"

	"./../crypt_main"
)

const (
	_URL           = "https://127.0.0.1:5002/status"
	CONFIGFILENAME = "config.json"
)

type Config struct {
	KeyStore        []byte `json:"KeyStore"`
	Command         string `json:"Command"`
	FileToSend      string `json:"FileToSend"`
	DirectoryToSend string `json:"DirectoryToSend"`
	IP              string `json:"IP"`
	Port            string `json:"Port`
}

type DataSctruct struct {
	Hash string `json:"Hash"`
	File string `json:"File"`
}

type Payload struct {
	Command string      `json:"Command"`
	Data    DataSctruct `json:"Data"`
}

func main() {
	fmt.Println("Client side")

	conf := Config{}
	conf.Init()

	//decrytText := conf.getRequestData()
	conn := conf.getTCPConnection()

	for {
		conn.Write([]byte("Simple message"))
		message, err := bufio.NewReader(conn).ReadString('\n')
		CheckErrors(err)
		fmt.Println(message)
		break
	}

}

func getJsonPayload(commandToSend string, hashOfData []byte, dataStr []byte) []byte {
	data := DataSctruct{
		Hash: base64.URLEncoding.EncodeToString(hashOfData),
		File: base64.URLEncoding.EncodeToString(dataStr)}
	pload := Payload{Command: commandToSend, Data: data}

	jsonData, err := json.Marshal(pload)
	crypt_main.CheckErrors(err)

	return jsonData
}

func (conf *Config) getRequestData() string {
	var (
		data []byte
		err  error
	)

	switch conf.Command {
	case "SaveFile":
		data, err = ioutil.ReadFile(conf.FileToSend)
		CheckErrors(err)
	case "ChangePass":
		data = crypt_main.GenerateKey()
	}

	hashData := crypt_main.GetHashFromString(string(data))

	jsonPayload := getJsonPayload(conf.Command, hashData, data)
	fmt.Println(len(jsonPayload))

	decryptText := crypt_main.Encrypt(conf.KeyStore, string(jsonPayload))

	return decryptText
}

func (c *Config) getTCPConnection() net.Conn {
	conn, err := net.Dial("tcp", c.IP+":"+c.Port)
	CheckErrors(err)
	return conn
}

func (c *Config) Init() {
	file, err := ioutil.ReadFile(CONFIGFILENAME)
	CheckErrors(err)

	var intermediate map[string]interface{}
	err = json.Unmarshal(file, &intermediate)

	c.Command = intermediate["Command"].(string)
	c.KeyStore, err = base64.URLEncoding.DecodeString(intermediate["KeyStore"].(string))
	CheckErrors(err)
	c.DirectoryToSend = intermediate["DirectoryToSend"].(string)
	c.FileToSend = intermediate["FileToSend"].(string)
	c.IP = intermediate["IP"].(string)
	c.Port = intermediate["Port"].(string)

}

func CheckErrors(err error) {
	if err != nil {
		panic(err)
	}
}
