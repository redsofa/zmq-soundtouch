package main

import (
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"github.com/redsofa/soundtouch-simulator/config"
	"os"
	"time"
)

func main() {
	config.ReadConf("./")

	clientSecretKey := config.ClientConf.ClientSecretKey
	serverPublicKey := config.ClientConf.ServerPublicKey
	clientPublicKey := config.ClientConf.ClientPublicKey

	//Start authentication engine
	zmq.AuthSetVerbose(true)
	zmq.AuthStart()
	zmq.AuthCurveAdd("*", string(clientPublicKey))

	//  Create and connect client socket
	client, err := zmq.NewSocket(zmq.PUSH)
	checkError(err)

	defer client.Close()

	client.ClientAuthCurve(string(serverPublicKey), string(clientPublicKey), string(clientSecretKey))
	client.Connect("tcp://" + config.ClientConf.PushServerIP + ":" + config.ClientConf.PushServerPort)

	count := 0
	for {
		msg := fmt.Sprintf("[%d]", count)
		_, err = client.SendMessage(msg)
		checkError(err)
		count++
		time.Sleep(100 * time.Millisecond)
	}

	zmq.AuthStop()
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
