package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

var Config Configs

func Setup(path string) {
	log.Println("Opening config.json file")
	err := ReadFileConfigs(path)
	if err != nil {
		panic(fmt.Errorf("Can't setup configurations! %v", err))
	}
	log.Println("Successfully took secrets from conf file")
}

func ReadFileConfigs(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &Config)
	if err != nil {
		return err
	}
	return nil
}
