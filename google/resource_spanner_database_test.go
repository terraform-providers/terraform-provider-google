package google

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/spanner/v1"
)

// Unit Tests

func TestDatabaseNameForApi(t *testing.T) {
	actual := databaseNameForApi("project123", "instance456", "db789")
	expected := "projects/project123/instances/instance456/databases/db789"
	expectEquals(t, expected, actual)
}

func TestExtractSpannerDBImportIds(t *testing.T) {
	sid, e := extractSpannerDatabaseImportIds("instance456/database789")
	if e != nil {
		t.Errorf("Error should have been nil")
	}
	expectEquals(t, "", sid.Project)
	expectEquals(t, "instance456", sid.Instance)
	expectEquals(t, "database789", sid.Database)
}

func TestExtractSpannerDBImportIds_withProject(t *testing.T) {
	sid, e := extractSpannerDatabaseImportIds("project123/instance456/database789")
	if e != nil {
		t.Errorf("Error should have been nil")
	}
	expectEquals(t, "project123", sid.Project)
	expectEquals(t, "instance456", sid.Instance)
	expectEquals(t, "database789", sid.Database)
}

func TestExtractSpannerDBImportIds_invalidLeadingSlash(t *testing.T) {
	sid, e := extractSpannerDatabaseImportIds("/instance456/database789")
	expectInvalidSpannerDbImportId(t, sid, e)
}

func TestExtractSpannerDBImportIds_invalidTrailingSlash(t *testing.T) {
	sid, e := extractSpannerDatabaseImportIds("instance456/database789/")
	expectInvalidSpannerDbImportId(t, sid, e)
}

func TestExtractSpannerDBImportIds_invalidSingleSlash(t *testing.T) {
	sid, e := extractSpannerDatabaseImportIds("/")
	expectInvalidSpannerDbImportId(t, sid, e)
}

func TestExtractSpannerDBImportIds_invalidMultiSlash(t *testing.T) {
	sid, e := extractSpannerDatabaseImportIds("project123/instance456/db789/next")
	expectInvalidSpannerDbImportId(t, sid, e)
}

func expectInvalidSpannerDbImportId(t *testing.T, sid *spannerDbImportId, e error) {
	if sid != nil {
		t.Errorf("Expected spannerDbImportId to be nil")
		return
	}
	if e == nil {
		t.Errorf("Expected an Error but did not get one")
		return
	}
	if !strings.HasPrefix(e.Error(), "Invalid spanner database specifier") {
		t.Errorf("Expecting Error starting with 'Invalid spanner database specifier'")
	}
}

// Acceptance Tests

func TestAccSpannerDatabase_basic(t *testing.T) {
	var db spanner.Database
	rnd := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckSpannerInstanceDestroy,
			testAccCheckSpannerDatabaseDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerDatabase_basic(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpannerDatabaseExists("google_spanner_database.basic", &db),

					resource.TestCheckResourceAttr("google_spanner_database.basic", "name", "my-db-"+rnd),
				),
			},
		},
	})
}

func TestAccSpannerDatabase_basicWithInitialDDL(t *testing.T) {
	var db spanner.Database
	rnd := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckSpannerInstanceDestroy,
			testAccCheckSpannerDatabaseDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerDatabase_basicWithInitialDDL(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpannerDatabaseExists("google_spanner_database.basic", &db),
				),
			},
		},
	})
}

func TestAccSpannerDatabase_duplicateNameError(t *testing.T) {
	var db spanner.Database
	rnd := acctest.RandString(10)
	dbName := fmt.Sprintf("spanner-test-%s", rnd)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckSpannerInstanceDestroy,
			testAccCheckSpannerDatabaseDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerDatabase_duplicateNameError_part1(rnd, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpannerDatabaseExists("google_spanner_database.basic1", &db),
				),
			},
			{
				Config: testAccSpannerDatabase_duplicateNameError_part2(rnd, dbName),
				ExpectError: regexp.MustCompile(
					fmt.Sprintf(".*A database with name %s already exists", dbName)),
			},
		},
	})
}

func testAccCheckSpannerDatabaseDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_spanner_database" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Unable to verify delete of spanner database, ID is empty")
		}

		project, err := getProjectId(rs, config)
		if err != nil {
			return err
		}

		dbName := rs.Primary.Attributes["name"]
		instanceName := rs.Primary.Attributes["instance"]
		_, err = config.clientSpanner.Projects.Instances.Databases.Get(
			databaseNameForApi(project, instanceName, dbName)).Do()

		if err != nil {
			if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == http.StatusNotFound {
				return nil
			}
			return fmt.Errorf("Error make GCP platform call to verify spanner database deleted: %s", err.Error())
		}
		return fmt.Errorf("Spanner database not destroyed - still exists")
	}

	return nil
}

func testAccCheckSpannerDatabaseExists(n string, instance *spanner.Database) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Terraform resource Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set for Spanner instance")
		}

		project, err := getProjectId(rs, config)
		if err != nil {
			return err
		}

		dbName := rs.Primary.Attributes["name"]
		instanceName := rs.Primary.Attributes["instance"]
		found, err := config.clientSpanner.Projects.Instances.Databases.Get(
			databaseNameForApi(project, instanceName, dbName)).Do()
		if err != nil {
			return err
		}

		fName := extractInstanceNameFromApi(found.Name)
		if fName != rs.Primary.ID {
			return fmt.Errorf("Spanner database %s not found, found %s instead", rs.Primary.ID, fName)
		}

		*instance = *found

		return nil
	}
}

func testAccSpannerDatabase_basic(rnd string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  name          = "my-instance-%s"
  config        = "regional-us-central1"
  display_name  = "my-displayname-%s"
  num_nodes     = 1
}

resource "google_spanner_database" "basic" {
  instance      = "${google_spanner_instance.basic.name}"
  name          = "my-db-%s"

}
`, rnd, rnd, rnd)
}

func testAccSpannerDatabase_basicWithInitialDDL(rnd string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  name          = "my-instance-%s"
  config        = "regional-us-central1"
  display_name  = "my-displayname-%s"
  num_nodes     = 1
}

resource "google_spanner_database" "basic" {
  instance      = "${google_spanner_instance.basic.name}"
  name          = "my-db-%s"
  ddl           =  [
     "CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)",
     "CREATE TABLE t2 (t2 INT64 NOT NULL,) PRIMARY KEY(t2)" ]
}
`, rnd, rnd, rnd)
}

func testAccSpannerDatabase_duplicateNameError_part1(rnd, dbName string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  name          = "my-instance-%s"
  config        = "regional-us-central1"
  display_name  = "my-displayname-%s"
  num_nodes     = 1
}

resource "google_spanner_database" "basic1" {
  instance      = "${google_spanner_instance.basic.name}"
  name          = "%s"

}
`, rnd, rnd, dbName)
}

func testAccSpannerDatabase_duplicateNameError_part2(rnd, dbName string) string {
	return fmt.Sprintf(`
%s

resource "google_spanner_database" "basic2" {
  instance      = "${google_spanner_instance.basic.name}"
  name          = "%s"
}
`, testAccSpannerDatabase_duplicateNameError_part1(rnd, dbName), dbName)
}

func testAccSpannerDatabase_basicImport(iname, dbname string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  name          = "%s"
  config        = "regional-us-central1"
  display_name  = "%s"
  num_nodes     = 1
}

resource "google_spanner_database" "basic" {
  instance      = "${google_spanner_instance.basic.name}"
  name          = "%s"

}
`, iname, iname, dbname)
}

func testAccSpannerDatabase_basicImportWithProject(project, iname, dbname string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  project       = "%s"
  name          = "%s"
  config        = "regional-us-central1"
  display_name  = "%s"
  num_nodes     = 1
}

resource "google_spanner_database" "basic" {
  project       = "%s"
  instance      = "${google_spanner_instance.basic.name}"
  name          = "%s"

}
`, project, iname, iname, project, dbname)
}
