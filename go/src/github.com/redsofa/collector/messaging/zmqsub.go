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
	"time"
	//"github.com/redsofa/collector/config"
	//"runtime"
)

//TODO : Config
const (
	PUB_URL = "tcp://127.0.0.1:7001"
)

type zmqSub struct {
	ctx     *zmq.Context
	msgCh   chan string
	doneCh  chan bool
	errCh   chan error
	timerCh <-chan time.Time
	client  *zmq.Socket
}

func NewZmqSub() *zmqSub {
	ctx, _ := zmq.NewContext()

	msgCh := make(chan string)
	doneCh := make(chan bool)
	errCh := make(chan error)
	client, err := ctx.NewSocket(zmq.SUB)

	if err != nil {
		logger.Error.Println("Error openinng zmq.SUB socket", err)
		os.Exit(1)
	}

	return &zmqSub{ctx: ctx,
		msgCh:  msgCh,
		doneCh: doneCh,
		errCh:  errCh,
		client: client,
	}
}

func (this *zmqSub) setTimeout(seconds int) {
	this.timerCh = time.NewTimer(time.Duration(seconds) * time.Second).C

	go func() {
		for {
			select {
			case <-this.timerCh:
				this.doneCh <- true
			}
		}
	}()
}

func (this *zmqSub) processMessages() {
	logger.Info.Println("Firing up zmqSub processMessages goroutine")
	for {
		select {

		//We receive a message on the message channel, process it
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

func (this *zmqSub) receiveMessages() {
	logger.Info.Println("Firing up zmqSub receiveMessages goroutine")
	for {
		select {

		default:
			subMsg, err := this.client.Recv(0)

			//If there's an error on receive, put it on the error channel and return
			//from this function, the processMessages method will log it and say we're done
			if err != nil {
				this.errCh <- err
				return
			}
			//put the reply on the channel
			this.msgCh <- subMsg //will be picked up by processMessages()
		}
	}
}

func (this *zmqSub) Start(timeout int) {
	logger.Info.Println("Starting ZMQ Sub ...")

	if timeout > 0 {
		logger.Info.Println("Setting Timeout to ", timeout, " seconds")
		this.setTimeout(timeout)
	}

	this.client.Connect(PUB_URL)
	this.client.SetSubscribe("")

	go this.processMessages()
	go this.receiveMessages()
	//Wait until we get something on the done chan
	<-this.doneCh
}
