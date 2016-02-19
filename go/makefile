GOPATH := $(shell pwd)


all: clean dist-dir build-collector build-soundtouch build-soundtouch-simulator

go-deps:
	go get github.com/pebbe/zmq4
	go get golang.org/x/net/websocket
	go get github.com/gorilla/mux

dist-dir:
	mkdir ./dist

build-collector:
	@echo Building collector
	@echo Using Go Path :
	@echo $(GOPATH)
	mkdir ./dist/collector
	go build -o ./dist/collector/server src/github.com/redsofa/collector/main/main.go
	cp src/github.com/redsofa/collector/main/config.json ./dist/collector/config.json
	cp -r src/github.com/redsofa/collector/main/webroot ./dist/collector/

build-soundtouch:
	@echo Building soundtouch
	@echo Using Go Path :
	@echo $(GOPATH)
	mkdir ./dist/soundtouch
	go build -o ./dist/soundtouch/soundtouch src/github.com/redsofa/soundtouch/main/main.go
	cp src/github.com/redsofa/soundtouch/main/config.json ./dist/soundtouch/config.json

build-soundtouch-simulator:
	@echo Building soundtouch-simulator
	@echo Using Go Path :
	@echo $(GOPATH)
	mkdir ./dist/soundtouch-simulator
	go build -o ./dist/soundtouch-simulator/simulator src/github.com/redsofa/soundtouch-simulator/main/main.go
	cp src/github.com/redsofa/soundtouch-simulator/main/config.json ./dist/soundtouch-simulator/config.json

clean:
	@echo Cleaning ...
	rm -fR ./dist

test: 
	@echo Testing ...
	@echo Coming ... 