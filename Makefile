# Go API Sample Makefile

# Variables
APP_NAME=go-api-server
DOCKER_IMAGE=go-api-server:latest
DOCKER_COMPOSE=docker-compose
DOCKER_COMPOSE_ENV=--env-file otel/.env

# Build binary
build:
	@echo "Building $(APP_NAME)..."
	go build -o bin/$(APP_NAME) .

# Build binary with vendor
build-vendor:
	@echo "Building $(APP_NAME) with vendor..."
	go build -mod=vendor -o bin/$(APP_NAME) .

# Run locally (requires PostgreSQL running)
run:
	@echo "Running $(APP_NAME) locally..."
	go run .

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Create vendor directory
vendor:
	@echo "Creating vendor directory..."
	go mod vendor

# Update vendor directory
vendor-update: deps vendor
	@echo "Vendor directory updated"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	rm -rf bin/
	rm -f coverage.out

# Clean vendor directory
clean-vendor:
	@echo "Cleaning vendor directory..."
	rm -rf vendor/

# Docker commands
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

docker-run: docker-build
	@echo "Running Docker container..."
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE)

# Docker Compose commands
docker-up:
	@echo "Starting services with Docker Compose..."
	$(DOCKER_COMPOSE) -f ./docker-compose-all.yml $(DOCKER_COMPOSE_ENV) up -d --no-deps --build

docker-down:
	@echo "Stopping services..."
	$(DOCKER_COMPOSE) -f ./docker-compose-all.yml down

docker-logs:
	@echo "Showing logs..."
	$(DOCKER_COMPOSE) -f ./docker-compose-all.yml logs -f

docker-restart:
	@echo "Restarting services..."
	$(DOCKER_COMPOSE) -f ./docker-compose-all.yml restart

# Development helpers
dev-up: docker-up
	@echo "Development environment is up!"
	@echo "API available at: http://localhost:8080"
	@echo "Health check: http://localhost:8080/health"
	@echo "You can check API by running: make test-api"

dev-down: docker-down
	@echo "Development environment stopped."

# Database helpers
db-reset:
	@echo "Resetting database..."
	$(DOCKER_COMPOSE) down postgres
	docker volume rm go-api-sample_postgres_data || true
	$(DOCKER_COMPOSE) up -d postgres

# API testing helpers
test-api:
	@echo "Testing API endpoints..."
	@echo "Health check:"
	curl -s http://localhost:8080/health | jq .
	@echo "\nAPI info:"
	curl -s http://localhost:8080/ | jq .

test-fetch-json:
	@echo "Testing JSON fetch with JSONPlaceholder API..."
	curl -s "http://localhost:8080/api/v1/fetch-json?link=https://jsonplaceholder.typicode.com/posts/1" | jq .

test-get-records:
	@echo "Getting stored records..."
	curl -s http://localhost:8080/api/v1/records | jq .

test-github-meta:
	@echo "Testing with GitHub Meta API..."
	./scripts/test-github-api.sh

curl-github-meta:
	@echo "Simple curl test with GitHub Meta API..."
	./scripts/curl-github-meta.sh

# Help
help:
	@echo "Available commands:"
	@echo "  build          - Build the application binary"
	@echo "  build-vendor   - Build with vendor mode"
	@echo "  run            - Run the application locally"
	@echo "  deps           - Install and tidy dependencies"
	@echo "  vendor         - Create vendor directory"
	@echo "  vendor-update  - Update vendor directory"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  clean          - Clean build artifacts"
	@echo "  clean-vendor   - Clean vendor directory"
	@echo ""
	@echo "Docker commands:"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  docker-up      - Start services with Docker Compose"
	@echo "  docker-down    - Stop Docker Compose services"
	@echo "  docker-logs    - Show Docker Compose logs"
	@echo "  docker-restart - Restart Docker Compose services"
	@echo ""
	@echo "Development:"
	@echo "  dev-up         - Start development environment"
	@echo "  dev-down       - Stop development environment"
	@echo "  db-reset       - Reset database"
	@echo ""
	@echo "Testing:"
	@echo "  test-api       - Test basic API endpoints"
	@echo "  test-fetch-json - Test JSON fetch endpoint"
	@echo "  test-get-records - Test get records endpoint"
	@echo "  test-github-meta - Test with GitHub Meta API (comprehensive)"
	@echo "  curl-github-meta - Simple curl test with GitHub Meta API"

.PHONY: build build-vendor run deps vendor vendor-update test test-coverage clean clean-vendor docker-build docker-run docker-up docker-down docker-logs docker-restart dev-up dev-down db-reset test-api test-fetch-json test-get-records test-github-meta curl-github-meta help