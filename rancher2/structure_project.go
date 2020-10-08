package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenProjectContainerResourceLimit(in *managementClient.ContainerResourceLimit) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	if len(in.LimitsCPU) > 0 {
		obj["limits_cpu"] = in.LimitsCPU
	}
	if len(in.LimitsMemory) > 0 {
		obj["limits_memory"] = in.LimitsMemory
	}
	if len(in.RequestsCPU) > 0 {
		obj["requests_cpu"] = in.RequestsCPU
	}
	if len(in.RequestsMemory) > 0 {
		obj["requests_memory"] = in.RequestsMemory
	}

	return []interface{}{obj}
}

func flattenProjectResourceQuotaLimit(in *managementClient.ResourceQuotaLimit) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	if len(in.ConfigMaps) > 0 {
		obj["config_maps"] = in.ConfigMaps
	}

	if len(in.LimitsCPU) > 0 {
		obj["limits_cpu"] = in.LimitsCPU
	}

	if len(in.LimitsMemory) > 0 {
		obj["limits_memory"] = in.LimitsMemory
	}

	if len(in.PersistentVolumeClaims) > 0 {
		obj["persistent_volume_claims"] = in.PersistentVolumeClaims
	}

	if len(in.Pods) > 0 {
		obj["pods"] = in.Pods
	}

	if len(in.ReplicationControllers) > 0 {
		obj["replication_controllers"] = in.ReplicationControllers
	}

	if len(in.RequestsCPU) > 0 {
		obj["requests_cpu"] = in.RequestsCPU
	}

	if len(in.RequestsMemory) > 0 {
		obj["requests_memory"] = in.RequestsMemory
	}

	if len(in.RequestsStorage) > 0 {
		obj["requests_storage"] = in.RequestsStorage
	}

	if len(in.Secrets) > 0 {
		obj["secrets"] = in.Secrets
	}

	if len(in.Services) > 0 {
		obj["services"] = in.Services
	}

	if len(in.ServicesLoadBalancers) > 0 {
		obj["services_load_balancers"] = in.ServicesLoadBalancers
	}

	if len(in.ServicesNodePorts) > 0 {
		obj["services_node_ports"] = in.ServicesNodePorts
	}

	return []interface{}{obj}
}

func flattenProjectResourceQuota(pQuota *managementClient.ProjectResourceQuota, nsQuota *managementClient.NamespaceResourceQuota) []interface{} {
	obj := make(map[string]interface{})
	if pQuota == nil || nsQuota == nil {
		return []interface{}{}
	}

	if pQuota.Limit != nil {
		limit := flattenProjectResourceQuotaLimit(pQuota.Limit)
		obj["project_limit"] = limit
	}

	if nsQuota.Limit != nil {
		limit := flattenProjectResourceQuotaLimit(nsQuota.Limit)
		obj["namespace_default_limit"] = limit
	}

	return []interface{}{obj}
}

func flattenProject(d *schema.ResourceData, in *managementClient.Project, monitoringInput *managementClient.MonitoringInput) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)
	d.Set("cluster_id", in.ClusterID)
	d.Set("name", in.Name)
	d.Set("description", in.Description)
	d.Set("enable_project_monitoring", in.EnableProjectMonitoring)

	if in.ContainerDefaultResourceLimit != nil {
		containerLimit := flattenProjectContainerResourceLimit(in.ContainerDefaultResourceLimit)
		err := d.Set("container_resource_limit", containerLimit)
		if err != nil {
			return err
		}
	}

	d.Set("pod_security_policy_template_id", in.PodSecurityPolicyTemplateName)

	if in.ResourceQuota != nil && in.NamespaceDefaultResourceQuota != nil {
		resourceQuota := flattenProjectResourceQuota(in.ResourceQuota, in.NamespaceDefaultResourceQuota)
		err := d.Set("resource_quota", resourceQuota)
		if err != nil {
			return err
		}
	}

	err := d.Set("project_monitoring_input", flattenMonitoringInput(monitoringInput))
	if err != nil {
		return err
	}

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

