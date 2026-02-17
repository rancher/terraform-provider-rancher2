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

func resourceRancher2ClusterProxyConfigV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2ClusterProxyConfigV2Create,
		Read:   resourceRancher2ClusterProxyConfigV2Read,
		Delete: resourceRancher2ClusterProxyConfigV2Delete,
		Update: resourceRancher2ClusterProxyConfigV2Update,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2ClusterProxyConfigV2Import,
		},
		Schema: clusterProxyConfigV2Fields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2ClusterProxyConfigV2Create(d *schema.ResourceData, meta interface{}) error {
	clusterID := d.Get("cluster_id").(string)

	log.Printf("[INFO] Creating ClusterProxyConfig for cluster %s", clusterID)

	clusterProxyConfigV2, err := expandClusterProxyConfigV2(d)
	if err != nil {
		return err
	}

	newClusterProxyConfigV2, err := createClusterProxyConfigV2(meta.(*Config), clusterID, clusterProxyConfigV2)
	if err != nil {
		return err
	}

	d.SetId(clusterID + "/" + clusterProxyConfigV2Name)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating"},
		Target:     []string{"active"},
		Refresh:    clusterProxyConfigV2StateRefreshFunc(meta.(*Config), clusterID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	res, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("waiting for cluster proxy config (%s) to be active: %s", newClusterProxyConfigV2.ID, waitErr)
	}
	log.Printf("[DEBUG] resourceRancher2ClusterProxyConfigV2Create: %#v", res)

	return resourceRancher2ClusterProxyConfigV2Read(d, meta)
}

func resourceRancher2ClusterProxyConfigV2Read(d *schema.ResourceData, meta interface{}) error {
	clusterID := d.Get("cluster_id").(string)
	clusterProxyConfigV2Id := clusterID + "/" + clusterProxyConfigV2Name

	log.Printf("[INFO] Refreshing ClusterProxyConfig for cluster %s", clusterID)

	return readClusterProxyConfigV2(clusterProxyConfigV2Id, d, meta)
}

