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

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"google.golang.org/api/googleapi"
)

func resourceNameSetFromSelfLinkSet(v interface{}) *schema.Set {
	if v == nil {
		return schema.NewSet(schema.HashString, nil)
	}
	vSet := v.(*schema.Set)
	ls := make([]interface{}, 0, vSet.Len())
	for _, v := range vSet.List() {
		if v == nil {
			continue
		}
		ls = append(ls, GetResourceNameFromSelfLink(v.(string)))
	}
	return schema.NewSet(schema.HashString, ls)
}

// drain_nat_ips MUST be set from (just set) previous values of nat_ips
// so this customizeDiff func makes sure drainNatIps values:
//   - aren't set at creation time
//   - are in old value of nat_ips but not in new values
func resourceComputeRouterNatDrainNatIpsCustomDiff(diff *schema.ResourceDiff, meta interface{}) error {
	o, n := diff.GetChange("drain_nat_ips")
	oSet := resourceNameSetFromSelfLinkSet(o)
	nSet := resourceNameSetFromSelfLinkSet(n)
	addDrainIps := nSet.Difference(oSet)

	// We don't care if there are no new drainNatIps
	if addDrainIps.Len() == 0 {
		return nil
	}

	// Resource hasn't been created yet - return error
	if diff.Id() == "" {
		return fmt.Errorf("New RouterNat cannot have drain_nat_ips, got values %+v", addDrainIps.List())
	}
	//
	o, n = diff.GetChange("nat_ips")
	oNatSet := resourceNameSetFromSelfLinkSet(o)
	nNatSet := resourceNameSetFromSelfLinkSet(n)

	// Resource is being updated - make sure new drainNatIps were in natIps prior d and no longer are in natIps.
	for _, v := range addDrainIps.List() {
		if !oNatSet.Contains(v) {
			return fmt.Errorf("drain_nat_ip %q was not previously set in nat_ips %+v", v.(string), oNatSet.List())
		}
		if nNatSet.Contains(v) {
			return fmt.Errorf("drain_nat_ip %q cannot be drained if still set in nat_ips %+v", v.(string), nNatSet.List())
		}
	}
	return nil
}

func computeRouterNatSubnetworkHash(v interface{}) int {
	obj := v.(map[string]interface{})
	name := obj["name"]
	sourceIpRanges := obj["source_ip_ranges_to_nat"]
	sourceIpRangesHash := 0
	if sourceIpRanges != nil {
		sourceIpSet := sourceIpRanges.(*schema.Set)

		for _, ipRange := range sourceIpSet.List() {
			sourceIpRangesHash += schema.HashString(ipRange.(string))
		}
	}

	secondaryIpRangeNames := obj["secondary_ip_range_names"]
	secondaryIpRangeHash := 0
	if secondaryIpRangeNames != nil {
		secondaryIpRangeSet := secondaryIpRangeNames.(*schema.Set)

		for _, secondaryIp := range secondaryIpRangeSet.List() {
			secondaryIpRangeHash += schema.HashString(secondaryIp.(string))
		}
	}

	return schema.HashString(NameFromSelfLinkStateFunc(name)) + sourceIpRangesHash + secondaryIpRangeHash
}

