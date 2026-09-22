# Go Arcade — Terminal Arcade Game Suite Makefile
# Provides targets for building, testing, vetting, and cross-compiling release binaries.

BINARY_NAME=arcade
CMD_DIR=./cmd/arcade
DIST_DIR=dist
GO=go
FLAGS=-buildvcs=false

.PHONY: all build test vet cross-compile clean

all: build

build:
	@echo "Building native binary $(BINARY_NAME)..."
	$(GO) build $(FLAGS) -o $(BINARY_NAME) $(CMD_DIR)

test:
	@echo "Running all unit tests headlessly..."
	$(GO) test -v ./...

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
	@echo "Cross-compilation complete:"
	@ls -la $(DIST_DIR)

clean:
	@echo "Cleaning artifacts..."
	@rm -rf $(BINARY_NAME) $(BINARY_NAME).exe $(DIST_DIR) arcade_export.json arcade_export.csv
