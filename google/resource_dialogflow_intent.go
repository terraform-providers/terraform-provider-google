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
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDialogflowIntent() *schema.Resource {
	return &schema.Resource{
		Create: resourceDialogflowIntentCreate,
		Read:   resourceDialogflowIntentRead,
		Update: resourceDialogflowIntentUpdate,
		Delete: resourceDialogflowIntentDelete,

		Importer: &schema.ResourceImporter{
			State: resourceDialogflowIntentImport,
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
				Description: `The name of this intent to be displayed on the console.`,
			},
			"action": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				Description: `The name of the action associated with the intent.
Note: The action name must not contain whitespaces.`,
			},
			"default_response_platforms": {
				Type:     schema.TypeList,
				Optional: true,
				Description: `The list of platforms for which the first responses will be copied from the messages in PLATFORM_UNSPECIFIED
(i.e. default platform). Possible values: ["FACEBOOK", "SLACK", "TELEGRAM", "KIK", "SKYPE", "LINE", "VIBER", "ACTIONS_ON_GOOGLE", "GOOGLE_HANGOUTS"]`,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"FACEBOOK", "SLACK", "TELEGRAM", "KIK", "SKYPE", "LINE", "VIBER", "ACTIONS_ON_GOOGLE", "GOOGLE_HANGOUTS"}, false),
				},
			},
			"events": {
				Type:     schema.TypeList,
				Optional: true,
				Description: `The collection of event names that trigger the intent. If the collection of input contexts is not empty, all of
the contexts must be present in the active user session for an event to trigger this intent. See the 
[events reference](https://cloud.google.com/dialogflow/docs/events-overview) for more details.`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"input_context_names": {
				Type:     schema.TypeList,
				Optional: true,
				Description: `The list of context names required for this intent to be triggered.
Format: projects/<Project ID>/agent/sessions/-/contexts/<Context ID>.`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"is_fallback": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: `Indicates whether this is a fallback intent.`,
			},
			"ml_disabled": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
				Description: `Indicates whether Machine Learning is disabled for the intent.
Note: If mlDisabled setting is set to true, then this intent is not taken into account during inference in ML
ONLY match mode. Also, auto-markup in the UI is turned off.`,
			},
			"parent_followup_intent_name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
				Description: `The unique identifier of the parent intent in the chain of followup intents.
Format: projects/<Project ID>/agent/intents/<Intent ID>.`,
			},
			"priority": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
				Description: `The priority of this intent. Higher numbers represent higher priorities.
  - If the supplied value is unspecified or 0, the service translates the value to 500,000, which corresponds
  to the Normal priority in the console.
  - If the supplied value is negative, the intent is ignored in runtime detect intent requests.`,
			},
			"reset_contexts": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: `Indicates whether to delete all contexts in the current session when this intent is matched.`,
			},
			"webhook_state": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"WEBHOOK_STATE_ENABLED", "WEBHOOK_STATE_ENABLED_FOR_SLOT_FILLING", ""}, false),
				Description: `Indicates whether webhooks are enabled for the intent.
* WEBHOOK_STATE_ENABLED: Webhook is enabled in the agent and in the intent.
* WEBHOOK_STATE_ENABLED_FOR_SLOT_FILLING: Webhook is enabled in the agent and in the intent. Also, each slot
filling prompt is forwarded to the webhook. Possible values: ["WEBHOOK_STATE_ENABLED", "WEBHOOK_STATE_ENABLED_FOR_SLOT_FILLING"]`,
			},
			"followup_intent_info": {
				Type:     schema.TypeList,
				Computed: true,
				Description: `Information about all followup intents that have this intent as a direct or indirect parent. We populate this field
only in the output.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"followup_intent_name": {
							Type:     schema.TypeString,
							Optional: true,
							Description: `The unique identifier of the followup intent.
Format: projects/<Project ID>/agent/intents/<Intent ID>.`,
						},
						"parent_followup_intent_name": {
							Type:     schema.TypeString,
							Optional: true,
							Description: `The unique identifier of the followup intent's parent.
Format: projects/<Project ID>/agent/intents/<Intent ID>.`,
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The unique identifier of this intent. 
Format: projects/<Project ID>/agent/intents/<Intent ID>.`,
			},
			"root_followup_intent_name": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The unique identifier of the root intent in the chain of followup intents. It identifies the correct followup
