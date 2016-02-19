zmq-soundtouch
---------------
A series of projects exploring secure ZeroMQ connections and SoundTouch event data collection. Result  is the live SoundTouch event data visualization seen below. 


![Screencast](https://github.com/redsofa/zmq-soundtouch/blob/master/docs/demo.gif "Screencast")


Main Project List : 
-------------------

1) - ./go/src/github.com/redsofa/`soundtouch` - Project that connects to Bose SoundTouch over WebSocket and pushes event notification messages to a secure ZeroMQ TCP PULL socket

2) - `./python/`publisher`.py` - Project that connects to secure TCP PUSH socket, receives SoundTouch notifications and broadcasts them to subscribers.

3) - ./python/`cache`.py - Project that connects to ZeroMQ TCP PUB socket, creates a cache containing a list of recent SoundTouch notification messages	and makes them available over a ZeroMQ Router socket.

4) - ./go/src/github.com/redsofa/`collector` - Project that connects to ZeroMQ TCP PUB socket, notificaton receives messages and passes them on to WebSocket clients. Project also connects to ZeroMQ Router socket to get a list of most recent messages. In addition, it serves static web content.


Project Topology :
-------------------
![Topology](https://github.com/redsofa/zmq-soundtouch/blob/master/docs/topology.png "Topology")