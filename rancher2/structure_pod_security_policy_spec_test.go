package rancher2

import (
	"reflect"
	"testing"

	policyv1 "k8s.io/api/policy/v1beta1"
)

var (
	testPodSecurityPolicyBoolPtr       *bool
	testPodSecurityPolicySpecConf      *policyv1.PodSecurityPolicySpec
	testPodSecurityPolicySpecInterface []interface{}
)

func init() {
	testPodSecurityPolicyBoolPtr := true
	testPodSecurityPolicySpecConf = &policyv1.PodSecurityPolicySpec{
		Privileged: true,
		DefaultAddCapabilities: testPodSecurityPolicyCapabilitiesConf,	
		RequiredDropCapabilities: testPodSecurityPolicyCapabilitiesConf,
		AllowedCapabilities: testPodSecurityPolicyCapabilitiesConf,
		Volumes: testPodSecurityPolicyVolumesConf,
		HostNetwork: true,
		HostPorts: testPodSecurityPolicyHostPortRangesConf,
		HostPID: false,
		HostIPC: true,
		SELinux: testPodSecurityPolicySELinuxStrategyConf,
		RunAsUser: testPodSecurityPolicyRunAsUserConf,
		RunAsGroup: testPodSecurityPolicyRunAsGroupConf,
		SupplementalGroups: testPodSecurityPolicySupplementalGroupsConf,
		FSGroup: testPodSecurityPolicyFSGroupConf,
		ReadOnlyRootFilesystem: false,
		DefaultAllowPrivilegeEscalation: &testPodSecurityPolicyBoolPtr,
		AllowPrivilegeEscalation: &testPodSecurityPolicyBoolPtr,
		AllowedHostPaths: testPodSecurityPolicyAllowedHostPathsConf,
		AllowedFlexVolumes: testPodSecurityPolicyAllowedFlexVolumesConf,
		AllowedUnsafeSysctls: []string{"foo", "bar"},
		ForbiddenSysctls: []string{"foo", "bar"},
		AllowedProcMountTypes: testPodSecurityPolicyAllowedProcMountTypesConf,
	}
	testPodSecurityPolicySpecInterface = []interface{}{
		map[string]interface{}{
			"privileged": true,
			"default_add_capabilities": testPodSecurityPolicyCapabilitiesSlice,
			"required_drop_capabilities": testPodSecurityPolicyCapabilitiesSlice,
			"allowed_capabilities": testPodSecurityPolicyCapabilitiesSlice,
			"volumes": testPodSecurityPolicyVolumesSlice,
			"host_network": true,
			"host_ports": testPodSecurityPolicyHostPortRangesInterface,
			"host_pid": false,
			"host_ipc": true,
			"se_linux": testPodSecurityPolicySELinuxStrategyInterface,
			"run_as_user": testPodSecurityPolicyRunAsUserInterface,
			"run_as_group": testPodSecurityPolicyRunAsGroupInterface,
			"supplemental_groups": testPodSecurityPolicySupplementalGroupsInterface,
			"fs_group": testPodSecurityPolicyFSGroupInterface,
			"read_only_root_filesystem": false,
			"default_allow_privilege_escalation": &testPodSecurityPolicyBoolPtr,
			"allow_privilege_escalation": &testPodSecurityPolicyBoolPtr,
			"allowed_host_paths": testPodSecurityPolicyAllowedHostPathsInterface,
			"allowed_flex_volumes": testPodSecurityPolicyAllowedFlexVolumesInterface,
			"allowed_unsafe_sysctls": toArrayInterface([]string{"foo", "bar"}),
			"forbidden_sysctls": toArrayInterface([]string{"foo", "bar"}),
			"allowed_proc_mount_types": testPodSecurityPolicyAllowedProcMountTypesSlice,
		},
	}
}

func TestFlattenPodSecurityPolicySpec(t *testing.T) {

	cases := []struct {
		Input          *policyv1.PodSecurityPolicySpec
		ExpectedOutput []interface{}
	}{
		{
			testPodSecurityPolicySpecConf,
			testPodSecurityPolicySpecInterface,
		},
	}

	for _, tc := range cases {
		output := flattenPodSecurityPolicySpec(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandPodSecurityPolicySpec(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *policyv1.PodSecurityPolicySpec
	}{
		{
			testPodSecurityPolicySpecInterface,
			testPodSecurityPolicySpecConf,
		},
	}

	for _, tc := range cases {
		output := expandPodSecurityPolicySpec(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
