// Contains functions that don't really belong anywhere else.

package google

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	computeBeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

type TerraformResourceData interface {
	HasChange(string) bool
	GetOk(string) (interface{}, bool)
	Set(string, interface{}) error
	SetId(string)
	Id() string
}

// getRegionFromZone returns the region from a zone for Google cloud.
func getRegionFromZone(zone string) string {
	if zone != "" && len(zone) > 2 {
		region := zone[:len(zone)-2]
		return region
	}
	return ""
}

// Infers the region based on the following (in order of priority):
// - `region` field in resource schema
// - region extracted from the `zone` field in resource schema
// - provider-level region
// - region extracted from the provider-level zone
func getRegion(d TerraformResourceData, config *Config) (string, error) {
	return getRegionFromSchema("region", "zone", d, config)
}

func getRegionFromInstanceState(is *terraform.InstanceState, config *Config) (string, error) {
	res, ok := is.Attributes["region"]

	if ok && res != "" {
		return res, nil
	}

	if config.Region != "" {
		return config.Region, nil
	}

	return "", fmt.Errorf("region: required field is not set")
}

// getProject reads the "project" field from the given resource data and falls
// back to the provider's value if not given. If the provider's value is not
// given, an error is returned.
func getProject(d TerraformResourceData, config *Config) (string, error) {
	return getProjectFromSchema("project", d, config)
}

func getProjectFromInstanceState(is *terraform.InstanceState, config *Config) (string, error) {
	res, ok := is.Attributes["project"]

	if ok && res != "" {
		return res, nil
	}

	if config.Project != "" {
		return config.Project, nil
	}

	return "", fmt.Errorf("project: required field is not set")
}

func getZonalResourceFromRegion(getResource func(string) (interface{}, error), region string, compute *compute.Service, project string) (interface{}, error) {
	zoneList, err := compute.Zones.List(project).Do()
	if err != nil {
		return nil, err
	}
	var resource interface{}
	for _, zone := range zoneList.Items {
		if strings.Contains(zone.Name, region) {
			resource, err = getResource(zone.Name)
			if err != nil {
				if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
					// Resource was not found in this zone
					continue
				}
				return nil, fmt.Errorf("Error reading Resource: %s", err)
			}
			// Resource was found
			return resource, nil
		}
	}
	// Resource does not exist in this region
	return nil, nil
}

func getZonalBetaResourceFromRegion(getResource func(string) (interface{}, error), region string, compute *computeBeta.Service, project string) (interface{}, error) {
	zoneList, err := compute.Zones.List(project).Do()
	if err != nil {
		return nil, err
	}
	var resource interface{}
	for _, zone := range zoneList.Items {
		if strings.Contains(zone.Name, region) {
			resource, err = getResource(zone.Name)
			if err != nil {
				if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
					// Resource was not found in this zone
					continue
				}
				return nil, fmt.Errorf("Error reading Resource: %s", err)
			}
			// Resource was found
			return resource, nil
		}
	}
	// Resource does not exist in this region
	return nil, nil
}

func getRouterLockName(region string, router string) string {
	return fmt.Sprintf("router/%s/%s", region, router)
}

func handleNotFoundError(err error, d *schema.ResourceData, resource string) error {
	if isGoogleApiErrorWithCode(err, 404) {
		log.Printf("[WARN] Removing %s because it's gone", resource)
		// The resource doesn't exist anymore
		d.SetId("")

		return nil
	}

	return fmt.Errorf("Error reading %s: %s", resource, err)
}

func isGoogleApiErrorWithCode(err error, errCode int) bool {
	gerr, ok := errwrap.GetType(err, &googleapi.Error{}).(*googleapi.Error)
	return ok && gerr != nil && gerr.Code == errCode
}

func isApiNotEnabledError(err error) bool {
	gerr, ok := errwrap.GetType(err, &googleapi.Error{}).(*googleapi.Error)
	if !ok {
		return false
	}
	if gerr == nil {
		return false
	}
	if gerr.Code != 403 {
		return false
	}
	for _, e := range gerr.Errors {
		if e.Reason == "accessNotConfigured" {
			return true
		}
	}
	return false
}

func isConflictError(err error) bool {
	if e, ok := err.(*googleapi.Error); ok && e.Code == 409 {
		return true
	} else if !ok && errwrap.ContainsType(err, &googleapi.Error{}) {
		e := errwrap.GetType(err, &googleapi.Error{}).(*googleapi.Error)
		if e.Code == 409 {
			return true
		}
	}
	return false
}

func linkDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	if GetResourceNameFromSelfLink(old) == new {
		return true
	}
	return false
}

func optionalPrefixSuppress(prefix string) schema.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *schema.ResourceData) bool {
		return prefix+old == new || prefix+new == old
	}
}

func optionalSurroundingSpacesSuppress(k, old, new string, d *schema.ResourceData) bool {
	return strings.TrimSpace(old) == strings.TrimSpace(new)
}

func emptyOrDefaultStringSuppress(defaultVal string) schema.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *schema.ResourceData) bool {
		return (old == "" && new == defaultVal) || (new == "" && old == defaultVal)
	}
}

func ipCidrRangeDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	// The range may be a:
	// A) single IP address (e.g. 10.2.3.4)
	// B) CIDR format string (e.g. 10.1.2.0/24)
	// C) netmask (e.g. /24)
	//
	// For A) and B), no diff to suppress, they have to match completely.
	// For C), The API picks a network IP address and this creates a diff of the form:
	// network_interface.0.alias_ip_range.0.ip_cidr_range: "10.128.1.0/24" => "/24"
	// We should only compare the mask portion for this case.
	if len(new) > 0 && new[0] == '/' {
		oldNetmaskStartPos := strings.LastIndex(old, "/")

		if oldNetmaskStartPos != -1 {
			oldNetmask := old[strings.LastIndex(old, "/"):]
			if oldNetmask == new {
				return true
			}
		}
	}

	return false
}

func caseDiffSuppress(_, old, new string, _ *schema.ResourceData) bool {
	return strings.ToUpper(old) == strings.ToUpper(new)
}

// Port range '80' and '80-80' is equivalent.
// `old` is read from the server and always has the full range format (e.g. '80-80', '1024-2048').
// `new` can be either a single port or a port range.
func portRangeDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	if old == new+"-"+new {
		return true
	}
	return false
}

// Single-digit hour is equivalent to hour with leading zero e.g. suppress diff 1:00 => 01:00.
// Assume either value could be in either format.
func rfc3339TimeDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	if (len(old) == 4 && "0"+old == new) || (len(new) == 4 && "0"+new == old) {
		return true
	}
	return false
}

// expandLabels pulls the value of "labels" out of a schema.ResourceData as a map[string]string.
func expandLabels(d *schema.ResourceData) map[string]string {
	return expandStringMap(d, "labels")
}

// expandStringMap pulls the value of key out of a schema.ResourceData as a map[string]string.
func expandStringMap(d *schema.ResourceData, key string) map[string]string {
	v, ok := d.GetOk(key)

	if !ok {
		return map[string]string{}
	}

	return convertStringMap(v.(map[string]interface{}))
}

func convertStringMap(v map[string]interface{}) map[string]string {
	m := make(map[string]string)
	for k, val := range v {
		m[k] = val.(string)
	}
	return m
}

func convertStringArr(ifaceArr []interface{}) []string {
	return convertAndMapStringArr(ifaceArr, func(s string) string { return s })
}

func convertAndMapStringArr(ifaceArr []interface{}, f func(string) string) []string {
	var arr []string
	for _, v := range ifaceArr {
		if v == nil {
			continue
		}
		arr = append(arr, f(v.(string)))
	}
	return arr
}

func convertStringArrToInterface(strs []string) []interface{} {
	arr := make([]interface{}, len(strs))
	for i, str := range strs {
		arr[i] = str
	}
	return arr
}

func convertStringSet(set *schema.Set) []string {
	s := make([]string, 0, set.Len())
	for _, v := range set.List() {
		s = append(s, v.(string))
	}
	return s
}

func mergeSchemas(a, b map[string]*schema.Schema) map[string]*schema.Schema {
	merged := make(map[string]*schema.Schema)

	for k, v := range a {
		merged[k] = v
	}

	for k, v := range b {
		merged[k] = v
	}

	return merged
}

func mergeResourceMaps(ms ...map[string]*schema.Resource) map[string]*schema.Resource {
	merged := make(map[string]*schema.Resource)

	for _, m := range ms {
		for k, v := range m {
			merged[k] = v
		}
	}

	return merged
}

func retry(retryFunc func() error) error {
	return retryTime(retryFunc, 1)
}

func retryTime(retryFunc func() error, minutes int) error {
	return resource.Retry(time.Duration(minutes)*time.Minute, func() *resource.RetryError {
		err := retryFunc()
		if err == nil {
			return nil
		}
		if gerr, ok := err.(*googleapi.Error); ok && (gerr.Code == 429 || gerr.Code == 500 || gerr.Code == 502 || gerr.Code == 503) {
			return resource.RetryableError(gerr)
		}
		return resource.NonRetryableError(err)
	})
}

func extractFirstMapConfig(m []interface{}) map[string]interface{} {
	if len(m) == 0 {
		return map[string]interface{}{}
	}

	return m[0].(map[string]interface{})
}

func lockedCall(lockKey string, f func() error) error {
	mutexKV.Lock(lockKey)
	defer mutexKV.Unlock(lockKey)

	return f()
}

// serviceAccountFQN will attempt to generate the fully qualified name in the format of:
// "projects/(-|<project_name>)/serviceAccounts/<service_account_id>@<project_name>.iam.gserviceaccount.com"
func serviceAccountFQN(serviceAccount, project string) string {
	// If the service account id isn't already the fully qualified name
	if !strings.HasPrefix(serviceAccount, "projects/") {
		// If the service account id is an email
		if strings.Contains(serviceAccount, "@") {
			serviceAccount = "projects/-/serviceAccounts/" + serviceAccount
		} else {
			// If the service account id doesn't contain the email, build it
			serviceAccount = fmt.Sprintf("projects/-/serviceAccounts/%s@%s.iam.gserviceaccount.com", serviceAccount, project)
		}
	}
	return serviceAccount
}
