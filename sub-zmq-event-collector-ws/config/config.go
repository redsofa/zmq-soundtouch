package config

import (
	"../logger"
	"encoding/json"
	"log"
	"os"
)

var ServerConfig Config

type Config struct {
	WebServerPort      string
	EventCollectorPort string
	LocalPrivateKey    string
	RemotePublicKey    string
}

func ReadServiceConfig(directory string) {
	logger.Info.Println("Reading Config : ", directory, "config.json")
	f, err := os.Open(directory + "config.json")
	defer f.Close()

	if err != nil {
		log.Println(err)
	}

	decoder := json.NewDecoder(f)

	err = decoder.Decode(&ServerConfig)
	if err != nil {
		log.Println(err)
	}
}
