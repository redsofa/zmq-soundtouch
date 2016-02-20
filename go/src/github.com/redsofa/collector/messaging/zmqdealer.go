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

	//TODO fix...
	if err != nil {
		fmt.Println(err)
	}

	return &dealer{ctx, msgChan, doneChan, errChan, client}
}

func (d *dealer) readMessages() {
	for {
		fmt.Println("buba")
		select {
		//We receive a message on the message channel
		case msg := <-d.msgChan:
			if strings.Compare(msg, "KTHXBYE") == 0 {
				fmt.Println("done in reading cache")
				d.doneChan <- true // for recieveMessages method
				return
			} else {
				fmt.Println(msg)
			}
		//We have an error
		case err := <-d.errChan:
			d.errChan <- err // for receiveMessages method
			return
		//We're done
		case <-d.doneChan:
			d.doneChan <- true // for receiveMessages method
			return
		}

	}
}

func (d *dealer) receiveMessages() {
	for {
		fmt.Println("buba jo")
		select {
		//We have an error
		case err := <-d.errChan:
			d.errChan <- err // for readMessages method
			return

		//We're done
		case <-d.doneChan:
			d.doneChan <- true // for readMessages method
			return

		// read data from socket connection (loop)
		default:
			reply, err := d.client.Recv(0)

			if err != nil {
				d.errChan <- err
			}

			d.msgChan <- reply
		}
	}
}

func (d *dealer) Start() {
	defer d.ctx.Term()

	//TODO log
	fmt.Println("Starting...")

	//TODO config
	d.client.Connect("tcp://127.0.0.1:8000")
	defer d.client.Close()

	//TODO config
	d.client.Send("ICANHAZ?", 0)

	go d.readMessages()
	d.receiveMessages()
	fmt.Println("ddddf")
}
