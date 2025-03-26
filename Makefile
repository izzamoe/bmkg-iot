# Makefile for BE_BMKG project

# Variables
BINARY_NAME=bmkg
BUILD_DIR=build
CMD_DIR=cmd/bmkg

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-s -w"

.PHONY: all build clean test run tidy vet fmt

all: clean build

build:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)

clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

test:
	$(GOTEST) -v ./...

run:
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	./$(BUILD_DIR)/$(BINARY_NAME)

tidy:
	$(GOMOD) tidy

vet:
	$(GOCMD) vet ./...

fmt:
	$(GOCMD) fmt ./...

# Dependencies management
deps:
	$(GOGET) -u ./...

# ini pocketbase jadi ada argument serve
serve:
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	./$(BUILD_DIR)/$(BINARY_NAME) serve
