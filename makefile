GOPATH := $(HOME)/dev/git/go

all: clean dist-dir build-sub-zmq-event-collector-ws build-ws-push-zmq-event-collector build-push-zmq-event-collector-tester

dist-dir:
	mkdir ./dist

build-sub-zmq-event-collector-ws:
	@echo Building sub-zmq-event-collector-ws
	@echo Using Go Path :
	@echo $(GOPATH)
	mkdir ./dist/sub-zmq-event-collector-ws
	go build -o ./dist/sub-zmq-event-collector-ws/sub-to-ws sub-zmq-event-collector-ws/main/main.go
	cp sub-zmq-event-collector-ws/main/config.json ./dist/sub-zmq-event-collector-ws/config.json
	cp -r sub-zmq-event-collector-ws/main/www ./dist/sub-zmq-event-collector-ws/

build-ws-push-zmq-event-collector:
	@echo Building ws-push-zmq-event-collector
	@echo Using Go Path :
	@echo $(GOPATH)
	mkdir ./dist/ws-push-zmq-event-collector
	go build -o ./dist/ws-push-zmq-event-collector/ws-to-push ws-push-zmq-event-collector/main/main.go
	cp ws-push-zmq-event-collector/main/config.json ./dist/ws-push-zmq-event-collector/config.json

build-push-zmq-event-collector-tester:
	@echo Building push-zmq-event-collector-tester
	@echo Using Go Path :
	@echo $(GOPATH)
	mkdir ./dist/push-zmq-event-collector-tester
	go build -o ./dist/push-zmq-event-collector-tester/push-test push-zmq-event-collector-tester/main/main.go
	cp push-zmq-event-collector-tester/main/config.json ./dist/push-zmq-event-collector-tester/config.json

clean:
	@echo Cleaning ...
	rm -fR ./dist

test: 
	@echo Testing ...
	@echo Coming ... 