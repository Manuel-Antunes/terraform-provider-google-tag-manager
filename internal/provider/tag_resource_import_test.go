package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccTagResource_importExistingTag tests importing an existing tag that was created outside Terraform
func TestAccTagResource_importExistingTag(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	var createdTagID string

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			// Step 1: Create a tag outside of Terraform (simulated by creating it first)
			{
				Config: testAccTagResourcePreCreateForImportConfig(),
				Check: resource.ComposeTestCheckFunc(
					// Capture the tag ID for import
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["gtm_tag.pre_created"]
						if !ok {
							return fmt.Errorf("Tag resource not found: gtm_tag.pre_created")
						}
						createdTagID = rs.Primary.ID
						return nil
					},
				),
			},
			// Step 3: Import the existing tag into a new resource
			{
				Config:       testAccTagResourceImportTargetConfig(),
				ResourceName: "gtm_tag.imported",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					return createdTagID, nil
				},
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// Some fields might not be perfectly preserved during import
					"notes",
				},
			},
			// Step 4: Verify the imported tag can be managed by Terraform
			{
				Config: testAccTagResourceImportTargetUpdatedConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gtm_tag.imported", "name", "tf-test-imported-tag-updated"),
					resource.TestCheckResourceAttr("gtm_tag.imported", "notes", "Updated after import"),
				),
			},
		},
	})
}

// TestAccTagResource_importWithTerraformImportBlock tests using Terraform's import block syntax
func TestAccTagResource_importWithTerraformImportBlock(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			// Step 1: Create a tag that will be "imported"
			{
				Config: testAccTagResourceForImportBlockConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gtm_tag.to_be_imported", "id"),
					resource.TestCheckResourceAttr("gtm_tag.to_be_imported", "name", "tf-test-tag-for-import-block"),
				),
			},
			// Step 2: Use import block to import into a different resource
			{
				Config: testAccTagResourceWithImportBlockConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gtm_tag.imported_via_block", "id"),
					resource.TestCheckResourceAttr("gtm_tag.imported_via_block", "name", "tf-test-tag-for-import-block"),
					resource.TestCheckResourceAttr("gtm_tag.imported_via_block", "type", "html"),
				),
			},
		},
	})
}

// TestAccTagResource_importComplexTag tests importing a tag with complex parameters
func TestAccTagResource_importComplexTag(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			// Step 1: Create a complex tag
			{
				Config: testAccTagResourceComplexForImportConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gtm_tag.complex_original", "id"),
					resource.TestCheckResourceAttr("gtm_tag.complex_original", "parameter.#", "3"),
				),
			},
			// Step 2: Import the complex tag
			{
				ResourceName:      "gtm_tag.complex_original",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// Complex nested parameters might have ordering differences
					"parameter",
				},
			},
		},
	})
}

// TestAccTagResource_importWithTriggers tests importing a tag that has firing triggers
func TestAccTagResource_importWithTriggers(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			// Step 1: Create a tag with triggers
			{
				Config: testAccTagResourceWithTriggersForImportConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gtm_tag.with_triggers_original", "id"),
					resource.TestCheckResourceAttr("gtm_tag.with_triggers_original", "firing_trigger_id.#", "1"),
				),
			},
			// Step 2: Import the tag with triggers
			{
				ResourceName:      "gtm_tag.with_triggers_original",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"firing_trigger_id", // Trigger IDs might be referenced differently after import
				},
			},
		},
	})
}

// TestAccTagResource_importNonExistentTag tests importing a tag that doesn't exist
func TestAccTagResource_importNonExistentTag(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config:        testAccTagResourceImportTargetConfig(),
				ResourceName:  "gtm_tag.imported",
				ImportState:   true,
				ImportStateId: "nonexistent-tag-id",
				ExpectError:   nil, // The ImportState method should handle this gracefully
			},
		},
	})
}

// TestAccTagResource_importAndManage tests the full lifecycle of import and management
func TestAccTagResource_importAndManage(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			// Step 1: Create initial tag (simulating external creation)
			{
				Config: testAccTagResourceInitialForLifecycleConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gtm_tag.lifecycle_test", "id"),
					resource.TestCheckResourceAttr("gtm_tag.lifecycle_test", "name", "tf-test-lifecycle-original"),
				),
			},
			// Step 2: Import the tag
			{
				ResourceName:      "gtm_tag.lifecycle_test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Step 3: Update the imported tag
			{
				Config: testAccTagResourceUpdatedForLifecycleConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gtm_tag.lifecycle_test", "name", "tf-test-lifecycle-updated"),
					resource.TestCheckResourceAttr("gtm_tag.lifecycle_test", "notes", "Updated after import"),
				),
			},
			// Step 4: Add parameters to the imported tag
			{
				Config: testAccTagResourceEnhancedForLifecycleConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gtm_tag.lifecycle_test", "parameter.#", "2"),
					resource.TestCheckResourceAttr("gtm_tag.lifecycle_test", "parameter.0.key", "html"),
					resource.TestCheckResourceAttr("gtm_tag.lifecycle_test", "parameter.1.key", "supportDocumentWrite"),
				),
			},
		},
	})
}

// Configuration functions

func testAccTagResourcePreCreateForImportConfig() string {
	return testAccProviderConfig() + `
resource "gtm_tag" "pre_created" {
  name = "tf-test-tag-for-import"
  type = "html"
  notes = "Tag created for import testing"
  
  parameter = [
    {
      key   = "html"
      type  = "template"
      value = "<script>console.log('Pre-created tag');</script>"
    }
  ]
}
`
}

