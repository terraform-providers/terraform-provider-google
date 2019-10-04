package google

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func resourceBigtableInstanceMigrateState(
	v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	if is.Empty() {
		log.Println("[DEBUG] Empty InstanceState; nothing to migrate.")
		return is, nil
	}

	switch v {
	case 0:
		log.Println("[INFO] Found Bigtable Instance State v0; migrating to v1")
		is, err := migrateBigtableInstanceStateV0toV1(is)
		if err != nil {
			return is, err
		}
		return is, nil
	default:
		return nil, fmt.Errorf("Unexpected schema version %d", v)
	}
}

// This state may have been produced by a version of this provider prior to 2.14.0.
// That version changed the schema to store items in the "cluster" field as a list
// instead of a set.  This changed the semantics of how Terraform indexes the items
// when storing them in the state's attributes, in a backwards-incompatible way.
// This migration remedies this by re-indexing the list of clusters by list index
// instead of Terraform-generated hash of the item.  Additionally, it cleans out
// some top-level fields that have been removed from the resource's schema.
func migrateBigtableInstanceStateV0toV1(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	log.Printf("[DEBUG] Attributes before migration: %#v", is.Attributes)

	// Extract hashes used to identify each cluster
	// TODO: Anything we can do about ordering?
	hashes := make(map[string]bool)
	for k := range is.Attributes {
		if strings.HasPrefix(k, "cluster.") && !strings.Contains(k, "#") {
			parts := strings.Split(k, ".")
			hash := parts[1]
			hashes[hash] = true
		}
	}
	// TODO: Only print keys
	log.Printf("[DEBUG] Found clusters for hashes: %#v", hashes)

	// Migrate each cluster's attributes to newly-indexed entries
	newAttributes := make(map[string]string)
	idx := 0
	fields := []string{"cluster_id", "num_nodes", "storage_type", "zone"}
	for hash, _ := range hashes {
		for _, field := range fields {
			oldAttrKey := fmt.Sprintf("cluster.%s.%s", hash, field)
			newAttrKey := fmt.Sprintf("cluster.%d.%s", idx, field)
			if _, exists := is.Attributes[oldAttrKey]; exists {
				log.Printf("[DEBUG] Migrating %s to %s (value: %s)", oldAttrKey, newAttrKey, is.Attributes[oldAttrKey])
				newAttributes[newAttrKey] = is.Attributes[oldAttrKey]
				delete(is.Attributes, oldAttrKey)
			}
		}
		idx++
	}
	for k, v := range newAttributes {
		is.Attributes[k] = v
	}
	// Also remove legacy cluster_id, zone, num_nodes, and storage_type attributes
	// in favor of nested cluster object
	for _, field := range fields {
		if v, exists := is.Attributes[field]; exists {
			log.Printf("[DEBUG] Deleting removed schema field %s (value: %s)", field, v)
			delete(is.Attributes, field)
		}
	}
	log.Printf("[DEBUG] Attributes after migration: %#v", is.Attributes)

	return is, nil
}
