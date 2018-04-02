package google

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/api/compute/v1"
)

const DEFAULT_MIN_CPU_TEST_VALUE = "Intel Haswell"

func TestAccComputeInstanceTemplate_basic(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeInstanceTemplate_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(
						"google_compute_instance_template.foobar", &instanceTemplate),
					testAccCheckComputeInstanceTemplateTag(&instanceTemplate, "foo"),
					testAccCheckComputeInstanceTemplateMetadata(&instanceTemplate, "foo", "bar"),
					testAccCheckComputeInstanceTemplateDisk(&instanceTemplate, "projects/debian-cloud/global/images/debian-8-jessie-v20160803", true, true),
					testAccCheckComputeInstanceTemplateContainsLabel(&instanceTemplate, "my_label", "foobar"),
					testAccCheckComputeInstanceTemplateContainsLabel(&instanceTemplate, "goog-partner-creation-tool", "terraform"),
				),
			},
			resource.TestStep{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceTemplate_preemptible(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeInstanceTemplate_preemptible(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(
						"google_compute_instance_template.foobar", &instanceTemplate),
					testAccCheckComputeInstanceTemplateAutomaticRestart(&instanceTemplate, false),
					testAccCheckComputeInstanceTemplatePreemptible(&instanceTemplate, true),
				),
			},
			resource.TestStep{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceTemplate_IP(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeInstanceTemplate_ip(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(
						"google_compute_instance_template.foobar", &instanceTemplate),
					testAccCheckComputeInstanceTemplateNetwork(&instanceTemplate),
				),
			},
			resource.TestStep{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceTemplate_networkIP(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate
	networkIP := "10.128.0.2"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeInstanceTemplate_networkIP(networkIP),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(
						"google_compute_instance_template.foobar", &instanceTemplate),
					testAccCheckComputeInstanceTemplateNetwork(&instanceTemplate),
					testAccCheckComputeInstanceTemplateNetworkIP(
						"google_compute_instance_template.foobar", networkIP, &instanceTemplate),
				),
			},
			resource.TestStep{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func TestAccComputeInstanceTemplate_address(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate
	address := "10.128.0.2"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeInstanceTemplate_address(address),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(
						"google_compute_instance_template.foobar", &instanceTemplate),
					testAccCheckComputeInstanceTemplateNetwork(&instanceTemplate),
					testAccCheckComputeInstanceTemplateAddress(
						"google_compute_instance_template.foobar", address, &instanceTemplate),
				),
			},
			resource.TestStep{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceTemplate_disks(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeInstanceTemplate_disks(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(
						"google_compute_instance_template.foobar", &instanceTemplate),
					testAccCheckComputeInstanceTemplateDisk(&instanceTemplate, "projects/debian-cloud/global/images/debian-8-jessie-v20160803", true, true),
					testAccCheckComputeInstanceTemplateDisk(&instanceTemplate, "terraform-test-foobar", false, false),
				),
			},
			resource.TestStep{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceTemplate_subnet_auto(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate
	network := "network-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeInstanceTemplate_subnet_auto(network),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(
						"google_compute_instance_template.foobar", &instanceTemplate),
					testAccCheckComputeInstanceTemplateNetworkName(&instanceTemplate, network),
				),
			},
			resource.TestStep{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceTemplate_subnet_custom(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeInstanceTemplate_subnet_custom(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(
						"google_compute_instance_template.foobar", &instanceTemplate),
					testAccCheckComputeInstanceTemplateSubnetwork(&instanceTemplate),
				),
			},
			resource.TestStep{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceTemplate_subnet_xpn(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate
	org := getTestOrgFromEnv(t)
	billingId := getTestBillingAccountFromEnv(t)
	projectName := fmt.Sprintf("tf-xpntest-%d", time.Now().Unix())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeInstanceTemplate_subnet_xpn(org, billingId, projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExistsInProject(
						"google_compute_instance_template.foobar", fmt.Sprintf("%s-service", projectName),
						&instanceTemplate),
					testAccCheckComputeInstanceTemplateSubnetwork(&instanceTemplate),
				),
			},
		},
	})
}

