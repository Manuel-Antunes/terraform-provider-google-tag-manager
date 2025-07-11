# Terraform Provider google-tag-manager

This provider allows you to manage Google Tag Manager resources using Terraform.

## Features

- Manage GTM Workspaces
- Manage GTM Tags
- Manage GTM Triggers
- Manage GTM Variables
- Import existing GTM resources into Terraform state

## Requirements

- Terraform 1.0+
- Go 1.20+ (for development)
- A Google Tag Manager account and container

## Using the Provider

### Configuration in Terraform Files

```hcl
terraform {
  required_providers {
    gtm = {
      source = "Manuel-Antunes/google-tag-manager"
    }
  }
}

provider "gtm" {
  credential_file = "/path/to/service-account.json"
  account_id      = "your-account-id"
  container_id    = "your-container-id"
  workspace_name  = "your-workspace"
}

resource "gtm_tag" "example" {
  name = "example-tag"
  type = "gaawe"
  
  parameter {
    key   = "eventName"
    type  = "template"
    value = "my_event"
  }
}
```

### Using Environment Variables

You can also configure the provider using environment variables:

```hcl
# No need to specify these attributes if using environment variables
provider "gtm" {}
```

Required environment variables:
- `GTM_CREDENTIAL_FILE`: Path to the Google Tag Manager service account credentials file
- `GTM_ACCOUNT_ID`: Your Google Tag Manager account ID
- `GTM_CONTAINER_ID`: Your Google Tag Manager container ID
- `GTM_WORKSPACE_NAME`: Your Google Tag Manager workspace name

Optional environment variables:
- `GTM_RETRY_LIMIT`: Number of retry attempts for API requests (default: 10)

You can use a `.env` file with your development environment to set these variables:

```bash
# Setup .env file from the provided example
make setup-env

# Edit the .env file with your credentials
nano .env
```

## Importing Existing Resources

This provider supports importing existing Google Tag Manager resources into your Terraform state, allowing you to manage pre-existing resources with Terraform.

### Import Syntax

To import a resource, use the `terraform import` command with the resource address and the resource ID:

```bash
terraform import [resource_type].[resource_name] [resource_id]
```

### Importable Resources

The following resources can be imported:

#### Tags

```bash
terraform import gtm_tag.example [tag_id]
```

#### Triggers

```bash
terraform import gtm_trigger.example [trigger_id]
```

#### Variables

```bash
terraform import gtm_variable.example [variable_id]
```

### Import Example

1. First, define an empty resource block in your configuration:

```hcl
resource "gtm_tag" "imported_tag" {
  # The attributes will be filled in by Terraform after import
}
```

2. Run the import command:

```bash
terraform import gtm_tag.imported_tag 1234567890
```

3. Run `terraform plan` to see what attributes Terraform detected and what might need to be adjusted:

```bash
terraform plan
```

4. Update your configuration with the necessary attributes to match the imported resource.

## Testing

The provider includes both unit and integration tests.

### Unit Tests

To run the unit tests:

```bash
make test
```

### Integration Tests

Integration tests validate the provider's functionality against the actual Google Tag Manager API.

To run the integration tests, you'll need to set up environment variables:

```bash
# Option 1: Set environment variables directly
export GTM_CREDENTIAL_FILE=/path/to/service-account.json
export GTM_ACCOUNT_ID=your-gtm-account-id
export GTM_CONTAINER_ID=your-gtm-container-id
export GTM_WORKSPACE_NAME=test-workspace

# Option 2: Use a .env file (recommended for development)
make setup-env   # Creates a .env file from .env.example if it doesn't exist
# Edit the .env file with your credentials
nano .env        # Add your credentials to the .env file

# Run API integration tests
make integration-test

# Run provider acceptance tests
make acceptance-test

# Run all tests (checks environment variables first)
make test-all
```

For more detailed information about running integration tests, see [docs/integration-tests.md](docs/integration-tests.md).
