# Show Service Makefile

.PHONY: help tidy fmt vet build test clean start-dynamo start stop restart logs

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

tidy: ## Clean up go.mod and go.sum
	go mod tidy

fmt: ## Format Go code
	go fmt ./...

vet: ## Check for suspicious code constructs
	go vet ./...

build: ## Build the application
	go build -o show-service ./cmd/server

test: ## Run tests
	go test -cover -count=1 ./internal/...

integration-test: ## Run integration tests (requires Docker for DynamoDB and app)
	@echo "Starting services..."
	docker-compose up -d --build
	@echo "Running integration tests..."
	go test ./test/integration_tests || (echo "Tests failed, stopping services..."; docker-compose down; exit 1)
	@echo "Stopping services..."
	docker-compose down

clean: ## Clean build artifacts
	rm -f show-service

# Docker targets
start-dynamo: ## Start DynamoDB Local only
	docker-compose up dynamodb-local --build

start: ## Start both DynamoDB Local and the application
	docker-compose up --build

stop: ## Stop all services
	docker-compose down

restart: ## Restart all services
	docker-compose restart

logs: ## Show logs from all services
	docker-compose logs -f

logs-dynamo: ## Show logs from DynamoDB Local only
	docker-compose logs -f dynamodb-local

logs-app: ## Show logs from the application only
	docker-compose logs -f show-service

# Development workflow
dev: tidy fmt vet test build ## Run development checks (tidy, format, vet, build)

# Production build
build-prod: ## Build for production
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o show-service ./cmd/server

# Quick start for development
quick-start: ## Quick start: build and start everything
	@echo "Building application..."
	@make build
	@echo "Starting services..."
	@make start

# Cleanup everything
reset: ## Reset everything (stop services, clean, rebuild)
	@echo "Stopping services..."
	@make stop
	@echo "Cleaning build artifacts..."
	@make clean
	@echo "Rebuilding..."
	@make build
	@echo "Starting fresh..."
	@make start
