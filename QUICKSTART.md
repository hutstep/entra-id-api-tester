# Quick Start Guide

## Getting Started in 5 Minutes

### Step 1: Create Your Configuration File

**Option 1: Copy the example** (recommended)

```bash
cp config.example.json config.json
```

**Option 2: Create from scratch**

Create a file named `config.json` in the project root:

```json
{
  "endpoints": [
    {
      "name": "My First API Test",
      "url": "https://your-api-endpoint.com/path",
      "method": "GET",
      "clientId": "your-client-id",
      "clientSecret": "your-client-secret",
      "tenantId": "your-tenant-id",
      "scope": "api://your-app-id/.default",
      "requestBody": null
    }
  ]
}
```

### Step 2: Build the Application

```bash
make build
# or
go build -o api-tester ./cmd/api-tester
```

### Step 3: Run Your Tests

```bash
./api-tester
```

## Getting Your Azure Credentials

You need three pieces of information from Azure Portal:

### 1. Tenant ID

- Go to Azure Portal → Entra ID → Overview
- Copy the "Tenant ID"

### 2. Client ID (Application ID)

- Go to Azure Portal → Entra ID → App registrations
- Select your application
- Copy the "Application (client) ID"

### 3. Client Secret

- In your app registration → Certificates & secrets
- Create a new client secret
- Copy the secret **value** (not the ID)

### 4. Scope

- Usually in the format: `api://<your-app-id>/.default`
- Or provided by the API documentation

## Example Output

```
Loaded configuration with 1 endpoint(s)
================================================================================

[1/1] Testing: My First API Test
    URL: https://your-api-endpoint.com/path
    Method: GET
    ✓ PASS - All checks passed (Duration: 1.234s)

================================================================================
SUMMARY
--------------------------------------------------------------------------------
Total Endpoints:           1
Passed:                    1 (100.0%)
Failed:                    0 (0.0%)

  • Authentication Failures:  0
  • Connectivity Failures:    0
  • Response Failures:        0
================================================================================
```

## Common Issues & Solutions

### Issue: "Failed to load configuration"

**Solution**: Make sure `config.json` exists and is valid JSON.

### Issue: "Authentication failed"

**Solution**:

- Verify your Client ID, Client Secret, and Tenant ID
- Make sure the client secret hasn't expired
- Check if your app has the required API permissions

### Issue: "Request failed: context deadline exceeded"

**Solution**: The API might be slow or unreachable. Check:

- Network connectivity
- API endpoint URL is correct
- Firewall rules

### Issue: "Unexpected status code: 403"

**Solution**: Your app doesn't have permission to access the API.

- Check API permissions in Azure Portal
- Make sure you've granted admin consent

## Testing Multiple Endpoints

Add more objects to the `endpoints` array:

```json
{
  "endpoints": [
    {
      "name": "Production API",
      "url": "https://api-prod.example.com/endpoint",
      "method": "GET",
      ...
    },
    {
      "name": "Staging API",
      "url": "https://api-staging.example.com/endpoint",
      "method": "GET",
      ...
    }
  ]
}
```

## Testing POST Requests

Include a `requestBody` field:

```json
{
  "name": "Create Resource",
  "url": "https://api.example.com/resources",
  "method": "POST",
  "clientId": "...",
  "clientSecret": "...",
  "tenantId": "...",
  "scope": "...",
  "requestBody": {
    "name": "Test Resource",
    "value": 123,
    "active": true
  }
}
```

## Command-Line Options

```bash
# Use a different config file
./api-tester -config my-config.json

# Verbose output (shows detailed steps)
./api-tester -verbose

# See all options
./api-tester -help
```

## Running Tests

```bash
# Run unit tests
make test

# Run with coverage
make test-coverage

# Clean and rebuild
make clean build
```

## Need Help?

Check these files:

- `README.md` - Full documentation
- `config.example.json` - Configuration example
- `PROJECT_SUMMARY.md` - Technical details

## Security Reminder

⚠️ **IMPORTANT**: Never commit `config.json` with real credentials to git!

The `.gitignore` file is already configured to ignore `config.json`, but always double-check before committing.
