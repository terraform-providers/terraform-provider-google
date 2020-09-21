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
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceStorageDefaultObjectAccessControl() *schema.Resource {
	return &schema.Resource{
		Create: resourceStorageDefaultObjectAccessControlCreate,
		Read:   resourceStorageDefaultObjectAccessControlRead,
		Update: resourceStorageDefaultObjectAccessControlUpdate,
		Delete: resourceStorageDefaultObjectAccessControlDelete,

		Importer: &schema.ResourceImporter{
			State: resourceStorageDefaultObjectAccessControlImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Update: schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(4 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      `The name of the bucket.`,
			},
			"entity": {
				Type:     schema.TypeString,
				Required: true,
				Description: `The entity holding the permission, in one of the following forms:
  * user-{{userId}}
  * user-{{email}} (such as "user-liz@example.com")
  * group-{{groupId}}
  * group-{{email}} (such as "group-example@googlegroups.com")
  * domain-{{domain}} (such as "domain-example.com")
  * project-team-{{projectId}}
  * allUsers
  * allAuthenticatedUsers`,
			},
			"role": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"OWNER", "READER"}, false),
				Description:  `The access permission for the entity. Possible values: ["OWNER", "READER"]`,
			},
			"object": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the object, if applied to an object.`,
			},
			"domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The domain associated with the entity.`,
			},
			"email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The email address associated with the entity.`,
			},
			"entity_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID for the entity`,
			},
			"generation": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The content generation of the object, if applied to an object.`,
			},
			"project_team": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The project team associated with the entity`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_number": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The project team associated with the entity`,
						},
						"team": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"editors", "owners", "viewers", ""}, false),
							Description:  `The team. Possible values: ["editors", "owners", "viewers"]`,
						},
					},
				},
			},
		},
	}
}

