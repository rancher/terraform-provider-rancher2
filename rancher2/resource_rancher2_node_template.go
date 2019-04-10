package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	norman "github.com/rancher/norman/types"
	managementClient "github.com/rancher/types/client/management/v3"
)

type NodeTemplate struct {
	managementClient.NodeTemplate
	Amazonec2Config    *amazonec2Config    `json:"amazonec2Config,omitempty" yaml:"amazonec2Config,omitempty"`
	AzureConfig        *azureConfig        `json:"azureConfig,omitempty" yaml:"azureConfig,omitempty"`
	DigitaloceanConfig *digitaloceanConfig `json:"digitaloceanConfig,omitempty" yaml:"digitaloceanConfig,omitempty"`
}

func nodeTemplateFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"amazonec2_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"azure_config", "digitalocean_config"},
			Elem: &schema.Resource{
				Schema: amazonec2ConfigFields(),
			},
		},
		"auth_certificate_authority": &schema.Schema{
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"auth_key": &schema.Schema{
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"azure_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"amazonec2_config", "digitalocean_config"},
			Elem: &schema.Resource{
				Schema: azureConfigFields(),
			},
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"digitalocean_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"amazonec2_config", "azure_config"},
			Elem: &schema.Resource{
				Schema: digitaloceanConfigFields(),
			},
		},
		"docker_version": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"driver": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"engine_env": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
		},
		"engine_insecure_registry": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"engine_install_url": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"engine_label": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
		},
		"engine_opt": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
		},
		"engine_registry_mirror": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"engine_storage_driver": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"use_internal_ip_address": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"annotations": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
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

func flattenNodeTemplate(d *schema.ResourceData, in *NodeTemplate) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)

	err := d.Set("name", in.Name)
	if err != nil {
		return err
	}

	err = d.Set("driver", in.Driver)
	if err != nil {
		return err
	}

	switch in.Driver {
	case amazonec2ConfigDriver:
		if in.Amazonec2Config == nil {
			return fmt.Errorf("[ERROR] Node template driver %s requires amazonec2_config", in.Driver)
		}
	case azureConfigDriver:
		if in.AzureConfig == nil {
			return fmt.Errorf("[ERROR] Node template driver %s requires azure_config", in.Driver)
		}
	case digitaloceanConfigDriver:
		if in.DigitaloceanConfig == nil {
			return fmt.Errorf("[ERROR] Node template driver %s requires digitalocean_config", in.Driver)
		}
	default:
		return fmt.Errorf("[ERROR] Unsupported driver on node template: %s", in.Driver)
	}

	if len(in.AuthCertificateAuthority) > 0 {
		err = d.Set("auth_certificate_authority", in.AuthCertificateAuthority)
		if err != nil {
			return err
		}
	}

	if len(in.AuthKey) > 0 {
		err = d.Set("auth_key", in.AuthKey)
		if err != nil {
			return err
		}
	}

	if len(in.Description) > 0 {
		err = d.Set("description", in.Description)
		if err != nil {
			return err
		}
	}

	if len(in.DockerVersion) > 0 {
		err = d.Set("docker_version", in.DockerVersion)
		if err != nil {
			return err
		}
	}

	if len(in.EngineEnv) > 0 {
		err = d.Set("engine_env", toMapInterface(in.EngineEnv))
		if err != nil {
			return err
		}
	}

	if len(in.EngineInsecureRegistry) > 0 {
		err = d.Set("engine_insecure_registry", toArrayInterface(in.EngineInsecureRegistry))
		if err != nil {
			return err
		}
	}

	if len(in.EngineInstallURL) > 0 {
		err = d.Set("engine_install_url", in.EngineInstallURL)
		if err != nil {
			return err
		}
	}

	if len(in.EngineLabel) > 0 {
		err = d.Set("engine_label", toMapInterface(in.EngineLabel))
		if err != nil {
			return err
		}
	}

	if len(in.EngineOpt) > 0 {
		err = d.Set("engine_opt", toMapInterface(in.EngineOpt))
		if err != nil {
			return err
		}
	}

	if len(in.EngineRegistryMirror) > 0 {
		err = d.Set("engine_registry_mirror", toArrayInterface(in.EngineRegistryMirror))
		if err != nil {
			return err
		}
	}

	if len(in.EngineStorageDriver) > 0 {
		err = d.Set("engine_storage_driver", in.EngineStorageDriver)
		if err != nil {
			return err
		}
	}

	err = d.Set("use_internal_ip_address", in.UseInternalIPAddress)
	if err != nil {
		return err
	}

	if len(in.Annotations) > 0 {
		err = d.Set("annotations", toMapInterface(in.Annotations))
		if err != nil {
			return err
		}
	}

	if len(in.Labels) > 0 {
		err = d.Set("labels", toMapInterface(in.Labels))
		if err != nil {
			return err
		}
	}

	return nil
}

// Expanders

