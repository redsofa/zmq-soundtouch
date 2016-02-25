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

type connection struct {
	websocket *websocket.Conn
}

type socketManager struct {
	wsUrlPattern    string
	broadcastCh     chan string
	msgQueCh        chan string
	connections     map[*connection]bool
	doneCh          chan bool
	quitProcMsgCh   chan bool
	quitReceivMsgCh chan bool
	timerCh         <-chan time.Time
	errCh           chan error
}

func newSocketManager(wsUrlPattern string) *socketManager {
	retVal := socketManager{
		wsUrlPattern:    wsUrlPattern,
		broadcastCh:     make(chan string),
		msgQueCh:        make(chan string, 100),
		connections:     make(map[*connection]bool),
		doneCh:          make(chan bool),
		quitProcMsgCh:   make(chan bool),
		quitReceivMsgCh: make(chan bool),
		errCh:           make(chan error),
	}
	//Bind retVal.wsOnConnect event handler to the WS url pattern retVal.wsUrlPattern
	//What this is saying is when user hits ws://<urlpattern> in the browser,
	//execute retVal.wsOnConnect
	http.Handle(retVal.wsUrlPattern, websocket.Handler(retVal.wsOnConnect))
	return &retVal
}

func (this *socketManager) wsOnConnect(ws *websocket.Conn) {
	var err error
	var p payload

	defer func() {
		if err = ws.Close(); err != nil {
			logger.Info.Println("Websocket could not be closed", err.Error())
		}
	}()
	//Client connected to websocket
	connection := connection{ws}
	//Keep a reference to the connection
	this.connections[&connection] = true
	logger.Info.Println("WebSocket Connection event ...")
	logger.Info.Println("Number of clients connected ...", len(this.connections))

	//Get the last messages from the zmqDealer TCP pipe
	this.getLastMessages(connection)
	//Wait a bit for the last messages to send...
	time.Sleep(3 * time.Second)
	// Loop to keep websocket open
	for {
		if err = websocket.JSON.Receive(ws, &p); err != nil {
			//Socket is closed if
			logger.Info.Println("Websocket Disconnected...", err.Error())
			//Remove connection from active connections
			delete(this.connections, &connection)
			logger.Info.Println("Number of clients connected ...", len(this.connections))
			return
		}
	}
}

func (this *socketManager) getLastMessages(connection connection) {
	cache := newDealer().start()

	for _, m := range cache {
		websocket.JSON.Send(connection.websocket, payload{"Soundtouch", m})
	}

}

func (this *socketManager) broadcastMessage(pl payload) {
	for cs, _ := range this.connections {
		websocket.JSON.Send(cs.websocket, pl)
	}
}

func (this *socketManager) stop() {
	logger.Info.Println("Stopping socketManager.")
	this.setTimeout(1)
	time.Sleep(1 * time.Second)
}

func (this *socketManager) setTimeout(seconds int) {
	this.timerCh = time.NewTimer(time.Duration(seconds) * time.Second).C

	go func() {
		for {
			select {
			case <-this.timerCh:
				logger.Info.Println("Timeout fired. Stopping socketmanager.")
				this.quitReceivMsgCh <- true
				this.quitProcMsgCh <- true
				this.doneCh <- true
			}
		}
	}()
}

func (this *socketManager) processMessages() {
	logger.Info.Println("Firing up socketManager processMessages goroutine")
	for {
		select {
		//We're done processing messages
		case <-this.quitProcMsgCh:
			logger.Info.Println("Quit Processing Messages")
			return
		//We received a message on the broadcast channel, process it
		case msg := <-this.broadcastCh:
			logger.Info.Println("Processing Message: " + msg)
			wsPayload := &payload{"SoundTouch", msg}
			this.broadcastMessage(*wsPayload)

		//We have an error on the error channel, log it and say we're done
		case err := <-this.errCh:
			logger.Error.Println("Error : ", err)
			this.doneCh <- true
		}
	}
}

func (this *socketManager) queueMessage(msg string) {
	this.msgQueCh <- msg
}

func (this *socketManager) receiveMessages() {
	logger.Info.Println("Firing up zmqSub receiveMessages goroutine")
	for {
		select {
		case <-this.quitReceivMsgCh:
			logger.Info.Println("Quit Receive Messages")
			return
		default:
			for msg := range this.msgQueCh {
				logger.Info.Println("message received ... : ", msg)
				this.broadcastCh <- msg
			}
		}
	}
}

func (this *socketManager) start() {
	logger.Info.Println("Starting web socket manager...")

	go this.processMessages()
	go this.receiveMessages()

	//Wait until we get something on the done chan
	<-this.doneCh
	logger.Info.Println("Socket manager process done...")
}
