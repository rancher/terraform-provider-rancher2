package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

//Schemas

func projectResourceQuotaLimitFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"config_maps": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"limits_cpu": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"limits_memory": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"persistent_volume_claims": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"pods": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"replication_controllers": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"requests_cpu": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"requests_memory": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"requests_storage": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"secrets": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"services_load_balancers": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"services_node_ports": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}

	return s
}

// Flatteners

func flattenProjectResourceQuotaLimit(in *managementClient.ResourceQuotaLimit) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
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

	return []interface{}{obj}, nil
}

// Expanders

func expandProjectResourceQuotaLimit(p []interface{}) (*managementClient.ResourceQuotaLimit, error) {
	obj := &managementClient.ResourceQuotaLimit{}

	if len(p) == 0 || p[0] == nil {
		return obj, nil
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

	return obj, nil
}
