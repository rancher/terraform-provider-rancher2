package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenClusterEKSConfigV2NodeGroupsLaunchTemplate(in *managementClient.LaunchTemplate, p []interface{}) []interface{} {
	if in == nil {
		return nil
	}
	obj := map[string]interface{}{}
	if len(p) != 0 && p[0] != nil {
		obj = p[0].(map[string]interface{})
	}
	if len(in.Name) > 0 {
		obj["name"] = in.Name
	}
	if in.Version != nil {
		obj["version"] = int(*in.Version)
	}

	return []interface{}{obj}
}

func flattenClusterEKSConfigV2NodeGroups(input []managementClient.NodeGroup, p []interface{}) []interface{} {
	if input == nil {
		return nil
	}
	out := make([]interface{}, len(input))
	for i, in := range input {
		obj := map[string]interface{}{}
		if i < len(p) && p[i] != nil {
			obj = p[i].(map[string]interface{})
		}

		if len(in.NodegroupName) > 0 {
			obj["name"] = in.NodegroupName
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
		if len(in.ImageID) > 0 {
			obj["image_id"] = in.ImageID
		}
		if len(in.InstanceType) > 0 {
			obj["instance_type"] = in.InstanceType
		}
		if len(in.Labels) > 0 {
			obj["labels"] = toMapInterface(in.Labels)
		}
		if in.LaunchTemplate != nil {
			v, ok := obj["launch_template"].([]interface{})
			if !ok {
				v = []interface{}{}
			}
			obj["launch_template"] = flattenClusterEKSConfigV2NodeGroupsLaunchTemplate(in.LaunchTemplate, v)
		}
		if in.MaxSize != nil {
			obj["max_size"] = int(*in.MaxSize)
		}
		if in.MinSize != nil {
			obj["min_size"] = int(*in.MinSize)
		}
		if in.RequestSpotInstances != nil {
			obj["request_spot_instances"] = *in.RequestSpotInstances
		}
		if len(in.ResourceTags) > 0 {
			obj["resource_tags"] = toMapInterface(in.ResourceTags)
		}
		if len(in.SpotInstanceTypes) > 0 {
			obj["spot_instance_types"] = toArrayInterface(in.SpotInstanceTypes)
		}
		if len(in.Subnets) > 0 {
			obj["subnets"] = toArrayInterface(in.Subnets)
		}
		if len(in.Tags) > 0 {
			obj["tags"] = toMapInterface(in.Tags)
		}
		if len(in.UserData) > 0 {
			obj["user_data"] = in.UserData
		}
		if len(in.Version) > 0 {
			obj["version"] = in.Version
		}
		out[i] = obj
	}

	return out
}

func flattenClusterEKSConfigV2(in *managementClient.EKSClusterConfigSpec, p []interface{}) []interface{} {
	if in == nil {
		return nil
	}

	obj := map[string]interface{}{}
	if len(p) != 0 && p[0] != nil {
		obj = p[0].(map[string]interface{})
	}

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
		v, ok := obj["node_groups"].([]interface{})
		if !ok {
			v = []interface{}{}
		}
		obj["node_groups"] = flattenClusterEKSConfigV2NodeGroups(in.NodeGroups, v)
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

func expandClusterEKSConfigV2NodeGroupsLaunchTemplate(p []interface{}) *managementClient.LaunchTemplate {
	obj := &managementClient.LaunchTemplate{}
	if len(p) == 0 || p[0] == nil {
		return nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["name"].(string); ok && len(v) > 0 {
		obj.Name = v
	}
	if v, ok := in["version"].(int); ok {
		ver := int64(v)
		obj.Version = &ver
	}
	return obj
}

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
		if v, ok := in["image_id"].(string); ok {
			obj.ImageID = v
		}
		if v, ok := in["labels"].(map[string]interface{}); ok {
			obj.Labels = toMapString(v)
		}
		if v, ok := in["launch_template"].([]interface{}); ok && len(v) > 0 {
			obj.LaunchTemplate = expandClusterEKSConfigV2NodeGroupsLaunchTemplate(v)
		}
		if v, ok := in["max_size"].(int); ok {
			size := int64(v)
			obj.MaxSize = &size
		}
		if v, ok := in["min_size"].(int); ok {
			size := int64(v)
			obj.MinSize = &size
		}
		if v, ok := in["request_spot_instances"].(bool); ok {
			obj.RequestSpotInstances = &v
		}
		if v, ok := in["resource_tags"].(map[string]interface{}); ok {
			obj.ResourceTags = toMapString(v)
		}
		if v, ok := in["spot_instance_types"].([]interface{}); ok {
			obj.SpotInstanceTypes = toArrayString(v)
		}
		// setting objSubnets from subnet var or from tf argument
		if subnets != nil {
			obj.Subnets = subnets
		}
		if v, ok := in["subnets"].([]interface{}); ok {
			obj.Subnets = toArrayString(v)
		}
		if v, ok := in["tags"].(map[string]interface{}); ok {
			obj.Tags = toMapString(v)
		}
		if v, ok := in["user_data"].(string); ok {
			obj.UserData = v
		}
		if len(version) > 0 {
			obj.Version = version
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

// This fix is required due to managementClient.LaunchTemplate struct doesn't contains ID field
func fixClusterEKSConfigV2(p []interface{}, values map[string]interface{}) map[string]interface{} {
	if len(p) == 0 || p[0] == nil {
		return nil
	}

	in := p[0].(map[string]interface{})

	v, ok := in["node_groups"].([]interface{})
	v2, ok2 := values["nodeGroups"].([]interface{})
	if ok && ok2 {
		values["nodeGroups"] = fixClusterEKSConfigV2NodeGroups(v, v2)
	}

	return values
}

func fixClusterEKSConfigV2NodeGroups(p []interface{}, values []interface{}) []interface{} {
	if len(p) == 0 || p[0] == nil {
		return nil
	}

	for i := range p {
		in := p[i].(map[string]interface{})

		if v, ok := in["launch_template"].([]interface{}); ok {
			values[i].(map[string]interface{})["launchTemplate"] = fixClusterEKSConfigV2NodeGroupsLaunchTemplate(v)
		}
	}

	return values
}

func fixClusterEKSConfigV2NodeGroupsLaunchTemplate(p []interface{}) map[string]interface{} {
	if len(p) == 0 || p[0] == nil {
		return nil
	}
	in := p[0].(map[string]interface{})

	obj := map[string]interface{}{}

	if v, ok := in["id"].(string); ok && len(v) > 0 {
		obj["id"] = v
	}
	if v, ok := in["name"].(string); ok && len(v) > 0 {
		obj["name"] = v
	}
	if v, ok := in["version"].(int); ok {
		ver := int64(v)
		obj["version"] = &ver
	}

	return obj
}
