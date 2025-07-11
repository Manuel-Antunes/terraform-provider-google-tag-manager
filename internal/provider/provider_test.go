package provider

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	ProviderNameEcho = "gtm"
)

func Context(t *testing.T) context.Context {
	t.Helper()
	ctx := t.Context()
	return ctx
}

func ProtoV6ProviderFactories(_ context.Context, providerNames ...string) map[string]func() (tfprotov6.ProviderServer, error) {
	factories := make(map[string]func() (tfprotov6.ProviderServer, error))

	for _, name := range providerNames {
		if name == ProviderNameEcho {
			factories[name] = providerserver.NewProtocol6WithError(New())
		}
	}

	return factories
}

func testAccPreCheck(t *testing.T) {
	// Wait before starting any test to prevent rate limiting
	GlobalTestCoordinator.WaitBeforeRequest()

	// Verify required environment variables are set for acceptance tests
	requiredEnvVars := []string{
		"GTM_CREDENTIAL_FILE",
		"GTM_ACCOUNT_ID",
		"GTM_CONTAINER_ID",
		"GTM_WORKSPACE_NAME",
	}

	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			t.Fatalf("Environment variable %s must be set for acceptance tests", envVar)
		}
	}
}

// TestProvider checks if the provider can be instantiated
func TestProvider(t *testing.T) {
	provider := New()
	if provider == nil {
		t.Fatal("Failed to create provider")
	}
}

// Test workspace creation and reading
func TestAccWorkspaceResource_createAndRead(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkspaceResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gtm_workspace.test", "id"),
					resource.TestCheckResourceAttr("gtm_workspace.test", "name", "tf-test-workspace"),
					resource.TestCheckResourceAttr("gtm_workspace.test", "description", "Created by Terraform"),
				),
			},
		},
	})
}

// Test workspace import
func TestAccWorkspaceResource_import(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkspaceResourceConfig(),
			},
			{
				ResourceName:      "gtm_workspace.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test workspace update
func TestAccWorkspaceResource_update(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkspaceResourceConfig(),
			},
			{
				Config: testAccWorkspaceResourceUpdateConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gtm_workspace.test", "id"),
					resource.TestCheckResourceAttr("gtm_workspace.test", "name", "tf-test-workspace-updated"),
					resource.TestCheckResourceAttr("gtm_workspace.test", "description", "Updated by Terraform"),
				),
			},
		},
	})
}

// Test tag creation and reading
func TestAccTagResource_createAndRead(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gtm_tag.test", "id"),
					resource.TestCheckResourceAttr("gtm_tag.test", "name", "tf-test-tag"),
					resource.TestCheckResourceAttr("gtm_tag.test", "type", "gaawe"),
					resource.TestCheckResourceAttr("gtm_tag.test", "notes", "Created by Terraform"),
					// Check parameters
					resource.TestCheckResourceAttr("gtm_tag.test", "parameter.0.key", "eventName"),
					resource.TestCheckResourceAttr("gtm_tag.test", "parameter.0.value", "test_event"),
					resource.TestCheckResourceAttr("gtm_tag.test", "parameter.1.key", "measurementIdOverride"),
					resource.TestCheckResourceAttr("gtm_tag.test", "parameter.1.value", "G-XXXXXX"),
				),
			},
		},
	})
}

// Test tag import
func TestAccTagResource_import(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceConfig(),
			},
			{
				ResourceName:      "gtm_tag.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test tag update
func TestAccTagResource_update(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceConfig(),
			},
			{
				Config: testAccTagResourceUpdateConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gtm_tag.test", "id"),
					resource.TestCheckResourceAttr("gtm_tag.test", "name", "tf-test-tag-updated"),
					resource.TestCheckResourceAttr("gtm_tag.test", "notes", "Updated by Terraform"),
					// Check updated parameters
					resource.TestCheckResourceAttr("gtm_tag.test", "parameter.0.key", "eventName"),
					resource.TestCheckResourceAttr("gtm_tag.test", "parameter.0.value", "updated_event"),
				),
			},
		},
	})
}

// Test tag with complex parameters
func TestAccTagResource_complexParameters(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceWithComplexParametersConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gtm_tag.complex", "id"),
					resource.TestCheckResourceAttr("gtm_tag.complex", "name", "tf-test-tag-complex"),
					resource.TestCheckResourceAttr("gtm_tag.complex", "type", "gaawe"),
					// Check the complex parameter structure
					resource.TestCheckResourceAttr("gtm_tag.complex", "parameter.2.key", "eventParameters"),
					resource.TestCheckResourceAttr("gtm_tag.complex", "parameter.2.type", "list"),
				),
			},
		},
	})
}

