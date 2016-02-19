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
	"log"
	"os"
)

var ClientConf Config

type Config struct {
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
