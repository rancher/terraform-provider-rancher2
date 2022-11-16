package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2ClusterAlertGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2ClusterAlertGroupCreate,
		Read:   resourceRancher2ClusterAlertGroupRead,
		Update: resourceRancher2ClusterAlertGroupUpdate,
		Delete: resourceRancher2ClusterAlertGroupDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2ClusterAlertGroupImport,
		},
		Schema: clusterAlertGroupFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2ClusterAlertGroupCreate(d *schema.ResourceData, meta interface{}) error {
	err := resourceRancher2ClusterAlertGroupRecients(d, meta)
	if err != nil {
		return err
	}
	clusterAlertGroup := expandClusterAlertGroup(d)

	log.Printf("[INFO] Creating Cluster Alert Group %s", clusterAlertGroup.Name)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	newClusterAlertGroup, err := client.ClusterAlertGroup.Create(clusterAlertGroup)
	if err != nil {
		return err
	}

	d.SetId(newClusterAlertGroup.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active"},
		Refresh:    clusterAlertGroupStateRefreshFunc(client, newClusterAlertGroup.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for cluster alert group (%s) to be created: %s", newClusterAlertGroup.ID, waitErr)
	}

	return resourceRancher2ClusterAlertGroupRead(d, meta)
}

func resourceRancher2ClusterAlertGroupRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Cluster Alert Group ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		clusterAlertGroup, err := client.ClusterAlertGroup.ByID(d.Id())
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] Cluster Alert Group ID %s not found.", d.Id())
				d.SetId("")
				return nil
			}
		}

		if err = flattenClusterAlertGroup(d, clusterAlertGroup); err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func resourceRancher2ClusterAlertGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Cluster Alert Group ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	clusterAlertGroup, err := client.ClusterAlertGroup.ByID(d.Id())
	if err != nil {
		return err
	}

	if d.HasChange("recipients") {
		err = resourceRancher2ClusterAlertGroupRecients(d, meta)
		if err != nil {
			return err
		}
	}

	update := map[string]interface{}{
		"clusterId":             d.Get("cluster_id").(string),
		"description":           d.Get("description").(string),
		"groupIntervalSeconds":  int64(d.Get("group_interval_seconds").(int)),
		"groupWaitSeconds":      int64(d.Get("group_wait_seconds").(int)),
		"name":                  d.Get("name").(string),
		"recipients":            expandRecipients(d.Get("recipients").([]interface{})),
		"repeatIntervalSeconds": int64(d.Get("repeat_interval_seconds").(int)),
		"annotations":           toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":                toMapString(d.Get("labels").(map[string]interface{})),
	}

	newClusterAlertGroup, err := client.ClusterAlertGroup.Update(clusterAlertGroup, update)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active"},
		Refresh:    clusterAlertGroupStateRefreshFunc(client, newClusterAlertGroup.ID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for cluster alert group (%s) to be updated: %s", newClusterAlertGroup.ID, waitErr)
	}

	return resourceRancher2ClusterAlertGroupRead(d, meta)
}

func resourceRancher2ClusterAlertGroupDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Cluster Alert Group ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	clusterAlertGroup, err := client.ClusterAlertGroup.ByID(id)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Cluster Alert Group ID %s not found.", id)
			d.SetId("")
			return nil
		}
		return err
	}

	err = client.ClusterAlertGroup.Delete(clusterAlertGroup)
	if err != nil {
		return fmt.Errorf("Error removing Cluster Alert Group: %s", err)
	}

	log.Printf("[DEBUG] Waiting for cluster alert group (%s) to be removed", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"removing"},
		Target:     []string{"removed"},
		Refresh:    clusterAlertGroupStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for cluster alert group (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

func resourceRancher2ClusterAlertGroupRecients(d *schema.ResourceData, meta interface{}) error {
	recipients, ok := d.Get("recipients").([]interface{})
	if !ok {
		return fmt.Errorf("[ERROR] Getting Cluster Alert Group Recipients")
	}

	if len(recipients) > 0 {
		log.Printf("[INFO] Getting Cluster Alert Group Recipients")

		for i := range recipients {
			in := recipients[i].(map[string]interface{})

			recipient, err := meta.(*Config).GetRecipientByNotifier(in["notifier_id"].(string))
			if err != nil {
				return err
			}

			in["notifier_type"] = recipient.NotifierType
			if v, ok := in["default_recipient"].(bool); ok && v {
				in["recipient"] = recipient.Recipient
			}

			recipients[i] = in
		}
		d.Set("recipients", recipients)
	}

	return nil
}

// clusterAlertGroupStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher ClusterAlertGroup.
func clusterAlertGroupStateRefreshFunc(client *managementClient.Client, clusterAlertGroupID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := client.ClusterAlertGroup.ByID(clusterAlertGroupID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}

		return obj, obj.State, nil
	}
}
