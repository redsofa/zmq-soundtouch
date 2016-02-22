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

package messaging

import (
	zmq "github.com/pebbe/zmq4"
	"github.com/redsofa/logger"
	"os"
	//"github.com/redsofa/collector/config"
	//"runtime"
)

const (
	PUB_URL = "tcp://127.0.0.1:7001"
)

type zmqSub struct {
	ctx      *zmq.Context
	msgChan  chan string
	doneChan chan bool
	errChan  chan error
	client   *zmq.Socket
}

func NewZmqSub() *zmqSub {
	ctx, _ := zmq.NewContext()

	msgChan := make(chan string)
	doneChan := make(chan bool)
	errChan := make(chan error)

	client, err := ctx.NewSocket(zmq.SUB)

	if err != nil {
		logger.Error.Println("Error openinng DEALER	socket", err)
		os.Exit(1)
	}

	return &zmqSub{ctx, msgChan, doneChan, errChan, client}
}

func (this *zmqSub) processMessages() {
	logger.Info.Println("Firing up zmqSub processMessages loop")
	for {
		select {
		//We receive a message on the message channel
		case msg := <-this.msgChan:
			logger.Info.Println("Processing Message: " + msg)
			//TODO : Relay cached messages over to WebSocket client. Should be method call on websocketclient

		//We have an error on the error channel
		case err := <-this.errChan:
			logger.Error.Println("Error : ", err)
			return
		//We're done
		case <-this.doneChan:
			logger.Info.Println("Done Processing Messages")
			return
		}
	}
}

func (this *zmqSub) receiveMessages() {
	logger.Info.Println("Firing up zmqSub receiveMessages loop")
	for {
		select {
		//We have an error
		case err := <-this.errChan:
			logger.Error.Println("Error :", err)
			this.errChan <- err // to notify processMessages()
			return

		//Read data from socket connection (loop)
		default:
			reply, err := this.client.Recv(0)

			if err != nil {
				this.errChan <- err
				return
			}

			this.msgChan <- reply //will be picked up by processMessages()

		}
	}
}

func (this *zmqSub) Start() {
	defer this.ctx.Term()

	logger.Info.Println("Starting ZMQ Sub ...")

	this.client.Connect(PUB_URL)
	this.client.SetSubscribe("")
	defer this.client.Close()

	go this.processMessages()
	this.receiveMessages()
}

/*
func PushZmqMessages(msgChan chan string, wsServer WebSocketServer) {
	for msg := range msgChan {
		log.Println("RECEIVED ON ZMQ SOCKET (from msgChan) : ", msg)
		wsMsg := &Payload{"SoundTouch", msg}
		wsServer.SendAll(wsMsg)
	}
}

func BindToZmqTcpPort(msgChan chan string) {
	port := config.ClientConf.EventCollectorPort

	log.Println("Binding to Local ZMQ Port :", port)

	checkError := func(err error) {
		if err != nil {
			log.SetFlags(0)
			_, filename, lineno, ok := runtime.Caller(1)
			if ok {
				log.Fatalf("%v:%v: %v", filename, lineno, err)
			} else {
				log.Fatalln(err)
			}
		}
	}

	//Get keys from config file
	localPrivateKey := config.ClientConf.LocalPrivateKey
	remotePublicKey := config.ClientConf.RemotePublicKey

	//Start authentication engine
	zmq.AuthSetVerbose(true)
	zmq.AuthStart()
	zmq.AuthAllow("*")
	zmq.AuthCurveAdd("*", remotePublicKey)

	//Create and bind server socket
	server, _ := zmq.NewSocket(zmq.PULL)
	//Need the PUSH client's private key
	server.ServerAuthCurve("*", localPrivateKey)
	server.Bind("tcp://*:" + config.ClientConf.EventCollectorPort)

	defer server.Close()
	defer close(msgChan)

	//When we receive messages on the TCP Socket,
	//put it on the msgChan channel
	for {
		msg, err := server.Recv(0)
		checkError(err)
		msgChan <- msg
	}
	//Stop Auth engine
	zmq.AuthStop()
}
*/
