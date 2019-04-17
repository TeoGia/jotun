	# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=jotun
BINARY_UNIX=$(BINARY_NAME)_unix
OS=$(shell uname -s)


all: runtest build run
.PHONY: build
build:
	$(info Building for: $(OS))
	$(GOBUILD) -o ./bin/$(BINARY_NAME) -v ./cmd/jotun/... 
runtest: 
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -f ./bin/$(BINARY_NAME)
	rm -f ./bin/$(BINARY_UNIX)
run:
	$(GOBUILD) -o ./bin/$(BINARY_NAME) -v ./cmd/jotun/...
	./bin/$(BINARY_NAME)
runprd:
	$(info Building for: $(OS))
	$(GOBUILD) -o ./bin/$(BINARY_NAME) -v ./cmd/jotun/...
	tar cvzf ./release/jotun-$(shell ./bin/$(BINARY_NAME) -v).tar.gz ./LICENSE ./jotun.man ./bin/jotun
	
