# AGENTS.md - AI Agent Context Documentation

## Project Overview

**Project Name:** API Tester (Go)  
**Language:** Go 1.25  
**Purpose:** Test Azure-protected API endpoints using Microsoft Entra ID authentication  
**Created:** October 2025

## Project Context

This is a production-ready Go application that tests API connectivity, authentication, and response validation for Azure-protected APIs.

### Primary Use Case

Testing multiple API endpoints (potentially across different environments/stages) to verify:

1. **Connectivity** - Can we reach the API?
2. **Authentication** - Can we successfully authenticate with Entra ID?
3. **Response** - Does the API return a successful response?

## Architecture & Design Decisions

### Core Design Principles

1. **Interface-based design** - All major components use interfaces for testability
2. **Configuration over code** - All endpoints and credentials in JSON config files
3. **Fail-continue pattern** - Tests all endpoints even if some fail
4. **Clear separation of concerns** - Auth, HTTP client, and config are separate packages
5. **Context-aware** - All operations support context cancellation and timeouts

### Package Structure

```
github.com/hutstep/entra-id-api-tester/
├── cmd/api-tester/              # Main application
│   └── main.go                  # Entry point, CLI, orchestration
├── internal/                    # Private packages
│   ├── auth/                    # Authentication logic
│   │   ├── auth.go              # Entra ID token acquisition
│   │   └── auth_test.go         # Unit tests
│   ├── client/                  # HTTP client
│   │   ├── client.go            # HTTP request handling
│   │   └── client_test.go       # Unit tests with httptest
│   └── config/                  # Configuration
│       ├── config.go            # JSON config loading & validation
│       └── config_test.go       # Unit tests
```

### Why `internal/` Package?

The `internal/` directory prevents these packages from being imported by other projects. This is intentional - these packages are specific to this application and not meant to be reusable libraries.

## Technical Details

### Dependencies

**Required (in go.mod):**

- `github.com/Azure/azure-sdk-for-go/sdk/azcore` v1.16.0 - Core Azure SDK functionality
- `github.com/Azure/azure-sdk-for-go/sdk/azidentity` v1.8.0 - Azure authentication

**Standard Library (no external deps):**

- `net/http` - HTTP client
- `encoding/json` - JSON parsing
- `context` - Context management
- `time` - Timeouts and durations
- `flag` - CLI argument parsing

### Authentication Flow

**Method:** OAuth 2.0 Client Credentials Flow  
**Library:** `azidentity.ClientSecretCredential`  
**Token Caching:** Handled automatically by Azure SDK

```go
// Simplified authentication flow
credential := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)
token := credential.GetToken(ctx, policy.TokenRequestOptions{Scopes: []string{scope}})
```

**Important:** The application does NOT cache tokens itself. The Azure SDK handles token caching and refresh automatically.

### Configuration Format

**File:** `config.json` (not committed to git)  
**Format:** JSON

```json
{
  "endpoints": [
    {
      "name": "string",              // Display name
      "url": "string",               // Full API URL
      "method": "string",            // GET|POST|PUT|PATCH|DELETE
      "clientId": "string",          // Entra ID client ID
      "clientSecret": "string",      // Entra ID client secret
      "tenantId": "string",          // Entra ID tenant ID
      "scope": "string",             // OAuth scope (usually api://app-id/.default)
      "requestBody": object|null     // Request body for POST/PUT/PATCH
    }
  ]
}
```

**Validation Rules:**

- All fields except `requestBody` are required
- `method` must be one of: GET, POST, PUT, PATCH, DELETE
- `requestBody` only used for POST, PUT, PATCH methods
- Empty endpoints array is invalid

### HTTP Client Behavior

**Timeouts:**

- Default: 30 seconds (configurable)
- Applied to entire request/response cycle
- Context-based cancellation supported

**Headers Sent:**

- `Authorization: Bearer <token>`
- `Content-Type: application/json` (for POST/PUT/PATCH)

**Success Criteria:**

