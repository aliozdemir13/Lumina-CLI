# Build the binary
build:
	go build -o lumina main.go

# Run all tests
test:
	go test -v ./internal/

# Run tests with coverage report
cover:
	go test -cover ./internal/

# Remove build artifacts
clean:
	rm -f lumina

# Build and run
run: build
	./lumina --help