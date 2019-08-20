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

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"google.golang.org/api/compute/v1"
)

func resourceComputeNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeNetworkCreate,
		Read:   resourceComputeNetworkRead,
		Update: resourceComputeNetworkUpdate,
		Delete: resourceComputeNetworkDelete,

		Importer: &schema.ResourceImporter{
			State: resourceComputeNetworkImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Update: schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(4 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"auto_create_subnetworks": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ipv4_range": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Legacy Networks are deprecated and you will no longer be able to create them using this field from Feb 1, 2020 onwards.",
				ForceNew:   true,
			},
			"routing_mode": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"REGIONAL", "GLOBAL", ""}, false),
			},

			"gateway_ipv4": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"delete_default_routes_on_create": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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

func resourceComputeNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	obj := make(map[string]interface{})
	descriptionProp, err := expandComputeNetworkDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !isEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	IPv4RangeProp, err := expandComputeNetworkIpv4_range(d.Get("ipv4_range"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("ipv4_range"); !isEmptyValue(reflect.ValueOf(IPv4RangeProp)) && (ok || !reflect.DeepEqual(v, IPv4RangeProp)) {
		obj["IPv4Range"] = IPv4RangeProp
	}
	nameProp, err := expandComputeNetworkName(d.Get("name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("name"); !isEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	autoCreateSubnetworksProp, err := expandComputeNetworkAutoCreateSubnetworks(d.Get("auto_create_subnetworks"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("auto_create_subnetworks"); !isEmptyValue(reflect.ValueOf(autoCreateSubnetworksProp)) && (ok || !reflect.DeepEqual(v, autoCreateSubnetworksProp)) {
		obj["autoCreateSubnetworks"] = autoCreateSubnetworksProp
	}
	routingConfigProp, err := expandComputeNetworkRoutingConfig(nil, d, config)
	if err != nil {
		return err
	} else if !isEmptyValue(reflect.ValueOf(routingConfigProp)) {
		obj["routingConfig"] = routingConfigProp
	}

	obj, err = resourceComputeNetworkEncoder(d, meta, obj)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/global/networks")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new Network: %#v", obj)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	res, err := sendRequestWithTimeout(config, "POST", project, url, obj, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating Network: %s", err)
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
		config.clientCompute, op, project, "Creating Network",
		int(d.Timeout(schema.TimeoutCreate).Minutes()))

	if waitErr != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error waiting to create Network: %s", waitErr)
	}

	log.Printf("[DEBUG] Finished creating Network %q: %#v", d.Id(), res)

	if d.Get("delete_default_routes_on_create").(bool) {
		token := ""
		for paginate := true; paginate; {
			networkLink := fmt.Sprintf("%s/%s", url, d.Get("name").(string))
			filter := fmt.Sprintf("(network=\"%s\") AND (destRange=\"0.0.0.0/0\")", networkLink)
			log.Printf("[DEBUG] Getting routes for network %q with filter '%q'", d.Get("name").(string), filter)
			resp, err := config.clientCompute.Routes.List(project).Filter(filter).Do()
			if err != nil {
				return fmt.Errorf("Error listing routes in proj: %s", err)
			}

			log.Printf("[DEBUG] Found %d routes rules in %q network", len(resp.Items), d.Get("name").(string))

			for _, route := range resp.Items {
				op, err := config.clientCompute.Routes.Delete(project, route.Name).Do()
				if err != nil {
					return fmt.Errorf("Error deleting route: %s", err)
				}
				err = computeSharedOperationWait(config.clientCompute, op, project, "Deleting Route")
				if err != nil {
					return err
				}
			}

			token = resp.NextPageToken
			paginate = token != ""
		}
	}

	return resourceComputeNetworkRead(d, meta)
}

func resourceComputeNetworkRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	url, err := replaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/global/networks/{{name}}")
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	res, err := sendRequest(config, "GET", project, url, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("ComputeNetwork %q", d.Id()))
	}

	// Explicitly set virtual fields to default values if unset
	if _, ok := d.GetOk("delete_default_routes_on_create"); !ok {
		d.Set("delete_default_routes_on_create", false)
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading Network: %s", err)
	}

	if err := d.Set("description", flattenComputeNetworkDescription(res["description"], d)); err != nil {
		return fmt.Errorf("Error reading Network: %s", err)
	}
	if err := d.Set("gateway_ipv4", flattenComputeNetworkGateway_ipv4(res["gatewayIPv4"], d)); err != nil {
		return fmt.Errorf("Error reading Network: %s", err)
	}
	if err := d.Set("ipv4_range", flattenComputeNetworkIpv4_range(res["IPv4Range"], d)); err != nil {
		return fmt.Errorf("Error reading Network: %s", err)
	}
	if err := d.Set("name", flattenComputeNetworkName(res["name"], d)); err != nil {
		return fmt.Errorf("Error reading Network: %s", err)
	}
	if err := d.Set("auto_create_subnetworks", flattenComputeNetworkAutoCreateSubnetworks(res["autoCreateSubnetworks"], d)); err != nil {
		return fmt.Errorf("Error reading Network: %s", err)
	}
	// Terraform must set the top level schema field, but since this object contains collapsed properties
	// it's difficult to know what the top level should be. Instead we just loop over the map returned from flatten.
	if flattenedProp := flattenComputeNetworkRoutingConfig(res["routingConfig"], d); flattenedProp != nil {
		casted := flattenedProp.([]interface{})[0]
		if casted != nil {
			for k, v := range casted.(map[string]interface{}) {
				d.Set(k, v)
			}
		}
	}
	if err := d.Set("self_link", ConvertSelfLinkToV1(res["selfLink"].(string))); err != nil {
		return fmt.Errorf("Error reading Network: %s", err)
	}

	return nil
}

func resourceComputeNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	d.Partial(true)

	if d.HasChange("routing_mode") {
		obj := make(map[string]interface{})
		routingConfigProp, err := expandComputeNetworkRoutingConfig(nil, d, config)
		if err != nil {
			return err
		} else if v, ok := d.GetOkExists("routing_config"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, routingConfigProp)) {
			obj["routingConfig"] = routingConfigProp
		}

		url, err := replaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/global/networks/{{name}}")
		if err != nil {
			return err
		}
		res, err := sendRequestWithTimeout(config, "PATCH", project, url, obj, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("Error updating Network %q: %s", d.Id(), err)
		}

		op := &compute.Operation{}
		err = Convert(res, op)
		if err != nil {
			return err
		}

		err = computeOperationWaitTime(
			config.clientCompute, op, project, "Updating Network",
			int(d.Timeout(schema.TimeoutUpdate).Minutes()))

		if err != nil {
			return err
		}

		d.SetPartial("routing_mode")
	}

	d.Partial(false)

	return resourceComputeNetworkRead(d, meta)
}

func resourceComputeNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/global/networks/{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	log.Printf("[DEBUG] Deleting Network %q", d.Id())

	res, err := sendRequestWithTimeout(config, "DELETE", project, url, obj, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return handleNotFoundError(err, d, "Network")
	}

	op := &compute.Operation{}
	err = Convert(res, op)
	if err != nil {
		return err
	}

	err = computeOperationWaitTime(
		config.clientCompute, op, project, "Deleting Network",
		int(d.Timeout(schema.TimeoutDelete).Minutes()))

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting Network %q: %#v", d.Id(), res)
	return nil
}

func resourceComputeNetworkImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/global/networks/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<name>[^/]+)",
		"(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	// Explicitly set virtual fields to default values on import
	d.Set("delete_default_routes_on_create", false)

	return []*schema.ResourceData{d}, nil
}

func flattenComputeNetworkDescription(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenComputeNetworkGateway_ipv4(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenComputeNetworkIpv4_range(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenComputeNetworkName(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenComputeNetworkAutoCreateSubnetworks(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenComputeNetworkRoutingConfig(v interface{}, d *schema.ResourceData) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["routing_mode"] =
		flattenComputeNetworkRoutingConfigRoutingMode(original["routingMode"], d)
	return []interface{}{transformed}
}
func flattenComputeNetworkRoutingConfigRoutingMode(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandComputeNetworkDescription(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeNetworkIpv4_range(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeNetworkName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeNetworkAutoCreateSubnetworks(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeNetworkRoutingConfig(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	transformed := make(map[string]interface{})
	transformedRoutingMode, err := expandComputeNetworkRoutingConfigRoutingMode(d.Get("routing_mode"), d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedRoutingMode); val.IsValid() && !isEmptyValue(val) {
		transformed["routingMode"] = transformedRoutingMode
	}

	return transformed, nil
}

func expandComputeNetworkRoutingConfigRoutingMode(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func resourceComputeNetworkEncoder(d *schema.ResourceData, meta interface{}, obj map[string]interface{}) (map[string]interface{}, error) {
	if _, ok := d.GetOk("ipv4_range"); !ok {
		obj["autoCreateSubnetworks"] = d.Get("auto_create_subnetworks")
	}

	return obj, nil
}
