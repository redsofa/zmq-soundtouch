package main

import (
	"github.com/redsofa/zmq-soundtouch/pull-zmq-event-collector-ws/config"
	"github.com/redsofa/zmq-soundtouch/pull-zmq-event-collector-ws/messaging"
	"log"
	"net/http"
)

func main() {
	zmqChan := make(chan string)

	config.ReadConf("./")

	port := config.ClientConf.WebServerPort

	log.SetFlags(log.Lshortfile)

	// websocket server
	server := messaging.NewServer("/entry")
	go server.Listen()

	go messaging.BindToZmqTcpPort(zmqChan)
	go messaging.PushZmqMessages(zmqChan, *server)

	http.Handle("/", http.FileServer(http.Dir("webroot")))

	log.Println("Web Sever Starting - Listing on port", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
