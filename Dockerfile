# Multi-platform build arguments (BuildKit feature)
ARG TARGETPLATFORM
ARG BUILDPLATFORM

# Build stage
FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM

WORKDIR /app

# Set GOTOOLCHAIN to auto to allow Go to download required version if needed
ENV GOTOOLCHAIN=auto

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN GOTOOLCHAIN=auto go mod download

# Copy source code
COPY . .

# Build the API server
# Parse TARGETPLATFORM to set GOOS and GOARCH
RUN TARGETOS=$(echo ${TARGETPLATFORM} | cut -d '/' -f1) && \
    TARGETARCH=$(echo ${TARGETPLATFORM} | cut -d '/' -f2) && \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} GOTOOLCHAIN=auto \
    go build -a -installsuffix cgo -o bin/api ./cmd/api

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/bin/api .

# Expose port
EXPOSE 8080

# Run the API server
CMD ["./api"]
