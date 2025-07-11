# Integration Tests

This document explains how to run the integration tests for the Google Tag Manager Terraform Provider.

## Prerequisites

To run the integration tests, you'll need:

1. A Google Cloud project with the Tag Manager API enabled
2. A Google Tag Manager account and container
3. A service account with permissions to access the Tag Manager API
4. A service account key file (JSON format)

## Environment Setup

Before running the integration tests, you need to set the following environment variables:

```bash
export GTM_CREDENTIAL_FILE=/path/to/your/service-account.json
export GTM_ACCOUNT_ID=your-gtm-account-id
export GTM_CONTAINER_ID=your-gtm-container-id
export GTM_WORKSPACE_NAME=test-workspace-for-integration-tests
```

## Running Tests

### API Integration Tests

These tests verify that the API client functions correctly interact with the Google Tag Manager API.

```bash
make integration-test
```

### Provider Acceptance Tests

These tests verify that the Terraform provider resources work correctly.

```bash
make acceptance-test
```

### All Tests

To run both API integration tests and provider acceptance tests:

```bash
make test-all
```

## Test Organization

### API Tests

- `internal/api/client_test.go`: Tests for the base API client
- `internal/api/clientInWorkspace_test.go`: Tests for the workspace-specific client operations

### Provider Tests

- `internal/provider/provider_test.go`: Tests for the Terraform provider and resources

## Notes

- Integration tests create real resources in your Google Tag Manager account, so they should be run in a test container.
- Tests automatically clean up resources they create, but if a test fails, some resources might be left behind.
- Running the full acceptance test suite can take several minutes.

## Handling Rate Limits

The Google Tag Manager API has rate limits that can impact test execution. To handle this:

1. Tests include built-in delays between API calls to avoid hitting rate limits
2. The client has configurable retry logic with exponential backoff
3. The test suite runs with reduced parallelism (2 tests at a time)
4. In CI environments, we recommend running tests one at a time

If you encounter rate limit errors:

```
Rate limit exceeded after X retries
```

Try one of these approaches:

- Run fewer tests in parallel by setting `-parallel=1` 
- Increase the retry limit in your tests
- Add longer delays between test runs
- Split your test runs into smaller batches
