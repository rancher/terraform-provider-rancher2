package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

// Flatteners

func flattenAmazonec2Config(in *amazonec2Config) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	if len(in.AccessKey) > 0 {
		obj["access_key"] = in.AccessKey
	}

	if len(in.Ami) > 0 {
		obj["ami"] = in.Ami
	}

	if len(in.BlockDurationMinutes) > 0 {
		obj["block_duration_minutes"] = in.BlockDurationMinutes
	}

	if len(in.DeviceName) > 0 {
		obj["device_name"] = in.DeviceName
	}

	if len(in.Endpoint) > 0 {
		obj["endpoint"] = in.Endpoint
	}

	if len(in.IamInstanceProfile) > 0 {
		obj["iam_instance_profile"] = in.IamInstanceProfile
	}

	obj["insecure_transport"] = in.InsecureTransport

	if len(in.InstanceType) > 0 {
		obj["instance_type"] = in.InstanceType
	}

	if len(in.KeypairName) > 0 {
		obj["keypair_name"] = in.KeypairName
	}

	obj["monitoring"] = in.Monitoring

	if len(in.OpenPort) > 0 {
		obj["open_port"] = toArrayInterface(in.OpenPort)
	}

	obj["private_address_only"] = in.PrivateAddressOnly

	if len(in.Region) > 0 {
		obj["region"] = in.Region
	}

	obj["request_spot_instance"] = in.RequestSpotInstance

	if len(in.Retries) > 0 {
		obj["retries"] = in.Retries
	}

	if len(in.RootSize) > 0 {
		obj["root_size"] = in.RootSize
	}

	if len(in.SecretKey) > 0 {
		obj["secret_key"] = in.SecretKey
	}

	if len(in.SecurityGroup) > 0 {
		obj["security_group"] = toArrayInterface(in.SecurityGroup)
	}

	obj["security_group_readonly"] = in.SecurityGroupReadonly

	if len(in.SessionToken) > 0 {
		obj["session_token"] = in.SessionToken
	}

	if len(in.SpotPrice) > 0 {
		obj["spot_price"] = in.SpotPrice
	}

	if len(in.SSHKeypath) > 0 {
		obj["ssh_keypath"] = in.SSHKeypath
	}

	if len(in.SSHUser) > 0 {
		obj["ssh_user"] = in.SSHUser
	}

	if len(in.SubnetID) > 0 {
		obj["subnet_id"] = in.SubnetID
	}

	if len(in.Tags) > 0 {
		obj["tags"] = in.Tags
	}

	obj["use_ebs_optimized_instance"] = in.UseEbsOptimizedInstance

	obj["use_private_address"] = in.UsePrivateAddress

	if len(in.Userdata) > 0 {
		obj["userdata"] = in.Userdata
	}

	if len(in.VolumeType) > 0 {
		obj["volume_type"] = in.VolumeType
	}

	if len(in.VpcID) > 0 {
		obj["vpc_id"] = in.VpcID
	}

	if len(in.Zone) > 0 {
		obj["zone"] = in.Zone
	}

	return []interface{}{obj}
}

func flattenAzureConfig(in *azureConfig) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	if len(in.AvailabilitySet) > 0 {
		obj["availability_set"] = in.AvailabilitySet
	}

	if len(in.ClientID) > 0 {
		obj["client_id"] = in.ClientID
	}

	if len(in.ClientSecret) > 0 {
		obj["client_secret"] = in.ClientSecret
	}

	if len(in.CustomData) > 0 {
		obj["custom_data"] = in.CustomData
	}

	if len(in.DNS) > 0 {
		obj["dns"] = in.DNS
	}

	if len(in.Environment) > 0 {
		obj["environment"] = in.Environment
	}

	if len(in.Image) > 0 {
		obj["image"] = in.Image
	}

	if len(in.Location) > 0 {
		obj["location"] = in.Location
	}

	obj["no_public_ip"] = in.NoPublicIP

	if len(in.OpenPort) > 0 {
		obj["open_port"] = toArrayInterface(in.OpenPort)
	}

	obj["private_address_only"] = in.PrivateAddressOnly

	if len(in.PrivateIPAddress) > 0 {
		obj["private_ip_address"] = in.PrivateIPAddress
	}

	if len(in.ResourceGroup) > 0 {
		obj["resource_group"] = in.ResourceGroup
	}

	if len(in.Size) > 0 {
		obj["size"] = in.Size
	}

	if len(in.SSHUser) > 0 {
		obj["ssh_user"] = in.SSHUser
	}

	obj["static_public_ip"] = in.StaticPublicIP

	if len(in.StorageType) > 0 {
		obj["storage_type"] = in.StorageType
	}

	if len(in.Subnet) > 0 {
		obj["subnet"] = in.Subnet
	}

	if len(in.SubnetPrefix) > 0 {
		obj["subnet_prefix"] = in.SubnetPrefix
	}

	if len(in.SubscriptionID) > 0 {
		obj["subscription_id"] = in.SubscriptionID
	}

	obj["use_private_ip"] = in.UsePrivateIP

	if len(in.Vnet) > 0 {
		obj["vnet"] = in.Vnet
	}

	return []interface{}{obj}
}

