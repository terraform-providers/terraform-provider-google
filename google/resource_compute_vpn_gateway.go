// TODO(https://github.com/GoogleCloudPlatform/magic-modules/issues/156): Re-enable code generation for this resource
package google

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	compute "google.golang.org/api/compute/v1"
)

func resourceComputeVpnGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeVpnGatewayCreate,
		Read:   resourceComputeVpnGatewayRead,
		Delete: resourceComputeVpnGatewayDelete,

		Importer: &schema.ResourceImporter{
			State: resourceComputeVpnGatewayImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(240 * time.Second),
			Delete: schema.DefaultTimeout(240 * time.Second),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"region": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				StateFunc:        NameFromSelfLinkStateFunc,
			},
			"creation_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceComputeVpnGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	descriptionProp, err := expandComputeVpnGatewayDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	}
	nameProp, err := expandComputeVpnGatewayName(d.Get("name"), d, config)
	if err != nil {
		return err
	}
	networkProp, err := expandComputeVpnGatewayNetwork(d.Get("network"), d, config)
	if err != nil {
		return err
	}
	regionProp, err := expandComputeVpnGatewayRegion(d.Get("region"), d, config)
	if err != nil {
		return err
	}

	obj := map[string]interface{}{
		"description": descriptionProp,
		"name":        nameProp,
		"network":     networkProp,
		"region":      regionProp,
	}

	url, err := replaceVars(d, config, "https://www.googleapis.com/compute/v1/projects/{{project}}/regions/{{region}}/targetVpnGateways")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new VpnGateway: %#v", obj)
	res, err := Post(config, url, obj)
	if err != nil {
		return fmt.Errorf("Error creating VpnGateway: %s", err)
	}

	// Store the ID now
	id, err := replaceVars(d, config, "{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	op := &compute.Operation{}
	err = Convert(res, op)
	if err != nil {
		return err
	}

	waitErr := computeOperationWaitTime(
		config.clientCompute, op, project, "Creating VpnGateway",
		int(d.Timeout(schema.TimeoutCreate).Minutes()))

	if waitErr != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error waiting to create VpnGateway: %s", waitErr)
	}

	return resourceComputeVpnGatewayRead(d, meta)
}

func resourceComputeVpnGatewayRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "https://www.googleapis.com/compute/v1/projects/{{project}}/regions/{{region}}/targetVpnGateways/{{name}}")
	if err != nil {
		return err
	}

	res, err := Get(config, url)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("ComputeVpnGateway %q", d.Id()))
	}
	if err := d.Set("creation_timestamp", flattenComputeVpnGatewayCreationTimestamp(res["creationTimestamp"])); err != nil {
		return fmt.Errorf("Error reading VpnGateway: %s", err)
	}
	if err := d.Set("description", flattenComputeVpnGatewayDescription(res["description"])); err != nil {
		return fmt.Errorf("Error reading VpnGateway: %s", err)
	}
	if err := d.Set("name", flattenComputeVpnGatewayName(res["name"])); err != nil {
		return fmt.Errorf("Error reading VpnGateway: %s", err)
	}
	if err := d.Set("network", flattenComputeVpnGatewayNetwork(res["network"])); err != nil {
		return fmt.Errorf("Error reading VpnGateway: %s", err)
	}
	if err := d.Set("region", flattenComputeVpnGatewayRegion(res["region"])); err != nil {
		return fmt.Errorf("Error reading VpnGateway: %s", err)
	}
	if err := d.Set("self_link", res["selfLink"]); err != nil {
		return fmt.Errorf("Error reading VpnGateway: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading VpnGateway: %s", err)
	}

	return nil
}

func resourceComputeVpnGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "https://www.googleapis.com/compute/v1/projects/{{project}}/regions/{{region}}/targetVpnGateways/{{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting VpnGateway %q", d.Id())
	res, err := Delete(config, url)
	if err != nil {
		return fmt.Errorf("Error deleting VpnGateway %q: %s", d.Id(), err)
	}

	op := &compute.Operation{}
	err = Convert(res, op)
	if err != nil {
		return err
	}

	err = computeOperationWaitTime(
		config.clientCompute, op, project, "Deleting VpnGateway",
		int(d.Timeout(schema.TimeoutDelete).Minutes()))

	if err != nil {
		return err
	}

	return nil
}

func resourceComputeVpnGatewayImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	parseImportId([]string{"projects/(?P<project>[^/]+)/regions/(?P<region>[^/]+)/targetVpnGateways/(?P<name>[^/]+)", "(?P<project>[^/]+)/(?P<region>[^/]+)/(?P<name>[^/]+)", "(?P<name>[^/]+)"}, d, config)

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenComputeVpnGatewayCreationTimestamp(v interface{}) interface{} {
	return v
}

func flattenComputeVpnGatewayDescription(v interface{}) interface{} {
	return v
}

func flattenComputeVpnGatewayName(v interface{}) interface{} {
	return v
}

func flattenComputeVpnGatewayNetwork(v interface{}) interface{} {
	return v
}

func flattenComputeVpnGatewayRegion(v interface{}) interface{} {
	return NameFromSelfLinkStateFunc(v)
}

func expandComputeVpnGatewayDescription(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeVpnGatewayName(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeVpnGatewayNetwork(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	f, err := parseGlobalFieldValue("networks", v.(string), "project", d, config, true)
	if err != nil {
		return nil, fmt.Errorf("Invalid value for network: %s", err)
	}
	return f.RelativeLink(), nil
}

func expandComputeVpnGatewayRegion(v interface{}, d *schema.ResourceData, config *Config) (interface{}, error) {
	f, err := parseGlobalFieldValue("regions", v.(string), "project", d, config, true)
	if err != nil {
		return nil, fmt.Errorf("Invalid value for region: %s", err)
	}
	return f.RelativeLink(), nil
}
