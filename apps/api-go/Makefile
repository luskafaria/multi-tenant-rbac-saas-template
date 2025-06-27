# Makefile

# Load environment variables from app.env
include app.env

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=api-server

.PHONY: dev db-up db-seed swag build run clean test test-coverage test-integration test-unit test-race tidy deps lint fmt vet check

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) ./cmd/server

# Run the application
run:
	$(GOCMD) run ./cmd/server

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Run integration tests specifically
test-integration:
	$(GOTEST) -v ./internal/integration/...

# Run unit tests (excluding integration tests)
test-unit:
	$(GOTEST) -v $(shell $(GOCMD) list ./... | grep -v /integration)

# Run tests with race detection
test-race:
	$(GOTEST) -v -race ./...

# Tidy dependencies
tidy:
	$(GOMOD) tidy

# Download dependencies
deps:
	$(GOMOD) download

# Development with hot reload (requires air)
dev:
	@echo "Starting the application in development mode with Air..."
	@air

db-up:
	@echo "Starting the database container..."
	@docker-compose up -d pg-go

db-seed:
	@echo "Seeding the database with test data..."
	@$(GOCMD) run ./cmd/seed

swag:
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/server/main.go

# Install development tools
install-tools:
	$(GOGET) github.com/swaggo/swag/cmd/swag@latest
	$(GOGET) github.com/cosmtrek/air@latest

# Format code
fmt:
	$(GOCMD) fmt ./...

# Vet code
vet:
	$(GOCMD) vet ./...

# Run all checks before commit
check: fmt vet test
