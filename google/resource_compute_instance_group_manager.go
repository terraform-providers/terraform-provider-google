package google

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"

	computeBeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
)

func resourceComputeInstanceGroupManager() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeInstanceGroupManagerCreate,
		Read:   resourceComputeInstanceGroupManagerRead,
		Update: resourceComputeInstanceGroupManagerUpdate,
		Delete: resourceComputeInstanceGroupManagerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"base_instance_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"instance_template": &schema.Schema{
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: compareSelfLinkRelativePaths,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"zone": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"fingerprint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"instance_group": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"named_port": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"port": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"self_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"update_strategy": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "RESTART",
			},

			"target_pools": &schema.Schema{
				Type:             schema.TypeSet,
				Optional:         true,
				DiffSuppressFunc: compareSelfLinkRelativePaths,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: selfLinkRelativePathHash,
			},

			"target_size": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func getNamedPorts(nps []interface{}) []*compute.NamedPort {
	namedPorts := make([]*compute.NamedPort, 0, len(nps))
	for _, v := range nps {
		np := v.(map[string]interface{})
		namedPorts = append(namedPorts, &compute.NamedPort{
			Name: np["name"].(string),
			Port: int64(np["port"].(int)),
		})
	}

	return namedPorts
}

func getNamedPortsBeta(nps []interface{}) []*computeBeta.NamedPort {
	namedPorts := make([]*computeBeta.NamedPort, 0, len(nps))
	for _, v := range nps {
		np := v.(map[string]interface{})
		namedPorts = append(namedPorts, &computeBeta.NamedPort{
			Name: np["name"].(string),
			Port: int64(np["port"].(int)),
		})
	}

	return namedPorts
}

var InstanceGroupManagerBaseVersion = v1

func resourceComputeInstanceGroupManagerCreate(d *schema.ResourceData, meta interface{}) error {
	computeApiVersion := getComputeApiVersion(d, InstanceGroupManagerBaseVersion, []Feature{})
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Get group size, default to 1 if not given
	var target_size int64 = 1
	if v, ok := d.GetOk("target_size"); ok {
		target_size = int64(v.(int))
	}

	// Build the parameter
	manager := &computeBeta.InstanceGroupManager{
		Name:             d.Get("name").(string),
		BaseInstanceName: d.Get("base_instance_name").(string),
		InstanceTemplate: d.Get("instance_template").(string),
		TargetSize:       target_size,
	}

	// Set optional fields
	if v, ok := d.GetOk("description"); ok {
		manager.Description = v.(string)
	}

	if v, ok := d.GetOk("named_port"); ok {
		manager.NamedPorts = getNamedPortsBeta(v.([]interface{}))
	}

	if attr := d.Get("target_pools").(*schema.Set); attr.Len() > 0 {
		var s []string
		for _, v := range attr.List() {
			s = append(s, v.(string))
		}
		manager.TargetPools = s
	}

	updateStrategy := d.Get("update_strategy").(string)
	if !(updateStrategy == "NONE" || updateStrategy == "RESTART") {
		return fmt.Errorf("Update strategy must be \"NONE\" or \"RESTART\"")
	}

	log.Printf("[DEBUG] InstanceGroupManager insert request: %#v", manager)
	op, err := config.clientComputeMultiversion.InsertInstanceGroupManager(project, d.Get("zone").(string), manager, computeApiVersion)
	if err != nil {
		return fmt.Errorf("Error creating InstanceGroupManager: %s", err)
	}

	// It probably maybe worked, so store the ID now
	d.SetId(manager.Name)

	// Wait for the operation to complete
	err = computeSharedOperationWaitZone(config, op, project, d.Get("zone").(string), "Creating InstanceGroupManager")
	if err != nil {
		return err
	}

	return resourceComputeInstanceGroupManagerRead(d, meta)
}

