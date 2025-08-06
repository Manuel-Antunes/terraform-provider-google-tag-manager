package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccTagResource_universalAnalytics tests Universal Analytics tag
func TestAccTagResource_universalAnalytics(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceUniversalAnalyticsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gtm_tag.ua", "id"),
					resource.TestCheckResourceAttr("gtm_tag.ua", "name", "tf-test-ua-tag"),
					resource.TestCheckResourceAttr("gtm_tag.ua", "type", "ua"),
					resource.TestCheckResourceAttr("gtm_tag.ua", "parameter.0.key", "trackingId"),
				),
			},
		},
	})
}

// TestAccTagResource_customHTML tests custom HTML tag
func TestAccTagResource_customHTML(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceCustomHTMLConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gtm_tag.custom_html", "id"),
					resource.TestCheckResourceAttr("gtm_tag.custom_html", "name", "tf-test-custom-html"),
					resource.TestCheckResourceAttr("gtm_tag.custom_html", "type", "html"),
				),
			},
		},
	})
}

// TestAccTagResource_customImage tests custom image tag
func TestAccTagResource_customImage(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceCustomImageConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gtm_tag.custom_image", "id"),
					resource.TestCheckResourceAttr("gtm_tag.custom_image", "name", "tf-test-custom-image"),
					resource.TestCheckResourceAttr("gtm_tag.custom_image", "type", "img"),
					resource.TestCheckResourceAttr("gtm_tag.custom_image", "parameter.0.key", "url"),
				),
			},
		},
	})
}

// TestAccTagResource_linkedinInsightTag tests LinkedIn Insight tag
func TestAccTagResource_linkedinInsightTag(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceLinkedInInsightTagConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gtm_tag.linkedin", "id"),
					resource.TestCheckResourceAttr("gtm_tag.linkedin", "name", "tf-test-linkedin-insight"),
					resource.TestCheckResourceAttr("gtm_tag.linkedin", "type", "html"),
				),
			},
		},
	})
}

// TestAccTagResource_twitterUniversalWebsiteTag tests Twitter Universal Website tag
func TestAccTagResource_twitterUniversalWebsiteTag(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceTwitterUWTConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gtm_tag.twitter", "id"),
					resource.TestCheckResourceAttr("gtm_tag.twitter", "name", "tf-test-twitter-uwt"),
					resource.TestCheckResourceAttr("gtm_tag.twitter", "type", "html"),
				),
			},
		},
	})
}

// TestAccTagResource_hotjarTracking tests Hotjar tracking tag
func TestAccTagResource_hotjarTracking(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceHotjarConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gtm_tag.hotjar", "id"),
					resource.TestCheckResourceAttr("gtm_tag.hotjar", "name", "tf-test-hotjar"),
					resource.TestCheckResourceAttr("gtm_tag.hotjar", "type", "html"),
				),
			},
		},
	})
}

// TestAccTagResource_googleOptimize tests Google Optimize tag
func TestAccTagResource_googleOptimize(t *testing.T) {
	testAccPreCheck(t)
	ctx := Context(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(ctx, ProviderNameEcho),
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceGoogleOptimizeConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gtm_tag.optimize", "id"),
					resource.TestCheckResourceAttr("gtm_tag.optimize", "name", "tf-test-google-optimize"),
					resource.TestCheckResourceAttr("gtm_tag.optimize", "type", "html"),
				),
			},
		},
	})
}

// Configuration functions for different tag types

func testAccTagResourceUniversalAnalyticsConfig() string {
	return testAccProviderConfig() + `
resource "gtm_tag" "ua" {
  name = "tf-test-ua-tag"
  type = "ua"
  notes = "Universal Analytics pageview tag"
  
  parameter = [
    {
      key   = "trackingId"
      type  = "template"
      value = "UA-12345678-1"
    },
    {
      key   = "trackType"
      type  = "template"
      value = "TRACK_PAGEVIEW"
    }
  ]
}
`
}

