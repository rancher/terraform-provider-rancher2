package rancher2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testCloudCredentialHarvesterConf, testCloudCredentialHarvesterConfImported           *harvesterCredentialConfig
	testCloudCredentialHarvesterInterface, testCloudCredentialHarvesterInterfaceImported []interface{}
)

func init() {
	testCloudCredentialHarvesterConf = &harvesterCredentialConfig{
		ClusterType:       "external",
		KubeconfigContent: "kubeconfigContent",
	}
	testCloudCredentialHarvesterInterface = []interface{}{
		map[string]interface{}{
			"cluster_type":       "external",
			"kubeconfig_content": "kubeconfigContent",
		},
	}
	testCloudCredentialHarvesterConfImported = &harvesterCredentialConfig{
		ClusterID:         "clusterId",
		ClusterType:       "imported",
		KubeconfigContent: "kubeconfigContent",
	}
	testCloudCredentialHarvesterInterfaceImported = []interface{}{
		map[string]interface{}{
			"cluster_id":         "clusterId",
			"cluster_type":       "imported",
			"kubeconfig_content": "kubeconfigContent",
		},
	}
}

func TestFlattenCloudCredentialHarvester(t *testing.T) {

	cases := []struct {
		Input          *harvesterCredentialConfig
		ExpectedOutput []interface{}
	}{
		{
			testCloudCredentialHarvesterConf,
			testCloudCredentialHarvesterInterface,
		},
		{
			testCloudCredentialHarvesterConfImported,
			testCloudCredentialHarvesterInterfaceImported,
		},
	}

	for _, tc := range cases {
		output := flattenCloudCredentialHarvester(tc.Input, tc.ExpectedOutput)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestExpandCloudCredentialHarvester(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *harvesterCredentialConfig
	}{
		{
			testCloudCredentialHarvesterInterface,
			testCloudCredentialHarvesterConf,
		},
	}

	for _, tc := range cases {
		output := expandCloudCredentialHarvester(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
