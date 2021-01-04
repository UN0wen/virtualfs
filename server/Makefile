# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=server

all: test build

.PHONY: build
build: 
		$(GOBUILD) -o ./bin/$(BINARY_NAME) -v

.PHONY: test
test: 
		$(GOTEST) -v ./...

.PHONY: clean
clean: 
		$(GOCLEAN)
		rm -f $(BINARY_NAME)

.PHONY: run
run: build
		./bin/$(BINARY_NAME)

.PHONY: get
get:
		$(GOMOD) vendor