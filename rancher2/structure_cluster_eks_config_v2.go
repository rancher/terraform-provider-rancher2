package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenClusterEKSConfigV2NodeGroups(input []managementClient.NodeGroup) []interface{} {
	if input == nil {
		return nil
	}
	out := make([]interface{}, len(input))
	for i, in := range input {
		obj := make(map[string]interface{})

		if len(in.NodegroupName) > 0 {
			obj["name"] = in.NodegroupName
		}
		if len(in.InstanceType) > 0 {
			obj["instance_type"] = in.InstanceType
		}
		if in.DesiredSize != nil {
			obj["desired_size"] = int(*in.DesiredSize)
		}
		if in.DiskSize != nil {
			obj["disk_size"] = int(*in.DiskSize)
		}
		if len(in.Ec2SshKey) > 0 {
			obj["ec2_ssh_key"] = in.Ec2SshKey
		}
		if in.Gpu != nil {
			obj["gpu"] = *in.Gpu
		}
		if len(in.Labels) > 0 {
			obj["labels"] = toMapInterface(in.Labels)
		}
		if len(in.Tags) > 0 {
			obj["tags"] = toMapInterface(in.Tags)
		}
		if in.MaxSize != nil {
			obj["max_size"] = int(*in.MaxSize)
		}
		if in.MinSize != nil {
			obj["min_size"] = int(*in.MinSize)
		}
		out[i] = obj
	}

	return out
}

func flattenClusterEKSConfigV2(in *managementClient.EKSClusterConfigSpec) []interface{} {
	if in == nil {
		return nil
	}
	obj := make(map[string]interface{})
	if len(in.AmazonCredentialSecret) > 0 {
		obj["cloud_credential_id"] = in.AmazonCredentialSecret
	}
	if len(in.DisplayName) > 0 {
		obj["name"] = in.DisplayName
	}
	if len(in.KubernetesVersion) > 0 {
		obj["kubernetes_version"] = in.KubernetesVersion
	}
	if len(in.NodeGroups) > 0 {
		obj["node_groups"] = flattenClusterEKSConfigV2NodeGroups(in.NodeGroups)
	}
	obj["imported"] = in.Imported
	if len(in.KmsKey) > 0 {
		obj["kms_key"] = in.KmsKey
	}
	if len(in.LoggingTypes) > 0 {
		obj["logging_types"] = toArrayInterface(in.LoggingTypes)
	}
	if in.PrivateAccess != nil {
		obj["private_access"] = *in.PrivateAccess
	}
	if in.PublicAccess != nil {
		obj["public_access"] = *in.PublicAccess
	}
	if len(in.PublicAccessSources) > 0 {
		obj["public_access_sources"] = toArrayInterface(in.PublicAccessSources)
	}
	if in.SecretsEncryption != nil {
		obj["secrets_encryption"] = *in.SecretsEncryption
	}
	if len(in.Region) > 0 {
		obj["region"] = in.Region
	}
	if in.SecretsEncryption != nil {
		obj["secrets_encryption"] = *in.SecretsEncryption
	}
	if len(in.SecurityGroups) > 0 {
		obj["security_groups"] = toArrayInterface(in.SecurityGroups)
	}
	if len(in.ServiceRole) > 0 {
		obj["service_role"] = in.ServiceRole
	}
	if len(in.Subnets) > 0 {
		obj["subnets"] = toArrayInterface(in.Subnets)
	}
	if len(in.Tags) > 0 {
		obj["tags"] = toMapInterface(in.Tags)
	}

	return []interface{}{obj}
}

// Expanders

func expandClusterEKSConfigV2NodeGroups(p []interface{}, subnets []string, version string) []managementClient.NodeGroup {
	if p == nil || len(p) == 0 {
		return []managementClient.NodeGroup{}
	}
	out := make([]managementClient.NodeGroup, len(p))
	for i := range p {
		in := p[i].(map[string]interface{})
		obj := managementClient.NodeGroup{}

		obj.NodegroupName = in["name"].(string)
		obj.InstanceType = in["instance_type"].(string)

		if v, ok := in["desired_size"].(int); ok {
			size := int64(v)
			obj.DesiredSize = &size
		}
		if v, ok := in["disk_size"].(int); ok {
			size := int64(v)
			obj.DiskSize = &size
		}
		if v, ok := in["ec2_ssh_key"].(string); ok {
			obj.Ec2SshKey = v
		}
		if v, ok := in["gpu"].(bool); ok {
			obj.Gpu = &v
		}
		if v, ok := in["labels"].(map[string]interface{}); ok {
			obj.Labels = toMapString(v)
		}
		if v, ok := in["max_size"].(int); ok {
			size := int64(v)
			obj.MaxSize = &size
		}
		if v, ok := in["min_size"].(int); ok {
			size := int64(v)
			obj.MinSize = &size
		}
		if subnets != nil {
			obj.Subnets = subnets
		}
		if len(version) > 0 {
			obj.Version = version
		}
		if v, ok := in["tags"].(map[string]interface{}); ok {
			obj.Tags = toMapString(v)
		}
		out[i] = obj
	}

	return out
}

func expandClusterEKSConfigV2(p []interface{}) *managementClient.EKSClusterConfigSpec {
	obj := &managementClient.EKSClusterConfigSpec{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	obj.AmazonCredentialSecret = in["cloud_credential_id"].(string)
	obj.DisplayName = in["name"].(string)
	obj.KubernetesVersion = in["kubernetes_version"].(string)

	if v, ok := in["subnets"].([]interface{}); ok {
		obj.Subnets = toArrayString(v)
	}
	if v, ok := in["node_groups"].([]interface{}); ok {
		obj.NodeGroups = expandClusterEKSConfigV2NodeGroups(v, obj.Subnets, obj.KubernetesVersion)
	}
	if v, ok := in["imported"].(bool); ok {
		obj.Imported = v
	}
	if v, ok := in["kms_key"].(string); ok && len(v) > 0 {
		obj.KmsKey = v
	}
	if v, ok := in["logging_types"].([]interface{}); ok {
		obj.LoggingTypes = toArrayString(v)
	}
	if v, ok := in["private_access"].(bool); ok {
		obj.PrivateAccess = &v
	}
	if v, ok := in["public_access"].(bool); ok {
		obj.PublicAccess = &v
	}
	if v, ok := in["public_access_sources"].([]interface{}); ok {
		obj.PublicAccessSources = toArrayString(v)
	}
	if v, ok := in["region"].(string); ok {
		obj.Region = v
	}
	if v, ok := in["secrets_encryption"].(bool); ok {
		obj.SecretsEncryption = &v
	}
	if v, ok := in["security_groups"].([]interface{}); ok {
		obj.SecurityGroups = toArrayString(v)
	}
	if v, ok := in["service_role"].(string); ok {
		obj.ServiceRole = v
	}
	obj.Tags = map[string]string{}
	if v, ok := in["tags"].(map[string]interface{}); ok {
		obj.Tags = toMapString(v)
	}

	return obj
}