func TestAccComputeInstanceTemplate_metadata_startup_script(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeInstanceTemplate_startup_script(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists(
						"google_compute_instance_template.foobar", &instanceTemplate),
					testAccCheckComputeInstanceTemplateStartupScript(&instanceTemplate, "echo 'Hello'"),
				),
			},
		},
	})
}
func TestAccComputeInstanceTemplate_primaryAliasIpRange(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeInstanceTemplate_primaryAliasIpRange(acctest.RandString(10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists("google_compute_instance_template.foobar", &instanceTemplate),
					testAccCheckComputeInstanceTemplateHasAliasIpRange(&instanceTemplate, "", "/24"),
				),
			},
			resource.TestStep{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceTemplate_secondaryAliasIpRange(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeInstanceTemplate_secondaryAliasIpRange(acctest.RandString(10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists("google_compute_instance_template.foobar", &instanceTemplate),
					testAccCheckComputeInstanceTemplateHasAliasIpRange(&instanceTemplate, "inst-test-secondary", "/24"),
				),
			},
			resource.TestStep{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeInstanceTemplate_guestAccelerator(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeInstanceTemplate_guestAccelerator(acctest.RandString(10), 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists("google_compute_instance_template.foobar", &instanceTemplate),
					testAccCheckComputeInstanceTemplateHasGuestAccelerator(&instanceTemplate, "nvidia-tesla-k80", 1),
				),
			},
			resource.TestStep{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

}

func TestAccComputeInstanceTemplate_guestAcceleratorSkip(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeInstanceTemplate_guestAccelerator(acctest.RandString(10), 0),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists("google_compute_instance_template.foobar", &instanceTemplate),
					testAccCheckComputeInstanceTemplateLacksGuestAccelerator(&instanceTemplate),
				),
			},
		},
	})

}

func TestAccComputeInstanceTemplate_minCpuPlatform(t *testing.T) {
	t.Parallel()

	var instanceTemplate compute.InstanceTemplate

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeInstanceTemplate_minCpuPlatform(acctest.RandString(10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceTemplateExists("google_compute_instance_template.foobar", &instanceTemplate),
					testAccCheckComputeInstanceTemplateHasMinCpuPlatform(&instanceTemplate, DEFAULT_MIN_CPU_TEST_VALUE),
				),
			},
			resource.TestStep{
				ResourceName:      "google_compute_instance_template.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckComputeInstanceTemplateDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_compute_instance_template" {
			continue
		}

		_, err := config.clientCompute.InstanceTemplates.Get(
			config.Project, rs.Primary.ID).Do()
		if err == nil {
			return fmt.Errorf("Instance template still exists")
		}
	}

	return nil
}

func testAccCheckComputeInstanceTemplateExists(n string, instanceTemplate *compute.InstanceTemplate) resource.TestCheckFunc {
	return testAccCheckComputeInstanceTemplateExistsInProject(n, getTestProjectFromEnv(), instanceTemplate)
}

func testAccCheckComputeInstanceTemplateExistsInProject(n, p string, instanceTemplate *compute.InstanceTemplate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		found, err := config.clientCompute.InstanceTemplates.Get(
			p, rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("Instance template not found")
		}

		*instanceTemplate = *found

		return nil
	}
}

func testAccCheckComputeInstanceTemplateMetadata(
	instanceTemplate *compute.InstanceTemplate,
	k string, v string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instanceTemplate.Properties.Metadata == nil {
			return fmt.Errorf("no metadata")
		}

		for _, item := range instanceTemplate.Properties.Metadata.Items {
			if k != item.Key {
				continue
			}

			if item.Value != nil && v == *item.Value {
				return nil
			}

			return fmt.Errorf("bad value for %s: %s", k, *item.Value)
		}

		return fmt.Errorf("metadata not found: %s", k)
	}
}

