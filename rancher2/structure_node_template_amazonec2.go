package rancher2

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

	obj["encrypt_ebs_volume"] = in.EncryptEbsVolume

	if len(in.Endpoint) > 0 {
		obj["endpoint"] = in.Endpoint
	}

	if len(in.HTTPEndpoint) > 0 {
		obj["http_endpoint"] = in.HTTPEndpoint
	}

	if len(in.HTTPTokens) > 0 {
		obj["http_tokens"] = in.HTTPTokens
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

	if len(in.KmsKey) > 0 {
		obj["kms_key"] = in.KmsKey
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

	if v, ok := in["encrypt_ebs_volume"].(bool); ok {
		obj.EncryptEbsVolume = v
	}

	if v, ok := in["endpoint"].(string); ok && len(v) > 0 {
		obj.Endpoint = v
	}

	if v, ok := in["http_endpoint"].(string); ok && len(v) > 0 {
		obj.HTTPEndpoint = v
	}

	if v, ok := in["http_tokens"].(string); ok && len(v) > 0 {
		obj.HTTPTokens = v
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

	if v, ok := in["kms_key"].(string); ok && len(v) > 0 {
		obj.KmsKey = v
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
