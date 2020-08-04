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
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func expandCloudIotDeviceRegistryHTTPConfig(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	original := v.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedHTTPEnabledState, err := expandCloudIotDeviceRegistryHTTPEnabledState(original["http_enabled_state"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedHTTPEnabledState); val.IsValid() && !isEmptyValue(val) {
		transformed["httpEnabledState"] = transformedHTTPEnabledState
	}

	return transformed, nil
}

func expandCloudIotDeviceRegistryHTTPEnabledState(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandCloudIotDeviceRegistryMqttConfig(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	original := v.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedMqttEnabledState, err := expandCloudIotDeviceRegistryMqttEnabledState(original["mqtt_enabled_state"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMqttEnabledState); val.IsValid() && !isEmptyValue(val) {
		transformed["mqttEnabledState"] = transformedMqttEnabledState
	}

	return transformed, nil
}

func expandCloudIotDeviceRegistryMqttEnabledState(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandCloudIotDeviceRegistryStateNotificationConfig(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	original := v.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedPubsubTopicName, err := expandCloudIotDeviceRegistryStateNotificationConfigPubsubTopicName(original["pubsub_topic_name"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedPubsubTopicName); val.IsValid() && !isEmptyValue(val) {
		transformed["pubsubTopicName"] = transformedPubsubTopicName
	}

	return transformed, nil
}

func expandCloudIotDeviceRegistryStateNotificationConfigPubsubTopicName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandCloudIotDeviceRegistryCredentials(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))

	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedPublicKeyCertificate, err := expandCloudIotDeviceRegistryCredentialsPublicKeyCertificate(original["public_key_certificate"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedPublicKeyCertificate); val.IsValid() && !isEmptyValue(val) {
			transformed["publicKeyCertificate"] = transformedPublicKeyCertificate
		}

		req = append(req, transformed)
	}

	return req, nil
}

func expandCloudIotDeviceRegistryCredentialsPublicKeyCertificate(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	original := v.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedFormat, err := expandCloudIotDeviceRegistryPublicKeyCertificateFormat(original["format"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedFormat); val.IsValid() && !isEmptyValue(val) {
		transformed["format"] = transformedFormat
	}

	transformedCertificate, err := expandCloudIotDeviceRegistryPublicKeyCertificateCertificate(original["certificate"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedCertificate); val.IsValid() && !isEmptyValue(val) {
		transformed["certificate"] = transformedCertificate
	}

	return transformed, nil
}

func expandCloudIotDeviceRegistryPublicKeyCertificateFormat(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandCloudIotDeviceRegistryPublicKeyCertificateCertificate(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func flattenCloudIotDeviceRegistryCredentials(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	log.Printf("[DEBUG] Flattening device resitry credentials: %q", d.Id())
	if v == nil {
		log.Printf("[DEBUG] The credentials array is nil: %q", d.Id())
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	for _, raw := range l {
		original := raw.(map[string]interface{})
		log.Printf("[DEBUG] Original credential: %+v", original)
		if len(original) < 1 {
			log.Printf("[DEBUG] Excluding empty credential that the API returned. %q", d.Id())
			continue
		}
		log.Printf("[DEBUG] Credentials array before appending a new credential: %+v", transformed)
		transformed = append(transformed, map[string]interface{}{
			"public_key_certificate": flattenCloudIotDeviceRegistryCredentialsPublicKeyCertificate(original["publicKeyCertificate"], d, config),
		})
		log.Printf("[DEBUG] Credentials array after appending a new credential: %+v", transformed)
	}
	return transformed
}

func flattenCloudIotDeviceRegistryCredentialsPublicKeyCertificate(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	log.Printf("[DEBUG] Flattening device resitry credentials public key certificate: %q", d.Id())
	if v == nil {
		log.Printf("[DEBUG] The public key certificate is nil: %q", d.Id())
		return v
	}

	original := v.(map[string]interface{})
	log.Printf("[DEBUG] Original public key certificate: %+v", original)
	transformed := make(map[string]interface{})

	transformedPublicKeyCertificateFormat := flattenCloudIotDeviceRegistryPublicKeyCertificateFormat(original["format"], d, config)
	transformed["format"] = transformedPublicKeyCertificateFormat

	transformedPublicKeyCertificateCertificate := flattenCloudIotDeviceRegistryPublicKeyCertificateCertificate(original["certificate"], d, config)
	transformed["certificate"] = transformedPublicKeyCertificateCertificate

	log.Printf("[DEBUG] Transformed public key certificate: %+v", transformed)

	return transformed
}

func flattenCloudIotDeviceRegistryPublicKeyCertificateFormat(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenCloudIotDeviceRegistryPublicKeyCertificateCertificate(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenCloudIotDeviceRegistryHTTPConfig(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return v
	}

	original := v.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedHTTPEnabledState := flattenCloudIotDeviceRegistryHTTPConfigHTTPEnabledState(original["httpEnabledState"], d, config)
	transformed["http_enabled_state"] = transformedHTTPEnabledState

	return transformed
}

func flattenCloudIotDeviceRegistryHTTPConfigHTTPEnabledState(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenCloudIotDeviceRegistryMqttConfig(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return v
	}

	original := v.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedMqttEnabledState := flattenCloudIotDeviceRegistryMqttConfigMqttEnabledState(original["mqttEnabledState"], d, config)
	transformed["mqtt_enabled_state"] = transformedMqttEnabledState

	return transformed
}

func flattenCloudIotDeviceRegistryMqttConfigMqttEnabledState(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenCloudIotDeviceRegistryStateNotificationConfig(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	log.Printf("[DEBUG] Flattening state notification config: %+v", v)
	if v == nil {
		return v
	}

	original := v.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedPubsubTopicName := flattenCloudIotDeviceRegistryStateNotificationConfigPubsubTopicName(original["pubsubTopicName"], d, config)
	if val := reflect.ValueOf(transformedPubsubTopicName); val.IsValid() && !isEmptyValue(val) {
		log.Printf("[DEBUG] pubsub topic name is not null: %v", d.Get("pubsub_topic_name"))
		transformed["pubsub_topic_name"] = transformedPubsubTopicName
	}

	return transformed
}

func flattenCloudIotDeviceRegistryStateNotificationConfigPubsubTopicName(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func validateCloudIotDeviceRegistryID(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if strings.HasPrefix(value, "goog") {
		errors = append(errors, fmt.Errorf(
			"%q (%q) can not start with \"goog\"", k, value))
	}
	if !regexp.MustCompile(CloudIoTIdRegex).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q (%q) doesn't match regexp %q", k, value, CloudIoTIdRegex))
	}
	return
}

func validateCloudIotDeviceRegistrySubfolderMatch(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if strings.HasPrefix(value, "/") {
		errors = append(errors, fmt.Errorf(
			"%q (%q) can not start with '/'", k, value))
	}
	return
}

func resourceCloudIotDeviceRegistry() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudIotDeviceRegistryCreate,
		Read:   resourceCloudIotDeviceRegistryRead,
		Update: resourceCloudIotDeviceRegistryUpdate,
		Delete: resourceCloudIotDeviceRegistryDelete,

		Importer: &schema.ResourceImporter{
			State: resourceCloudIotDeviceRegistryImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Update: schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(4 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateCloudIotDeviceRegistryID,
				Description:  `A unique name for the resource, required by device registry.`,
			},
			"event_notification_configs": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Description: `List of configurations for event notifications, such as PubSub topics
to publish device events to.`,
				MaxItems: 10,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pubsub_topic_name": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: compareSelfLinkOrResourceName,
							Description:      `PubSub topic name to publish device events.`,
						},
						"subfolder_matches": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateCloudIotDeviceRegistrySubfolderMatch,
							Description: `If the subfolder name matches this string exactly, this
configuration will be used. The string must not include the
leading '/' character. If empty, all strings are matched. Empty
value can only be used for the last 'event_notification_configs'
item.`,
						},
					},
				},
			},
			"log_level": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringInSlice([]string{"NONE", "ERROR", "INFO", "DEBUG", ""}, false),
				DiffSuppressFunc: emptyOrDefaultStringSuppress("NONE"),
				Description: `The default logging verbosity for activity from devices in this
registry. Specifies which events should be written to logs. For
example, if the LogLevel is ERROR, only events that terminate in
errors will be logged. LogLevel is inclusive; enabling INFO logging
will also enable ERROR logging. Default value: "NONE" Possible values: ["NONE", "ERROR", "INFO", "DEBUG"]`,
				Default: "NONE",
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
				Description: `The region in which the created registry should reside.
If it is not provided, the provider region is used.`,
			},
			"state_notification_config": {
				Type:        schema.TypeMap,
				Description: `A PubSub topic to publish device state updates.`,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pubsub_topic_name": {
							Type:             schema.TypeString,
							Description:      `PubSub topic name to publish device state updates.`,
							Required:         true,
							DiffSuppressFunc: compareSelfLinkOrResourceName,
						},
					},
				},
			},
			"mqtt_config": {
				Type:        schema.TypeMap,
				Description: `Activate or deactivate MQTT.`,
				Computed:    true,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mqtt_enabled_state": {
							Type:        schema.TypeString,
							Description: `The field allows MQTT_ENABLED or MQTT_DISABLED`,
							Required:    true,
							ValidateFunc: validation.StringInSlice(
								[]string{"MQTT_DISABLED", "MQTT_ENABLED"}, false),
						},
					},
				},
			},
			"http_config": {
				Type:        schema.TypeMap,
				Description: `Activate or deactivate HTTP.`,
				Computed:    true,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"http_enabled_state": {
							Type:        schema.TypeString,
							Description: `The field allows HTTP_ENABLED or HTTP_DISABLED`,
							Required:    true,
							ValidateFunc: validation.StringInSlice(
								[]string{"HTTP_DISABLED", "HTTP_ENABLED"}, false),
						},
					},
				},
			},
			"credentials": {
				Type:        schema.TypeList,
				Description: `List of public key certificates to authenticate devices.`,
				Optional:    true,
				MaxItems:    10,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"public_key_certificate": {
							Type:        schema.TypeMap,
							Description: `A public key certificate format and data.`,
							Required:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"format": {
										Type:        schema.TypeString,
										Description: `The field allows only X509_CERTIFICATE_PEM.`,
										Required:    true,
										ValidateFunc: validation.StringInSlice(
											[]string{"X509_CERTIFICATE_PEM"}, false),
									},
									"certificate": {
										Type:        schema.TypeString,
										Description: `The certificate data.`,
										Required:    true,
									},
								},
							},
						},
					},
				},
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

func resourceCloudIotDeviceRegistryCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	obj := make(map[string]interface{})
	idProp, err := expandCloudIotDeviceRegistryName(d.Get("name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("name"); !isEmptyValue(reflect.ValueOf(idProp)) && (ok || !reflect.DeepEqual(v, idProp)) {
		obj["id"] = idProp
	}
	eventNotificationConfigsProp, err := expandCloudIotDeviceRegistryEventNotificationConfigs(d.Get("event_notification_configs"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("event_notification_configs"); !isEmptyValue(reflect.ValueOf(eventNotificationConfigsProp)) && (ok || !reflect.DeepEqual(v, eventNotificationConfigsProp)) {
		obj["eventNotificationConfigs"] = eventNotificationConfigsProp
	}
	logLevelProp, err := expandCloudIotDeviceRegistryLogLevel(d.Get("log_level"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("log_level"); !isEmptyValue(reflect.ValueOf(logLevelProp)) && (ok || !reflect.DeepEqual(v, logLevelProp)) {
		obj["logLevel"] = logLevelProp
	}

	obj, err = resourceCloudIotDeviceRegistryEncoder(d, meta, obj)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{CloudIotBasePath}}projects/{{project}}/locations/{{region}}/registries")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new DeviceRegistry: %#v", obj)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	res, err := sendRequestWithTimeout(config, "POST", project, url, obj, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating DeviceRegistry: %s", err)
	}

	// Store the ID now
	id, err := replaceVars(d, config, "projects/{{project}}/locations/{{region}}/registries/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating DeviceRegistry %q: %#v", d.Id(), res)

	return resourceCloudIotDeviceRegistryRead(d, meta)
}

func resourceCloudIotDeviceRegistryRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	url, err := replaceVars(d, config, "{{CloudIotBasePath}}projects/{{project}}/locations/{{region}}/registries/{{name}}")
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	res, err := sendRequest(config, "GET", project, url, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("CloudIotDeviceRegistry %q", d.Id()))
	}

	res, err = resourceCloudIotDeviceRegistryDecoder(d, meta, res)
	if err != nil {
		return err
	}

	if res == nil {
		// Decoding the object has resulted in it being gone. It may be marked deleted
		log.Printf("[DEBUG] Removing CloudIotDeviceRegistry because it no longer exists.")
		d.SetId("")
		return nil
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading DeviceRegistry: %s", err)
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}
	if err := d.Set("region", region); err != nil {
		return fmt.Errorf("Error reading DeviceRegistry: %s", err)
	}

	if err := d.Set("name", flattenCloudIotDeviceRegistryName(res["id"], d, config)); err != nil {
		return fmt.Errorf("Error reading DeviceRegistry: %s", err)
	}
	if err := d.Set("event_notification_configs", flattenCloudIotDeviceRegistryEventNotificationConfigs(res["eventNotificationConfigs"], d, config)); err != nil {
		return fmt.Errorf("Error reading DeviceRegistry: %s", err)
	}
	if err := d.Set("log_level", flattenCloudIotDeviceRegistryLogLevel(res["logLevel"], d, config)); err != nil {
		return fmt.Errorf("Error reading DeviceRegistry: %s", err)
	}

	return nil
}

func resourceCloudIotDeviceRegistryUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	eventNotificationConfigsProp, err := expandCloudIotDeviceRegistryEventNotificationConfigs(d.Get("event_notification_configs"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("event_notification_configs"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, eventNotificationConfigsProp)) {
		obj["eventNotificationConfigs"] = eventNotificationConfigsProp
	}
	logLevelProp, err := expandCloudIotDeviceRegistryLogLevel(d.Get("log_level"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("log_level"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, logLevelProp)) {
		obj["logLevel"] = logLevelProp
	}

	obj, err = resourceCloudIotDeviceRegistryEncoder(d, meta, obj)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{CloudIotBasePath}}projects/{{project}}/locations/{{region}}/registries/{{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating DeviceRegistry %q: %#v", d.Id(), obj)
	updateMask := []string{}

	if d.HasChange("event_notification_configs") {
		updateMask = append(updateMask, "eventNotificationConfigs")
	}

	if d.HasChange("log_level") {
		updateMask = append(updateMask, "logLevel")
	}
	// updateMask is a URL parameter but not present in the schema, so replaceVars
	// won't set it
	url, err = addQueryParams(url, map[string]string{"updateMask": strings.Join(updateMask, ",")})
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] updateMask before adding extra schema entries %q: %v", d.Id(), updateMask)

	log.Printf("[DEBUG] Pre-update on state notification config: %q", d.Id())
	if d.HasChange("state_notification_config") {
		log.Printf("[DEBUG] %q stateNotificationConfig.pubsubTopicName has a change. Adding it to the update mask", d.Id())
		updateMask = append(updateMask, "stateNotificationConfig.pubsubTopicName")
	}

	log.Printf("[DEBUG] Pre-update on MQTT config: %q", d.Id())
	if d.HasChange("mqtt_config") {
		log.Printf("[DEBUG] %q mqttConfig.mqttEnabledState has a change. Adding it to the update mask", d.Id())
		updateMask = append(updateMask, "mqttConfig.mqttEnabledState")
	}

	log.Printf("[DEBUG] Pre-update on HTTP config: %q", d.Id())
	if d.HasChange("http_config") {
		log.Printf("[DEBUG] %q httpConfig.httpEnabledState has a change. Adding it to the update mask", d.Id())
		updateMask = append(updateMask, "httpConfig.httpEnabledState")
	}

	log.Printf("[DEBUG] Pre-update on credentials: %q", d.Id())
	if d.HasChange("credentials") {
		log.Printf("[DEBUG] %q credentials has a change. Adding it to the update mask", d.Id())
		updateMask = append(updateMask, "credentials")
	}

	log.Printf("[DEBUG] updateMask after adding extra schema entries %q: %v", d.Id(), updateMask)

	// Refreshing updateMask after adding extra schema entries
	url, err = addQueryParams(url, map[string]string{"updateMask": strings.Join(updateMask, ",")})
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Update URL %q: %v", d.Id(), url)
	res, err := sendRequestWithTimeout(config, "PATCH", project, url, obj, d.Timeout(schema.TimeoutUpdate))

	if err != nil {
		return fmt.Errorf("Error updating DeviceRegistry %q: %s", d.Id(), err)
	} else {
		log.Printf("[DEBUG] Finished updating DeviceRegistry %q: %#v", d.Id(), res)
	}

	return resourceCloudIotDeviceRegistryRead(d, meta)
}

func resourceCloudIotDeviceRegistryDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{CloudIotBasePath}}projects/{{project}}/locations/{{region}}/registries/{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	log.Printf("[DEBUG] Deleting DeviceRegistry %q", d.Id())

	res, err := sendRequestWithTimeout(config, "DELETE", project, url, obj, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return handleNotFoundError(err, d, "DeviceRegistry")
	}

	log.Printf("[DEBUG] Finished deleting DeviceRegistry %q: %#v", d.Id(), res)
	return nil
}

func resourceCloudIotDeviceRegistryImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"(?P<project>[^/]+)/locations/(?P<region>[^/]+)/registries/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<region>[^/]+)/(?P<name>[^/]+)",
		"(?P<region>[^/]+)/(?P<name>[^/]+)",
		"(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "projects/{{project}}/locations/{{region}}/registries/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenCloudIotDeviceRegistryName(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenCloudIotDeviceRegistryEventNotificationConfigs(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed = append(transformed, map[string]interface{}{
			"subfolder_matches": flattenCloudIotDeviceRegistryEventNotificationConfigsSubfolderMatches(original["subfolderMatches"], d, config),
			"pubsub_topic_name": flattenCloudIotDeviceRegistryEventNotificationConfigsPubsubTopicName(original["pubsubTopicName"], d, config),
		})
	}
	return transformed
}
func flattenCloudIotDeviceRegistryEventNotificationConfigsSubfolderMatches(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenCloudIotDeviceRegistryEventNotificationConfigsPubsubTopicName(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenCloudIotDeviceRegistryLogLevel(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func expandCloudIotDeviceRegistryName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandCloudIotDeviceRegistryEventNotificationConfigs(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedSubfolderMatches, err := expandCloudIotDeviceRegistryEventNotificationConfigsSubfolderMatches(original["subfolder_matches"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedSubfolderMatches); val.IsValid() && !isEmptyValue(val) {
			transformed["subfolderMatches"] = transformedSubfolderMatches
		}

		transformedPubsubTopicName, err := expandCloudIotDeviceRegistryEventNotificationConfigsPubsubTopicName(original["pubsub_topic_name"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedPubsubTopicName); val.IsValid() && !isEmptyValue(val) {
			transformed["pubsubTopicName"] = transformedPubsubTopicName
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandCloudIotDeviceRegistryEventNotificationConfigsSubfolderMatches(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandCloudIotDeviceRegistryEventNotificationConfigsPubsubTopicName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandCloudIotDeviceRegistryLogLevel(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func resourceCloudIotDeviceRegistryEncoder(d *schema.ResourceData, meta interface{}, obj map[string]interface{}) (map[string]interface{}, error) {
	config := meta.(*Config)

	log.Printf("[DEBUG] Resource data before encoding extra schema entries %q: %#v", d.Id(), obj)

	log.Printf("[DEBUG] Encoding state notification config: %q", d.Id())
	stateNotificationConfigProp, err := expandCloudIotDeviceRegistryStateNotificationConfig(d.Get("state_notification_config"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("state_notification_config"); !isEmptyValue(reflect.ValueOf(stateNotificationConfigProp)) && (ok || !reflect.DeepEqual(v, stateNotificationConfigProp)) {
		log.Printf("[DEBUG] Encoding %q. Setting stateNotificationConfig: %#v", d.Id(), stateNotificationConfigProp)
		obj["stateNotificationConfig"] = stateNotificationConfigProp
	}

	log.Printf("[DEBUG] Encoding HTTP config: %q", d.Id())
	httpConfigProp, err := expandCloudIotDeviceRegistryHTTPConfig(d.Get("http_config"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("http_config"); !isEmptyValue(reflect.ValueOf(httpConfigProp)) && (ok || !reflect.DeepEqual(v, httpConfigProp)) {
		log.Printf("[DEBUG] Encoding %q. Setting httpConfig: %#v", d.Id(), httpConfigProp)
		obj["httpConfig"] = httpConfigProp
	}

	log.Printf("[DEBUG] Encoding MQTT config: %q", d.Id())
	mqttConfigProp, err := expandCloudIotDeviceRegistryMqttConfig(d.Get("mqtt_config"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("mqtt_config"); !isEmptyValue(reflect.ValueOf(mqttConfigProp)) && (ok || !reflect.DeepEqual(v, mqttConfigProp)) {
		log.Printf("[DEBUG] Encoding %q. Setting mqttConfig: %#v", d.Id(), mqttConfigProp)
		obj["mqttConfig"] = mqttConfigProp
	}

	log.Printf("[DEBUG] Encoding credentials: %q", d.Id())
	credentialsProp, err := expandCloudIotDeviceRegistryCredentials(d.Get("credentials"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("credentials"); !isEmptyValue(reflect.ValueOf(credentialsProp)) && (ok || !reflect.DeepEqual(v, credentialsProp)) {
		log.Printf("[DEBUG] Encoding %q. Setting credentials: %#v", d.Id(), credentialsProp)
		obj["credentials"] = credentialsProp
	}

	log.Printf("[DEBUG] Resource data after encoding extra schema entries %q: %#v", d.Id(), obj)

	return obj, nil
}

func resourceCloudIotDeviceRegistryDecoder(d *schema.ResourceData, meta interface{}, res map[string]interface{}) (map[string]interface{}, error) {
	config := meta.(*Config)

	log.Printf("[DEBUG] Decoding state notification config: %q", d.Id())
	log.Printf("[DEBUG] State notification config before decoding: %v", d.Get("state_notification_config"))
	if err := d.Set("state_notification_config", flattenCloudIotDeviceRegistryStateNotificationConfig(res["stateNotificationConfig"], d, config)); err != nil {
		return nil, fmt.Errorf("Error reading DeviceRegistry: %s", err)
	}
	log.Printf("[DEBUG] State notification config after decoding: %v", d.Get("state_notification_config"))

	log.Printf("[DEBUG] Decoding HTTP config: %q", d.Id())
	log.Printf("[DEBUG] HTTP config before decoding: %v", d.Get("http_config"))
	if err := d.Set("http_config", flattenCloudIotDeviceRegistryHTTPConfig(res["httpConfig"], d, config)); err != nil {
		return nil, fmt.Errorf("Error reading DeviceRegistry: %s", err)
	}
	log.Printf("[DEBUG] HTTP config after decoding: %v", d.Get("http_config"))

	log.Printf("[DEBUG] Decoding MQTT config: %q", d.Id())
	log.Printf("[DEBUG] MQTT config before decoding: %v", d.Get("mqtt_config"))
	if err := d.Set("mqtt_config", flattenCloudIotDeviceRegistryMqttConfig(res["mqttConfig"], d, config)); err != nil {
		return nil, fmt.Errorf("Error reading DeviceRegistry: %s", err)
	}
	log.Printf("[DEBUG] MQTT config after decoding: %v", d.Get("mqtt_config"))

	log.Printf("[DEBUG] Decoding credentials: %q", d.Id())
	log.Printf("[DEBUG] credentials before decoding: %v", d.Get("credentials"))
	if err := d.Set("credentials", flattenCloudIotDeviceRegistryCredentials(res["credentials"], d, config)); err != nil {
		return nil, fmt.Errorf("Error reading DeviceRegistry: %s", err)
	}
	log.Printf("[DEBUG] credentials after decoding: %v", d.Get("credentials"))

	return res, nil
}