func flattenNamedPortsBeta(namedPorts []*computeBeta.NamedPort) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(namedPorts))
	for _, namedPort := range namedPorts {
		namedPortMap := make(map[string]interface{})
		namedPortMap["name"] = namedPort.Name
		namedPortMap["port"] = namedPort.Port
		result = append(result, namedPortMap)
	}
	return result

}

func resourceComputeInstanceGroupManagerRead(d *schema.ResourceData, meta interface{}) error {
	computeApiVersion := getComputeApiVersion(d, InstanceGroupManagerBaseVersion, []Feature{})
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	manager := &computeBeta.InstanceGroupManager{}
	switch computeApiVersion {
	case v1:
		getInstanceGroupManager := func(zone string) (interface{}, error) {
			return config.clientCompute.InstanceGroupManagers.Get(project, zone, d.Id()).Do()
		}

		var v1Manager *compute.InstanceGroupManager
		var e error
		if zone, ok := d.GetOk("zone"); ok {
			v1Manager, e = config.clientCompute.InstanceGroupManagers.Get(project, zone.(string), d.Id()).Do()

			if e != nil {
				return handleNotFoundError(e, d, fmt.Sprintf("Instance Group Manager %q", d.Get("name").(string)))
			}
		} else {
			// If the resource was imported, the only info we have is the ID. Try to find the resource
			// by searching in the region of the project.
			var resource interface{}
			resource, e = getZonalResourceFromRegion(getInstanceGroupManager, region, config.clientCompute, project)

			if e != nil {
				return e
			}

			v1Manager = resource.(*compute.InstanceGroupManager)
		}

		if v1Manager == nil {
			log.Printf("[WARN] Removing Instance Group Manager %q because it's gone", d.Get("name").(string))

			// The resource doesn't exist anymore
			d.SetId("")
			return nil
		}

		err = Convert(v1Manager, manager)
		if err != nil {
			return err
		}

	case v0beta:
		getInstanceGroupManager := func(zone string) (interface{}, error) {
			return config.clientComputeBeta.InstanceGroupManagers.Get(project, zone, d.Id()).Do()
		}

		var v0betaManager *computeBeta.InstanceGroupManager
		var e error
		if zone, ok := d.GetOk("zone"); ok {
			v0betaManager, e = config.clientComputeBeta.InstanceGroupManagers.Get(project, zone.(string), d.Id()).Do()

			if e != nil {
				return handleNotFoundError(e, d, fmt.Sprintf("Instance Group Manager %q", d.Get("name").(string)))
			}
		} else {
			// If the resource was imported, the only info we have is the ID. Try to find the resource
			// by searching in the region of the project.
			var resource interface{}
			resource, e = getZonalBetaResourceFromRegion(getInstanceGroupManager, region, config.clientComputeBeta, project)
			if e != nil {
				return e
			}

			v0betaManager = resource.(*computeBeta.InstanceGroupManager)
		}

		if v0betaManager == nil {
			log.Printf("[WARN] Removing Instance Group Manager %q because it's gone", d.Get("name").(string))

			// The resource doesn't exist anymore
			d.SetId("")
			return nil
		}

		manager = v0betaManager
	}

	zoneUrl := strings.Split(manager.Zone, "/")
	d.Set("base_instance_name", manager.BaseInstanceName)
	d.Set("instance_template", manager.InstanceTemplate)
	d.Set("name", manager.Name)
	d.Set("zone", zoneUrl[len(zoneUrl)-1])
	d.Set("description", manager.Description)
	d.Set("project", project)
	d.Set("target_size", manager.TargetSize)
	d.Set("target_pools", manager.TargetPools)
	d.Set("named_port", flattenNamedPortsBeta(manager.NamedPorts))
	d.Set("fingerprint", manager.Fingerprint)
	d.Set("instance_group", manager.InstanceGroup)
	d.Set("target_size", manager.TargetSize)
	d.Set("self_link", manager.SelfLink)
	update_strategy, ok := d.GetOk("update_strategy")
	if !ok {
		update_strategy = "RESTART"
	}
	d.Set("update_strategy", update_strategy.(string))

	return nil
}

