#!/bin/bash
set -e

# Install required dependencies
echo "Installing required testing dependencies..."
go get github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource
go get github.com/hashicorp/terraform-plugin-testing/helper/resource
go get github.com/stretchr/testify/suite

echo "Running API integration tests..."
go test -v ./internal/api/...

echo "Running provider integration tests..."
# Note: To run the provider tests, you'll need to set the following environment variables:
# export GTM_CREDENTIAL_FILE=path/to/your/service-account.json
# export GTM_ACCOUNT_ID=your-account-id
# export GTM_CONTAINER_ID=your-container-id
# export GTM_WORKSPACE_NAME=your-workspace-name
TF_ACC=1 go test -v -timeout 30m ./internal/provider/...

echo "All integration tests completed successfully!"
