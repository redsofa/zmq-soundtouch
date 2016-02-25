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
	"github.com/redsofa/collector/config"
	"github.com/redsofa/logger"
)

type collector struct {
	socketManager *socketManager
	queChan       chan string
	zmqSub        *zmqSub
}

func NewCollector() *collector {
	logger.Info.Println("Setting Up collector...")
	logger.Info.Println("Using WebSocket EndPoint : ", config.ServerConfig.SocketEndPoint)
	socketManager := newSocketManager(config.ServerConfig.SocketEndPoint)
	msgQueuCh := make(chan string, 100)
	return &collector{
		socketManager: socketManager,
		queChan:       msgQueuCh,
		zmqSub:        newZmqSub(msgQueuCh),
	}
}

func (this *collector) Start() {
	//Start up the socket manager in a goroutine
	go this.socketManager.start()
	//start up the ZeroMQ TCP client - no timeout
	go this.zmqSub.start(0)

	//The queChan memeber variable is where messages
	//coming from the zmqSub member variable are stuffed.
	//Whenever there is a message available on this channel,
	//queue it up on the socketManager.
	//The socketManager will then broadcast the message
	//to all WebSocket client connections
	for m := range this.queChan {
		this.socketManager.queueMessage(m)
	}
}
