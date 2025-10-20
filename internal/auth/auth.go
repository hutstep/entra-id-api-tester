// Package auth provides authentication functionality for Microsoft Entra ID
// using the Azure Identity SDK. It implements OAuth 2.0 client credentials
// flow for service-to-service authentication.
package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

// TokenProvider defines the interface for token acquisition
type TokenProvider interface {
	GetAccessToken(ctx context.Context, clientID, clientSecret, tenantID, scope string) (string, error)
}

// EntraIDTokenProvider implements TokenProvider using Azure Identity SDK
type EntraIDTokenProvider struct {
	timeout time.Duration
}

// NewEntraIDTokenProvider creates a new EntraIDTokenProvider with default timeout
func NewEntraIDTokenProvider() *EntraIDTokenProvider {
	return &EntraIDTokenProvider{
		timeout: 30 * time.Second,
	}
}

// NewEntraIDTokenProviderWithTimeout creates a new EntraIDTokenProvider with custom timeout
func NewEntraIDTokenProviderWithTimeout(timeout time.Duration) *EntraIDTokenProvider {
	return &EntraIDTokenProvider{
		timeout: timeout,
	}
}

// GetAccessToken acquires an access token using client credentials flow
func (p *EntraIDTokenProvider) GetAccessToken(ctx context.Context, clientID, clientSecret, tenantID, scope string) (string, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	// Create client secret credential
	credential, err := azidentity.NewClientSecretCredential(
		tenantID,
		clientID,
		clientSecret,
		nil, // Use default options
	)
	if err != nil {
		return "", fmt.Errorf("failed to create credential: %w", err)
	}

	// Acquire token
	token, err := credential.GetToken(ctx, policy.TokenRequestOptions{
		Scopes: []string{scope},
	})
	if err != nil {
		return "", fmt.Errorf("failed to acquire token: %w", err)
	}

	if token.Token == "" {
		return "", fmt.Errorf("received empty token")
	}

	return token.Token, nil
}
