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
package main

import (
	"github.com/gorilla/mux"
	"github.com/redsofa/collector/config"
	"github.com/redsofa/logger"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	SERVER_VERSION = "0.0.1"
)

func init() {
	logger.InitLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	err := config.ReadServiceConfig("./")
	if err != nil {
		os.Exit(1)
	}
}

func main() {
	//The port our server listens on
	listenPort := config.ServerConfig.WebServerPort

	logger.Info.Printf("Sever Starting - Listing on port %s - (Version - %s)", listenPort, SERVER_VERSION)

	router := mux.NewRouter()
	//Our static content
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./webroot")))

	//Listen for connections and serve content

	//TODO add logger middleware
	logger.Info.Println(http.ListenAndServe(":"+listenPort, router))
}
