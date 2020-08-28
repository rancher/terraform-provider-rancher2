package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenPodSecurityPolicyTemplate(d *schema.ResourceData, in *managementClient.PodSecurityPolicyTemplate) error {
	if in == nil {
		return fmt.Errorf("[ERROR] flattening pod security policy template: Input setting is nil")
	}

	d.SetId(in.ID)
	d.Set("name", in.Name)

	if len(in.Description) > 0 {
		d.Set("description", in.Description)
	}

	err := d.Set("annotations", toMapInterface(in.Annotations))
	if err != nil {
		return err
	}

	err = d.Set("labels", toMapInterface(in.Labels))
	if err != nil {
		return err
	}

	if in.AllowPrivilegeEscalation != nil {
		d.Set("allow_privilege_escalation", *in.AllowPrivilegeEscalation)
	}

	if len(in.AllowedCapabilities) > 0 {
		d.Set("allowed_capabilities", toArrayInterface(in.AllowedCapabilities))
	}

	if len(in.AllowedCSIDrivers) > 0 {
		d.Set("allowed_csi_driver", flattenPodSecurityPolicyAllowedCSIDrivers(in.AllowedCSIDrivers))
	}

	if len(in.AllowedFlexVolumes) > 0 {
		d.Set("allowed_flex_volume", flattenPodSecurityPolicyAllowedFlexVolumes(in.AllowedFlexVolumes))
	}

	if len(in.AllowedHostPaths) > 0 {
		d.Set("allowed_host_path", flattenPodSecurityPolicyAllowedHostPaths(in.AllowedHostPaths))
	}

	if len(in.AllowedProcMountTypes) > 0 {
		d.Set("allowed_proc_mount_types", toArrayInterface(in.AllowedProcMountTypes))
	}

	if len(in.AllowedUnsafeSysctls) > 0 {
		d.Set("allowed_unsafe_sysctls", toArrayInterface(in.AllowedUnsafeSysctls))
	}

	if len(in.DefaultAddCapabilities) > 0 {
		d.Set("default_add_capabilities", toArrayInterface(in.DefaultAddCapabilities))
	}

	if in.DefaultAllowPrivilegeEscalation != nil {
		d.Set("default_allow_privilege_escalation", *in.DefaultAllowPrivilegeEscalation)
	}

	if len(in.ForbiddenSysctls) > 0 {
		d.Set("forbidden_sysctls", toArrayInterface(in.ForbiddenSysctls))
	}

	d.Set("fs_group", flattenPodSecurityPolicyFSGroup(in.FSGroup))
	d.Set("host_ipc", in.HostIPC)
	d.Set("host_network", in.HostNetwork)
	d.Set("host_pid", in.HostPID)

	if len(in.HostPorts) > 0 {
		d.Set("host_port", flattenPodSecurityPolicyHostPortRanges(in.HostPorts))
	}

	d.Set("privileged", in.Privileged)
	d.Set("read_only_root_filesystem", in.ReadOnlyRootFilesystem)

	if len(in.RequiredDropCapabilities) > 0 {
		d.Set("required_drop_capabilities", toArrayInterface(in.RequiredDropCapabilities))
	}

	d.Set("run_as_user", flattenPodSecurityPolicyRunAsUser(in.RunAsUser))

	if in.RunAsGroup != nil {
		d.Set("run_as_group", flattenPodSecurityPolicyRunAsGroup(in.RunAsGroup))
	}

	d.Set("runtime_class", flattenPodSecurityPolicyRuntimeClassStrategy(in.RuntimeClass))
	d.Set("se_linux", flattenPodSecurityPolicySELinuxStrategy(in.SELinux))
	d.Set("supplemental_group", flattenPodSecurityPolicySupplementalGroups(in.SupplementalGroups))
	d.Set("volumes", toArrayInterface(in.Volumes))

	return nil
}

