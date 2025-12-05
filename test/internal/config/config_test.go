package config_test

import (
	"os"
	"testing"

	"bff-go-mvp/internal/config"
)

func TestLoad(t *testing.T) {
	// Save original env values
	originalDiscoverAddr := os.Getenv("DISCOVER_SERVICE_ADDRESS")
	originalPort := os.Getenv("API_PORT")

	// Clean up after test
	defer func() {
		if originalDiscoverAddr != "" {
			os.Setenv("DISCOVER_SERVICE_ADDRESS", originalDiscoverAddr)
		} else {
			os.Unsetenv("DISCOVER_SERVICE_ADDRESS")
		}
		if originalPort != "" {
			os.Setenv("API_PORT", originalPort)
		} else {
			os.Unsetenv("API_PORT")
		}
	}()

	// Test with defaults
	cfg := config.Load()
	if cfg.GRPC.DiscoverServiceAddress != "localhost:8081" {
		t.Errorf("Expected default discover service address 'localhost:8081', got '%s'", cfg.GRPC.DiscoverServiceAddress)
	}
	if cfg.API.Port != "8000" {
		t.Errorf("Expected default API port '8000', got '%s'", cfg.API.Port)
	}

	// Test with environment variables
	os.Setenv("DISCOVER_SERVICE_ADDRESS", "custom-discover:8082")
	os.Setenv("API_PORT", "9090")

	cfg = config.Load()
	if cfg.GRPC.DiscoverServiceAddress != "custom-discover:8082" {
		t.Errorf("Expected discover service address 'custom-discover:8082', got '%s'", cfg.GRPC.DiscoverServiceAddress)
	}
	if cfg.API.Port != "9090" {
		t.Errorf("Expected API port '9090', got '%s'", cfg.API.Port)
	}
}
