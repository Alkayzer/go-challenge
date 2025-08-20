.PHONY: help build run test clean docker-build docker-run docker-stop lint format

# Default target
help:
	@echo "Available commands:"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application locally"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run with Docker Compose"
	@echo "  docker-stop  - Stop Docker Compose services"
	@echo "  lint         - Run linter"
	@echo "  format       - Format code"

# Build the application
build:
	go build -o bin/events-service ./cmd/server

# Run the application locally
run:
	go run ./cmd/server



# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Build Docker image
docker-build:
	docker build -t events-service .

# Run with Docker Compose
docker-run:
	docker-compose up -d

# Stop Docker Compose services
docker-stop:
	docker-compose down

# Run linter
lint:
	golangci-lint run

# Format code
format:
	go fmt ./...
	go vet ./...

# Install dependencies
deps:
	go mod download
	go mod tidy

# Generate go.sum
sum:
	go mod tidy
	go mod verify
