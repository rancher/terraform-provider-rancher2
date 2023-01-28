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
	if in.ID != nil && len(*in.ID) > 0 {
		obj["id"] = *in.ID
	}
	if in.Name != nil && len(*in.Name) > 0 {
		obj["name"] = *in.Name
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

		if in.NodegroupName != nil && len(*in.NodegroupName) > 0 {
			obj["name"] = *in.NodegroupName
		}
		if in.DesiredSize != nil {
			obj["desired_size"] = int(*in.DesiredSize)
		}
		if in.DiskSize != nil {
			obj["disk_size"] = int(*in.DiskSize)
		}
		if in.Ec2SshKey != nil && len(*in.Ec2SshKey) > 0 {
			obj["ec2_ssh_key"] = *in.Ec2SshKey
		}
		if in.Gpu != nil {
			obj["gpu"] = *in.Gpu
		}
		if in.ImageID != nil && len(*in.ImageID) > 0 {
			obj["image_id"] = *in.ImageID
		}
		if in.InstanceType != nil && len(*in.InstanceType) > 0 {
			obj["instance_type"] = *in.InstanceType
		}
		if in.Labels != nil && len(in.Labels) > 0 {
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
		if in.NodeRole != nil && len(*in.NodeRole) > 0 {
			obj["node_role"] = *in.NodeRole
		}
		if in.RequestSpotInstances != nil {
			obj["request_spot_instances"] = *in.RequestSpotInstances
		}
		if in.ResourceTags != nil && len(in.ResourceTags) > 0 {
			obj["resource_tags"] = toMapInterface(in.ResourceTags)
		}
		if in.SpotInstanceTypes != nil && len(in.SpotInstanceTypes) > 0 {
			obj["spot_instance_types"] = toArrayInterfaceSorted(in.SpotInstanceTypes)
		}
		if in.Subnets != nil && len(in.Subnets) > 0 {
			obj["subnets"] = toArrayInterfaceSorted(in.Subnets)
		}
		if in.Tags != nil && len(in.Tags) > 0 {
			obj["tags"] = toMapInterface(in.Tags)
		}
		if in.UserData != nil && len(*in.UserData) > 0 {
			obj["user_data"] = *in.UserData
		}
		if in.Version != nil && len(*in.Version) > 0 {
			obj["version"] = *in.Version
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
	if in.KubernetesVersion != nil && len(*in.KubernetesVersion) > 0 {
		obj["kubernetes_version"] = *in.KubernetesVersion
	}
	if in.NodeGroups != nil && len(in.NodeGroups) > 0 {
		v, ok := obj["node_groups"].([]interface{})
		if !ok {
			v = []interface{}{}
		}
		obj["node_groups"] = flattenClusterEKSConfigV2NodeGroups(in.NodeGroups, v)
	}
	obj["imported"] = in.Imported
	if in.KmsKey != nil && len(*in.KmsKey) > 0 {
		obj["kms_key"] = *in.KmsKey
	}
	if in.LoggingTypes != nil && len(in.LoggingTypes) > 0 {
		obj["logging_types"] = toArrayInterfaceSorted(in.LoggingTypes)
	}
	if in.PrivateAccess != nil {
		obj["private_access"] = *in.PrivateAccess
	}
	if in.PublicAccess != nil {
		obj["public_access"] = *in.PublicAccess
	}
	if in.PublicAccessSources != nil && len(in.PublicAccessSources) > 0 {
		obj["public_access_sources"] = toArrayInterfaceSorted(in.PublicAccessSources)
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
	if in.SecurityGroups != nil && len(in.SecurityGroups) > 0 {
		obj["security_groups"] = toArrayInterfaceSorted(in.SecurityGroups)
	}
	if in.ServiceRole != nil && len(*in.ServiceRole) > 0 {
		obj["service_role"] = *in.ServiceRole
	}
	if in.Subnets != nil && len(in.Subnets) > 0 {
		obj["subnets"] = toArrayInterfaceSorted(in.Subnets)
	}
	if in.Tags != nil && len(in.Tags) > 0 {
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

	if v, ok := in["id"].(string); ok && len(v) > 0 {
		obj.ID = &v
	}
	if v, ok := in["name"].(string); ok && len(v) > 0 {
		obj.Name = &v
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

		if v, ok := in["name"].(string); ok {
			obj.NodegroupName = &v
		}
		if v, ok := in["instance_type"].(string); ok {
			obj.InstanceType = &v
		}
		if v, ok := in["desired_size"].(int); ok {
			size := int64(v)
			obj.DesiredSize = &size
		}
		if v, ok := in["disk_size"].(int); ok {
			size := int64(v)
			obj.DiskSize = &size
		}
		if v, ok := in["ec2_ssh_key"].(string); ok {
			obj.Ec2SshKey = &v
		}
		if v, ok := in["gpu"].(bool); ok {
			obj.Gpu = &v
		}
		if v, ok := in["image_id"].(string); ok {
			obj.ImageID = &v
		}
		if v, ok := in["labels"].(map[string]interface{}); ok {
			labels := toMapString(v)
			obj.Labels = labels
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
		if v, ok := in["node_role"].(string); ok {
			obj.NodeRole = &v
		}
		if v, ok := in["request_spot_instances"].(bool); ok {
			obj.RequestSpotInstances = &v
		}
		if v, ok := in["resource_tags"].(map[string]interface{}); ok {
			resourceTags := toMapString(v)
			obj.ResourceTags = resourceTags
		}
		if v, ok := in["spot_instance_types"].([]interface{}); ok {
			spotInstanceTypes := toArrayStringSorted(v)
			obj.SpotInstanceTypes = spotInstanceTypes
		}
		// setting objSubnets from subnet var or from tf argument
		if subnets != nil {
			obj.Subnets = subnets
		}
		if v, ok := in["subnets"].([]interface{}); ok {
			nets := toArrayStringSorted(v)
			obj.Subnets = nets
		}
		if v, ok := in["tags"].(map[string]interface{}); ok {
			tags := toMapString(v)
			obj.Tags = tags
		}
		if v, ok := in["user_data"].(string); ok {
			obj.UserData = &v
		}
		if len(version) > 0 {
			obj.Version = &version
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
	k8sVersion := ""
	if v, ok := in["kubernetes_version"].(string); ok && len(v) > 0 {
		k8sVersion = v
		obj.KubernetesVersion = &k8sVersion
	}
	subnets := []string{}
	if v, ok := in["subnets"].([]interface{}); ok {
		subnets = toArrayStringSorted(v)
		obj.Subnets = subnets
	}
	if v, ok := in["node_groups"].([]interface{}); ok {
		nodeGroups := expandClusterEKSConfigV2NodeGroups(v, subnets, k8sVersion)
		obj.NodeGroups = nodeGroups
	}
	if v, ok := in["imported"].(bool); ok {
		obj.Imported = v
	}
	if v, ok := in["kms_key"].(string); ok && len(v) > 0 {
		obj.KmsKey = newString(v)
	}
	if v, ok := in["logging_types"].([]interface{}); ok {
		loggingTypes := toArrayStringSorted(v)
		obj.LoggingTypes = loggingTypes
	}
	if v, ok := in["private_access"].(bool); ok && !obj.Imported {
		obj.PrivateAccess = &v
	}
	if v, ok := in["public_access"].(bool); ok && !obj.Imported {
		obj.PublicAccess = &v
	}
	if v, ok := in["public_access_sources"].([]interface{}); ok {
		publicAccessSources := toArrayStringSorted(v)
		obj.PublicAccessSources = publicAccessSources
	}
	if v, ok := in["region"].(string); ok {
		obj.Region = v
	}
	if v, ok := in["secrets_encryption"].(bool); ok && !obj.Imported {
		obj.SecretsEncryption = &v
	}
	if v, ok := in["security_groups"].([]interface{}); ok {
		securityGroups := toArrayStringSorted(v)
		obj.SecurityGroups = securityGroups
	}
	if v, ok := in["service_role"].(string); ok {
		obj.ServiceRole = newString(v)
	}
	obj.Tags = map[string]string{}
	if v, ok := in["tags"].(map[string]interface{}); ok {
		tags := toMapString(v)
		obj.Tags = tags
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
