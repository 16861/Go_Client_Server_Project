//go crypt test file

package crypt_main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

const docFileName string = "deadbeef-0.7.2.zip"

type Config struct {
	Key_store string `json:"Key_store"`
}

func main() {
	originalText, err := ioutil.ReadFile(docFileName)
	CheckErrors(err)
	originalhash := sha256.New()
	originalhash.Write(originalText)
	hashOrigin := originalhash.Sum(nil)
	// fmt.Println("Orginal text: ", oraginalText)

	config, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	var jsonDat map[string]interface{}
	json.Unmarshal(config, &jsonDat)

	keyS, err := base64.URLEncoding.DecodeString(jsonDat["Key_store"].(string))
	key := []byte(keyS)

	// key := GenerateKey()
	// conf := Config{Key_store: base64.URLEncoding.EncodeToString(key)}
	// fmt.Println("Conf: ", conf)
	// keyJSON, err := json.Marshal(conf)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("Key JSON: ", keyJSON)
	// fmt.Println("Key len: ", string(keyJSON))
	// err = ioutil.WriteFile("config.json", keyJSON, 0644)
	// if err != nil {
	// 	panic(err)
	// }

	encryptedText := Encrypt(key, string(originalText))
	err = ioutil.WriteFile("deryptOutput", []byte(encryptedText), 0644)

	fileData, err := ioutil.ReadFile("deryptOutput")
	CheckErrors(err)

	decryptedText := Decrypt(key, string(fileData))
	decryptHasher := sha256.New()
	decryptHasher.Write([]byte(decryptedText))
	hashDec := decryptHasher.Sum(nil)
	fmt.Println("Original hash: ", string(hashOrigin))
	fmt.Println("Decrypt hash: ", string(hashDec))
	if CompareByteSlices(hashDec, hashOrigin) {
		fmt.Println("True")
	}
	err = ioutil.WriteFile("archive.zip", []byte(decryptedText), 0644)
	// fmt.Println("Decrypted text: ", decryptedText)

}

func CompareByteSlices(a []byte, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, x := range a {
		if x != b[i] {
			return false
		}
	}
	return true
}

func GetHashFromString(str string) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(str))
	return hasher.Sum(nil)
}

func Encrypt(key []byte, originalMessage string) string {
	plaintext := []byte(originalMessage)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	cipherText := make([]byte, aes.BlockSize+len(plaintext))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(cipherText)
}

func DecryptBytes(key []byte, encryptedTextb64 []byte) []byte {
	ciphertext := make([]byte, len(encryptedTextb64))
	copy(ciphertext, encryptedTextb64)
	// _, err := base64.URLEncoding.Decode(ciphertext, encryptedTextb64)
	// CheckErrors(err)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext is too short!")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext

}

func Decrypt(key []byte, encryptedTextb64 string) string {
	fmt.Println(len(encryptedTextb64), encryptedTextb64[:100])
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedTextb64)
	CheckErrors(err)
	fmt.Println(len(ciphertext))

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext is too short!")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)

}

func GenerateKey() []byte {
	key := make([]byte, 32)

	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}

	return key
}

func CheckErrors(err error) {
	if err != nil {
		panic(err)
	}
}
