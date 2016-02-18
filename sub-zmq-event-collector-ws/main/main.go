package main

import (
	"../config"
	"../logger"
	"io/ioutil"
	"os"
)

const (
	SERVER_VERSION = "0.0.1"
)

func init() {
	logger.InitLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	config.ReadServiceConfig("./")
}

func main() {

	//The port our server listens on
	listenPort := config.ServerConfig.WebServerPort

	logger.Info.Printf("Sever Starting - Listing on port %s - (Version - %s)", listenPort, SERVER_VERSION)
}
