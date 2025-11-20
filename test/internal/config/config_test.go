package config_test

import (
	"os"
	"testing"

	"bff-go-mvp/internal/config"
)

func TestLoad(t *testing.T) {
	// Save original env values
	originalHost := os.Getenv("TEMPORAL_HOST")
	originalNamespace := os.Getenv("TEMPORAL_NAMESPACE")
	originalTaskQueue := os.Getenv("TEMPORAL_TASK_QUEUE")
	originalGRPCAddr := os.Getenv("GRPC_SERVICE_ADDRESS")
	originalPort := os.Getenv("API_PORT")

	// Clean up after test
	defer func() {
		if originalHost != "" {
			os.Setenv("TEMPORAL_HOST", originalHost)
		} else {
			os.Unsetenv("TEMPORAL_HOST")
		}
		if originalNamespace != "" {
			os.Setenv("TEMPORAL_NAMESPACE", originalNamespace)
		} else {
			os.Unsetenv("TEMPORAL_NAMESPACE")
		}
		if originalTaskQueue != "" {
			os.Setenv("TEMPORAL_TASK_QUEUE", originalTaskQueue)
		} else {
			os.Unsetenv("TEMPORAL_TASK_QUEUE")
		}
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
	if cfg.Temporal.Host != "localhost:7233" {
		t.Errorf("Expected default Temporal host 'localhost:7233', got '%s'", cfg.Temporal.Host)
	}
	if cfg.Temporal.Namespace != "default" {
		t.Errorf("Expected default namespace 'default', got '%s'", cfg.Temporal.Namespace)
	}
	if cfg.GRPC.ServiceAddress != "localhost:50051" {
		t.Errorf("Expected default gRPC address 'localhost:50051', got '%s'", cfg.GRPC.ServiceAddress)
	}
	if cfg.API.Port != "8080" {
		t.Errorf("Expected default API port '8080', got '%s'", cfg.API.Port)
	}

	// Test with environment variables
	os.Setenv("TEMPORAL_HOST", "custom-host:7234")
	os.Setenv("TEMPORAL_NAMESPACE", "custom-namespace")
	os.Setenv("TEMPORAL_TASK_QUEUE", "CUSTOM_QUEUE")
	os.Setenv("GRPC_SERVICE_ADDRESS", "custom-grpc:50052")
	os.Setenv("API_PORT", "9090")

	cfg = config.Load()
	if cfg.Temporal.Host != "custom-host:7234" {
		t.Errorf("Expected Temporal host 'custom-host:7234', got '%s'", cfg.Temporal.Host)
	}
	if cfg.Temporal.Namespace != "custom-namespace" {
		t.Errorf("Expected namespace 'custom-namespace', got '%s'", cfg.Temporal.Namespace)
	}
	if cfg.Temporal.TaskQueue != "CUSTOM_QUEUE" {
		t.Errorf("Expected task queue 'CUSTOM_QUEUE', got '%s'", cfg.Temporal.TaskQueue)
	}
	if cfg.GRPC.ServiceAddress != "custom-grpc:50052" {
		t.Errorf("Expected gRPC address 'custom-grpc:50052', got '%s'", cfg.GRPC.ServiceAddress)
	}
	if cfg.API.Port != "9090" {
		t.Errorf("Expected API port '9090', got '%s'", cfg.API.Port)
	}
}
