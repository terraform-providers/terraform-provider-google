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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceStorageHmacKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceStorageHmacKeyCreate,
		Read:   resourceStorageHmacKeyRead,
		Update: resourceStorageHmacKeyUpdate,
		Delete: resourceStorageHmacKeyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceStorageHmacKeyImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Update: schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(4 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"service_account_email": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The email address of the key's associated service account.`,
			},
			"state": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ACTIVE", "INACTIVE", ""}, false),
				Description:  `The state of the key. Can be set to one of ACTIVE, INACTIVE. Default value: "ACTIVE" Possible values: ["ACTIVE", "INACTIVE"]`,
				Default:      "ACTIVE",
			},
			"access_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The access ID of the HMAC Key.`,
			},
			"secret": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `HMAC secret key material.`,
				Sensitive:   true,
			},
			"time_created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `'The creation time of the HMAC key in RFC 3339 format. '`,
			},
			"updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `'The last modification time of the HMAC key metadata in RFC 3339 format.'`,
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

func resourceStorageHmacKeyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	serviceAccountEmailProp, err := expandStorageHmacKeyServiceAccountEmail(d.Get("service_account_email"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("service_account_email"); !isEmptyValue(reflect.ValueOf(serviceAccountEmailProp)) && (ok || !reflect.DeepEqual(v, serviceAccountEmailProp)) {
		obj["serviceAccountEmail"] = serviceAccountEmailProp
	}
	stateProp, err := expandStorageHmacKeyState(d.Get("state"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("state"); !isEmptyValue(reflect.ValueOf(stateProp)) && (ok || !reflect.DeepEqual(v, stateProp)) {
		obj["state"] = stateProp
	}

	url, err := replaceVars(d, config, "{{StorageBasePath}}projects/{{project}}/hmacKeys?serviceAccountEmail={{service_account_email}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new HmacKey: %#v", obj)
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

	res, err := sendRequestWithTimeout(config, "POST", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating HmacKey: %s", err)
	}

	// Store the ID now
	id, err := replaceVars(d, config, "projects/{{project}}/hmacKeys/{{access_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating HmacKey %q: %#v", d.Id(), res)

	// `secret` and `access_id` are generated by the API upon successful CREATE. The following
	// ensures terraform has the correct values based on the Projects.hmacKeys response object.
	secret, ok := res["secret"].(string)
	if !ok {
		return fmt.Errorf("The response to CREATE was missing an expected field. Your create did not work.")
	}

	if err := d.Set("secret", secret); err != nil {
		return fmt.Errorf("Error setting secret: %s", err)
	}

	metadata := res["metadata"].(map[string]interface{})
	accessId, ok := metadata["accessId"].(string)
	if !ok {
		return fmt.Errorf("The response to CREATE was missing an expected field. Your create did not work.")
	}

	if err := d.Set("access_id", accessId); err != nil {
		return fmt.Errorf("Error setting access_id: %s", err)
	}

	id, err = replaceVars(d, config, "projects/{{project}}/hmacKeys/{{access_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}

	d.SetId(id)

	return resourceStorageHmacKeyRead(d, meta)
}

func resourceStorageHmacKeyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{StorageBasePath}}projects/{{project}}/hmacKeys/{{access_id}}")
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

	res, err := sendRequest(config, "GET", billingProject, url, userAgent, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("StorageHmacKey %q", d.Id()))
	}

	res, err = resourceStorageHmacKeyDecoder(d, meta, res)
	if err != nil {
		return err
	}

	if res == nil {
		// Decoding the object has resulted in it being gone. It may be marked deleted
		log.Printf("[DEBUG] Removing StorageHmacKey because it no longer exists.")
		d.SetId("")
		return nil
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading HmacKey: %s", err)
	}

	if err := d.Set("service_account_email", flattenStorageHmacKeyServiceAccountEmail(res["serviceAccountEmail"], d, config)); err != nil {
		return fmt.Errorf("Error reading HmacKey: %s", err)
	}
	if err := d.Set("state", flattenStorageHmacKeyState(res["state"], d, config)); err != nil {
		return fmt.Errorf("Error reading HmacKey: %s", err)
	}
	if err := d.Set("access_id", flattenStorageHmacKeyAccessId(res["accessId"], d, config)); err != nil {
		return fmt.Errorf("Error reading HmacKey: %s", err)
	}
	if err := d.Set("time_created", flattenStorageHmacKeyTimeCreated(res["timeCreated"], d, config)); err != nil {
		return fmt.Errorf("Error reading HmacKey: %s", err)
	}
	if err := d.Set("updated", flattenStorageHmacKeyUpdated(res["updated"], d, config)); err != nil {
		return fmt.Errorf("Error reading HmacKey: %s", err)
	}

	return nil
}

func resourceStorageHmacKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	billingProject = project

	d.Partial(true)

	if d.HasChange("state") {
		obj := make(map[string]interface{})

		getUrl, err := replaceVars(d, config, "{{StorageBasePath}}projects/{{project}}/hmacKeys/{{access_id}}")
		if err != nil {
			return err
		}

		// err == nil indicates that the billing_project value was found
		if bp, err := getBillingProject(d, config); err == nil {
			billingProject = bp
		}

		getRes, err := sendRequest(config, "GET", billingProject, getUrl, userAgent, nil)
		if err != nil {
			return handleNotFoundError(err, d, fmt.Sprintf("StorageHmacKey %q", d.Id()))
		}

		obj["etag"] = getRes["etag"]

		stateProp, err := expandStorageHmacKeyState(d.Get("state"), d, config)
		if err != nil {
			return err
		} else if v, ok := d.GetOkExists("state"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, stateProp)) {
			obj["state"] = stateProp
		}

		url, err := replaceVars(d, config, "{{StorageBasePath}}projects/{{project}}/hmacKeys/{{access_id}}")
		if err != nil {
			return err
		}

		// err == nil indicates that the billing_project value was found
		if bp, err := getBillingProject(d, config); err == nil {
			billingProject = bp
		}

		res, err := sendRequestWithTimeout(config, "PUT", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("Error updating HmacKey %q: %s", d.Id(), err)
		} else {
			log.Printf("[DEBUG] Finished updating HmacKey %q: %#v", d.Id(), res)
		}

	}

	d.Partial(false)

	return resourceStorageHmacKeyRead(d, meta)
}

func resourceStorageHmacKeyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	billingProject = project

	url, err := replaceVars(d, config, "{{StorageBasePath}}projects/{{project}}/hmacKeys/{{access_id}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	getUrl, err := replaceVars(d, config, "{{StorageBasePath}}projects/{{project}}/hmacKeys/{{access_id}}")
	if err != nil {
		return err
	}

	getRes, err := sendRequest(config, "GET", project, getUrl, userAgent, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("StorageHmacKey %q", d.Id()))
	}

	// HmacKeys need to be INACTIVE to be deleted and the API doesn't accept noop
	// updates
	if v := getRes["state"]; v == "ACTIVE" {
		getRes["state"] = "INACTIVE"
		updateUrl, err := replaceVars(d, config, "{{StorageBasePath}}projects/{{project}}/hmacKeys/{{access_id}}")
		if err != nil {
			return err
		}

		log.Printf("[DEBUG] Deactivating HmacKey %q: %#v", d.Id(), getRes)
		_, err = sendRequestWithTimeout(config, "PUT", project, updateUrl, userAgent, getRes, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("Error deactivating HmacKey %q: %s", d.Id(), err)
		}
	}
	log.Printf("[DEBUG] Deleting HmacKey %q", d.Id())

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequestWithTimeout(config, "DELETE", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return handleNotFoundError(err, d, "HmacKey")
	}

	log.Printf("[DEBUG] Finished deleting HmacKey %q: %#v", d.Id(), res)
	return nil
}

func resourceStorageHmacKeyImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/hmacKeys/(?P<access_id>[^/]+)",
		"(?P<project>[^/]+)/(?P<access_id>[^/]+)",
		"(?P<access_id>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "projects/{{project}}/hmacKeys/{{access_id}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenStorageHmacKeyServiceAccountEmail(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenStorageHmacKeyState(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenStorageHmacKeyAccessId(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenStorageHmacKeyTimeCreated(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenStorageHmacKeyUpdated(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func expandStorageHmacKeyServiceAccountEmail(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandStorageHmacKeyState(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func resourceStorageHmacKeyDecoder(d *schema.ResourceData, meta interface{}, res map[string]interface{}) (map[string]interface{}, error) {
	if v := res["state"]; v == "DELETED" {
		return nil, nil
	}

	return res, nil
}
