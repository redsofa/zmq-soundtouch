zmq-soundtouch
---------------
A series of projects exploring secure ZeroMQ connections and SoundTouch event data collection.


ws-push-zmq-event-collector
----------------------------
Project that connects to Bose SoundTouch over WebSocketw and pushes the WebSocket messages to a secure ZeroMQ TCP PULL socket


pull-zmq-event-collector-ws
---------------------------
Project that connects to ZeroMQ TCP PUSH socket, receives messages and passes them on to WebSocket clients


push-zmq-event-collector-tester
-------------------------------
Test project that pushes events (numbers incrementing) to ZeroMQ TCP PULL socket 