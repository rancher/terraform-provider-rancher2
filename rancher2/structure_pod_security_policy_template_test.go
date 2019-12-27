package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

var (
	testPodSecurityPolicyBool                         bool
	testPodSecurityPolicyTemplateConf                 *managementClient.PodSecurityPolicyTemplate
	testPodSecurityPolicyTemplateInterface            map[string]interface{}
)

func init() {
	testPodSecurityPolicyBool = true
	testPodSecurityPolicyTemplateConf = &managementClient.PodSecurityPolicyTemplate{
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
		SupplementalGroups: testPodSecurityPolicySupplementalGroupsConf,
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
	testPodSecurityPolicyTemplateInterface = map[string]interface{}{
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
		"supplemental_groups": testPodSecurityPolicySupplementalGroupsInterface,
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
	}
}

func TestFlattenPodSecurityPolicyTemplate(t *testing.T) {

	cases := []struct {
		Input          *managementClient.PodSecurityPolicyTemplate
		ExpectedOutput map[string]interface{}
	}{
		{
			testPodSecurityPolicyTemplateConf,
			testPodSecurityPolicyTemplateInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, podSecurityPolicyTemplateFields(), map[string]interface{}{})
		err := flattenPodSecurityPolicyTemplate(output, tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		given := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			given[k] = output.Get(k)
		}
		if !reflect.DeepEqual(given, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
			tc.ExpectedOutput, given)
		}
	}
}

func TestExpandPodSecurityPolicyTemplate(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *managementClient.PodSecurityPolicyTemplate
	}{
		{
			testPodSecurityPolicyTemplateInterface,
			testPodSecurityPolicyTemplateConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, podSecurityPolicyTemplateFields(), tc.Input)
		output := expandPodSecurityPolicyTemplate(inputResourceData)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
