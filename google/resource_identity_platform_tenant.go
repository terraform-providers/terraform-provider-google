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
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIdentityPlatformTenant() *schema.Resource {
	return &schema.Resource{
		Create: resourceIdentityPlatformTenantCreate,
		Read:   resourceIdentityPlatformTenantRead,
		Update: resourceIdentityPlatformTenantUpdate,
		Delete: resourceIdentityPlatformTenantDelete,

		Importer: &schema.ResourceImporter{
			State: resourceIdentityPlatformTenantImport,
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
				Description: `Human friendly display name of the tenant.`,
			},
			"allow_password_signup": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to allow email/password user authentication.`,
			},
			"disable_auth": {
				Type:     schema.TypeBool,
				Optional: true,
				Description: `Whether authentication is disabled for the tenant. If true, the users under
the disabled tenant are not allowed to sign-in. Admins of the disabled tenant
are not able to manage its users.`,
			},
			"enable_email_link_signin": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to enable email link user authentication.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the tenant that is generated by the server`,
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

func resourceIdentityPlatformTenantCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	displayNameProp, err := expandIdentityPlatformTenantDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("display_name"); !isEmptyValue(reflect.ValueOf(displayNameProp)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}
	allowPasswordSignupProp, err := expandIdentityPlatformTenantAllowPasswordSignup(d.Get("allow_password_signup"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("allow_password_signup"); !isEmptyValue(reflect.ValueOf(allowPasswordSignupProp)) && (ok || !reflect.DeepEqual(v, allowPasswordSignupProp)) {
		obj["allowPasswordSignup"] = allowPasswordSignupProp
	}
	enableEmailLinkSigninProp, err := expandIdentityPlatformTenantEnableEmailLinkSignin(d.Get("enable_email_link_signin"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("enable_email_link_signin"); !isEmptyValue(reflect.ValueOf(enableEmailLinkSigninProp)) && (ok || !reflect.DeepEqual(v, enableEmailLinkSigninProp)) {
		obj["enableEmailLinkSignin"] = enableEmailLinkSigninProp
	}
	disableAuthProp, err := expandIdentityPlatformTenantDisableAuth(d.Get("disable_auth"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("disable_auth"); !isEmptyValue(reflect.ValueOf(disableAuthProp)) && (ok || !reflect.DeepEqual(v, disableAuthProp)) {
		obj["disableAuth"] = disableAuthProp
	}

	url, err := replaceVars(d, config, "{{IdentityPlatformBasePath}}projects/{{project}}/tenants")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new Tenant: %#v", obj)
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
		return fmt.Errorf("Error creating Tenant: %s", err)
	}
	if err := d.Set("name", flattenIdentityPlatformTenantName(res["name"], d, config)); err != nil {
		return fmt.Errorf(`Error setting computed identity field "name": %s`, err)
	}

	// Store the ID now
	id, err := replaceVars(d, config, "projects/{{project}}/tenants/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating Tenant %q: %#v", d.Id(), res)

	// `name` is autogenerated from the api so needs to be set post-create
	name, ok := res["name"]
	if !ok {
		return fmt.Errorf("Create response didn't contain critical fields. Create may not have succeeded.")
	}
	if err := d.Set("name", GetResourceNameFromSelfLink(name.(string))); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	// Store the ID now that we have set the computed name
	id, err = replaceVars(d, config, "projects/{{project}}/tenants/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return resourceIdentityPlatformTenantRead(d, meta)
}

func resourceIdentityPlatformTenantRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{IdentityPlatformBasePath}}projects/{{project}}/tenants/{{name}}")
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
		return handleNotFoundError(err, d, fmt.Sprintf("IdentityPlatformTenant %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading Tenant: %s", err)
	}

	if err := d.Set("name", flattenIdentityPlatformTenantName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading Tenant: %s", err)
	}
	if err := d.Set("display_name", flattenIdentityPlatformTenantDisplayName(res["displayName"], d, config)); err != nil {
		return fmt.Errorf("Error reading Tenant: %s", err)
	}
	if err := d.Set("allow_password_signup", flattenIdentityPlatformTenantAllowPasswordSignup(res["allowPasswordSignup"], d, config)); err != nil {
		return fmt.Errorf("Error reading Tenant: %s", err)
	}
	if err := d.Set("enable_email_link_signin", flattenIdentityPlatformTenantEnableEmailLinkSignin(res["enableEmailLinkSignin"], d, config)); err != nil {
		return fmt.Errorf("Error reading Tenant: %s", err)
	}
	if err := d.Set("disable_auth", flattenIdentityPlatformTenantDisableAuth(res["disableAuth"], d, config)); err != nil {
		return fmt.Errorf("Error reading Tenant: %s", err)
	}

	return nil
}

func resourceIdentityPlatformTenantUpdate(d *schema.ResourceData, meta interface{}) error {
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
	displayNameProp, err := expandIdentityPlatformTenantDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("display_name"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}
	allowPasswordSignupProp, err := expandIdentityPlatformTenantAllowPasswordSignup(d.Get("allow_password_signup"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("allow_password_signup"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, allowPasswordSignupProp)) {
		obj["allowPasswordSignup"] = allowPasswordSignupProp
	}
	enableEmailLinkSigninProp, err := expandIdentityPlatformTenantEnableEmailLinkSignin(d.Get("enable_email_link_signin"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("enable_email_link_signin"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, enableEmailLinkSigninProp)) {
		obj["enableEmailLinkSignin"] = enableEmailLinkSigninProp
	}
	disableAuthProp, err := expandIdentityPlatformTenantDisableAuth(d.Get("disable_auth"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("disable_auth"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, disableAuthProp)) {
		obj["disableAuth"] = disableAuthProp
	}

	url, err := replaceVars(d, config, "{{IdentityPlatformBasePath}}projects/{{project}}/tenants/{{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating Tenant %q: %#v", d.Id(), obj)
	updateMask := []string{}

	if d.HasChange("display_name") {
		updateMask = append(updateMask, "displayName")
	}

	if d.HasChange("allow_password_signup") {
		updateMask = append(updateMask, "allowPasswordSignup")
	}

	if d.HasChange("enable_email_link_signin") {
		updateMask = append(updateMask, "enableEmailLinkSignin")
	}

	if d.HasChange("disable_auth") {
		updateMask = append(updateMask, "disableAuth")
	}
	// updateMask is a URL parameter but not present in the schema, so replaceVars
	// won't set it
	url, err = addQueryParams(url, map[string]string{"updateMask": strings.Join(updateMask, ",")})
	if err != nil {
		return err
	}

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequestWithTimeout(config, "PATCH", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutUpdate))

	if err != nil {
		return fmt.Errorf("Error updating Tenant %q: %s", d.Id(), err)
	} else {
		log.Printf("[DEBUG] Finished updating Tenant %q: %#v", d.Id(), res)
	}

	return resourceIdentityPlatformTenantRead(d, meta)
}

func resourceIdentityPlatformTenantDelete(d *schema.ResourceData, meta interface{}) error {
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

	url, err := replaceVars(d, config, "{{IdentityPlatformBasePath}}projects/{{project}}/tenants/{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	log.Printf("[DEBUG] Deleting Tenant %q", d.Id())

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequestWithTimeout(config, "DELETE", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return handleNotFoundError(err, d, "Tenant")
	}

	log.Printf("[DEBUG] Finished deleting Tenant %q: %#v", d.Id(), res)
	return nil
}

func resourceIdentityPlatformTenantImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/tenants/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<name>[^/]+)",
		"(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "projects/{{project}}/tenants/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenIdentityPlatformTenantName(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return v
	}
	return NameFromSelfLinkStateFunc(v)
}

func flattenIdentityPlatformTenantDisplayName(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenIdentityPlatformTenantAllowPasswordSignup(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenIdentityPlatformTenantEnableEmailLinkSignin(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenIdentityPlatformTenantDisableAuth(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func expandIdentityPlatformTenantDisplayName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandIdentityPlatformTenantAllowPasswordSignup(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandIdentityPlatformTenantEnableEmailLinkSignin(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandIdentityPlatformTenantDisableAuth(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}
