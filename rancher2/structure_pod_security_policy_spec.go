package rancher2

import (
    managementClient "github.com/rancher/types/client/management/v3"
)

// Flatteners

func flattenPodSecurityPolicySpec(in *managementClient.PodSecurityPolicySpec) []interface{} {
    spec := make(map[string]interface{})

    if in.AllowPrivilegeEscalation != nil {
		spec["allow_privilege_escalation"] = *in.AllowPrivilegeEscalation
	}

	if len(in.AllowedCapabilities) > 0 {
		spec["allowed_capabilities"] = toArrayInterface(in.AllowedCapabilities)
	}

	if len(in.AllowedFlexVolumes) > 0 {
		spec["allowed_flex_volumes"] = flattenPodSecurityPolicyAllowedFlexVolumes(in.AllowedFlexVolumes)
	}

	if len(in.AllowedHostPaths) > 0 {
		spec["allowed_host_paths"] = flattenPodSecurityPolicyAllowedHostPaths(in.AllowedHostPaths)
	}

	if len(in.AllowedProcMountTypes) > 0 {
		spec["allowed_proc_mount_types"] = toArrayInterface(in.AllowedProcMountTypes)
	}

	if len(in.AllowedUnsafeSysctls) > 0 {
		spec["allowed_unsafe_sysctls"] = toArrayInterface(in.AllowedUnsafeSysctls)
	}

	if len(in.DefaultAddCapabilities) > 0 {
		spec["default_add_capabilities"] = toArrayInterface(in.DefaultAddCapabilities)
	}

	if in.DefaultAllowPrivilegeEscalation != nil {
		spec["default_allow_privilege_escalation"] = *in.DefaultAllowPrivilegeEscalation
	}

	if len(in.ForbiddenSysctls) > 0 {
		spec["forbidden_sysctls"] = toArrayInterface(in.ForbiddenSysctls)
	}

	spec["fs_group"] = flattenPodSecurityPolicyFSGroup(in.FSGroup)
	spec["host_ipc"] = in.HostIPC
	spec["host_network"] = in.HostNetwork
	spec["host_pid"] = in.HostPID

	if len(in.HostPorts) > 0 {
		spec["host_ports"] = flattenPodSecurityPolicyHostPortRanges(in.HostPorts)
	}

	spec["privileged"] = in.Privileged
	spec["read_only_root_filesystem"] = in.ReadOnlyRootFilesystem

	if len(in.RequiredDropCapabilities) > 0 {
		spec["required_drop_capabilities"] = toArrayInterface(in.RequiredDropCapabilities)
	}

	spec["run_as_user"] = flattenPodSecurityPolicyRunAsUser(in.RunAsUser)

	if in.RunAsGroup != nil {
		spec["run_as_group"] = flattenPodSecurityPolicyRunAsGroup(in.RunAsGroup)
	}

	spec["se_linux"] = flattenPodSecurityPolicySELinuxStrategy(in.SELinux)
	spec["supplemental_groups"] = flattenPodSecurityPolicySupplementalGroups(in.SupplementalGroups)
	spec["volumes"] = toArrayInterface(in.Volumes)

    return []interface{}{spec}
}

func expandPodSecurityPolicySpec(in []interface{}) *managementClient.PodSecurityPolicySpec {
	spec := &managementClient.PodSecurityPolicySpec{}
	if len(in) == 0 || in[0] == nil {
		return spec
	}

	m, _ := in[0].(map[string]interface{})
	/*if !ok {
		return spec, fmt.Errorf("failed to expand PodSecurityPolicy.Spec: malformed input")
	}*/

	if v, ok := m["allow_privilege_escalation"].(bool); ok {
		spec.AllowPrivilegeEscalation = &v
	}

	if v, ok := m["allowed_capabilities"].([]interface{}); ok && len(v) > 0 {
		spec.AllowedCapabilities = toArrayString(v)
	}

	if v, ok := m["allowed_flex_volumes"].([]interface{}); ok && len(v) > 0 {
		spec.AllowedFlexVolumes = expandPodSecurityPolicyAllowedFlexVolumes(v)
	}

	if v, ok := m["allowed_host_paths"].([]interface{}); ok && len(v) > 0 {
		spec.AllowedHostPaths = expandPodSecurityPolicyAllowedHostPaths(v)
	}

	if v, ok := m["allowed_proc_mount_types"].([]interface{}); ok && len(v) > 0 {
		spec.AllowedProcMountTypes = toArrayString(v)
	}

	if v, ok := m["allowed_unsafe_sysctls"].([]interface{}); ok && len(v) > 0 {
		spec.AllowedUnsafeSysctls = toArrayString(v)
	}

	if v, ok := m["default_add_capabilities"].([]interface{}); ok && len(v) > 0 {
		spec.DefaultAddCapabilities = toArrayString(v)
	}

	if v, ok := m["default_allow_privilege_escalation"].(bool); ok {
		spec.DefaultAllowPrivilegeEscalation = &v
	}

	if v, ok := m["forbidden_sysctls"].([]interface{}); ok && len(v) > 0 {
		spec.ForbiddenSysctls = toArrayString(v)
	}

	if v, ok := m["fs_group"].([]interface{}); ok && len(v) > 0 {
		spec.FSGroup = expandPodSecurityPolicyFSGroup(v)
	}

	if v, ok := m["host_ipc"].(bool); ok {
		spec.HostIPC = v
	}

	if v, ok := m["host_network"].(bool); ok {
		spec.HostNetwork = v
	}

	if v, ok := m["host_pid"].(bool); ok {
		spec.HostPID = v
	}

	if v, ok := m["host_ports"].([]interface{}); ok && len(v) > 0 {
		spec.HostPorts = expandPodSecurityPolicyHostPortRanges(v)
	}

	if v, ok := m["privileged"].(bool); ok {
		spec.Privileged = v
	}

	if v, ok := m["read_only_root_filesystem"].(bool); ok {
		spec.ReadOnlyRootFilesystem = v
	}

	if v, ok := m["required_drop_capabilities"].([]interface{}); ok && len(v) > 0 {
		spec.RequiredDropCapabilities = toArrayString(v)
	}

	if v, ok := m["run_as_user"].([]interface{}); ok && len(v) > 0 {
		spec.RunAsUser = expandPodSecurityPolicyRunAsUser(v)
	}

	if v, ok := m["run_as_group"].([]interface{}); ok && len(v) > 0 {
		spec.RunAsGroup = expandPodSecurityPolicyRunAsGroup(v)
	}

	if v, ok := m["se_linux"].([]interface{}); ok && len(v) > 0 {
		spec.SELinux = expandPodSecurityPolicySELinuxStrategy(v)
	}

	if v, ok := m["supplemental_groups"].([]interface{}); ok && len(v) > 0 {
		spec.SupplementalGroups = expandPodSecurityPolicySupplementalGroups(v)
	}

	if v, ok := m["volumes"].([]interface{}); ok && len(v) > 0 {
		spec.Volumes = toArrayString(v)
	}

	return spec
}