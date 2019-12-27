package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/types/client/management/v3"
)

var (
	testPodSecurityPolicyBool                         bool
	testPodSecurityPolicySpecConf                     *managementClient.PodSecurityPolicySpec
	testPodSecurityPolicySpecInterface                []interface{}
	testPodSecurityPolicySupplementalGroups2Conf      *managementClient.SupplementalGroupsStrategyOptions
	testPodSecurityPolicySupplementalGroups2Interface []interface{}
	testPodSecurityPolicyIDRanges4Conf                []managementClient.IDRange
	testPodSecurityPolicyIDRanges4Interface           []interface{}
)

func init() {
	testPodSecurityPolicyIDRanges4Conf = []managementClient.IDRange{
		{
			Min: int64(1),
			Max: int64(3000),
		},
		{
			Min: int64(0),
			Max: int64(5000),
		},
	}
	testPodSecurityPolicyIDRanges4Interface = []interface{}{
		map[string]interface{}{
			"min": 1,
			"max": 3000,
		},
		map[string]interface{}{
			"min": 0,
			"max": 5000,
		},
	}
	testPodSecurityPolicySupplementalGroups2Conf = &managementClient.SupplementalGroupsStrategyOptions{
		Rule: "RunAsAny",
		Ranges: testPodSecurityPolicyIDRanges4Conf,
	}
	testPodSecurityPolicySupplementalGroups2Interface = []interface{}{
		map[string]interface{}{
			"rule": "RunAsAny",
			"ranges": testPodSecurityPolicyIDRanges4Interface,
		},
	}

	testPodSecurityPolicyBool = true
	testPodSecurityPolicySpecConf = &managementClient.PodSecurityPolicySpec{
		Privileged: true,
		DefaultAddCapabilities: []string{"NET_ADMIN"},
		RequiredDropCapabilities: []string{"NET_ADMIN"},
		AllowedCapabilities: []string{"NET_ADMIN"},
		Volumes: []string{"hostPath", "emptyDir"},
		HostNetwork: true,
		HostPorts: testPodSecurityPolicyHostPortRangesConf,
		HostPID: false,
		HostIPC: true,
		SELinux: testPodSecurityPolicySELinuxStrategyConf,
		RunAsUser: testPodSecurityPolicyRunAsUserConf,
		RunAsGroup: testPodSecurityPolicyRunAsGroupConf,
		SupplementalGroups: testPodSecurityPolicySupplementalGroups2Conf,
		FSGroup: testPodSecurityPolicyFSGroupConf,
		ReadOnlyRootFilesystem: false,
		DefaultAllowPrivilegeEscalation: &testPodSecurityPolicyBool,
		AllowPrivilegeEscalation: &testPodSecurityPolicyBool,
		AllowedHostPaths: testPodSecurityPolicyAllowedHostPathsConf,
		AllowedFlexVolumes: testPodSecurityPolicyAllowedFlexVolumesConf,
		AllowedCSIDrivers: testPodSecurityPolicyAllowedCSIDriversConf,
		AllowedUnsafeSysctls: []string{"foo", "bar"},
		ForbiddenSysctls: []string{"foo", "bar"},
		AllowedProcMountTypes: []string{"Default", "Unmasked"},
		RuntimeClass: testPodSecurityPolicyRuntimeClassStrategyConf,
	}
	testPodSecurityPolicySpecInterface = []interface{}{
		map[string]interface{}{
			"privileged": true,
			"default_add_capabilities": toArrayInterface([]string{"NET_ADMIN"}),
			"required_drop_capabilities": toArrayInterface([]string{"NET_ADMIN"}),
			"allowed_capabilities": toArrayInterface([]string{"NET_ADMIN"}),
			"volumes": toArrayInterface([]string{"hostPath", "emptyDir"}),
			"host_network": true,
			"host_ports": testPodSecurityPolicyHostPortRangesInterface,
			"host_pid": false,
			"host_ipc": true,
			"se_linux": testPodSecurityPolicySELinuxStrategyInterface,
			"run_as_user": testPodSecurityPolicyRunAsUserInterface,
			"run_as_group": testPodSecurityPolicyRunAsGroupInterface,
			"supplemental_groups": testPodSecurityPolicySupplementalGroups2Interface,
			"fs_group": testPodSecurityPolicyFSGroupInterface,
			"read_only_root_filesystem": false,
			"default_allow_privilege_escalation": testPodSecurityPolicyBool,
			"allow_privilege_escalation": testPodSecurityPolicyBool,
			"allowed_host_paths": testPodSecurityPolicyAllowedHostPathsInterface,
			"allowed_flex_volumes": testPodSecurityPolicyAllowedFlexVolumesInterface,
			"allowed_csi_drivers": testPodSecurityPolicyAllowedCSIDriversInterface,
			"allowed_unsafe_sysctls": toArrayInterface([]string{"foo", "bar"}),
			"forbidden_sysctls": toArrayInterface([]string{"foo", "bar"}),
			"allowed_proc_mount_types": toArrayInterface([]string{"Default", "Unmasked"}),
			"runtime_class": testPodSecurityPolicyRuntimeClassStrategyInterface,
		},
	}
}

func TestFlattenPodSecurityPolicySpec(t *testing.T) {

	cases := []struct {
		Input          *managementClient.PodSecurityPolicySpec
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
		ExpectedOutput *managementClient.PodSecurityPolicySpec
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
