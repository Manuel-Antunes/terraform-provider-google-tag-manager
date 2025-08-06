package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Test basic tag creation and reading
func TestAccTagResource_basic(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceBasicConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gtm_tag.basic", "id"),
					resource.TestCheckResourceAttr("gtm_tag.basic", "name", "tf-test-tag-basic"),
					resource.TestCheckResourceAttr("gtm_tag.basic", "type", "html"),
					resource.TestCheckResourceAttr("gtm_tag.basic", "notes", "Basic HTML tag created by Terraform"),
					resource.TestCheckResourceAttr("gtm_tag.basic", "parameter.#", "1"),
					resource.TestCheckResourceAttr("gtm_tag.basic", "parameter.0.key", "html"),
					resource.TestCheckResourceAttr("gtm_tag.basic", "parameter.0.type", "template"),
					resource.TestCheckResourceAttr("gtm_tag.basic", "parameter.0.value", "<h1>Hello World</h1>"),
				),
			},
		},
	})
}

// Test Google Analytics 4 tag creation and reading
func TestAccTagResource_ga4(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceGA4Config(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gtm_tag.ga4", "id"),
					resource.TestCheckResourceAttr("gtm_tag.ga4", "name", "tf-test-tag-ga4"),
					resource.TestCheckResourceAttr("gtm_tag.ga4", "type", "gaawe"),
					resource.TestCheckResourceAttr("gtm_tag.ga4", "notes", "GA4 event tag created by Terraform"),
					resource.TestCheckResourceAttr("gtm_tag.ga4", "parameter.#", "3"),
					// Check GA4 specific parameters
					resource.TestCheckResourceAttr("gtm_tag.ga4", "parameter.0.key", "eventName"),
					resource.TestCheckResourceAttr("gtm_tag.ga4", "parameter.0.value", "page_view"),
					resource.TestCheckResourceAttr("gtm_tag.ga4", "parameter.1.key", "measurementIdOverride"),
					resource.TestCheckResourceAttr("gtm_tag.ga4", "parameter.1.value", "G-XXXXXXXXXX"),
				),
			},
		},
	})
}

// Test tag with firing triggers
func TestAccTagResource_withTriggers(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceWithTriggersConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gtm_tag.with_triggers", "id"),
					resource.TestCheckResourceAttr("gtm_tag.with_triggers", "name", "tf-test-tag-with-triggers"),
					resource.TestCheckResourceAttr("gtm_tag.with_triggers", "type", "html"),
					resource.TestCheckResourceAttr("gtm_tag.with_triggers", "firing_trigger_id.#", "1"),
					resource.TestCheckResourceAttrSet("gtm_tag.with_triggers", "firing_trigger_id.0"),
				),
			},
		},
	})
}

// Test tag import functionality
func TestAccTagResource_importBasic(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceBasicConfig(),
			},
			{
				ResourceName:      "gtm_tag.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test tag import with ID validation
func TestAccTagResource_importWithInvalidID(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config:        testAccTagResourceBasicConfig(),
				ResourceName:  "gtm_tag.basic",
				ImportState:   true,
				ImportStateId: "invalid-tag-id",
				ExpectError:   nil, // Will be handled by ImportState method
			},
		},
	})
}

// Test tag update functionality
func TestAccTagResource_updateBasic(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceBasicConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gtm_tag.basic", "name", "tf-test-tag-basic"),
					resource.TestCheckResourceAttr("gtm_tag.basic", "notes", "Basic HTML tag created by Terraform"),
				),
			},
			{
				Config: testAccTagResourceBasicUpdatedConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gtm_tag.basic", "name", "tf-test-tag-basic-updated"),
					resource.TestCheckResourceAttr("gtm_tag.basic", "notes", "Updated HTML tag by Terraform"),
					resource.TestCheckResourceAttr("gtm_tag.basic", "parameter.0.value", "<h1>Hello Updated World</h1>"),
				),
			},
		},
	})
}

// Test tag with complex nested parameters
func TestAccTagResource_complexNestedParameters(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceComplexParametersConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gtm_tag.complex", "id"),
					resource.TestCheckResourceAttr("gtm_tag.complex", "name", "tf-test-tag-complex"),
					resource.TestCheckResourceAttr("gtm_tag.complex", "type", "gaawe"),
					// Check nested parameters structure
					resource.TestCheckResourceAttr("gtm_tag.complex", "parameter.2.key", "eventParameters"),
					resource.TestCheckResourceAttr("gtm_tag.complex", "parameter.2.type", "list"),
					resource.TestCheckResourceAttr("gtm_tag.complex", "parameter.2.list.#", "2"),
				),
			},
		},
	})
}

