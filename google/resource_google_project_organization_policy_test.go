package google

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/api/cloudresourcemanager/v1"
)

func TestAccProjectOrganizationPolicy_boolean(t *testing.T) {
	t.Parallel()

	projectId := acctest.RandomWithPrefix("tf-test")

	org := getTestOrgFromEnv(t)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleProjectOrganizationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				// Test creation of an enforced boolean policy
				Config: testAccProjectOrganizationPolicy_boolean(org, projectId, true),
				Check:  testAccCheckGoogleProjectOrganizationBooleanPolicy("bool", true),
			},
			{
				// Test update from enforced to not
				Config: testAccProjectOrganizationPolicy_boolean(org, projectId, false),
				Check:  testAccCheckGoogleProjectOrganizationBooleanPolicy("bool", false),
			},
			{
				Config:  " ",
				Destroy: true,
			},
			{
				// Test creation of a not enforced boolean policy
				Config: testAccProjectOrganizationPolicy_boolean(org, projectId, false),
				Check:  testAccCheckGoogleProjectOrganizationBooleanPolicy("bool", false),
			},
			{
				// Test update from not enforced to enforced
				Config: testAccProjectOrganizationPolicy_boolean(org, projectId, true),
				Check:  testAccCheckGoogleProjectOrganizationBooleanPolicy("bool", true),
			},
		},
	})
}

func TestAccProjectOrganizationPolicy_list_allowAll(t *testing.T) {
	t.Parallel()

	projectId := acctest.RandomWithPrefix("tf-test")

	org := getTestOrgFromEnv(t)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleProjectOrganizationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectOrganizationPolicy_list_allowAll(org, projectId),
				Check:  testAccCheckGoogleProjectOrganizationListPolicyAll("list", "ALLOW"),
			},
		},
	})
}

func TestAccProjectOrganizationPolicy_list_allowSome(t *testing.T) {
	t.Parallel()

	projectId := acctest.RandomWithPrefix("tf-test")
	org := getTestOrgFromEnv(t)
	project := getTestProjectFromEnv()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleProjectOrganizationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectOrganizationPolicy_list_allowSome(org, projectId, project),
				Check:  testAccCheckGoogleProjectOrganizationListPolicyAllowedValues("list", []string{project}),
			},
		},
	})
}

func TestAccProjectOrganizationPolicy_list_denySome(t *testing.T) {
	t.Parallel()

	projectId := acctest.RandomWithPrefix("tf-test")
	org := getTestOrgFromEnv(t)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleProjectOrganizationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectOrganizationPolicy_list_denySome(org, projectId),
				Check:  testAccCheckGoogleProjectOrganizationListPolicyDeniedValues("list", DENIED_ORG_POLICIES),
			},
		},
	})
}

func TestAccProjectOrganizationPolicy_list_update(t *testing.T) {
	t.Parallel()

	projectId := acctest.RandomWithPrefix("tf-test")
	org := getTestOrgFromEnv(t)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleProjectOrganizationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectOrganizationPolicy_list_allowAll(org, projectId),
				Check:  testAccCheckGoogleProjectOrganizationListPolicyAll("list", "ALLOW"),
			},
			{
				Config: testAccProjectOrganizationPolicy_list_denySome(org, projectId),
				Check:  testAccCheckGoogleProjectOrganizationListPolicyDeniedValues("list", DENIED_ORG_POLICIES),
			},
		},
	})
}

func testAccCheckGoogleProjectOrganizationPolicyDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_project_organization_policy" {
			continue
		}

		projectId := rs.Primary.Attributes["projectId"]
		constraint := canonicalOrgPolicyConstraint(rs.Primary.Attributes["constraint"])
		policy, err := config.clientResourceManager.Projects.GetOrgPolicy(projectId, &cloudresourcemanager.GetOrgPolicyRequest{
			Constraint: constraint,
		}).Do()

		if err != nil {
			return err
		}

		if policy.ListPolicy != nil || policy.BooleanPolicy != nil {
			return fmt.Errorf("Org policy with constraint '%s' hasn't been cleared", constraint)
		}
	}
	return nil
}

func testAccCheckGoogleProjectOrganizationBooleanPolicy(n string, enforced bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		policy, err := getGoogleProjectOrganizationPolicyTestResource(s, n)
		if err != nil {
			return err
		}

		if policy.BooleanPolicy.Enforced != enforced {
			return fmt.Errorf("Expected boolean policy enforcement to be '%t', got '%t'", enforced, policy.BooleanPolicy.Enforced)
		}

		return nil
	}
}

