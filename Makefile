.PHONY: build clean test coverage lint

BINARY_NAME=keytool
BUILD_DIR=bin

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/keytool

clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

test:
	go test -v ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run ./...
