package rancher2

// Flatteners

func flattenOpenstackConfig(in *openstackConfig) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	obj["active_timeout"] = in.ActiveTimeout
	obj["auth_url"] = in.AuthURL
	obj["availability_zone"] = in.AvailabilityZone
	obj["cacert"] = in.CaCert
	obj["config_drive"] = in.ConfigDrive
	obj["domain_id"] = in.DomainID
	obj["domain_name"] = in.DomainName
	obj["endpoint_type"] = in.EndpointType
	obj["flavor_id"] = in.FlavorID
	obj["flavor_name"] = in.FlavorName
	obj["floating_ip_pool"] = in.FloatingIPPool
	obj["image_id"] = in.ImageID
	obj["image_name"] = in.ImageName
	obj["insecure"] = in.Insecure
	obj["ip_version"] = in.IPVersion
	obj["keypair_name"] = in.KeypairName
	obj["net_id"] = in.NetID
	obj["net_name"] = in.NetName
	obj["nova_network"] = in.NovaNetwork
	obj["password"] = in.Password
	obj["private_key_file"] = in.PrivateKeyFile
	obj["region"] = in.Region
	obj["sec_groups"] = in.SecGroups
	obj["ssh_port"] = in.SSHPort
	obj["ssh_user"] = in.SSHUser
	obj["tenant_id"] = in.TenantID
	obj["tenant_name"] = in.TenantName
	obj["user_data_file"] = in.UserDataFile
	obj["username"] = in.Username
	obj["application_credential_id"] = in.ApplicationCredentialID
	obj["application_credential_name"] = in.ApplicationCredentialName
	obj["application_credential_secret"] = in.ApplicationCredentialSecret
	obj["boot_from_volume"] = in.BootFromVolume
	obj["volume_size"] = in.VolumeSize
	obj["volume_type"] = in.VolumeType
	obj["volume_id"] = in.VolumeID
	obj["volume_name"] = in.VolumeName
	obj["volume_device_path"] = in.VolumeDevicePath

	return []interface{}{obj}
}

// Expanders

func expandOpenstackConfig(p []interface{}) *openstackConfig {
	obj := &openstackConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["active_timeout"].(string); ok && len(v) > 0 {
		obj.ActiveTimeout = v
	}
	if v, ok := in["auth_url"].(string); ok && len(v) > 0 {
		obj.AuthURL = v
	}
	if v, ok := in["availability_zone"].(string); ok && len(v) > 0 {
		obj.AvailabilityZone = v
	}
	if v, ok := in["cacert"].(string); ok && len(v) > 0 {
		obj.CaCert = v
	}
	if v, ok := in["config_drive"].(bool); ok {
		obj.ConfigDrive = v
	}
	if v, ok := in["domain_id"].(string); ok && len(v) > 0 {
		obj.DomainID = v
	}
	if v, ok := in["domain_name"].(string); ok && len(v) > 0 {
		obj.DomainName = v
	}
	if v, ok := in["endpoint_type"].(string); ok && len(v) > 0 {
		obj.EndpointType = v
	}
	if v, ok := in["flavor_id"].(string); ok && len(v) > 0 {
		obj.FlavorID = v
	}
	if v, ok := in["flavor_name"].(string); ok && len(v) > 0 {
		obj.FlavorName = v
	}
	if v, ok := in["floating_ip_pool"].(string); ok && len(v) > 0 {
		obj.FloatingIPPool = v
	}
	if v, ok := in["ip_version"].(string); ok && len(v) > 0 {
		obj.IPVersion = v
	}
	if v, ok := in["image_id"].(string); ok && len(v) > 0 {
		obj.ImageID = v
	}
	if v, ok := in["image_name"].(string); ok && len(v) > 0 {
		obj.ImageName = v
	}
	if v, ok := in["insecure"].(bool); ok {
		obj.Insecure = v
	}
	if v, ok := in["ip_version"].(string); ok && len(v) > 0 {
		obj.IPVersion = v
	}
	if v, ok := in["keypair_name"].(string); ok && len(v) > 0 {
		obj.KeypairName = v
	}
	if v, ok := in["net_id"].(string); ok && len(v) > 0 {
		obj.NetID = v
	}
	if v, ok := in["net_name"].(string); ok && len(v) > 0 {
		obj.NetName = v
	}
	if v, ok := in["nova_network"].(bool); ok {
		obj.NovaNetwork = v
	}
	if v, ok := in["password"].(string); ok && len(v) > 0 {
		obj.Password = v
	}
	if v, ok := in["private_key_file"].(string); ok && len(v) > 0 {
		obj.PrivateKeyFile = v
	}
	if v, ok := in["region"].(string); ok && len(v) > 0 {
		obj.Region = v
	}
	if v, ok := in["sec_groups"].(string); ok && len(v) > 0 {
		obj.SecGroups = v
	}
	if v, ok := in["ssh_port"].(string); ok && len(v) > 0 {
		obj.SSHPort = v
	}
	if v, ok := in["ssh_user"].(string); ok && len(v) > 0 {
		obj.SSHUser = v
	}
	if v, ok := in["tenant_id"].(string); ok && len(v) > 0 {
		obj.TenantID = v
	}
	if v, ok := in["tenant_name"].(string); ok && len(v) > 0 {
		obj.TenantName = v
	}
	if v, ok := in["user_data_file"].(string); ok && len(v) > 0 {
		obj.UserDataFile = v
	}
	if v, ok := in["username"].(string); ok && len(v) > 0 {
		obj.Username = v
	}
	if v, ok := in["application_credential_id"].(string); ok && len(v) > 0 {
		obj.ApplicationCredentialID = v
	}
	if v, ok := in["application_credential_name"].(string); ok && len(v) > 0 {
		obj.ApplicationCredentialName = v
	}
	if v, ok := in["application_credential_secret"].(string); ok && len(v) > 0 {
		obj.ApplicationCredentialSecret = v
	}
	if v, ok := in["boot_from_volume"].(bool); ok {
		obj.BootFromVolume = v
	}
	if v, ok := in["volume_size"].(string); ok && len(v) > 0 {
		obj.VolumeSize = v
	}
	if v, ok := in["volume_type"].(string); ok && len(v) > 0 {
		obj.VolumeType = v
	}
	if v, ok := in["volume_id"].(string); ok && len(v) > 0 {
		obj.VolumeID = v
	}
	if v, ok := in["volume_name"].(string); ok && len(v) > 0 {
		obj.VolumeName = v
	}
	if v, ok := in["volume_device_path"].(string); ok && len(v) > 0 {
		obj.VolumeDevicePath = v
	}
	return obj
}
