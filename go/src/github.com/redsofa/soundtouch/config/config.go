package config

import (
	"encoding/json"
	"log"
	"os"
)

var ClientConf Config

type Config struct {
	SoundTouchIP    string
	SoundTouchPort  string
	PushServerIP    string
	PushServerPort  string
	ClientSecretKey string
	ServerPublicKey string
	ClientPublicKey string
}

func ReadConf(directory string) {
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
