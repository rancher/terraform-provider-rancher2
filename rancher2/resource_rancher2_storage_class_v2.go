package rancher2

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/rancher/norman/types"
)

func resourceRancher2StorageClassV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2StorageClassV2Create,
		Read:   resourceRancher2StorageClassV2Read,
		Update: resourceRancher2StorageClassV2Update,
		Delete: resourceRancher2StorageClassV2Delete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2StorageClassV2Import,
		},
		Schema: storageClassV2Fields(),
		CustomizeDiff: func(d *schema.ResourceDiff, i interface{}) error {
			if d.HasChange("mount_options") {
				old, new := d.GetChange("mount_options")
				oldObj := toArrayStringSorted(old.([]interface{}))
				newObj := toArrayStringSorted(new.([]interface{}))
				if reflect.DeepEqual(oldObj, newObj) {
					d.Clear("mount_options")
				} else {
					err := d.SetNew("mount_options", toArrayInterfaceSorted(newObj))
					if err != nil {
						return err
					}
				}
			}
			return nil
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2StorageClassV2Create(d *schema.ResourceData, meta interface{}) error {
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	storageClass := expandStorageClassV2(d)

	log.Printf("[INFO] Creating StorageClass V2 %s", name)

	newStorageClass, err := createStorageClassV2(meta.(*Config), clusterID, storageClass)
	if err != nil {
		return err
	}
	d.SetId(clusterID + storageClassV2ClusterIDsep + newStorageClass.ID)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active"},
		Refresh:    storageClassV2StateRefreshFunc(meta, clusterID, newStorageClass.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for storageClass (%s) to be active: %s", newStorageClass.ID, waitErr)
	}
	return resourceRancher2StorageClassV2Read(d, meta)
}

func resourceRancher2StorageClassV2Read(d *schema.ResourceData, meta interface{}) error {
	clusterID, rancherID := splitID(d.Id())
	log.Printf("[INFO] Refreshing StorageClass V2 %s at Cluster ID %s", rancherID, clusterID)

	return resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		storageClass, err := getStorageClassV2ByID(meta.(*Config), clusterID, rancherID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] StorageClass V2 %s not found at cluster ID %s", rancherID, clusterID)
				d.SetId("")
				return nil
			}
			return resource.NonRetryableError(err)
		}
		if err = flattenStorageClassV2(d, storageClass); err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func resourceRancher2StorageClassV2Update(d *schema.ResourceData, meta interface{}) error {
	clusterID, rancherID := splitID(d.Id())
	storageClass := expandStorageClassV2(d)
	log.Printf("[INFO] Updating StorageClass V2 %s at Cluster ID %s", rancherID, clusterID)

	newStorageClass, err := updateStorageClassV2(meta.(*Config), clusterID, rancherID, storageClass)
	if err != nil {
		return err
	}
	d.SetId(clusterID + storageClassV2ClusterIDsep + newStorageClass.ID)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active"},
		Refresh:    storageClassV2StateRefreshFunc(meta, clusterID, newStorageClass.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for storageClass (%s) to be active: %s", newStorageClass.ID, waitErr)
	}
	return resourceRancher2StorageClassV2Read(d, meta)
}

func resourceRancher2StorageClassV2Delete(d *schema.ResourceData, meta interface{}) error {
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	log.Printf("[INFO] Deleting StorageClass V2 %s", name)

	_, rancherID := splitID(d.Id())
	storageClass, err := getStorageClassV2ByID(meta.(*Config), clusterID, rancherID)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			d.SetId("")
			return nil
		}
	}
	err = deleteStorageClassV2(meta.(*Config), clusterID, storageClass)
	if err != nil {
		return err
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"removed"},
		Refresh:    storageClassV2StateRefreshFunc(meta, clusterID, storageClass.ID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for storageClass (%s) to be active: %s", storageClass.ID, waitErr)
	}
	d.SetId("")
	return nil
}

