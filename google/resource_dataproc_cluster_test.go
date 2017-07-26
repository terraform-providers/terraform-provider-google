package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"google.golang.org/api/dataproc/v1"
	"google.golang.org/api/googleapi"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
)

const base10 = 10

func TestExtractLastResourceFromUri_withUrl(t *testing.T) {
	actual := extractLastResourceFromUri("http://something.com/one/two/three")
	expected := "three"
	if actual != expected {
		t.Fatalf("Expected %s, but got %s", expected, actual)
	}
}

func TestExtractLastResourceFromUri_WithStaticValue(t *testing.T) {
	actual := extractLastResourceFromUri("three")
	expected := "three"
	if actual != expected {
		t.Fatalf("Expected %s, but got %s", expected, actual)
	}
}

func TestExtractInitTimeout(t *testing.T) {
	actual, err := extractInitTimeout("500s")
	expected := 500
	if err != nil {
		t.Fatalf("Expected %d, but got error %v", expected, err)
	}
	if actual != expected {
		t.Fatalf("Expected %d, but got %d", expected, actual)
	}
}

func TestExtractInitTimeout_badFormat(t *testing.T) {
	_, err := extractInitTimeout("5m")
	expected := "Unexpected init timeout format expecting in seconds e.g. ZZZs, found : 5m"
	if err != nil && err.Error() == expected {
		return
	}
	t.Fatalf("Expected an error with message '%s', but got %v", expected, err)
}

func TestExtractInitTimeout_empty(t *testing.T) {
	_, err := extractInitTimeout("")
	expected := "Cannot extract init timeout from empty string"
	if err != nil && err.Error() == expected {
		return
	}
	t.Fatalf("Expected an error with message '%s', but got %v", expected, err)
}

func TestAccDataprocCluster_missingZoneGlobalRegion(t *testing.T) {
	rnd := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckDataproc_missingZoneGlobalRegion(rnd),
				ExpectError: regexp.MustCompile("zone is mandatory when region is set to 'global'"),
			},
		},
	})
}

func TestAccDataprocCluster_basic(t *testing.T) {
	rnd := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_basic(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterAttrMatch(
						"google_dataproc_cluster.basic"),
				),
			},
		},
	})
}

func TestAccDataprocCluster_singleNodeCluster(t *testing.T) {
	rnd := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_singleNodeCluster(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterAttrMatch(
						"google_dataproc_cluster.single_node_cluster"),
					resource.TestCheckResourceAttr(
						"google_dataproc_cluster.single_node_cluster",
						"master_config.0.num_masters",
						"1"),
					resource.TestCheckResourceAttr(
						"google_dataproc_cluster.single_node_cluster",
						"worker_config.0.num_workers",
						"0"),
				),
			},
		},
	})
}

func TestAccDataprocCluster_withBucketRef(t *testing.T) {
	rnd := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withBucketAndCluster(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterAttrMatch(
						"google_dataproc_cluster.with_bucket"),
				),
			},
			{
				// Simulate destroy of cluster by removing it
				Config: testAccDataprocCluster_withBucket(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterDeletedButNotCustomBucket(
						"us-central1",
						fmt.Sprintf("dproc-cluster-test-%s", rnd),
						fmt.Sprintf("dproc-cluster-test-%s-bucket", rnd),
					),
				),
			},
		},
	})
}

func TestAccDataprocCluster_withInitAction(t *testing.T) {
	rnd := acctest.RandString(10)
	bucketName := fmt.Sprintf("dproc-cluster-test-%s-init-bucket", rnd)
	objectName := "msg.txt"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withInitAction(rnd, bucketName, objectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterAttrMatch(
						"google_dataproc_cluster.with_init_action"),
					testInitActionSucceeded(
						bucketName, objectName),
				),
			},
		},
	})
}

func TestAccDataprocCluster_withConfigOverrides(t *testing.T) {
	rnd := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withConfigOverrides(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterAttrMatch(
						"google_dataproc_cluster.with_config_overrides"),
				),
			},
		},
	})
}

func TestAccDataprocCluster_withServiceAcc(t *testing.T) {

	saEmail := os.Getenv("GOOGLE_SERVICE_ACCOUNT")
	var cluster dataproc.Cluster
	rnd := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithServiceAccount(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withServiceAcc(saEmail, rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(
						"google_dataproc_cluster.with_service_account", &cluster),
					testAccCheckDataprocClusterHasServiceScopes(t, &cluster,
						"https://www.googleapis.com/auth/cloud.useraccounts.readonly",
						"https://www.googleapis.com/auth/devstorage.read_write",
						"https://www.googleapis.com/auth/logging.write",
						"https://www.googleapis.com/auth/monitoring",
					),
					testAccCheckDataprocClusterAttrMatch(
						"google_dataproc_cluster.with_service_account"),
				),
			},
		},
	})
}

