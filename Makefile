# Complaint Portal API Makefile

.PHONY: run build test clean demo help

# Default target
help:
	@echo "Available commands:"
	@echo "  run     - Run the application"
	@echo "  build   - Build the application"
	@echo "  test    - Run tests"
	@echo "  demo    - Run client demo"
	@echo "  clean   - Clean build artifacts"
	@echo "  help    - Show this help message"

# Run the application
run:
	@echo "Starting Complaint Portal API..."
	go run main.go

# Build the application
build:
	@echo "Building Complaint Portal API..."
	go build -o complaint-portal.exe main.go
	@echo "Build completed: complaint-portal.exe"

# Run tests
test:
	@echo "Running tests..."
	@echo "Starting server in background for testing..."
	@start /B go run main.go
	@timeout /t 3 /nobreak > nul
	go test -v
	@echo "Tests completed"

# Run client demo
demo:
	@echo "Starting server in background for demo..."
	@start /B go run main.go
	@timeout /t 3 /nobreak > nul
	@echo "Running client demo..."
	go run client_demo.go

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@if exist complaint-portal.exe del complaint-portal.exe
	@echo "Clean completed"

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Code formatted"

# Check for Go modules
mod:
	@echo "Tidying Go modules..."
	go mod tidy
	@echo "Modules tidied"