// storageClassV2StateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher StorageClass v2.
func storageClassV2StateRefreshFunc(meta interface{}, clusterID, storageClassID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := getStorageClassV2ByID(meta.(*Config), clusterID, storageClassID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}
		return obj, "active", nil
	}
}

// Rancher2 StorageClass V2 API CRUD functions
func createStorageClassV2(c *Config, clusterID string, obj *StorageClassV2) (*StorageClassV2, error) {
	if c == nil {
		return nil, fmt.Errorf("Creating storageClass V2: Provider config is nil")
	}
	if len(clusterID) == 0 {
		return nil, fmt.Errorf("Creating storageClass V2: Cluster ID is empty")
	}
	if obj == nil {
		return nil, fmt.Errorf("Creating storageClass V2: StorageClass V2 is nil")
	}
	resp := &StorageClassV2{}
	err := c.createObjectV2(clusterID, storageClassV2APIType, obj, resp)
	if err != nil {
		return nil, fmt.Errorf("Creating storageClass V2: %s", err)
	}
	return resp, nil
}

func deleteStorageClassV2(c *Config, clusterID string, obj *StorageClassV2) error {
	if c == nil {
		return fmt.Errorf("Deleting storageClass V2: Provider config is nil")
	}
	if len(clusterID) == 0 {
		return fmt.Errorf("Deleting storageClass V2: Cluster ID is empty")
	}
	if obj == nil {
		return fmt.Errorf("Deleting storageClass V2: StorageClass V2 is nil")
	}
	resource := &types.Resource{
		ID:      obj.ID,
		Type:    storageClassV2APIType,
		Links:   obj.Links,
		Actions: obj.Actions,
	}
	return c.deleteObjectV2(clusterID, resource)
}

func getStorageClassV2ByID(c *Config, clusterID, id string) (*StorageClassV2, error) {
	if c == nil {
		return nil, fmt.Errorf("Getting storageClass V2: Provider config is nil")
	}
	if len(clusterID) == 0 {
		return nil, fmt.Errorf("Getting storageClass V2: Cluster ID is empty")
	}
	if len(id) == 0 {
		return nil, fmt.Errorf("Getting storageClass V2: StorageClass V2 ID is empty")
	}
	resp := &StorageClassV2{}
	err := c.getObjectV2ByID(clusterID, id, storageClassV2APIType, resp)
	if err != nil {
		if !IsServerError(err) && !IsNotFound(err) && !IsForbidden(err) {
			return nil, fmt.Errorf("Getting storageClass V2: %s", err)
		}
		return nil, err
	}
	return resp, nil
}

func updateStorageClassV2(c *Config, clusterID, id string, obj *StorageClassV2) (*StorageClassV2, error) {
	if c == nil {
		return nil, fmt.Errorf("Updating storageClass V2: Provider config is nil")
	}
	if len(clusterID) == 0 {
		return nil, fmt.Errorf("Updating storageClass V2: Cluster ID is empty")
	}
	if len(id) == 0 {
		return nil, fmt.Errorf("Updating storageClass V2: StorageClass V2 ID is empty")
	}
	if obj == nil {
		return nil, fmt.Errorf("Updating storageClass V2: StorageClass V2 is nil")
	}
	resp := &StorageClassV2{}
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()
	for {
		err := c.updateObjectV2(clusterID, id, storageClassV2APIType, obj, resp)
		if err == nil {
			return resp, err
		}
		if !IsServerError(err) && !IsUnknownSchemaType(err) && !IsConflict(err) {
			return nil, err
		}
		if IsConflict(err) {
			// Read storageClass again and update ObjectMeta.ResourceVersion before retry
			newObj := &StorageClassV2{}
			err = c.getObjectV2ByID(clusterID, id, storageClassV2APIType, newObj)
			if err != nil {
				return nil, err
			}
			obj.ObjectMeta.ResourceVersion = newObj.ObjectMeta.ResourceVersion
		}
		select {
		case <-time.After(rancher2RetriesWait * time.Second):
		case <-ctx.Done():
			return nil, fmt.Errorf("Timeout updating storageClass V2 ID %s: %v", id, err)
		}
	}
}
