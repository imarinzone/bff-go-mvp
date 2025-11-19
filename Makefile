.PHONY: build run-api run-worker test clean generate docker-build docker-up docker-down docker-logs docker-clean env

# Build all binaries
build:
	@echo "Building..."
	@go build -o bin/api cmd/api/main.go
	@go build -o bin/worker cmd/worker/main.go

# Run API server
run-api:
	@echo "Starting API server..."
	@go run cmd/api/main.go

# Run Temporal worker
run-worker:
	@echo "Starting Temporal worker..."
	@go run cmd/worker/main.go

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Generate protobuf code
generate:
	@echo "Generating protobuf code..."
	@mkdir -p proto/discovery/gen
	@protoc --go_out=proto/discovery/gen \
		--go_opt=paths=source_relative \
		--go-grpc_out=proto/discovery/gen \
		--go-grpc_opt=paths=source_relative \
		proto/discovery/discovery.proto

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf proto/discovery/gen

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Create .env file from .env.example
env:
	@if [ ! -f .env ]; then \
		echo "Creating .env file from .env.example..."; \
		cp .env.example .env; \
		echo ".env file created successfully!"; \
		echo "Please update .env with your configuration values."; \
	else \
		echo ".env file already exists. Skipping..."; \
	fi

# Docker commands
docker-build:
	@echo "Building Docker images..."
	@docker build -f Dockerfile.api -t bff-api:latest .
	@docker build -f Dockerfile.worker -t bff-worker:latest .

docker-up:
	@echo "Starting Docker Compose services..."
	@docker-compose up -d

docker-down:
	@echo "Stopping Docker Compose services..."
	@docker-compose down

docker-logs:
	@docker-compose logs -f

docker-clean:
	@echo "Cleaning Docker resources..."
	@docker-compose down -v
	@docker rmi bff-api:latest bff-worker:latest 2>/dev/null || true


