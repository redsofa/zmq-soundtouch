package messaging

import (
	zmq "github.com/pebbe/zmq4"
	"github.com/redsofa/zmq-soundtouch/pull-zmq-event-collector-ws/config"
	"log"
	"runtime"
)

func PushZmqMessages(msgChan chan string, wsServer Server) {
	for msg := range msgChan {
		log.Println("RECEIVED ON ZMQ TCP (from msgChan) : ", msg)
		wsMsg := &Message{"SoundTouch", msg}
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
