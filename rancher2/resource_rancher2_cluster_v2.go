package rancher2

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	norman "github.com/rancher/norman/types"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2ClusterV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2ClusterV2Create,
		Read:   resourceRancher2ClusterV2Read,
		Update: resourceRancher2ClusterV2Update,
		Delete: resourceRancher2ClusterV2Delete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2ClusterV2Import,
		},
		Schema: clusterV2Fields(),
		CustomizeDiff: func(d *schema.ResourceDiff, i interface{}) error {
			if d.HasChange("rke_config") {
				oldObj, newObj := d.GetChange("rke_config")
				//return fmt.Errorf("\n%#v\n%#v\n", oldObj, newObj)
				oldInterface, oldOk := oldObj.([]interface{})
				newInterface, newOk := newObj.([]interface{})
				if oldOk && newOk && len(newInterface) > 0 {
					oldConfig := expandClusterV2RKEConfig(oldInterface)
					newConfig := expandClusterV2RKEConfig(newInterface)
					if reflect.DeepEqual(oldConfig, newConfig) {
						d.Clear("rke_config")
					} else {
						d.SetNew("rke_config", flattenClusterV2RKEConfig(newConfig))
					}
				}
			}
			if d.HasChange("local_auth_endpoint") {
				oldObj, newObj := d.GetChange("local_auth_endpoint")
				//return fmt.Errorf("\n%#v\n%#v\n", oldObj, newObj)
				oldInterface, oldOk := oldObj.([]interface{})
				newInterface, newOk := newObj.([]interface{})
				if oldOk && newOk && len(newInterface) > 0 {
					oldConfig := expandClusterV2LocalAuthEndpoint(oldInterface)
					newConfig := expandClusterV2LocalAuthEndpoint(newInterface)
					if reflect.DeepEqual(oldConfig, newConfig) {
						d.Clear("local_auth_endpoint")
					} else {
						d.SetNew("local_auth_endpoint", flattenClusterV2LocalAuthEndpoint(newConfig))
					}
				}
			}
			return nil
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
	}
}

func resourceRancher2ClusterV2Create(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	cluster := expandClusterV2(d)

	log.Printf("[INFO] Creating Cluster V2 %s", name)

	newCluster, err := createClusterV2(meta.(*Config), cluster)
	if err != nil {
		return err
	}
	d.SetId(newCluster.ID)
	newCluster, err = waitForClusterV2State(meta.(*Config), newCluster.ID, clusterV2CreatedCondition, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	// Waiting for cluster v2 active if it has machine pools defined
	if newCluster.Spec.RKEConfig != nil && newCluster.Spec.RKEConfig.MachinePools != nil && len(newCluster.Spec.RKEConfig.MachinePools) > 0 {
		newCluster, err = waitForClusterV2State(meta.(*Config), newCluster.ID, clusterV2ActiveCondition, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return err
		}
	}

	return resourceRancher2ClusterV2Read(d, meta)
}

func resourceRancher2ClusterV2Read(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Cluster V2 %s", d.Id())

	cluster, err := getClusterV2ByID(meta.(*Config), d.Id())
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) || IsNotAccessibleByID(err) {
			log.Printf("[INFO] Cluster V2 %s not found", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}
	d.Set("cluster_v1_id", cluster.Status.ClusterName)
	err = setClusterV2LegacyData(d, meta.(*Config))
	if err != nil {
		return err
	}
	return flattenClusterV2(d, cluster)
}

func resourceRancher2ClusterV2Update(d *schema.ResourceData, meta interface{}) error {
	cluster := expandClusterV2(d)
	log.Printf("[INFO] Updating Cluster V2 %s", d.Id())

	newCluster, err := updateClusterV2(meta.(*Config), d.Id(), cluster)
	if err != nil {
		return err
	}
	// Waiting for cluster v2 active if it has machine pools defined
	if newCluster.Spec.RKEConfig != nil && newCluster.Spec.RKEConfig.MachinePools != nil && len(newCluster.Spec.RKEConfig.MachinePools) > 0 {
		newCluster, err = waitForClusterV2State(meta.(*Config), newCluster.ID, clusterV2ActiveCondition, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return err
		}
	}
	return resourceRancher2ClusterV2Read(d, meta)
}

func resourceRancher2ClusterV2Delete(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	log.Printf("[INFO] Deleting Cluster V2 %s", name)

	cluster, err := getClusterV2ByID(meta.(*Config), d.Id())
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			d.SetId("")
			return nil
		}
	}
	err = deleteClusterV2(meta.(*Config), cluster)
	if err != nil {
		return err
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"removed"},
		Refresh:    clusterV2StateRefreshFunc(meta, cluster.ID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for cluster (%s) to be removed: %s", cluster.ID, waitErr)
	}
	d.SetId("")
	return nil
}

