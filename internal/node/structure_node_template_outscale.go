package rancher2

// Flatteners

func flattenOutscaleConfig(in *outscaleConfig) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	if len(in.AccessKey) > 0 {
		obj["access_key"] = in.AccessKey
	}

	if len(in.SourceOmi) > 0 {
		obj["source_omi"] = in.SourceOmi
	}

	if len(in.InstanceType) > 0 {
		obj["instance_type"] = in.InstanceType
	}

	if len(in.Region) > 0 {
		obj["region"] = in.Region
	}

	if len(in.SecretKey) > 0 {
		obj["secret_key"] = in.SecretKey
	}

	if len(in.SecurityGroupIds) > 0 {
		obj["security_group_ids"] = toArrayInterface(in.SecurityGroupIds)
	}

	if len(in.ExtraTagsAll) > 0 {
		obj["extra_tags_all"] = toArrayInterface(in.ExtraTagsAll)
	}

	if len(in.ExtraTagsInstances) > 0 {
		obj["extra_tags_instances"] = toArrayInterface(in.ExtraTagsInstances)
	}

	if len(in.RootDiskType) > 0 {
		obj["root_disk_type"] = in.RootDiskType
	}

	obj["root_disk_size"] = in.RootDiskSize
	obj["root_disk_iops"] = in.RootDiskIops

	return []interface{}{obj}
}

// Expanders

func expandOutscaleConfig(p []interface{}) *outscaleConfig {
	obj := &outscaleConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["access_key"].(string); ok && len(v) > 0 {
		obj.AccessKey = v
	}

	if v, ok := in["source_omi"].(string); ok && len(v) > 0 {
		obj.SourceOmi = v
	}

	if v, ok := in["instance_type"].(string); ok && len(v) > 0 {
		obj.InstanceType = v
	}

	if v, ok := in["region"].(string); ok && len(v) > 0 {
		obj.Region = v
	}

	if v, ok := in["secret_key"].(string); ok && len(v) > 0 {
		obj.SecretKey = v
	}

	if v, ok := in["security_group_ids"].([]interface{}); ok && len(v) > 0 {
		obj.SecurityGroupIds = toArrayString(v)
	}

	if v, ok := in["extra_tags_all"].([]interface{}); ok && len(v) > 0 {
		obj.ExtraTagsAll = toArrayString(v)
	}

	if v, ok := in["extra_tags_instances"].([]interface{}); ok && len(v) > 0 {
		obj.ExtraTagsInstances = toArrayString(v)
	}

	if v, ok := in["root_disk_type"].(string); ok && len(v) > 0 {
		obj.RootDiskType = v
	}

	if v, ok := in["root_disk_size"].(int); ok {
		obj.RootDiskSize = v
	}

	if v, ok := in["root_disk_iops"].(int); ok {
		obj.RootDiskIops = v
	}

	return obj
}
