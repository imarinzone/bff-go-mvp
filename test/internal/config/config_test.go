package config_test

import (
	"os"
	"testing"

	"bff-go-mvp/internal/config"
)

func TestLoad(t *testing.T) {
	// Save original env values
	originalServiceAddr := os.Getenv("GRPC_SERVICE_ADDRESS")
	originalPort := os.Getenv("API_PORT")

	// Clean up after test
	defer func() {
		if originalServiceAddr != "" {
			os.Setenv("GRPC_SERVICE_ADDRESS", originalServiceAddr)
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
	if cfg.GRPC.ServiceAddress != "localhost:9090" {
		t.Errorf("Expected default service address 'localhost:9090', got '%s'", cfg.GRPC.ServiceAddress)
	}
	if cfg.API.Port != "8000" {
		t.Errorf("Expected default API port '8000', got '%s'", cfg.API.Port)
	}

	// Test with environment variables
	os.Setenv("GRPC_SERVICE_ADDRESS", "custom-service:8082")
	os.Setenv("API_PORT", "9090")

	cfg = config.Load()
	if cfg.GRPC.ServiceAddress != "custom-service:8082" {
		t.Errorf("Expected service address 'custom-service:8082', got '%s'", cfg.GRPC.ServiceAddress)
	}
	if cfg.API.Port != "9090" {
		t.Errorf("Expected API port '9090', got '%s'", cfg.API.Port)
	}
}