func TestAccDataprocCluster_withImageVersion(t *testing.T) {
	rnd := acctest.RandString(10)
	var cluster dataproc.Cluster
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_withImageVersion(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(
						"google_dataproc_cluster.with_image_version", &cluster),
					testAccCheckDataprocClusterImageVersion(
						&cluster, "1.0.44"),
				),
			},
		},
	})
}

func testAccCheckDataprocClusterHasServiceScopes(t *testing.T, cluster *dataproc.Cluster, scopes ...string) func(s *terraform.State) error {
	return func(s *terraform.State) error {

		if !reflect.DeepEqual(scopes, cluster.Config.GceClusterConfig.ServiceAccountScopes) {
			return fmt.Errorf("Cluster does not contain expected set of service account scopes : %v : instead %v",
				scopes, cluster.Config.GceClusterConfig.ServiceAccountScopes)
		}
		return nil
	}
}
func testAccCheckDataprocClusterImageVersion(cluster *dataproc.Cluster, expected string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		actual := ""
		if cluster.Config != nil && cluster.Config.SoftwareConfig != nil {
			actual = cluster.Config.SoftwareConfig.ImageVersion
		}

		if actual != expected {
			return fmt.Errorf("Cluster Image version set to %s, but expected: %s",
				actual, expected)
		}
		return nil
	}
}

func TestAccDataprocCluster_network(t *testing.T) {
	rnd := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataprocClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocCluster_networkRef(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterAttrMatch(
						"google_dataproc_cluster.with_net_ref_by_url"),
					testAccCheckDataprocClusterAttrMatch(
						"google_dataproc_cluster.with_net_ref_by_name"),
				),
			},
		},
	})
}

func testAccCheckDataprocClusterDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_dataproc_cluster" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Unable to verify delete of dataproc cluster, ID is empty")
		}
		attributes := rs.Primary.Attributes

		validateClusterDeleted(config.Project, attributes["region"], rs.Primary.ID, config)
		validateAutoBucketsDeleted(attributes["staging_bucket"], attributes["bucket"], config)

	}

	return nil
}

func validateClusterDeleted(project, region, clusterName string, config *Config) error {
	_, err := config.clientDataproc.Projects.Regions.Clusters.Get(
		project, region, clusterName).Do()

	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return nil
		} else if ok {
			return fmt.Errorf("Error make GCP platform call to verify dataproc cluster deleted: http code error : %d, http message error: %s", gerr.Code, gerr.Message)
		}
		return fmt.Errorf("Error make GCP platform call to verify dataproc cluster deleted: %s", err.Error())
	}
	return fmt.Errorf("Dataproc cluster still exists")
}

func validateAutoBucketsDeleted(stagingBucketName, bucket string, config *Config) error {
	if stagingBucketName == "" {
		log.Printf("[DEBUG] explicit bucket specified %s (for dataproc cluster) leaving alone: \n\n", bucket)
		return nil
	}

	log.Printf("[DEBUG] validating autogen bucket %s (for dataproc cluster) is deleted \n\n", bucket)
	return validateBucketDoesNotExist(bucket, config)
}

func validateBucketDoesNotExist(bucket string, config *Config) error {
	_, err := config.clientStorage.Buckets.Get(bucket).Do()

	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return nil
		} else if ok {
			return fmt.Errorf("Error make GCP platform call to verify if bucket deleted: http code error : %d, http message error: %s", gerr.Code, gerr.Message)
		}
		return fmt.Errorf("Error make GCP platform call to verify if bucket deleted: %s", err.Error())
	}
	return fmt.Errorf("bucket still exists")
}

func validateBucketExists(bucket string, config *Config) (bool, error) {
	_, err := config.clientStorage.Buckets.Get(bucket).Do()

	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return false, nil
		} else if ok {
			return false, fmt.Errorf("Error make GCP platform call to verify if bucket deleted: http code error : %d, http message error: %s", gerr.Code, gerr.Message)
		}
		return false, fmt.Errorf("Error make GCP platform call to verify if bucket deleted: %s", err.Error())
	}
	return true, nil
}

