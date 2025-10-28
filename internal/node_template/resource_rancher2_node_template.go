package rancher2

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	norman "github.com/rancher/norman/types"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2NodeTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2NodeTemplateCreate,
		Read:   resourceRancher2NodeTemplateRead,
		Update: resourceRancher2NodeTemplateUpdate,
		Delete: resourceRancher2NodeTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2NodeTemplateImport,
		},

		Schema: nodeTemplateFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2NodeTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	nodeTemplate := expandNodeTemplate(d)

	log.Printf("[INFO] Creating Node Template %s", nodeTemplate.Name)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	driverID := d.Get("driver_id").(string)

	err = meta.(*Config).activateNodeDriver(driverID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	newNodeTemplate := &NodeTemplate{}

	err = client.APIBaseClient.Create(managementClient.NodeTemplateType, nodeTemplate, newNodeTemplate)
	if err != nil {
		return err
	}

	d.SetId(newNodeTemplate.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active"},
		Refresh:    nodeTemplateStateRefreshFunc(client, newNodeTemplate.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for node template (%s) to be created: %s", newNodeTemplate.ID, waitErr)
	}

	return resourceRancher2NodeTemplateRead(d, meta)
}

func resourceRancher2NodeTemplateRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Node Template ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		// Normalize node-template ID due to API change
		d.SetId(meta.(*Config).fixNodeTemplateID(d.Id()))

		nodeTemplate := &NodeTemplate{}

		err = client.APIBaseClient.ByID(managementClient.NodeTemplateType, d.Id(), nodeTemplate)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] Node template ID %s not found.", d.Id())
				d.SetId("")
				return nil
			}
			return resource.NonRetryableError(err)
		}

		if err = flattenNodeTemplate(d, nodeTemplate); err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func resourceRancher2NodeTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Node Template ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	nodeTemplate := &norman.Resource{}
	err = client.APIBaseClient.ByID(managementClient.NodeTemplateType, d.Id(), nodeTemplate)
	if err != nil {
		return err
	}

	useInternalIPAddress := d.Get("use_internal_ip_address").(bool)
	update := map[string]interface{}{
		"name":                     d.Get("name").(string),
		"authCertificateAuthority": d.Get("auth_certificate_authority").(string),
		"authKey":                  d.Get("auth_key").(string),
		"cloudCredentialId":        d.Get("cloud_credential_id").(string),
		"description":              d.Get("description").(string),
		"engineEnv":                toMapString(d.Get("engine_env").(map[string]interface{})),
		"engineInsecureRegistry":   toArrayString(d.Get("engine_insecure_registry").([]interface{})),
		"engineInstallURL":         d.Get("engine_install_url").(string),
		"engineLabel":              toMapString(d.Get("engine_label").(map[string]interface{})),
		"engineOpt":                toMapString(d.Get("engine_opt").(map[string]interface{})),
		"engineRegistryMirror":     toArrayString(d.Get("engine_registry_mirror").([]interface{})),
		"engineStorageDriver":      d.Get("engine_storage_driver").(string),
		"nodeTaints":               expandTaints(d.Get("node_taints").([]interface{})),
		"useInternalIpAddress":     &useInternalIPAddress,
		"annotations":              toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":                   toMapString(d.Get("labels").(map[string]interface{})),
	}

	switch driver := d.Get("driver").(string); driver {
	case amazonec2ConfigDriver:
		update["amazonec2Config"] = expandAmazonec2Config(d.Get("amazonec2_config").([]interface{}))
	case azureConfigDriver:
		update["azureConfig"] = expandAzureConfig(d.Get("azure_config").([]interface{}))
	case digitaloceanConfigDriver:
		update["digitaloceanConfig"] = expandDigitaloceanConfig(d.Get("digitalocean_config").([]interface{}))
	case harvesterConfigDriver:
		update["harvesterConfig"] = expandHarvestercloudConfig(d.Get("harvester_config").([]interface{}))
	case hetznerConfigDriver:
		update["hetznerConfig"] = expandHetznercloudConfig(d.Get("hetzner_config").([]interface{}))
	case linodeConfigDriver:
		update["linodeConfig"] = expandLinodeConfig(d.Get("linode_config").([]interface{}))
	case openstackConfigDriver:
		update["openstackConfig"] = expandOpenstackConfig(d.Get("openstack_config").([]interface{}))
	case opennebulaConfigDriver:
		update["opennebulaConfig"] = expandOpennebulaConfig(d.Get("opennebula_config").([]interface{}))
	case vmwarevsphereConfigDriver:
		update["vmwarevsphereConfig"] = expandVsphereConfig(d.Get("vsphere_config").([]interface{}))
	case outscaleConfigDriver:
		update["outscaleConfig"] = expandOutscaleConfig(d.Get("outscale_config").([]interface{}))
	}

	newNodeTemplate := &NodeTemplate{}
	err = client.APIBaseClient.Update(managementClient.NodeTemplateType, nodeTemplate, update, newNodeTemplate)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    nodeTemplateStateRefreshFunc(client, newNodeTemplate.ID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for node template (%s) to be updated: %s", newNodeTemplate.ID, waitErr)
	}

	return resourceRancher2NodeTemplateRead(d, meta)
}

func resourceRancher2NodeTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Node Template ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	nodeTemplate := &norman.Resource{}
	err = client.APIBaseClient.ByID(managementClient.NodeTemplateType, id, nodeTemplate)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Node Template ID %s not found.", id)
			d.SetId("")
			return nil
		}
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), meta.(*Config).Timeout)
	defer cancel()
	for {
		err = client.APIBaseClient.Delete(nodeTemplate)
		if err == nil {
			break
		}
		if !IsNotAllowed(err) {
			return fmt.Errorf("[ERROR] removing Node Template: %s", err)
		}
		select {
		case <-time.After(rancher2RetriesWait * time.Second):
		case <-ctx.Done():
			return fmt.Errorf("[ERROR] timeout removing Node Template: %s", err)
		}
	}

	log.Printf("[DEBUG] Waiting for node template (%s) to be removed", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"removing"},
		Target:     []string{"removed"},
		Refresh:    nodeTemplateStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for node template (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

// nodeTemplateStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher NodeTemplate.
func nodeTemplateStateRefreshFunc(client *managementClient.Client, nodePoolID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj := &NodeTemplate{}
		err := client.APIBaseClient.ByID(managementClient.NodeTemplateType, nodePoolID, obj)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}

		return obj, obj.State, nil
	}
}