func testAccCheckComputeInstanceTemplateNetwork(instanceTemplate *compute.InstanceTemplate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, i := range instanceTemplate.Properties.NetworkInterfaces {
			for _, c := range i.AccessConfigs {
				if c.NatIP == "" {
					return fmt.Errorf("no NAT IP")
				}
			}
		}

		return nil
	}
}

func testAccCheckComputeInstanceTemplateNetworkName(instanceTemplate *compute.InstanceTemplate, network string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, i := range instanceTemplate.Properties.NetworkInterfaces {
			if !strings.Contains(i.Network, network) {
				return fmt.Errorf("Network doesn't match expected value, Expected: %s Actual: %s", network, i.Network[strings.LastIndex("/", i.Network)+1:])
			}
		}

		return nil
	}
}

func testAccCheckComputeInstanceTemplateDisk(instanceTemplate *compute.InstanceTemplate, source string, delete bool, boot bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instanceTemplate.Properties.Disks == nil {
			return fmt.Errorf("no disks")
		}

		for _, disk := range instanceTemplate.Properties.Disks {
			if disk.InitializeParams == nil {
				// Check disk source
				if disk.Source == source {
					if disk.AutoDelete == delete && disk.Boot == boot {
						return nil
					}
				}
			} else {
				// Check source image
				if disk.InitializeParams.SourceImage == source {
					if disk.AutoDelete == delete && disk.Boot == boot {
						return nil
					}
				}
			}
		}

		return fmt.Errorf("Disk not found: %s", source)
	}
}

func testAccCheckComputeInstanceTemplateSubnetwork(instanceTemplate *compute.InstanceTemplate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, i := range instanceTemplate.Properties.NetworkInterfaces {
			if i.Subnetwork == "" {
				return fmt.Errorf("no subnet")
			}
		}

		return nil
	}
}

func testAccCheckComputeInstanceTemplateTag(instanceTemplate *compute.InstanceTemplate, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instanceTemplate.Properties.Tags == nil {
			return fmt.Errorf("no tags")
		}

		for _, k := range instanceTemplate.Properties.Tags.Items {
			if k == n {
				return nil
			}
		}

		return fmt.Errorf("tag not found: %s", n)
	}
}

func testAccCheckComputeInstanceTemplatePreemptible(instanceTemplate *compute.InstanceTemplate, preemptible bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instanceTemplate.Properties.Scheduling.Preemptible != preemptible {
			return fmt.Errorf("Expected preemptible value %v, got %v", preemptible, instanceTemplate.Properties.Scheduling.Preemptible)
		}
		return nil
	}
}

func testAccCheckComputeInstanceTemplateAutomaticRestart(instanceTemplate *compute.InstanceTemplate, automaticRestart bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ar := instanceTemplate.Properties.Scheduling.AutomaticRestart
		if ar == nil {
			return fmt.Errorf("Expected to see a value for AutomaticRestart, but got nil")
		}
		if *ar != automaticRestart {
			return fmt.Errorf("Expected automatic restart value %v, got %v", automaticRestart, ar)
		}
		return nil
	}
}

func testAccCheckComputeInstanceTemplateStartupScript(instanceTemplate *compute.InstanceTemplate, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instanceTemplate.Properties.Metadata == nil && n == "" {
			return nil
		} else if instanceTemplate.Properties.Metadata == nil && n != "" {
			return fmt.Errorf("Expected metadata.startup-script to be '%s', metadata wasn't set at all", n)
		}
		for _, item := range instanceTemplate.Properties.Metadata.Items {
			if item.Key != "startup-script" {
				continue
			}
			if item.Value != nil && *item.Value == n {
				return nil
			} else if item.Value == nil && n == "" {
				return nil
			} else if item.Value == nil && n != "" {
				return fmt.Errorf("Expected metadata.startup-script to be '%s', wasn't set", n)
			} else if *item.Value != n {
				return fmt.Errorf("Expected metadata.startup-script to be '%s', got '%s'", n, *item.Value)
			}
		}
		return fmt.Errorf("This should never be reached.")
	}
}

