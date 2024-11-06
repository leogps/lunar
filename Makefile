#/*
 # * Copyright (c) 2024, Paul Gundarapu.
 # *
 # * Permission is hereby granted, free of charge, to any person obtaining a copy
 # * of this software and associated documentation files (the "Software"), to deal
 # * in the Software without restriction, including without limitation the rights
 # * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 # * copies of the Software, and to permit persons to whom the Software is
 # * furnished to do so, subject to the following conditions:
 # *
 # * The above copyright notice and this permission notice shall be included in
 # * all copies or substantial portions of the Software.
 # *
 # * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 # * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 # * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 # * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 # * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 # * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 # * THE SOFTWARE.
 # *
 # */
 
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

tidy:
	go mod tidy

test:
	go test -v ./...

# Build binaries for all platforms
build-all: clean | tidy test
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