# Go Arcade — Terminal Arcade Game Suite Makefile
# Provides targets for building, testing, vetting, and cross-compiling release binaries.

export DEVELOPER_DIR ?= /Library/Developer/CommandLineTools

BINARY_NAME=arcade
CMD_DIR=./cmd/arcade
DIST_DIR=dist
GO=go
VERSION ?= v1.0.0
FLAGS=-buildvcs=false -ldflags "-s -w -X main.version=$(VERSION)"

.PHONY: all build test vet cross-compile clean

all: build

build:
	@echo "Building native binary $(BINARY_NAME)..."
	$(GO) build $(FLAGS) -o $(BINARY_NAME) $(CMD_DIR)

test:
	@echo "Running all unit tests headlessly..."
	$(GO) test -buildvcs=false -v ./...

vet:
	@echo "Running static analysis..."
	$(GO) vet ./...

cross-compile: clean
	@echo "Building release binaries across platforms into $(DIST_DIR)/..."
	@mkdir -p $(DIST_DIR)
	GOOS=linux   GOARCH=amd64 $(GO) build $(FLAGS) -o $(DIST_DIR)/$(BINARY_NAME)_linux_amd64       $(CMD_DIR)
	GOOS=linux   GOARCH=arm64 $(GO) build $(FLAGS) -o $(DIST_DIR)/$(BINARY_NAME)_linux_arm64       $(CMD_DIR)
	GOOS=darwin  GOARCH=amd64 $(GO) build $(FLAGS) -o $(DIST_DIR)/$(BINARY_NAME)_darwin_amd64      $(CMD_DIR)
	GOOS=darwin  GOARCH=arm64 $(GO) build $(FLAGS) -o $(DIST_DIR)/$(BINARY_NAME)_darwin_arm64      $(CMD_DIR)
	GOOS=windows GOARCH=amd64 $(GO) build $(FLAGS) -o $(DIST_DIR)/$(BINARY_NAME)_windows_amd64.exe  $(CMD_DIR)
	@echo "Generating SHA256 checksums..."
	@cd $(DIST_DIR) && (command -v sha256sum >/dev/null 2>&1 && sha256sum $(BINARY_NAME)_* > checksums.txt || shasum -a 256 $(BINARY_NAME)_* > checksums.txt)
	@echo "Cross-compilation complete:"
	@ls -la $(DIST_DIR)
	@echo "Checksums:"
	@cat $(DIST_DIR)/checksums.txt

clean:
	@echo "Cleaning artifacts..."
	@rm -rf $(BINARY_NAME) $(BINARY_NAME).exe $(DIST_DIR) arcade_export.json arcade_export.csv

PREFIX ?= /usr/local
INSTALL_BIN ?= $(PREFIX)/bin

install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_BIN)..."
	@mkdir -p $(INSTALL_BIN)
	@cp $(BINARY_NAME) $(INSTALL_BIN)/$(BINARY_NAME)
	@chmod 755 $(INSTALL_BIN)/$(BINARY_NAME)
	@echo "Installation complete! Run 'arcade' to play."

uninstall:
	@echo "Removing $(BINARY_NAME) from $(INSTALL_BIN)..."
	@rm -f $(INSTALL_BIN)/$(BINARY_NAME)
	@echo "Uninstalled."
