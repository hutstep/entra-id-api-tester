package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hutstep/entra-id-api-tester/internal/auth"
	"github.com/hutstep/entra-id-api-tester/internal/client"
	"github.com/hutstep/entra-id-api-tester/internal/config"
)

const (
	defaultConfigPath = "config.json"
)

// TestResult represents the result of testing an endpoint
type TestResult struct {
	EndpointName    string
	ErrorMessage    string
	Duration        time.Duration
	StatusCode      int
	Success         bool
	AuthSuccess     bool
	ConnectSuccess  bool
	ResponseSuccess bool
}

func main() {
	// Parse command-line flags
	configPath := flag.String("config", defaultConfigPath, "Path to configuration file")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	fmt.Printf("Loaded configuration with %d endpoint(s)\n", len(cfg.Endpoints))
	fmt.Println("=" + repeat("=", 78))

	// Initialize auth and API clients
	tokenProvider := auth.NewEntraIDTokenProvider()
	apiClient := client.NewAPIClient()

	// Test each endpoint
	results := make([]TestResult, 0, len(cfg.Endpoints))
	for i := range cfg.Endpoints {
		endpoint := &cfg.Endpoints[i]
		fmt.Printf("\n[%d/%d] Testing: %s\n", i+1, len(cfg.Endpoints), endpoint.Name)
		fmt.Printf("    URL: %s\n", endpoint.URL)
		fmt.Printf("    Method: %s\n", endpoint.Method)

		result := testEndpoint(context.Background(), endpoint, tokenProvider, apiClient, *verbose)
		results = append(results, result)

		printTestResult(result)
	}

	// Print summary
	fmt.Println("\n" + repeat("=", 80))
	printSummary(results)

	// Exit with appropriate code
	if hasFailures(results) {
		os.Exit(1)
	}
}

// testEndpoint tests a single endpoint
func testEndpoint(ctx context.Context, endpoint *config.Endpoint, tokenProvider auth.TokenProvider, apiClient *client.APIClient, verbose bool) TestResult {
	result := TestResult{
		EndpointName: endpoint.Name,
	}

	startTime := time.Now()

	// Step 1: Authenticate
	if verbose {
		fmt.Println("    → Authenticating...")
	}

	token, err := tokenProvider.GetAccessToken(ctx, endpoint.ClientID, endpoint.ClientSecret, endpoint.TenantID, endpoint.Scope)
	if err != nil {
		result.ErrorMessage = fmt.Sprintf("Authentication failed: %v", err)
		result.Duration = time.Since(startTime)
		return result
	}

	result.AuthSuccess = true
	if verbose {
		fmt.Println("    ✓ Authentication successful")
	}

	// Step 2: Make API call
	if verbose {
		fmt.Println("    → Making API request...")
	}

	response, err := apiClient.CallAPI(ctx, endpoint.Method, endpoint.URL, token, endpoint.RequestBody)
	if err != nil {
		result.ErrorMessage = fmt.Sprintf("Request failed: %v", err)
		result.Duration = time.Since(startTime)
		return result
	}

	result.ConnectSuccess = true
	result.StatusCode = response.StatusCode
	result.Duration = time.Since(startTime)

	if verbose {
		fmt.Printf("    ✓ Request completed (Status: %d)\n", response.StatusCode)
	}

	// Step 3: Check response
	if response.IsSuccessStatusCode() {
		result.ResponseSuccess = true
		result.Success = true
	} else {
		result.ErrorMessage = fmt.Sprintf("Unexpected status code: %d", response.StatusCode)
		if verbose && len(response.Body) > 0 {
			fmt.Printf("    Response body: %s\n", response.GetBodyAsString())
		}
	}

	return result
}

// printTestResult prints the result of a single test
func printTestResult(result TestResult) {
	if result.Success {
		fmt.Printf("    ✓ PASS - All checks passed (Duration: %v)\n", result.Duration)
	} else {
		fmt.Printf("    ✗ FAIL - %s (Duration: %v)\n", result.ErrorMessage, result.Duration)
		if !result.AuthSuccess {
			fmt.Println("      • Authentication: FAILED")
		} else {
			fmt.Println("      • Authentication: PASSED")
		}
		if result.ConnectSuccess {
			fmt.Println("      • Connectivity: PASSED")
			fmt.Printf("      • Response Status: FAILED (Status Code: %d)\n", result.StatusCode)
		} else {
			fmt.Println("      • Connectivity: FAILED")
		}
	}
}

// printSummary prints a summary of all test results
func printSummary(results []TestResult) {
	total := len(results)
	passed := 0
	authFailed := 0
	connectFailed := 0
	responseFailed := 0

	for _, result := range results {
		if result.Success {
			passed++
		} else {
			switch {
			case !result.AuthSuccess:
				authFailed++
			case !result.ConnectSuccess:
				connectFailed++
			default:
				responseFailed++
			}
		}
	}

	fmt.Println("SUMMARY")
	fmt.Println(repeat("-", 80))
	fmt.Printf("Total Endpoints:           %d\n", total)
	fmt.Printf("Passed:                    %d (%.1f%%)\n", passed, float64(passed)/float64(total)*100)
	fmt.Printf("Failed:                    %d (%.1f%%)\n", total-passed, float64(total-passed)/float64(total)*100)
	fmt.Println()
	fmt.Printf("  • Authentication Failures:  %d\n", authFailed)
	fmt.Printf("  • Connectivity Failures:    %d\n", connectFailed)
	fmt.Printf("  • Response Failures:        %d\n", responseFailed)
	fmt.Println(repeat("=", 80))
}

// hasFailures checks if any tests failed
func hasFailures(results []TestResult) bool {
	for _, result := range results {
		if !result.Success {
			return true
		}
	}
	return false
}

// repeat repeats a string n times
func repeat(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}