func testAccCheckDataprocClusterDeletedButNotCustomBucket(region, clusterName, bucketName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Check cluster is gone
		config := testAccProvider.Meta().(*Config)
		err := validateClusterDeleted(config.Project, region, clusterName, config)
		if err != nil {
			return err
		}

		// Check Original Custom Bucket still exists
		exists, err := validateBucketExists(bucketName, config)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Bucket: %s does not exist", bucketName)
		}
		return nil
	}
}

func testAccCheckDataprocClusterAttrMatch(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attributes, err := getResourceAttributes(n, s)
		if err != nil {
			return err
		}

		config := testAccProvider.Meta().(*Config)
		cluster, err := config.clientDataproc.Projects.Regions.Clusters.Get(
			config.Project, attributes["region"], attributes["name"]).Do()
		if err != nil {
			return err
		}

		if cluster.ClusterName != attributes["name"] {
			return fmt.Errorf("Cluster %s not found, found %s instead", attributes["name"], cluster.ClusterName)
		}

		type clusterTestField struct {
			tf_attr  string
			gcp_attr interface{}
		}

		clusterTests := []clusterTestField{

			{"bucket", cluster.Config.ConfigBucket},
			{"image_version", cluster.Config.SoftwareConfig.ImageVersion},
			{"zone", extractLastResourceFromUri(cluster.Config.GceClusterConfig.ZoneUri)},

			{"network", extractLastResourceFromUri(cluster.Config.GceClusterConfig.NetworkUri)},
			{"subnetwork", extractLastResourceFromUri(cluster.Config.GceClusterConfig.SubnetworkUri)},
			{"service_account", cluster.Config.GceClusterConfig.ServiceAccount},
			{"service_account_scopes", cluster.Config.GceClusterConfig.ServiceAccountScopes},
			{"metadata", cluster.Config.GceClusterConfig.Metadata},
			{"labels", cluster.Labels},
			{"tags", cluster.Config.GceClusterConfig.Tags},
		}

		if cluster.Config.MasterConfig != nil {
			clusterTests = append(clusterTests,
				clusterTestField{"master_config.0.num_masters", strconv.FormatInt(cluster.Config.MasterConfig.NumInstances, base10)},
				clusterTestField{"master_config.0.boot_disk_size_gb", strconv.FormatInt(cluster.Config.MasterConfig.DiskConfig.BootDiskSizeGb, base10)},
				clusterTestField{"master_config.0.num_local_ssds", strconv.FormatInt(cluster.Config.MasterConfig.DiskConfig.NumLocalSsds, base10)},
				clusterTestField{"master_config.0.machine_type", extractLastResourceFromUri(cluster.Config.MasterConfig.MachineTypeUri)})
		}

		if cluster.Config.WorkerConfig != nil {
			clusterTests = append(clusterTests,
				clusterTestField{"worker_config.0.num_workers", strconv.FormatInt(cluster.Config.WorkerConfig.NumInstances, base10)},
				clusterTestField{"worker_config.0.boot_disk_size_gb", strconv.FormatInt(cluster.Config.WorkerConfig.DiskConfig.BootDiskSizeGb, base10)},
				clusterTestField{"worker_config.0.num_local_ssds", strconv.FormatInt(cluster.Config.WorkerConfig.DiskConfig.NumLocalSsds, base10)},
				clusterTestField{"worker_config.0.machine_type", extractLastResourceFromUri(cluster.Config.WorkerConfig.MachineTypeUri)})
		}

		if cluster.Config.SecondaryWorkerConfig != nil {
			clusterTests = append(clusterTests,
				clusterTestField{"worker_config.0.preemptible_num_workers", strconv.FormatInt(cluster.Config.SecondaryWorkerConfig.NumInstances, base10)},
				clusterTestField{"worker_config.0.preemptible_boot_disk_size_gb", strconv.FormatInt(cluster.Config.SecondaryWorkerConfig.DiskConfig.BootDiskSizeGb, base10)})
		}

		extracted := false
		if len(cluster.Config.InitializationActions) > 0 {
			actions := []string{}
			for _, v := range cluster.Config.InitializationActions {
				actions = append(actions, v.ExecutableFile)

				if !extracted && len(v.ExecutionTimeout) > 0 {
					tsec, err := extractInitTimeout(v.ExecutionTimeout)
					if err != nil {
						return err
					}
					clusterTests = append(clusterTests, clusterTestField{"initialization_action_timeout_sec", strconv.Itoa(tsec)})
					extracted = true
				}
			}
			clusterTests = append(clusterTests, clusterTestField{"initialization_actions", actions})
		}

		for _, attrs := range clusterTests {
			if c := checkMatch(attributes, attrs.tf_attr, attrs.gcp_attr); c != "" {
				return fmt.Errorf(c)
			}
		}

		return nil
	}
}