// Test variable creation and reading
func TestAccVariableResource_createAndRead(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config: testAccVariableResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gtm_variable.test", "id"),
					resource.TestCheckResourceAttr("gtm_variable.test", "name", "tf-test-variable"),
					resource.TestCheckResourceAttr("gtm_variable.test", "type", "v"),
					resource.TestCheckResourceAttr("gtm_variable.test", "notes", "Created by Terraform"),
					// Check parameters
					resource.TestCheckResourceAttr("gtm_variable.test", "parameter.0.key", "name"),
					resource.TestCheckResourceAttr("gtm_variable.test", "parameter.0.value", "test-param"),
					resource.TestCheckResourceAttr("gtm_variable.test", "parameter.1.key", "value"),
					resource.TestCheckResourceAttr("gtm_variable.test", "parameter.1.value", "test-value"),
				),
			},
		},
	})
}

// Test variable import
func TestAccVariableResource_import(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config: testAccVariableResourceConfig(),
			},
			{
				ResourceName:      "gtm_variable.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test variable update
func TestAccVariableResource_update(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config: testAccVariableResourceConfig(),
			},
			{
				Config: testAccVariableResourceUpdateConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gtm_variable.test", "123"),
					resource.TestCheckResourceAttr("gtm_variable.test", "name", "tf-test-variable-updated"),
					resource.TestCheckResourceAttr("gtm_variable.test", "notes", "Updated by Terraform"),
					// Check updated parameters
					resource.TestCheckResourceAttr("gtm_variable.test", "parameter.0.key", "name"),
					resource.TestCheckResourceAttr("gtm_variable.test", "parameter.0.value", "updated-param"),
					resource.TestCheckResourceAttr("gtm_variable.test", "parameter.1.key", "value"),
					resource.TestCheckResourceAttr("gtm_variable.test", "parameter.1.value", "updated-value"),
				),
			},
		},
	})
}

// Test trigger creation and reading
func TestAccTriggerResource_createAndRead(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config: testAccTriggerResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gtm_trigger.test", "id"),
					resource.TestCheckResourceAttr("gtm_trigger.test", "name", "tf-test-trigger"),
					resource.TestCheckResourceAttr("gtm_trigger.test", "type", "customEvent"),
					resource.TestCheckResourceAttr("gtm_trigger.test", "notes", "Created by Terraform"),
					// Check custom event filter
					resource.TestCheckResourceAttr("gtm_trigger.test", "custom_event_filter.0.type", "equals"),
					resource.TestCheckResourceAttr("gtm_trigger.test", "custom_event_filter.0.parameter.0.key", "arg0"),
					resource.TestCheckResourceAttr("gtm_trigger.test", "custom_event_filter.0.parameter.0.value", "{{_event}}"),
					resource.TestCheckResourceAttr("gtm_trigger.test", "custom_event_filter.0.parameter.1.key", "arg1"),
					resource.TestCheckResourceAttr("gtm_trigger.test", "custom_event_filter.0.parameter.1.value", "test-event"),
				),
			},
		},
	})
}

// Test trigger import
func TestAccTriggerResource_import(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config: testAccTriggerResourceConfig(),
			},
			{
				ResourceName:      "gtm_trigger.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test trigger update (changing from custom event to click trigger)
func TestAccTriggerResource_update(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config: testAccTriggerResourceConfig(),
			},
			{
				Config: testAccTriggerResourceUpdateConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gtm_trigger.test", "id"),
					resource.TestCheckResourceAttr("gtm_trigger.test", "name", "tf-test-trigger-updated"),
					resource.TestCheckResourceAttr("gtm_trigger.test", "type", "click"),
					resource.TestCheckResourceAttr("gtm_trigger.test", "notes", "Updated by Terraform"),
					// Check click parameters
					resource.TestCheckResourceAttr("gtm_trigger.test", "parameter.0.key", "clickText"),
					resource.TestCheckResourceAttr("gtm_trigger.test", "parameter.0.value", "UpdatedButton"),
				),
			},
		},
	})
}

// Test configurations for each resource type
func testAccProviderConfig() string {
	retryLimit := 15
	return fmt.Sprintf(`
provider "gtm" {
  credential_file = %q
  account_id      = %q
  container_id    = %q
  workspace_name  = %q
  retry_limit     = %d  # Higher retry limit for tests to handle rate limits
}
`,
		os.Getenv("GTM_CREDENTIAL_FILE"),
		os.Getenv("GTM_ACCOUNT_ID"),
		os.Getenv("GTM_CONTAINER_ID"),
		os.Getenv("GTM_WORKSPACE_NAME"),
		retryLimit,
	)
}

