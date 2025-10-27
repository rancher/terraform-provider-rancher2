package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Flatteners

func flattenNodeTemplate(d *schema.ResourceData, in *NodeTemplate) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)
	d.Set("name", in.Name)
	d.Set("driver", in.Driver)

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
	case linodeConfigDriver:
		if in.LinodeConfig == nil {
			return fmt.Errorf("[ERROR] Node template driver %s requires linode_config", in.Driver)
		}
	case openstackConfigDriver:
		if in.OpenstackConfig == nil {
			return fmt.Errorf("[ERROR] Node template driver %s requires openstack_config", in.Driver)
		}
	case hetznerConfigDriver:
		if in.HetznerConfig == nil {
			return fmt.Errorf("[ERROR] Node template driver %s requires hetzner_config", in.Driver)
		}
	case harvesterConfigDriver:
		if in.HarvesterConfig == nil {
			return fmt.Errorf("[ERROR] Node template driver %s requires harvester_config", in.Driver)
		}
	case vmwarevsphereConfigDriver:
		if in.VmwarevsphereConfig == nil {
			return fmt.Errorf("[ERROR] Node template driver %s requires vsphere_config", in.Driver)
		}
	case opennebulaConfigDriver:
		if in.OpennebulaConfig == nil {
			return fmt.Errorf("[ERROR] Node template driver %s requires opennebula_config", in.Driver)
		}
	case outscaleConfigDriver:
		if in.OutscaleConfig == nil {
			return fmt.Errorf("[ERROR] Node template driver %s requires outscale_config", in.Driver)
		}
	default:
		return fmt.Errorf("[ERROR] Unsupported driver on node template: %s", in.Driver)
	}

	if len(in.AuthCertificateAuthority) > 0 {
		d.Set("auth_certificate_authority", in.AuthCertificateAuthority)
	}

	if len(in.AuthKey) > 0 {
		d.Set("auth_key", in.AuthKey)
	}

	if len(in.CloudCredentialID) > 0 {
		d.Set("cloud_credential_id", in.CloudCredentialID)
	}

	if len(in.Description) > 0 {
		d.Set("description", in.Description)
	}

	if len(in.EngineEnv) > 0 {
		err := d.Set("engine_env", toMapInterface(in.EngineEnv))
		if err != nil {
			return err
		}
	}

	if len(in.EngineInsecureRegistry) > 0 {
		err := d.Set("engine_insecure_registry", toArrayInterface(in.EngineInsecureRegistry))
		if err != nil {
			return err
		}
	}

	if len(in.EngineInstallURL) > 0 {
		d.Set("engine_install_url", in.EngineInstallURL)
	}

	if len(in.EngineLabel) > 0 {
		err := d.Set("engine_label", toMapInterface(in.EngineLabel))
		if err != nil {
			return err
		}
	}

	if len(in.EngineOpt) > 0 {
		err := d.Set("engine_opt", toMapInterface(in.EngineOpt))
		if err != nil {
			return err
		}
	}

	if len(in.EngineRegistryMirror) > 0 {
		err := d.Set("engine_registry_mirror", toArrayInterface(in.EngineRegistryMirror))
		if err != nil {
			return err
		}
	}

	if len(in.EngineStorageDriver) > 0 {
		d.Set("engine_storage_driver", in.EngineStorageDriver)
	}

	err := d.Set("node_taints", flattenTaints(in.NodeTaints))
	if err != nil {
		return err
	}

	d.Set("use_internal_ip_address", *in.UseInternalIPAddress)

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

	if v, ok := in.Get("cloud_credential_id").(string); ok && len(v) > 0 {
		obj.CloudCredentialID = v
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

	if v, ok := in.Get("node_taints").([]interface{}); ok && len(v) > 0 {
		obj.NodeTaints = expandTaints(v)
	}

	if v, ok := in.Get("linode_config").([]interface{}); ok && len(v) > 0 {
		obj.LinodeConfig = expandLinodeConfig(v)
		obj.Driver = linodeConfigDriver
	}

	if v, ok := in.Get("openstack_config").([]interface{}); ok && len(v) > 0 {
		obj.OpenstackConfig = expandOpenstackConfig(v)
		obj.Driver = openstackConfigDriver
	}

	if v, ok := in.Get("opennebula_config").([]interface{}); ok && len(v) > 0 {
		obj.OpennebulaConfig = expandOpennebulaConfig(v)
		obj.Driver = opennebulaConfigDriver
	}

	if v, ok := in.Get("outscale_config").([]interface{}); ok && len(v) > 0 {
		obj.OutscaleConfig = expandOutscaleConfig(v)
		obj.Driver = outscaleConfigDriver
	}

	if v, ok := in.Get("hetzner_config").([]interface{}); ok && len(v) > 0 {
		obj.HetznerConfig = expandHetznercloudConfig(v)
		obj.Driver = hetznerConfigDriver
	}

	if v, ok := in.Get("harvester_config").([]interface{}); ok && len(v) > 0 {
		obj.HarvesterConfig = expandHarvestercloudConfig(v)
		obj.Driver = harvesterConfigDriver
	}

	if v, ok := in.Get("use_internal_ip_address").(bool); ok {
		obj.UseInternalIPAddress = &v
	}

	if v, ok := in.Get("vsphere_config").([]interface{}); ok && len(v) > 0 {
		obj.VmwarevsphereConfig = expandVsphereConfig(v)
		obj.Driver = vmwarevsphereConfigDriver
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	// Computing driver_id if empty
	if v, ok := in.Get("driver_id").(string); ok && len(v) == 0 {
		in.Set("driver_id", obj.Driver)
	}

	return obj
}
