package config

import (
	"os"
)

// Config holds application configuration
type Config struct {
	GRPC GRPCConfig
	API  APIConfig
}

// GRPCConfig holds gRPC client configuration
type GRPCConfig struct {
	DiscoverServiceAddress string
}

// APIConfig holds API server configuration
type APIConfig struct {
	Port string
}

// Load loads configuration from environment variables with defaults
func Load() *Config {
	return &Config{
		GRPC: GRPCConfig{
			DiscoverServiceAddress: getEnv("DISCOVER_SERVICE_ADDRESS", "localhost:8081"),
		},
		API: APIConfig{
			Port: getEnv("API_PORT", "8000"),
		},
	}
}

// getEnv gets environment variable or returns default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
