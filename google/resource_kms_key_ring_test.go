package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"log"
)

func TestAccGoogleKmsKeyRing_basic(t *testing.T) {
	skipIfEnvNotSet(t,
		[]string{
			"GOOGLE_ORG",
			"GOOGLE_BILLING_ACCOUNT",
		}...,
	)

	projectId := "terraform-" + acctest.RandString(10)
	projectOrg := multiEnvSearch([]string{
		"GOOGLE_ORG",
	})
	projectBillingAccount := multiEnvSearch([]string{
		"GOOGLE_BILLING_ACCOUNT",
	})
	keyRingName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testGoogleKmsKeyRing_recreate(keyRingName),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testGoogleKmsKeyRing_basic(projectId, projectOrg, projectBillingAccount, keyRingName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleKmsKeyRingExists("google_kms_key_ring.key_ring"),
				),
			},
		},
	})
}

func testAccCheckGoogleKmsKeyRingExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		location := rs.Primary.Attributes["location"]
		project := rs.Primary.Attributes["project"]

		parent := kmsResourceParentString(project, location)
		keyRingName := kmsResourceParentKeyRingName(project, location, name)

		listKeyRingsResponse, err := config.clientKms.Projects.Locations.KeyRings.List(parent).Do()
		if err != nil {
			return fmt.Errorf("Error listing KeyRings: %s", err)
		}

		for _, keyRing := range listKeyRingsResponse.KeyRings {
			log.Printf("[DEBUG] Found KeyRing: %s", keyRing.Name)

			if keyRing.Name == keyRingName {
				return nil
			}
		}

		return fmt.Errorf("KeyRing not found: %s", keyRingName)
	}
}

// TODO
// KMS KeyRings cannot be deleted. This will test if the resource can be added back to state after being removed
func testGoogleKmsKeyRing_recreate(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return nil
	}
}

func testGoogleKmsKeyRing_basic(projectId, projectOrg, projectBillingAccount, keyRingName string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
    name            = "%s"
	project_id      = "%s"
    org_id          = "%s"
	billing_account = "%s"
}

resource "google_project_services" "acceptance" {
    project  = "${google_project.acceptance.project_id}"
    services = [
        "cloudkms.googleapis.com"
    ]
}

resource "google_kms_key_ring" "key_ring" {
	project  = "${google_project_services.acceptance.project}"
	name     = "%s"
	location = "us-central1"
}
	`, projectId, projectId, projectOrg, projectBillingAccount, keyRingName)
}
