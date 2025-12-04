package config_test

import (
	"os"
	"testing"

	"bff-go-mvp/internal/config"
)

func TestLoad(t *testing.T) {
	// Save original env values
	originalGRPCAddr := os.Getenv("GRPC_SERVICE_ADDRESS")
	originalPort := os.Getenv("API_PORT")

	// Clean up after test
	defer func() {
		if originalGRPCAddr != "" {
			os.Setenv("GRPC_SERVICE_ADDRESS", originalGRPCAddr)
		} else {
			os.Unsetenv("GRPC_SERVICE_ADDRESS")
		}
		if originalPort != "" {
			os.Setenv("API_PORT", originalPort)
		} else {
			os.Unsetenv("API_PORT")
		}
	}()

	// Test with defaults
	cfg := config.Load()
	if cfg.GRPC.ServiceAddress != "localhost:50051" {
		t.Errorf("Expected default gRPC address 'localhost:50051', got '%s'", cfg.GRPC.ServiceAddress)
	}
	if cfg.API.Port != "8080" {
		t.Errorf("Expected default API port '8080', got '%s'", cfg.API.Port)
	}

	// Test with environment variables
	os.Setenv("GRPC_SERVICE_ADDRESS", "custom-grpc:50052")
	os.Setenv("API_PORT", "9090")

	cfg = config.Load()
	if cfg.GRPC.ServiceAddress != "custom-grpc:50052" {
		t.Errorf("Expected gRPC address 'custom-grpc:50052', got '%s'", cfg.GRPC.ServiceAddress)
	}
	if cfg.API.Port != "9090" {
		t.Errorf("Expected API port '9090', got '%s'", cfg.API.Port)
	}
}