func expandNodeTemplate(in *schema.ResourceData) *NodeTemplate {
	obj := &NodeTemplate{}
	if in == nil {
		return nil
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}
	obj.Name = in.Get("name").(string)

	if v, ok := in.Get("amazonec2_config").([]interface{}); ok && len(v) > 0 {
		obj.Amazonec2Config = expandAmazonec2Config(v)
		obj.Driver = amazonec2ConfigDriver
	}

	if v, ok := in.Get("auth_certificate_authority").(string); ok && len(v) > 0 {
		obj.AuthCertificateAuthority = v
	}

	if v, ok := in.Get("auth_key").(string); ok && len(v) > 0 {
		obj.AuthKey = v
	}

	if v, ok := in.Get("azure_config").([]interface{}); ok && len(v) > 0 {
		obj.AzureConfig = expandAzureConfig(v)
		obj.Driver = azureConfigDriver
	}

	if v, ok := in.Get("description").(string); ok && len(v) > 0 {
		obj.Description = v
	}

	if v, ok := in.Get("digitalocean_config").([]interface{}); ok && len(v) > 0 {
		obj.DigitaloceanConfig = expandDigitaloceanConfig(v)
		obj.Driver = digitaloceanConfigDriver
	}

	if v, ok := in.Get("engine_env").(map[string]interface{}); ok && len(v) > 0 {
		obj.EngineEnv = toMapString(v)
	}

	if v, ok := in.Get("engine_insecure_registry").([]interface{}); ok && len(v) > 0 {
		obj.EngineInsecureRegistry = toArrayString(v)
	}

	if v, ok := in.Get("engine_install_url").(string); ok && len(v) > 0 {
		obj.EngineInstallURL = v
	}

	if v, ok := in.Get("engine_label").(map[string]interface{}); ok && len(v) > 0 {
		obj.EngineLabel = toMapString(v)
	}

	if v, ok := in.Get("engine_opt").(map[string]interface{}); ok && len(v) > 0 {
		obj.EngineOpt = toMapString(v)
	}

	if v, ok := in.Get("engine_registry_mirror").([]interface{}); ok && len(v) > 0 {
		obj.EngineRegistryMirror = toArrayString(v)
	}

	if v, ok := in.Get("engine_storage_driver").(string); ok && len(v) > 0 {
		obj.EngineStorageDriver = v
	}

	if v, ok := in.Get("use_internal_ip_address").(bool); ok {
		obj.UseInternalIPAddress = v
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj
}

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

	err = meta.(*Config).activateNodeDriver(nodeTemplate.Driver)
	if err != nil {
		return err
	}

	newNodeTemplate := &NodeTemplate{}

	err = client.APIBaseClient.Create(managementClient.NodeTemplateType, nodeTemplate, newNodeTemplate)
	if err != nil {
		return err
	}

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

	d.SetId(newNodeTemplate.ID)

	return resourceRancher2NodeTemplateRead(d, meta)
}

func resourceRancher2NodeTemplateRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Node Template ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	nodeTemplate := &NodeTemplate{}

	err = client.APIBaseClient.ByID(managementClient.NodeTemplateType, d.Id(), nodeTemplate)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Node template ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = flattenNodeTemplate(d, nodeTemplate)
	if err != nil {
		return err
	}

	return nil
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

	update := map[string]interface{}{
		"name":                     d.Get("name").(string),
		"authCertificateAuthority": d.Get("auth_certificate_authority").(string),
		"authKey":                  d.Get("auth_key").(string),
		"description":              d.Get("description").(string),
		"dockerVersion":            d.Get("docker_version").(string),
		"engineEnv":                toMapString(d.Get("engine_env").(map[string]interface{})),
		"engineInsecureRegistry":   toArrayString(d.Get("engine_insecure_registry").([]interface{})),
		"engineInstallURL":         d.Get("engine_install_url").(string),
		"engineLabel":              toMapString(d.Get("engine_label").(map[string]interface{})),
		"engineOpt":                toMapString(d.Get("engine_opt").(map[string]interface{})),
		"engineRegistryMirror":     toArrayString(d.Get("engine_registry_mirror").([]interface{})),
		"engineStorageDriver":      d.Get("engine_storage_driver").(string),
		"annotations":              toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":                   toMapString(d.Get("labels").(map[string]interface{})),
	}

	switch driver := d.Get("driver").(string); driver {
	case amazonec2ConfigDriver:
		update["amazonec2Config"] = expandAmazonec2Config(d.Get("amazonec2_config").([]interface{}))
	case azureConfigDriver:
		update["azureConfig"] = expandAzureConfig(d.Get("azure_config").([]interface{}))
	case digitaloceanConfigDriver:
		update["digitaloceanConfig"] = expandAzureConfig(d.Get("digitalocean_config").([]interface{}))
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
		if IsNotFound(err) {
			log.Printf("[INFO] Node Template ID %s not found.", id)
			d.SetId("")
			return nil
		}
		return err
	}

	err = client.APIBaseClient.Delete(nodeTemplate)
	if err != nil {
		return fmt.Errorf("Error removing Node Template: %s", err)
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

func resourceRancher2NodeTemplateImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return []*schema.ResourceData{}, err
	}
	nodeTemplate := &NodeTemplate{}
	err = client.APIBaseClient.ByID(managementClient.NodeTemplateType, d.Id(), nodeTemplate)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	err = flattenNodeTemplate(d, nodeTemplate)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}

// nodeTemplateStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher NodeTemplate.
func nodeTemplateStateRefreshFunc(client *managementClient.Client, nodePoolID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj := &NodeTemplate{}
		err := client.APIBaseClient.ByID(managementClient.NodeTemplateType, nodePoolID, obj)
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