func testAccCheckComputeInstanceTemplateNetworkIP(n, networkIP string, instanceTemplate *compute.InstanceTemplate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ip := instanceTemplate.Properties.NetworkInterfaces[0].NetworkIP
		err := resource.TestCheckResourceAttr(n, "network_interface.0.network_ip", ip)(s)
		if err != nil {
			return err
		}
		return resource.TestCheckResourceAttr(n, "network_interface.0.network_ip", networkIP)(s)
	}
}

func testAccCheckComputeInstanceTemplateAddress(n, address string, instanceTemplate *compute.InstanceTemplate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ip := instanceTemplate.Properties.NetworkInterfaces[0].NetworkIP
		err := resource.TestCheckResourceAttr(n, "network_interface.0.network_ip", ip)(s)
		if err != nil {
			return err
		}
		return resource.TestCheckResourceAttr(n, "network_interface.0.network_ip", address)(s)
	}
}

func testAccCheckComputeInstanceTemplateContainsLabel(instanceTemplate *compute.InstanceTemplate, key string, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		v, ok := instanceTemplate.Properties.Labels[key]
		if !ok {
			return fmt.Errorf("Expected label with key '%s' not found", key)
		}
		if v != value {
			return fmt.Errorf("Incorrect label value for key '%s': expected '%s' but found '%s'", key, value, v)
		}
		return nil
	}
}

func testAccCheckComputeInstanceTemplateHasAliasIpRange(instanceTemplate *compute.InstanceTemplate, subnetworkRangeName, iPCidrRange string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, networkInterface := range instanceTemplate.Properties.NetworkInterfaces {
			for _, aliasIpRange := range networkInterface.AliasIpRanges {
				if aliasIpRange.SubnetworkRangeName == subnetworkRangeName && (aliasIpRange.IpCidrRange == iPCidrRange || ipCidrRangeDiffSuppress("ip_cidr_range", aliasIpRange.IpCidrRange, iPCidrRange, nil)) {
					return nil
				}
			}
		}

		return fmt.Errorf("Alias ip range with name %s and cidr %s not present", subnetworkRangeName, iPCidrRange)
	}
}

func testAccCheckComputeInstanceTemplateHasGuestAccelerator(instanceTemplate *compute.InstanceTemplate, acceleratorType string, acceleratorCount int64) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(instanceTemplate.Properties.GuestAccelerators) != 1 {
			return fmt.Errorf("Expected only one guest accelerator")
		}

		if !strings.HasSuffix(instanceTemplate.Properties.GuestAccelerators[0].AcceleratorType, acceleratorType) {
			return fmt.Errorf("Wrong accelerator type: expected %v, got %v", acceleratorType, instanceTemplate.Properties.GuestAccelerators[0].AcceleratorType)
		}

		if instanceTemplate.Properties.GuestAccelerators[0].AcceleratorCount != acceleratorCount {
			return fmt.Errorf("Wrong accelerator acceleratorCount: expected %d, got %d", acceleratorCount, instanceTemplate.Properties.GuestAccelerators[0].AcceleratorCount)
		}

		return nil
	}
}

func testAccCheckComputeInstanceTemplateLacksGuestAccelerator(instanceTemplate *compute.InstanceTemplate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(instanceTemplate.Properties.GuestAccelerators) > 0 {
			return fmt.Errorf("Expected no guest accelerators")
		}

		return nil
	}
}

func testAccCheckComputeInstanceTemplateHasMinCpuPlatform(instanceTemplate *compute.InstanceTemplate, minCpuPlatform string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instanceTemplate.Properties.MinCpuPlatform != minCpuPlatform {
			return fmt.Errorf("Wrong minimum CPU platform: expected %s, got %s", minCpuPlatform, instanceTemplate.Properties.MinCpuPlatform)
		}

		return nil
	}
}