func resourceComputeInstanceGroupManagerUpdate(d *schema.ResourceData, meta interface{}) error {
	computeApiVersion := getComputeApiVersionUpdate(d, InstanceGroupManagerBaseVersion, []Feature{}, []Feature{})
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	d.Partial(true)

	// If target_pools changes then update
	if d.HasChange("target_pools") {
		var targetPools []string
		if attr := d.Get("target_pools").(*schema.Set); attr.Len() > 0 {
			for _, v := range attr.List() {
				targetPools = append(targetPools, v.(string))
			}
		}

		// Build the parameter
		setTargetPools := &computeBeta.InstanceGroupManagersSetTargetPoolsRequest{
			Fingerprint: d.Get("fingerprint").(string),
			TargetPools: targetPools,
		}

		var op interface{}
		switch computeApiVersion {
		case v1:
			setTargetPoolsV1 := &compute.InstanceGroupManagersSetTargetPoolsRequest{}
			err := Convert(setTargetPools, setTargetPoolsV1)
			if err != nil {
				return err
			}

			op, err = config.clientCompute.InstanceGroupManagers.SetTargetPools(
				project, d.Get("zone").(string), d.Id(), setTargetPoolsV1).Do()
		case v0beta:
			setTargetPoolsV0beta := &computeBeta.InstanceGroupManagersSetTargetPoolsRequest{}
			err := Convert(setTargetPools, setTargetPoolsV0beta)
			if err != nil {
				return err
			}

			op, err = config.clientComputeBeta.InstanceGroupManagers.SetTargetPools(
				project, d.Get("zone").(string), d.Id(), setTargetPoolsV0beta).Do()
		}

		if err != nil {
			return fmt.Errorf("Error updating InstanceGroupManager: %s", err)
		}

		// Wait for the operation to complete
		err = computeSharedOperationWaitZone(config, op, project, d.Get("zone").(string), "Updating InstanceGroupManager")
		if err != nil {
			return err
		}

		d.SetPartial("target_pools")
	}

	// If instance_template changes then update
	if d.HasChange("instance_template") {
		// Build the parameter
		setInstanceTemplate := &computeBeta.InstanceGroupManagersSetInstanceTemplateRequest{
			InstanceTemplate: d.Get("instance_template").(string),
		}

		var op interface{}
		switch computeApiVersion {
		case v1:
			setInstanceTemplateV1 := &compute.InstanceGroupManagersSetInstanceTemplateRequest{}
			err := Convert(setInstanceTemplate, setInstanceTemplateV1)
			if err != nil {
				return err
			}

			op, err = config.clientCompute.InstanceGroupManagers.SetInstanceTemplate(
				project, d.Get("zone").(string), d.Id(), setInstanceTemplateV1).Do()
		case v0beta:
			setInstanceTemplateV0beta := &computeBeta.InstanceGroupManagersSetInstanceTemplateRequest{}
			err := Convert(setInstanceTemplate, setInstanceTemplateV0beta)
			if err != nil {
				return err
			}

			op, err = config.clientComputeBeta.InstanceGroupManagers.SetInstanceTemplate(
				project, d.Get("zone").(string), d.Id(), setInstanceTemplateV0beta).Do()
		}

		if err != nil {
			return fmt.Errorf("Error updating InstanceGroupManager: %s", err)
		}

		// Wait for the operation to complete
		err = computeSharedOperationWaitZone(config, op, project, d.Get("zone").(string), "Updating InstanceGroupManager")
		if err != nil {
			return err
		}

		if d.Get("update_strategy").(string) == "RESTART" {
			managedInstances := &computeBeta.InstanceGroupManagersListManagedInstancesResponse{}
			switch computeApiVersion {
			case v1:
				managedInstancesV1, err := config.clientCompute.InstanceGroupManagers.ListManagedInstances(
					project, d.Get("zone").(string), d.Id()).Do()
				if err != nil {
					return fmt.Errorf("Error getting instance group managers instances: %s", err)
				}

				err = Convert(managedInstancesV1, managedInstances)
				if err != nil {
					return err
				}
			case v0beta:
				managedInstancesV0beta, err := config.clientComputeBeta.InstanceGroupManagers.ListManagedInstances(
					project, d.Get("zone").(string), d.Id()).Do()
				if err != nil {
					return fmt.Errorf("Error getting instance group managers instances: %s", err)
				}

				err = Convert(managedInstancesV0beta, managedInstances)
				if err != nil {
					return err
				}
			}

			managedInstanceCount := len(managedInstances.ManagedInstances)
			instances := make([]string, managedInstanceCount)
			for i, v := range managedInstances.ManagedInstances {
				instances[i] = v.Instance
			}

			recreateInstances := &computeBeta.InstanceGroupManagersRecreateInstancesRequest{
				Instances: instances,
			}

			var op interface{}
			switch computeApiVersion {
			case v1:
				recreateInstancesV1 := &compute.InstanceGroupManagersRecreateInstancesRequest{}
				err := Convert(recreateInstances, recreateInstancesV1)
				if err != nil {
					return err
				}

				op, err = config.clientCompute.InstanceGroupManagers.RecreateInstances(
					project, d.Get("zone").(string), d.Id(), recreateInstancesV1).Do()
				if err != nil {
					return fmt.Errorf("Error restarting instance group managers instances: %s", err)
				}
			case v0beta:
				recreateInstancesV0beta := &computeBeta.InstanceGroupManagersRecreateInstancesRequest{}
				err := Convert(recreateInstances, recreateInstancesV0beta)
				if err != nil {
					return err
				}

				op, err = config.clientComputeBeta.InstanceGroupManagers.RecreateInstances(
					project, d.Get("zone").(string), d.Id(), recreateInstancesV0beta).Do()
				if err != nil {
					return fmt.Errorf("Error restarting instance group managers instances: %s", err)
				}
			}

			// Wait for the operation to complete
			err = computeSharedOperationWaitZoneTime(config, op, project, d.Get("zone").(string),
				managedInstanceCount*4, "Restarting InstanceGroupManagers instances")
			if err != nil {
				return err
			}
		}

		d.SetPartial("instance_template")
	}

	// If named_port changes then update:
	if d.HasChange("named_port") {

		// Build the parameters for a "SetNamedPorts" request:
		namedPorts := getNamedPortsBeta(d.Get("named_port").([]interface{}))
		setNamedPorts := &computeBeta.InstanceGroupsSetNamedPortsRequest{
			NamedPorts: namedPorts,
		}

		// Make the request:
		var op interface{}
		switch computeApiVersion {
		case v1:
			setNamedPortsV1 := &compute.InstanceGroupsSetNamedPortsRequest{}
			err := Convert(setNamedPorts, setNamedPortsV1)
			if err != nil {
				return err
			}

			op, err = config.clientCompute.InstanceGroups.SetNamedPorts(
				project, d.Get("zone").(string), d.Id(), setNamedPortsV1).Do()
		case v0beta:
			setNamedPortsV0beta := &computeBeta.InstanceGroupsSetNamedPortsRequest{}
			err := Convert(setNamedPorts, setNamedPortsV0beta)
			if err != nil {
				return err
			}

			op, err = config.clientComputeBeta.InstanceGroups.SetNamedPorts(
				project, d.Get("zone").(string), d.Id(), setNamedPortsV0beta).Do()
		}

		if err != nil {
			return fmt.Errorf("Error updating InstanceGroupManager: %s", err)
		}

		// Wait for the operation to complete:
		err = computeSharedOperationWaitZone(config, op, project, d.Get("zone").(string), "Updating InstanceGroupManager")
		if err != nil {
			return err
		}

		d.SetPartial("named_port")
	}

	// If size changes trigger a resize
	if d.HasChange("target_size") {
		if v, ok := d.GetOk("target_size"); ok {
			// Only do anything if the new size is set
			target_size := int64(v.(int))

			var op interface{}
			switch computeApiVersion {
			case v1:
				op, err = config.clientCompute.InstanceGroupManagers.Resize(
					project, d.Get("zone").(string), d.Id(), target_size).Do()
			case v0beta:
				op, err = config.clientComputeBeta.InstanceGroupManagers.Resize(
					project, d.Get("zone").(string), d.Id(), target_size).Do()
			}

			if err != nil {
				return fmt.Errorf("Error updating InstanceGroupManager: %s", err)
			}

			// Wait for the operation to complete
			err = computeSharedOperationWaitZone(config, op, project, d.Get("zone").(string), "Updating InstanceGroupManager")
			if err != nil {
				return err
			}
		}

		d.SetPartial("target_size")
	}

	d.Partial(false)

	return resourceComputeInstanceGroupManagerRead(d, meta)
}