// clusterV2StateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher Cluster v2.
func clusterV2StateRefreshFunc(meta interface{}, objID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := getClusterV2ByID(meta.(*Config), objID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) || IsNotAccessibleByID(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}
		return obj, "active", nil
	}
}

// Rancher2 Cluster V2 API CRUD functions
func createClusterV2(c *Config, obj *ClusterV2) (*ClusterV2, error) {
	if c == nil {
		return nil, fmt.Errorf("Creating cluster V2: Provider config is nil")
	}
	if obj == nil {
		return nil, fmt.Errorf("Creating cluster V2: Cluster V2 is nil")
	}
	resp := &ClusterV2{}
	err := c.createObjectV2(rancher2DefaultLocalClusterID, clusterV2APIType, obj, resp)
	if err != nil {
		return nil, fmt.Errorf("Creating cluster V2: %s", err)
	}
	return resp, nil
}

func deleteClusterV2(c *Config, obj *ClusterV2) error {
	if c == nil {
		return fmt.Errorf("Deleting cluster V2: Provider config is nil")
	}
	if obj == nil {
		return fmt.Errorf("Deleting cluster V2: Cluster V2 is nil")
	}
	resource := &norman.Resource{
		ID:      obj.ID,
		Type:    clusterV2APIType,
		Links:   obj.Links,
		Actions: obj.Actions,
	}
	return c.deleteObjectV2(rancher2DefaultLocalClusterID, resource)
}

func getClusterV2ByID(c *Config, id string) (*ClusterV2, error) {
	if c == nil {
		return nil, fmt.Errorf("Getting cluster V2: Provider config is nil")
	}
	if len(id) == 0 {
		return nil, fmt.Errorf("Getting cluster V2: Cluster V2 ID is empty")
	}
	resp := &ClusterV2{}
	err := c.getObjectV2ByID(rancher2DefaultLocalClusterID, id, clusterV2APIType, resp)
	if err != nil {
		if !IsServerError(err) && !IsNotFound(err) && !IsForbidden(err) {
			return nil, fmt.Errorf("Getting cluster V2: %s", err)
		}
		return nil, err
	}
	return resp, nil
}

func updateClusterV2(c *Config, id string, obj *ClusterV2) (*ClusterV2, error) {
	if c == nil {
		return nil, fmt.Errorf("Updating cluster V2: Provider config is nil")
	}
	if len(id) == 0 {
		return nil, fmt.Errorf("Updating cluster V2: Cluster V2 ID is empty")
	}
	if obj == nil {
		return nil, fmt.Errorf("Updating cluster V2: Cluster V2 is nil")
	}
	resp := &ClusterV2{}
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()
	for {
		err := c.updateObjectV2(rancher2DefaultLocalClusterID, id, clusterV2APIType, obj, resp)
		if err == nil {
			return resp, err
		}
		if !IsServerError(err) && !IsUnknownSchemaType(err) && !IsConflict(err) {
			return nil, err
		}
		if IsConflict(err) {
			// Read cluster again and update ObjectMeta.ResourceVersion before retry
			newObj := &ClusterV2{}
			err = c.getObjectV2ByID(rancher2DefaultLocalClusterID, id, clusterV2APIType, newObj)
			if err != nil {
				return nil, err
			}
			obj.ObjectMeta.ResourceVersion = newObj.ObjectMeta.ResourceVersion
		}
		select {
		case <-time.After(rancher2RetriesWait * time.Second):
		case <-ctx.Done():
			return nil, fmt.Errorf("Timeout updating cluster V2 ID %s: %v", id, err)
		}
	}
}

