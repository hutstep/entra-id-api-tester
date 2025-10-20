package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// MockHTTPClient is a mock implementation of HTTPClient for testing
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestNewAPIClient(t *testing.T) {
	client := NewAPIClient()
	if client == nil {
		t.Fatal("Expected client to be non-nil")
	}
	if client.timeout != 30*time.Second {
		t.Errorf("Expected timeout to be 30s, got %v", client.timeout)
	}
}

func TestNewAPIClientWithTimeout(t *testing.T) {
	customTimeout := 60 * time.Second
	client := NewAPIClientWithTimeout(customTimeout)
	if client == nil {
		t.Fatal("Expected client to be non-nil")
	}
	if client.timeout != customTimeout {
		t.Errorf("Expected timeout to be %v, got %v", customTimeout, client.timeout)
	}
}

func TestCallAPI_GET_Success(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			t.Errorf("Expected Bearer token in Authorization header, got %s", authHeader)
		}

		// Send response
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}))
	defer server.Close()

	// Create client and make request
	client := NewAPIClient()
	resp, err := client.CallAPI(context.Background(), "GET", server.URL, "test-token", nil)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	if !resp.IsSuccessStatusCode() {
		t.Error("Expected IsSuccessStatusCode to return true")
	}
}

func TestCallAPI_POST_WithBody(t *testing.T) {
	requestBody := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
	}

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", contentType)
		}

		// Read and verify body
		var receivedBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&receivedBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if receivedBody["key1"] != "value1" {
			t.Errorf("Expected key1=value1, got %v", receivedBody["key1"])
		}

		// Send response
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "created"})
	}))
	defer server.Close()

	// Create client and make request
	client := NewAPIClient()
	resp, err := client.CallAPI(context.Background(), "POST", server.URL, "test-token", requestBody)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}
}

func TestCallAPI_AllHTTPMethods(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != method {
					t.Errorf("Expected %s method, got %s", method, r.Method)
				}
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			client := NewAPIClient()
			requestBody := map[string]interface{}{}
			if method == "POST" || method == "PUT" || method == "PATCH" {
				requestBody["test"] = "data"
			}

			resp, err := client.CallAPI(context.Background(), method, server.URL, "test-token", requestBody)
			if err != nil {
				t.Fatalf("Unexpected error for %s: %v", method, err)
			}
			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status 200 for %s, got %d", method, resp.StatusCode)
			}
		})
	}
}

func TestCallAPI_ErrorStatusCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
	}))
	defer server.Close()

	client := NewAPIClient()
	resp, err := client.CallAPI(context.Background(), "GET", server.URL, "test-token", nil)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}
	if resp.IsSuccessStatusCode() {
		t.Error("Expected IsSuccessStatusCode to return false for 404")
	}
}

func TestCallAPI_ContextCancellation(t *testing.T) {
	// Create a server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create client with short timeout
	client := NewAPIClientWithTimeout(100 * time.Millisecond)

	_, err := client.CallAPI(context.Background(), "GET", server.URL, "test-token", nil)

	if err == nil {
		t.Error("Expected error due to timeout, got nil")
	}
}

func TestCallAPI_InvalidURL(t *testing.T) {
	client := NewAPIClient()
	_, err := client.CallAPI(context.Background(), "GET", "://invalid-url", "test-token", nil)

	if err == nil {
		t.Error("Expected error with invalid URL, got nil")
	}
}

func TestResponse_GetBodyAsString(t *testing.T) {
	expectedBody := "test response body"
	resp := &Response{
		StatusCode: 200,
		Body:       []byte(expectedBody),
	}

	bodyStr := resp.GetBodyAsString()
	if bodyStr != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, bodyStr)
	}
}

func TestResponse_GetBodyAsJSON(t *testing.T) {
	t.Run("valid json", func(t *testing.T) {
		data := map[string]string{"key": "value"}
		jsonData, _ := json.Marshal(data)

		resp := &Response{
			StatusCode: 200,
			Body:       jsonData,
		}

		var result map[string]string
		err := resp.GetBodyAsJSON(&result)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if result["key"] != "value" {
			t.Errorf("Expected key=value, got %v", result["key"])
		}
	})

	t.Run("empty body", func(t *testing.T) {
		resp := &Response{
			StatusCode: 200,
			Body:       []byte{},
		}

		var result map[string]string
		err := resp.GetBodyAsJSON(&result)
		if err == nil {
			t.Error("Expected error with empty body, got nil")
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		resp := &Response{
			StatusCode: 200,
			Body:       []byte("not json"),
		}

		var result map[string]string
		err := resp.GetBodyAsJSON(&result)
		if err == nil {
			t.Error("Expected error with invalid JSON, got nil")
		}
	})
}

func TestResponse_IsSuccessStatusCode(t *testing.T) {
	tests := []struct {
		statusCode int
		expected   bool
	}{
		{200, true},
		{201, true},
		{204, true},
		{299, true},
		{199, false},
		{300, false},
		{400, false},
		{404, false},
		{500, false},
	}

	for _, tt := range tests {
		t.Run(string(rune(tt.statusCode)), func(t *testing.T) {
			resp := &Response{StatusCode: tt.statusCode}
			result := resp.IsSuccessStatusCode()
			if result != tt.expected {
				t.Errorf("StatusCode %d: expected %v, got %v", tt.statusCode, tt.expected, result)
			}
		})
	}
}

func TestCallAPI_WithMockHTTPClient(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(`{"status":"ok"}`)),
				Header:     make(http.Header),
			}, nil
		},
	}

	client := NewAPIClientWithHTTPClient(mockClient, 30*time.Second)
	resp, err := client.CallAPI(context.Background(), "GET", "http://example.com", "test-token", nil)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}
