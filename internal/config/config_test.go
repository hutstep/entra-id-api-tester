package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig_ValidFile(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	configContent := `{
		"endpoints": [
			{
				"name": "Test Endpoint",
				"url": "https://api.example.com",
				"method": "GET",
				"clientId": "test-client-id",
				"clientSecret": "test-secret",
				"tenantId": "test-tenant",
				"scope": "test-scope"
			}
		]
	}`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Load config
	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(config.Endpoints) != 1 {
		t.Errorf("Expected 1 endpoint, got %d", len(config.Endpoints))
	}

	endpoint := config.Endpoints[0]
	if endpoint.Name != "Test Endpoint" {
		t.Errorf("Expected name 'Test Endpoint', got '%s'", endpoint.Name)
	}
	if endpoint.Method != "GET" {
		t.Errorf("Expected method 'GET', got '%s'", endpoint.Method)
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := LoadConfig("nonexistent.json")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.json")

	if err := os.WriteFile(configPath, []byte("not valid json"), 0644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	_, err := LoadConfig(configPath)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestLoadConfig_EmptyEndpoints(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "empty.json")

	configContent := `{"endpoints": []}`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	_, err := LoadConfig(configPath)
	if err == nil {
		t.Error("Expected error for empty endpoints, got nil")
	}
}

func TestEndpointValidate_ValidEndpoint(t *testing.T) {
	endpoint := Endpoint{
		Name:         "Test",
		URL:          "https://api.example.com",
		Method:       "GET",
		ClientID:     "client-id",
		ClientSecret: "secret",
		TenantID:     "tenant",
		Scope:        "scope",
	}

	err := endpoint.Validate()
	if err != nil {
		t.Errorf("Unexpected error for valid endpoint: %v", err)
	}
}

func TestEndpointValidate_MissingFields(t *testing.T) {
	tests := []struct {
		name     string
		endpoint Endpoint
	}{
		{
			"missing name",
			Endpoint{URL: "url", Method: "GET", ClientID: "id", ClientSecret: "secret", TenantID: "tenant", Scope: "scope"},
		},
		{
			"missing url",
			Endpoint{Name: "test", Method: "GET", ClientID: "id", ClientSecret: "secret", TenantID: "tenant", Scope: "scope"},
		},
		{
			"missing method",
			Endpoint{Name: "test", URL: "url", ClientID: "id", ClientSecret: "secret", TenantID: "tenant", Scope: "scope"},
		},
		{
			"missing clientId",
			Endpoint{Name: "test", URL: "url", Method: "GET", ClientSecret: "secret", TenantID: "tenant", Scope: "scope"},
		},
		{
			"missing clientSecret",
			Endpoint{Name: "test", URL: "url", Method: "GET", ClientID: "id", TenantID: "tenant", Scope: "scope"},
		},
		{
			"missing tenantId",
			Endpoint{Name: "test", URL: "url", Method: "GET", ClientID: "id", ClientSecret: "secret", Scope: "scope"},
		},
		{
			"missing scope",
			Endpoint{Name: "test", URL: "url", Method: "GET", ClientID: "id", ClientSecret: "secret", TenantID: "tenant"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.endpoint.Validate()
			if err == nil {
				t.Errorf("Expected error for %s, got nil", tt.name)
			}
		})
	}
}

func TestEndpointValidate_InvalidHTTPMethod(t *testing.T) {
	endpoint := Endpoint{
		Name:         "Test",
		URL:          "https://api.example.com",
		Method:       "INVALID",
		ClientID:     "client-id",
		ClientSecret: "secret",
		TenantID:     "tenant",
		Scope:        "scope",
	}

	err := endpoint.Validate()
	if err == nil {
		t.Error("Expected error for invalid HTTP method, got nil")
	}
}

func TestEndpointValidate_AllValidHTTPMethods(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			endpoint := Endpoint{
				Name:         "Test",
				URL:          "https://api.example.com",
				Method:       method,
				ClientID:     "client-id",
				ClientSecret: "secret",
				TenantID:     "tenant",
				Scope:        "scope",
			}

			err := endpoint.Validate()
			if err != nil {
				t.Errorf("Unexpected error for valid method %s: %v", method, err)
			}
		})
	}
}

func TestLoadConfig_MultipleEndpoints(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	configContent := `{
		"endpoints": [
			{
				"name": "Endpoint 1",
				"url": "https://api1.example.com",
				"method": "GET",
				"clientId": "client1",
				"clientSecret": "secret1",
				"tenantId": "tenant1",
				"scope": "scope1"
			},
			{
				"name": "Endpoint 2",
				"url": "https://api2.example.com",
				"method": "POST",
				"clientId": "client2",
				"clientSecret": "secret2",
				"tenantId": "tenant2",
				"scope": "scope2",
				"requestBody": {"key": "value"}
			}
		]
	}`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(config.Endpoints) != 2 {
		t.Errorf("Expected 2 endpoints, got %d", len(config.Endpoints))
	}

	if config.Endpoints[1].RequestBody == nil {
		t.Error("Expected request body for second endpoint")
	}
}
