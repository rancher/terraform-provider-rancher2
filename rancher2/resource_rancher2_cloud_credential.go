package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	norman "github.com/rancher/norman/types"
	managementClient "github.com/rancher/types/client/management/v3"
)

func resourceRancher2CloudCredential() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2CloudCredentialCreate,
		Read:   resourceRancher2CloudCredentialRead,
		Update: resourceRancher2CloudCredentialUpdate,
		Delete: resourceRancher2CloudCredentialDelete,

		Schema: cloudCredentialFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2CloudCredentialCreate(d *schema.ResourceData, meta interface{}) error {
	cloudCredential := expandCloudCredential(d)

	log.Printf("[INFO] Creating Cloud Credential %s", cloudCredential.Name)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	nodeDriver := d.Get("driver").(string)
	err = meta.(*Config).activateNodeDriver(nodeDriver)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active"},
		Refresh:    nodeDriverStateRefreshFunc(client, nodeDriver),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for cloud credential (%s) to be activated: %s", nodeDriver, waitErr)
	}

	newCloudCredential := &CloudCredential{}
	err = client.APIBaseClient.Create(managementClient.CloudCredentialType, cloudCredential, newCloudCredential)
	if err != nil {
		return err
	}

	stateConf = &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active"},
		Refresh:    cloudCredentialStateRefreshFunc(client, newCloudCredential.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr = stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for cloud credential (%s) to be created: %s", newCloudCredential.ID, waitErr)
	}

	d.SetId(newCloudCredential.ID)

	return resourceRancher2CloudCredentialRead(d, meta)
}

func resourceRancher2CloudCredentialRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Cloud Credential ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	cloudCredential := &CloudCredential{}
	err = client.APIBaseClient.ByID(managementClient.CloudCredentialType, d.Id(), cloudCredential)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Cloud Credential ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = flattenCloudCredential(d, cloudCredential)
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2CloudCredentialUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Cloud Credential ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	cloudCredential := &norman.Resource{}
	err = client.APIBaseClient.ByID(managementClient.CloudCredentialType, d.Id(), cloudCredential)
	if err != nil {
		return err
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
	case openstackConfigDriver:
		update["openstackcredentialConfig"] = expandCloudCredentialOpenstack(d.Get("openstack_credential_config").([]interface{}))
	case vmwarevsphereConfigDriver:
		update["vmwarevspherecredentialConfig"] = expandCloudCredentialVsphere(d.Get("vsphere_credential_config").([]interface{}))
	default:
		return fmt.Errorf("[ERROR] updating cloud credential: Unsupported driver \"%s\"", driver)
	}

	newCloudCredential := &CloudCredential{}
	err = client.APIBaseClient.Update(managementClient.CloudCredentialType, cloudCredential, update, newCloudCredential)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    cloudCredentialStateRefreshFunc(client, newCloudCredential.ID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for cloud credential (%s) to be updated: %s", newCloudCredential.ID, waitErr)
	}

	return resourceRancher2CloudCredentialRead(d, meta)
}

func resourceRancher2CloudCredentialDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Cloud Credential ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	cloudCredential := &norman.Resource{}
	err = client.APIBaseClient.ByID(managementClient.CloudCredentialType, id, cloudCredential)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Cloud Credential ID %s not found.", id)
			d.SetId("")
			return nil
		}
		return err
	}

	err = client.APIBaseClient.Delete(cloudCredential)
	if err != nil {
		return fmt.Errorf("Error removing Cloud Credential: %s", err)
	}

	log.Printf("[DEBUG] Waiting for cloud credential (%s) to be removed", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"removed"},
		Refresh:    cloudCredentialStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for cloud credential (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

// cloudCredentialStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher CloudCredential.
func cloudCredentialStateRefreshFunc(client *managementClient.Client, credentialID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj := &CloudCredential{}
		err := client.APIBaseClient.ByID(managementClient.CloudCredentialType, credentialID, obj)
		if err != nil {
			if IsNotFound(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}

		return obj, "active", nil
	}
}
