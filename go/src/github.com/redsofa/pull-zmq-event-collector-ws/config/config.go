package config

import (
	"encoding/json"
	"log"
	"os"
)

var ClientConf Config

type Config struct {
	WebServerPort      string
	EventCollectorPort string
	LocalPrivateKey    string
	RemotePublicKey    string
}

func ReadConf(directory string) {
	log.Println("Reading Config : ", directory, "config.json")
	f, err := os.Open(directory + "config.json")
	defer f.Close()

	if err != nil {
		log.Println(err)
	}

	decoder := json.NewDecoder(f)

	err = decoder.Decode(&ClientConf)
	if err != nil {
		log.Println(err)
	}
}