func testAccComputeInstanceTemplate_basic() string {
	return fmt.Sprintf(`
resource "google_compute_instance_template" "foobar" {
	name = "instancet-test-%s"
	machine_type = "n1-standard-1"
	can_ip_forward = false
	tags = ["foo", "bar"]

	disk {
		source_image = "debian-8-jessie-v20160803"
		auto_delete = true
		boot = true
	}

	network_interface {
		network = "default"
	}

	scheduling {
		preemptible = false
		automatic_restart = true
	}

	metadata {
		foo = "bar"
	}

	service_account {
		scopes = ["userinfo-email", "compute-ro", "storage-ro"]
	}

    labels {
        my_label = "foobar"
    }
}`, acctest.RandString(10))
}

func testAccComputeInstanceTemplate_preemptible() string {
	return fmt.Sprintf(`
resource "google_compute_instance_template" "foobar" {
	name = "instancet-test-%s"
	machine_type = "n1-standard-1"
	can_ip_forward = false
	tags = ["foo", "bar"]

	disk {
		source_image = "debian-8-jessie-v20160803"
		auto_delete = true
		boot = true
	}

	network_interface {
		network = "default"
	}

	scheduling {
		preemptible = true
		automatic_restart = false
	}

	metadata {
		foo = "bar"
	}

	service_account {
		scopes = ["userinfo-email", "compute-ro", "storage-ro"]
	}
}`, acctest.RandString(10))
}

func testAccComputeInstanceTemplate_ip() string {
	return fmt.Sprintf(`
resource "google_compute_address" "foo" {
	name = "instancet-test-%s"
}

resource "google_compute_instance_template" "foobar" {
	name = "instancet-test-%s"
	machine_type = "n1-standard-1"
	tags = ["foo", "bar"]

	disk {
		source_image = "debian-8-jessie-v20160803"
	}

	network_interface {
		network = "default"
		access_config {
			nat_ip = "${google_compute_address.foo.address}"
		}
	}

	metadata {
		foo = "bar"
	}
}`, acctest.RandString(10), acctest.RandString(10))
}

func testAccComputeInstanceTemplate_networkIP(networkIP string) string {
	return fmt.Sprintf(`
resource "google_compute_instance_template" "foobar" {
	name = "instancet-test-%s"
	machine_type = "n1-standard-1"
	tags = ["foo", "bar"]

	disk {
		source_image = "debian-8-jessie-v20160803"
	}

	network_interface {
		network    = "default"
		network_ip = "%s"
	}

	metadata {
		foo = "bar"
	}
}`, acctest.RandString(10), networkIP)
}

func testAccComputeInstanceTemplate_address(address string) string {
	return fmt.Sprintf(`
resource "google_compute_instance_template" "foobar" {
	name = "instancet-test-%s"
	machine_type = "n1-standard-1"
	tags = ["foo", "bar"]

	disk {
		source_image = "debian-8-jessie-v20160803"
	}

	network_interface {
		network    = "default"
		address    = "%s"
	}

	metadata {
		foo = "bar"
	}
}`, acctest.RandString(10), address)
}

func testAccComputeInstanceTemplate_disks() string {
	return fmt.Sprintf(`
resource "google_compute_disk" "foobar" {
	name = "instancet-test-%s"
	image = "debian-8-jessie-v20160803"
	size = 10
	type = "pd-ssd"
	zone = "us-central1-a"
}

resource "google_compute_instance_template" "foobar" {
	name = "instancet-test-%s"
	machine_type = "n1-standard-1"

	disk {
		source_image = "debian-8-jessie-v20160803"
		auto_delete = true
		disk_size_gb = 100
		boot = true
	}

	disk {
		source = "terraform-test-foobar"
		auto_delete = false
		boot = false
	}

	network_interface {
		network = "default"
	}

	metadata {
		foo = "bar"
	}
}`, acctest.RandString(10), acctest.RandString(10))
}

