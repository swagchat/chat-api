GOCMD=$(GOROOT)/bin/go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=chat-api

all: test build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v
test:
	$(GOTEST) -covermode=count -coverprofile=coverage.out ./... && $(GOCMD) tool cover -html=coverage.out
test-func:
	$(GOTEST) -covermode=count -coverprofile=coverage.out ./... && $(GOCMD) tool cover -func=coverage.out
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
run:
	$(GOBUILD)
	./$(BINARY_NAME)
