# Go parameters
GO := go
GOBUILD := $(GO) build
GOCLEAN := $(GO) clean
GOTEST := $(GO) test
GOGET := $(GO) get
BUF := buf

# Directories
CMD_DIR := cmd
PROTO_DIR := proto
GEN_DIR := gen/go
BIN_DIR := bin

# Binary names
CLIENT_BINARY := $(BIN_DIR)/client
SERVER_BINARY := $(BIN_DIR)/server

.PHONY: all build test clean generate get-deps run-client run-server buf-generate

all: build

build: build-client build-server

build-client:
	@echo "Building client..."
	@cd $(CMD_DIR)/client && $(GOBUILD) -o ../../$(CLIENT_BINARY)

build-server:
	@echo "Building server..."
	@cd $(CMD_DIR)/server && $(GOBUILD) -o ../../$(SERVER_BINARY)

test:
	@echo "Running tests..."
	@$(GOTEST) -v ./...

clean:
	@echo "Cleaning..."
	@$(GOCLEAN)
	@rm -f $(CLIENT_BINARY) $(SERVER_BINARY)

get-deps:
	@echo "Getting dependencies..."
	@$(GOGET) ./...

run-client: build-client
	@echo "Running client..."
	@$(CLIENT_BINARY)

run-server: build-server
	@echo "Running server..."
	@$(SERVER_BINARY)

buf-generate:
	@echo "Installing buf dependencies..."
	@$(BUF) mod update
	@echo "Running buf generate..."
	@$(BUF) generate
	@echo "Getting Go dependencies..."
	@$(GO) mod download
	@$(GO) mod vendor

tidy:
	go mod tidy
	go mod vendor