func testInitActionSucceeded(bucket, object string) resource.TestCheckFunc {

	// The init script will have created an object in the specified bucket.
	// Ensure it exists
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		_, err := config.clientStorage.Objects.Get(bucket, object).Do()
		if err != nil {
			return fmt.Errorf("Unable to verify init action success: Error reading object %s in bucket %s: %v", object, bucket, err)
		}

		return nil
	}
}

func testAccCheckDataproc_missingZoneGlobalRegion(rnd string) string {
	return fmt.Sprintf(`
resource "google_dataproc_cluster" "basic" {
	name   = "dproc-cluster-test-%s"
	region = "global"
}
`, rnd)
}

func testAccDataprocCluster_basic(rnd string) string {
	return fmt.Sprintf(`
resource "google_dataproc_cluster" "basic" {
	name   = "dproc-cluster-test-%s"
	region = "us-central1"
}
`, rnd)
}

func testAccDataprocCluster_singleNodeCluster(rnd string) string {
	return fmt.Sprintf(`
resource "google_dataproc_cluster" "single_node_cluster" {
	name   = "dproc-cluster-test-%s"
	region = "us-central1"

    properties = {
        "dataproc:dataproc.allow.zero.workers" = "true"
    }

    # Because of current restrictions with computed AND default
    # [list|Set] properties, we need to add this empty config
    # here otherwise if you plan straight away afterwards you
    # will get a diff. If you have actual config values that is
    # fine, but if you were hoping to use the defaults, this is
    # required
    master_config { }
    worker_config { }
}
`, rnd)
}

func testAccDataprocCluster_withConfigOverrides(rnd string) string {
	return fmt.Sprintf(`
resource "google_dataproc_cluster" "with_config_overrides" {
	name   = "dproc-cluster-test-%s"
	region = "us-central1"

	master_config {
		num_masters       = 1
		machine_type      = "n1-standard-1"
		boot_disk_size_gb = 10
		num_local_ssds    = 0
	}

	worker_config {
	    num_workers       = 2
		machine_type      = "n1-standard-1"
		boot_disk_size_gb = 10
		num_local_ssds    = 0

		preemptible_num_workers       = 1
		preemptible_boot_disk_size_gb = 10
	}
}`, rnd)
}

func testAccDataprocCluster_withInitAction(rnd, bucket, objName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "init_bucket" {
    name          = "%s"
    force_destroy = "true"
}

resource "google_storage_bucket_object" "init_script" {
  name           = "dproc-cluster-test-%s-init-script.sh"
  bucket         = "${google_storage_bucket.init_bucket.name}"
  content        = <<EOL
#!/bin/bash
ROLE=$$(/usr/share/google/get_metadata_value attributes/dataproc-role)
if [[ "$${ROLE}" == 'Master' ]]; then
  echo "on the master" >> /tmp/%s
  gsutil cp /tmp/%s ${google_storage_bucket.init_bucket.url}
else
  echo "on the worker" >> /tmp/msg.txt
fi
EOL

}

resource "google_dataproc_cluster" "with_init_action" {
	name   = "dproc-cluster-test-%s"
	region = "us-central1"
	initialization_action_timeout_sec = 500
	initialization_actions = [
	   "${google_storage_bucket.init_bucket.url}/${google_storage_bucket_object.init_script.name}"
	]

    # Keep the costs down with smallest config we can get away with
    properties = {
        "dataproc:dataproc.allow.zero.workers" = "true"
    }
    worker_config { }
	master_config {
		machine_type      = "n1-standard-1"
		boot_disk_size_gb = 10
	}
}`, bucket, rnd, objName, objName, rnd)
}

func testAccDataprocCluster_withBucket(rnd string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
    name          = "dproc-cluster-test-%s-bucket"
    force_destroy = "true"
}`, rnd)
}

