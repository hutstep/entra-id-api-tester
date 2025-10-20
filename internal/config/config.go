package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Endpoint represents a single API endpoint to test
type Endpoint struct {
	Name         string                 `json:"name"`
	URL          string                 `json:"url"`
	Method       string                 `json:"method"`
	ClientID     string                 `json:"clientId"`
	ClientSecret string                 `json:"clientSecret"`
	TenantID     string                 `json:"tenantId"`
	Scope        string                 `json:"scope"`
	RequestBody  map[string]interface{} `json:"requestBody"`
}

// Config represents the complete configuration
type Config struct {
	Endpoints []Endpoint `json:"endpoints"`
}

// LoadConfig loads the configuration from a JSON file
func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if len(c.Endpoints) == 0 {
		return fmt.Errorf("no endpoints defined in configuration")
	}

	for i, endpoint := range c.Endpoints {
		if err := endpoint.Validate(); err != nil {
			return fmt.Errorf("endpoint %d (%s): %w", i, endpoint.Name, err)
		}
	}

	return nil
}

// Validate checks if an endpoint configuration is valid
func (e *Endpoint) Validate() error {
	if e.Name == "" {
		return fmt.Errorf("name is required")
	}
	if e.URL == "" {
		return fmt.Errorf("url is required")
	}
	if e.Method == "" {
		return fmt.Errorf("method is required")
	}

	// Validate HTTP method
	validMethods := map[string]bool{
		"GET": true, "POST": true, "PUT": true,
		"PATCH": true, "DELETE": true,
	}
	if !validMethods[e.Method] {
		return fmt.Errorf("invalid HTTP method: %s (must be GET, POST, PUT, PATCH, or DELETE)", e.Method)
	}

	if e.ClientID == "" {
		return fmt.Errorf("clientId is required")
	}
	if e.ClientSecret == "" {
		return fmt.Errorf("clientSecret is required")
	}
	if e.TenantID == "" {
		return fmt.Errorf("tenantId is required")
	}
	if e.Scope == "" {
		return fmt.Errorf("scope is required")
	}

	return nil
}