- HTTP status codes 200-299 considered successful
- All other status codes are failures (but logged, not fatal)

### Test Results Structure

```go
type TestResult struct {
    EndpointName    string        // Name from config
    Success         bool          // Overall success
    AuthSuccess     bool          // Token acquired?
    ConnectSuccess  bool          // API reachable?
    ResponseSuccess bool          // 2xx status?
    StatusCode      int           // HTTP status code
    ErrorMessage    string        // Error if any
    Duration        time.Duration // Total test time
}
```

**Failure Hierarchy:**

1. If auth fails → ConnectSuccess = false, ResponseSuccess = false
2. If connection fails → ResponseSuccess = false
3. If response not 2xx → ResponseSuccess = false, but log status code

### Command-Line Interface

**Binary Name:** `api-tester`

**Flags:**

- `-config <path>` - Path to config file (default: `config.json`)
- `-verbose` - Enable verbose output (shows auth & request steps)

**Exit Codes:**

- `0` - All tests passed
- `1` - One or more tests failed

### Output Format

**Console Output:**

- Progress indicators during testing
- Pass/fail status with checkmarks (✓/✗)
- Timing information for each endpoint
- Summary statistics at the end
- Detailed failure breakdown

**No File Output:** Currently only console output. Future enhancement could add JSON/CSV export.

## Testing Strategy

### Unit Test Coverage

**auth package (76.9%):**

- Provider initialization
- Invalid credentials handling
- Context cancellation
- Empty parameter validation
- Mock provider implementation

**client package (93.3%):**

- All HTTP methods (GET, POST, PUT, PATCH, DELETE)
- Request body handling
- Error status codes
- Context cancellation/timeout
- Invalid URLs
- Response parsing (string, JSON)
- Success status code checks

**config package (97.1%):**

- Valid/invalid JSON parsing
- Missing required fields
- Invalid HTTP methods
- Multiple endpoints
- File not found errors

### Test Patterns Used

1. **Table-driven tests** - For testing multiple similar cases
2. **httptest.Server** - For mocking HTTP servers
3. **Mock interfaces** - For auth and HTTP client mocking
4. **Temporary files** - For config file testing
5. **Context cancellation** - For timeout testing

### Running Tests

```bash
# All tests
go test ./...

# With coverage
go test -cover ./...

# Verbose
go test -v ./...

# Specific package
go test ./internal/auth
```

## Code Patterns & Conventions

### Error Handling

```go
// Always wrap errors with context
return fmt.Errorf("failed to do X: %w", err)

// Check errors immediately
result, err := doSomething()
if err != nil {
    return fmt.Errorf("descriptive message: %w", err)
}
```

### Context Usage

```go
// Create context with timeout for external operations
ctx, cancel := context.WithTimeout(ctx, timeout)
defer cancel()
```

### Struct Validation

```go
// Each struct has a Validate() error method
func (e *Endpoint) Validate() error {
    if e.Name == "" {
        return fmt.Errorf("name is required")
    }
    // ... more validation
}
```

### Interface Design

```go
// Interfaces are defined in the package that uses them
type TokenProvider interface {
    GetAccessToken(ctx context.Context, ...) (string, error)
}

// Implementations can be in same or different packages
type EntraIDTokenProvider struct { ... }
```

## Important Files

### `.gitignore`

**Purpose:** Prevent committing secrets  
**Key entries:**

- `config.json` - Contains client secrets
- `api-tester` - Binary
- `*.out` - Coverage files

### `go.mod`

**Module name:** `github.com/hutstep/entra-id-api-tester`  
**Note:** This is the import path, not necessarily where the code lives

### `Makefile`

**Targets:**

- `build` - Compile binary
- `test` - Run tests
- `test-coverage` - Generate HTML coverage report
- `clean` - Remove artifacts
- `run` - Build and execute
- `lint` - Run go vet and gofmt

## Known Limitations & Design Choices

### 1. Sequential Testing Only