intents chain for this intent.
Format: projects/<Project ID>/agent/intents/<Intent ID>.`,
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

func resourceDialogflowIntentCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	displayNameProp, err := expandDialogflowIntentDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("display_name"); !isEmptyValue(reflect.ValueOf(displayNameProp)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}
	webhookStateProp, err := expandDialogflowIntentWebhookState(d.Get("webhook_state"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("webhook_state"); !isEmptyValue(reflect.ValueOf(webhookStateProp)) && (ok || !reflect.DeepEqual(v, webhookStateProp)) {
		obj["webhookState"] = webhookStateProp
	}
	priorityProp, err := expandDialogflowIntentPriority(d.Get("priority"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("priority"); !isEmptyValue(reflect.ValueOf(priorityProp)) && (ok || !reflect.DeepEqual(v, priorityProp)) {
		obj["priority"] = priorityProp
	}
	isFallbackProp, err := expandDialogflowIntentIsFallback(d.Get("is_fallback"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("is_fallback"); !isEmptyValue(reflect.ValueOf(isFallbackProp)) && (ok || !reflect.DeepEqual(v, isFallbackProp)) {
		obj["isFallback"] = isFallbackProp
	}
	mlDisabledProp, err := expandDialogflowIntentMlDisabled(d.Get("ml_disabled"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("ml_disabled"); !isEmptyValue(reflect.ValueOf(mlDisabledProp)) && (ok || !reflect.DeepEqual(v, mlDisabledProp)) {
		obj["mlDisabled"] = mlDisabledProp
	}
	inputContextNamesProp, err := expandDialogflowIntentInputContextNames(d.Get("input_context_names"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("input_context_names"); !isEmptyValue(reflect.ValueOf(inputContextNamesProp)) && (ok || !reflect.DeepEqual(v, inputContextNamesProp)) {
		obj["inputContextNames"] = inputContextNamesProp
	}
	eventsProp, err := expandDialogflowIntentEvents(d.Get("events"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("events"); !isEmptyValue(reflect.ValueOf(eventsProp)) && (ok || !reflect.DeepEqual(v, eventsProp)) {
		obj["events"] = eventsProp
	}
	actionProp, err := expandDialogflowIntentAction(d.Get("action"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("action"); !isEmptyValue(reflect.ValueOf(actionProp)) && (ok || !reflect.DeepEqual(v, actionProp)) {
		obj["action"] = actionProp
	}
	resetContextsProp, err := expandDialogflowIntentResetContexts(d.Get("reset_contexts"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("reset_contexts"); !isEmptyValue(reflect.ValueOf(resetContextsProp)) && (ok || !reflect.DeepEqual(v, resetContextsProp)) {
		obj["resetContexts"] = resetContextsProp
	}
	defaultResponsePlatformsProp, err := expandDialogflowIntentDefaultResponsePlatforms(d.Get("default_response_platforms"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("default_response_platforms"); !isEmptyValue(reflect.ValueOf(defaultResponsePlatformsProp)) && (ok || !reflect.DeepEqual(v, defaultResponsePlatformsProp)) {
		obj["defaultResponsePlatforms"] = defaultResponsePlatformsProp
	}
	parentFollowupIntentNameProp, err := expandDialogflowIntentParentFollowupIntentName(d.Get("parent_followup_intent_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("parent_followup_intent_name"); !isEmptyValue(reflect.ValueOf(parentFollowupIntentNameProp)) && (ok || !reflect.DeepEqual(v, parentFollowupIntentNameProp)) {
		obj["parentFollowupIntentName"] = parentFollowupIntentNameProp
	}

	url, err := replaceVars(d, config, "{{DialogflowBasePath}}projects/{{project}}/agent/intents/")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new Intent: %#v", obj)
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
		return fmt.Errorf("Error creating Intent: %s", err)
	}
	if err := d.Set("name", flattenDialogflowIntentName(res["name"], d, config)); err != nil {
		return fmt.Errorf(`Error setting computed identity field "name": %s`, err)
	}

	// Store the ID now
	id, err := replaceVars(d, config, "{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating Intent %q: %#v", d.Id(), res)

	// `name` is autogenerated from the api so needs to be set post-create
	name, ok := res["name"]
	if !ok {
		respBody, ok := res["response"]
		if !ok {
			return fmt.Errorf("Create response didn't contain critical fields. Create may not have succeeded.")
		}

		name, ok = respBody.(map[string]interface{})["name"]
		if !ok {
			return fmt.Errorf("Create response didn't contain critical fields. Create may not have succeeded.")
		}
	}
	if err := d.Set("name", name.(string)); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	d.SetId(name.(string))

	return resourceDialogflowIntentRead(d, meta)
}

func resourceDialogflowIntentRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{DialogflowBasePath}}{{name}}")
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
		return handleNotFoundError(err, d, fmt.Sprintf("DialogflowIntent %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading Intent: %s", err)
	}

	if err := d.Set("name", flattenDialogflowIntentName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading Intent: %s", err)
	}
	if err := d.Set("display_name", flattenDialogflowIntentDisplayName(res["displayName"], d, config)); err != nil {
		return fmt.Errorf("Error reading Intent: %s", err)
	}
	if err := d.Set("webhook_state", flattenDialogflowIntentWebhookState(res["webhookState"], d, config)); err != nil {
		return fmt.Errorf("Error reading Intent: %s", err)
	}
	if err := d.Set("priority", flattenDialogflowIntentPriority(res["priority"], d, config)); err != nil {
		return fmt.Errorf("Error reading Intent: %s", err)
	}
	if err := d.Set("is_fallback", flattenDialogflowIntentIsFallback(res["isFallback"], d, config)); err != nil {
		return fmt.Errorf("Error reading Intent: %s", err)
	}
	if err := d.Set("ml_disabled", flattenDialogflowIntentMlDisabled(res["mlDisabled"], d, config)); err != nil {
		return fmt.Errorf("Error reading Intent: %s", err)
	}
	if err := d.Set("input_context_names", flattenDialogflowIntentInputContextNames(res["inputContextNames"], d, config)); err != nil {
		return fmt.Errorf("Error reading Intent: %s", err)
	}
	if err := d.Set("events", flattenDialogflowIntentEvents(res["events"], d, config)); err != nil {
		return fmt.Errorf("Error reading Intent: %s", err)
	}
	if err := d.Set("action", flattenDialogflowIntentAction(res["action"], d, config)); err != nil {
		return fmt.Errorf("Error reading Intent: %s", err)
	}
	if err := d.Set("reset_contexts", flattenDialogflowIntentResetContexts(res["resetContexts"], d, config)); err != nil {
		return fmt.Errorf("Error reading Intent: %s", err)
	}
	if err := d.Set("default_response_platforms", flattenDialogflowIntentDefaultResponsePlatforms(res["defaultResponsePlatforms"], d, config)); err != nil {
		return fmt.Errorf("Error reading Intent: %s", err)
	}
	if err := d.Set("root_followup_intent_name", flattenDialogflowIntentRootFollowupIntentName(res["rootFollowupIntentName"], d, config)); err != nil {
		return fmt.Errorf("Error reading Intent: %s", err)
	}
	if err := d.Set("parent_followup_intent_name", flattenDialogflowIntentParentFollowupIntentName(res["parentFollowupIntentName"], d, config)); err != nil {
		return fmt.Errorf("Error reading Intent: %s", err)
	}
	if err := d.Set("followup_intent_info", flattenDialogflowIntentFollowupIntentInfo(res["followupIntentInfo"], d, config)); err != nil {
		return fmt.Errorf("Error reading Intent: %s", err)
	}

	return nil
}

func resourceDialogflowIntentUpdate(d *schema.ResourceData, meta interface{}) error {
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

	obj := make(map[string]interface{})
	displayNameProp, err := expandDialogflowIntentDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("display_name"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}
	webhookStateProp, err := expandDialogflowIntentWebhookState(d.Get("webhook_state"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("webhook_state"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, webhookStateProp)) {
		obj["webhookState"] = webhookStateProp
	}
	priorityProp, err := expandDialogflowIntentPriority(d.Get("priority"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("priority"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, priorityProp)) {
		obj["priority"] = priorityProp
	}
	isFallbackProp, err := expandDialogflowIntentIsFallback(d.Get("is_fallback"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("is_fallback"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, isFallbackProp)) {
		obj["isFallback"] = isFallbackProp
	}
	mlDisabledProp, err := expandDialogflowIntentMlDisabled(d.Get("ml_disabled"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("ml_disabled"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, mlDisabledProp)) {
		obj["mlDisabled"] = mlDisabledProp
	}
	inputContextNamesProp, err := expandDialogflowIntentInputContextNames(d.Get("input_context_names"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("input_context_names"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, inputContextNamesProp)) {
		obj["inputContextNames"] = inputContextNamesProp
	}
	eventsProp, err := expandDialogflowIntentEvents(d.Get("events"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("events"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, eventsProp)) {
		obj["events"] = eventsProp
	}
	actionProp, err := expandDialogflowIntentAction(d.Get("action"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("action"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, actionProp)) {
		obj["action"] = actionProp
	}
	resetContextsProp, err := expandDialogflowIntentResetContexts(d.Get("reset_contexts"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("reset_contexts"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, resetContextsProp)) {
		obj["resetContexts"] = resetContextsProp
	}
	defaultResponsePlatformsProp, err := expandDialogflowIntentDefaultResponsePlatforms(d.Get("default_response_platforms"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("default_response_platforms"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, defaultResponsePlatformsProp)) {
		obj["defaultResponsePlatforms"] = defaultResponsePlatformsProp
	}

	url, err := replaceVars(d, config, "{{DialogflowBasePath}}{{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating Intent %q: %#v", d.Id(), obj)

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequestWithTimeout(config, "PATCH", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutUpdate))

	if err != nil {
		return fmt.Errorf("Error updating Intent %q: %s", d.Id(), err)
	} else {
		log.Printf("[DEBUG] Finished updating Intent %q: %#v", d.Id(), res)
	}

	return resourceDialogflowIntentRead(d, meta)
}

func resourceDialogflowIntentDelete(d *schema.ResourceData, meta interface{}) error {
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

	url, err := replaceVars(d, config, "{{DialogflowBasePath}}{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	log.Printf("[DEBUG] Deleting Intent %q", d.Id())

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequestWithTimeout(config, "DELETE", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return handleNotFoundError(err, d, "Intent")
	}

	log.Printf("[DEBUG] Finished deleting Intent %q: %#v", d.Id(), res)
	return nil
}

func resourceDialogflowIntentImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	config := meta.(*Config)

	// current import_formats can't import fields with forward slashes in their value
	if err := parseImportId([]string{"(?P<name>.+)"}, d, config); err != nil {
		return nil, err
	}

	stringParts := strings.Split(d.Get("name").(string), "/")
	if len(stringParts) < 2 {
		return nil, fmt.Errorf(
			"Could not split project from name: %s",
			d.Get("name"),
		)
	}

	if err := d.Set("project", stringParts[1]); err != nil {
		return nil, fmt.Errorf("Error setting project: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}

func flattenDialogflowIntentName(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenDialogflowIntentDisplayName(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenDialogflowIntentWebhookState(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenDialogflowIntentPriority(v interface{}, d *schema.ResourceData, config *Config) interface{} {
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

func flattenDialogflowIntentIsFallback(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenDialogflowIntentMlDisabled(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenDialogflowIntentInputContextNames(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenDialogflowIntentEvents(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenDialogflowIntentAction(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenDialogflowIntentResetContexts(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenDialogflowIntentDefaultResponsePlatforms(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenDialogflowIntentRootFollowupIntentName(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenDialogflowIntentParentFollowupIntentName(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenDialogflowIntentFollowupIntentInfo(v interface{}, d *schema.ResourceData, config *Config) interface{} {
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
			"followup_intent_name":        flattenDialogflowIntentFollowupIntentInfoFollowupIntentName(original["followupIntentName"], d, config),
			"parent_followup_intent_name": flattenDialogflowIntentFollowupIntentInfoParentFollowupIntentName(original["parentFollowupIntentName"], d, config),
		})
	}
	return transformed
}
func flattenDialogflowIntentFollowupIntentInfoFollowupIntentName(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenDialogflowIntentFollowupIntentInfoParentFollowupIntentName(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func expandDialogflowIntentDisplayName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDialogflowIntentWebhookState(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDialogflowIntentPriority(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDialogflowIntentIsFallback(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDialogflowIntentMlDisabled(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDialogflowIntentInputContextNames(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDialogflowIntentEvents(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDialogflowIntentAction(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDialogflowIntentResetContexts(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDialogflowIntentDefaultResponsePlatforms(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDialogflowIntentParentFollowupIntentName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}