func testAccComputeInstanceTemplate_subnet_auto(network string) string {
	return fmt.Sprintf(`
	resource "google_compute_network" "auto-network" {
		name = "%s"
		auto_create_subnetworks = true
	}

	resource "google_compute_instance_template" "foobar" {
		name = "instance-tpl-%s"
		machine_type = "n1-standard-1"

		disk {
			source_image = "debian-8-jessie-v20160803"
			auto_delete = true
			disk_size_gb = 10
			boot = true
		}

		network_interface {
			network = "${google_compute_network.auto-network.name}"
		}

		metadata {
			foo = "bar"
		}
	}`, network, acctest.RandString(10))
}

func testAccComputeInstanceTemplate_subnet_custom() string {
	return fmt.Sprintf(`
resource "google_compute_network" "network" {
	name = "network-%s"
	auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnetwork" {
	name = "subnetwork-%s"
	ip_cidr_range = "10.0.0.0/24"
	region = "us-central1"
	network = "${google_compute_network.network.self_link}"
}

resource "google_compute_instance_template" "foobar" {
	name = "instance-test-%s"
	machine_type = "n1-standard-1"
	region = "us-central1"

	disk {
		source_image = "debian-8-jessie-v20160803"
		auto_delete = true
		disk_size_gb = 10
		boot = true
	}

	network_interface {
		subnetwork = "${google_compute_subnetwork.subnetwork.name}"
	}

	metadata {
		foo = "bar"
	}
}`, acctest.RandString(10), acctest.RandString(10), acctest.RandString(10))
}

func testAccComputeInstanceTemplate_subnet_xpn(org, billingId, projectName string) string {
	return fmt.Sprintf(`
	resource "google_project" "host_project" {
		name = "Test Project XPN Host"
		project_id = "%s-host"
		org_id = "%s"
		billing_account = "%s"
	}

	resource "google_project_service" "host_project" {
		project = "${google_project.host_project.project_id}"
		service = "compute.googleapis.com"
	}

	resource "google_compute_shared_vpc_host_project" "host_project" {
		project = "${google_project_service.host_project.project}"
	}
	
	resource "google_project" "service_project" {
		name = "Test Project XPN Service"
		project_id = "%s-service"
		org_id = "%s"
		billing_account = "%s"
	}

	resource "google_project_service" "service_project" {
		project = "${google_project.service_project.project_id}"
		service = "compute.googleapis.com"
	}

	resource "google_compute_shared_vpc_service_project" "service_project" {
		host_project = "${google_compute_shared_vpc_host_project.host_project.project}"
		service_project = "${google_project_service.service_project.project}"
	}

	resource "google_compute_network" "network" {
		name = "network-%s"
		auto_create_subnetworks = false
		project = "${google_compute_shared_vpc_host_project.host_project.project}"
	}

	resource "google_compute_subnetwork" "subnetwork" {
		name = "subnetwork-%s"
		ip_cidr_range = "10.0.0.0/24"
		region = "us-central1"
		network = "${google_compute_network.network.self_link}"
		project = "${google_compute_shared_vpc_host_project.host_project.project}"
	}

	resource "google_compute_instance_template" "foobar" {
		name = "instance-test-%s"
		machine_type = "n1-standard-1"
		region = "us-central1"

		disk {
			source_image = "debian-8-jessie-v20160803"
			auto_delete = true
			disk_size_gb = 10
			boot = true
		}

		network_interface {
			subnetwork = "${google_compute_subnetwork.subnetwork.name}"
			subnetwork_project = "${google_compute_subnetwork.subnetwork.project}"
		}

		metadata {
			foo = "bar"
		}
		project = "${google_compute_shared_vpc_service_project.service_project.service_project}"
	}`, projectName, org, billingId, projectName, org, billingId, acctest.RandString(10), acctest.RandString(10), acctest.RandString(10))
}

