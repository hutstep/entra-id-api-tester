# API Tester

A robust, Go-based API testing tool designed to verify connectivity, authentication, and successful API responses for Azure-protected endpoints using Microsoft Entra ID (formerly Azure Active Directory) authentication.

## Features

- ✅ **Entra ID Authentication**: Client credentials flow using Azure Identity SDK
- ✅ **Multi-Endpoint Testing**: Test multiple APIs with different configurations
- ✅ **All HTTP Methods**: Support for GET, POST, PUT, PATCH, DELETE
- ✅ **Comprehensive Testing**: Checks connectivity, authentication, and response status
- ✅ **Detailed Reporting**: Console output with pass/fail status and error details
- ✅ **Configuration-Based**: JSON configuration file for endpoints and credentials
- ✅ **Well-Tested**: Comprehensive unit tests for all components
- ✅ **Best Practices**: Built following Go and Azure SDK best practices

## Prerequisites

- Go 1.23 or higher
- Azure service principal credentials (Client ID, Client Secret, Tenant ID)
- Access to the APIs you want to test

## Installation

1. Clone or download this repository
2. Navigate to the project directory
3. Install dependencies:

```bash
go mod download
```

## Configuration

Create a `config.json` file in the project root with your API endpoints.

**Quick Start:** Copy the example configuration:

```bash
cp config.example.json config.json
```

Then edit `config.json` with your actual credentials and endpoints:

```json
{
  "endpoints": [
    {
      "name": "My API - Production",
      "url": "https://api.example.com/endpoint",
      "method": "GET",
      "clientId": "your-client-id",
      "clientSecret": "your-client-secret",
      "tenantId": "your-tenant-id",
      "scope": "api://your-app-id/.default",
      "requestBody": null
    },
    {
      "name": "My API - POST Example",
      "url": "https://api.example.com/resource",
      "method": "POST",
      "clientId": "your-client-id",
      "clientSecret": "your-client-secret",
      "tenantId": "your-tenant-id",
      "scope": "api://your-app-id/.default",
      "requestBody": {
        "key": "value",
        "foo": "bar"
      }
    }
  ]
}
```

### Configuration Fields

| Field | Required | Description |
| --- | --- | --- |
| `name` | Yes | Descriptive name for the endpoint |
| `url` | Yes | Full URL of the API endpoint to test |
| `method` | Yes | HTTP method (GET, POST, PUT, PATCH, DELETE) |
| `clientId` | Yes | Azure AD application (client) ID |
| `clientSecret` | Yes | Azure AD client secret |
| `tenantId` | Yes | Azure AD tenant ID |
| `scope` | Yes | OAuth scope (typically `api://<app-id>/.default`) |
| `requestBody` | No | JSON object for POST/PUT/PATCH requests |

## Usage

### Build the Application

```bash
go build -o api-tester ./cmd/api-tester
```

### Run the Tests

```bash
# Using default config.json
./api-tester

# Using custom config file
./api-tester -config path/to/config.json

# Verbose output
./api-tester -verbose
```

### Command-Line Flags

- `-config`: Path to configuration file (default: `config.json`)
- `-verbose`: Enable verbose output showing detailed test steps

## Example Output

```
Loaded configuration with 2 endpoint(s)
================================================================================

[1/2] Testing: My API - Production
    URL: https://api.example.com/endpoint
    Method: GET
    ✓ PASS - All checks passed (Duration: 1.234s)

[2/2] Testing: My API - POST Example
    URL: https://api.example.com/resource
    Method: POST
    ✗ FAIL - Unexpected status code: 401 (Duration: 567ms)
      • Authentication: PASSED
      • Connectivity: PASSED
      • Response Status: FAILED (Status Code: 401)

================================================================================
SUMMARY
--------------------------------------------------------------------------------
Total Endpoints:           2
Passed:                    1 (50.0%)
Failed:                    1 (50.0%)

  • Authentication Failures:  0
  • Connectivity Failures:    0
  • Response Failures:        1
================================================================================
```

## Testing

Run the unit tests:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run tests for a specific package
go test ./internal/auth
go test ./internal/client
go test ./internal/config
```

## Project Structure

```
.
├── cmd/
│   └── api-tester/
│       └── main.go              # Main application entry point
├── internal/
│   ├── auth/
│   │   ├── auth.go              # Authentication logic
│   │   └── auth_test.go         # Authentication tests
│   ├── client/
│   │   ├── client.go            # HTTP client logic
│   │   └── client_test.go       # HTTP client tests
│   └── config/
│       ├── config.go            # Configuration handling
│       └── config_test.go       # Configuration tests
├── config.example.json          # Example configuration file
├── go.mod                       # Go module definition
├── go.sum                       # Go module checksums
└── README.md                    # This file
```

## How It Works

1. **Configuration Loading**: Reads and validates the JSON configuration file
2. **Authentication**: For each endpoint, acquires an access token from Microsoft Entra ID using client credentials flow
3. **API Request**: Makes an HTTP request to the endpoint with the Bearer token
4. **Response Validation**: Checks if the response status code indicates success (2xx)
5. **Reporting**: Outputs detailed results for each endpoint and a summary

## Security Best Practices

- ✅ Never commit `config.json` with real credentials to version control
- ✅ Use environment-specific configuration files
- ✅ Rotate client secrets regularly
- ✅ Use Azure Key Vault for production deployments
- ✅ Follow the principle of least privilege when assigning API permissions

## Troubleshooting

### Authentication Failures

- Verify your Client ID, Client Secret, and Tenant ID are correct
- Ensure the service principal has the necessary permissions
- Check that the scope matches your API's application ID

### Connectivity Failures

- Verify the URL is correct and accessible
- Check network connectivity and firewall rules
- Ensure DNS resolution is working

### Response Failures

- Check if the API requires specific headers or query parameters
- Verify the HTTP method is correct
- Review the API's documentation for required request format

## Contributing

Contributions are welcome! Here are some ways you can contribute:

### Reporting Issues

- Use GitHub Issues to report bugs or request features
- Include clear reproduction steps for bugs
- Provide sample configurations (without secrets!) when relevant

### Pull Requests

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass (`make test`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

### Development Guidelines

- Follow Go best practices and idiomatic code style
- Maintain or improve test coverage
- Update documentation for new features
- Run `make lint` before committing

### Ideas for Contributions

- Support for custom headers
- Certificate-based authentication
- Response body validation (JSON schema)
- Export results to JSON/CSV formats
- Parallel endpoint testing
- Retry logic with exponential backoff
- Docker container support
- CI/CD pipeline examples

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Related Documentation

- [Azure Identity SDK for Go](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity)
- [Microsoft Entra ID Authentication](https://learn.microsoft.com/en-us/azure/developer/go/sdk/authentication/authentication-overview)
- [OAuth 2.0 Client Credentials Flow](https://learn.microsoft.com/en-us/entra/identity-platform/v2-oauth2-client-creds-grant-flow)
