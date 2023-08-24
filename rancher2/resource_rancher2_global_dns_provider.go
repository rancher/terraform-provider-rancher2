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

func resourceRancher2GlobalDNSProvider() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRancher2GlobalDNSProviderCreate,
		ReadContext:   resourceRancher2GlobalDNSProviderRead,
		UpdateContext: resourceRancher2GlobalDNSProviderUpdate,
		DeleteContext: resourceRancher2GlobalDNSProviderDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRancher2GlobalDNSProviderImport,
		},

		Schema: globalDNSProviderFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceRancher2GlobalDNSProviderCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	globalDNSProvider := expandGlobalDNSProvider(d)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating Global DNS Provider %s", globalDNSProvider.Name)

	newGlobalDNSProvider, err := client.GlobalDnsProvider.Create(globalDNSProvider)
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    globalDNSProviderStateRefreshFunc(client, newGlobalDNSProvider.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return diag.Errorf(
			"[ERROR] waiting for global dns provider (%s) to be created: %s", newGlobalDNSProvider.ID, waitErr)
	}

	err = flattenGlobalDNSProvider(d, newGlobalDNSProvider)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceRancher2GlobalDNSProviderRead(ctx, d, meta)
}

func resourceRancher2GlobalDNSProviderRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Refreshing Global DNS Provider ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		globalDNSProvider, err := client.GlobalDnsProvider.ByID(d.Id())
		if err != nil {
			if IsNotFound(err) {
				log.Printf("[INFO] Global DNS Provider ID %s not found.", d.Id())
				d.SetId("")
				return nil
			}
			return retry.NonRetryableError(err)
		}

		if err = flattenGlobalDNSProvider(d, globalDNSProvider); err != nil {
			return retry.NonRetryableError(err)
		}

		return nil
	}))
}

func resourceRancher2GlobalDNSProviderUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Updating Global DNS Provider ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	globalDNSProvider, err := client.GlobalDnsProvider.ByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	update := map[string]interface{}{
		"annotations": toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":      toMapString(d.Get("labels").(map[string]interface{})),
	}

	newGlobalDNSProvider, err := client.GlobalDnsProvider.Update(globalDNSProvider, update)
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    globalDNSProviderStateRefreshFunc(client, newGlobalDNSProvider.ID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return diag.Errorf(
			"[ERROR] waiting for global dns provider (%s) to be updated: %s", newGlobalDNSProvider.ID, waitErr)
	}

	return resourceRancher2GlobalDNSProviderRead(ctx, d, meta)
}

func resourceRancher2GlobalDNSProviderDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Deleting Global DNS Provider ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	globalDNSProvider, err := client.GlobalDnsProvider.ByID(id)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Global DNS Provider ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	err = client.GlobalDnsProvider.Delete(globalDNSProvider)
	if err != nil {
		return diag.Errorf("Error removing Global DNS Provider: %s", err)
	}

	log.Printf("[DEBUG] Waiting for global dns provider (%s) to be removed", id)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"removed"},
		Refresh:    globalDNSProviderStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return diag.Errorf(
			"[ERROR] waiting for global dns provider (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

// globalDNSProviderStateRefreshFunc returns a retry.StateRefreshFunc, used to watch a Rancher Global DNS Provider.
func globalDNSProviderStateRefreshFunc(client *managementClient.Client, globalDNSProviderID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := client.GlobalDnsProvider.ByID(globalDNSProviderID)
		if err != nil {
			if IsNotFound(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}

		if obj.Removed != "" {
			return obj, "removed", nil
		}

		return obj, "active", nil
	}
}