func testAccDataprocCluster_withBucketAndCluster(rnd string) string {
	return fmt.Sprintf(`
%s

resource "google_dataproc_cluster" "with_bucket" {
	name   = "dproc-cluster-test-%s"
	region = "us-central1"
	staging_bucket = "${google_storage_bucket.bucket.name}"

    # Keep the costs down with smallest config we can get away with
    properties = {
        "dataproc:dataproc.allow.zero.workers" = "true"
    }
    worker_config { }
	master_config {
		machine_type      = "n1-standard-1"
		boot_disk_size_gb = 10
	}
}`, testAccDataprocCluster_withBucket(rnd), rnd)
}

func testAccDataprocCluster_withImageVersion(rnd string) string {
	return fmt.Sprintf(`
resource "google_dataproc_cluster" "with_image_version" {
	name   = "dproc-cluster-test-%s"
	region = "us-central1"
	image_version = "1.0.44"
}`, rnd)
}

func testAccDataprocCluster_withServiceAcc(saEmail string, rnd string) string {
	return fmt.Sprintf(`
resource "google_dataproc_cluster" "with_service_account" {
	name   = "dproc-cluster-test-%s"
	region = "us-central1"

    # Keep the costs down with smallest config we can get away with
    properties = {
        "dataproc:dataproc.allow.zero.workers" = "true"
    }
    worker_config { }
	master_config {
		machine_type      = "n1-standard-1"
		boot_disk_size_gb = 10
	}

	service_account = "%s"

	service_account_scopes = [
        #    The following scopes necessary for the cluster to function properly are
		#	always added, even if not explicitly specified:
		#		useraccounts-ro: https://www.googleapis.com/auth/cloud.useraccounts.readonly
		#		storage-rw:      https://www.googleapis.com/auth/devstorage.read_write
		#		logging-write:   https://www.googleapis.com/auth/logging.write
        #
		#	So user is expected to add these explicitly (in this order) otherwise terraform
		#   will think there is a change to resource
		"useraccounts-ro","storage-rw","logging-write",

	    # Additional ones specifically desired by user (Note for now must be in alpha order
	    # of fully qualified scope name)
	    "https://www.googleapis.com/auth/monitoring"

	]

}`, rnd, saEmail)
}

func testAccDataprocCluster_networkRef(rnd string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "dataproc_network" {
	name = "dproc-cluster-test-%s-net"
	auto_create_subnetworks = true
}

resource "google_compute_firewall" "dataproc_network_firewall" {
	name = "dproc-cluster-test-%s-allow-internal"
	description = "Firewall rules for dataproc Terraform acceptance testing"
	network = "${google_compute_network.dataproc_network.name}"

	allow {
	    protocol = "icmp"
	}

	allow {
		protocol = "tcp"
		ports    = ["0-65535"]
	}

	allow {
		protocol = "udp"
		ports    = ["0-65535"]
	}
}

resource "google_dataproc_cluster" "with_net_ref_by_name" {
	name   = "dproc-cluster-test-%s-name"
	region = "us-central1"
	depends_on = ["google_compute_firewall.dataproc_network_firewall"]

    # Keep the costs down with smallest config we can get away with
    properties = {
        "dataproc:dataproc.allow.zero.workers" = "true"
    }
    worker_config { }
	master_config {
		machine_type      = "n1-standard-1"
		boot_disk_size_gb = 10
	}

	network = "${google_compute_network.dataproc_network.name}"
}

resource "google_dataproc_cluster" "with_net_ref_by_url" {
	name   = "dproc-cluster-test-%s-url"
	region = "us-central1"
    depends_on = ["google_compute_firewall.dataproc_network_firewall"]

    # Keep the costs down with smallest config we can get away with
    properties = {
        "dataproc:dataproc.allow.zero.workers" = "true"
    }
    worker_config { }
	master_config {
		machine_type      = "n1-standard-1"
		boot_disk_size_gb = 10
	}

	network = "${google_compute_network.dataproc_network.self_link}"
}

`, rnd, rnd, rnd, rnd)
}

func testAccPreCheckWithServiceAccount(t *testing.T) {
	testAccPreCheck(t)
	if v := os.Getenv("GOOGLE_SERVICE_ACCOUNT"); v == "" {
		t.Skipf("GOOGLE_SERVICE_ACCOUNT must be set for the dataproc acceptance test testing service account functionality")
	}

}

func testAccCheckDataprocClusterExists(n string, cluster *dataproc.Cluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		found, err := config.clientDataproc.Projects.Regions.Clusters.Get(
			config.Project, rs.Primary.Attributes["region"], rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		if found.ClusterName != rs.Primary.ID {
			return fmt.Errorf("Cluster not found")
		}

		*cluster = *found

		return nil
	}
}
