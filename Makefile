.PHONY: build run-api run-api-dev test clean generate swagger swagger-clean docker-build docker-up docker-down docker-logs docker-clean env

# Build all binaries
build:
	@echo "Building..."
	@go build -o bin/api cmd/api/main.go

# Run API server
run-api:
	@echo "Starting API server..."
	@go run cmd/api/main.go

# Run API server in dev mode with auto-reload (requires air: go install github.com/air-verse/air@latest)
run-api-dev:
	@echo "Starting API server with air (auto-reload)..."
	@air -c .air.toml

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Generate protobuf code
generate:
	@echo "Generating protobuf code..."
	@mkdir -p proto/common/gen proto/search/gen
	@PATH=$$(go env GOPATH)/bin:$$PATH protoc --proto_path=proto/schemas \
		--go_out=proto/common/gen \
		--go_opt=paths=source_relative \
		proto/schemas/common/context.proto
	@PATH=$$(go env GOPATH)/bin:$$PATH protoc --proto_path=proto/schemas \
		--go_out=proto/search/gen \
		--go_opt=paths=source_relative \
		--go-grpc_out=proto/search/gen \
		--go-grpc_opt=paths=source_relative \
		proto/schemas/search/search.proto

# Generate Swagger docs using swaggo
swagger:
	@echo "Generating Swagger docs with swag..."
	@swag init -g cmd/api/main.go -o internal/docs

# Clean generated Swagger docs
swagger-clean:
	@echo "Cleaning generated Swagger docs..."
	@rm -rf internal/docs/*.go

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf proto/common/gen proto/search/gen

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
	@docker rmi bff-api:latest 2>/dev/null || true


