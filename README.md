zmq-soundtouch
---------------
A series of projects exploring secure ZeroMQ connections and SoundTouch event data collection. Result  is the live SoundTouch event data visualization seen below. 


![Screencast](https://github.com/redsofa/zmq-soundtouch/blob/master/docs/demo.gif "Screencast")


Project List : 
--------------

1) - `ws-push-zmq-event-collector` :
Project that connects to Bose SoundTouch over WebSocket and pushes the WebSocket messages to a secure ZeroMQ TCP PULL socket


2) - `pull-zmq-event-collector-ws` : 
Project that connects to ZeroMQ TCP PUSH socket, receives messages and passes them on to WebSocket clients


3) - `push-zmq-event-collector-tester` :
Test project that pushes events (numbers incrementing) to ZeroMQ TCP PULL socket 


Project Topology :
-------------------
![Topology](https://github.com/redsofa/zmq-soundtouch/blob/master/docs/topology.png "Topology")