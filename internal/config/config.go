// Package config handles loading and validation of API testing configuration
// from JSON files. It defines the structure for endpoints, credentials, and
// request parameters needed for testing Entra ID-protected APIs.
package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Endpoint represents a single API endpoint to test
type Endpoint struct {
	RequestBody  map[string]interface{}
	Name         string
	URL          string
	Method       string
	ClientID     string
	ClientSecret string
	TenantID     string
	Scope        string
}

// Config represents the complete configuration
type Config struct {
	Endpoints []Endpoint `json:"endpoints"`
}

// LoadConfig loads the configuration from a JSON file
func LoadConfig(filePath string) (*Config, error) {
	// Open the file
	file, err := os.Open(filePath) // #nosec G304 - file path is provided by user via CLI flag
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close config file: %w", closeErr)
		}
	}()

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
