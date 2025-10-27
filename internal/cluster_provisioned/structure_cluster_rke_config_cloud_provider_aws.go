package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenClusterRKEConfigCloudProviderAwsGlobal(in *managementClient.GlobalAwsOpts) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	obj["disable_security_group_ingress"] = in.DisableSecurityGroupIngress
	obj["disable_strict_zone_check"] = in.DisableStrictZoneCheck

	if len(in.ElbSecurityGroup) > 0 {
		obj["elb_security_group"] = in.ElbSecurityGroup
	}

	if len(in.KubernetesClusterID) > 0 {
		obj["kubernetes_cluster_id"] = in.KubernetesClusterID
	}

	if len(in.KubernetesClusterTag) > 0 {
		obj["kubernetes_cluster_tag"] = in.KubernetesClusterTag
	}

	if len(in.RoleARN) > 0 {
		obj["role_arn"] = in.RoleARN
	}

	if len(in.RouteTableID) > 0 {
		obj["route_table_id"] = in.RouteTableID
	}

	if len(in.SubnetID) > 0 {
		obj["subnet_id"] = in.SubnetID
	}

	if len(in.VPC) > 0 {
		obj["vpc"] = in.VPC
	}

	if len(in.Zone) > 0 {
		obj["zone"] = in.Zone
	}

	return []interface{}{obj}
}

func flattenClusterRKEConfigCloudProviderAwsServiceOverride(in map[string]managementClient.ServiceOverride) []interface{} {
	if len(in) == 0 {
		return []interface{}{}
	}

	out := make([]interface{}, len(in))
	i := 0
	for key := range in {
		obj := make(map[string]interface{})
		if len(in[key].Region) > 0 {
			obj["region"] = in[key].Region
		}

		if len(in[key].Service) > 0 {
			obj["service"] = in[key].Service
		}

		if len(in[key].SigningMethod) > 0 {
			obj["signing_method"] = in[key].SigningMethod
		}

		if len(in[key].SigningName) > 0 {
			obj["signing_name"] = in[key].SigningName
		}

		if len(in[key].SigningRegion) > 0 {
			obj["signing_region"] = in[key].SigningRegion
		}

		if len(in[key].URL) > 0 {
			obj["url"] = in[key].URL
		}
		out[i] = obj
		i++
	}

	return out
}

func flattenClusterRKEConfigCloudProviderAws(in *managementClient.AWSCloudProvider) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if in.Global != nil {
		obj["global"] = flattenClusterRKEConfigCloudProviderAwsGlobal(in.Global)
	}

	if len(in.ServiceOverride) > 0 {
		obj["service_override"] = flattenClusterRKEConfigCloudProviderAwsServiceOverride(in.ServiceOverride)
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandClusterRKEConfigCloudProviderAwsGlobal(p []interface{}) *managementClient.GlobalAwsOpts {
	obj := &managementClient.GlobalAwsOpts{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["disable_security_group_ingress"].(bool); ok {
		obj.DisableSecurityGroupIngress = v
	}

	if v, ok := in["disable_strict_zone_check"].(bool); ok {
		obj.DisableStrictZoneCheck = v
	}

	if v, ok := in["elb_security_group"].(string); ok && len(v) > 0 {
		obj.ElbSecurityGroup = v
	}

	if v, ok := in["kubernetes_cluster_id"].(string); ok && len(v) > 0 {
		obj.KubernetesClusterID = v
	}

	if v, ok := in["kubernetes_cluster_tag"].(string); ok && len(v) > 0 {
		obj.KubernetesClusterTag = v
	}

	if v, ok := in["role_arn"].(string); ok && len(v) > 0 {
		obj.RoleARN = v
	}

	if v, ok := in["route_table_id"].(string); ok && len(v) > 0 {
		obj.RouteTableID = v
	}

	if v, ok := in["subnet_id"].(string); ok && len(v) > 0 {
		obj.SubnetID = v
	}

	if v, ok := in["vpc"].(string); ok && len(v) > 0 {
		obj.VPC = v
	}

	if v, ok := in["zone"].(string); ok && len(v) > 0 {
		obj.Zone = v
	}

	return obj
}

func expandClusterRKEConfigCloudProviderAwsServiceOverride(p []interface{}) map[string]managementClient.ServiceOverride {
	if len(p) == 0 || p[0] == nil {
		return map[string]managementClient.ServiceOverride{}
	}

	obj := make(map[string]managementClient.ServiceOverride)

	for i := range p {
		in := p[i].(map[string]interface{})
		aux := managementClient.ServiceOverride{}
		key := in["service"].(string)

		if v, ok := in["region"].(string); ok && len(v) > 0 {
			aux.Region = v
		}

		if v, ok := in["service"].(string); ok && len(v) > 0 {
			aux.Service = v
		}

		if v, ok := in["signing_method"].(string); ok && len(v) > 0 {
			aux.SigningMethod = v
		}

		if v, ok := in["signing_name"].(string); ok && len(v) > 0 {
			aux.SigningName = v
		}

		if v, ok := in["signing_region"].(string); ok && len(v) > 0 {
			aux.SigningRegion = v
		}

		if v, ok := in["url"].(string); ok && len(v) > 0 {
			aux.URL = v
		}
		obj[key] = aux
	}
	return obj
}

func expandClusterRKEConfigCloudProviderAws(p []interface{}) (*managementClient.AWSCloudProvider, error) {
	obj := &managementClient.AWSCloudProvider{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["global"].([]interface{}); ok && len(v) > 0 {
		obj.Global = expandClusterRKEConfigCloudProviderAwsGlobal(v)
	}

	if v, ok := in["service_override"].([]interface{}); ok && len(v) > 0 {
		obj.ServiceOverride = expandClusterRKEConfigCloudProviderAwsServiceOverride(v)
	}

	return obj, nil
}