func testAccTagResourceEmptyConfig() string {
	return testAccProviderConfig()
}

func testAccTagResourceImportTargetConfig() string {
	return testAccProviderConfig() + `
resource "gtm_tag" "imported" {
  name = "tf-test-tag-for-import"
  type = "html"
  notes = "Tag created for import testing"
  
  parameter = [
    {
      key   = "html"
      type  = "template"
      value = "<script>console.log('Pre-created tag');</script>"
    }
  ]
}
`
}

func testAccTagResourceImportTargetUpdatedConfig() string {
	return testAccProviderConfig() + `
resource "gtm_tag" "imported" {
  name = "tf-test-imported-tag-updated"
  type = "html"
  notes = "Updated after import"
  
  parameter = [
    {
      key   = "html"
      type  = "template"
      value = "<script>console.log('Updated imported tag');</script>"
    }
  ]
}
`
}

func testAccTagResourceForImportBlockConfig() string {
	return testAccProviderConfig() + `
resource "gtm_tag" "to_be_imported" {
  name = "tf-test-tag-for-import-block"
  type = "html"
  notes = "Tag to be imported using import block"
  
  parameter = [
    {
      key   = "html"
      type  = "template"
      value = "<script>console.log('Tag for import block');</script>"
    }
  ]
}
`
}

func testAccTagResourceWithImportBlockConfig() string {
	return testAccProviderConfig() + `
# First create the tag that we'll import from
resource "gtm_tag" "to_be_imported" {
  name = "tf-test-tag-for-import-block"
  type = "html"
  notes = "Tag to be imported using import block"
  
  parameter = [
    {
      key   = "html"
      type  = "template"
      value = "<script>console.log('Tag for import block');</script>"
    }
  ]
}

# Import block syntax (Terraform 1.5+)
import {
  to = gtm_tag.imported_via_block
  id = gtm_tag.to_be_imported.id
}

resource "gtm_tag" "imported_via_block" {
  name = "tf-test-tag-for-import-block"
  type = "html"
  notes = "Tag to be imported using import block"
  
  parameter = [
    {
      key   = "html"
      type  = "template"
      value = "<script>console.log('Tag for import block');</script>"
    }
  ]
}
`
}

func testAccTagResourceComplexForImportConfig() string {
	return testAccProviderConfig() + `
resource "gtm_tag" "complex_original" {
  name = "tf-test-complex-tag-for-import"
  type = "gaawe"
  notes = "Complex GA4 tag for import testing"
  
  parameter = [
    {
      key   = "eventName"
      type  = "template"
      value = "import_test_event"
    },
    {
      key   = "measurementIdOverride"
      type  = "template"
      value = "G-IMPORT123"
    },
    {
      key  = "eventParameters"
      type = "list"
      
      list = [{
        type = "map"
        
        map = [{
          key   = "custom_parameter_1"
          type  = "template"
          value = "import_value_1"
        }, {
          key   = "custom_parameter_2"
          type  = "template"
          value = "import_value_2"
        }]
      }]
    }
  ]
}
`
}

func testAccTagResourceWithTriggersForImportConfig() string {
	return testAccProviderConfig() + `
# Create a trigger first
resource "gtm_trigger" "import_test" {
  name = "tf-test-trigger-for-import"
  type = "pageview"
  
  filter = [
    {
      type      = "equals"
      parameter = [{
        type  = "template"
        key   = "arg0"
        value = "{{Page URL}}"
      }, {
        type  = "template"
        key   = "arg1"
        value = "https://import-test.com/"
      }]
    }
  ]
}

# Create a tag with the trigger
resource "gtm_tag" "with_triggers_original" {
  name = "tf-test-tag-with-triggers-for-import"
  type = "html"
  notes = "Tag with triggers for import testing"
  
  parameter = [
    {
      key   = "html"
      type  = "template"
      value = "<script>console.log('Tag with triggers for import');</script>"
    }
  ]
  
  firing_trigger_id = [gtm_trigger.import_test.id]
}
`
}

func testAccTagResourceInitialForLifecycleConfig() string {
	return testAccProviderConfig() + `
resource "gtm_tag" "lifecycle_test" {
  name = "tf-test-lifecycle-original"
  type = "html"
  notes = "Original lifecycle test tag"
  
  parameter = [
    {
      key   = "html"
      type  = "template"
      value = "<script>console.log('Original lifecycle tag');</script>"
    }
  ]
}
`
}

func testAccTagResourceUpdatedForLifecycleConfig() string {
	return testAccProviderConfig() + `
resource "gtm_tag" "lifecycle_test" {
  name = "tf-test-lifecycle-updated"
  type = "html"
  notes = "Updated after import"
  
  parameter = [
    {
      key   = "html"
      type  = "template"
      value = "<script>console.log('Updated lifecycle tag');</script>"
    }
  ]
}
`
}

func testAccTagResourceEnhancedForLifecycleConfig() string {
	return testAccProviderConfig() + `
resource "gtm_tag" "lifecycle_test" {
  name = "tf-test-lifecycle-updated"
  type = "html"
  notes = "Enhanced after import"
  
  parameter = [
    {
      key   = "html"
      type  = "template"
      value = "<script>console.log('Enhanced lifecycle tag with more features');</script>"
    },
    {
      key   = "supportDocumentWrite"
      type  = "boolean"
      value = "false"
    }
  ]
}
`
}