func readClusterProxyConfigV2(id string, d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	resp := &ClusterProxyConfigV2{}
	err := config.getObjectV2ByID(rancher2DefaultLocalClusterID, id, clusterProxyConfigV2ApiType, resp)
	if err != nil {
		if IsNotFound(err) || IsNotAccessibleByID(err) {
			log.Printf("[INFO] Cluster V2 %s not found", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	return flattenClusterProxyConfigV2(d, resp)
}

func resourceRancher2ClusterProxyConfigV2Delete(d *schema.ResourceData, meta interface{}) error {
	clusterID := d.Get("cluster_id").(string)
	log.Printf("[INFO] Deleting ClusterProxyConfig for cluster %s", clusterID)

	clusterProxyConfigV2Id := clusterID + "/" + clusterProxyConfigV2Name
	// First, get the current object to have the resource info for deletion
	obj := &ClusterProxyConfigV2{}
	err := meta.(*Config).getObjectV2ByID(rancher2DefaultLocalClusterID, clusterProxyConfigV2Id, clusterProxyConfigV2ApiType, obj)
	if err != nil {
		if IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	err = deleteClusterProxyConfigV2(meta.(*Config), obj)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"removing"},
		Target:     []string{"removed"},
		Refresh:    clusterProxyConfigV2StateRefreshFunc(meta.(*Config), clusterID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	res, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("waiting for cluster proxy config (%s) to be removed: %s", clusterProxyConfigV2Id, waitErr)
	}
	log.Printf("[DEBUG] resourceRancher2ClusterProxyConfigV2Delete: %#v", res)

	d.SetId("")

	return nil
}

// clusterProxyConfigV2StateRefreshFunc returns a resource.StateRefreshFunc, used to watch a ClusterProxyConfig.
func clusterProxyConfigV2StateRefreshFunc(c *Config, clusterID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		clusterProxyConfigV2Id := clusterID + "/" + clusterProxyConfigV2Name
		obj := &ClusterProxyConfigV2{}
		err := c.getObjectV2ByID(rancher2DefaultLocalClusterID, clusterProxyConfigV2Id, clusterProxyConfigV2ApiType, obj)
		if err != nil {
			if IsNotFound(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}

		return obj, "active", nil
	}
}

func createClusterProxyConfigV2(c *Config, clusterID string, obj *ClusterProxyConfigV2) (*ClusterProxyConfigV2, error) {
	if c == nil {
		return nil, fmt.Errorf("creating ClusterProxyConfig V2: Provider config is nil")
	}
	if clusterID == "" {
		return nil, fmt.Errorf("creating ClusterProxyConfig V2: Cluster ID is empty")
	}
	if obj == nil {
		return nil, fmt.Errorf("creating ClusterProxyConfig V2: ClusterProxyConfig V2 is nil")
	}

	resp := &ClusterProxyConfigV2{}

	err := c.createObjectV2(rancher2DefaultLocalClusterID, clusterProxyConfigV2ApiType, obj, resp)
	if err != nil {
		return nil, fmt.Errorf("creating ClusterProxyConfig V2: %s", err)
	}

	return resp, nil
}

func resourceRancher2ClusterProxyConfigV2Update(d *schema.ResourceData, meta interface{}) error {
	clusterID := d.Get("cluster_id").(string)
	clusterProxyConfigV2Id := clusterID + "/" + clusterProxyConfigV2Name

	log.Printf("[INFO] Updating ClusterProxyConfig for cluster %s", clusterID)

	clusterProxyConfigV2, err := expandClusterProxyConfigV2(d)
	if err != nil {
		return err
	}

	_, err = updateClusterProxyConfigV2(meta.(*Config), clusterProxyConfigV2Id, clusterProxyConfigV2)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"updating"},
		Target:     []string{"active"},
		Refresh:    clusterProxyConfigV2StateRefreshFunc(meta.(*Config), clusterID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	res, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("waiting for cluster proxy config (%s) to be updated: %s", clusterProxyConfigV2Id, waitErr)
	}
	log.Printf("[DEBUG] resourceRancher2ClusterProxyConfigV2Update: %#v", res)

	return resourceRancher2ClusterProxyConfigV2Read(d, meta)
}

func updateClusterProxyConfigV2(c *Config, id string, obj *ClusterProxyConfigV2) (*ClusterProxyConfigV2, error) {
	if c == nil {
		return nil, fmt.Errorf("updating ClusterProxyConfig V2: Provider config is nil")
	}

	if id == "" {
		return nil, fmt.Errorf("updating ClusterProxyConfig V2: ID is empty")
	}

	if obj == nil {
		return nil, fmt.Errorf("updating ClusterProxyConfig V2: ClusterProxyConfig V2 is nil")
	}

	resp := &ClusterProxyConfigV2{}
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	existingObj := &ClusterProxyConfigV2{}
	if err := c.getObjectV2ByID(rancher2DefaultLocalClusterID, id, clusterProxyConfigV2ApiType, existingObj); err != nil {
		return nil, err
	}
	obj.ObjectMeta.ResourceVersion = existingObj.ResourceVersion

	for {
		err := c.updateObjectV2(rancher2DefaultLocalClusterID, id, clusterProxyConfigV2ApiType, obj, resp)
		if err == nil {
			return resp, err
		}

		if !IsServerError(err) && !IsUnknownSchemaType(err) && !IsConflict(err) {
			return nil, err
		}

		if IsConflict(err) {
			// Read object again and update ObjectMeta.ResourceVersion before retry
			newObj := &ClusterProxyConfigV2{}
			if err := c.getObjectV2ByID(rancher2DefaultLocalClusterID, id, clusterProxyConfigV2ApiType, newObj); err != nil {
				return nil, err
			}

			obj.ObjectMeta.ResourceVersion = newObj.ObjectMeta.ResourceVersion
		}

		select {
		case <-time.After(rancher2RetriesWait * time.Second):
		case <-ctx.Done():
			return nil, fmt.Errorf("timeout updating ClusterProxyConfig V2 ID %s: %v", id, err)
		}
	}
}

func deleteClusterProxyConfigV2(c *Config, obj *ClusterProxyConfigV2) error {
	if c == nil {
		return fmt.Errorf("deleting ClusterProxyConfig V2: Provider config is nil")
	}
	if obj == nil {
		return fmt.Errorf("deleting ClusterProxyConfig V2: ClusterProxyConfig V2 is nil")
	}

	resource := &types.Resource{
		ID:      obj.ID,
		Type:    clusterProxyConfigV2ApiType,
		Links:   obj.Links,
		Actions: obj.Actions,
	}
	return c.deleteObjectV2(rancher2DefaultLocalClusterID, resource)
}
