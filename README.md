# BFF Go MVP

A Backend for Frontend (BFF) service built with Go that handles REST API requests and communicates with downstream microservices via gRPC using Temporal for workflow orchestration.

## Architecture

- **REST API**: Gorilla Mux router for handling HTTP requests
- **Logging**: Zap for structured, high-performance logging
- **Temporal**: Workflow orchestration for async processing
- **gRPC**: Communication with downstream microservices
- **Protobuf**: Service definitions and data contracts

## Project Structure

```
bff-go-mvp/
├── cmd/
│   ├── api/          # REST API server
│   └── worker/       # Temporal worker
├── internal/
│   ├── api/          # API handlers
│   ├── temporal/     # Temporal workflows & activities
│   ├── grpc/         # gRPC client
│   ├── config/       # Configuration
│   └── logger/       # Logger utilities
├── proto/            # Protobuf definitions
└── pkg/              # Shared packages
    └── models/       # Data models
```

## Prerequisites

### For Local Development
- Go 1.24 or higher
- Temporal server running (default: localhost:7233)
- Protocol Buffers compiler (`protoc`) - for generating code from `.proto` files
- Go protobuf plugins:
  - `protoc-gen-go` - generates Go code from protobuf
  - `protoc-gen-go-grpc` - generates gRPC Go code

### For Docker
- Docker 20.10 or higher
- Docker Compose 2.0 or higher

## Setup

1. Install dependencies:
```bash
go mod tidy
```

2. Create environment file:
```bash
# Create .env file from template
make env
# or manually
cp .env.example .env
```

Edit `.env` file with your configuration values if needed (defaults are provided).

3. Install Protocol Buffers tools (if not already installed):

**On macOS (using Homebrew):**
```bash
# Install protoc compiler
brew install protobuf

# Install Go protobuf plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# IMPORTANT: Ensure Go bin directory is in PATH
# Add this to your ~/.zshrc or ~/.bash_profile:
export PATH="$PATH:$(go env GOPATH)/bin"

# Or run it in your current shell:
export PATH="$PATH:$(go env GOPATH)/bin"

# Verify installation:
protoc --version
protoc-gen-go --version
protoc-gen-go-grpc --version
```

**On Linux (Ubuntu/Debian):**
```bash
# Install protoc compiler
sudo apt-get update
sudo apt-get install -y protobuf-compiler

# Install Go protobuf plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# IMPORTANT: Ensure Go bin directory is in PATH
# Add this to your ~/.bashrc or ~/.profile:
export PATH="$PATH:$(go env GOPATH)/bin"

# Or run it in your current shell:
export PATH="$PATH:$(go env GOPATH)/bin"

# Verify installation:
protoc --version
protoc-gen-go --version
protoc-gen-go-grpc --version
```

**On Windows:**
```bash
# Install protoc compiler
# Download from: https://github.com/protocolbuffers/protobuf/releases
# Extract and add to PATH

# Install Go protobuf plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Ensure Go bin directory is in PATH
# Add %USERPROFILE%\go\bin to your PATH
```

4. Generate protobuf code:
```bash
make generate
```

**Note:** The `make generate` command requires:
- `protoc` - Protocol Buffers compiler (must be in PATH)
- `protoc-gen-go` - Go protobuf plugin (must be in PATH)
- `protoc-gen-go-grpc` - Go gRPC plugin (must be in PATH)

**Troubleshooting:**
If you get "command not found" errors, ensure:
1. All tools are installed (see installation steps above)
2. Go bin directory is in your PATH: `export PATH="$PATH:$(go env GOPATH)/bin"`
3. Verify with: `which protoc protoc-gen-go protoc-gen-go-grpc`

5. Start Temporal server (if not already running):
```bash
# Using Docker Compose (includes Temporal, PostgreSQL, API, and Worker)
docker-compose up -d

# Or using Docker directly
docker run -p 7233:7233 temporalio/auto-setup:latest
```

## Running the Application

### Using Docker Compose (Recommended)

Start all services (Temporal, PostgreSQL, API, and Worker) with a single command:

```bash
docker-compose up -d
```

View logs:
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f api
docker-compose logs -f worker
docker-compose logs -f temporal
```

Stop all services:
```bash
docker-compose down
```

Stop and remove volumes:
```bash
docker-compose down -v
```

### Local Development

#### Start the Temporal Worker

In one terminal:
```bash
make run-worker
# or
go run cmd/worker/main.go
```

#### Start the API Server

In another terminal:
```bash
make run-api
# or
go run cmd/api/main.go
```

The API server will start on `http://localhost:8080`

### Docker Commands

Build individual images:
```bash
# Build API image
docker build -f Dockerfile.api -t bff-api:latest .

# Build Worker image
docker build -f Dockerfile.worker -t bff-worker:latest .
```

Run individual containers:
```bash
# Run API container
docker run -p 8080:8080 \
  -e TEMPORAL_HOST=host.docker.internal:7233 \
  -e TEMPORAL_NAMESPACE=default \
  -e TEMPORAL_TASK_QUEUE=DISCOVERY_TASK_QUEUE \
  bff-api:latest

# Run Worker container
docker run \
  -e TEMPORAL_HOST=host.docker.internal:7233 \
  -e TEMPORAL_NAMESPACE=default \
  -e TEMPORAL_TASK_QUEUE=DISCOVERY_TASK_QUEUE \
  bff-worker:latest
```

## API Endpoints

### POST /discovery

Discover services based on location and context.

**Request Body:**
```json
{
  "context": {
    "version": "1.0.0",
    "action": "on_discover",
    "domain": "mobility",
    "location": {
      "country": {
        "code": "IND"
      },
      "city": {
        "code": "std:080"
      }
    },
    "bap_id": "bap-123",
    "bap_uri": "https://bap.example.com",
    "bpp_id": "bpp-456",
    "bpp_uri": "https://bpp.example.com",
    "transaction_id": "txn-789",
    "message_id": "msg-001",
    "timestamp": "2024-01-01T00:00:00Z",
    "ttl": "PT30S"
  },
  "message": {
    "catalogs": [...]
  }
}
```

**Response:**
Returns the discovery response from the downstream gRPC service.

## Configuration

Configuration can be set via environment variables. The recommended approach is to use a `.env` file:

1. Create `.env` file from template:
   ```bash
   make env
   ```

2. Edit `.env` file with your configuration values.

### Environment Variables

- `ENV`: Environment mode - "development" or "dev" for dev logger, otherwise production (default: production)
- `TEMPORAL_HOST`: Temporal server host (default: localhost:7233)
- `TEMPORAL_NAMESPACE`: Temporal namespace (default: default)
- `TEMPORAL_TASK_QUEUE`: Temporal task queue name (default: DISCOVERY_TASK_QUEUE)
- `GRPC_SERVICE_ADDRESS`: gRPC service address (default: localhost:50051)
- `API_PORT`: API server port (default: 8080)

### Using .env File

When running locally, the application will automatically read from `.env` file if you use a tool like `godotenv` or export the variables:

```bash
# Export variables from .env file
export $(cat .env | xargs)
go run cmd/api/main.go
```

Docker Compose automatically loads variables from `.env` file, so no additional setup is needed for Docker.

## Development

### Build
```bash
# Build binaries locally
make build

# Build Docker images
make docker-build
```

### Run Tests
```bash
make test
```

### Generate Protobuf Code
```bash
make generate
```

### Docker Development

Rebuild and restart services:
```bash
docker-compose up -d --build
```

Access Temporal UI:
- Open `http://localhost:8088` in your browser

View service status:
```bash
docker-compose ps
```

## License

MIT


