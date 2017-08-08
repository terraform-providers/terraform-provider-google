package google

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"strings"

	compute "google.golang.org/api/compute/v1"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/terraform"
)

func resourceComputeInstanceMigrateState(
	v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	if is.Empty() {
		log.Println("[DEBUG] Empty InstanceState; nothing to migrate.")
		return is, nil
	}

	switch v {
	case 0:
		log.Println("[INFO] Found Compute Instance State v0; migrating to v1")
		is, err := migrateStateV0toV1(is)
		if err != nil {
			return is, err
		}
		fallthrough
	case 1:
		log.Println("[INFO] Found Compute Instance State v1; migrating to v2")
		is, err := migrateStateV1toV2(is)
		if err != nil {
			return is, err
		}
		return is, nil
	case 2:
		log.Println("[INFO] Found Compute Instance State v2; migrating to v3")
		is, err := migrateStateV2toV3(is)
		if err != nil {
			return is, err
		}
		return is, nil
	case 3:
		log.Println("[INFO] Found Compute Instance State v3; migrating to v4")
		is, err := migrateStateV3toV4(is, meta)
		if err != nil {
			return is, err
		}
		return is, nil
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateStateV0toV1(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	log.Printf("[DEBUG] Attributes before migration: %#v", is.Attributes)

	// Delete old count
	delete(is.Attributes, "metadata.#")

	newMetadata := make(map[string]string)

	for k, v := range is.Attributes {
		if !strings.HasPrefix(k, "metadata.") {
			continue
		}

		// We have a key that looks like "metadata.*" and we know it's not
		// metadata.# because we deleted it above, so it must be metadata.<N>.<key>
		// from the List of Maps. Just need to convert it to a single Map by
		// ditching the '<N>' field.
		kParts := strings.SplitN(k, ".", 3)

		// Sanity check: all three parts should be there and <N> should be a number
		badFormat := false
		if len(kParts) != 3 {
			badFormat = true
		} else if _, err := strconv.Atoi(kParts[1]); err != nil {
			badFormat = true
		}

		if badFormat {
			return is, fmt.Errorf(
				"migration error: found metadata key in unexpected format: %s", k)
		}

		// Rejoin as "metadata.<key>"
		newK := strings.Join([]string{kParts[0], kParts[2]}, ".")
		newMetadata[newK] = v
		delete(is.Attributes, k)
	}

	for k, v := range newMetadata {
		is.Attributes[k] = v
	}

	log.Printf("[DEBUG] Attributes after migration: %#v", is.Attributes)
	return is, nil
}

func migrateStateV1toV2(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	log.Printf("[DEBUG] Attributes before migration: %#v", is.Attributes)

	// Maps service account index to list of scopes for that sccount
	newScopesMap := make(map[string][]string)

	for k, v := range is.Attributes {
		if !strings.HasPrefix(k, "service_account.") {
			continue
		}

		if k == "service_account.#" {
			continue
		}

		if strings.HasSuffix(k, ".scopes.#") {
			continue
		}

		if strings.HasSuffix(k, ".email") {
			continue
		}

		// Key is now of the form service_account.%d.scopes.%d
		kParts := strings.Split(k, ".")

		// Sanity check: all three parts should be there and <N> should be a number
		badFormat := false
		if len(kParts) != 4 {
			badFormat = true
		} else if _, err := strconv.Atoi(kParts[1]); err != nil {
			badFormat = true
		}

		if badFormat {
			return is, fmt.Errorf(
				"migration error: found scope key in unexpected format: %s", k)
		}

		newScopesMap[kParts[1]] = append(newScopesMap[kParts[1]], v)

		delete(is.Attributes, k)
	}

	for service_acct_index, newScopes := range newScopesMap {
		for _, newScope := range newScopes {
			hash := hashcode.String(canonicalizeServiceScope(newScope))
			newKey := fmt.Sprintf("service_account.%s.scopes.%d", service_acct_index, hash)
			is.Attributes[newKey] = newScope
		}
	}

	log.Printf("[DEBUG] Attributes after migration: %#v", is.Attributes)
	return is, nil
}

func migrateStateV2toV3(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	log.Printf("[DEBUG] Attributes before migration: %#v", is.Attributes)
	is.Attributes["create_timeout"] = "4"
	log.Printf("[DEBUG] Attributes after migration: %#v", is.Attributes)
	return is, nil
}

func migrateStateV3toV4(is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	log.Printf("[DEBUG] Attributes before migration: %#v", is.Attributes)

	// Read instance from GCP. Since disks are not necessarily returned from the API in the order they were set,
	// we have no other way to know which source belongs to which attached disk.
	// Also note that the following code modifies the returned instance- if you need immutability, please change
	// this to make a copy of the needed data.
	config := meta.(*Config)
	instance, err := getInstanceFromInstanceState(config, is)
	if err != nil {
		return is, fmt.Errorf("migration error: %s", err)
	}
	diskList, err := getAllDisksFromInstanceState(config, is)
	if err != nil {
		return is, fmt.Errorf("migration error: %s", err)
	}
	allDisks := make(map[string]*compute.Disk)
	for _, disk := range diskList {
		allDisks[disk.Name] = disk
	}

	hasBootDisk := is.Attributes["boot_disk.#"] == "1"

	scratchDisks := 0
	if v := is.Attributes["scratch_disk.#"]; v != "" {
		scratchDisks, err = strconv.Atoi(v)
		if err != nil {
			return is, fmt.Errorf("migration error: found scratch_disk.# value in unexpected format", err)
		}
	}

	attachedDisks := 0
	if v := is.Attributes["attached_disk.#"]; v != "" {
		attachedDisks, err = strconv.Atoi(v)
		if err != nil {
			return is, fmt.Errorf("migration error: found attached_disk.# value in unexpected format", err)
		}
	}

	disks, err := strconv.Atoi(is.Attributes["disk.#"])
	if err != nil {
		return is, fmt.Errorf("migration error: found disk.# value in unexpected format", err)
	}

	for i := 0; i < disks; i++ {
		if !hasBootDisk && i == 0 {
			is.Attributes["boot_disk.#"] = "1"

			// Note: the GCP API does not allow for scratch disks to be boot disks, so this situation
			// should never occur.
			if is.Attributes["disk.0.scratch_disk"] == "true" {
				return is, fmt.Errorf("migration error: found scratch disk at index 0")
			}

			for _, disk := range instance.Disks {
				if disk.Boot {
					sourceUrl := strings.Split(disk.Source, "/")
					is.Attributes["boot_disk.0.source"] = sourceUrl[len(sourceUrl)-1]
					is.Attributes["boot_disk.0.device_name"] = disk.DeviceName
					break
				}
			}
			is.Attributes["boot_disk.0.auto_delete"] = is.Attributes["disk.0.auto_delete"]
			is.Attributes["boot_disk.0.disk_encryption_key_raw"] = is.Attributes["disk.0.disk_encryption_key_raw"]
			is.Attributes["boot_disk.0.disk_encryption_key_sha256"] = is.Attributes["disk.0.disk_encryption_key_sha256"]

			// Don't worry about initialize_params, since the disk has already been created.
		} else if is.Attributes[fmt.Sprintf("disk.%d.scratch", i)] == "true" {
			// Note: the GCP API does not allow for scratch disks without auto_delete, so this situation
			// should never occur.
			if is.Attributes[fmt.Sprintf("disk.%d.auto_delete", i)] != "true" {
				return is, fmt.Errorf("migration error: attempted to migrate scratch disk where auto_delete is not true")
			}

			is.Attributes[fmt.Sprintf("scratch_disk.%d.interface", scratchDisks)] = "SCSI"

			scratchDisks++
		} else {
			// If disk is neither boot nor scratch, then it is attached.

			disk, err := getDiskFromAttributes(config, instance, allDisks, is.Attributes, i)
			if err != nil {
				return is, fmt.Errorf("migration error: %s", err)
			}

			is.Attributes[fmt.Sprintf("attached_disk.%d.source", attachedDisks)] = disk.Source
			is.Attributes[fmt.Sprintf("attached_disk.%d.device_name", attachedDisks)] = disk.DeviceName
			is.Attributes[fmt.Sprintf("attached_disk.%d.disk_encryption_key_raw", attachedDisks)] = is.Attributes[fmt.Sprintf("disk.%d.disk_encryption_key_raw", i)]
			is.Attributes[fmt.Sprintf("attached_disk.%d.disk_encryption_key_sha256", attachedDisks)] = is.Attributes[fmt.Sprintf("disk.%d.disk_encryption_key_sha256", i)]

			attachedDisks++
		}
	}

	for k, _ := range is.Attributes {
		if !strings.HasPrefix(k, "disk.") {
			continue
		}

		delete(is.Attributes, k)
	}
	if scratchDisks > 0 {
		is.Attributes["scratch_disk.#"] = strconv.Itoa(scratchDisks)
	}
	if attachedDisks > 0 {
		is.Attributes["attached_disk.#"] = strconv.Itoa(attachedDisks)
	}

	log.Printf("[DEBUG] Attributes after migration: %#v", is.Attributes)
	return is, nil
}

func getInstanceFromInstanceState(config *Config, is *terraform.InstanceState) (*compute.Instance, error) {
	project, ok := is.Attributes["project"]
	if !ok && config.Project == "" {
		return nil, fmt.Errorf("could not determine 'project'")
	}
	project = config.Project

	zone, ok := is.Attributes["zone"]
	if !ok {
		return nil, fmt.Errorf("could not determine 'zone'")
	}

	instance, err := config.clientCompute.Instances.Get(
		project, zone, is.ID).Do()
	if err != nil {
		return nil, fmt.Errorf("error reading instance: %s", err)
	}

	return instance, nil
}

func getAllDisksFromInstanceState(config *Config, is *terraform.InstanceState) ([]*compute.Disk, error) {
	project, ok := is.Attributes["project"]
	if !ok && config.Project == "" {
		return nil, fmt.Errorf("could not determine 'project'")
	}
	project = config.Project

	zone, ok := is.Attributes["zone"]
	if !ok {
		return nil, fmt.Errorf("could not determine 'zone'")
	}

	diskList := []*compute.Disk{}
	token := ""
	for {
		disks, err := config.clientCompute.Disks.List(project, zone).PageToken(token).Do()
		if err != nil {
			return nil, fmt.Errorf("error reading disks: %s", err)
		}
		diskList = append(diskList, disks.Items...)
		token = disks.NextPageToken
		if token == "" {
			break
		}
	}

	return diskList, nil
}

func getDiskFromAttributes(config *Config, instance *compute.Instance, allDisks map[string]*compute.Disk, attributes map[string]string, i int) (*compute.AttachedDisk, error) {
	if diskSource := attributes[fmt.Sprintf("disk.%d.disk", i)]; diskSource != "" {
		return getDiskFromSource(instance, diskSource)
	}

	if deviceName := attributes[fmt.Sprintf("disk.%d.device_name", i)]; deviceName != "" {
		return getDiskFromDeviceName(instance, deviceName)
	}

	if encryptionKey := attributes[fmt.Sprintf("disk.%d.disk_encryption_key_raw", i)]; encryptionKey != "" {
		return getDiskFromEncryptionKey(instance, encryptionKey)
	}

	autoDelete, err := strconv.ParseBool(attributes[fmt.Sprintf("disk.%d.auto_delete", i)])
	if err != nil {
		return nil, fmt.Errorf("error parsing auto_delete attribute of disk %d", i)
	}
	image := attributes[fmt.Sprintf("disk.%d.image", i)]

	// We know project and zone are set because we used them to read the instance
	project := attributes["project"]
	zone := attributes["zone"]
	return getDiskFromAutoDeleteAndImage(config, instance, allDisks, autoDelete, image, project, zone)
}

func getDiskFromSource(instance *compute.Instance, source string) (*compute.AttachedDisk, error) {
	for _, disk := range instance.Disks {
		if disk.Boot == true || disk.Type == "SCRATCH" {
			// Ignore boot/scratch disks since this is just for finding attached disks
			continue
		}
		// we can just compare suffixes because terraform only allows setting "disk" by name and uses
		// the zone of the instance so we know there can be no duplicate names.
		if strings.HasSuffix(disk.Source, "/"+source) {
			return disk, nil
		}
	}
	return nil, fmt.Errorf("could not find attached disk with source %q", source)
}

func getDiskFromDeviceName(instance *compute.Instance, deviceName string) (*compute.AttachedDisk, error) {
	for _, disk := range instance.Disks {
		if disk.Boot == true || disk.Type == "SCRATCH" {
			// Ignore boot/scratch disks since this is just for finding attached disks
			continue
		}
		if disk.DeviceName == deviceName {
			return disk, nil
		}
	}
	return nil, fmt.Errorf("could not find attached disk with deviceName %q", deviceName)
}

func getDiskFromEncryptionKey(instance *compute.Instance, encryptionKey string) (*compute.AttachedDisk, error) {
	sha := sha256.Sum256([]byte(encryptionKey))
	encryptionSha := base64.StdEncoding.EncodeToString(sha[:])
	for _, disk := range instance.Disks {
		if disk.Boot == true || disk.Type == "SCRATCH" {
			// Ignore boot/scratch disks since this is just for finding attached disks
			continue
		}
		if disk.DiskEncryptionKey.Sha256 == encryptionSha {
			return disk, nil
		}
	}
	return nil, fmt.Errorf("could not find attached disk with encryption hash %q", encryptionSha)
}

func getDiskFromAutoDeleteAndImage(config *Config, instance *compute.Instance, allDisks map[string]*compute.Disk, autoDelete bool, image, project, zone string) (*compute.AttachedDisk, error) {

	log.Printf("[DEBUG] AllDisks is %+v", allDisks)
	for i, disk := range instance.Disks {
		if disk.Boot == true || disk.Type == "SCRATCH" {
			// Ignore boot/scratch disks since this is just for finding attached disks
			continue
		}
		if disk.AutoDelete == autoDelete {
			log.Printf("[DEBUG] Checking disk %+v", disk)
			// Read the disk to check if its image matches
			sourceUrl := strings.Split(disk.Source, "/")
			fullDisk := allDisks[sourceUrl[len(sourceUrl)-1]]
			fullDiskImageUrl := strings.Split(fullDisk.SourceImage, "/")
			fullDiskImage := fullDiskImageUrl[len(fullDiskImageUrl)-1]

			if fullDiskImage == image {
				// Delete this disk because there might be multiple that match
				instance.Disks = append(instance.Disks[:i], instance.Disks[i+1:]...)
				return disk, nil
			}
		}
	}

	// We're not done! It's possible the disk was created with an image family rather than the image itself.
	// Now, do the exact same iteration but do some prefix matching to check if the families match.
	// This assumes that all disks with a given family are of the format family-name-{diskname}, because
	// there's not much else we can do at this point.
	for i, disk := range instance.Disks {
		if disk.Boot == true || disk.Type == "SCRATCH" {
			// Ignore boot/scratch disks since this is just for finding attached disks
			continue
		}
		if disk.AutoDelete == autoDelete {
			// Read the disk to check if its image matches
			sourceUrl := strings.Split(disk.Source, "/")
			fullDisk := allDisks[sourceUrl[len(sourceUrl)-1]]
			fullDiskImageUrl := strings.Split(fullDisk.SourceImage, "/")
			fullDiskImage := fullDiskImageUrl[len(fullDiskImageUrl)-1]

			if strings.HasPrefix(fullDiskImage, image+"-") {
				// Delete this disk because there might be multiple that match
				instance.Disks = append(instance.Disks[:i], instance.Disks[i+1:]...)
				return disk, nil
			}
		}
	}

	return nil, fmt.Errorf("could not find attached disk with image %q", image)
}
