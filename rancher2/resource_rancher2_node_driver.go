package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

//Schemas

func nodeDriverFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"active": &schema.Schema{
			Type:     schema.TypeBool,
			Required: true,
		},

		"builtin": &schema.Schema{
			Type:     schema.TypeBool,
			Required: true,
		},
		"checksum": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"external_id": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"ui_url": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"url": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"whitelist_domains": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
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

func flattenNodeDriver(d *schema.ResourceData, in *managementClient.NodeDriver) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)

	err := d.Set("active", in.Active)
	if err != nil {
		return err
	}

	err = d.Set("builtin", in.Builtin)
	if err != nil {
		return err
	}

	err = d.Set("checksum", in.Checksum)
	if err != nil {
		return err
	}

	err = d.Set("description", in.Description)
	if err != nil {
		return err
	}

	err = d.Set("name", in.Name)
	if err != nil {
		return err
	}

	err = d.Set("external_id", in.ExternalID)
	if err != nil {
		return err
	}

	err = d.Set("ui_url", in.UIURL)
	if err != nil {
		return err
	}

	err = d.Set("url", in.URL)
	if err != nil {
		return err
	}

	err = d.Set("whitelist_domains", toArrayInterface(in.WhitelistDomains))
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

func expandNodeDriver(in *schema.ResourceData) *managementClient.NodeDriver {
	obj := &managementClient.NodeDriver{}
	if in == nil {
		return nil
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	obj.Active = in.Get("active").(bool)
	obj.Builtin = in.Get("builtin").(bool)
	obj.Checksum = in.Get("checksum").(string)
	obj.Description = in.Get("description").(string)
	obj.ExternalID = in.Get("external_id").(string)
	obj.Name = in.Get("name").(string)
	obj.UIURL = in.Get("ui_url").(string)
	obj.URL = in.Get("url").(string)

	if v, ok := in.Get("whitelist_domains").([]interface{}); ok && len(v) > 0 {
		obj.WhitelistDomains = toArrayString(v)
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj
}

func resourceRancher2NodeDriver() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2NodeDriverCreate,
		Read:   resourceRancher2NodeDriverRead,
		Update: resourceRancher2NodeDriverUpdate,
		Delete: resourceRancher2NodeDriverDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2NodeDriverImport,
		},
		Schema: nodeDriverFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2NodeDriverCreate(d *schema.ResourceData, meta interface{}) error {
	nodeDriver := expandNodeDriver(d)

	log.Printf("[INFO] Creating Node Driver %s", nodeDriver.Name)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	newNodeDriver, err := client.NodeDriver.Create(nodeDriver)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"downloading", "activating"},
		Target:     []string{"active", "inactive"},
		Refresh:    nodeDriverStateRefreshFunc(client, newNodeDriver.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for node driver (%s) to be created: %s", newNodeDriver.ID, waitErr)
	}

	d.SetId(newNodeDriver.ID)

	return resourceRancher2NodeDriverRead(d, meta)
}

func resourceRancher2NodeDriverRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Node Driver ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	nodeDriver, err := client.NodeDriver.ByID(d.Id())
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Node Driver ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
	}

	return flattenNodeDriver(d, nodeDriver)
}

func resourceRancher2NodeDriverUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Node Driver ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	nodeDriver, err := client.NodeDriver.ByID(d.Id())
	if err != nil {
		return err
	}

	update := map[string]interface{}{
		"active":           d.Get("active").(bool),
		"builtin":          d.Get("builtin").(bool),
		"checksum":         d.Get("checksum").(string),
		"description":      d.Get("description").(string),
		"externalId":       d.Get("external_id").(string),
		"name":             d.Get("name").(string),
		"uiUrl":            d.Get("ui_url").(string),
		"url":              d.Get("url").(string),
		"whitelistDomains": toArrayString(d.Get("whitelist_domains").([]interface{})),
		"annotations":      toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":           toMapString(d.Get("labels").(map[string]interface{})),
	}

	newNodeDriver, err := client.NodeDriver.Update(nodeDriver, update)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active", "inactive", "downloading", "activating", "deactivating"},
		Target:     []string{"active", "inactive"},
		Refresh:    nodeDriverStateRefreshFunc(client, newNodeDriver.ID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for node driver (%s) to be updated: %s", newNodeDriver.ID, waitErr)
	}

	return resourceRancher2NodeDriverRead(d, meta)
}

func resourceRancher2NodeDriverDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Node Driver ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	nodeDriver, err := client.NodeDriver.ByID(id)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Node Driver ID %s not found.", id)
			d.SetId("")
			return nil
		}
		return err
	}

	err = client.NodeDriver.Delete(nodeDriver)
	if err != nil {
		return fmt.Errorf("Error removing Node Driver: %s", err)
	}

	log.Printf("[DEBUG] Waiting for node driver (%s) to be removed", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"removing"},
		Target:     []string{"removed"},
		Refresh:    nodeDriverStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for node driver (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

// nodeDriverStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher NodeDriver.
func nodeDriverStateRefreshFunc(client *managementClient.Client, nodeDriverID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := client.NodeDriver.ByID(nodeDriverID)
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
