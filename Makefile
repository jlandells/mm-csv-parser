# Variables
APP_NAME=mm-csv-parse

# Build for Linux
build-linux-amd64:
	@echo "Building for Linux AMD64..."
	GOOS=linux GOARCH=amd64 go build -o $(APP_NAME)_linux_amd64

build-linux-arm64:
	@echo "Building for Linux ARM64..."
	GOOS=linux GOARCH=arm64 go build -o $(APP_NAME)_linux_arm64

# Build for macOS
build-macos-apple:
	@echo "Building for macOS (Apple Silicon)..."
	GOOS=darwin GOARCH=arm64 go build -o $(APP_NAME)_macos_apple

build-macos-intel:
	@echo "Building for macOS (intel)..."
	GOOS=darwin GOARCH=amd64 go build -o $(APP_NAME)_macos_intel

# Build for Windows
build-windows:
	@echo "Building for Windows..."
	GOOS=darwin GOARCH=amd64 go build -o $(APP_NAME)_windows.exe

build-all: build-linux-amd64 build-linux-arm64 build-macos-apple build-macos-intel build-windows

.PHONY: build-linux-amd64 build-linux-arm64 build-macos-apple build-macos-intel build-windows build-all

# Clean Up
clean:
	@echo "Cleaning up..."
	rm -f $(APP_NAME)_*
