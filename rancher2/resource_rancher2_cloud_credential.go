package rancher2

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	norman "github.com/rancher/norman/types"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2CloudCredential() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRancher2CloudCredentialCreate,
		ReadContext:   resourceRancher2CloudCredentialRead,
		UpdateContext: resourceRancher2CloudCredentialUpdate,
		DeleteContext: resourceRancher2CloudCredentialDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRancher2CloudCredentialsImport,
		},
		Schema: cloudCredentialFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2CloudCredentialCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cloudCredential := expandCloudCredential(d)

	log.Printf("[INFO] Creating Cloud Credential %s", cloudCredential.Name)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	if nodeDriver, ok := d.Get("driver").(string); ok && nodeDriver != s3ConfigDriver {
		err = meta.(*Config).activateDriver(nodeDriver, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	newCloudCredential := &CloudCredential{}
	err = client.APIBaseClient.Create(managementClient.CloudCredentialType, cloudCredential, newCloudCredential)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newCloudCredential.ID)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active"},
		Refresh:    cloudCredentialStateRefreshFunc(client, newCloudCredential.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return diag.Errorf("[ERROR] waiting for cloud credential (%s) to be created: %s", newCloudCredential.ID, waitErr)
	}

	return resourceRancher2CloudCredentialRead(ctx, d, meta)
}

func resourceRancher2CloudCredentialRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Refreshing Cloud Credential ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		cloudCredential := &CloudCredential{}
		err = client.APIBaseClient.ByID(managementClient.CloudCredentialType, d.Id(), cloudCredential)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] Cloud Credential ID %s not found.", d.Id())
				d.SetId("")
				return nil
			}
			return retry.NonRetryableError(err)
		}

		if err = flattenCloudCredential(d, cloudCredential); err != nil {
			return retry.NonRetryableError(err)
		}

		return nil
	}))
}

func resourceRancher2CloudCredentialUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Updating Cloud Credential ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudCredential := &norman.Resource{}
	err = client.APIBaseClient.ByID(managementClient.CloudCredentialType, d.Id(), cloudCredential)
	if err != nil {
		return diag.FromErr(err)
	}

	update := map[string]interface{}{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
		"annotations": toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":      toMapString(d.Get("labels").(map[string]interface{})),
	}

	switch driver := d.Get("driver").(string); driver {
	case amazonec2ConfigDriver:
		update["amazonec2credentialConfig"] = expandCloudCredentialAmazonec2(d.Get("amazonec2_credential_config").([]interface{}))
	case azureConfigDriver:
		update["azurecredentialConfig"] = expandCloudCredentialAzure(d.Get("azure_credential_config").([]interface{}))
	case digitaloceanConfigDriver:
		update["digitaloceancredentialConfig"] = expandCloudCredentialDigitalocean(d.Get("digitalocean_credential_config").([]interface{}))
	case googleConfigDriver:
		update["googlecredentialConfig"] = expandCloudCredentialGoogle(d.Get("google_credential_config").([]interface{}))
	case harvesterConfigDriver:
		update["harvestercredentialConfig"] = expandCloudCredentialHarvester(d.Get("harvester_credential_config").([]interface{}))
	case linodeConfigDriver:
		update["linodecredentialConfig"] = expandCloudCredentialLinode(d.Get("linode_credential_config").([]interface{}))
	case openstackConfigDriver:
		update["openstackcredentialConfig"] = expandCloudCredentialOpenstack(d.Get("openstack_credential_config").([]interface{}))
	case s3ConfigDriver:
		update["s3credentialConfig"] = expandCloudCredentialS3(d.Get("s3_credential_config").([]interface{}))
	case vmwarevsphereConfigDriver:
		update["vmwarevspherecredentialConfig"] = expandCloudCredentialVsphere(d.Get("vsphere_credential_config").([]interface{}))
	default:
		return diag.Errorf("[ERROR] updating cloud credential: Unsupported driver \"%s\"", driver)
	}

	newCloudCredential := &CloudCredential{}
	err = client.APIBaseClient.Update(managementClient.CloudCredentialType, cloudCredential, update, newCloudCredential)
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    cloudCredentialStateRefreshFunc(client, newCloudCredential.ID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return diag.Errorf(
			"[ERROR] waiting for cloud credential (%s) to be updated: %s", newCloudCredential.ID, waitErr)
	}

	return resourceRancher2CloudCredentialRead(ctx, d, meta)
}

func resourceRancher2CloudCredentialDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Deleting Cloud Credential ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudCredential := &norman.Resource{}
	err = client.APIBaseClient.ByID(managementClient.CloudCredentialType, id, cloudCredential)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Cloud Credential ID %s not found.", id)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	err = client.APIBaseClient.Delete(cloudCredential)
	if err != nil {
		return diag.Errorf("Error removing Cloud Credential: %s", err)
	}

	log.Printf("[DEBUG] Waiting for cloud credential (%s) to be removed", id)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"removed"},
		Refresh:    cloudCredentialStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return diag.Errorf("[ERROR] waiting for cloud credential (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

// cloudCredentialStateRefreshFunc returns a retry.StateRefreshFunc, used to watch a Rancher CloudCredential.
func cloudCredentialStateRefreshFunc(client *managementClient.Client, credentialID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj := &CloudCredential{}
		err := client.APIBaseClient.ByID(managementClient.CloudCredentialType, credentialID, obj)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}

		return obj, "active", nil
	}
}
