package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testClusterRKE2ConfigConf      *managementClient.Rke2Config
	testClusterRKE2ConfigInterface []interface{}
)

func init() {
	testClusterRKE2ConfigConf = &managementClient.Rke2Config{
		ClusterUpgradeStrategy: testClusterK3SUpgradeStrategyConfigConf,
		Version:                "version",
	}
	testClusterRKE2ConfigInterface = []interface{}{
		map[string]interface{}{
			"upgrade_strategy": testClusterK3SUpgradeStrategyConfigInterface,
			"version":          "version",
		},
	}
}

func TestFlattenClusterRKE2Config(t *testing.T) {

	cases := []struct {
		Input          *managementClient.Rke2Config
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKE2ConfigConf,
			testClusterRKE2ConfigInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterRKE2Config(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKE2Config(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.Rke2Config
	}{
		{
			testClusterRKE2ConfigInterface,
			testClusterRKE2ConfigConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterRKE2Config(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
