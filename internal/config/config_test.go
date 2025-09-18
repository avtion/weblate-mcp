package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Set environment variables for test
	os.Setenv("WEBLATE_API_URL", "https://test.weblate.com")
	os.Setenv("WEBLATE_API_TOKEN", "test-token")

	config, err := Load()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if config.WeblateAPIURL != "https://test.weblate.com/api" {
		t.Errorf("Expected WeblateAPIURL to be 'https://test.weblate.com/api', got '%s'", config.WeblateAPIURL)
	}

	if config.WeblateAPIToken != "test-token" {
		t.Errorf("Expected WeblateAPIToken to be 'test-token', got '%s'", config.WeblateAPIToken)
	}

	// Clean up
	os.Unsetenv("WEBLATE_API_URL")
	os.Unsetenv("WEBLATE_API_TOKEN")
}

func TestLoadMissingRequired(t *testing.T) {
	// Ensure environment variables are not set
	os.Unsetenv("WEBLATE_API_URL")
	os.Unsetenv("WEBLATE_API_TOKEN")

	_, err := Load()
	if err == nil {
		t.Fatal("Expected error for missing required configuration, got none")
	}
}

func TestAPIURLNormalization(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"https://test.weblate.com", "https://test.weblate.com/api"},
		{"https://test.weblate.com/", "https://test.weblate.com/api"},
		{"https://test.weblate.com/api", "https://test.weblate.com/api"},
		{"https://test.weblate.com/api/", "https://test.weblate.com/api"},
	}

	for _, test := range tests {
		os.Setenv("WEBLATE_API_URL", test.input)
		os.Setenv("WEBLATE_API_TOKEN", "test-token")

		config, err := Load()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if config.WeblateAPIURL != test.expected {
			t.Errorf("For input '%s', expected '%s', got '%s'", test.input, test.expected, config.WeblateAPIURL)
		}

		os.Unsetenv("WEBLATE_API_URL")
		os.Unsetenv("WEBLATE_API_TOKEN")
	}
}