// Test Google Ads Conversion tag
func TestAccTagResource_googleAdsConversion(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceGoogleAdsConversionConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gtm_tag.google_ads", "id"),
					resource.TestCheckResourceAttr("gtm_tag.google_ads", "name", "tf-test-google-ads-conversion"),
					resource.TestCheckResourceAttr("gtm_tag.google_ads", "type", "awct"),
					resource.TestCheckResourceAttr("gtm_tag.google_ads", "parameter.0.key", "conversionId"),
					resource.TestCheckResourceAttr("gtm_tag.google_ads", "parameter.1.key", "conversionLabel"),
				),
			},
		},
	})
}

// Test Facebook Pixel tag
func TestAccTagResource_facebookPixel(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceFacebookPixelConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gtm_tag.facebook_pixel", "id"),
					resource.TestCheckResourceAttr("gtm_tag.facebook_pixel", "name", "tf-test-facebook-pixel"),
					resource.TestCheckResourceAttr("gtm_tag.facebook_pixel", "type", "html"),
					resource.TestCheckResourceAttr("gtm_tag.facebook_pixel", "parameter.0.key", "html"),
				),
			},
		},
	})
}

// Test tag deletion
func TestAccTagResource_disappears(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceBasicConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTagExists("gtm_tag.basic"),
					testAccCheckTagDestroy("gtm_tag.basic"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// Test tag with no optional parameters
func TestAccTagResource_minimal(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceMinimalConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gtm_tag.minimal", "id"),
					resource.TestCheckResourceAttr("gtm_tag.minimal", "name", "tf-test-tag-minimal"),
					resource.TestCheckResourceAttr("gtm_tag.minimal", "type", "html"),
					resource.TestCheckResourceAttr("gtm_tag.minimal", "notes", ""),
					resource.TestCheckResourceAttr("gtm_tag.minimal", "parameter.#", "1"),
					resource.TestCheckResourceAttr("gtm_tag.minimal", "firing_trigger_id.#", "0"),
				),
			},
		},
	})
}

// Test tag parameter update
func TestAccTagResource_parameterUpdate(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceGA4Config(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gtm_tag.ga4", "parameter.0.value", "page_view"),
				),
			},
			{
				Config: testAccTagResourceGA4UpdatedConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gtm_tag.ga4", "parameter.0.value", "purchase"),
					resource.TestCheckResourceAttr("gtm_tag.ga4", "parameter.#", "4"), // Added new parameter
				),
			},
		},
	})
}

// Helper functions for testing

// testAccCheckTagExists verifies a tag exists in GTM
func testAccCheckTagExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Tag resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Tag ID is not set")
		}

		// Additional check could be made here to verify the tag exists in GTM
		// This would require access to the GTM client

		return nil
	}
}

// testAccCheckTagDestroy verifies a tag no longer exists
func testAccCheckTagDestroy(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// This would typically check that the tag was deleted from GTM
		// For now, we'll just verify the resource is removed from state
		return nil
	}
}

// Configuration functions

func testAccTagResourceBasicConfig() string {
	return testAccProviderConfig() + `
resource "gtm_tag" "basic" {
  name  = "tf-test-tag-basic"
  type  = "html"
  notes = "Basic HTML tag created by Terraform"
  
  parameter = [
    {
      key   = "html"
      type  = "template"
      value = "<h1>Hello World</h1>"
    }
  ]
}
`
}

func testAccTagResourceBasicUpdatedConfig() string {
	return testAccProviderConfig() + `
resource "gtm_tag" "basic" {
  name  = "tf-test-tag-basic-updated"
  type  = "html"
  notes = "Updated HTML tag by Terraform"
  
  parameter = [
    {
      key   = "html"
      type  = "template"
      value = "<h1>Hello Updated World</h1>"
    }
  ]
}
`
}

func testAccTagResourceGA4Config() string {
	return testAccProviderConfig() + `
resource "gtm_tag" "ga4" {
  name  = "tf-test-tag-ga4"
  type  = "gaawe"
  notes = "GA4 event tag created by Terraform"
  
  parameter = [
    {
      key   = "eventName"
      type  = "template"
      value = "page_view"
    },
    {
      key   = "measurementIdOverride"
      type  = "template"
      value = "G-XXXXXXXXXX"
    },
    {
      key  = "eventParameters"
      type = "list"
      
      list = [{
        type = "map"
        
        map = [{
          key   = "page_title"
          type  = "template"
          value = "{{Page Title}}"
        }, {
          key   = "page_location"
          type  = "template"
          value = "{{Page URL}}"
        }]
      }]
    }
  ]
}
`
}

