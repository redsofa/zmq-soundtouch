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
//zmq "github.com/pebbe/zmq4"
//"github.com/redsofa/collector/config"
//"log"
//"runtime"
)

//TODO : Implement
type zmqSub struct {
}

func NewZmqSub() {

}

/*
func PushZmqMessages(msgChan chan string, wsServer WebSocketServer) {
	for msg := range msgChan {
		log.Println("RECEIVED ON ZMQ SOCKET (from msgChan) : ", msg)
		wsMsg := &Payload{"SoundTouch", msg}
		wsServer.SendAll(wsMsg)
	}
}

func BindToZmqTcpPort(msgChan chan string) {
	port := config.ClientConf.EventCollectorPort

	log.Println("Binding to Local ZMQ Port :", port)

	checkError := func(err error) {
		if err != nil {
			log.SetFlags(0)
			_, filename, lineno, ok := runtime.Caller(1)
			if ok {
				log.Fatalf("%v:%v: %v", filename, lineno, err)
			} else {
				log.Fatalln(err)
			}
		}
	}

	//Get keys from config file
	localPrivateKey := config.ClientConf.LocalPrivateKey
	remotePublicKey := config.ClientConf.RemotePublicKey

	//Start authentication engine
	zmq.AuthSetVerbose(true)
	zmq.AuthStart()
	zmq.AuthAllow("*")
	zmq.AuthCurveAdd("*", remotePublicKey)

	//Create and bind server socket
	server, _ := zmq.NewSocket(zmq.PULL)
	//Need the PUSH client's private key
	server.ServerAuthCurve("*", localPrivateKey)
	server.Bind("tcp://*:" + config.ClientConf.EventCollectorPort)

	defer server.Close()
	defer close(msgChan)

	//When we receive messages on the TCP Socket,
	//put it on the msgChan channel
	for {
		msg, err := server.Recv(0)
		checkError(err)
		msgChan <- msg
	}
	//Stop Auth engine
	zmq.AuthStop()
}
*/
