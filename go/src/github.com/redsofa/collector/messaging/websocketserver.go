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
	"github.com/redsofa/logger"
	"golang.org/x/net/websocket"
	"net/http"
	"time"
)

type webSocketServer struct {
	urlPattern string
	wsClients  map[int]*webSocketClient
	addClCh    chan *webSocketClient
	delClCh    chan *webSocketClient
	sendAllCh  chan *Payload
	doneCh     chan bool
	errCh      chan error
	timerCh    <-chan time.Time
}

func NewWebSocketServer(urlPattern string) *webSocketServer {
	logger.Info.Println("Creating new WebSocket server...")
	wsClients := make(map[int]*webSocketClient)
	addClCh := make(chan *webSocketClient)
	delClCh := make(chan *webSocketClient)
	sendAllCh := make(chan *Payload)
	doneCh := make(chan bool)
	errCh := make(chan error)

	return &webSocketServer{
		urlPattern: urlPattern,
		wsClients:  wsClients,
		addClCh:    addClCh,
		delClCh:    delClCh,
		sendAllCh:  sendAllCh,
		doneCh:     doneCh,
		errCh:      errCh,
	}
}

func (this *webSocketServer) setTimeout(seconds int) {
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

func (this *webSocketServer) Error(err error) {
	this.errCh <- err
}

func (this *webSocketServer) AddWSClient(client *webSocketClient) {
	this.addClCh <- client
}

func (this *webSocketServer) RemoveWSClient(client *webSocketClient) {
	this.delClCh <- client
}

func (this *webSocketServer) Broadcast(payload *Payload) {
	logger.Info.Println("Broadcasting ...", payload)
	this.sendAllCh <- payload
}

func (this *webSocketServer) sendAll(payload *Payload) {
	for _, client := range this.wsClients {
		client.Write(payload)
	}
}

func (this *webSocketServer) Shutdown() {
	this.doneCh <- true
}

func (this *webSocketServer) process() {

	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				this.errCh <- err
			}
		}()
		logger.Info.Println("onConnect WebSocket handler fired...")
		logger.Info.Println("WebSocket client connecting...")
		wsClient := NewWebSocketClient(ws, this)
		this.AddWSClient(wsClient)
		wsClient.Start()
	}

	http.Handle(this.urlPattern, websocket.Handler(onConnected))

	for {
		select {

		// Add a new WebSocket client
		case client := <-this.addClCh:
			logger.Info.Println("Added new WebSocket client")
			this.wsClients[client.id] = client
			logger.Info.Println("Number of conneced clients : ", len(this.wsClients))

		// Delete WebSocket client
		case client := <-this.delClCh:
			logger.Info.Println("Deleting WebSocket client")
			delete(this.wsClients, client.id)

		// Send message to all WebSocket clients
		case payload := <-this.sendAllCh:
			logger.Info.Println("Sending all :", payload)
			this.sendAll(payload)

		case err := <-this.errCh:
			logger.Error.Println("WebSocket Error : ", err)

		case <-this.doneCh:
			logger.Info.Println("WebSocket Server Shutdown")
			return
		}
	}
}

func (this *webSocketServer) Start(timeout int) {
	logger.Info.Println("Starting Web Socket Server...")

	if timeout > 0 {
		logger.Info.Println("Setting Timeout to ", timeout, " seconds")
		this.setTimeout(timeout)
	}

	go this.process()

	//Wait until we get something on the done chan
	<-this.doneCh
}
