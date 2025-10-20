package auth

import (
	"context"
	"testing"
	"time"
)

// MockTokenProvider is a mock implementation of TokenProvider for testing
type MockTokenProvider struct {
	TokenToReturn string
	ErrorToReturn error
}

func (m *MockTokenProvider) GetAccessToken(ctx context.Context, clientID, clientSecret, tenantID, scope string) (string, error) {
	if m.ErrorToReturn != nil {
		return "", m.ErrorToReturn
	}
	return m.TokenToReturn, nil
}

func TestNewEntraIDTokenProvider(t *testing.T) {
	provider := NewEntraIDTokenProvider()
	if provider == nil {
		t.Fatal("Expected provider to be non-nil")
	}
	if provider.timeout != 30*time.Second {
		t.Errorf("Expected timeout to be 30s, got %v", provider.timeout)
	}
}

func TestNewEntraIDTokenProviderWithTimeout(t *testing.T) {
	customTimeout := 60 * time.Second
	provider := NewEntraIDTokenProviderWithTimeout(customTimeout)
	if provider == nil {
		t.Fatal("Expected provider to be non-nil")
	}
	if provider.timeout != customTimeout {
		t.Errorf("Expected timeout to be %v, got %v", customTimeout, provider.timeout)
	}
}

func TestGetAccessToken_InvalidCredentials(t *testing.T) {
	provider := NewEntraIDTokenProviderWithTimeout(5 * time.Second)
	ctx := context.Background()

	// Test with invalid credentials (should fail quickly)
	token, err := provider.GetAccessToken(
		ctx,
		"invalid-client-id",
		"invalid-client-secret",
		"invalid-tenant-id",
		"invalid-scope",
	)

	if err == nil {
		t.Error("Expected error with invalid credentials, got nil")
	}
	if token != "" {
		t.Errorf("Expected empty token with invalid credentials, got %s", token)
	}
}

func TestGetAccessToken_EmptyParameters(t *testing.T) {
	provider := NewEntraIDTokenProvider()
	ctx := context.Background()

	tests := []struct {
		name         string
		clientID     string
		clientSecret string
		tenantID     string
		scope        string
	}{
		{"empty clientID", "", "secret", "tenant", "scope"},
		{"empty clientSecret", "client", "", "tenant", "scope"},
		{"empty tenantID", "client", "secret", "", "scope"},
		{"empty scope", "client", "secret", "tenant", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := provider.GetAccessToken(ctx, tt.clientID, tt.clientSecret, tt.tenantID, tt.scope)
			if err == nil {
				t.Error("Expected error with empty parameter, got nil")
			}
			if token != "" {
				t.Errorf("Expected empty token, got %s", token)
			}
		})
	}
}

func TestGetAccessToken_ContextCancellation(t *testing.T) {
	provider := NewEntraIDTokenProvider()

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	token, err := provider.GetAccessToken(
		ctx,
		"client-id",
		"client-secret",
		"tenant-id",
		"scope",
	)

	if err == nil {
		t.Error("Expected error with cancelled context, got nil")
	}
	if token != "" {
		t.Errorf("Expected empty token with cancelled context, got %s", token)
	}
}

func TestMockTokenProvider(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		expectedToken := "mock-token-12345"
		mock := &MockTokenProvider{
			TokenToReturn: expectedToken,
		}

		token, err := mock.GetAccessToken(context.Background(), "id", "secret", "tenant", "scope")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if token != expectedToken {
			t.Errorf("Expected token %s, got %s", expectedToken, token)
		}
	})

	t.Run("error case", func(t *testing.T) {
		mock := &MockTokenProvider{
			ErrorToReturn: context.DeadlineExceeded,
		}

		token, err := mock.GetAccessToken(context.Background(), "id", "secret", "tenant", "scope")
		if err == nil {
			t.Error("Expected error, got nil")
		}
		if token != "" {
			t.Errorf("Expected empty token, got %s", token)
		}
	})
}
