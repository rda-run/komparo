.PHONY: build build-linux build-windows build-darwin build-all clean test lint run help

APP_NAME = komparo
DIST_DIR = dist

# Read version from .version file, default to 'dev' if missing
VERSION ?= $(shell cat .version 2>/dev/null || echo "dev")
LDFLAGS = -s -w -X 'github.com/rda-run/komparo/cmd.Version=$(VERSION)'

# Aggressive optimizations for production:
# - CGO_ENABLED=0: Statically linked binaries
# - -trimpath: Removes absolute file system paths from the compiled executable
# - -ldflags="-s -w": Strips symbol table and DWARF debugging information to reduce binary size

build: build-linux ## Build the highly optimized binary for Linux (production default)

build-linux: ## Build highly optimized binary for Linux
	@echo "Building $(APP_NAME) for Linux..."
	@mkdir -p $(DIST_DIR)
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(APP_NAME)-linux-amd64 .

build-windows: ## Build highly optimized binary for Windows
	@echo "Building $(APP_NAME) for Windows..."
	@mkdir -p $(DIST_DIR)
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(APP_NAME)-windows-amd64.exe .

build-darwin: ## Build highly optimized binary for macOS (Intel & ARM64/Apple Silicon)
	@echo "Building $(APP_NAME) for macOS (Intel)..."
	@mkdir -p $(DIST_DIR)
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(APP_NAME)-darwin-amd64 .
	@echo "Building $(APP_NAME) for macOS (ARM64)..."
	@CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(APP_NAME)-darwin-arm64 .

build-all: build-linux build-windows build-darwin ## Build optimized binaries for all operating systems

clean: ## Remove the dist folder and previous builds
	@echo "Cleaning up..."
	@rm -rf $(DIST_DIR)

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

lint: ## Run go vet and go fmt
	@echo "Running linter..."
	@go vet ./...
	@go fmt ./...

run: build ## Build and run the linux binary
	@$(DIST_DIR)/$(APP_NAME)-linux-amd64

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
