package main

import (
	"github.com/gorilla/mux"
	"github.com/redsofa/collector/config"
	"github.com/redsofa/collector/logger"
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
