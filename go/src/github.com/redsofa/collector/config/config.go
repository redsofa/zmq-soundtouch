/*
Copyright 2016 Rene Richard

This file is part of zmq-soundtouch.

zmq-soundtouch is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

zmq-soundtouch is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with zmq-soundtouch.  If not, see <http://www.gnu.org/licenses/>.
*/

package config

import (
	"encoding/json"
	"github.com/redsofa/logger"
	"os"
)

var ServerConfig Config

type Config struct {
	WebServerPort      string
	EventCollectorPort string
	LocalPrivateKey    string
	RemotePublicKey    string
	RouterUrl          string
	CacheEndToken      string
	CacheStartToken    string
	ZmqPubURL          string
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
