package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
	"strings"
)

func resourceGoogleOrganizationPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceGoogleOrganizationPolicyCreate,
		Read:   resourceGoogleOrganizationPolicyRead,
		Update: resourceGoogleOrganizationPolicyUpdate,
		Delete: resourceGoogleOrganizationPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceGoogleOrganizationPolicyImportState,
		},

		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"constraint": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: linkDiffSuppress,
			},
			"boolean_policy": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"list_policy"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enforced": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},
			"list_policy": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"boolean_policy"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow": {
							Type:          schema.TypeList,
							Optional:      true,
							MaxItems:      1,
							ConflictsWith: []string{"list_policy.0.deny"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"all": {
										Type:          schema.TypeBool,
										Optional:      true,
										Default:       false,
										ConflictsWith: []string{"list_policy.0.allow.0.values"},
									},
									"values": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Set:      schema.HashString,
									},
								},
							},
						},
						"deny": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"all": {
										Type:          schema.TypeBool,
										Optional:      true,
										Default:       false,
										ConflictsWith: []string{"list_policy.0.deny.0.values"},
									},
									"values": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Set:      schema.HashString,
									},
								},
							},
						},
						"suggested_value": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"version": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"etag": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGoogleOrganizationPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	if err := setOrganizationPolicy(d, meta); err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s:%s", d.Get("org_id"), d.Get("constraint").(string)))

	return resourceGoogleOrganizationPolicyRead(d, meta)
}

func resourceGoogleOrganizationPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	org := "organizations/" + d.Get("org_id").(string)

	policy, err := config.clientResourceManager.Organizations.GetOrgPolicy(org, &cloudresourcemanager.GetOrgPolicyRequest{
		Constraint: canonicalOrgPolicyConstraint(d.Get("constraint").(string)),
	}).Do()

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Organization policy for %s", org))
	}

	d.Set("constraint", policy.Constraint)
	d.Set("boolean_policy", flattenBooleanOrganizationPolicy(policy.BooleanPolicy))
	d.Set("list_policy", flattenListOrganizationPolicy(policy.ListPolicy))
	d.Set("version", policy.Version)
	d.Set("etag", policy.Etag)
	d.Set("update_time", policy.UpdateTime)

	return nil
}

func resourceGoogleOrganizationPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	if err := setOrganizationPolicy(d, meta); err != nil {
		return err
	}

	return resourceGoogleOrganizationPolicyRead(d, meta)
}

func resourceGoogleOrganizationPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	org := "organizations/" + d.Get("org_id").(string)

	_, err := config.clientResourceManager.Organizations.ClearOrgPolicy(org, &cloudresourcemanager.ClearOrgPolicyRequest{
		Constraint: canonicalOrgPolicyConstraint(d.Get("constraint").(string)),
	}).Do()

	if err != nil {
		return err
	}

	return nil
}

func resourceGoogleOrganizationPolicyImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid id format. Expecting {org_id}:{constraint}, got '%s' instead.", d.Id())
	}

	d.Set("org_id", parts[0])
	d.Set("constraint", parts[1])

	return []*schema.ResourceData{d}, nil
}

func setOrganizationPolicy(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	org := "organizations/" + d.Get("org_id").(string)

	listPolicy, err := expandListOrganizationPolicy(d.Get("list_policy").([]interface{}))
	if err != nil {
		return err
	}

	_, err = config.clientResourceManager.Organizations.SetOrgPolicy(org, &cloudresourcemanager.SetOrgPolicyRequest{
		Policy: &cloudresourcemanager.OrgPolicy{
			Constraint:    canonicalOrgPolicyConstraint(d.Get("constraint").(string)),
			BooleanPolicy: expandBooleanOrganizationPolicy(d.Get("boolean_policy").([]interface{})),
			ListPolicy:    listPolicy,
			Version:       int64(d.Get("version").(int)),
			Etag:          d.Get("etag").(string),
		},
	}).Do()

	return err
}

func flattenBooleanOrganizationPolicy(policy *cloudresourcemanager.BooleanPolicy) []map[string]interface{} {
	bPolicies := make([]map[string]interface{}, 0, 1)

	if policy == nil {
		return bPolicies
	}

	bPolicies = append(bPolicies, map[string]interface{}{
		"enforced": policy.Enforced,
	})

	return bPolicies
}

func expandBooleanOrganizationPolicy(configured []interface{}) *cloudresourcemanager.BooleanPolicy {
	if len(configured) == 0 {
		return nil
	}

	booleanPolicy := configured[0].(map[string]interface{})
	return &cloudresourcemanager.BooleanPolicy{
		Enforced: booleanPolicy["enforced"].(bool),
	}
}

func flattenListOrganizationPolicy(policy *cloudresourcemanager.ListPolicy) []map[string]interface{} {
	lPolicies := make([]map[string]interface{}, 0, 1)

	if policy == nil {
		return lPolicies
	}

	listPolicy := map[string]interface{}{}
	switch {
	case policy.AllValues == "ALLOW":
		listPolicy["allow"] = []interface{}{map[string]interface{}{
			"all": true,
		}}
	case policy.AllValues == "DENY":
		listPolicy["deny"] = []interface{}{map[string]interface{}{
			"all": true,
		}}
	case len(policy.AllowedValues) > 0:
		listPolicy["allow"] = []interface{}{map[string]interface{}{
			"values": schema.NewSet(schema.HashString, convertStringArrToInterface(policy.AllowedValues)),
		}}
	case len(policy.DeniedValues) > 0:
		listPolicy["deny"] = []interface{}{map[string]interface{}{
			"values": schema.NewSet(schema.HashString, convertStringArrToInterface(policy.DeniedValues)),
		}}
	}

	lPolicies = append(lPolicies, listPolicy)

	return lPolicies
}

func expandListOrganizationPolicy(configured []interface{}) (*cloudresourcemanager.ListPolicy, error) {
	if len(configured) == 0 {
		return nil, nil
	}

	listPolicyMap := configured[0].(map[string]interface{})

	allow := listPolicyMap["allow"].([]interface{})
	deny := listPolicyMap["deny"].([]interface{})

	var allValues string
	var allowedValues []string
	var deniedValues []string
	if len(allow) > 0 {
		allowMap := allow[0].(map[string]interface{})
		all := allowMap["all"].(bool)
		values := allowMap["values"].(*schema.Set)

		if all {
			allValues = "ALLOW"
		} else {
			allowedValues = convertStringArr(values.List())
		}
	}

	if len(deny) > 0 {
		denyMap := deny[0].(map[string]interface{})
		all := denyMap["all"].(bool)
		values := denyMap["values"].(*schema.Set)

		if all {
			allValues = "DENY"
		} else {
			deniedValues = convertStringArr(values.List())
		}
	}

	listPolicy := configured[0].(map[string]interface{})
	return &cloudresourcemanager.ListPolicy{
		AllValues:      allValues,
		AllowedValues:  allowedValues,
		DeniedValues:   deniedValues,
		SuggestedValue: listPolicy["suggested_value"].(string),
	}, nil
}

func canonicalOrgPolicyConstraint(constraint string) string {
	if strings.HasPrefix(constraint, "constraints/") {
		return constraint
	}
	return "constraints/" + constraint
}
