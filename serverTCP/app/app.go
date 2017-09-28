package app

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"..//../crypt_main"
)

const (
	configFilename = "config.json"
)

type App struct {
	Key          string `json:"Key"`
	OutputFolder string `json:"OutputFolder"`
}

func (a *App) Init() {
	config, err := ioutil.ReadFile("config.json")
	a.CheckErrors(err)

	var jsonDat map[string]interface{}
	json.Unmarshal(config, &jsonDat)
	temp, err := base64.URLEncoding.DecodeString(jsonDat["Key_store"].(string))
	a.CheckErrors(err)
	a.Key = string(temp)

	//keyS, err := base64.URLEncoding.DecodeString(jsonDat["Key_store"].(string))
	//key := []byte(keyS)
}

func (a *App) Run(incomingPayload []byte) string {
	message := make([]byte, base64.URLEncoding.DecodedLen(len(incomingPayload)))
	_, err := base64.URLEncoding.Decode(message, incomingPayload)
	a.CheckErrors(err)
	message = deleteBlankSymbold(message)
	fmt.Println(len(message))
	ioutil.WriteFile("temp", message, 0644)
	originalMessage := crypt_main.DecryptBytes([]byte(a.Key), message)
	ioutil.WriteFile("tempDec", message, 0644)

	var intermediate map[string]interface{}
	err = json.Unmarshal(originalMessage, &intermediate)
	a.CheckErrors(err)

	command := intermediate["Command"].(string)

	switch command {
	case "status":
		fmt.Println("Check status")
	case "change_secret":
		fmt.Println("Changing secret")
	case "SaveFile":
		fmt.Println("Saving file")
		data := intermediate["Data"].(map[string]interface{})
		fmt.Println(data["File"])
		file, err := base64.URLEncoding.DecodeString(data["File"].(string))
		a.CheckErrors(err)
		err = ioutil.WriteFile("incomingFile", file, 0644)
		a.CheckErrors(err)

	default:
		return "Error"
	}

	return "OK"

	// data, err := base64.URLEncoding.DecodeString(intermediate["Data"].(string))
	// hashDara, err := base64.URLEncoding.DecodeString(intermediate["Hash"].(string))
}

func (a *App) CheckErrors(err error) {
	if err != nil {
		panic(err)
	}
}

func deleteBlankSymbold(message []byte) []byte {
	lastIndex := len(message) - 1
	for message[lastIndex] == '\x00' {
		lastIndex -= 1
	}

	return message[:lastIndex+1]
}

// hashDec := crypt_main.GetHashFromString(originalText)
// crypt_main.CheckErrors(err)
// recievedHash, err := base64.URLEncoding.DecodeString(intermediate["Hash"].(string))
// if err != nil {
// 	fmt.Println(err)
// }
// if crypt_main.CompareByteSlices(hashDec, recievedHash) {
// 	fmt.Println("hash sum of files equals. writing it on hard drive")
// 	w.Write([]byte("OK"))
// 	err = ioutil.WriteFile("archive.zip", []byte(originalText), 0644)
// 	crypt_main.CheckErrors(err)
// 	fmt.Println("Success!")
// } else {
// 	fmt.Println("Intermediate: ", intermediate["Hash"].(string))
// 	fmt.Println("hashDec: ", string(hashDec))
// 	w.Write([]byte("Error!"))
// }
