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
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"strings"
)

const (
	ROUTER_URL        = "tcp://127.0.0.1:8000"
	CACHE_END_TOKEN   = "KTHXBYE"
	CACHE_START_TOKEN = "ICANHAZ?"
)

type dealer struct {
	ctx      *zmq.Context
	msgChan  chan string
	doneChan chan bool
	errChan  chan error
	client   *zmq.Socket
}

func NewDealer() *dealer {
	ctx, _ := zmq.NewContext()

	msgChan := make(chan string)
	doneChan := make(chan bool)
	errChan := make(chan error)

	client, err := ctx.NewSocket(zmq.DEALER)

	//TODO fix...log and bail
	if err != nil {
		fmt.Println(err)
	}

	return &dealer{ctx, msgChan, doneChan, errChan, client}
}

func (this *dealer) processMessages() {
	fmt.Println("Firing up processMessages loop")
	for {
		select {
		//We receive a message on the message channel
		case msg := <-this.msgChan:
			fmt.Println("Processing Message: " + msg)
		//We have an error on the error channel
		case err := <-this.errChan:
			fmt.Println("Error : ", err)
			return
		//We're done
		case <-this.doneChan:
			fmt.Println("Done Processing Messages")
			return
		}
	}
}

func (this *dealer) receiveMessages() {
	fmt.Println("Firing up receiveMessages loop")
	for {
		select {
		//We have an error
		case err := <-this.errChan:
			fmt.Println("Error :", err)
			this.errChan <- err // to notify processMessages()
			return

		//Read data from socket connection (loop)
		default:
			reply, err := this.client.Recv(0)

			if err != nil {
				this.errChan <- err
				return
			}
			if strings.Compare(reply, CACHE_END_TOKEN) == 0 {
				fmt.Println("Done reading cache")
				this.doneChan <- true //to notify processMessages()
				return
			} else {
				this.msgChan <- reply //will be picked up by processMessages()
			}
		}
	}
}

func (this *dealer) Start() {
	defer this.ctx.Term()

	//TODO log
	fmt.Println("Starting Dealer ...")

	//TODO config
	this.client.Connect(ROUTER_URL)
	defer this.client.Close()

	fmt.Println("Sending request for cache conntents")
	this.client.Send(CACHE_START_TOKEN, 0)

	go this.processMessages()
	this.receiveMessages()
}