**Current:** Tests run one endpoint at a time  
**Why:** Simpler implementation, easier debugging  
**Future:** Could add parallel testing with goroutines

### 2. No Token Caching at App Level

**Current:** Relies on Azure SDK's internal caching  
**Why:** Azure SDK handles this better than we could  
**Note:** Each test run may acquire new tokens if SDK cache expires

### 3. No Retry Logic

**Current:** Single attempt per endpoint  
**Why:** This is a testing tool, not a production client  
**Future:** Could add configurable retry with exponential backoff

### 4. Console Output Only

**Current:** Results printed to stdout  
**Why:** Simple and sufficient for most use cases  
**Future:** Could add JSON/CSV export for CI/CD integration

### 5. Request Body Must Be JSON

**Current:** Only supports JSON request bodies  
**Why:** Covers 99% of modern APIs  
**Future:** Could add support for form data, XML, etc.

### 6. No Custom Headers

**Current:** Only sends Authorization and Content-Type headers  
**Why:** Simplifies configuration  
**Future:** Could add custom headers in config

## Common Modifications

### Adding a New HTTP Method

```go
// 1. Add to validation in config/config.go
validMethods := map[string]bool{
    "GET": true, "POST": true, "PUT": true,
    "PATCH": true, "DELETE": true,
    "HEAD": true, // <-- Add new method
}

// 2. No changes needed in client.go (already handles any method)
// 3. Add test in client_test.go
```

### Adding Custom Headers

```go
// In config/config.go, add to Endpoint struct:
type Endpoint struct {
    // ... existing fields
    Headers map[string]string `json:"headers"`
}

// In client/client.go, in CallAPI():
for key, value := range headers {
    req.Header.Set(key, value)
}
```

### Adding Response Body Validation

```go
// In config/config.go:
type Endpoint struct {
    // ... existing fields
    ExpectedStatus int                    `json:"expectedStatus"`
    ExpectedBody   map[string]interface{} `json:"expectedBody"`
}

// In main.go, add validation after API call
if endpoint.ExpectedBody != nil {
    var actual map[string]interface{}
    json.Unmarshal(response.Body, &actual)
    // Compare actual vs expected
}
```

### Adding Parallel Execution

```go
// In main.go:
var wg sync.WaitGroup
resultsChan := make(chan TestResult, len(cfg.Endpoints))

for _, endpoint := range cfg.Endpoints {
    wg.Add(1)
    go func(ep config.Endpoint) {
        defer wg.Done()
        result := testEndpoint(ctx, ep, tokenProvider, apiClient, *verbose)
        resultsChan <- result
    }(endpoint)
}

wg.Wait()
close(resultsChan)

// Collect results
for result := range resultsChan {
    results = append(results, result)
}
```

## Troubleshooting Guide

### "could not import" errors during build

**Cause:** Dependencies not downloaded  
**Fix:** `go mod download` or `go mod tidy`

### Tests timeout

**Cause:** Network issues or slow API  
**Fix:** Increase timeout in test or use mock server

### "no required module provides package"

**Cause:** go.mod is out of sync  
**Fix:** `go mod tidy`

### Authentication fails with valid credentials

**Possible causes:**

1. Client secret expired
2. Scope is incorrect (should usually end with `/.default`)
3. App doesn't have API permissions in Azure Portal
4. Tenant ID is wrong

### All tests fail with "context deadline exceeded"

**Cause:** Timeout too short or network issues  
**Fix:** Check network, increase timeout, or use verbose mode to see where it hangs

## Security Considerations

### Secrets Management

**Current:** Secrets in `config.json` file  
**Production:** Should use:

- Azure Key Vault
- Environment variables
- Kubernetes secrets
- Managed identities (where possible)

### Client Secret Rotation

The application doesn't cache secrets, so rotating is simple:

1. Update `config.json` with new secret
2. Restart application

### Audit Logging

Currently no audit logging. Future enhancement could add:

- Request logging
- Token acquisition logging
- File-based audit trail