func expandPodSecurityPolicyTemplate(in *schema.ResourceData) *managementClient.PodSecurityPolicyTemplate {

	if in == nil {
		return nil
	}

	obj := &managementClient.PodSecurityPolicyTemplate{}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	if v, ok := in.Get("name").(string); ok && len(v) > 0 {
		obj.Name = v
	}

	if v, ok := in.Get("description").(string); ok && len(v) > 0 {
		obj.Description = v
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	if v, ok := in.Get("allow_privilege_escalation").(bool); ok {
		obj.AllowPrivilegeEscalation = &v
	}

	if v, ok := in.Get("allowed_capabilities").([]interface{}); ok && len(v) > 0 {
		obj.AllowedCapabilities = toArrayString(v)
	}

	if v, ok := in.Get("allowed_csi_driver").([]interface{}); ok && len(v) > 0 {
		obj.AllowedCSIDrivers = expandPodSecurityPolicyAllowedCSIDrivers(v)
	}

	if v, ok := in.Get("allowed_flex_volume").([]interface{}); ok && len(v) > 0 {
		obj.AllowedFlexVolumes = expandPodSecurityPolicyAllowedFlexVolumes(v)
	}

	if v, ok := in.Get("allowed_host_path").([]interface{}); ok && len(v) > 0 {
		obj.AllowedHostPaths = expandPodSecurityPolicyAllowedHostPaths(v)
	}

	if v, ok := in.Get("allowed_proc_mount_types").([]interface{}); ok && len(v) > 0 {
		obj.AllowedProcMountTypes = toArrayString(v)
	}

	if v, ok := in.Get("allowed_unsafe_sysctls").([]interface{}); ok && len(v) > 0 {
		obj.AllowedUnsafeSysctls = toArrayString(v)
	}

	if v, ok := in.Get("default_add_capabilities").([]interface{}); ok && len(v) > 0 {
		obj.DefaultAddCapabilities = toArrayString(v)
	}

	if v, ok := in.Get("default_allow_privilege_escalation").(bool); ok {
		obj.DefaultAllowPrivilegeEscalation = &v
	}

	if v, ok := in.Get("forbidden_sysctls").([]interface{}); ok && len(v) > 0 {
		obj.ForbiddenSysctls = toArrayString(v)
	}

	if v, ok := in.Get("fs_group").([]interface{}); ok && len(v) > 0 {
		obj.FSGroup = expandPodSecurityPolicyFSGroup(v)
	}

	if v, ok := in.Get("host_ipc").(bool); ok {
		obj.HostIPC = v
	}

	if v, ok := in.Get("host_network").(bool); ok {
		obj.HostNetwork = v
	}

	if v, ok := in.Get("host_pid").(bool); ok {
		obj.HostPID = v
	}

	if v, ok := in.Get("host_port").([]interface{}); ok && len(v) > 0 {
		obj.HostPorts = expandPodSecurityPolicyHostPortRanges(v)
	}

	if v, ok := in.Get("privileged").(bool); ok {
		obj.Privileged = v
	}

	if v, ok := in.Get("read_only_root_filesystem").(bool); ok {
		obj.ReadOnlyRootFilesystem = v
	}

	if v, ok := in.Get("required_drop_capabilities").([]interface{}); ok && len(v) > 0 {
		obj.RequiredDropCapabilities = toArrayString(v)
	}

	if v, ok := in.Get("run_as_user").([]interface{}); ok && len(v) > 0 {
		obj.RunAsUser = expandPodSecurityPolicyRunAsUser(v)
	}

	if v, ok := in.Get("run_as_group").([]interface{}); ok && len(v) > 0 {
		obj.RunAsGroup = expandPodSecurityPolicyRunAsGroup(v)
	}

	if v, ok := in.Get("runtime_class").([]interface{}); ok && len(v) > 0 {
		obj.RuntimeClass = expandPodSecurityPolicyRuntimeClassStrategy(v)
	}

	if v, ok := in.Get("se_linux").([]interface{}); ok && len(v) > 0 {
		obj.SELinux = expandPodSecurityPolicySELinuxStrategy(v)
	}

	if v, ok := in.Get("supplemental_group").([]interface{}); ok && len(v) > 0 {
		obj.SupplementalGroups = expandPodSecurityPolicySupplementalGroups(v)
	}

	if v, ok := in.Get("volumes").([]interface{}); ok && len(v) > 0 {
		obj.Volumes = toArrayString(v)
	}

	return obj
}
