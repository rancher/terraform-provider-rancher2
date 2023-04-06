package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2Notifier() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2NotifierCreate,
		Read:   resourceRancher2NotifierRead,
		Update: resourceRancher2NotifierUpdate,
		Delete: resourceRancher2NotifierDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2NotifierImport,
		},

		Schema: notifierFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2NotifierCreate(d *schema.ResourceData, meta interface{}) error {
	notifier, err := expandNotifier(d)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating Notifier %s", notifier.Name)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	newNotifier, err := client.Notifier.Create(notifier)
	if err != nil {
		return err
	}

	d.SetId(newNotifier.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active"},
		Refresh:    notifierStateRefreshFunc(client, newNotifier.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for notifier (%s) to be created: %s", newNotifier.ID, waitErr)
	}

	return resourceRancher2NotifierRead(d, meta)
}

func resourceRancher2NotifierRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Notifier ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		notifier, err := client.Notifier.ByID(d.Id())
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] Notifier ID %s not found.", d.Id())
				d.SetId("")
				return nil
			}
			return resource.NonRetryableError(err)
		}

		if err = flattenNotifier(d, notifier); err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func resourceRancher2NotifierUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Notifier ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	notifier, err := client.Notifier.ByID(d.Id())
	if err != nil {
		return err
	}

	newNotifier, err := expandNotifier(d)
	if err != nil {
		return err
	}
	newNotifier.Links = notifier.Links
	newNotifier, err = client.Notifier.Replace(newNotifier)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    notifierStateRefreshFunc(client, newNotifier.ID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for notifier (%s) to be updated: %s", newNotifier.ID, waitErr)
	}

	return resourceRancher2NotifierRead(d, meta)
}

func resourceRancher2NotifierDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Notifier ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	notifier, err := client.Notifier.ByID(id)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Notifier ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = client.Notifier.Delete(notifier)
	if err != nil {
		return fmt.Errorf("Error removing Notifier: %s", err)
	}

	log.Printf("[DEBUG] Waiting for notifier (%s) to be removed", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"removing"},
		Target:     []string{"removed"},
		Refresh:    notifierStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for notifier (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

// notifierStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher Notifier.
func notifierStateRefreshFunc(client *managementClient.Client, notifierID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := client.Notifier.ByID(notifierID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}

		return obj, obj.State, nil
	}
}