func resourceComputeInstanceGroupManagerDelete(d *schema.ResourceData, meta interface{}) error {
	computeApiVersion := getComputeApiVersion(d, InstanceGroupManagerBaseVersion, []Feature{})
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	zone := d.Get("zone").(string)

	var op interface{}
	switch computeApiVersion {
	case v1:
		op, err = config.clientCompute.InstanceGroupManagers.Delete(project, zone, d.Id()).Do()
		attempt := 0
		for err != nil && attempt < 20 {
			attempt++
			time.Sleep(2000 * time.Millisecond)
			op, err = config.clientCompute.InstanceGroupManagers.Delete(project, zone, d.Id()).Do()
		}
	case v0beta:
		op, err = config.clientComputeBeta.InstanceGroupManagers.Delete(project, zone, d.Id()).Do()
		attempt := 0
		for err != nil && attempt < 20 {
			attempt++
			time.Sleep(2000 * time.Millisecond)
			op, err = config.clientComputeBeta.InstanceGroupManagers.Delete(project, zone, d.Id()).Do()
		}
	}

	if err != nil {
		return fmt.Errorf("Error deleting instance group manager: %s", err)
	}

	currentSize := int64(d.Get("target_size").(int))

	// Wait for the operation to complete
	err = computeSharedOperationWaitZone(config, op, project, d.Get("zone").(string), "Deleting InstanceGroupManager")

	for err != nil && currentSize > 0 {
		if !strings.Contains(err.Error(), "timeout") {
			return err
		}

		var instanceGroupSize int64
		switch computeApiVersion {
		case v1:
			instanceGroup, err := config.clientCompute.InstanceGroups.Get(
				project, d.Get("zone").(string), d.Id()).Do()
			if err != nil {
				return fmt.Errorf("Error getting instance group size: %s", err)
			}

			instanceGroupSize = instanceGroup.Size
		case v0beta:
			instanceGroup, err := config.clientComputeBeta.InstanceGroups.Get(
				project, d.Get("zone").(string), d.Id()).Do()
			if err != nil {
				return fmt.Errorf("Error getting instance group size: %s", err)
			}

			instanceGroupSize = instanceGroup.Size
		}

		if instanceGroupSize >= currentSize {
			return fmt.Errorf("Error, instance group isn't shrinking during delete")
		}

		log.Printf("[INFO] timeout occured, but instance group is shrinking (%d < %d)", instanceGroupSize, currentSize)
		currentSize = instanceGroupSize
		err = computeSharedOperationWaitZone(config, op, project, d.Get("zone").(string), "Deleting InstanceGroupManager")
	}

	d.SetId("")
	return nil
}