func waitForClusterV2State(c *Config, id, state string, interval time.Duration) (*ClusterV2, error) {
	if id == "" || state == "" {
		return nil, fmt.Errorf("Cluster V2 ID and/or condition is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), interval)
	defer cancel()
	for {
		obj, err := getClusterV2ByID(c, id)
		if err != nil {
			log.Printf("[DEBUG] Retrying on error Refreshing Cluster V2 %s: %v", id, err)
			if !IsNotFound(err) && !IsForbidden(err) && !IsNotAccessibleByID(err) {
				return nil, fmt.Errorf("Getting cluster V2 ID (%s): %v", id, err)
			}
			if IsNotAccessibleByID(err) {
				// Restarting clients to update RBAC
				c.RestartClients()
			}
		}
		if obj != nil {
			for i := range obj.Status.Conditions {
				if obj.Status.Conditions[i].Type == state {
					// Status of the condition, one of True, False, Unknown.
					if obj.Status.Conditions[i].Status == "Unknown" {
						break
					}
					if obj.Status.Conditions[i].Status == "True" {
						return obj, nil
					}
					// When cluster condition is false, retrying if it has been updated for last rancher2WaitFalseCond seconds
					lastUpdate, err := time.Parse(time.RFC3339, obj.Status.Conditions[i].LastUpdateTime)
					if err == nil && time.Since(lastUpdate) < rancher2WaitFalseCond*time.Second {
						break
					}
					return nil, fmt.Errorf("Cluster V2 ID %s: %s", id, obj.Status.Conditions[i].Message)
				}
			}
		}
		select {
		case <-time.After(rancher2RetriesWait * time.Second):
		case <-ctx.Done():
			return nil, fmt.Errorf("Timeout waiting for cluster V2 ID %s", id)
		}
	}
}

func setClusterV2LegacyData(d *schema.ResourceData, c *Config) error {
	if c == nil {
		return fmt.Errorf("Setting cluster V2 legacy data: Provider config is nil")
	}
	clusterV1ID := d.Get("cluster_v1_id").(string)
	if len(clusterV1ID) == 0 {
		return fmt.Errorf("Setting cluster V2 legacy data: cluster_v1_id is empty")
	}

	client, err := c.ManagementClient()
	if err != nil {
		return fmt.Errorf("Setting cluster V2 legacy data: %v", err)
	}

	cluster := &Cluster{}
	err = client.APIBaseClient.ByID(managementClient.ClusterType, clusterV1ID, cluster)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Cluster ID %s not found.", cluster.ID)
			return nil
		}
		return fmt.Errorf("Setting cluster V2 legacy data: %v", err)
	}

	clusterRegistrationToken, err := findClusterRegistrationToken(client, cluster.ID)
	if err != nil && !IsForbidden(err) {
		return fmt.Errorf("Setting cluster V2 legacy data: %v", err)
	}
	regToken, _ := flattenClusterRegistationToken(clusterRegistrationToken)
	err = d.Set("cluster_registration_token", regToken)
	if err != nil {
		return fmt.Errorf("Setting cluster V2 legacy data: %v", err)
	}

	kubeConfig, err := getClusterKubeconfig(c, cluster.ID, d.Get("kube_config").(string))
	if err != nil {
		return fmt.Errorf("Setting cluster V2 legacy data: %v", err)
	}
	d.Set("kube_config", kubeConfig.Config)

	return nil
}