func testAccTagResourceCustomHTMLConfig() string {
	return testAccProviderConfig() + `
resource "gtm_tag" "custom_html" {
  name = "tf-test-custom-html"
  type = "html"
  notes = "Custom HTML tag with JavaScript"
  
  parameter = [
    {
      key   = "html"
      type  = "template"
      value = <<-EOT
        <script>
          (function() {
            console.log('Custom HTML tag fired');
            
            // Custom tracking logic
            if (typeof dataLayer !== 'undefined') {
              dataLayer.push({
                'event': 'custom_html_tag_fired',
                'tag_name': 'tf-test-custom-html'
              });
            }
            
            // Custom analytics
            var customData = {
              timestamp: new Date().toISOString(),
              userAgent: navigator.userAgent,
              url: window.location.href
            };
            
            console.log('Custom data:', customData);
          })();
        </script>
      EOT
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

func testAccTagResourceCustomImageConfig() string {
	return testAccProviderConfig() + `
resource "gtm_tag" "custom_image" {
  name = "tf-test-custom-image"
  type = "img"
  notes = "Custom image tag for tracking pixels"
  
  parameter = [
    {
      key   = "url"
      type  = "template"
      value = "https://example.com/tracking-pixel.png?utm_source={{utm_source}}&utm_medium={{utm_medium}}"
    },
    {
      key   = "cacheBusterQueryParam"
      type  = "template"
      value = "cb"
    }
  ]
}
`
}

func testAccTagResourceLinkedInInsightTagConfig() string {
	return testAccProviderConfig() + `
resource "gtm_tag" "linkedin" {
  name = "tf-test-linkedin-insight"
  type = "html"
  notes = "LinkedIn Insight Tag for conversion tracking"
  
  parameter = [
    {
      key   = "html"
      type  = "template"
      value = <<-EOT
        <script type="text/javascript">
        _linkedin_partner_id = "123456";
        window._linkedin_data_partner_ids = window._linkedin_data_partner_ids || [];
        window._linkedin_data_partner_ids.push(_linkedin_partner_id);
        </script><script type="text/javascript">
        (function(){var s = document.getElementsByTagName("script")[0];
        var b = document.createElement("script");
        b.type = "text/javascript";b.async = true;
        b.src = "https://snap.licdn.com/li.lms-analytics/insight.min.js";
        s.parentNode.insertBefore(b, s);})();
        </script>
        <noscript>
        <img height="1" width="1" style="display:none;" alt="" src="https://px.ads.linkedin.com/collect/?pid=123456&fmt=gif" />
        </noscript>
      EOT
    }
  ]
}
`
}

func testAccTagResourceTwitterUWTConfig() string {
	return testAccProviderConfig() + `
resource "gtm_tag" "twitter" {
  name = "tf-test-twitter-uwt"
  type = "html"
  notes = "Twitter Universal Website Tag"
  
  parameter = [
    {
      key   = "html"
      type  = "template"
      value = <<-EOT
        <script>
        !function(e,t,n,s,u,a){e.twq||(s=e.twq=function(){s.exe?s.exe.apply(s,arguments):s.queue.push(arguments);
        },s.version='1.1',s.queue=[],u=t.createElement(n),u.async=!0,u.src='//static.ads-twitter.com/uwt.js',
        a=t.getElementsByTagName(n)[0],a.parentNode.insertBefore(u,a))}(window,document,'script');
        twq('init','o1234');
        twq('track','PageView');
        </script>
      EOT
    }
  ]
}
`
}

func testAccTagResourceHotjarConfig() string {
	return testAccProviderConfig() + `
resource "gtm_tag" "hotjar" {
  name = "tf-test-hotjar"
  type = "html"
  notes = "Hotjar Tracking Code"
  
  parameter = [
    {
      key   = "html"
      type  = "template"
      value = <<-EOT
        <script>
            (function(h,o,t,j,a,r){
                h.hj=h.hj||function(){(h.hj.q=h.hj.q||[]).push(arguments)};
                h._hjSettings={hjid:123456,hjsv:6};
                a=o.getElementsByTagName('head')[0];
                r=o.createElement('script');r.async=1;
                r.src=t+h._hjSettings.hjid+j+h._hjSettings.hjsv;
                a.appendChild(r);
            })(window,document,'https://static.hotjar.com/c/hotjar-','.js?sv=');
        </script>
      EOT
    }
  ]
}
`
}

func testAccTagResourceGoogleOptimizeConfig() string {
	return testAccProviderConfig() + `
resource "gtm_tag" "optimize" {
  name = "tf-test-google-optimize"
  type = "html"
  notes = "Google Optimize Tag"
  
  parameter = [
    {
      key   = "html"
      type  = "template"
      value = <<-EOT
        <script src="https://www.googleoptimize.com/optimize.js?id=GTM-XXXXXXX"></script>
        <script>
          gtag('config', 'GA_TRACKING_ID', {
            optimize_id: 'GTM-XXXXXXX'
          });
        </script>
      EOT
    }
  ]
}
`
}
