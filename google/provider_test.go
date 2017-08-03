package google

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"google": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("GOOGLE_CREDENTIALS_FILE"); v != "" {
		creds, err := ioutil.ReadFile(v)
		if err != nil {
			t.Fatalf("Error reading GOOGLE_CREDENTIALS_FILE path: %s", err)
		}
		os.Setenv("GOOGLE_CREDENTIALS", string(creds))
	}

	multiEnvSearch := func(ks []string) string {
		for _, k := range ks {
			if v := os.Getenv(k); v != "" {
				return v
			}
		}
		return ""
	}

	creds := []string{
		"GOOGLE_CREDENTIALS",
		"GOOGLE_CLOUD_KEYFILE_JSON",
		"GCLOUD_KEYFILE_JSON",
		"GOOGLE_USE_DEFAULT_CREDENTIALS",
	}
	if v := multiEnvSearch(creds); v == "" {
		t.Fatalf("One of %s must be set for acceptance tests", strings.Join(creds, ", "))
	}

	projs := []string{
		"GOOGLE_PROJECT",
		"GCLOUD_PROJECT",
		"CLOUDSDK_CORE_PROJECT",
	}
	if v := multiEnvSearch(projs); v == "" {
		t.Fatalf("One of %s must be set for acceptance tests", strings.Join(projs, ", "))
	}

	regs := []string{
		"GOOGLE_REGION",
		"GCLOUD_REGION",
		"CLOUDSDK_COMPUTE_REGION",
	}
	if v := multiEnvSearch(regs); v != "us-central1" {
		t.Fatalf("One of %s must be set to us-central1 for acceptance tests", strings.Join(regs, ", "))
	}

	if v := os.Getenv("GOOGLE_XPN_HOST_PROJECT"); v == "" {
		t.Fatal("GOOGLE_XPN_HOST_PROJECT must be set for acceptance tests")
	}
}

func TestProvider_getRegionFromZone(t *testing.T) {
	expected := "us-central1"
	actual := getRegionFromZone("us-central1-f")
	if expected != actual {
		t.Fatalf("Region (%s) did not match expected value: %s", actual, expected)
	}
}

// getTestRegion has the same logic as the provider's getRegion, to be used in tests.
func getTestRegion(is *terraform.InstanceState, config *Config) (string, error) {
	if res, ok := is.Attributes["region"]; ok {
		return res, nil
	}
	if config.Region != "" {
		return config.Region, nil
	}
	return "", fmt.Errorf("%q: required field is not set", "region")
}

// getTestProject has the same logic as the provider's getProject, to be used in tests.
func getTestProject(is *terraform.InstanceState, config *Config) (string, error) {
	if res, ok := is.Attributes["project"]; ok {
		return res, nil
	}
	if config.Project != "" {
		return config.Project, nil
	}
	return "", fmt.Errorf("%q: required field is not set", "project")
}

func TestConvertStringArr(t *testing.T) {

	input := make([]interface{}, 3)
	input[0] = "aaa"
	input[1] = "bbb"
	input[2] = "aaa"

	expected := []string{"aaa", "bbb", "ccc"}
	actual := convertStringArr(input)

	if reflect.DeepEqual(expected, actual) {
		t.Fatalf("(%s) did not match expected value: %s", actual, expected)
	}
}

func TestConvertAndMapStringArr(t *testing.T) {

	input := make([]interface{}, 3)
	input[0] = "aaa"
	input[1] = "bbb"
	input[2] = "aaa"

	expected := []string{"AAA", "BBB", "CCC"}
	actual := convertAndMapStringArr(input, strings.ToUpper)

	if reflect.DeepEqual(expected, actual) {
		t.Fatalf("(%s) did not match expected value: %s", actual, expected)
	}
}
