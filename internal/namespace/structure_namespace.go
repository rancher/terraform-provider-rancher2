package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	clusterClient "github.com/rancher/rancher/pkg/client/generated/cluster/v3"
)

// Flatteners

func flattenNamespaceContainerResourceLimit(in *clusterClient.ContainerResourceLimit) []interface{} {
	obj := make(map[string]interface{})
	empty := clusterClient.ContainerResourceLimit{}
	if in == nil || *in == empty {
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

func flattenNamespaceResourceQuotaLimit(in *clusterClient.ResourceQuotaLimit) []interface{} {
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

func flattenNamespaceResourceQuota(in *clusterClient.NamespaceResourceQuota) []interface{} {
	obj := make(map[string]interface{})
	empty := clusterClient.NamespaceResourceQuota{}
	if in == nil || *in == empty {
		return []interface{}{}
	}

	if in.Limit != nil {
		limit := flattenNamespaceResourceQuotaLimit(in.Limit)
		obj["limit"] = limit
	}

	return []interface{}{obj}
}

func flattenNamespace(d *schema.ResourceData, in *clusterClient.Namespace) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)

	if len(in.ProjectID) > 0 {
		d.Set("project_id", in.ProjectID)
	}

	d.Set("name", in.Name)
	d.Set("description", in.Description)

	containerLimit := flattenNamespaceContainerResourceLimit(in.ContainerDefaultResourceLimit)
	err := d.Set("container_resource_limit", containerLimit)
	if err != nil {
		return err
	}

	resourceQuota := flattenNamespaceResourceQuota(in.ResourceQuota)
	err = d.Set("resource_quota", resourceQuota)
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

func expandNamespaceContainerResourceLimit(p []interface{}) *clusterClient.ContainerResourceLimit {
	obj := &clusterClient.ContainerResourceLimit{}
	if len(p) == 0 || p[0] == nil {
		return obj
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

func expandNamespaceResourceQuotaLimit(p []interface{}) *clusterClient.ResourceQuotaLimit {
	obj := &clusterClient.ResourceQuotaLimit{}

	if len(p) == 0 || p[0] == nil {
		return obj
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

func expandNamespaceResourceQuota(p []interface{}) *clusterClient.NamespaceResourceQuota {
	obj := &clusterClient.NamespaceResourceQuota{}

	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["limit"].([]interface{}); ok && len(v) > 0 {
		limit := expandNamespaceResourceQuotaLimit(v)
		obj.Limit = limit
	}

	return obj
}

func expandNamespace(in *schema.ResourceData) *clusterClient.Namespace {
	obj := &clusterClient.Namespace{}
	if in == nil {
		return nil
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	_, projectID := splitProjectID(in.Get("project_id").(string))
	obj.ProjectID = projectID
	obj.Name = in.Get("name").(string)
	obj.Description = in.Get("description").(string)

	containerLimit := expandNamespaceContainerResourceLimit(in.Get("container_resource_limit").([]interface{}))
	obj.ContainerDefaultResourceLimit = containerLimit

	if v, ok := in.Get("resource_quota").([]interface{}); ok && len(v) > 0 {
		resourceQuota := expandNamespaceResourceQuota(v)
		obj.ResourceQuota = resourceQuota
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj
}
