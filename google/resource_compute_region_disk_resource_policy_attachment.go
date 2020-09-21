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
	"log"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceComputeRegionDiskResourcePolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeRegionDiskResourcePolicyAttachmentCreate,
		Read:   resourceComputeRegionDiskResourcePolicyAttachmentRead,
		Delete: resourceComputeRegionDiskResourcePolicyAttachmentDelete,

		Importer: &schema.ResourceImporter{
			State: resourceComputeRegionDiskResourcePolicyAttachmentImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(4 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"disk": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      `The name of the regional disk in which the resource policies are attached to.`,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: `The resource policy to be attached to the disk for scheduling snapshot
creation. Do not specify the self link.`,
			},
			"region": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      `A reference to the region where the disk resides.`,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceComputeRegionDiskResourcePolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	var m providerMeta

	err := d.GetProviderMeta(&m)
	if err != nil {
		return err
	}

	config := meta.(*Config)
	config.userAgent = fmt.Sprintf("%s %s", config.userAgent, m.ModuleName)

	obj := make(map[string]interface{})
	nameProp, err := expandNestedComputeRegionDiskResourcePolicyAttachmentName(d.Get("name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("name"); !isEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}

	obj, err = resourceComputeRegionDiskResourcePolicyAttachmentEncoder(d, meta, obj)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/regions/{{region}}/disks/{{disk}}/addResourcePolicies")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new RegionDiskResourcePolicyAttachment: %#v", obj)
	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequestWithTimeout(config, "POST", billingProject, url, obj, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating RegionDiskResourcePolicyAttachment: %s", err)
	}

	// Store the ID now
	id, err := replaceVars(d, config, "{{project}}/{{region}}/{{disk}}/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = computeOperationWaitTime(
		config, res, project, "Creating RegionDiskResourcePolicyAttachment",
		d.Timeout(schema.TimeoutCreate))

	if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error waiting to create RegionDiskResourcePolicyAttachment: %s", err)
	}

	log.Printf("[DEBUG] Finished creating RegionDiskResourcePolicyAttachment %q: %#v", d.Id(), res)

	return resourceComputeRegionDiskResourcePolicyAttachmentRead(d, meta)
}

func resourceComputeRegionDiskResourcePolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	var m providerMeta

	err := d.GetProviderMeta(&m)
	if err != nil {
		return err
	}

	config := meta.(*Config)
	config.userAgent = fmt.Sprintf("%s %s", config.userAgent, m.ModuleName)

	url, err := replaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/regions/{{region}}/disks/{{disk}}")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequest(config, "GET", billingProject, url, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("ComputeRegionDiskResourcePolicyAttachment %q", d.Id()))
	}

	res, err = flattenNestedComputeRegionDiskResourcePolicyAttachment(d, meta, res)
	if err != nil {
		return err
	}

	if res == nil {
		// Object isn't there any more - remove it from the state.
		log.Printf("[DEBUG] Removing ComputeRegionDiskResourcePolicyAttachment because it couldn't be matched.")
		d.SetId("")
		return nil
	}

	res, err = resourceComputeRegionDiskResourcePolicyAttachmentDecoder(d, meta, res)
	if err != nil {
		return err
	}

	if res == nil {
		// Decoding the object has resulted in it being gone. It may be marked deleted
		log.Printf("[DEBUG] Removing ComputeRegionDiskResourcePolicyAttachment because it no longer exists.")
		d.SetId("")
		return nil
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading RegionDiskResourcePolicyAttachment: %s", err)
	}

	if err := d.Set("name", flattenNestedComputeRegionDiskResourcePolicyAttachmentName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading RegionDiskResourcePolicyAttachment: %s", err)
	}

	return nil
}

func resourceComputeRegionDiskResourcePolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	var m providerMeta

	err := d.GetProviderMeta(&m)
	if err != nil {
		return err
	}

	config := meta.(*Config)
	config.userAgent = fmt.Sprintf("%s %s", config.userAgent, m.ModuleName)

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	billingProject = project

	url, err := replaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/regions/{{region}}/disks/{{disk}}/removeResourcePolicies")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	obj = make(map[string]interface{})

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}
	if region == "" {
		return fmt.Errorf("region must be non-empty - set in resource or at provider-level")
	}

	name, err := expandNestedComputeDiskResourcePolicyAttachmentName(d.Get("name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("name"); !isEmptyValue(reflect.ValueOf(name)) && (ok || !reflect.DeepEqual(v, name)) {
		obj["resourcePolicies"] = []interface{}{fmt.Sprintf("projects/%s/regions/%s/resourcePolicies/%s", project, region, name)}
	}
	log.Printf("[DEBUG] Deleting RegionDiskResourcePolicyAttachment %q", d.Id())

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequestWithTimeout(config, "POST", billingProject, url, obj, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return handleNotFoundError(err, d, "RegionDiskResourcePolicyAttachment")
	}

	err = computeOperationWaitTime(
		config, res, project, "Deleting RegionDiskResourcePolicyAttachment",
		d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting RegionDiskResourcePolicyAttachment %q: %#v", d.Id(), res)
	return nil
}

func resourceComputeRegionDiskResourcePolicyAttachmentImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/regions/(?P<region>[^/]+)/disks/(?P<disk>[^/]+)/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<region>[^/]+)/(?P<disk>[^/]+)/(?P<name>[^/]+)",
		"(?P<region>[^/]+)/(?P<disk>[^/]+)/(?P<name>[^/]+)",
		"(?P<disk>[^/]+)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "{{project}}/{{region}}/{{disk}}/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenNestedComputeRegionDiskResourcePolicyAttachmentName(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func expandNestedComputeRegionDiskResourcePolicyAttachmentName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func resourceComputeRegionDiskResourcePolicyAttachmentEncoder(d *schema.ResourceData, meta interface{}, obj map[string]interface{}) (map[string]interface{}, error) {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return nil, err
	}
	if region == "" {
		return nil, fmt.Errorf("region must be non-empty - set in resource or at provider-level")
	}

	obj["resourcePolicies"] = []interface{}{fmt.Sprintf("projects/%s/regions/%s/resourcePolicies/%s", project, region, obj["name"])}
	delete(obj, "name")
	return obj, nil
}

func flattenNestedComputeRegionDiskResourcePolicyAttachment(d *schema.ResourceData, meta interface{}, res map[string]interface{}) (map[string]interface{}, error) {
	var v interface{}
	var ok bool

	v, ok = res["resourcePolicies"]
	if !ok || v == nil {
		return nil, nil
	}

	switch v.(type) {
	case []interface{}:
		break
	case map[string]interface{}:
		// Construct list out of single nested resource
		v = []interface{}{v}
	default:
		return nil, fmt.Errorf("expected list or map for value resourcePolicies. Actual value: %v", v)
	}

	_, item, err := resourceComputeRegionDiskResourcePolicyAttachmentFindNestedObjectInList(d, meta, v.([]interface{}))
	if err != nil {
		return nil, err
	}
	return item, nil
}

func resourceComputeRegionDiskResourcePolicyAttachmentFindNestedObjectInList(d *schema.ResourceData, meta interface{}, items []interface{}) (index int, item map[string]interface{}, err error) {
	expectedName, err := expandNestedComputeRegionDiskResourcePolicyAttachmentName(d.Get("name"), d, meta.(*Config))
	if err != nil {
		return -1, nil, err
	}
	expectedFlattenedName := flattenNestedComputeRegionDiskResourcePolicyAttachmentName(expectedName, d, meta.(*Config))

	// Search list for this resource.
	for idx, itemRaw := range items {
		if itemRaw == nil {
			continue
		}
		// List response only contains the ID - construct a response object.
		item := map[string]interface{}{
			"name": itemRaw,
		}

		// Decode list item before comparing.
		item, err := resourceComputeRegionDiskResourcePolicyAttachmentDecoder(d, meta, item)
		if err != nil {
			return -1, nil, err
		}

		itemName := flattenNestedComputeRegionDiskResourcePolicyAttachmentName(item["name"], d, meta.(*Config))
		// isEmptyValue check so that if one is nil and the other is "", that's considered a match
		if !(isEmptyValue(reflect.ValueOf(itemName)) && isEmptyValue(reflect.ValueOf(expectedFlattenedName))) && !reflect.DeepEqual(itemName, expectedFlattenedName) {
			log.Printf("[DEBUG] Skipping item with name= %#v, looking for %#v)", itemName, expectedFlattenedName)
			continue
		}
		log.Printf("[DEBUG] Found item for resource %q: %#v)", d.Id(), item)
		return idx, item, nil
	}
	return -1, nil, nil
}
func resourceComputeRegionDiskResourcePolicyAttachmentDecoder(d *schema.ResourceData, meta interface{}, res map[string]interface{}) (map[string]interface{}, error) {
	res["name"] = GetResourceNameFromSelfLink(res["name"].(string))
	return res, nil
}
