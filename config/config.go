package config

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	HttpPort            int
	ConnectionHost      string
	ConnectionPort      int
	ConnectionType      string
	ChannelMessagesSize int
}

//ReadFile
func (c *Configuration) ReadFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c)
	if err != nil {
		log.Fatal(err)
	}
}

//NewConfiguration
func NewConfiguration(path string) *Configuration {
	var configuration Configuration
	configuration.ReadFile(path)
	return &configuration
}
