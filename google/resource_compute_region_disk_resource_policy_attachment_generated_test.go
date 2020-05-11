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

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccComputeRegionDiskResourcePolicyAttachment_regionDiskResourcePolicyAttachmentBasicExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeRegionDiskResourcePolicyAttachmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRegionDiskResourcePolicyAttachment_regionDiskResourcePolicyAttachmentBasicExample(context),
			},
			{
				ResourceName:            "google_compute_region_disk_resource_policy_attachment.attachment",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"disk", "region"},
			},
		},
	})
}

func testAccComputeRegionDiskResourcePolicyAttachment_regionDiskResourcePolicyAttachmentBasicExample(context map[string]interface{}) string {
	return Nprintf(`
resource "google_compute_region_disk_resource_policy_attachment" "attachment" {
  name = google_compute_resource_policy.policy.name
  disk = google_compute_region_disk.ssd.name
  region = "us-central1"
}

resource "google_compute_disk" "disk" {
  name  = "tf-test-my-base-disk%{random_suffix}"
  image = "debian-cloud/debian-9"
  size  = 50
  type  = "pd-ssd"
  zone  = "us-central1-a"
}

resource "google_compute_snapshot" "snapdisk" {
  name  = "tf-test-my-snapshot%{random_suffix}"
  source_disk = google_compute_disk.disk.name
  zone        = "us-central1-a"
}

resource "google_compute_region_disk" "ssd" {
  name  = "tf-test-my-disk%{random_suffix}"
  replica_zones = ["us-central1-a", "us-central1-f"]
  snapshot = google_compute_snapshot.snapdisk.id
  size  = 50
  type  = "pd-ssd"
  region  = "us-central1"
}

resource "google_compute_resource_policy" "policy" {
  name = "tf-test-my-resource-policy%{random_suffix}"
  region = "us-central1"
  snapshot_schedule_policy {
    schedule {
      daily_schedule {
        days_in_cycle = 1
        start_time = "04:00"
      }
    }
  }
}

data "google_compute_image" "my_image" {
  family  = "debian-9"
  project = "debian-cloud"
}
`, context)
}

func testAccCheckComputeRegionDiskResourcePolicyAttachmentDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_region_disk_resource_policy_attachment" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := googleProviderConfig(t)

			url, err := replaceVarsForTest(config, rs, "{{ComputeBasePath}}projects/{{project}}/regions/{{region}}/disks/{{disk}}")
			if err != nil {
				return err
			}

			_, err = sendRequest(config, "GET", "", url, nil)
			if err == nil {
				return fmt.Errorf("ComputeRegionDiskResourcePolicyAttachment still exists at %s", url)
			}
		}

		return nil
	}
}