func resourceComputeRouterNat() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeRouterNatCreate,
		Read:   resourceComputeRouterNatRead,
		Update: resourceComputeRouterNatUpdate,
		Delete: resourceComputeRouterNatDelete,

		Importer: &schema.ResourceImporter{
			State: resourceComputeRouterNatImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: resourceComputeRouterNatDrainNatIpsCustomDiff,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateRFC1035Name(2, 63),
				Description: `Name of the NAT service. The name must be 1-63 characters long and
comply with RFC1035.`,
			},
			"nat_ip_allocate_option": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"MANUAL_ONLY", "AUTO_ONLY"}, false),
				Description: `How external IPs should be allocated for this NAT. Valid values are
'AUTO_ONLY' for only allowing NAT IPs allocated by Google Cloud
Platform, or 'MANUAL_ONLY' for only user-allocated NAT IP addresses.`,
			},
			"router": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      `The name of the Cloud Router in which this NAT will be configured.`,
			},
			"source_subnetwork_ip_ranges_to_nat": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALL_SUBNETWORKS_ALL_IP_RANGES", "ALL_SUBNETWORKS_ALL_PRIMARY_IP_RANGES", "LIST_OF_SUBNETWORKS"}, false),
				Description: `How NAT should be configured per Subnetwork.
If 'ALL_SUBNETWORKS_ALL_IP_RANGES', all of the
IP ranges in every Subnetwork are allowed to Nat.
If 'ALL_SUBNETWORKS_ALL_PRIMARY_IP_RANGES', all of the primary IP
ranges in every Subnetwork are allowed to Nat.
'LIST_OF_SUBNETWORKS': A list of Subnetworks are allowed to Nat
(specified in the field subnetwork below). Note that if this field
contains ALL_SUBNETWORKS_ALL_IP_RANGES or
ALL_SUBNETWORKS_ALL_PRIMARY_IP_RANGES, then there should not be any
other RouterNat section in any Router for this network in this region.`,
			},
			"icmp_idle_timeout_sec": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Timeout (in seconds) for ICMP connections. Defaults to 30s if not set.`,
				Default:     30,
			},
			"log_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Configuration for logging on NAT`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: `Indicates whether or not to export logs.`,
						},
						"filter": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"ERRORS_ONLY", "TRANSLATIONS_ONLY", "ALL"}, false),
							Description: `Specifies the desired filtering of logs on this NAT. Valid
values are: '"ERRORS_ONLY"', '"TRANSLATIONS_ONLY"', '"ALL"'`,
						},
					},
				},
			},
			"min_ports_per_vm": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Minimum number of ports allocated to a VM from this NAT.`,
			},
			"nat_ips": {
				Type:     schema.TypeSet,
				Optional: true,
				Description: `Self-links of NAT IPs. Only valid if natIpAllocateOption
is set to MANUAL_ONLY.`,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					DiffSuppressFunc: compareSelfLinkOrResourceName,
				},
				// Default schema.HashSchema is used.
			},
			"region": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      `Region where the router and NAT reside.`,
			},
			"subnetwork": {
				Type:     schema.TypeSet,
				Optional: true,
				Description: `One or more subnetwork NAT configurations. Only used if
'source_subnetwork_ip_ranges_to_nat' is set to 'LIST_OF_SUBNETWORKS'`,
				Elem: computeRouterNatSubnetworkSchema(),
				Set:  computeRouterNatSubnetworkHash,
			},
			"tcp_established_idle_timeout_sec": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: `Timeout (in seconds) for TCP established connections.
Defaults to 1200s if not set.`,
				Default: 1200,
			},
			"tcp_transitory_idle_timeout_sec": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: `Timeout (in seconds) for TCP transitory connections.
Defaults to 30s if not set.`,
				Default: 30,
			},
			"udp_idle_timeout_sec": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Timeout (in seconds) for UDP connections. Defaults to 30s if not set.`,
				Default:     30,
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

func computeRouterNatSubnetworkSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      `Self-link of subnetwork to NAT`,
			},
			"source_ip_ranges_to_nat": {
				Type:     schema.TypeSet,
				Required: true,
				Description: `List of options for which source IPs in the subnetwork
should have NAT enabled. Supported values include:
'ALL_IP_RANGES', 'LIST_OF_SECONDARY_IP_RANGES',
'PRIMARY_IP_RANGE'.`,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"secondary_ip_range_names": {
				Type:     schema.TypeSet,
				Optional: true,
				Description: `List of the secondary ranges of the subnetwork that are allowed
to use NAT. This can be populated only if
'LIST_OF_SECONDARY_IP_RANGES' is one of the values in
sourceIpRangesToNat`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
		},
	}
}

func resourceComputeRouterNatCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	obj := make(map[string]interface{})
	nameProp, err := expandComputeRouterNatName(d.Get("name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("name"); !isEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	natIpAllocateOptionProp, err := expandComputeRouterNatNatIpAllocateOption(d.Get("nat_ip_allocate_option"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("nat_ip_allocate_option"); !isEmptyValue(reflect.ValueOf(natIpAllocateOptionProp)) && (ok || !reflect.DeepEqual(v, natIpAllocateOptionProp)) {
		obj["natIpAllocateOption"] = natIpAllocateOptionProp
	}
	natIpsProp, err := expandComputeRouterNatNatIps(d.Get("nat_ips"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("nat_ips"); ok || !reflect.DeepEqual(v, natIpsProp) {
		obj["natIps"] = natIpsProp
	}
	sourceSubnetworkIpRangesToNatProp, err := expandComputeRouterNatSourceSubnetworkIpRangesToNat(d.Get("source_subnetwork_ip_ranges_to_nat"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("source_subnetwork_ip_ranges_to_nat"); !isEmptyValue(reflect.ValueOf(sourceSubnetworkIpRangesToNatProp)) && (ok || !reflect.DeepEqual(v, sourceSubnetworkIpRangesToNatProp)) {
		obj["sourceSubnetworkIpRangesToNat"] = sourceSubnetworkIpRangesToNatProp
	}
	subnetworksProp, err := expandComputeRouterNatSubnetwork(d.Get("subnetwork"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("subnetwork"); ok || !reflect.DeepEqual(v, subnetworksProp) {
		obj["subnetworks"] = subnetworksProp
	}
	minPortsPerVmProp, err := expandComputeRouterNatMinPortsPerVm(d.Get("min_ports_per_vm"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("min_ports_per_vm"); !isEmptyValue(reflect.ValueOf(minPortsPerVmProp)) && (ok || !reflect.DeepEqual(v, minPortsPerVmProp)) {
		obj["minPortsPerVm"] = minPortsPerVmProp
	}
	udpIdleTimeoutSecProp, err := expandComputeRouterNatUdpIdleTimeoutSec(d.Get("udp_idle_timeout_sec"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("udp_idle_timeout_sec"); !isEmptyValue(reflect.ValueOf(udpIdleTimeoutSecProp)) && (ok || !reflect.DeepEqual(v, udpIdleTimeoutSecProp)) {
		obj["udpIdleTimeoutSec"] = udpIdleTimeoutSecProp
	}
	icmpIdleTimeoutSecProp, err := expandComputeRouterNatIcmpIdleTimeoutSec(d.Get("icmp_idle_timeout_sec"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("icmp_idle_timeout_sec"); !isEmptyValue(reflect.ValueOf(icmpIdleTimeoutSecProp)) && (ok || !reflect.DeepEqual(v, icmpIdleTimeoutSecProp)) {
		obj["icmpIdleTimeoutSec"] = icmpIdleTimeoutSecProp
	}
	tcpEstablishedIdleTimeoutSecProp, err := expandComputeRouterNatTcpEstablishedIdleTimeoutSec(d.Get("tcp_established_idle_timeout_sec"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("tcp_established_idle_timeout_sec"); !isEmptyValue(reflect.ValueOf(tcpEstablishedIdleTimeoutSecProp)) && (ok || !reflect.DeepEqual(v, tcpEstablishedIdleTimeoutSecProp)) {
		obj["tcpEstablishedIdleTimeoutSec"] = tcpEstablishedIdleTimeoutSecProp
	}
	tcpTransitoryIdleTimeoutSecProp, err := expandComputeRouterNatTcpTransitoryIdleTimeoutSec(d.Get("tcp_transitory_idle_timeout_sec"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("tcp_transitory_idle_timeout_sec"); !isEmptyValue(reflect.ValueOf(tcpTransitoryIdleTimeoutSecProp)) && (ok || !reflect.DeepEqual(v, tcpTransitoryIdleTimeoutSecProp)) {
		obj["tcpTransitoryIdleTimeoutSec"] = tcpTransitoryIdleTimeoutSecProp
	}
	logConfigProp, err := expandComputeRouterNatLogConfig(d.Get("log_config"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("log_config"); !isEmptyValue(reflect.ValueOf(logConfigProp)) && (ok || !reflect.DeepEqual(v, logConfigProp)) {
		obj["logConfig"] = logConfigProp
	}

	lockName, err := replaceVars(d, config, "router/{{region}}/{{router}}")
	if err != nil {
		return err
	}
	mutexKV.Lock(lockName)
	defer mutexKV.Unlock(lockName)

	url, err := replaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/regions/{{region}}/routers/{{router}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new RouterNat: %#v", obj)

	obj, err = resourceComputeRouterNatPatchCreateEncoder(d, meta, obj)
	if err != nil {
		return err
	}
	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	res, err := sendRequestWithTimeout(config, "PATCH", project, url, obj, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating RouterNat: %s", err)
	}

	// Store the ID now
	id, err := replaceVars(d, config, "{{project}}/{{region}}/{{router}}/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = computeOperationWaitTime(
		config, res, project, "Creating RouterNat",
		int(d.Timeout(schema.TimeoutCreate).Minutes()))

	if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error waiting to create RouterNat: %s", err)
	}

	log.Printf("[DEBUG] Finished creating RouterNat %q: %#v", d.Id(), res)

	return resourceComputeRouterNatRead(d, meta)
}

func resourceComputeRouterNatRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	url, err := replaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/regions/{{region}}/routers/{{router}}")
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	res, err := sendRequest(config, "GET", project, url, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("ComputeRouterNat %q", d.Id()))
	}

	res, err = flattenNestedComputeRouterNat(d, meta, res)
	if err != nil {
		return err
	}

	if res == nil {
		// Object isn't there any more - remove it from the state.
		log.Printf("[DEBUG] Removing ComputeRouterNat because it couldn't be matched.")
		d.SetId("")
		return nil
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading RouterNat: %s", err)
	}

	if err := d.Set("name", flattenComputeRouterNatName(res["name"], d)); err != nil {
		return fmt.Errorf("Error reading RouterNat: %s", err)
	}
	if err := d.Set("nat_ip_allocate_option", flattenComputeRouterNatNatIpAllocateOption(res["natIpAllocateOption"], d)); err != nil {
		return fmt.Errorf("Error reading RouterNat: %s", err)
	}
	if err := d.Set("nat_ips", flattenComputeRouterNatNatIps(res["natIps"], d)); err != nil {
		return fmt.Errorf("Error reading RouterNat: %s", err)
	}
	if err := d.Set("source_subnetwork_ip_ranges_to_nat", flattenComputeRouterNatSourceSubnetworkIpRangesToNat(res["sourceSubnetworkIpRangesToNat"], d)); err != nil {
		return fmt.Errorf("Error reading RouterNat: %s", err)
	}
	if err := d.Set("subnetwork", flattenComputeRouterNatSubnetwork(res["subnetworks"], d)); err != nil {
		return fmt.Errorf("Error reading RouterNat: %s", err)
	}
	if err := d.Set("min_ports_per_vm", flattenComputeRouterNatMinPortsPerVm(res["minPortsPerVm"], d)); err != nil {
		return fmt.Errorf("Error reading RouterNat: %s", err)
	}
	if err := d.Set("udp_idle_timeout_sec", flattenComputeRouterNatUdpIdleTimeoutSec(res["udpIdleTimeoutSec"], d)); err != nil {
		return fmt.Errorf("Error reading RouterNat: %s", err)
	}
	if err := d.Set("icmp_idle_timeout_sec", flattenComputeRouterNatIcmpIdleTimeoutSec(res["icmpIdleTimeoutSec"], d)); err != nil {
		return fmt.Errorf("Error reading RouterNat: %s", err)
	}
	if err := d.Set("tcp_established_idle_timeout_sec", flattenComputeRouterNatTcpEstablishedIdleTimeoutSec(res["tcpEstablishedIdleTimeoutSec"], d)); err != nil {
		return fmt.Errorf("Error reading RouterNat: %s", err)
	}
	if err := d.Set("tcp_transitory_idle_timeout_sec", flattenComputeRouterNatTcpTransitoryIdleTimeoutSec(res["tcpTransitoryIdleTimeoutSec"], d)); err != nil {
		return fmt.Errorf("Error reading RouterNat: %s", err)
	}
	if err := d.Set("log_config", flattenComputeRouterNatLogConfig(res["logConfig"], d)); err != nil {
		return fmt.Errorf("Error reading RouterNat: %s", err)
	}

	return nil
}

func resourceComputeRouterNatUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	natIpAllocateOptionProp, err := expandComputeRouterNatNatIpAllocateOption(d.Get("nat_ip_allocate_option"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("nat_ip_allocate_option"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, natIpAllocateOptionProp)) {
		obj["natIpAllocateOption"] = natIpAllocateOptionProp
	}
	natIpsProp, err := expandComputeRouterNatNatIps(d.Get("nat_ips"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("nat_ips"); ok || !reflect.DeepEqual(v, natIpsProp) {
		obj["natIps"] = natIpsProp
	}
	sourceSubnetworkIpRangesToNatProp, err := expandComputeRouterNatSourceSubnetworkIpRangesToNat(d.Get("source_subnetwork_ip_ranges_to_nat"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("source_subnetwork_ip_ranges_to_nat"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, sourceSubnetworkIpRangesToNatProp)) {
		obj["sourceSubnetworkIpRangesToNat"] = sourceSubnetworkIpRangesToNatProp
	}
	subnetworksProp, err := expandComputeRouterNatSubnetwork(d.Get("subnetwork"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("subnetwork"); ok || !reflect.DeepEqual(v, subnetworksProp) {
		obj["subnetworks"] = subnetworksProp
	}
	minPortsPerVmProp, err := expandComputeRouterNatMinPortsPerVm(d.Get("min_ports_per_vm"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("min_ports_per_vm"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, minPortsPerVmProp)) {
		obj["minPortsPerVm"] = minPortsPerVmProp
	}
	udpIdleTimeoutSecProp, err := expandComputeRouterNatUdpIdleTimeoutSec(d.Get("udp_idle_timeout_sec"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("udp_idle_timeout_sec"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, udpIdleTimeoutSecProp)) {
		obj["udpIdleTimeoutSec"] = udpIdleTimeoutSecProp
	}
	icmpIdleTimeoutSecProp, err := expandComputeRouterNatIcmpIdleTimeoutSec(d.Get("icmp_idle_timeout_sec"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("icmp_idle_timeout_sec"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, icmpIdleTimeoutSecProp)) {
		obj["icmpIdleTimeoutSec"] = icmpIdleTimeoutSecProp
	}
	tcpEstablishedIdleTimeoutSecProp, err := expandComputeRouterNatTcpEstablishedIdleTimeoutSec(d.Get("tcp_established_idle_timeout_sec"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("tcp_established_idle_timeout_sec"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, tcpEstablishedIdleTimeoutSecProp)) {
		obj["tcpEstablishedIdleTimeoutSec"] = tcpEstablishedIdleTimeoutSecProp
	}
	tcpTransitoryIdleTimeoutSecProp, err := expandComputeRouterNatTcpTransitoryIdleTimeoutSec(d.Get("tcp_transitory_idle_timeout_sec"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("tcp_transitory_idle_timeout_sec"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, tcpTransitoryIdleTimeoutSecProp)) {
		obj["tcpTransitoryIdleTimeoutSec"] = tcpTransitoryIdleTimeoutSecProp
	}
	logConfigProp, err := expandComputeRouterNatLogConfig(d.Get("log_config"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("log_config"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, logConfigProp)) {
		obj["logConfig"] = logConfigProp
	}

	lockName, err := replaceVars(d, config, "router/{{region}}/{{router}}")
	if err != nil {
		return err
	}
	mutexKV.Lock(lockName)
	defer mutexKV.Unlock(lockName)

	url, err := replaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/regions/{{region}}/routers/{{router}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating RouterNat %q: %#v", d.Id(), obj)

	obj, err = resourceComputeRouterNatPatchUpdateEncoder(d, meta, obj)
	if err != nil {
		return err
	}
	res, err := sendRequestWithTimeout(config, "PATCH", project, url, obj, d.Timeout(schema.TimeoutUpdate))

	if err != nil {
		return fmt.Errorf("Error updating RouterNat %q: %s", d.Id(), err)
	}

	err = computeOperationWaitTime(
		config, res, project, "Updating RouterNat",
		int(d.Timeout(schema.TimeoutUpdate).Minutes()))

	if err != nil {
		return err
	}

	return resourceComputeRouterNatRead(d, meta)
}

func resourceComputeRouterNatDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	lockName, err := replaceVars(d, config, "router/{{region}}/{{router}}")
	if err != nil {
		return err
	}
	mutexKV.Lock(lockName)
	defer mutexKV.Unlock(lockName)

	url, err := replaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/regions/{{region}}/routers/{{router}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}

	obj, err = resourceComputeRouterNatPatchDeleteEncoder(d, meta, obj)
	if err != nil {
		return handleNotFoundError(err, d, "RouterNat")
	}
	log.Printf("[DEBUG] Deleting RouterNat %q", d.Id())

	res, err := sendRequestWithTimeout(config, "PATCH", project, url, obj, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return handleNotFoundError(err, d, "RouterNat")
	}

	err = computeOperationWaitTime(
		config, res, project, "Deleting RouterNat",
		int(d.Timeout(schema.TimeoutDelete).Minutes()))

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting RouterNat %q: %#v", d.Id(), res)
	return nil
}

func resourceComputeRouterNatImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/regions/(?P<region>[^/]+)/routers/(?P<router>[^/]+)/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<region>[^/]+)/(?P<router>[^/]+)/(?P<name>[^/]+)",
		"(?P<region>[^/]+)/(?P<router>[^/]+)/(?P<name>[^/]+)",
		"(?P<router>[^/]+)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "{{project}}/{{region}}/{{router}}/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenComputeRouterNatName(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenComputeRouterNatNatIpAllocateOption(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenComputeRouterNatNatIps(v interface{}, d *schema.ResourceData) interface{} {
	if v == nil {
		return v
	}
	return convertAndMapStringArr(v.([]interface{}), ConvertSelfLinkToV1)
}

func flattenComputeRouterNatSourceSubnetworkIpRangesToNat(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenComputeRouterNatSubnetwork(v interface{}, d *schema.ResourceData) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := schema.NewSet(computeRouterNatSubnetworkHash, []interface{}{})
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed.Add(map[string]interface{}{
			"name":                     flattenComputeRouterNatSubnetworkName(original["name"], d),
			"source_ip_ranges_to_nat":  flattenComputeRouterNatSubnetworkSourceIpRangesToNat(original["sourceIpRangesToNat"], d),
			"secondary_ip_range_names": flattenComputeRouterNatSubnetworkSecondaryIpRangeNames(original["secondaryIpRangeNames"], d),
		})
	}
	return transformed
}
func flattenComputeRouterNatSubnetworkName(v interface{}, d *schema.ResourceData) interface{} {
	if v == nil {
		return v
	}
	return ConvertSelfLinkToV1(v.(string))
}

func flattenComputeRouterNatSubnetworkSourceIpRangesToNat(v interface{}, d *schema.ResourceData) interface{} {
	if v == nil {
		return v
	}
	return schema.NewSet(schema.HashString, v.([]interface{}))
}

func flattenComputeRouterNatSubnetworkSecondaryIpRangeNames(v interface{}, d *schema.ResourceData) interface{} {
	if v == nil {
		return v
	}
	return schema.NewSet(schema.HashString, v.([]interface{}))
}

func flattenComputeRouterNatMinPortsPerVm(v interface{}, d *schema.ResourceData) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := strconv.ParseInt(strVal, 10, 64); err == nil {
			return intVal
		} // let terraform core handle it if we can't convert the string to an int.
	}
	return v
}

func flattenComputeRouterNatUdpIdleTimeoutSec(v interface{}, d *schema.ResourceData) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := strconv.ParseInt(strVal, 10, 64); err == nil {
			return intVal
		} // let terraform core handle it if we can't convert the string to an int.
	}
	return v
}

func flattenComputeRouterNatIcmpIdleTimeoutSec(v interface{}, d *schema.ResourceData) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := strconv.ParseInt(strVal, 10, 64); err == nil {
			return intVal
		} // let terraform core handle it if we can't convert the string to an int.
	}
	return v
}

func flattenComputeRouterNatTcpEstablishedIdleTimeoutSec(v interface{}, d *schema.ResourceData) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := strconv.ParseInt(strVal, 10, 64); err == nil {
			return intVal
		} // let terraform core handle it if we can't convert the string to an int.
	}
	return v
}

func flattenComputeRouterNatTcpTransitoryIdleTimeoutSec(v interface{}, d *schema.ResourceData) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := strconv.ParseInt(strVal, 10, 64); err == nil {
			return intVal
		} // let terraform core handle it if we can't convert the string to an int.
	}
	return v
}

func flattenComputeRouterNatLogConfig(v interface{}, d *schema.ResourceData) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["enable"] =
		flattenComputeRouterNatLogConfigEnable(original["enable"], d)
	transformed["filter"] =
		flattenComputeRouterNatLogConfigFilter(original["filter"], d)
	return []interface{}{transformed}
}
func flattenComputeRouterNatLogConfigEnable(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenComputeRouterNatLogConfigFilter(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandComputeRouterNatName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeRouterNatNatIpAllocateOption(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeRouterNatNatIps(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	v = v.(*schema.Set).List()
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			return nil, fmt.Errorf("Invalid value for nat_ips: nil")
		}
		f, err := parseRegionalFieldValue("addresses", raw.(string), "project", "region", "zone", d, config, true)
		if err != nil {
			return nil, fmt.Errorf("Invalid value for nat_ips: %s", err)
		}
		req = append(req, f.RelativeLink())
	}
	return req, nil
}

func expandComputeRouterNatSourceSubnetworkIpRangesToNat(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeRouterNatSubnetwork(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	v = v.(*schema.Set).List()
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedName, err := expandComputeRouterNatSubnetworkName(original["name"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedName); val.IsValid() && !isEmptyValue(val) {
			transformed["name"] = transformedName
		}

		transformedSourceIpRangesToNat, err := expandComputeRouterNatSubnetworkSourceIpRangesToNat(original["source_ip_ranges_to_nat"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedSourceIpRangesToNat); val.IsValid() && !isEmptyValue(val) {
			transformed["sourceIpRangesToNat"] = transformedSourceIpRangesToNat
		}

		transformedSecondaryIpRangeNames, err := expandComputeRouterNatSubnetworkSecondaryIpRangeNames(original["secondary_ip_range_names"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedSecondaryIpRangeNames); val.IsValid() && !isEmptyValue(val) {
			transformed["secondaryIpRangeNames"] = transformedSecondaryIpRangeNames
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandComputeRouterNatSubnetworkName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	f, err := parseRegionalFieldValue("subnetworks", v.(string), "project", "region", "zone", d, config, true)
	if err != nil {
		return nil, fmt.Errorf("Invalid value for name: %s", err)
	}
	return f.RelativeLink(), nil
}

func expandComputeRouterNatSubnetworkSourceIpRangesToNat(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	v = v.(*schema.Set).List()
	return v, nil
}

func expandComputeRouterNatSubnetworkSecondaryIpRangeNames(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	v = v.(*schema.Set).List()
	return v, nil
}

func expandComputeRouterNatMinPortsPerVm(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeRouterNatUdpIdleTimeoutSec(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeRouterNatIcmpIdleTimeoutSec(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeRouterNatTcpEstablishedIdleTimeoutSec(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeRouterNatTcpTransitoryIdleTimeoutSec(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeRouterNatLogConfig(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedEnable, err := expandComputeRouterNatLogConfigEnable(original["enable"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedEnable); val.IsValid() && !isEmptyValue(val) {
		transformed["enable"] = transformedEnable
	}

	transformedFilter, err := expandComputeRouterNatLogConfigFilter(original["filter"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedFilter); val.IsValid() && !isEmptyValue(val) {
		transformed["filter"] = transformedFilter
	}

	return transformed, nil
}

func expandComputeRouterNatLogConfigEnable(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeRouterNatLogConfigFilter(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func flattenNestedComputeRouterNat(d *schema.ResourceData, meta interface{}, res map[string]interface{}) (map[string]interface{}, error) {
	var v interface{}
	var ok bool

	v, ok = res["nats"]
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
		return nil, fmt.Errorf("expected list or map for value nats. Actual value: %v", v)
	}

	_, item, err := resourceComputeRouterNatFindNestedObjectInList(d, meta, v.([]interface{}))
	if err != nil {
		return nil, err
	}
	return item, nil
}

func resourceComputeRouterNatFindNestedObjectInList(d *schema.ResourceData, meta interface{}, items []interface{}) (index int, item map[string]interface{}, err error) {
	expectedName, err := expandComputeRouterNatName(d.Get("name"), d, meta.(*Config))
	if err != nil {
		return -1, nil, err
	}

	// Search list for this resource.
	for idx, itemRaw := range items {
		if itemRaw == nil {
			continue
		}
		item := itemRaw.(map[string]interface{})

		itemName := flattenComputeRouterNatName(item["name"], d)
		if !reflect.DeepEqual(itemName, expectedName) {
			log.Printf("[DEBUG] Skipping item with name= %#v, looking for %#v)", itemName, expectedName)
			continue
		}
		log.Printf("[DEBUG] Found item for resource %q: %#v)", d.Id(), item)
		return idx, item, nil
	}
	return -1, nil, nil
}

// PatchCreateEncoder handles creating request data to PATCH parent resource
// with list including new object.
func resourceComputeRouterNatPatchCreateEncoder(d *schema.ResourceData, meta interface{}, obj map[string]interface{}) (map[string]interface{}, error) {
	currItems, err := resourceComputeRouterNatListForPatch(d, meta)
	if err != nil {
		return nil, err
	}

	_, found, err := resourceComputeRouterNatFindNestedObjectInList(d, meta, currItems)
	if err != nil {
		return nil, err
	}

	// Return error if item already created.
	if found != nil {
		return nil, fmt.Errorf("Unable to create RouterNat, existing object already found: %+v", found)
	}

	// Return list with the resource to create appended
	return map[string]interface{}{
		"nats": append(currItems, obj),
	}, nil
}

// PatchUpdateEncoder handles creating request data to PATCH parent resource
// with list including updated object.
func resourceComputeRouterNatPatchUpdateEncoder(d *schema.ResourceData, meta interface{}, obj map[string]interface{}) (map[string]interface{}, error) {
	items, err := resourceComputeRouterNatListForPatch(d, meta)
	if err != nil {
		return nil, err
	}

	idx, item, err := resourceComputeRouterNatFindNestedObjectInList(d, meta, items)
	if err != nil {
		return nil, err
	}

	// Return error if item to update does not exist.
	if item == nil {
		return nil, fmt.Errorf("Unable to update RouterNat %q - not found in list", d.Id())
	}

	// Merge new object into old.
	for k, v := range obj {
		item[k] = v
	}
	items[idx] = item

	// Return list with new item added
	return map[string]interface{}{
		"nats": items,
	}, nil
}

// PatchDeleteEncoder handles creating request data to PATCH parent resource
// with list excluding object to delete.
func resourceComputeRouterNatPatchDeleteEncoder(d *schema.ResourceData, meta interface{}, obj map[string]interface{}) (map[string]interface{}, error) {
	currItems, err := resourceComputeRouterNatListForPatch(d, meta)
	if err != nil {
		return nil, err
	}

	idx, item, err := resourceComputeRouterNatFindNestedObjectInList(d, meta, currItems)
	if err != nil {
		return nil, err
	}
	if item == nil {
		// Spoof 404 error for proper handling by Delete (i.e. no-op)
		return nil, &googleapi.Error{
			Code:    404,
			Message: "RouterNat not found in list",
		}
	}

	updatedItems := append(currItems[:idx], currItems[idx+1:]...)
	return map[string]interface{}{
		"nats": updatedItems,
	}, nil
}

// ListForPatch handles making API request to get parent resource and
// extracting list of objects.
func resourceComputeRouterNatListForPatch(d *schema.ResourceData, meta interface{}) ([]interface{}, error) {
	config := meta.(*Config)
	url, err := replaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/regions/{{region}}/routers/{{router}}")
	if err != nil {
		return nil, err
	}
	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}
	res, err := sendRequest(config, "GET", project, url, nil)
	if err != nil {
		return nil, err
	}

	v, ok := res["nats"]
	if ok && v != nil {
		ls, lsOk := v.([]interface{})
		if !lsOk {
			return nil, fmt.Errorf(`expected list for nested field "nats"`)
		}
		return ls, nil
	}
	return nil, nil
}
