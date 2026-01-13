package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database
	DatabaseURL string

	// API Server
	APIPort string
	APIHost string

	// Kubernetes
	KubernetesInCluster bool
	Kubeconfig          string

	// Registry
	RegistryURL string

	// Environment
	Environment string
	LogLevel    string
}

func Load() (*Config, error) {
	// Try to load .env file (ignore error if not found)
	_ = godotenv.Load()

	cfg := &Config{
		DatabaseURL:         getEnv("DATABASE_URL", ""),
		APIPort:             getEnv("API_PORT", "8080"),
		APIHost:             getEnv("API_HOST", "0.0.0.0"),
		KubernetesInCluster: getEnvBool("KUBERNETES_IN_CLUSTER", false),
		Kubeconfig:          getEnv("KUBECONFIG", ""),
		RegistryURL:         getEnv("REGISTRY_URL", "registry.superfly-system.svc.cluster.local:5000"),
		Environment:         getEnv("ENV", "development"),
		LogLevel:            getEnv("LOG_LEVEL", "info"),
	}

	// Validate required fields
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return defaultValue
		}
		return boolVal
	}
	return defaultValue
}