func flattenDigitaloceanConfig(in *digitaloceanConfig) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	if len(in.AccessToken) > 0 {
		obj["access_token"] = in.AccessToken
	}

	obj["backups"] = in.Backups

	if len(in.Image) > 0 {
		obj["image"] = in.Image
	}

	obj["ipv6"] = in.IPV6
	obj["monitoring"] = in.Monitoring
	obj["private_networking"] = in.PrivateNetworking

	if len(in.Region) > 0 {
		obj["region"] = in.Region
	}

	if len(in.Size) > 0 {
		obj["size"] = in.Size
	}

	if len(in.SSHKeyFingerprint) > 0 {
		obj["ssh_key_fingerprint"] = in.SSHKeyFingerprint
	}

	if len(in.SSHKeyPath) > 0 {
		obj["ssh_key_path"] = in.SSHKeyPath
	}

	if len(in.SSHPort) > 0 {
		obj["ssh_port"] = in.SSHPort
	}

	if len(in.SSHUser) > 0 {
		obj["ssh_user"] = in.SSHUser
	}

	if len(in.Tags) > 0 {
		obj["tags"] = in.Tags
	}

	if len(in.Userdata) > 0 {
		obj["userdata"] = in.Userdata
	}

	return []interface{}{obj}
}

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

func expandAmazonec2Config(p []interface{}) *amazonec2Config {
	obj := &amazonec2Config{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["access_key"].(string); ok && len(v) > 0 {
		obj.AccessKey = v
	}

	if v, ok := in["ami"].(string); ok && len(v) > 0 {
		obj.Ami = v
	}

	if v, ok := in["block_duration_minutes"].(string); ok && len(v) > 0 {
		obj.BlockDurationMinutes = v
	}

	if v, ok := in["device_name"].(string); ok && len(v) > 0 {
		obj.DeviceName = v
	}

	if v, ok := in["endpoint"].(string); ok && len(v) > 0 {
		obj.Endpoint = v
	}

	if v, ok := in["iam_instance_profile"].(string); ok && len(v) > 0 {
		obj.IamInstanceProfile = v
	}

	if v, ok := in["insecure_transport"].(bool); ok {
		obj.InsecureTransport = v
	}

	if v, ok := in["instance_type"].(string); ok && len(v) > 0 {
		obj.InstanceType = v
	}

	if v, ok := in["keypair_name"].(string); ok && len(v) > 0 {
		obj.KeypairName = v
	}

	if v, ok := in["monitoring"].(bool); ok {
		obj.Monitoring = v
	}

	if v, ok := in["open_port"].([]interface{}); ok && len(v) > 0 {
		obj.OpenPort = toArrayString(v)
	}

	if v, ok := in["private_address_only"].(bool); ok {
		obj.PrivateAddressOnly = v
	}

	if v, ok := in["region"].(string); ok && len(v) > 0 {
		obj.Region = v
	}

	if v, ok := in["request_spot_instance"].(bool); ok {
		obj.RequestSpotInstance = v
	}

	if v, ok := in["retries"].(string); ok && len(v) > 0 {
		obj.Retries = v
	}

	if v, ok := in["root_size"].(string); ok && len(v) > 0 {
		obj.RootSize = v
	}

	if v, ok := in["secret_key"].(string); ok && len(v) > 0 {
		obj.SecretKey = v
	}

	if v, ok := in["security_group"].([]interface{}); ok && len(v) > 0 {
		obj.SecurityGroup = toArrayString(v)
	}

	if v, ok := in["security_group_readonly"].(bool); ok {
		obj.SecurityGroupReadonly = v
	}

	if v, ok := in["session_token"].(string); ok && len(v) > 0 {
		obj.SessionToken = v
	}

	if v, ok := in["spot_price"].(string); ok && len(v) > 0 {
		obj.SpotPrice = v
	}

	if v, ok := in["ssh_keypath"].(string); ok && len(v) > 0 {
		obj.SSHKeypath = v
	}

	if v, ok := in["ssh_user"].(string); ok && len(v) > 0 {
		obj.SSHUser = v
	}

	if v, ok := in["subnet_id"].(string); ok && len(v) > 0 {
		obj.SubnetID = v
	}

	if v, ok := in["tags"].(string); ok && len(v) > 0 {
		obj.Tags = v
	}

	if v, ok := in["use_ebs_optimized_instance"].(bool); ok {
		obj.UseEbsOptimizedInstance = v
	}

	if v, ok := in["use_private_address"].(bool); ok {
		obj.UsePrivateAddress = v
	}

	if v, ok := in["userdata"].(string); ok && len(v) > 0 {
		obj.Userdata = v
	}

	if v, ok := in["volume_type"].(string); ok && len(v) > 0 {
		obj.VolumeType = v
	}

	if v, ok := in["vpc_id"].(string); ok && len(v) > 0 {
		obj.VpcID = v
	}

	if v, ok := in["zone"].(string); ok && len(v) > 0 {
		obj.Zone = v
	}

	return obj
}

func expandAzureConfig(p []interface{}) *azureConfig {
	obj := &azureConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["availability_set"].(string); ok && len(v) > 0 {
		obj.AvailabilitySet = v
	}

	if v, ok := in["client_id"].(string); ok && len(v) > 0 {
		obj.ClientID = v
	}

	if v, ok := in["client_secret"].(string); ok && len(v) > 0 {
		obj.ClientSecret = v
	}

	if v, ok := in["custom_data"].(string); ok && len(v) > 0 {
		obj.CustomData = v
	}

	if v, ok := in["dns"].(string); ok && len(v) > 0 {
		obj.DNS = v
	}

	if v, ok := in["environment"].(string); ok && len(v) > 0 {
		obj.Environment = v
	}

	if v, ok := in["image"].(string); ok && len(v) > 0 {
		obj.Image = v
	}

	if v, ok := in["location"].(string); ok && len(v) > 0 {
		obj.Location = v
	}

	if v, ok := in["no_public_ip"].(bool); ok {
		obj.NoPublicIP = v
	}

	if v, ok := in["open_port"].([]interface{}); ok && len(v) > 0 {
		obj.OpenPort = toArrayString(v)
	}

	if v, ok := in["private_address_only"].(bool); ok {
		obj.PrivateAddressOnly = v
	}

	if v, ok := in["private_ip_address"].(string); ok && len(v) > 0 {
		obj.PrivateIPAddress = v
	}

	if v, ok := in["resource_group"].(string); ok && len(v) > 0 {
		obj.ResourceGroup = v
	}

	if v, ok := in["size"].(string); ok && len(v) > 0 {
		obj.Size = v
	}

	if v, ok := in["ssh_user"].(string); ok && len(v) > 0 {
		obj.SSHUser = v
	}

	if v, ok := in["static_public_ip"].(bool); ok {
		obj.StaticPublicIP = v
	}

	if v, ok := in["storage_type"].(string); ok && len(v) > 0 {
		obj.StorageType = v
	}

	if v, ok := in["subnet"].(string); ok && len(v) > 0 {
		obj.Subnet = v
	}

	if v, ok := in["subnet_prefix"].(string); ok && len(v) > 0 {
		obj.SubnetPrefix = v
	}

	if v, ok := in["subscription_id"].(string); ok && len(v) > 0 {
		obj.SubscriptionID = v
	}

	if v, ok := in["use_private_ip"].(bool); ok {
		obj.UsePrivateIP = v
	}

	if v, ok := in["vnet"].(string); ok && len(v) > 0 {
		obj.Vnet = v
	}

	return obj
}

func expandDigitaloceanConfig(p []interface{}) *digitaloceanConfig {
	obj := &digitaloceanConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["access_token"].(string); ok && len(v) > 0 {
		obj.AccessToken = v
	}

	if v, ok := in["backups"].(bool); ok {
		obj.Backups = v
	}

	if v, ok := in["image"].(string); ok && len(v) > 0 {
		obj.Image = v
	}

	if v, ok := in["ipv6"].(bool); ok {
		obj.IPV6 = v
	}

	if v, ok := in["monitoring"].(bool); ok {
		obj.Monitoring = v
	}
	if v, ok := in["private_networking"].(bool); ok {
		obj.PrivateNetworking = v
	}

	if v, ok := in["region"].(string); ok && len(v) > 0 {
		obj.Region = v
	}

	if v, ok := in["size"].(string); ok && len(v) > 0 {
		obj.Size = v
	}

	if v, ok := in["ssh_key_fingerprint"].(string); ok && len(v) > 0 {
		obj.SSHKeyFingerprint = v
	}

	if v, ok := in["ssh_key_path"].(string); ok && len(v) > 0 {
		obj.SSHKeyPath = v
	}

	if v, ok := in["ssh_port"].(string); ok && len(v) > 0 {
		obj.SSHPort = v
	}

	if v, ok := in["ssh_user"].(string); ok && len(v) > 0 {
		obj.SSHUser = v
	}

	if v, ok := in["tags"].(string); ok && len(v) > 0 {
		obj.Tags = v
	}

	if v, ok := in["userdata"].(string); ok && len(v) > 0 {
		obj.Userdata = v
	}

	return obj
}

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
