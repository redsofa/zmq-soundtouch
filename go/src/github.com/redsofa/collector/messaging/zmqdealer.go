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
}

func NewDealer() *dealer {
	ctx, _ := zmq.NewContext()

	msgChan := make(chan string)
	doneChan := make(chan bool)

	return &dealer{ctx, msgChan, doneChan}
}

func (d *dealer) MakeCacheRequest() {
	d.msgChan <- "ICANHAZ"
}

func (d *dealer) Done() {
	d.doneChan <- true
}

func (d *dealer) Start() {
	defer d.ctx.Term()

	//TODO log
	fmt.Println("Starting...")
	client, err := d.ctx.NewSocket(zmq.DEALER)

	if err != nil {
		fmt.Println(err)
	}

	client.Connect("tcp://127.0.0.1:8000")
	defer client.Close()

	client.Send("ICANHAZ?", 0)

	for {
		reply, _ := client.Recv(0)
		fmt.Println("received : ", reply)
		if strings.Compare(reply, "KTHXBYE") == 0 {
			return
		}
	}

}
