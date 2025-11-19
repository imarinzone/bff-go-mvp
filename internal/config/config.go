package config

import (
	"os"
)

// Config holds application configuration
type Config struct {
	Temporal TemporalConfig
	GRPC     GRPCConfig
	API      APIConfig
}

// TemporalConfig holds Temporal client configuration
type TemporalConfig struct {
	Host      string
	Namespace string
	TaskQueue string
}

// GRPCConfig holds gRPC client configuration
type GRPCConfig struct {
	ServiceAddress string
}

// APIConfig holds API server configuration
type APIConfig struct {
	Port string
}

// Load loads configuration from environment variables with defaults
func Load() *Config {
	return &Config{
		Temporal: TemporalConfig{
			Host:      getEnv("TEMPORAL_HOST", "localhost:7233"),
			Namespace: getEnv("TEMPORAL_NAMESPACE", "default"),
			TaskQueue: getEnv("TEMPORAL_TASK_QUEUE", "DISCOVERY_TASK_QUEUE"),
		},
		GRPC: GRPCConfig{
			ServiceAddress: getEnv("GRPC_SERVICE_ADDRESS", "localhost:50051"),
		},
		API: APIConfig{
			Port: getEnv("API_PORT", "8080"),
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

