package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func podSecurityPolicyTemplateFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Pod Security Policy template policy name",
			ForceNew:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Pod Security Policy template policy description",
		},
		"allow_privilege_escalation": {
			Type:        schema.TypeBool,
			Description: "allowPrivilegeEscalation determines if a pod can request to allow privilege escalation. If unspecified, defaults to true.",
			Optional:    true,
			Computed:    true,
		},
		"allowed_capabilities": {
			Type:        schema.TypeList,
			Description: "allowedCapabilities is a list of capabilities that can be requested to add to the container. Capabilities in this field may be added at the pod author's discretion. You must not list a capability in both allowedCapabilities and requiredDropCapabilities.",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"allowed_csi_driver": {
			Type:        schema.TypeList,
			Description: "AllowedCSIDrivers is a whitelist of inline CSI drivers that must be explicitly set to be embedded within a pod spec. An empty value indicates that any CSI driver can be used for inline ephemeral volumes. This is an alpha field, and is only honored if the API server enables the CSIInlineVolume feature gate.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: podSecurityPolicyAllowedCSIDriverFields(),
			},
		},
		"allowed_flex_volume": {
			Type:        schema.TypeList,
			Description: "allowedFlexVolumes is a whitelist of allowed Flexvolumes.  Empty or nil indicates that all Flexvolumes may be used.  This parameter is effective only when the usage of the Flexvolumes is allowed in the \"volumes\" field.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: podSecurityPolicyAllowedFlexVolumesFields(),
			},
		},
		"allowed_host_path": {
			Type:        schema.TypeList,
			Description: "allowedHostPaths is a white list of allowed host paths. Empty indicates that all host paths may be used.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: podSecurityPolicyAllowedHostPathsFields(),
			},
		},
		"allowed_proc_mount_types": {
			Type:        schema.TypeList,
			Description: "AllowedProcMountTypes is a whitelist of allowed ProcMountTypes. Empty or nil indicates that only the DefaultProcMountType may be used. This requires the ProcMountType feature flag to be enabled.",
			Optional:    true,
			Elem:        podSecurityPolicyProcMountTypeFields(),
		},
		"allowed_unsafe_sysctls": {
			Type:        schema.TypeList,
			Description: "allowedUnsafeSysctls is a list of explicitly allowed unsafe sysctls, defaults to none. Each entry is either a plain sysctl name or ends in \"*\" in which case it is considered as a prefix of allowed sysctls. Single * means all unsafe sysctls are allowed. Kubelet has to whitelist all allowed unsafe sysctls explicitly to avoid rejection.",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"default_add_capabilities": {
			Type:        schema.TypeList,
			Description: "defaultAddCapabilities is the default set of capabilities that will be added to the container unless the pod spec specifically drops the capability.  You may not list a capability in both defaultAddCapabilities and requiredDropCapabilities. Capabilities added here are implicitly allowed, and need not be included in the allowedCapabilities list.",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"default_allow_privilege_escalation": {
			Type:        schema.TypeBool,
			Description: "defaultAllowPrivilegeEscalation controls the default setting for whether a process can gain more privileges than its parent process.",
			Optional:    true,
		},
		"forbidden_sysctls": {
			Type:        schema.TypeList,
			Description: "forbiddenSysctls is a list of explicitly forbidden sysctls, defaults to none. Each entry is either a plain sysctl name or ends in \"*\" in which case it is considered as a prefix of forbidden sysctls. Single * means all sysctls are forbidden.",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"fs_group": {
			Type:        schema.TypeList,
			Description: "fsGroup is the strategy that will dictate what fs group is used by the SecurityContext.",
			Optional:    true,
			Computed:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: podSecurityPolicyAllowedFsGroupFields(),
			},
		},
		"host_ipc": {
			Type:        schema.TypeBool,
			Description: "hostIPC determines if the policy allows the use of HostIPC in the pod spec.",
			Optional:    true,
			Computed:    true,
		},
		"host_network": {
			Type:        schema.TypeBool,
			Description: "hostNetwork determines if the policy allows the use of HostNetwork in the pod spec.",
			Optional:    true,
			Computed:    true,
		},
		"host_pid": {
			Type:        schema.TypeBool,
			Description: "hostPID determines if the policy allows the use of HostPID in the pod spec.",
			Optional:    true,
			Computed:    true,
		},
		"host_port": {
			Type:        schema.TypeList,
			Description: "hostPorts determines which host port ranges are allowed to be exposed.",
			Optional:    true,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: podSecurityPolicyHostPortRangeFields(),
			},
		},
		"privileged": {
			Type:        schema.TypeBool,
			Description: "privileged determines if a pod can request to be run as privileged.",
			Optional:    true,
			Computed:    true,
		},
		"read_only_root_filesystem": {
			Type:        schema.TypeBool,
			Description: "readOnlyRootFilesystem when set to true will force containers to run with a read only root file system.  If the container specifically requests to run with a non-read only root file system the PSP should deny the pod. If set to false the container may run with a read only root file system if it wishes but it will not be forced to.",
			Optional:    true,
			Computed:    true,
		},
		"required_drop_capabilities": {
			Type:        schema.TypeList,
			Description: "requiredDropCapabilities are the capabilities that will be dropped from the container.  These are required to be dropped and cannot be added.",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"run_as_user": {
			Type:        schema.TypeList,
			Description: "runAsUser is the strategy that will dictate the allowable RunAsUser values that may be set.",
			Optional:    true,
			Computed:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: podSecurityPolicyRunAsUserFields(),
			},
		},
		"run_as_group": {
			Type:        schema.TypeList,
			Description: "RunAsGroup is the strategy that will dictate the allowable RunAsGroup values that may be set. If this field is omitted, the pod's RunAsGroup can take any value. This field requires the RunAsGroup feature gate to be enabled.",
			Optional:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: podSecurityPolicyRunAsGroupFields(),
			},
		},
		"runtime_class": {
			Type:        schema.TypeList,
			Description: "runtimeClass is the strategy that will dictate the allowable RuntimeClasses for a pod. If this field is omitted, the pod's runtimeClassName field is unrestricted. Enforcement of this field depends on the RuntimeClass feature gate being enabled.",
			Optional:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: podSecurityPolicyRuntimeClassFields(),
			},
		},
		"se_linux": {
			Type:        schema.TypeList,
			Description: "seLinux is the strategy that will dictate the allowable labels that may be set.",
			Optional:    true,
			Computed:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: podSecurityPolicySELinuxFields(),
			},
		},
		"supplemental_group": {
			Type:        schema.TypeList,
			Description: "supplementalGroups is the strategy that will dictate what supplemental groups are used by the SecurityContext.",
			Optional:    true,
			Computed:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: podSecurityPolicySupplementalGroupsFields(),
			},
		},
		"volumes": {
			Type:        schema.TypeList,
			Description: "volumes is a white list of allowed volume plugins. Empty indicates that no volumes may be used. To allow all volumes you may use '*'",
			Optional:    true,
			Computed:    true,
			Elem:        podSecurityPolicyVolumesFields(),
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}
