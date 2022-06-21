package rancher2

import (
	"time"

	provisionv1 "github.com/rancher/rancher/pkg/apis/provisioning.cattle.io/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// Flatteners

func flattenClusterV2RKEConfigMachinePoolMachineConfig(in *corev1.ObjectReference) []interface{} {
	if in == nil {
		return nil
	}

	obj := make(map[string]interface{})

	obj["kind"] = in.Kind
	obj["name"] = in.Name

	return []interface{}{obj}
}

func flattenClusterV2RKEConfigMachinePoolRollingUpdate(in *provisionv1.RKEMachinePoolRollingUpdate) []interface{} {
	if in == nil {
		return nil
	}

	obj := make(map[string]interface{})

	if in.MaxSurge != nil {
		obj["max_surge"] = in.MaxSurge.String()
	}
	if in.MaxUnavailable != nil {
		obj["max_unavailable"] = in.MaxUnavailable.String()
	}

	return []interface{}{obj}
}

func flattenClusterV2RKEConfigMachinePools(p []provisionv1.RKEMachinePool) []interface{} {
	if p == nil {
		return nil
	}
	out := make([]interface{}, len(p))
	for i, in := range p {
		obj := map[string]interface{}{}

		obj["name"] = in.Name
		if len(in.CloudCredentialSecretName) > 0 {
			obj["cloud_credential_secret_name"] = in.CloudCredentialSecretName
		}
		if in.NodeConfig != nil {
			obj["machine_config"] = flattenClusterV2RKEConfigMachinePoolMachineConfig(in.NodeConfig)
		}
		obj["control_plane_role"] = in.ControlPlaneRole
		obj["etcd_role"] = in.EtcdRole
		obj["drain_before_delete"] = in.DrainBeforeDelete

		if len(in.MachineDeploymentAnnotations) > 0 {
			obj["annotations"] = toMapInterface(in.MachineDeploymentAnnotations)
		}
		if len(in.MachineDeploymentLabels) > 0 {
			obj["labels"] = toMapInterface(in.MachineDeploymentLabels)
		}
		if len(in.Labels) > 0 {
			obj["machine_labels"] = toMapInterface(in.Labels)
		}
		obj["paused"] = in.Paused
		if in.Quantity != nil {
			obj["quantity"] = int(*in.Quantity)
		}
		if in.RollingUpdate != nil {
			obj["rolling_update"] = flattenClusterV2RKEConfigMachinePoolRollingUpdate(in.RollingUpdate)
		}
		if len(in.Taints) > 0 {
			obj["taints"] = flattenTaintsV2(in.Taints)
		}
		obj["worker_role"] = in.WorkerRole
		out[i] = obj

		if in.NodeStartupTimeout != nil {
			obj["node_startup_timeout_seconds"] = int(in.NodeStartupTimeout.Seconds())
		}
		if in.UnhealthyNodeTimeout != nil {
			obj["unhealthy_node_timeout_seconds"] = int(in.UnhealthyNodeTimeout.Seconds())
		}
		if in.DrainBeforeDeleteTimeout != nil {
			obj["node_drain_timeout"] = int(in.DrainBeforeDeleteTimeout.Seconds())
		}
		if in.MaxUnhealthy != nil {
			obj["max_unhealthy"] = *in.MaxUnhealthy
		}
		if in.UnhealthyRange != nil {
			obj["unhealthy_range"] = *in.UnhealthyRange
		}
	}

	return out
}

// Expanders

func expandClusterV2RKEConfigMachinePoolMachineConfig(p []interface{}) *corev1.ObjectReference {
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil
	}

	obj := &corev1.ObjectReference{}

	in := p[0].(map[string]interface{})

	if v, ok := in["kind"].(string); ok {
		obj.Kind = v
	}
	if v, ok := in["name"].(string); ok {
		obj.Name = v
	}

	return obj
}

func expandClusterV2RKEConfigMachinePoolRollingUpdate(p []interface{}) *provisionv1.RKEMachinePoolRollingUpdate {
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil
	}

	obj := &provisionv1.RKEMachinePoolRollingUpdate{}

	in := p[0].(map[string]interface{})

	if v, ok := in["max_surge"].(string); ok && len(v) > 0 {
		maxSurge := intstr.FromString(v)
		obj.MaxSurge = &maxSurge
	}
	if v, ok := in["max_unavailable"].(string); ok && len(v) > 0 {
		maxUnavailable := intstr.FromString(v)
		obj.MaxUnavailable = &maxUnavailable
	}

	return obj
}

func expandClusterV2RKEConfigMachinePools(p []interface{}) []provisionv1.RKEMachinePool {
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil
	}

	out := make([]provisionv1.RKEMachinePool, len(p))
	for i := range p {
		in := p[i].(map[string]interface{})
		obj := provisionv1.RKEMachinePool{}

		if v, ok := in["name"].(string); ok {
			obj.Name = v
			obj.DisplayName = v
		}
		if v, ok := in["cloud_credential_secret_name"].(string); ok && len(v) > 0 {
			obj.CloudCredentialSecretName = v
		}
		if v, ok := in["machine_config"].([]interface{}); ok && len(v) > 0 {
			obj.NodeConfig = expandClusterV2RKEConfigMachinePoolMachineConfig(v)
		}
		if v, ok := in["control_plane_role"].(bool); ok {
			obj.ControlPlaneRole = v
		}
		if v, ok := in["etcd_role"].(bool); ok {
			obj.EtcdRole = v
		}
		if v, ok := in["drain_before_delete"].(bool); ok {
			obj.DrainBeforeDelete = v
		}
		if v, ok := in["annotations"].(map[string]interface{}); ok && len(v) > 0 {
			obj.MachineDeploymentAnnotations = toMapString(v)
		}
		if v, ok := in["labels"].(map[string]interface{}); ok && len(v) > 0 {
			obj.MachineDeploymentLabels = toMapString(v)
		}
		if v, ok := in["machine_labels"].(map[string]interface{}); ok && len(v) > 0 {
			obj.Labels = toMapString(v)
		}
		if v, ok := in["paused"].(bool); ok {
			obj.Paused = v
		}
		if v, ok := in["quantity"].(int); ok {
			quantity := int32(v)
			obj.Quantity = &quantity
		}
		if v, ok := in["rolling_update"].([]interface{}); ok && len(v) > 0 {
			obj.RollingUpdate = expandClusterV2RKEConfigMachinePoolRollingUpdate(v)
		}
		if v, ok := in["taints"].([]interface{}); ok && len(v) > 0 {
			obj.Taints = expandTaintsV2(v)
		}
		if v, ok := in["worker_role"].(bool); ok {
			obj.WorkerRole = v
		}
		if v, ok := in["node_startup_timeout_seconds"].(int); ok && v > 0 {
			d := metav1.Duration{Duration: time.Duration(v) * time.Second}
			obj.NodeStartupTimeout = &d
		}
		if v, ok := in["unhealthy_node_timeout_seconds"].(int); ok && v > 0 {
			d := metav1.Duration{Duration: time.Duration(v) * time.Second}
			obj.UnhealthyNodeTimeout = &d
		}
		if v, ok := in["node_drain_timeout"].(int); ok && v > 0 {
			d := metav1.Duration{Duration: time.Duration(v) * time.Second}
			obj.DrainBeforeDeleteTimeout = &d
		}
		if v, ok := in["max_unhealthy"].(string); ok && len(v) > 0 {
			obj.MaxUnhealthy = &v
		}
		if v, ok := in["unhealthy_range"].(string); ok && len(v) > 0 {
			obj.UnhealthyRange = &v
		}

		out[i] = obj
	}

	return out
}
