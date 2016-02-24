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
	"time"
)

type collector struct {
	zmqDealer       *dealer
	zmqSub          *zmqSub
	webSocketServer *webSocketServer
}

func NewCollector() *collector {
	zmqDealer := NewDealer()
	zmqSub := NewZmqSub()
	webSocketServer := NewWebSocketServer("/socket")

	return &collector{zmqDealer: zmqDealer, zmqSub: zmqSub, webSocketServer: webSocketServer}
}

func (this *collector) Start(timeout int) {

	//The zmqDealer is the thing that get the last <x> messages
	this.zmqDealer.Start()

	time.Sleep(time.Second)

	if timeout > 0 {
		this.zmqSub.Start(timeout)
	} else {
		this.zmqSub.Start(0)
	}

	//Start up the websocket server

	//server := NewWebSocketServer("/ws")
	// if timeout > 0 {
	// 	this.webSocketServer.Start(timeout)
	// } else {
	// 	this.webSocketServer.Start(0)
	// }

}
