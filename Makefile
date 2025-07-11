default: install

# Include .env file if it exists
-include .env

generate:
	go generate ./...

install:
	go install .

test:
	go test -count=1 -parallel=4 ./...

# Run all tests including integration tests and acceptance tests
test-all: check-env integration-test acceptance-test

# Check required environment variables for tests
check-env:
	@echo "Checking required environment variables..."
	@test -n "$(GTM_CREDENTIAL_FILE)" || (echo "Error: GTM_CREDENTIAL_FILE is not set" && exit 1)
	@test -n "$(GTM_ACCOUNT_ID)" || (echo "Error: GTM_ACCOUNT_ID is not set" && exit 1)
	@test -n "$(GTM_CONTAINER_ID)" || (echo "Error: GTM_CONTAINER_ID is not set" && exit 1)
	@test -n "$(GTM_WORKSPACE_NAME)" || (echo "Error: GTM_WORKSPACE_NAME is not set" && exit 1)

testacc:
	TF_ACC=1 go test -count=1 -parallel=2 -timeout 15m -v ./...

# Integration tests
integration-test: install-test-deps
	@echo "Running integration tests..."
	@go test -v -parallel=2 -timeout=15m ./internal/api/...
	@echo "API integration tests completed successfully!"

# Provider acceptance tests (requires environment variables)
acceptance-test: install-test-deps
	@echo "Running provider acceptance tests..."
	@echo "Note: Using environment variables (GTM_CREDENTIAL_FILE, GTM_ACCOUNT_ID, GTM_CONTAINER_ID, GTM_WORKSPACE_NAME)"
	@TF_ACC=1 go test -v -timeout 30m -parallel=2 ./internal/provider/...
	@echo "Provider acceptance tests completed successfully!"

# Install testing dependencies
install-test-deps:
	@echo "Installing testing dependencies..."
	@go get github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource
	@go get github.com/hashicorp/terraform-plugin-testing/helper/resource
	@go get github.com/stretchr/testify/suite

# Setup initial .env file from example if it doesn't exist
setup-env:
	@if [ ! -f .env ]; then \
		echo "Creating initial .env file from .env.example..."; \
		cp .env.example .env; \
		echo "Please edit .env file with your credentials"; \
	else \
		echo ".env file already exists"; \
	fi
