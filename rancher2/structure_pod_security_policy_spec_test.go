package rancher2

import (
	"reflect"
	"testing"

	policyv1 "k8s.io/api/policy/v1beta1"
)

var (
	testPodSecurityPolicyBool                         bool
	testPodSecurityPolicySpecConf                     *policyv1.PodSecurityPolicySpec
	testPodSecurityPolicySpecInterface                []interface{}
	testPodSecurityPolicySupplementalGroups2Conf      policyv1.SupplementalGroupsStrategyOptions
	testPodSecurityPolicySupplementalGroups2Interface []interface{}
	testPodSecurityPolicyIDRanges4Conf                []policyv1.IDRange
	testPodSecurityPolicyIDRanges4Interface           []interface{}
)

func init() {
	testPodSecurityPolicyIDRanges4Conf = []policyv1.IDRange{
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
	testPodSecurityPolicySupplementalGroups2Conf = policyv1.SupplementalGroupsStrategyOptions{
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
	testPodSecurityPolicySpecConf = &policyv1.PodSecurityPolicySpec{
		Privileged: true,
		DefaultAddCapabilities: testPodSecurityPolicyCapabilitiesConf,	
		RequiredDropCapabilities: testPodSecurityPolicyCapabilitiesConf,
		AllowedCapabilities: testPodSecurityPolicyCapabilitiesConf,
		Volumes: []policyv1.FSType{"hostPath", "emptyDir"},
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
			"volumes": []interface{}{"hostPath", "emptyDir"},
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
