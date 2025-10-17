package config

import (
	"os"
	"testing"
)

func TestDetermineEnvironment(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected string
	}{
		{
			name:     "ECS environment",
			envVars:  map[string]string{"ECS_CONTAINER_METADATA_URI": "some-value"},
			expected: string(EnvDev),
		},
		{
			name:     "AWS execution environment",
			envVars:  map[string]string{"AWS_EXECUTION_ENV": "some-value"},
			expected: string(EnvDev),
		},
		{
			name:     "APP_ENV set to dev",
			envVars:  map[string]string{"APP_ENV": "dev"},
			expected: string(EnvDev),
		},
		{
			name:     "APP_ENV set to DEV",
			envVars:  map[string]string{"APP_ENV": "DEV"},
			expected: string(EnvDev),
		},
		{
			name:     "APP_ENV set to local",
			envVars:  map[string]string{"APP_ENV": "local"},
			expected: string(EnvLocal),
		},
		{
			name:     "APP_ENV set to LOCAL",
			envVars:  map[string]string{"APP_ENV": "LOCAL"},
			expected: string(EnvLocal),
		},
		{
			name:     "No environment variables",
			envVars:  map[string]string{},
			expected: string(EnvLocal),
		},
		{
			name:     "APP_ENV set to invalid value",
			envVars:  map[string]string{"APP_ENV": "prod"},
			expected: string(EnvLocal),
		},
		{
			name:     "APP_ENV set to empty",
			envVars:  map[string]string{"APP_ENV": ""},
			expected: string(EnvLocal),
		},
		{
			name:     "APP_ENV with spaces",
			envVars:  map[string]string{"APP_ENV": " dev "},
			expected: string(EnvDev),
		},
		{
			name:     "Multiple ECS indicators",
			envVars:  map[string]string{"ECS_CONTAINER_METADATA_URI": "value1", "AWS_EXECUTION_ENV": "value2"},
			expected: string(EnvDev),
		},
		{
			name:     "APP_ENV overrides ECS",
			envVars:  map[string]string{"ECS_CONTAINER_METADATA_URI": "value", "APP_ENV": "local"},
			expected: string(EnvLocal),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear relevant environment variables
			envKeys := []string{"ECS_CONTAINER_METADATA_URI", "AWS_EXECUTION_ENV", "APP_ENV"}
			for _, key := range envKeys {
				os.Unsetenv(key)
			}

			// Set test environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			// Call the function
			result := determineEnvironment()

			// Assert
			if result != tt.expected {
				t.Errorf("determineEnvironment() = %v, want %v", result, tt.expected)
			}

			// Clean up
			for key := range tt.envVars {
				os.Unsetenv(key)
			}
		})
	}
}
func TestLoad(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		setupFiles  func() // Optional function to set up config files
		expectedEnv Env
		expectError bool
		validate    func(t *testing.T, cfg *Config)
	}{
		{
			name:        "default local environment",
			envVars:     map[string]string{},
			expectedEnv: EnvLocal,
			expectError: false,
			validate: func(t *testing.T, cfg *Config) {
				if cfg.Env != EnvLocal {
					t.Errorf("Expected Env to be %v, got %v", EnvLocal, cfg.Env)
				}
				if cfg.Log.Level != "info" {
					t.Errorf("Expected Log.Level to be 'info', got %v", cfg.Log.Level)
				}
				if cfg.DynamoDB.Region != "ap-southeast-2" {
					t.Errorf("Expected DynamoDB.Region to be 'ap-southeast-2', got %v", cfg.DynamoDB.Region)
				}
			},
		},
		{
			name:        "dev environment from APP_ENV",
			envVars:     map[string]string{"APP_ENV": "dev"},
			expectedEnv: EnvDev,
			expectError: false,
			validate: func(t *testing.T, cfg *Config) {
				if cfg.Env != EnvDev {
					t.Errorf("Expected Env to be %v, got %v", EnvDev, cfg.Env)
				}
				if cfg.DynamoDB.ShowsTable != "shows-dev" {
					t.Errorf("Expected ShowsTable to be 'shows-dev', got %v", cfg.DynamoDB.ShowsTable)
				}
			},
		},
		{
			name:        "ECS environment",
			envVars:     map[string]string{"ECS_CONTAINER_METADATA_URI": "some-uri"},
			expectedEnv: EnvDev,
			expectError: false,
			validate: func(t *testing.T, cfg *Config) {
				if cfg.Env != EnvDev {
					t.Errorf("Expected Env to be %v, got %v", EnvDev, cfg.Env)
				}
			},
		},
		{
			name:    "environment variable overrides",
			envVars: map[string]string{"APP_DYNAMODB__REGION": "us-west-2", "APP_LOG__LEVEL": "debug"},
			validate: func(t *testing.T, cfg *Config) {
				if cfg.DynamoDB.Region != "us-west-2" {
					t.Errorf("Expected DynamoDB.Region to be 'us-west-2', got %v", cfg.DynamoDB.Region)
				}
				if cfg.Log.Level != "debug" {
					t.Errorf("Expected Log.Level to be 'debug', got %v", cfg.Log.Level)
				}
			},
			expectError: false,
		},
		{
			name:        "invalid APP_ENV",
			envVars:     map[string]string{"APP_ENV": "prod"},
			expectedEnv: EnvLocal, // Should default to local
			expectError: false,
		},
		{
			name:        "empty APP_ENV",
			envVars:     map[string]string{"APP_ENV": ""},
			expectedEnv: EnvLocal,
			expectError: false,
		},
		{
			name:        "case insensitive APP_ENV",
			envVars:     map[string]string{"APP_ENV": "DEV"},
			expectedEnv: EnvDev,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear relevant environment variables
			envKeys := []string{"APP_ENV", "ECS_CONTAINER_METADATA_URI", "AWS_EXECUTION_ENV", "APP_DYNAMODB__REGION", "APP_LOG__LEVEL"}
			for _, key := range envKeys {
				os.Unsetenv(key)
			}

			// Set test environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			// Setup files if provided
			if tt.setupFiles != nil {
				tt.setupFiles()
			}

			// Call the function
			cfg, err := Load()

			// Assert error expectation
			if tt.expectError && err == nil {
				t.Error("Expected an error, but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// If no error expected, validate the config
			if !tt.expectError && cfg != nil {
				if tt.validate != nil {
					tt.validate(t, cfg)
				}
			}

			// Clean up
			for key := range tt.envVars {
				os.Unsetenv(key)
			}
		})
	}
}
