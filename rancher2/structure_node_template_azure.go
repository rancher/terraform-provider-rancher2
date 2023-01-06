package rancher2

// Flatteners

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

	if len(in.DiskSize) > 0 {
		obj["disk_size"] = in.DiskSize
	}

	if len(in.DNS) > 0 {
		obj["dns"] = in.DNS
	}

	if len(in.Environment) > 0 {
		obj["environment"] = in.Environment
	}

	if len(in.FaultDomainCount) > 0 {
		obj["fault_domain_count"] = in.FaultDomainCount
	}

	if len(in.Image) > 0 {
		obj["image"] = in.Image
	}

	if len(in.Location) > 0 {
		obj["location"] = in.Location
	}

	obj["managed_disks"] = in.ManagedDisks
	obj["no_public_ip"] = in.NoPublicIP

	if len(in.NSG) > 0 {
		obj["nsg"] = in.NSG
	}

	if len(in.Plan) > 0 {
		obj["plan"] = in.Plan
	}

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

	if len(in.UpdateDomainCount) > 0 {
		obj["update_domain_count"] = in.UpdateDomainCount
	}

	obj["use_private_ip"] = in.UsePrivateIP
	obj["use_public_ip_standard_sku"] = in.UsePublicIPStandardSKU
	obj["accelerated_networking"] = in.AcceleratedNetworking

	if len(in.Vnet) > 0 {
		obj["vnet"] = in.Vnet
	}

	if len(in.AvailabilityZone) > 0 {
		obj["availability_zone"] = in.AvailabilityZone
	}

	if len(in.Tags) > 0 {
		obj["tags"] = in.Tags
	}

	return []interface{}{obj}
}

// Expanders

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

	if v, ok := in["disk_size"].(string); ok && len(v) > 0 {
		obj.DiskSize = v
	}

	if v, ok := in["dns"].(string); ok && len(v) > 0 {
		obj.DNS = v
	}

	if v, ok := in["environment"].(string); ok && len(v) > 0 {
		obj.Environment = v
	}

	if v, ok := in["fault_domain_count"].(string); ok && len(v) > 0 {
		obj.FaultDomainCount = v
	}

	if v, ok := in["image"].(string); ok && len(v) > 0 {
		obj.Image = v
	}

	if v, ok := in["location"].(string); ok && len(v) > 0 {
		obj.Location = v
	}

	if v, ok := in["managed_disks"].(bool); ok {
		obj.ManagedDisks = v
	}

	if v, ok := in["no_public_ip"].(bool); ok {
		obj.NoPublicIP = v
	}

	if v, ok := in["nsg"].(string); ok && len(v) > 0 {
		obj.NSG = v
	}

	if v, ok := in["plan"].(string); ok && len(v) > 0 {
		obj.Plan = v
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

	if v, ok := in["update_domain_count"].(string); ok && len(v) > 0 {
		obj.UpdateDomainCount = v
	}

	if v, ok := in["use_private_ip"].(bool); ok {
		obj.UsePrivateIP = v
	}

	if v, ok := in["vnet"].(string); ok && len(v) > 0 {
		obj.Vnet = v
	}

	if v, ok := in["use_public_ip_standard_sku"].(bool); ok {
		obj.UsePublicIPStandardSKU = v
	}

	if v, ok := in["accelerated_networking"].(bool); ok {
		obj.AcceleratedNetworking = v
	}

	if v, ok := in["availability_zone"].(string); ok && len(v) > 0 {
		obj.AvailabilityZone = v
	}

	if v, ok := in["tags"].(string); ok && len(v) > 0 {
		obj.Tags = v
	}

	return obj
}
