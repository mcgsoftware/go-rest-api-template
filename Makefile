.PHONY: build run test clean docker-build docker-up docker-down swagger sqlc generate help

# Binary name
BINARY_NAME=server
BUILD_DIR=bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GOGET=$(GOCMD) get

# Default target
.DEFAULT_GOAL := help

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

## build: Build the binary
build:
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server

## run: Run the server locally
run: build
	./$(BUILD_DIR)/$(BINARY_NAME) serve

## test: Run tests
test:
	$(GOTEST) -v -race ./...

## test-coverage: Run tests with coverage
test-coverage:
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

## clean: Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

## deps: Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

## sqlc: Generate sqlc code
sqlc:
	sqlc generate -f db/sqlc.yaml

## swagger: Generate Swagger docs
swagger:
	swag init -g cmd/server/main.go -o docs

## generate: Run all code generation (sqlc + swagger)
generate: sqlc swagger

## docker-build: Build Docker image
docker-build:
	docker build -t go-rest-api-template:latest .

## docker-up: Start services with docker-compose
docker-up:
	docker-compose up -d

## docker-down: Stop services
docker-down:
	docker-compose down

## docker-logs: View docker-compose logs
docker-logs:
	docker-compose logs -f

## docker-clean: Stop services and remove volumes
docker-clean:
	docker-compose down -v

## lint: Run linter
lint:
	golangci-lint run ./...

## fmt: Format code
fmt:
	$(GOCMD) fmt ./...
