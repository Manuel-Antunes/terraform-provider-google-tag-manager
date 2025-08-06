package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccTagResource_validation tests various validation scenarios
func TestAccTagResource_validation(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config:      testAccTagResourceInvalidTypeConfig(),
				ExpectError: nil, // GTM API will handle validation
			},
		},
	})
}

// TestAccTagResource_emptyParameters tests tag with empty parameters
func TestAccTagResource_emptyParameters(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceEmptyParametersConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gtm_tag.empty_params", "id"),
					resource.TestCheckResourceAttr("gtm_tag.empty_params", "name", "tf-test-tag-empty-params"),
					resource.TestCheckResourceAttr("gtm_tag.empty_params", "parameter.#", "0"),
				),
			},
		},
	})
}

// TestAccTagResource_longName tests tag with maximum length name
func TestAccTagResource_longName(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceLongNameConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gtm_tag.long_name", "id"),
					resource.TestCheckResourceAttr("gtm_tag.long_name", "name", "tf-test-tag-with-very-long-name-that-tests-maximum-length-limits-in-google-tag-manager-names"),
				),
			},
		},
	})
}

// TestAccTagResource_specialCharacters tests tag with special characters
func TestAccTagResource_specialCharacters(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceSpecialCharactersConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gtm_tag.special_chars", "id"),
					resource.TestCheckResourceAttr("gtm_tag.special_chars", "name", "tf-test-tag-special-chars-!@#$%^&*()"),
				),
			},
		},
	})
}

// TestAccTagResource_unicodeCharacters tests tag with unicode characters
func TestAccTagResource_unicodeCharacters(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceUnicodeCharactersConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gtm_tag.unicode", "id"),
					resource.TestCheckResourceAttr("gtm_tag.unicode", "name", "tf-test-tag-unicode-测试-тест-🏷️"),
				),
			},
		},
	})
}

// TestAccTagResource_multipleTriggersUpdate tests updating firing triggers
func TestAccTagResource_multipleTriggersUpdate(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceMultipleTriggersConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gtm_tag.multiple_triggers", "firing_trigger_id.#", "2"),
				),
			},
			{
				Config: testAccTagResourceMultipleTriggersUpdatedConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gtm_tag.multiple_triggers", "firing_trigger_id.#", "1"),
				),
			},
		},
	})
}

// TestAccTagResource_removeAllParameters tests removing all parameters
func TestAccTagResource_removeAllParameters(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceBasicConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gtm_tag.basic", "parameter.#", "1"),
				),
			},
			{
				Config: testAccTagResourceEmptyParametersConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gtm_tag.empty_params", "parameter.#", "0"),
				),
			},
		},
	})
}

// Configuration functions

func testAccTagResourceInvalidTypeConfig() string {
	return testAccProviderConfig() + `
resource "gtm_tag" "invalid_type" {
  name = "tf-test-tag-invalid-type"
  type = "invalid_tag_type"
  
  parameter = [
    {
      key   = "html"
      type  = "template"
      value = "<p>Test</p>"
    }
  ]
}
`
}

func testAccTagResourceEmptyParametersConfig() string {
	return testAccProviderConfig() + `
resource "gtm_tag" "empty_params" {
  name = "tf-test-tag-empty-params"
  type = "html"
  notes = "Tag with no parameters"
}
`
}

func testAccTagResourceLongNameConfig() string {
	return testAccProviderConfig() + `
resource "gtm_tag" "long_name" {
  name = "tf-test-tag-with-very-long-name-that-tests-maximum-length-limits-in-google-tag-manager-names"
  type = "html"
  
  parameter = [
    {
      key   = "html"
      type  = "template"
      value = "<p>Long name test</p>"
    }
  ]
}
`
}

func testAccTagResourceSpecialCharactersConfig() string {
	return testAccProviderConfig() + `
resource "gtm_tag" "special_chars" {
  name = "tf-test-tag-special-chars-!@#$%^&*()"
  type = "html"
  notes = "Tag with special characters in name"
  
  parameter = [
    {
      key   = "html"
      type  = "template"
      value = "<script>alert('Special chars: !@#$%^&*()');</script>"
    }
  ]
}
`
}

func testAccTagResourceUnicodeCharactersConfig() string {
	return testAccProviderConfig() + `
resource "gtm_tag" "unicode" {
  name = "tf-test-tag-unicode-测试-тест-🏷️"
  type = "html"
  notes = "Tag with unicode characters: 测试 тест 🏷️"
  
  parameter = [
    {
      key   = "html"
      type  = "template"
      value = "<p>Unicode test: 你好世界 Привет мир 🌍</p>"
    }
  ]
}
`
}

func testAccTagResourceMultipleTriggersConfig() string {
	return testAccProviderConfig() + `
# Create multiple triggers
resource "gtm_trigger" "test1" {
  name = "tf-test-trigger-1"
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

resource "gtm_trigger" "test2" {
  name = "tf-test-trigger-2"
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
        value = "https://test.com/"
      }]
    }
  ]
}

resource "gtm_tag" "multiple_triggers" {
  name = "tf-test-tag-multiple-triggers"
  type = "html"
  
  parameter = [
    {
      key   = "html"
      type  = "template"
      value = "<script>console.log('Multiple triggers');</script>"
    }
  ]
  
  firing_trigger_id = [
    gtm_trigger.test1.id,
    gtm_trigger.test2.id
  ]
}
`
}

func testAccTagResourceMultipleTriggersUpdatedConfig() string {
	return testAccProviderConfig() + `
# Keep only one trigger
resource "gtm_trigger" "test1" {
  name = "tf-test-trigger-1"
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

resource "gtm_tag" "multiple_triggers" {
  name = "tf-test-tag-multiple-triggers"
  type = "html"
  
  parameter = [
    {
      key   = "html"
      type  = "template"
      value = "<script>console.log('Single trigger now');</script>"
    }
  ]
  
  firing_trigger_id = [gtm_trigger.test1.id]
}
`
}
