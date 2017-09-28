//client go

package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"./../crypt_main"
)

const (
	_URL           = "https://127.0.0.1:5002/status"
	CONFIGFILENAME = "config.json"
)

type Client struct {
	conf Config
}

type Config struct {
	KeyStore        []byte `json:"KeyStore"`
	Command         string `json:"Command"`
	FileToSend      string `json:"FileToSend"`
	DirectoryToSend string `json:"DirectoryToSend"`
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
	fmt.Println("Starting client...")

	c := Client{}
	c.conf.Init()

	conf := Config{}
	conf.Init()

	decryptText := c.getRequestData()

	resp := sendRequest(decryptText, "POST")
	fmt.Println("response status: ", resp.StatusCode)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Body: ", string(body))

	for {
		rsp = sendRequest()
	}

}

func sendRequest(data string, method string) http.Response {
	req, err := http.NewRequest("POST", _URL, bytes.NewBuffer([]byte(data)))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "plain/text")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	return *resp

}

func (c *Client) getJsonPayload(commandToSend string, hashOfData []byte, dataStr []byte) []byte {
	data := DataSctruct{
		Hash: base64.URLEncoding.EncodeToString(hashOfData),
		File: base64.URLEncoding.EncodeToString(dataStr)}
	pload := Payload{Command: commandToSend, Data: data}

	jsonData, err := json.Marshal(pload)
	crypt_main.CheckErrors(err)

	return jsonData
}

func (c *Client) getRequestData() string {
	var (
		data []byte
		err  error
	)

	switch c.conf.Command {
	case "SaveFile":
		data, err = ioutil.ReadFile(c.conf.FileToSend)
		CheckErrors(err)
	case "ChangePass":
		data = crypt_main.GenerateKey()
	}

	hashData := crypt_main.GetHashFromString(string(data))

	jsonPayload := c.getJsonPayload(c.conf.Command, hashData, data)
	fmt.Println(len(jsonPayload))

	decryptText := crypt_main.Encrypt(c.conf.KeyStore, string(jsonPayload))

	return decryptText
}

func (c *Config) Init() {
	file, err := ioutil.ReadFile(CONFIGFILENAME)
	CheckErrors(err)

	var intermediate map[string]interface{}
	err = json.Unmarshal(file, &intermediate)

	c.Command = intermediate["Command"].(string)
	c.KeyStore, err = base64.URLEncoding.DecodeString(intermediate["Key_store"].(string))
	CheckErrors(err)
	c.DirectoryToSend = intermediate["DirectoryToSend"].(string)
	c.FileToSend = intermediate["FileToSend"].(string)

}

func CheckErrors(err error) {
	if err != nil {
		panic(err)
	}
}
