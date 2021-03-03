package service

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// Config holds the list of services to run/keep track of
var Config []Service

// Service is an application/binary managed by PicoInit
type Service struct {
	Name       string `json:"name"`
	WorkingDir string `json:"workdir"`
	Command    string `json:"cmd"`
	Restart    string `json:"restart"`
}

func init() {
	// load and parse config file
	jsonFile, err := os.Open("./picoinit_config.json")
	if err != nil {
		log.Println("./picoinit_config.json not found.")
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &Config)
	if len(Config) < 1 {
		log.Fatalln("Please add some services to ./picoinit_config.json")
	}
}
