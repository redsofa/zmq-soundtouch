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
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"github.com/redsofa/logger"
	"github.com/redsofa/soundtouch/config"
	"github.com/redsofa/soundtouch/version"
	"golang.org/x/net/websocket"
	"io"
	"io/ioutil"
	"os"
)

func init() {
	logger.InitLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	err := config.ReadConf("./")
	if err != nil {
		os.Exit(1)
	}
}

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
				logger.Error.Println(err)
				close(msgChan)
				break
			}
			logger.Error.Println("Couldn't receive msg " + err.Error())
			close(msgChan)
			break
		}
		msgChan <- msg
	}
	close(msgChan)
}

func main() {
	msgChan := make(chan string)

	logger.Info.Println(fmt.Sprintf("SoundTouch - Connecting to SoundTouch at IP : %s on port : %s - (Version - %s)",
		config.ClientConf.SoundTouchIP,
		config.ClientConf.SoundTouchPort,
		version.APP_VERSION))

	logger.Info.Println(fmt.Sprintf("Connecting to Push Server IP : %s on port : %s",
		config.ClientConf.PushServerIP,
		config.ClientConf.PushServerPort))

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

		logger.Info.Println("Sent : ", msg)
	}

	zmq.AuthStop()
}

func checkError(err error) {
	if err != nil {
		logger.Error.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