func testAccCheckGoogleProjectOrganizationListPolicyAll(n, policyType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		policy, err := getGoogleProjectOrganizationPolicyTestResource(s, n)
		if err != nil {
			return err
		}

		if len(policy.ListPolicy.AllowedValues) > 0 || len(policy.ListPolicy.DeniedValues) > 0 {
			return fmt.Errorf("The `values` field shouldn't be set")
		}

		if policy.ListPolicy.AllValues != policyType {
			return fmt.Errorf("The list policy should %s all values", policyType)
		}

		return nil
	}
}

func testAccCheckGoogleProjectOrganizationListPolicyAllowedValues(n string, values []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		policy, err := getGoogleProjectOrganizationPolicyTestResource(s, n)
		if err != nil {
			return err
		}

		sort.Strings(policy.ListPolicy.AllowedValues)
		sort.Strings(values)
		if !reflect.DeepEqual(policy.ListPolicy.AllowedValues, values) {
			return fmt.Errorf("Expected the list policy to allow '%s', instead allowed '%s'", values, policy.ListPolicy.AllowedValues)
		}

		return nil
	}
}

func testAccCheckGoogleProjectOrganizationListPolicyDeniedValues(n string, values []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		policy, err := getGoogleProjectOrganizationPolicyTestResource(s, n)
		if err != nil {
			return err
		}

		sort.Strings(policy.ListPolicy.DeniedValues)
		sort.Strings(values)
		if !reflect.DeepEqual(policy.ListPolicy.DeniedValues, values) {
			return fmt.Errorf("Expected the list policy to deny '%s', instead denied '%s'", values, policy.ListPolicy.DeniedValues)
		}

		return nil
	}
}

func getGoogleProjectOrganizationPolicyTestResource(s *terraform.State, n string) (*cloudresourcemanager.OrgPolicy, error) {
	rn := "google_project_organization_policy." + n
	rs, ok := s.RootModule().Resources[rn]
	if !ok {
		return nil, fmt.Errorf("Not found: %s", rn)
	}

	if rs.Primary.ID == "" {
		return nil, fmt.Errorf("No ID is set")
	}

	config := testAccProvider.Meta().(*Config)
	//TODO: fix this to return an ACTUAL projectID from state
	projectId := rs.Primary.Attributes["projectId"]

	return config.clientResourceManager.Projects.GetOrgPolicy(projectId, &cloudresourcemanager.GetOrgPolicyRequest{
		Constraint: rs.Primary.Attributes["constraint"],
	}).Do()
}

func testAccProjectOrganizationPolicy_boolean(name, pid string, enforced bool) string {
	return fmt.Sprintf(`
resource "google_project" "orgpolicy" {
  name = "%s"
  project_id = "%s"
}

resource "google_project_organization_policy" "bool" {
  project    = "${google_project.orgpolicy.project_id}"
  constraint = "constraints/compute.disableSerialPortAccess"

  boolean_policy {
    enforced = %t
  }
}
`, name, pid, enforced)
}

func testAccProjectOrganizationPolicy_list_allowAll(name, pid string) string {
	return fmt.Sprintf(`
resource "google_project" "orgpolicy" {
  name = "%s"
  project_id = "%s"
}

resource "google_project_organization_policy" "list" {
  project    = "${google_project.orgpolicy.project_id}"
  constraint = "constraints/serviceuser.services"

  list_policy {
    allow {
      all = true
    }
  }
}
`, name, pid)
}

func testAccProjectOrganizationPolicy_list_allowSome(name, pid, project string) string {
	return fmt.Sprintf(`
resource "google_project" "orgpolicy" {
  name = "%s"
  project_id = "%s"
}

resource "google_project_organization_policy" "list" {
  project    = "${google_project.orgpolicy.project_id}"
  constraint = "constraints/compute.trustedImageProjects"

  list_policy {
    allow {
      values = ["%s"]
    }
  }
}
`, name, pid, project)
}

func testAccProjectOrganizationPolicy_list_denySome(name, pid string) string {
	return fmt.Sprintf(`
resource "google_project" "orgpolicy" {
  name = "%s"
  project_id = "%s"
}

resource "google_project_organization_policy" "list" {
  project    = "${google_project.orgpolicy.project_id}"
  constraint = "compute.vmExternalIpAccess"

  list_policy {
    deny {
      all = true
    }
  }
}
`, name, pid)
}
