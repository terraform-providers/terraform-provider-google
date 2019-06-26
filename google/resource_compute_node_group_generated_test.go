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
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccComputeNodeGroup_nodeGroupBasicExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(10),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeNodeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeNodeGroup_nodeGroupBasicExample(context),
			},
			{
				ResourceName:      "google_compute_node_group.nodes",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeNodeGroup_nodeGroupBasicExample(context map[string]interface{}) string {
	return Nprintf(`
data "google_compute_node_types" "central1a" {
  zone = "us-central1-a"
}

resource "google_compute_node_template" "soletenant-tmpl" {
  name = "soletenant-tmpl-%{random_suffix}"
  region = "us-central1"
  node_type = "${data.google_compute_node_types.central1a.names[0]}"
}

resource "google_compute_node_group" "nodes" {
  name = "soletenant-group-%{random_suffix}"
  zone = "us-central1-a"
  description = "example google_compute_node_group for Terraform Google Provider"

  size = 1
  node_template = "${google_compute_node_template.soletenant-tmpl.self_link}"
}
`, context)
}

func testAccCheckComputeNodeGroupDestroy(s *terraform.State) error {
	for name, rs := range s.RootModule().Resources {
		if rs.Type != "google_compute_node_group" {
			continue
		}
		if strings.HasPrefix(name, "data.") {
			continue
		}

		config := testAccProvider.Meta().(*Config)

		url, err := replaceVarsForTest(config, rs, "{{ComputeBasePath}}projects/{{project}}/zones/{{zone}}/nodeGroups/{{name}}")
		if err != nil {
			return err
		}

		_, err = sendRequest(config, "GET", url, nil)
		if err == nil {
			return fmt.Errorf("ComputeNodeGroup still exists at %s", url)
		}
	}

	return nil
}
