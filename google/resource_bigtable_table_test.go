package google

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccBigtableTable_basic(t *testing.T) {
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	tableName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableTable(instanceName, tableName),
				Check: resource.ComposeTestCheckFunc(
					testAccBigtableTableExists(
						"google_bigtable_table.table"),
				),
			},
		},
	})
}

func TestAccBigtableTable_splitKeys(t *testing.T) {
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	tableName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableTable_splitKeys(instanceName, tableName),
				Check: resource.ComposeTestCheckFunc(
					testAccBigtableTableExists(
						"google_bigtable_table.table"),
				),
			},
		},
	})
}

func TestAccBigtableTable_family(t *testing.T) {
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	tableName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	family := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableTable_family(instanceName, tableName, family),
				Check: resource.ComposeTestCheckFunc(
					testAccBigtableTableExists(
						"google_bigtable_table.table"),
				),
			},
		},
	})
}

func TestAccBigtableTable_familyMany(t *testing.T) {
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	tableName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	family := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableTable_familyMany(instanceName, tableName, family),
				Check: resource.ComposeTestCheckFunc(
					testAccBigtableTableExists(
						"google_bigtable_table.table"),
				),
			},
		},
	})
}

func testAccCheckBigtableTableDestroy(s *terraform.State) error {
	var ctx = context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_bigtable_table" {
			continue
		}

		config := testAccProvider.Meta().(*Config)
		c, err := config.bigtableClientFactory.NewAdminClient(config.Project, rs.Primary.Attributes["instance_name"])
		if err != nil {
			// The instance is already gone
			return nil
		}

		_, err = c.TableInfo(ctx, rs.Primary.Attributes["name"])
		if err == nil {
			return fmt.Errorf("table still present. Found %s in %s.", rs.Primary.Attributes["name"], rs.Primary.Attributes["instance_name"])
		}

		c.Close()
	}

	return nil
}

func testAccBigtableTableExists(n string) resource.TestCheckFunc {
	var ctx = context.Background()
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		config := testAccProvider.Meta().(*Config)
		c, err := config.bigtableClientFactory.NewAdminClient(config.Project, rs.Primary.Attributes["instance_name"])
		if err != nil {
			return fmt.Errorf("error starting admin client. %s", err)
		}

		_, err = c.TableInfo(ctx, rs.Primary.Attributes["name"])
		if err != nil {
			return fmt.Errorf("error retrieving table. Could not find %s in %s.", rs.Primary.Attributes["name"], rs.Primary.Attributes["instance_name"])
		}

		c.Close()

		return nil
	}
}

func testAccBigtableTable(instanceName, tableName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name          = "%s"
  instance_type = "DEVELOPMENT"
  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
  }
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = "${google_bigtable_instance.instance.name}"
}
`, instanceName, instanceName, tableName)
}

func testAccBigtableTable_splitKeys(instanceName, tableName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name          = "%s"
  instance_type = "DEVELOPMENT"
  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
  }
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = "${google_bigtable_instance.instance.name}"
  split_keys    = ["a", "b", "c"]
}
`, instanceName, instanceName, tableName)
}

func testAccBigtableTable_family(instanceName, tableName, family string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name          = "%s"

  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
  }

  instance_type = "DEVELOPMENT"
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = "${google_bigtable_instance.instance.name}"

  column_family {
    family = "%s"
  }
}
`, instanceName, instanceName, tableName, family)
}

func testAccBigtableTable_familyMany(instanceName, tableName, family string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name          = "%s"

  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
  }

  instance_type = "DEVELOPMENT"
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = "${google_bigtable_instance.instance.name}"

  column_family {
    family = "%s-first"
  }

  column_family {
    family = "%s-second"
  }
}
`, instanceName, instanceName, tableName, family, family)
}
