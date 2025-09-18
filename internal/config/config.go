package config

import (
	"fmt"
	"os"
	"strings"
)

// Config holds all configuration for the Weblate MCP server
type Config struct {
	WeblateAPIURL   string
	WeblateAPIToken string
	LogLevel        string
	MCPServerName   string
	MCPServerVersion string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	config := &Config{
		WeblateAPIURL:    os.Getenv("WEBLATE_API_URL"),
		WeblateAPIToken:  os.Getenv("WEBLATE_API_TOKEN"),
		LogLevel:         getEnvWithDefault("LOG_LEVEL", "info"),
		MCPServerName:    getEnvWithDefault("MCP_SERVER_NAME", "weblate-mcp-server"),
		MCPServerVersion: getEnvWithDefault("MCP_SERVER_VERSION", "1.3.0"),
	}

	// Validate required configuration
	if config.WeblateAPIURL == "" {
		return nil, fmt.Errorf("WEBLATE_API_URL is required")
	}
	if config.WeblateAPIToken == "" {
		return nil, fmt.Errorf("WEBLATE_API_TOKEN is required")
	}

	// Ensure the API URL ends with /api
	if !strings.HasSuffix(config.WeblateAPIURL, "/api") {
		config.WeblateAPIURL += "/api"
	}

	return config, nil
}

// getEnvWithDefault returns the value of the environment variable or a default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}