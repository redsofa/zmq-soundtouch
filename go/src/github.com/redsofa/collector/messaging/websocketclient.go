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
	"errors"
	"fmt"
	"github.com/redsofa/logger"
	"golang.org/x/net/websocket"
	"io"
	"os"
)

const payloadBuffer = 100

var lastClId int = 0

type webSocketClient struct {
	id        int
	payloadCh chan *Payload
	wsConn    *websocket.Conn
	wsServer  *webSocketServer
	doneCh    chan bool
}

// Create new chat client.
func NewWebSocketClient(wsConn *websocket.Conn, wsServer *webSocketServer) *webSocketClient {

	if wsConn == nil {
		logger.Error.Println("Error : wsConn is null")
		os.Exit(1)
	}

	if wsServer == nil {
		logger.Error.Println("Error : wsServer is null")
		os.Exit(1)
	}

	lastClId++

	payloadCh := make(chan *Payload, payloadBuffer)
	doneCh := make(chan bool)

	return &webSocketClient{
		id:        lastClId,
		payloadCh: payloadCh,
		wsConn:    wsConn,
		wsServer:  wsServer,
		doneCh:    doneCh}
}

func (this *webSocketClient) Start() {
	logger.Info.Println("Starting Web Socket Client...")

	go this.write()
	go this.read()

	<-this.doneCh
}

func (this *webSocketClient) Write(payload *Payload) {
	select {
	case this.payloadCh <- payload:
	default:
		this.wsServer.RemoveWSClient(this)
		err := errors.New(fmt.Sprintf("client %d is disconnected.", this.id))
		this.wsServer.Error(err)
	}
}

func (this *webSocketClient) write() {
	logger.Info.Println("Starting Write to client goroutine")
	for {
		select {
		//send payload to WebSocket client
		case payload := <-this.payloadCh:
			logger.Info.Println("Send:", payload)
			websocket.JSON.Send(this.wsConn, payload)

		// receive done request bail
		case <-this.doneCh:
			logger.Info.Println("Client WS client connection done")
			this.wsServer.RemoveWSClient(this)
			return
		}
	}
}

func (this *webSocketClient) read() {
	logger.Info.Println("Starting Read from client goroutine")
	for {
		select {

		// read data from websocket connection
		default:
			var payload Payload
			err := websocket.JSON.Receive(this.wsConn, &payload)
			if err == io.EOF {
				this.doneCh <- true
			} else if err != nil {
				this.wsServer.Error(err)
				os.Exit(1)
			} else {
				this.wsServer.Broadcast(&payload)
			}
		}
	}
}

func (this *webSocketClient) Done() {
	this.doneCh <- true
}