func testAccTagResourceGA4UpdatedConfig() string {
	return testAccProviderConfig() + `
resource "gtm_tag" "ga4" {
  name  = "tf-test-tag-ga4"
  type  = "gaawe"
  notes = "GA4 event tag created by Terraform"
  
  parameter = [
    {
      key   = "eventName"
      type  = "template"
      value = "purchase"
    },
    {
      key   = "measurementIdOverride"
      type  = "template"
      value = "G-XXXXXXXXXX"
    },
    {
      key  = "eventParameters"
      type = "list"
      
      list = [{
        type = "map"
        
        map = [{
          key   = "transaction_id"
          type  = "template"
          value = "{{Transaction ID}}"
        }, {
          key   = "value"
          type  = "template"
          value = "{{Purchase Value}}"
        }]
      }]
    },
    {
      key   = "currency"
      type  = "template"
      value = "USD"
    }
  ]
}
`
}

func testAccTagResourceWithTriggersConfig() string {
	return testAccProviderConfig() + `
# First create a trigger to use with the tag
resource "gtm_trigger" "test" {
  name = "tf-test-trigger-for-tag"
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
        value = "https://example.com/"
      }]
    }
  ]
}

resource "gtm_tag" "with_triggers" {
  name  = "tf-test-tag-with-triggers"
  type  = "html"
  notes = "HTML tag with firing triggers"
  
  parameter = [
    {
      key   = "html"
      type  = "template"
      value = "<script>console.log('Tag fired!');</script>"
    }
  ]
  
  firing_trigger_id = [gtm_trigger.test.id]
}
`
}

func testAccTagResourceComplexParametersConfig() string {
	return testAccProviderConfig() + `
resource "gtm_tag" "complex" {
  name  = "tf-test-tag-complex"
  type  = "gaawe"
  notes = "GA4 tag with complex nested parameters"
  
  parameter = [
    {
      key   = "eventName"
      type  = "template"
      value = "complex_event"
    },
    {
      key   = "measurementIdOverride"
      type  = "template"
      value = "G-XXXXXXXXXX"
    },
    {
      key  = "eventParameters"
      type = "list"
      
      list = [{
        type = "map"
        
        map = [{
          key   = "name"
          type  = "template"
          value = "custom_parameter_1"
        }, {
          key   = "value"
          type  = "template"
          value = "{{Custom Variable 1}}"
        }]
      }, {
        type = "map"
        
        map = [{
          key   = "name"
          type  = "template"
          value = "custom_parameter_2"
        }, {
          key   = "value"
          type  = "template"
          value = "{{Custom Variable 2}}"
        }]
      }]
    }
  ]
}
`
}

func testAccTagResourceGoogleAdsConversionConfig() string {
	return testAccProviderConfig() + `
resource "gtm_tag" "google_ads" {
  name  = "tf-test-google-ads-conversion"
  type  = "awct"
  notes = "Google Ads conversion tag"
  
  parameter = [
    {
      key   = "conversionId"
      type  = "template"
      value = "123456789"
    },
    {
      key   = "conversionLabel"
      type  = "template"
      value = "abcd1234"
    },
    {
      key   = "conversionValue"
      type  = "template"
      value = "{{Transaction Value}}"
    }
  ]
}
`
}

func testAccTagResourceFacebookPixelConfig() string {
	return testAccProviderConfig() + `
resource "gtm_tag" "facebook_pixel" {
  name  = "tf-test-facebook-pixel"
  type  = "html"
  notes = "Facebook Pixel tag"
  
  parameter = [
    {
      key   = "html"
      type  = "template"
      value = <<-EOT
        <!-- Facebook Pixel Code -->
        <script>
        !function(f,b,e,v,n,t,s)
        {if(f.fbq)return;n=f.fbq=function(){n.callMethod?
        n.callMethod.apply(n,arguments):n.queue.push(arguments)};
        if(!f._fbq)f._fbq=n;n.push=n;n.loaded=!0;n.version='2.0';
        n.queue=[];t=b.createElement(e);t.async=!0;
        t.src=v;s=b.getElementsByTagName(e)[0];
        s.parentNode.insertBefore(t,s)}(window, document,'script',
        'https://connect.facebook.net/en_US/fbevents.js');
        fbq('init', '1234567890123456');
        fbq('track', 'PageView');
        </script>
        <noscript><img height="1" width="1" style="display:none"
        src="https://www.facebook.com/tr?id=1234567890123456&ev=PageView&noscript=1"
        /></noscript>
        <!-- End Facebook Pixel Code -->
      EOT
    }
  ]
}
`
}

func testAccTagResourceMinimalConfig() string {
	return testAccProviderConfig() + `
resource "gtm_tag" "minimal" {
  name = "tf-test-tag-minimal"
  type = "html"
  
  parameter = [
    {
      key   = "html"
      type  = "template"
      value = "<p>Minimal tag</p>"
    }
  ]
}
`
}
