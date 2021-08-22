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

//NewConfiguration read file, return configuration
func NewConfiguration(path string, port int) *Configuration {
	var configuration Configuration
	configuration.InitConfigParams(path)
	if port != 0 {
		configuration.HttpPort = port
	}
	return &configuration
}

func (c *Configuration) InitConfigParams(path string) {
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

func (c *Configuration) Params() *Configuration {
	return c
}
