.PHONY: all build test clean install

# Build settings
BINARY_NAME=remotesync
GO=go
GOFLAGS=-v
VERSION=1.0.0

# Output directories
BIN_DIR=bin
DIST_DIR=dist

all: clean build

build:
	mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS) -o $(BIN_DIR)/$(BINARY_NAME)-server ./cmd/server
	$(GO) build $(GOFLAGS) -o $(BIN_DIR)/$(BINARY_NAME)-client ./cmd/client

test:
	$(GO) test -v ./...
	$(GO) test -bench=. ./internal/performance

dist: build
	mkdir -p $(DIST_DIR)
	tar czf $(DIST_DIR)/$(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz -C $(BIN_DIR) .
	zip -j $(DIST_DIR)/$(BINARY_NAME)-$(VERSION)-linux-amd64.zip $(BIN_DIR)/*

clean:
	rm -rf $(BIN_DIR)
	rm -rf $(DIST_DIR)

install: build
	install -d $(DESTDIR)/usr/local/bin/
	install -m 755