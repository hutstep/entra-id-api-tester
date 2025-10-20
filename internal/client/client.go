// Package client provides HTTP client functionality for making API requests
// with Bearer token authentication. It supports all standard HTTP methods
// (GET, POST, PUT, PATCH, DELETE) and includes response parsing utilities.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClient defines the interface for making HTTP requests
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// APIClient handles API requests with authentication
type APIClient struct {
	httpClient HTTPClient
	timeout    time.Duration
}

// NewAPIClient creates a new APIClient with default settings
func NewAPIClient() *APIClient {
	return &APIClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		timeout: 30 * time.Second,
	}
}

// NewAPIClientWithTimeout creates a new APIClient with custom timeout
func NewAPIClientWithTimeout(timeout time.Duration) *APIClient {
	return &APIClient{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		timeout: timeout,
	}
}

// NewAPIClientWithHTTPClient creates a new APIClient with a custom HTTP client
func NewAPIClientWithHTTPClient(httpClient HTTPClient, timeout time.Duration) *APIClient {
	return &APIClient{
		httpClient: httpClient,
		timeout:    timeout,
	}
}

// Response represents an API response
type Response struct {
	Headers    http.Header
	Body       []byte
	StatusCode int
}

// CallAPI makes an HTTP request to the specified endpoint
func (c *APIClient) CallAPI(ctx context.Context, method, url, accessToken string, requestBody map[string]interface{}) (*Response, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// Prepare request body
	var bodyReader io.Reader
	if requestBody != nil && (method == "POST" || method == "PUT" || method == "PATCH") {
		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	if requestBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			// Log error but don't override the main error
			_ = closeErr
		}
	}()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Body:       body,
		Headers:    resp.Header,
	}, nil
}

// IsSuccessStatusCode checks if the status code indicates success
func (r *Response) IsSuccessStatusCode() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// GetBodyAsString returns the response body as a string
func (r *Response) GetBodyAsString() string {
	return string(r.Body)
}

// GetBodyAsJSON attempts to unmarshal the response body as JSON
func (r *Response) GetBodyAsJSON(v interface{}) error {
	if len(r.Body) == 0 {
		return fmt.Errorf("empty response body")
	}
	return json.Unmarshal(r.Body, v)
}
