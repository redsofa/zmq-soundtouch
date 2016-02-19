GOPATH := $(shell pwd)


all: dist-collector dist-soundtouch dist-soundtouch-simulator

go-deps:
	go get github.com/pebbe/zmq4
	go get golang.org/x/net/websocket
	go get github.com/gorilla/mux

dist-dir:
	mkdir ./dist

run-collector:
	@echo Running collector
	@echo Using Go Path :
	@echo $(GOPATH)
	$ (cd src/github.com/redsofa/collector/main; go run main.go)	

build-collector: clean dist-dir 
	@echo Building collector
	@echo Using Go Path :
	@echo $(GOPATH)
	mkdir ./dist/collector
	go build -o ./dist/collector/collector src/github.com/redsofa/collector/main/main.go

dist-collector: build-collector
	@echo Distributing collector
	cp src/github.com/redsofa/collector/main/config.json ./dist/collector/config.json
	cp -r src/github.com/redsofa/collector/main/webroot ./dist/collector/

build-soundtouch: clean dist-dir
	@echo Building soundtouch
	@echo Using Go Path :
	@echo $(GOPATH)
	mkdir ./dist/soundtouch
	go build -o ./dist/soundtouch/soundtouch src/github.com/redsofa/soundtouch/main/main.go

dist-soundtouch: build-soundtouch
	@echo Distributing soundtouch
	cp src/github.com/redsofa/soundtouch/main/config.json ./dist/soundtouch/config.json


build-soundtouch-simulator: clean dist-dir
	@echo Building soundtouch-simulator
	@echo Using Go Path :
	@echo $(GOPATH)
	mkdir ./dist/soundtouch-simulator
	go build -o ./dist/soundtouch-simulator/soundtouch-simulator src/github.com/redsofa/soundtouch-simulator/main/main.go

dist-soundtouch-simulator: build-soundtouch-simulator
	@echo Distributing soundtouch-simulator
	cp src/github.com/redsofa/soundtouch-simulator/main/config.json ./dist/soundtouch-simulator/config.json

clean:
	@echo Cleaning ...
	rm -fR ./dist

test: 
	@echo Testing ...
	@echo Coming ... 