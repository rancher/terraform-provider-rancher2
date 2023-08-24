package rancher2

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2Notifier() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRancher2NotifierCreate,
		ReadContext:   resourceRancher2NotifierRead,
		UpdateContext: resourceRancher2NotifierUpdate,
		DeleteContext: resourceRancher2NotifierDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRancher2NotifierImport,
		},

		Schema: notifierFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2NotifierCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	notifier, err := expandNotifier(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating Notifier %s", notifier.Name)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	newNotifier, err := client.Notifier.Create(notifier)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newNotifier.ID)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active"},
		Refresh:    notifierStateRefreshFunc(client, newNotifier.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return diag.Errorf(
			"[ERROR] waiting for notifier (%s) to be created: %s", newNotifier.ID, waitErr)
	}

	return resourceRancher2NotifierRead(ctx, d, meta)
}

func resourceRancher2NotifierRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Refreshing Notifier ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		notifier, err := client.Notifier.ByID(d.Id())
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] Notifier ID %s not found.", d.Id())
				d.SetId("")
				return nil
			}
			return retry.NonRetryableError(err)
		}

		if err = flattenNotifier(d, notifier); err != nil {
			return retry.NonRetryableError(err)
		}

		return nil
	}))
}

func resourceRancher2NotifierUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Updating Notifier ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	notifier, err := client.Notifier.ByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	newNotifier, err := expandNotifier(d)
	if err != nil {
		return diag.FromErr(err)
	}
	newNotifier.Links = notifier.Links
	newNotifier, err = client.Notifier.Replace(newNotifier)
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    notifierStateRefreshFunc(client, newNotifier.ID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return diag.Errorf(
			"[ERROR] waiting for notifier (%s) to be updated: %s", newNotifier.ID, waitErr)
	}

	return resourceRancher2NotifierRead(ctx, d, meta)
}

func resourceRancher2NotifierDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Deleting Notifier ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	notifier, err := client.Notifier.ByID(id)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Notifier ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	err = client.Notifier.Delete(notifier)
	if err != nil {
		return diag.Errorf("Error removing Notifier: %s", err)
	}

	log.Printf("[DEBUG] Waiting for notifier (%s) to be removed", id)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"removing"},
		Target:     []string{"removed"},
		Refresh:    notifierStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return diag.Errorf(
			"[ERROR] waiting for notifier (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

// notifierStateRefreshFunc returns a retry.StateRefreshFunc, used to watch a Rancher Notifier.
func notifierStateRefreshFunc(client *managementClient.Client, notifierID string) retry.StateRefreshFunc {
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
