package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenClusterRKEConfigNodeDrainInput(in *managementClient.NodeDrainInput) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	obj["delete_local_data"] = in.DeleteLocalData
	obj["force"] = in.Force
	obj["grace_period"] = int(in.GracePeriod)
	if in.IgnoreDaemonSets != nil {
		obj["ignore_daemon_sets"] = *in.IgnoreDaemonSets
	}
	obj["timeout"] = int(in.Timeout)

	return []interface{}{obj}
}

func flattenClusterRKEConfigNodeUpgradeStrategy(in *managementClient.NodeUpgradeStrategy) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	if in.Drain != nil {
		obj["drain"] = *in.Drain
	}

	if in.DrainInput != nil {
		obj["drain_input"] = flattenClusterRKEConfigNodeDrainInput(in.DrainInput)
	}

	if len(in.MaxUnavailableControlplane) > 0 {
		obj["max_unavailable_controlplane"] = in.MaxUnavailableControlplane
	}

	if len(in.MaxUnavailableWorker) > 0 {
		obj["max_unavailable_worker"] = in.MaxUnavailableWorker
	}

	return []interface{}{obj}
}

func flattenClusterRKEConfigNodes(p []managementClient.RKEConfigNode) ([]interface{}, error) {
	out := []interface{}{}

	for _, in := range p {
		obj := make(map[string]interface{})

		if len(in.Address) > 0 {
			obj["address"] = in.Address
		}

		if len(in.DockerSocket) > 0 {
			obj["docker_socket"] = in.DockerSocket
		}

		if len(in.HostnameOverride) > 0 {
			obj["hostname_override"] = in.HostnameOverride
		}

		if len(in.InternalAddress) > 0 {
			obj["internal_address"] = in.InternalAddress
		}

		if len(in.Labels) > 0 {
			obj["labels"] = toMapInterface(in.Labels)
		}

		if len(in.NodeID) > 0 {
			obj["node_id"] = in.NodeID
		}

		if len(in.Port) > 0 {
			obj["port"] = in.Port
		}

		if len(in.Role) > 0 {
			obj["role"] = toArrayInterface(in.Role)
		}

		obj["ssh_agent_auth"] = in.SSHAgentAuth

		if len(in.SSHKey) > 0 {
			obj["ssh_key"] = in.SSHKey
		}

		if len(in.SSHKeyPath) > 0 {
			obj["ssh_key_path"] = in.SSHKeyPath
		}

		if len(in.User) > 0 {
			obj["user"] = in.User
		}

		out = append(out, obj)
	}

	return out, nil
}

// Expanders

func expandClusterRKEConfigNodeDrainInput(p []interface{}) *managementClient.NodeDrainInput {
	obj := &managementClient.NodeDrainInput{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["delete_local_data"].(bool); ok {
		obj.DeleteLocalData = v
	}

	if v, ok := in["force"].(bool); ok {
		obj.Force = v
	}

	if v, ok := in["grace_period"].(int); ok {
		obj.GracePeriod = int64(v)
	}

	if v, ok := in["ignore_daemon_sets"].(bool); ok {
		obj.IgnoreDaemonSets = &v
	}

	if v, ok := in["timeout"].(int); ok {
		obj.Timeout = int64(v)
	}

	return obj
}

func expandClusterRKEConfigNodeUpgradeStrategy(p []interface{}) *managementClient.NodeUpgradeStrategy {
	obj := &managementClient.NodeUpgradeStrategy{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["drain"].(bool); ok {
		obj.Drain = &v
	}

	if v, ok := in["drain_input"].([]interface{}); ok {
		obj.DrainInput = expandClusterRKEConfigNodeDrainInput(v)
	}

	if v, ok := in["max_unavailable_controlplane"].(string); ok && len(v) > 0 {
		obj.MaxUnavailableControlplane = v
	}

	if v, ok := in["max_unavailable_worker"].(string); ok && len(v) > 0 {
		obj.MaxUnavailableWorker = v
	}

	return obj
}

func expandClusterRKEConfigNodes(p []interface{}) ([]managementClient.RKEConfigNode, error) {
	out := []managementClient.RKEConfigNode{}
	if len(p) == 0 || p[0] == nil {
		return out, nil
	}

	for i := range p {
		in := p[i].(map[string]interface{})
		obj := managementClient.RKEConfigNode{}

		if v, ok := in["address"].(string); ok && len(v) > 0 {
			obj.Address = v
		}

		if v, ok := in["docker_socket"].(string); ok && len(v) > 0 {
			obj.DockerSocket = v
		}

		if v, ok := in["hostname_override"].(string); ok && len(v) > 0 {
			obj.HostnameOverride = v
		}

		if v, ok := in["internal_address"].(string); ok && len(v) > 0 {
			obj.InternalAddress = v
		}

		if v, ok := in["labels"].(map[string]interface{}); ok && len(v) > 0 {
			obj.Labels = toMapString(v)
		}

		if v, ok := in["node_id"].(string); ok && len(v) > 0 {
			obj.NodeID = v
		}

		if v, ok := in["port"].(string); ok && len(v) > 0 {
			obj.Port = v
		}

		if v, ok := in["role"].([]interface{}); ok && len(v) > 0 {
			obj.Role = toArrayString(v)
		}

		if v, ok := in["ssh_agent_auth"].(bool); ok {
			obj.SSHAgentAuth = v
		}

		if v, ok := in["ssh_key"].(string); ok && len(v) > 0 {
			obj.SSHKey = v
		}

		if v, ok := in["ssh_key_path"].(string); ok && len(v) > 0 {
			obj.SSHKeyPath = v
		}

		if v, ok := in["user"].(string); ok && len(v) > 0 {
			obj.User = v
		}

		out = append(out, obj)
	}

	return out, nil
}
