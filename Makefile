.PHONY: all build test clean install

# Build settings
BINARY_NAME=remotesync
GO=go
GOFLAGS=-v -mod=vendor
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
	$(GO) test $(GOFLAGS) ./...

clean:
	rm -rf $(BIN_DIR)
	rm -rf $(DIST_DIR)
	rm -rf vendor/

install: build
	install -d $(DESTDIR)/usr/local/bin/
	install -m 755 $(BIN_DIR)/$(BINARY_NAME)-server $(DESTDIR)/usr/local/bin/
	install -m 755 $(BIN_DIR)/$(BINARY_NAME)-client $(DESTDIR)/usr/local/bin/