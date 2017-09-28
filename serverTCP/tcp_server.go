// golang tsl web server

package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"./app"
)

func StartPageHandler(w http.ResponseWriter, req *http.Request) {
	app := app.App{}
	app.Init()
	body, err := ioutil.ReadAll(req.Body)
	app.CheckErrors(err)
	fmt.Println(len(body))

	resp := app.Run(body)
	w.Write([]byte(resp))
}

// func StatusHandler(w http.ResponseWriter, req *http.Request) {
// 	fmt.Println("Recieved status command...")
// 	config, err := ioutil.ReadFile("config.json")
// 	if err != nil {
// 		panic(err)
// 	}
// 	var jsonDat map[string]interface{}
// 	json.Unmarshal(config, &jsonDat)

// 	keyS, err := base64.URLEncoding.DecodeString(jsonDat["Key_store"].(string))
// 	key := []byte(keyS)

// 	body, err := ioutil.ReadAll(req.Body)
// 	if err != nil {
// 		panic(err)
// 	}
// 	var intermediate map[string]interface{}
// 	err = json.Unmarshal(body, &intermediate)
// 	crypt_main.CheckErrors(err)

// 	originalText := crypt_main.Decrypt(key, intermediate["Data"].(string))
// 	hashDec := crypt_main.GetHashFromString(originalText)
// 	crypt_main.CheckErrors(err)
// 	recievedHash, err := base64.URLEncoding.DecodeString(intermediate["Hash"].(string))
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	if crypt_main.CompareByteSlices(hashDec, recievedHash) {
// 		fmt.Println("hash sum of files equals. writing it on hard drive")
// 		w.Write([]byte("OK"))
// 		err = ioutil.WriteFile("archive.zip", []byte(originalText), 0644)
// 		crypt_main.CheckErrors(err)
// 		fmt.Println("Success!")
// 	} else {
// 		fmt.Println("Intermediate: ", intermediate["Hash"].(string))
// 		fmt.Println("hashDec: ", string(hashDec))
// 		w.Write([]byte("Error!"))
// 	}

// }

func main() {
	fmt.Println("Starting server...")
	ln, err := net.Listen("tcp", "127.0.0.1:5002")
	if err != nil {
		panic(err)
	}
	conn, _ := ln.Accept()
	buf := make([]byte, 2048)
	defer conn.Close()

	for {
		conn.Read(buf)
		fmt.Println("Message recieved: ", string(buf))
		newmes := "OK"
		conn.Write([]byte(newmes + "\n"))
		break
	}
}
