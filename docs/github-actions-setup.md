# Setting Up GitHub Actions for Integration Tests

This document explains how to set up GitHub Actions secrets to enable running integration tests for the Google Tag Manager Terraform Provider in CI/CD.

## Required Secrets

To run integration tests in GitHub Actions, you need to set up the following secrets in your repository:

1. `GTM_CREDENTIALS` - The entire contents of your Google service account JSON key file
2. `GTM_CREDENTIAL_FILE_PATH` - Path where the credentials file will be stored during the workflow (e.g., `/tmp/credentials.json`)
3. `GTM_ACCOUNT_ID` - Your Google Tag Manager account ID
4. `GTM_CONTAINER_ID` - Your Google Tag Manager container ID
5. `GTM_WORKSPACE_NAME` - Name for the test workspace

## Setting Up Secrets

1. Go to your GitHub repository
2. Click on **Settings**
3. In the left sidebar, click on **Secrets and variables** > **Actions**
4. Click on **New repository secret**
5. Add each of the secrets listed above

## Service Account Setup

The service account used for testing should have the following permissions:

1. Tag Manager API access enabled
2. Permissions to read/write/create/delete GTM resources

## Running Acceptance Tests

By default, the acceptance tests (which interact with the real GTM API) only run on pushes to the main branch and when the commit message contains `[run-acc-tests]`. 

To run acceptance tests in a pull request, include `[run-acc-tests]` in your commit message:

```
git commit -m "Feature: Add support for new tag types [run-acc-tests]"
```

## Security Considerations

- The service account used for testing should be restricted to a test GTM container only
- Consider using GitHub environments to restrict which workflows can access these secrets
- Regularly rotate your service account credentials
