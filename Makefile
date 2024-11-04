.PHONY: clean install build-all

APP_NAME := lunar
BIN_DIR := ./bin
SRC_FILE := main.go

# Target platforms and architectures
PLATFORMS := \
	linux/arm64 \
	windows/arm64 \
	darwin/arm64 \
	linux/amd64 \
	windows/amd64 \
	darwin/amd64

# Clean up old binaries
clean:
	@echo "Cleaning binaries..."
	@rm -rf $(BIN_DIR)

# Build binaries for all platforms
build-all: clean
	@mkdir -p $(BIN_DIR)
	@for platform in $(PLATFORMS); do \
		GOOS=$$(echo $$platform | cut -d'/' -f1); \
		GOARCH=$$(echo $$platform | cut -d'/' -f2); \
		OUTPUT_NAME=$(BIN_DIR)/$(APP_NAME)-$$GOOS-$$GOARCH; \
		if [ "$$GOOS" = "windows" ]; then OUTPUT_NAME=$$OUTPUT_NAME.exe; fi; \
		echo "Building $$OUTPUT_NAME..."; \
		GOOS=$$GOOS GOARCH=$$GOARCH go build -o $$OUTPUT_NAME $(SRC_FILE); \
	done

# Default install for the current OS/Arch
install:
	@echo "Building and installing $(APP_NAME) for the current platform..."
	@go build -o $(BIN_DIR)/$(APP_NAME) $(SRC_FILE)