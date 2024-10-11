package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2NodePool() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2NodePoolCreate,
		Read:   resourceRancher2NodePoolRead,
		Update: resourceRancher2NodePoolUpdate,
		Delete: resourceRancher2NodePoolDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2NodePoolImport,
		},

		Schema: nodePoolFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2NodePoolCreate(d *schema.ResourceData, meta interface{}) error {
	nodePool := expandNodePool(d)

	log.Printf("[INFO] Creating Node Pool %s", nodePool.Name)

	err := meta.(*Config).ClusterExist(nodePool.ClusterID)
	if err != nil {
		return err
	}

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	newNodePool, err := client.NodePool.Create(nodePool)
	if err != nil {
		return err
	}

	d.SetId(newNodePool.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active"},
		Refresh:    nodePoolStateRefreshFunc(client, newNodePool.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for node pool (%s) to be created: %s", newNodePool.ID, waitErr)
	}

	return resourceRancher2NodePoolRead(d, meta)
}

func resourceRancher2NodePoolRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Node Pool ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		nodePool, err := client.NodePool.ByID(d.Id())
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] Node Pool ID %s not found.", d.Id())
				d.SetId("")
				return nil
			}
			return resource.NonRetryableError(err)
		}

		if err = flattenNodePool(d, nodePool); err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func resourceRancher2NodePoolUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Node Pool ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	nodePool, err := client.NodePool.ByID(d.Id())
	if err != nil {
		return err
	}

	update := map[string]interface{}{
		"hostnamePrefix":          d.Get("hostname_prefix").(string),
		"deleteNotReadyAfterSecs": int64(d.Get("delete_not_ready_after_secs").(int)),
		"drainBeforeDelete":       d.Get("drain_before_delete").(bool),
		"nodeTemplateId":          d.Get("node_template_id").(string),
		"nodeTaints":              expandTaints(d.Get("node_taints").([]interface{})),
		"quantity":                int64(d.Get("quantity").(int)),
		"controlPlane":            d.Get("control_plane").(bool),
		"etcd":                    d.Get("etcd").(bool),
		"worker":                  d.Get("worker").(bool),
		"annotations":             toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":                  toMapString(d.Get("labels").(map[string]interface{})),
	}

	newNodePool, err := client.NodePool.Update(nodePool, update)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    nodePoolStateRefreshFunc(client, newNodePool.ID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for node pool (%s) to be updated: %s", newNodePool.ID, waitErr)
	}

	return resourceRancher2NodePoolRead(d, meta)
}

func resourceRancher2NodePoolDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Node Pool ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	nodePool, err := client.NodePool.ByID(id)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Node Pool ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = client.NodePool.Delete(nodePool)
	if err != nil {
		return fmt.Errorf("Error removing Node Pool: %s", err)
	}

	log.Printf("[DEBUG] Waiting for node pool (%s) to be removed", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"removing"},
		Target:     []string{"removed"},
		Refresh:    nodePoolStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for node pool (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

// nodePoolStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher NodePool.
func nodePoolStateRefreshFunc(client *managementClient.Client, nodePoolID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := client.NodePool.ByID(nodePoolID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}

		return obj, obj.State, nil
	}
}
