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

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceMonitoringNotificationChannel() *schema.Resource {
	return &schema.Resource{
		Create: resourceMonitoringNotificationChannelCreate,
		Read:   resourceMonitoringNotificationChannelRead,
		Update: resourceMonitoringNotificationChannelUpdate,
		Delete: resourceMonitoringNotificationChannelDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMonitoringNotificationChannelImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Update: schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(4 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `An optional human-readable name for this notification channel. It is recommended that you specify a non-empty and unique name in order to make it easier to identify the channels in your project, though this is not enforced. The display name is limited to 512 Unicode characters.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the notification channel. This field matches the value of the NotificationChannelDescriptor.type field. See https://cloud.google.com/monitoring/api/ref_v3/rest/v3/projects.notificationChannelDescriptors/list to get the list of valid values such as "email", "slack", etc...`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `An optional human-readable description of this notification channel. This description may provide additional details, beyond the display name, for the channel. This may not exceed 1024 Unicode characters.`,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether notifications are forwarded to the described channel. This makes it possible to disable delivery of notifications to a particular channel without removing the channel from all alerting policies that reference the channel. This is a more convenient approach when the change is temporary and you want to receive notifications from the same set of alerting policies on the channel at some point in the future.`,
				Default:     true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Description: `Configuration fields that define the channel and its behavior. The
permissible and required labels are specified in the
NotificationChannelDescriptor corresponding to the type field.

**Note**: Some NotificationChannelDescriptor labels are
sensitive and the API will return an partially-obfuscated value.
For example, for '"type": "slack"' channels, an 'auth_token'
label with value "SECRET" will be obfuscated as "**CRET". In order
to avoid a diff, Terraform will use the state value if it appears
that the obfuscated value matches the state value in
length/unobfuscated characters. However, Terraform will not detect a
diff if the obfuscated portion of the value was changed outside of
Terraform.`,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"user_labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: `User-supplied key/value data that does not need to conform to the corresponding NotificationChannelDescriptor's schema, unlike the labels field. This field is intended to be used for organizing and identifying the NotificationChannel objects.The field can contain up to 64 entries. Each key and value is limited to 63 Unicode characters or 128 bytes, whichever is smaller. Labels and values can contain only lowercase letters, numerals, underscores, and dashes. Keys must begin with a letter.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The full REST resource name for this channel. The syntax is:
projects/[PROJECT_ID]/notificationChannels/[CHANNEL_ID]
The [CHANNEL_ID] is automatically assigned by the server on creation.`,
			},
			"verification_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates whether this channel has been verified or not. On a ListNotificationChannels or GetNotificationChannel operation, this field is expected to be populated.If the value is UNVERIFIED, then it indicates that the channel is non-functioning (it both requires verification and lacks verification); otherwise, it is assumed that the channel works.If the channel is neither VERIFIED nor UNVERIFIED, it implies that the channel is of a type that does not require verification or that this specific channel has been exempted from verification because it was created prior to verification being required for channels of this type.This field cannot be modified using a standard UpdateNotificationChannel operation. To change the value of this field, you must call VerifyNotificationChannel.`,
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

func resourceMonitoringNotificationChannelCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	obj := make(map[string]interface{})
	labelsProp, err := expandMonitoringNotificationChannelLabels(d.Get("labels"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("labels"); !isEmptyValue(reflect.ValueOf(labelsProp)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}
	typeProp, err := expandMonitoringNotificationChannelType(d.Get("type"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("type"); !isEmptyValue(reflect.ValueOf(typeProp)) && (ok || !reflect.DeepEqual(v, typeProp)) {
		obj["type"] = typeProp
	}
	userLabelsProp, err := expandMonitoringNotificationChannelUserLabels(d.Get("user_labels"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("user_labels"); !isEmptyValue(reflect.ValueOf(userLabelsProp)) && (ok || !reflect.DeepEqual(v, userLabelsProp)) {
		obj["userLabels"] = userLabelsProp
	}
	descriptionProp, err := expandMonitoringNotificationChannelDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !isEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	displayNameProp, err := expandMonitoringNotificationChannelDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("display_name"); !isEmptyValue(reflect.ValueOf(displayNameProp)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}
	enabledProp, err := expandMonitoringNotificationChannelEnabled(d.Get("enabled"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("enabled"); ok || !reflect.DeepEqual(v, enabledProp) {
		obj["enabled"] = enabledProp
	}

	lockName, err := replaceVars(d, config, "stackdriver/notifications/{{project}}")
	if err != nil {
		return err
	}
	mutexKV.Lock(lockName)
	defer mutexKV.Unlock(lockName)

	url, err := replaceVars(d, config, "{{MonitoringBasePath}}projects/{{project}}/notificationChannels")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new NotificationChannel: %#v", obj)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	res, err := sendRequestWithTimeout(config, "POST", project, url, obj, d.Timeout(schema.TimeoutCreate), isMonitoringRetryableError)
	if err != nil {
		return fmt.Errorf("Error creating NotificationChannel: %s", err)
	}

	// Store the ID now
	id, err := replaceVars(d, config, "{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating NotificationChannel %q: %#v", d.Id(), res)

	// `name` is autogenerated from the api so needs to be set post-create
	name, ok := res["name"]
	if !ok {
		return fmt.Errorf("Create response didn't contain critical fields. Create may not have succeeded.")
	}
	d.Set("name", name.(string))
	d.SetId(name.(string))

	return resourceMonitoringNotificationChannelRead(d, meta)
}

func resourceMonitoringNotificationChannelRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	url, err := replaceVars(d, config, "{{MonitoringBasePath}}{{name}}")
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	res, err := sendRequest(config, "GET", project, url, nil, isMonitoringRetryableError)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("MonitoringNotificationChannel %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading NotificationChannel: %s", err)
	}

	if err := d.Set("labels", flattenMonitoringNotificationChannelLabels(res["labels"], d)); err != nil {
		return fmt.Errorf("Error reading NotificationChannel: %s", err)
	}
	if err := d.Set("name", flattenMonitoringNotificationChannelName(res["name"], d)); err != nil {
		return fmt.Errorf("Error reading NotificationChannel: %s", err)
	}
	if err := d.Set("verification_status", flattenMonitoringNotificationChannelVerificationStatus(res["verificationStatus"], d)); err != nil {
		return fmt.Errorf("Error reading NotificationChannel: %s", err)
	}
	if err := d.Set("type", flattenMonitoringNotificationChannelType(res["type"], d)); err != nil {
		return fmt.Errorf("Error reading NotificationChannel: %s", err)
	}
	if err := d.Set("user_labels", flattenMonitoringNotificationChannelUserLabels(res["userLabels"], d)); err != nil {
		return fmt.Errorf("Error reading NotificationChannel: %s", err)
	}
	if err := d.Set("description", flattenMonitoringNotificationChannelDescription(res["description"], d)); err != nil {
		return fmt.Errorf("Error reading NotificationChannel: %s", err)
	}
	if err := d.Set("display_name", flattenMonitoringNotificationChannelDisplayName(res["displayName"], d)); err != nil {
		return fmt.Errorf("Error reading NotificationChannel: %s", err)
	}
	if err := d.Set("enabled", flattenMonitoringNotificationChannelEnabled(res["enabled"], d)); err != nil {
		return fmt.Errorf("Error reading NotificationChannel: %s", err)
	}

	return nil
}

func resourceMonitoringNotificationChannelUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	labelsProp, err := expandMonitoringNotificationChannelLabels(d.Get("labels"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("labels"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}
	typeProp, err := expandMonitoringNotificationChannelType(d.Get("type"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("type"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, typeProp)) {
		obj["type"] = typeProp
	}
	userLabelsProp, err := expandMonitoringNotificationChannelUserLabels(d.Get("user_labels"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("user_labels"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, userLabelsProp)) {
		obj["userLabels"] = userLabelsProp
	}
	descriptionProp, err := expandMonitoringNotificationChannelDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	displayNameProp, err := expandMonitoringNotificationChannelDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("display_name"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}
	enabledProp, err := expandMonitoringNotificationChannelEnabled(d.Get("enabled"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("enabled"); ok || !reflect.DeepEqual(v, enabledProp) {
		obj["enabled"] = enabledProp
	}

	lockName, err := replaceVars(d, config, "stackdriver/notifications/{{project}}")
	if err != nil {
		return err
	}
	mutexKV.Lock(lockName)
	defer mutexKV.Unlock(lockName)

	url, err := replaceVars(d, config, "{{MonitoringBasePath}}{{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating NotificationChannel %q: %#v", d.Id(), obj)
	_, err = sendRequestWithTimeout(config, "PATCH", project, url, obj, d.Timeout(schema.TimeoutUpdate), isMonitoringRetryableError)

	if err != nil {
		return fmt.Errorf("Error updating NotificationChannel %q: %s", d.Id(), err)
	}

	return resourceMonitoringNotificationChannelRead(d, meta)
}

func resourceMonitoringNotificationChannelDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	lockName, err := replaceVars(d, config, "stackdriver/notifications/{{project}}")
	if err != nil {
		return err
	}
	mutexKV.Lock(lockName)
	defer mutexKV.Unlock(lockName)

	url, err := replaceVars(d, config, "{{MonitoringBasePath}}{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	log.Printf("[DEBUG] Deleting NotificationChannel %q", d.Id())

	res, err := sendRequestWithTimeout(config, "DELETE", project, url, obj, d.Timeout(schema.TimeoutDelete), isMonitoringRetryableError)
	if err != nil {
		return handleNotFoundError(err, d, "NotificationChannel")
	}

	log.Printf("[DEBUG] Finished deleting NotificationChannel %q: %#v", d.Id(), res)
	return nil
}

func resourceMonitoringNotificationChannelImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	config := meta.(*Config)

	// current import_formats can't import fields with forward slashes in their value
	if err := parseImportId([]string{"(?P<name>.+)"}, d, config); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

// Some labels are obfuscated for monitoring channels
// e.g. if the value is "SECRET", the server will return "**CRET"
// This method checks to see if the value read from the server looks like
// the obfuscated version of the state value. If so, it will just use the state
// value to avoid permadiff.
func flattenMonitoringNotificationChannelLabels(v interface{}, d *schema.ResourceData) interface{} {
	if v == nil {
		return v
	}
	readLabels := v.(map[string]interface{})

	stateLabelsRaw, ok := d.GetOk("labels")
	if !ok {
		return v
	}
	stateLabels := stateLabelsRaw.(map[string]interface{})

	for k, serverV := range readLabels {
		stateV, ok := stateLabels[k]
		if !ok {
			continue
		}
		useStateV := isMonitoringNotificationChannelLabelsObfuscated(serverV.(string), stateV.(string))
		if useStateV {
			readLabels[k] = stateV.(string)
		}
	}
	return readLabels
}

func isMonitoringNotificationChannelLabelsObfuscated(serverLabel, stateLabel string) bool {
	if stateLabel == serverLabel {
		return false
	}

	if len(stateLabel) != len(serverLabel) {
		return false
	}

	// Check if value read from GCP has either the same character or replaced
	// it with '*'.
	for i := 0; i < len(stateLabel); i++ {
		if serverLabel[i] != '*' && stateLabel[i] != serverLabel[i] {
			return false
		}
	}
	return true
}

func flattenMonitoringNotificationChannelName(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenMonitoringNotificationChannelVerificationStatus(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenMonitoringNotificationChannelType(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenMonitoringNotificationChannelUserLabels(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenMonitoringNotificationChannelDescription(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenMonitoringNotificationChannelDisplayName(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenMonitoringNotificationChannelEnabled(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandMonitoringNotificationChannelLabels(v interface{}, d TerraformResourceData, config *Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func expandMonitoringNotificationChannelType(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandMonitoringNotificationChannelUserLabels(v interface{}, d TerraformResourceData, config *Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func expandMonitoringNotificationChannelDescription(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandMonitoringNotificationChannelDisplayName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandMonitoringNotificationChannelEnabled(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}