func testAccWorkspaceResourceConfig() string {
	return testAccProviderConfig() + `
resource "gtm_workspace" "test" {
  name        = "tf-test-workspace"
  description = "Created by Terraform"
}
`
}

func testAccWorkspaceResourceUpdateConfig() string {
	return testAccProviderConfig() + `
resource "gtm_workspace" "test" {
  name        = "tf-test-workspace-updated"
  description = "Updated by Terraform"
}
`
}

func testAccTagResourceConfig() string {
	return testAccProviderConfig() + `
resource "gtm_tag" "test" {
  name  = "tf-test-tag"
  type  = "gaawe"
  notes = "Created by Terraform"
  
  parameter = [
	{
    key   = "eventName"
    type  = "template"
    value = "test_event"
  },
	 {
    key   = "measurementIdOverride"
    type  = "template"
    value = "G-XXXXXX"
  },
	{
    key  = "eventParameters"
    type = "list"
    list = [{
      type = "map"
      
      map = [{
        key   = "name"
        type  = "template"
        value = "param-name"
      },{
        key   = "value"
        type  = "template"
        value = "param-value"
      }]
    }]
  }
	]
}
`
}

func testAccTagResourceUpdateConfig() string {
	return testAccProviderConfig() + `
resource "gtm_tag" "test" {
  name  = "tf-test-tag-updated"
  type  = "gaawe"
  notes = "Updated by Terraform"
  
  parameter = [
	{
    key   = "eventName"
    type  = "template"
    value = "updated_event"
  }, {
    key   = "measurementIdOverride"
    type  = "template"
    value = "G-XXXXXX"
  }]
}
`
}

func testAccVariableResourceConfig() string {
	return testAccProviderConfig() + `
resource "gtm_variable" "test" {
  name  = "tf-test-variable"
  type  = "v"
  notes = "Created by Terraform"
  
  parameter = [{
    key   = "name"
    type  = "template"
    value = "test-param"
  }, {
    key   = "value"
    type  = "template"
    value = "test-value"
  }]
}
`
}

func testAccVariableResourceUpdateConfig() string {
	return testAccProviderConfig() + `
resource "gtm_variable" "test" {
  name  = "tf-test-variable-updated"
  type  = "v"
  notes = "Updated by Terraform"
  
  parameter = [{
    key   = "name"
    type  = "template"
    value = "updated-param"
  },
  {
    key   = "value"
    type  = "template"
    value = "updated-value"
  }]
}
`
}

func testAccTriggerResourceConfig() string {
	return testAccProviderConfig() + `
resource "gtm_trigger" "test" {
  name  = "tf-test-trigger"
  type  = "customEvent"
  notes = "Created by Terraform"
	custom_event_filter = [
    {
      type = "equals",
      parameter = [
        {
          type  = "template",
          key   = "arg0",
          value = "{{_event}}"
        },
        {
          type  = "template",
          key   = "arg1",
          value = "test-event"
        }
      ]
    }
  ]
}
`
}

func testAccTriggerResourceUpdateConfig() string {
	return testAccProviderConfig() + `
resource "gtm_trigger" "test" {
  name  = "tf-test-trigger-updated"
  type  = "click"
  notes = "Updated by Terraform"
  
  parameter {
    key   = "clickText"
    value = "UpdatedButton"
    type  = "template"
  }
}
`
}

func testAccTagResourceWithComplexParametersConfig() string {
	return testAccProviderConfig() + `
resource "gtm_tag" "complex" {
  name  = "tf-test-tag-complex"
  type  = "gaawe"
  notes = "Created by Terraform with complex parameters"
  
  parameter  = [
		{
			key   = "eventName"
			type  = "template"
			value = "complex_event"
		},
		{
			key   = "measurementIdOverride"
			type  = "template"
			value = "G-XXXXXX"
		}, 
		{
			key  = "eventParameters"
			type = "list"
			
			list = [{
				type = "map"
				
				map = [{
					key   = "name"
					type  = "template"
					value = "param1"
				},{
					key   = "value"
					type  = "template"
					value = "value1"
				}]
			},{
				type = "map"
				
				map = [{
					key   = "name"
					type  = "template"
					value = "param2"
				}, {
					key   = "value"
					type  = "template"
					value = "value2"
				}]
			}]
		}
	]
`
}
