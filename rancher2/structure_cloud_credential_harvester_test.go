package rancher2

import (
	"reflect"
	"testing"
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
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
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
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
