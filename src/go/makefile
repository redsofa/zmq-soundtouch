GOPATH := $(shell pwd)


all: dist-collector dist-soundtouch

go-deps:
	@echo Using Go Path :
	@echo $(GOPATH)	
	go get github.com/pebbe/zmq4
	go get golang.org/x/net/websocket
	go get github.com/gorilla/mux

dist-dir:
	mkdir ./dist

godoc:
	godoc -http=:6060
	$ (cd src/github.com/redsofa/; godoc -http=:6060)

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

clean:
	@echo Cleaning ...
	rm -fR ./dist

test: 
	@echo Testing ...
	go test ./src/github.com/redsofa/collector/messaging/. -v