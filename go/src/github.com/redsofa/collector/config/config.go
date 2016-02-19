package config

import (
	"encoding/json"
	"github.com/redsofa/collector/logger"
	"os"
)

var ServerConfig Config

type Config struct {
	WebServerPort      string
	EventCollectorPort string
	LocalPrivateKey    string
	RemotePublicKey    string
}

func ReadServiceConfig(directory string) error {
	logger.Info.Println("Reading Config : ", directory, "config.json")
	f, err := os.Open(directory + "config.json")
	defer f.Close()

	if err != nil {
		logger.Error.Println(err)
		return err
	}

	decoder := json.NewDecoder(f)

	err = decoder.Decode(&ServerConfig)
	if err != nil {
		logger.Error.Println(err)
		return err
	}
	return nil
}
