package utils

import "testing"

func TestGetConfigPath(t *testing.T) {
	// case for configPath "development"
	configPath := "development"
	expectedPath := "./config/config-dev"
	result := GetConfigPath(configPath)
	if result != expectedPath {
		t.Errorf("For configPath %s, got %s, expected %s", configPath, result, expectedPath)
	}

	// case for configPath "heroku"
	configPath = "heroku"
	expectedPath = "./config/config-heroku"
	result = GetConfigPath(configPath)
	if result != expectedPath {
		t.Errorf("For configPath %s, got %s, expected %s", configPath, result, expectedPath)
	}

	// case for configPath "production"
	configPath = "production"
	expectedPath = "./config/config-prod"
	result = GetConfigPath(configPath)
	if result != expectedPath {
		t.Errorf("For configPath %s, got %s, expected %s", configPath, result, expectedPath)
	}

	// Case for unknown configPath
	configPath = "unknown"
	expectedPath = "./config/cfg"
	result = GetConfigPath(configPath)
	if result != expectedPath {
		t.Errorf("For configPath %s, got %s, expected %s", configPath, result, expectedPath)
	}
}