## Performance Characteristics

### Memory Usage

- **Minimal** - No significant data structures
- **Per request:** ~1-2 MB
- **Token caching:** Handled by Azure SDK

### Execution Time

- **Per endpoint:** 0.5-2 seconds typical
  - Token acquisition: 0.3-1 second (cached: <10ms)
  - HTTP request: 0.2-1 second (depends on API)
- **Total:** Sequential, so sum of all endpoints

### Scalability

- **Endpoints:** Can handle 100+ endpoints (limited by time, not resources)
- **Concurrent:** Currently sequential only
- **Rate limiting:** No built-in rate limiting

## Documentation Files

### For Users

- **README.md** - Complete user documentation
- **QUICKSTART.md** - 5-minute getting started guide
- **config.example.json** - Configuration template

### For Developers

- **PROJECT_SUMMARY.md** - Technical overview and architecture
- **AGENTS.md** - This file - AI agent context
- **Makefile** - Build automation

### Code Documentation

- Inline comments for complex logic
- Package-level documentation in each `package` declaration
- Exported functions have doc comments

## Version History & Evolution

### Initial Version (October 2025)

- Core functionality: auth, client, config
- Console output
- Sequential testing
- Comprehensive unit tests
- Full documentation

### Future Considerations

- [ ] Parallel execution
- [ ] JSON/CSV export
- [ ] Custom headers support
- [ ] Certificate authentication
- [ ] Response validation
- [ ] Retry logic
- [ ] Prometheus metrics
- [ ] Docker support

## Build & Deployment

### Building

```bash
# Development build
go build -o api-tester ./cmd/api-tester

# Production build (smaller binary)
CGO_ENABLED=0 go build -ldflags="-s -w" -o api-tester ./cmd/api-tester

# Cross-compile for Linux
GOOS=linux GOARCH=amd64 go build -o api-tester-linux ./cmd/api-tester

# Cross-compile for Windows
GOOS=windows GOARCH=amd64 go build -o api-tester.exe ./cmd/api-tester
```

### Deployment Options

1. **Binary distribution** - Copy `api-tester` binary and `config.json`
2. **Container** - Could create Dockerfile (not included yet)
3. **CI/CD** - Run as part of deployment validation

### Dependencies for Running

- **None!** - Single static binary
- Only needs `config.json` at runtime

## AI Agent Instructions

### When Asked to Modify This Project

1. **Read this file first** - Contains all context
2. **Check test coverage** - Maintain >75% coverage
3. **Follow existing patterns** - Interface-based design, error wrapping
4. **Update documentation** - Keep README, QUICKSTART, and this file in sync
5. **Run tests before committing** - `go test ./...`

### When Asked to Debug

1. **Check `get_errors`** - See compiler/lint errors
2. **Run with `-verbose`** - See detailed execution
3. **Check test output** - Tests often reveal issues
4. **Verify go.mod** - Run `go mod tidy` if imports fail

### When Asked to Add Features

1. **Start with tests** - TDD approach preferred
2. **Update config first** - If feature needs config changes
3. **Maintain backwards compatibility** - Old configs should still work
4. **Document in README** - User-facing features need docs

### When Asked About Design Decisions

Refer to these sections:

- **Architecture & Design Decisions** - Why things are structured this way
- **Known Limitations & Design Choices** - Intentional limitations
- **Technical Details** - Implementation specifics

## Contact & Maintenance

**Original Author:** Created via AI assistance  
**Date:** October 2025  
**Language:** Go 1.25  
**Azure SDK Version:** azidentity v1.8.0, azcore v1.16.0

**Maintenance Notes:**

- Azure SDK is actively maintained by Microsoft
- Go 1.25 is current as of October 2025
- No external dependencies beyond Azure SDK
- Tests ensure reliability across updates

---

**Last Updated:** October 20, 2025  
**Document Version:** 1.0  
**AI Agent Version:** Claude 3.5 Sonnet
