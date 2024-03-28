package rancher2

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/rancher/norman/types"
)

func resourceRancher2ConfigMapV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2ConfigMapV2Create,
		Read:   resourceRancher2ConfigMapV2Read,
		Update: resourceRancher2ConfigMapV2Update,
		Delete: resourceRancher2ConfigMapV2Delete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2ConfigMapV2Import,
		},
		Schema: configMapV2Fields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2ConfigMapV2Create(d *schema.ResourceData, meta interface{}) error {
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	configMap := expandConfigMapV2(d)

	log.Printf("[INFO] Creating ConfigMap V2 %s", name)

	newConfigMap, err := createConfigMapV2(meta.(*Config), clusterID, configMap)
	if err != nil {
		return err
	}
	d.SetId(clusterID + configMapV2ClusterIDsep + newConfigMap.ID)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active"},
		Refresh:    configMapV2StateRefreshFunc(meta, clusterID, newConfigMap.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for configMap (%s) to be active: %s", newConfigMap.ID, waitErr)
	}
	return resourceRancher2ConfigMapV2Read(d, meta)
}

func resourceRancher2ConfigMapV2Read(d *schema.ResourceData, meta interface{}) error {
	clusterID, rancherID := splitID(d.Id())
	log.Printf("[INFO] Refreshing ConfigMap V2 %s at Cluster ID %s", rancherID, clusterID)

	return resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		configMap, err := getConfigMapV2ByID(meta.(*Config), clusterID, rancherID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] ConfigMap V2 %s not found at cluster ID %s", rancherID, clusterID)
				d.SetId("")
				return nil
			}
			return resource.NonRetryableError(err)
		}
		if err = flattenConfigMapV2(d, configMap); err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func resourceRancher2ConfigMapV2Update(d *schema.ResourceData, meta interface{}) error {
	clusterID, rancherID := splitID(d.Id())
	configMap := expandConfigMapV2(d)
	log.Printf("[INFO] Updating ConfigMap V2 %s at Cluster ID %s", rancherID, clusterID)

	newConfigMap, err := updateConfigMapV2(meta.(*Config), clusterID, rancherID, configMap)
	if err != nil {
		return err
	}
	d.SetId(clusterID + configMapV2ClusterIDsep + newConfigMap.ID)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active"},
		Refresh:    configMapV2StateRefreshFunc(meta, clusterID, newConfigMap.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for configMap (%s) to be active: %s", newConfigMap.ID, waitErr)
	}
	return resourceRancher2ConfigMapV2Read(d, meta)
}

func resourceRancher2ConfigMapV2Delete(d *schema.ResourceData, meta interface{}) error {
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	log.Printf("[INFO] Deleting ConfigMap V2 %s", name)

	_, rancherID := splitID(d.Id())
	configMap, err := getConfigMapV2ByID(meta.(*Config), clusterID, rancherID)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			d.SetId("")
			return nil
		}
	}
	err = deleteConfigMapV2(meta.(*Config), clusterID, configMap)
	if err != nil {
		return err
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"removed"},
		Refresh:    configMapV2StateRefreshFunc(meta, clusterID, configMap.ID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for configMap (%s) to be removed: %s", configMap.ID, waitErr)
	}
	d.SetId("")
	return nil
}

// configMapV2StateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher ConfigMap v2.
func configMapV2StateRefreshFunc(meta interface{}, clusterID, configMapID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := getConfigMapV2ByID(meta.(*Config), clusterID, configMapID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}
		return obj, "active", nil
	}
}

// Rancher2 ConfigMap V2 API CRUD functions
func createConfigMapV2(c *Config, clusterID string, obj *ConfigMapV2) (*ConfigMapV2, error) {
	if c == nil {
		return nil, fmt.Errorf("Creating configMap V2: Provider config is nil")
	}
	if len(clusterID) == 0 {
		return nil, fmt.Errorf("Creating configMap V2: Cluster ID is empty")
	}
	if obj == nil {
		return nil, fmt.Errorf("Creating configMap V2: ConfigMap V2 is nil")
	}
	resp := &ConfigMapV2{}
	err := c.createObjectV2(clusterID, configMapV2APIType, obj, resp)
	if err != nil {
		return nil, fmt.Errorf("Creating configMap V2: %s", err)
	}
	return resp, nil
}

func deleteConfigMapV2(c *Config, clusterID string, obj *ConfigMapV2) error {
	if c == nil {
		return fmt.Errorf("Deleting configMap V2: Provider config is nil")
	}
	if len(clusterID) == 0 {
		return fmt.Errorf("Deleting configMap V2: Cluster ID is empty")
	}
	if obj == nil {
		return fmt.Errorf("Deleting configMap V2: ConfigMap V2 is nil")
	}
	resource := &types.Resource{
		ID:      obj.ID,
		Type:    configMapV2APIType,
		Links:   obj.Links,
		Actions: obj.Actions,
	}
	return c.deleteObjectV2(clusterID, resource)
}

func getConfigMapV2ByID(c *Config, clusterID, id string) (*ConfigMapV2, error) {
	if c == nil {
		return nil, fmt.Errorf("Getting configMap V2: Provider config is nil")
	}
	if len(clusterID) == 0 {
		return nil, fmt.Errorf("Getting configMap V2: Cluster ID is empty")
	}
	if len(id) == 0 {
		return nil, fmt.Errorf("Getting configMap V2: ConfigMap V2 ID is empty")
	}
	resp := &ConfigMapV2{}
	err := c.getObjectV2ByID(clusterID, id, configMapV2APIType, resp)
	if err != nil {
		if !IsServerError(err) && !IsNotFound(err) && !IsForbidden(err) {
			return nil, fmt.Errorf("Getting configMap V2: %s", err)
		}
		return nil, err
	}
	return resp, nil
}

func updateConfigMapV2(c *Config, clusterID, id string, obj *ConfigMapV2) (*ConfigMapV2, error) {
	if c == nil {
		return nil, fmt.Errorf("Updating configMap V2: Provider config is nil")
	}
	if len(clusterID) == 0 {
		return nil, fmt.Errorf("Updating configMap V2: Cluster ID is empty")
	}
	if len(id) == 0 {
		return nil, fmt.Errorf("Updating configMap V2: ConfigMap V2 ID is empty")
	}
	if obj == nil {
		return nil, fmt.Errorf("Updating configMap V2: ConfigMap V2 is nil")
	}
	resp := &ConfigMapV2{}
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()
	for {
		err := c.updateObjectV2(clusterID, id, configMapV2APIType, obj, resp)
		if err == nil {
			return resp, err
		}
		if !IsServerError(err) && !IsUnknownSchemaType(err) && !IsConflict(err) {
			return nil, err
		}
		if IsConflict(err) {
			// Read clusterRepo again and update ObjectMeta.ResourceVersion before retry
			newObj := &ConfigMapV2{}
			err = c.getObjectV2ByID(clusterID, id, configMapV2APIType, newObj)
			if err != nil {
				return nil, err
			}
			obj.ObjectMeta.ResourceVersion = newObj.ObjectMeta.ResourceVersion
		}
		select {
		case <-time.After(rancher2RetriesWait * time.Second):
		case <-ctx.Done():
			return nil, fmt.Errorf("Timeout updating ConfigMap V2 ID %s: %v", id, err)
		}
	}
}
