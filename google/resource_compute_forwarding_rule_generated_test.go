// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccComputeForwardingRule_forwardingRuleBasicExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(10),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeForwardingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeForwardingRule_forwardingRuleBasicExample(context),
			},
			{
				ResourceName:      "google_compute_forwarding_rule.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeForwardingRule_forwardingRuleBasicExample(context map[string]interface{}) string {
	return Nprintf(`
resource "google_compute_forwarding_rule" "default" {
  name       = "website-forwarding-rule-%{random_suffix}"
  target     = "${google_compute_target_pool.default.self_link}"
  port_range = "80"
}

resource "google_compute_target_pool" "default" {
  name = "website-target-pool-%{random_suffix}"
}
`, context)
}

func testAccCheckComputeForwardingRuleDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_compute_forwarding_rule" {
			continue
		}

		config := testAccProvider.Meta().(*Config)

		url, err := replaceVarsForTest(rs, "https://www.googleapis.com/compute/v1/projects/{{project}}/regions/{{region}}/forwardingRules/{{name}}")
		if err != nil {
			return err
		}

		_, err = sendRequest(config, "GET", url, nil)
		if err == nil {
			return fmt.Errorf("ComputeForwardingRule still exists at %s", url)
		}
	}

	return nil
}
