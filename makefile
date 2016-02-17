GOPATH := $(HOME)/dev/git/go

all: clean dist-dir build-pull-zmq-event-collector-ws build-ws-push-zmq-event-collector

dist-dir:
	mkdir ./dist

build-pull-zmq-event-collector-ws:
	@echo Building pull-zmq-event-collector-ws
	@echo Using Go Path :
	@echo $(GOPATH)
	mkdir ./dist/pull-zmq-event-collector-ws
	go build -o ./dist/pull-zmq-event-collector-ws/pull-to-ws pull-zmq-event-collector-ws/main/main.go
	cp pull-zmq-event-collector-ws/main/config.json ./dist/pull-zmq-event-collector-ws/config.json
	cp -r pull-zmq-event-collector-ws/main/webroot ./dist/pull-zmq-event-collector-ws/

build-ws-push-zmq-event-collector:
	@echo Building ws-push-zmq-event-collector
	@echo Using Go Path :
	@echo $(GOPATH)
	mkdir ./dist/ws-push-zmq-event-collector
	go build -o ./dist/ws-push-zmq-event-collector/ws-to-push ws-push-zmq-event-collector/main/main.go
	cp ws-push-zmq-event-collector/main/config.json ./dist/ws-push-zmq-event-collector/config.json

clean:
	@echo Cleaning ...
	rm -fR ./dist

test: 
	@echo Testing ...
	@echo Coming ... 