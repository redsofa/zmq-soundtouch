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

//TODO: Might want to return payload slice and have timeout option

package messaging

import (
	zmq "github.com/pebbe/zmq4"
	"github.com/redsofa/logger"
	"os"
	"strings"
)

//TODO: Config
const (
	ROUTER_URL        = "tcp://127.0.0.1:8000"
	CACHE_END_TOKEN   = "KTHXBYE"
	CACHE_START_TOKEN = "ICANHAZ?"
)

type dealer struct {
	ctx    *zmq.Context
	msgCh  chan string
	doneCh chan bool
	errCh  chan error
	client *zmq.Socket
}

func NewDealer() *dealer {
	ctx, _ := zmq.NewContext()

	msgCh := make(chan string)
	doneCh := make(chan bool)
	errCh := make(chan error)

	client, err := ctx.NewSocket(zmq.DEALER)

	if err != nil {
		logger.Error.Println("Error openinng DEALER	socket", err)
		os.Exit(1)
	}

	return &dealer{ctx, msgCh, doneCh, errCh, client}
}

func (this *dealer) processMessages() {
	logger.Info.Println("Firing up dealer processMessages goroutine")
	for {
		select {
		//We receive a message on the message channel
		case msg := <-this.msgCh:
			logger.Info.Println("Processing Message: " + msg)
			//TODO : Relay cached messages over to WebSocket client. Should be method call on websocketclient

		//We have an error on the error channel, log it and say we're done
		case err := <-this.errCh:
			logger.Error.Println("Error : ", err)
			this.doneCh <- true

		//We're done, clean up and say we're done
		case <-this.doneCh:
			logger.Info.Println("Done Processing Messages")
			this.ctx.Term()
			this.client.Close()
			return
		}
	}
}

func (this *dealer) receiveMessages() {
	logger.Info.Println("Firing up dealer receiveMessages goroutine")
	for {
		select {

		//Read data from socket connection (loop)
		default:
			reply, err := this.client.Recv(0)

			if err != nil {
				this.errCh <- err
				return
			}
			if strings.Compare(reply, CACHE_END_TOKEN) == 0 {
				logger.Info.Println("Done reading cache")
				this.doneCh <- true //to notify processMessages()
			} else {
				this.msgCh <- reply //will be picked up by processMessages()
			}
		}
	}
}

func (this *dealer) Start() {
	logger.Info.Println("Starting Dealer ...")
	this.client.Connect(ROUTER_URL)
	logger.Info.Println("Sending request for cache conntents")
	this.client.Send(CACHE_START_TOKEN, 0)

	go this.processMessages()
	go this.receiveMessages()

	//Wait until we get something on the done channel
	<-this.doneCh
}