func testAccComputeInstanceTemplate_startup_script() string {
	return fmt.Sprintf(`
resource "google_compute_instance_template" "foobar" {
	name = "instance-test-%s"
	machine_type = "n1-standard-1"

	disk {
		source_image = "debian-8-jessie-v20160803"
		auto_delete = true
		disk_size_gb = 10
		boot = true
	}

	metadata {
		foo = "bar"
	}

	network_interface{
		network = "default"
	}

	metadata_startup_script = "echo 'Hello'"
}`, acctest.RandString(10))
}

func testAccComputeInstanceTemplate_primaryAliasIpRange(i string) string {
	return fmt.Sprintf(`
resource "google_compute_instance_template" "foobar" {
	name = "instance-test-%s"
	machine_type = "n1-standard-1"

	disk {
		source_image = "debian-8-jessie-v20160803"
		auto_delete = true
		disk_size_gb = 10
		boot = true
	}

	metadata {
		foo = "bar"
	}

	network_interface{
		network = "default"
		alias_ip_range {
			ip_cidr_range = "/24"
		}
	}
}`, i)
}

func testAccComputeInstanceTemplate_secondaryAliasIpRange(i string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "inst-test-network" {
	name = "inst-test-network-%s"
}
resource "google_compute_subnetwork" "inst-test-subnetwork" {
	name          = "inst-test-subnetwork-%s"
	ip_cidr_range = "10.0.0.0/16"
	region        = "us-east1"
	network       = "${google_compute_network.inst-test-network.self_link}"
	secondary_ip_range {
		range_name = "inst-test-secondary"
		ip_cidr_range = "172.16.0.0/20"
	}
}
resource "google_compute_instance_template" "foobar" {
	name = "instance-test-%s"
	machine_type = "n1-standard-1"

	disk {
		source_image = "debian-8-jessie-v20160803"
		auto_delete = true
		disk_size_gb = 10
		boot = true
	}

	metadata {
		foo = "bar"
	}

  	network_interface {
		subnetwork = "${google_compute_subnetwork.inst-test-subnetwork.self_link}"

		// Note that unlike compute instances, instance templates seem to be
		// only able to specify the netmask here. Trying a full CIDR string
		// results in: 
		// Invalid value for field 'resource.properties.networkInterfaces[0].aliasIpRanges[0].ipCidrRange':
		// '172.16.0.0/24'. Alias IP CIDR range must be a valid netmask starting with '/' (e.g. '/24')
		alias_ip_range {
			subnetwork_range_name = "${google_compute_subnetwork.inst-test-subnetwork.secondary_ip_range.0.range_name}"
			ip_cidr_range = "/24"
		}
	}
}`, i, i, i)
}

func testAccComputeInstanceTemplate_guestAccelerator(i string, count uint8) string {
	return fmt.Sprintf(`
resource "google_compute_instance_template" "foobar" {
	name = "instance-test-%s"
	machine_type = "n1-standard-1"

	disk {
		source_image = "debian-8-jessie-v20160803"
		auto_delete = true
		disk_size_gb = 10
		boot = true
	}

	network_interface {
		network = "default"
	}

	scheduling {
		# Instances with guest accelerators do not support live migration.
		on_host_maintenance = "TERMINATE"
	}

	guest_accelerator {
		count = %d
		type = "nvidia-tesla-k80"
	}
}`, i, count)
}

func testAccComputeInstanceTemplate_minCpuPlatform(i string) string {
	return fmt.Sprintf(`
resource "google_compute_instance_template" "foobar" {
	name = "instance-test-%s"
	machine_type = "n1-standard-1"

	disk {
		source_image = "debian-8-jessie-v20160803"
		auto_delete = true
		disk_size_gb = 10
		boot = true
	}

	network_interface {
		network = "default"
	}

	scheduling {
		# Instances with guest accelerators do not support live migration.
		on_host_maintenance = "TERMINATE"
	}

	min_cpu_platform = "%s"
}`, i, DEFAULT_MIN_CPU_TEST_VALUE)
}
