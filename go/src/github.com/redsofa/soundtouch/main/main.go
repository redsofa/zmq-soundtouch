package main

import (
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"github.com/redsofa/soundtouch/config"
	"golang.org/x/net/websocket"
	"io"
	"os"
)

//TODO : need to be able to control websocket connection (with other channel?)
func connectWS(msgChan chan string, soundTouchIp, soundTouchPort string) {
	conn, err := websocket.Dial("ws://"+soundTouchIp+
		":"+soundTouchPort, "gabbo", "http://redsofa.ca")
	checkError(err)

	var msg string
	for {
		err := websocket.Message.Receive(conn, &msg)
		if err != nil {
			if err == io.EOF {
				fmt.Println(err)
				close(msgChan)
				break
			}
			fmt.Println("Couldn't receive msg " + err.Error())
			close(msgChan)
			break
		}
		msgChan <- msg
	}
	close(msgChan)
}

func main() {
	msgChan := make(chan string)

	config.ReadConf("./")

	go connectWS(msgChan, config.ClientConf.SoundTouchIP, config.ClientConf.SoundTouchPort)

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

	//While we're getting messages on the msgChan channel send them to the push sever
	for msg := range msgChan {
		_, err = client.SendMessage(msg)
		checkError(err)

		fmt.Println("Sent : ", msg)
	}

	zmq.AuthStop()
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