func resourceStorageDefaultObjectAccessControlCreate(d *schema.ResourceData, meta interface{}) error {
	var m providerMeta

	err := d.GetProviderMeta(&m)
	if err != nil {
		return err
	}

	config := meta.(*Config)
	config.userAgent = fmt.Sprintf("%s %s", config.userAgent, m.ModuleName)

	obj := make(map[string]interface{})
	bucketProp, err := expandStorageDefaultObjectAccessControlBucket(d.Get("bucket"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("bucket"); !isEmptyValue(reflect.ValueOf(bucketProp)) && (ok || !reflect.DeepEqual(v, bucketProp)) {
		obj["bucket"] = bucketProp
	}
	entityProp, err := expandStorageDefaultObjectAccessControlEntity(d.Get("entity"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("entity"); !isEmptyValue(reflect.ValueOf(entityProp)) && (ok || !reflect.DeepEqual(v, entityProp)) {
		obj["entity"] = entityProp
	}
	objectProp, err := expandStorageDefaultObjectAccessControlObject(d.Get("object"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("object"); !isEmptyValue(reflect.ValueOf(objectProp)) && (ok || !reflect.DeepEqual(v, objectProp)) {
		obj["object"] = objectProp
	}
	roleProp, err := expandStorageDefaultObjectAccessControlRole(d.Get("role"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("role"); !isEmptyValue(reflect.ValueOf(roleProp)) && (ok || !reflect.DeepEqual(v, roleProp)) {
		obj["role"] = roleProp
	}

	lockName, err := replaceVars(d, config, "storage/buckets/{{bucket}}")
	if err != nil {
		return err
	}
	mutexKV.Lock(lockName)
	defer mutexKV.Unlock(lockName)

	url, err := replaceVars(d, config, "{{StorageBasePath}}b/{{bucket}}/defaultObjectAcl")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new DefaultObjectAccessControl: %#v", obj)
	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequestWithTimeout(config, "POST", billingProject, url, obj, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating DefaultObjectAccessControl: %s", err)
	}

	// Store the ID now
	id, err := replaceVars(d, config, "{{bucket}}/{{entity}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating DefaultObjectAccessControl %q: %#v", d.Id(), res)

	return resourceStorageDefaultObjectAccessControlRead(d, meta)
}

func resourceStorageDefaultObjectAccessControlRead(d *schema.ResourceData, meta interface{}) error {
	var m providerMeta

	err := d.GetProviderMeta(&m)
	if err != nil {
		return err
	}

	config := meta.(*Config)
	config.userAgent = fmt.Sprintf("%s %s", config.userAgent, m.ModuleName)

	url, err := replaceVars(d, config, "{{StorageBasePath}}b/{{bucket}}/defaultObjectAcl/{{entity}}")
	if err != nil {
		return err
	}

	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequest(config, "GET", billingProject, url, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("StorageDefaultObjectAccessControl %q", d.Id()))
	}

	if err := d.Set("domain", flattenStorageDefaultObjectAccessControlDomain(res["domain"], d, config)); err != nil {
		return fmt.Errorf("Error reading DefaultObjectAccessControl: %s", err)
	}
	if err := d.Set("email", flattenStorageDefaultObjectAccessControlEmail(res["email"], d, config)); err != nil {
		return fmt.Errorf("Error reading DefaultObjectAccessControl: %s", err)
	}
	if err := d.Set("entity", flattenStorageDefaultObjectAccessControlEntity(res["entity"], d, config)); err != nil {
		return fmt.Errorf("Error reading DefaultObjectAccessControl: %s", err)
	}
	if err := d.Set("entity_id", flattenStorageDefaultObjectAccessControlEntityId(res["entityId"], d, config)); err != nil {
		return fmt.Errorf("Error reading DefaultObjectAccessControl: %s", err)
	}
	if err := d.Set("generation", flattenStorageDefaultObjectAccessControlGeneration(res["generation"], d, config)); err != nil {
		return fmt.Errorf("Error reading DefaultObjectAccessControl: %s", err)
	}
	if err := d.Set("object", flattenStorageDefaultObjectAccessControlObject(res["object"], d, config)); err != nil {
		return fmt.Errorf("Error reading DefaultObjectAccessControl: %s", err)
	}
	if err := d.Set("project_team", flattenStorageDefaultObjectAccessControlProjectTeam(res["projectTeam"], d, config)); err != nil {
		return fmt.Errorf("Error reading DefaultObjectAccessControl: %s", err)
	}
	if err := d.Set("role", flattenStorageDefaultObjectAccessControlRole(res["role"], d, config)); err != nil {
		return fmt.Errorf("Error reading DefaultObjectAccessControl: %s", err)
	}

	return nil
}

func resourceStorageDefaultObjectAccessControlUpdate(d *schema.ResourceData, meta interface{}) error {
	var m providerMeta

	err := d.GetProviderMeta(&m)
	if err != nil {
		return err
	}

	config := meta.(*Config)
	config.userAgent = fmt.Sprintf("%s %s", config.userAgent, m.ModuleName)

	billingProject := ""

	obj := make(map[string]interface{})
	bucketProp, err := expandStorageDefaultObjectAccessControlBucket(d.Get("bucket"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("bucket"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, bucketProp)) {
		obj["bucket"] = bucketProp
	}
	entityProp, err := expandStorageDefaultObjectAccessControlEntity(d.Get("entity"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("entity"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, entityProp)) {
		obj["entity"] = entityProp
	}
	objectProp, err := expandStorageDefaultObjectAccessControlObject(d.Get("object"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("object"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, objectProp)) {
		obj["object"] = objectProp
	}
	roleProp, err := expandStorageDefaultObjectAccessControlRole(d.Get("role"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("role"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, roleProp)) {
		obj["role"] = roleProp
	}

	lockName, err := replaceVars(d, config, "storage/buckets/{{bucket}}")
	if err != nil {
		return err
	}
	mutexKV.Lock(lockName)
	defer mutexKV.Unlock(lockName)

	url, err := replaceVars(d, config, "{{StorageBasePath}}b/{{bucket}}/defaultObjectAcl/{{entity}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating DefaultObjectAccessControl %q: %#v", d.Id(), obj)

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequestWithTimeout(config, "PUT", billingProject, url, obj, d.Timeout(schema.TimeoutUpdate))

	if err != nil {
		return fmt.Errorf("Error updating DefaultObjectAccessControl %q: %s", d.Id(), err)
	} else {
		log.Printf("[DEBUG] Finished updating DefaultObjectAccessControl %q: %#v", d.Id(), res)
	}

	return resourceStorageDefaultObjectAccessControlRead(d, meta)
}

func resourceStorageDefaultObjectAccessControlDelete(d *schema.ResourceData, meta interface{}) error {
	var m providerMeta

	err := d.GetProviderMeta(&m)
	if err != nil {
		return err
	}

	config := meta.(*Config)
	config.userAgent = fmt.Sprintf("%s %s", config.userAgent, m.ModuleName)

	billingProject := ""

	lockName, err := replaceVars(d, config, "storage/buckets/{{bucket}}")
	if err != nil {
		return err
	}
	mutexKV.Lock(lockName)
	defer mutexKV.Unlock(lockName)

	url, err := replaceVars(d, config, "{{StorageBasePath}}b/{{bucket}}/defaultObjectAcl/{{entity}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	log.Printf("[DEBUG] Deleting DefaultObjectAccessControl %q", d.Id())

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequestWithTimeout(config, "DELETE", billingProject, url, obj, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return handleNotFoundError(err, d, "DefaultObjectAccessControl")
	}

	log.Printf("[DEBUG] Finished deleting DefaultObjectAccessControl %q: %#v", d.Id(), res)
	return nil
}

func resourceStorageDefaultObjectAccessControlImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"(?P<bucket>[^/]+)/(?P<entity>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "{{bucket}}/{{entity}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenStorageDefaultObjectAccessControlDomain(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenStorageDefaultObjectAccessControlEmail(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenStorageDefaultObjectAccessControlEntity(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenStorageDefaultObjectAccessControlEntityId(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenStorageDefaultObjectAccessControlGeneration(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := strconv.ParseInt(strVal, 10, 64); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenStorageDefaultObjectAccessControlObject(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenStorageDefaultObjectAccessControlProjectTeam(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["project_number"] =
		flattenStorageDefaultObjectAccessControlProjectTeamProjectNumber(original["projectNumber"], d, config)
	transformed["team"] =
		flattenStorageDefaultObjectAccessControlProjectTeamTeam(original["team"], d, config)
	return []interface{}{transformed}
}
func flattenStorageDefaultObjectAccessControlProjectTeamProjectNumber(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenStorageDefaultObjectAccessControlProjectTeamTeam(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenStorageDefaultObjectAccessControlRole(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func expandStorageDefaultObjectAccessControlBucket(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandStorageDefaultObjectAccessControlEntity(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandStorageDefaultObjectAccessControlObject(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandStorageDefaultObjectAccessControlRole(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}
