# This asks Git for the current tag. If no tag exists, it uses the commit hash.
VERSION := $(shell git describe --tags --always)

# Build the binary
build:
	go build -ldflags="-X github.com/aliozdemir13/Lumina/cmd.Version=$(VERSION)" -o lumina main.go

# Run all tests
test:
	go test -v ./...

# Run tests with coverage report
cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

# Remove build artifacts
clean:
	rm -f lumina
	rm -f coverage.out

# Build and run
run: build
	./lumina --help

# to prevent conflicts with files named 'test' or 'build'
.PHONY: build test cover clean run