func expandProjectContainerResourceLimit(p []interface{}) *managementClient.ContainerResourceLimit {
	obj := &managementClient.ContainerResourceLimit{}
	if len(p) == 0 || p[0] == nil {
		return nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["limits_cpu"].(string); ok && len(v) > 0 {
		obj.LimitsCPU = v
	}
	if v, ok := in["limits_memory"].(string); ok && len(v) > 0 {
		obj.LimitsMemory = v
	}
	if v, ok := in["requests_cpu"].(string); ok && len(v) > 0 {
		obj.RequestsCPU = v
	}
	if v, ok := in["requests_memory"].(string); ok && len(v) > 0 {
		obj.RequestsMemory = v
	}

	return obj
}

func expandProjectResourceQuotaLimit(p []interface{}) *managementClient.ResourceQuotaLimit {
	obj := &managementClient.ResourceQuotaLimit{}

	if len(p) == 0 || p[0] == nil {
		return nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["config_maps"].(string); ok && len(v) > 0 {
		obj.ConfigMaps = v
	}

	if v, ok := in["limits_cpu"].(string); ok && len(v) > 0 {
		obj.LimitsCPU = v
	}

	if v, ok := in["limits_memory"].(string); ok && len(v) > 0 {
		obj.LimitsMemory = v
	}

	if v, ok := in["persistent_volume_claims"].(string); ok && len(v) > 0 {
		obj.PersistentVolumeClaims = v
	}

	if v, ok := in["pods"].(string); ok && len(v) > 0 {
		obj.Pods = v
	}

	if v, ok := in["replication_controllers"].(string); ok && len(v) > 0 {
		obj.ReplicationControllers = v
	}

	if v, ok := in["requests_cpu"].(string); ok && len(v) > 0 {
		obj.RequestsCPU = v
	}

	if v, ok := in["requests_memory"].(string); ok && len(v) > 0 {
		obj.RequestsMemory = v
	}

	if v, ok := in["requests_storage"].(string); ok && len(v) > 0 {
		obj.RequestsStorage = v
	}

	if v, ok := in["secrets"].(string); ok && len(v) > 0 {
		obj.Secrets = v
	}

	if v, ok := in["services"].(string); ok && len(v) > 0 {
		obj.Services = v
	}

	if v, ok := in["services_load_balancers"].(string); ok && len(v) > 0 {
		obj.ServicesLoadBalancers = v
	}

	if v, ok := in["services_node_ports"].(string); ok && len(v) > 0 {
		obj.ServicesNodePorts = v
	}

	return obj
}

func expandProjectResourceQuota(p []interface{}) (*managementClient.ProjectResourceQuota, *managementClient.NamespaceResourceQuota) {
	pQuota := &managementClient.ProjectResourceQuota{}
	nsQuota := &managementClient.NamespaceResourceQuota{}

	if len(p) == 0 || p[0] == nil {
		return nil, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["project_limit"].([]interface{}); ok && len(v) > 0 {
		pLimit := expandProjectResourceQuotaLimit(v)
		pQuota.Limit = pLimit
	}

	if v, ok := in["namespace_default_limit"].([]interface{}); ok && len(v) > 0 {
		nsLimit := expandProjectResourceQuotaLimit(v)
		nsQuota.Limit = nsLimit
	}

	return pQuota, nsQuota
}

func expandProject(in *schema.ResourceData) *managementClient.Project {
	obj := &managementClient.Project{}
	if in == nil {
		return nil
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	obj.ClusterID = in.Get("cluster_id").(string)
	obj.Name = in.Get("name").(string)
	obj.Description = in.Get("description").(string)

	if v, ok := in.Get("container_resource_limit").([]interface{}); ok && len(v) > 0 {
		containerLimit := expandProjectContainerResourceLimit(v)
		obj.ContainerDefaultResourceLimit = containerLimit
	}

	if v, ok := in.Get("enable_project_monitoring").(bool); ok {
		obj.EnableProjectMonitoring = v
	}

	obj.PodSecurityPolicyTemplateName = in.Get("pod_security_policy_template_id").(string)

	if v, ok := in.Get("resource_quota").([]interface{}); ok && len(v) > 0 {
		resourceQuota, nsResourceQuota := expandProjectResourceQuota(v)
		obj.ResourceQuota = resourceQuota
		obj.NamespaceDefaultResourceQuota = nsResourceQuota
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj
}
