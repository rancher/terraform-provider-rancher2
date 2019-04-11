package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

// Schema

func nodePoolFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cluster_id": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"hostname_prefix": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"node_template_id": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"quantity": &schema.Schema{
			Type:     schema.TypeInt,
			Required: true,
		},
		"control_plane": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		"etcd": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		"worker": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		"annotations": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"labels": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
	}

	return s
}

// Flatteners

func flattenNodePool(d *schema.ResourceData, in *managementClient.NodePool) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)

	err := d.Set("cluster_id", in.ClusterID)
	if err != nil {
		return err
	}

	err = d.Set("name", in.Name)
	if err != nil {
		return err
	}

	err = d.Set("hostname_prefix", in.HostnamePrefix)
	if err != nil {
		return err
	}

	err = d.Set("node_template_id", in.NodeTemplateID)
	if err != nil {
		return err
	}

	err = d.Set("quantity", int(in.Quantity))
	if err != nil {
		return err
	}

	err = d.Set("control_plane", in.ControlPlane)
	if err != nil {
		return err
	}

	err = d.Set("etcd", in.Etcd)
	if err != nil {
		return err
	}

	err = d.Set("worker", in.Worker)
	if err != nil {
		return err
	}

	err = d.Set("annotations", toMapInterface(in.Annotations))
	if err != nil {
		return err
	}

	err = d.Set("labels", toMapInterface(in.Labels))
	if err != nil {
		return err
	}
	return nil
}

// Expanders

func expandNodePool(in *schema.ResourceData) *managementClient.NodePool {
	obj := &managementClient.NodePool{}
	if in == nil {
		return nil
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	obj.ClusterID = in.Get("cluster_id").(string)
	obj.Name = in.Get("name").(string)
	obj.HostnamePrefix = in.Get("hostname_prefix").(string)
	obj.NodeTemplateID = in.Get("node_template_id").(string)
	obj.Quantity = int64(in.Get("quantity").(int))
	obj.ControlPlane = in.Get("control_plane").(bool)
	obj.Etcd = in.Get("etcd").(bool)
	obj.Worker = in.Get("worker").(bool)

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj
}

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

	d.SetId(newNodePool.ID)

	return resourceRancher2NodePoolRead(d, meta)
}

func resourceRancher2NodePoolRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Node Pool ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	nodePool, err := client.NodePool.ByID(d.Id())
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Node Pool ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = flattenNodePool(d, nodePool)
	if err != nil {
		return err
	}

	return nil
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
		"clusterId":      d.Get("cluster_id").(string),
		"hostnamePrefix": d.Get("hostname_prefix").(string),
		"nodeTemplateId": d.Get("node_template_id").(string),
		"quantity":       int64(d.Get("quantity").(int)),
		"controlPlane":   d.Get("control_plane").(bool),
		"etcd":           d.Get("etcd").(bool),
		"worker":         d.Get("worker").(bool),
		"annotations":    toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":         toMapString(d.Get("labels").(map[string]interface{})),
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
		if IsNotFound(err) {
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
			if IsNotFound(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}

		if obj.Removed != "" {
			return obj, "removed", nil
		}

		return obj, obj.State, nil
	}
}
