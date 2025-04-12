# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

LDFLAGS=
BUILD_DIR=./bin
BINARY_NAME?=server
PACKAGE?=server
MAIN_PACKAGE=./cmd/$(PACKAGE)

.PHONY: all build clean test run mod-tidy help install dev

# Default target
all: build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@$(GOBUILD) $(LDFLAGS) -buildvcs=false -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	@$(GOCLEAN)
	@rm -rf $(BUILD_DIR)

# Run tests
test:
	@echo "Running tests..."
	@$(GOTEST) -v ./...

# Build and run
run: build
	@echo "Running $(BINARY_NAME)..."
	@$(BUILD_DIR)/$(BINARY_NAME) $(ARGS)

# Clean up dependencies
mod-tidy:
	@echo "Tidying Go modules..."
	@$(GOCMD) mod tidy


# Cross-compile for multiple platforms
cross-build:
	@echo "Cross-compiling for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PACKAGE)
	@GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PACKAGE)
	@GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PACKAGE)

# Format code
fmt:
	@echo "Formatting code..."
	@$(GOCMD) fmt ./...

# Vet code for potential issues
vet:
	@echo "Vetting code..."
	@$(GOCMD) vet ./...

# Development workflow: format, vet, build
dev: fmt vet build