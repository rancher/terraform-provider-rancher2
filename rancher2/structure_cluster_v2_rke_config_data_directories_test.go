package rancher2

import (
	"testing"

	rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
	"github.com/stretchr/testify/assert"
)

var (
	testClusterV2RKEConfigDataDirectoriesConf      rkev1.DataDirectories
	testClusterV2RKEConfigDataDirectoriesInterface []interface{}
)

func init() {
	testClusterV2RKEConfigDataDirectoriesConf = rkev1.DataDirectories{
		SystemAgent:  "/tmp/test/agent",
		Provisioning: "/tmp/test/provisioning",
		K8sDistro:    "/tmp/test/distro",
	}

	testClusterV2RKEConfigDataDirectoriesInterface = []any{
		map[string]any{
			"system_agent": "/tmp/test/agent",
			"provisioning": "/tmp/test/provisioning",
			"k8s_distro":   "/tmp/test/distro",
		},
	}
}

func Test_flattenClusterV2RKEConfigDataDirectories(t *testing.T) {
	cases := []struct {
		Input          rkev1.DataDirectories
		ExpectedOutput []interface{}
	}{
		{
			testClusterV2RKEConfigDataDirectoriesConf,
			testClusterV2RKEConfigDataDirectoriesInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterV2RKEConfigDataDirectories(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestExpandClusterV2RKEConfigDataDirectories(t *testing.T) {
	cases := []struct {
		Input          []interface{}
		ExpectedOutput rkev1.DataDirectories
	}{
		{
			testClusterV2RKEConfigDataDirectoriesInterface,
			testClusterV2RKEConfigDataDirectoriesConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterV2RKEConfigDataDirectories(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
