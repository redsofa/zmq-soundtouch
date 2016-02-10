GOPATH := $(HOME)/dev/git/go

all: clean build-pull-zmq-event-collector-ws copy-configs

build-pull-zmq-event-collector-ws:
	@echo Building pull-zmq-event-collector-ws
	@echo Using Go Path :
	@echo $(GOPATH)


clean:
	@echo Cleaning ...


copy-configs:
	@echo Copying Configurations ...


test: 
	@echo Testing ... 