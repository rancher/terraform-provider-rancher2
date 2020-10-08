package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testClusterK3SUpgradeStrategyConfigConf      *managementClient.ClusterUpgradeStrategy
	testClusterK3SUpgradeStrategyConfigInterface []interface{}
	testClusterK3SConfigConf                     *managementClient.K3sConfig
	testClusterK3SConfigInterface                []interface{}
)

func init() {
	testClusterK3SUpgradeStrategyConfigConf = &managementClient.ClusterUpgradeStrategy{
		DrainServerNodes:  true,
		DrainWorkerNodes:  true,
		ServerConcurrency: 2,
		WorkerConcurrency: 2,
	}
	testClusterK3SUpgradeStrategyConfigInterface = []interface{}{
		map[string]interface{}{
			"drain_server_nodes": true,
			"drain_worker_nodes": true,
			"server_concurrency": 2,
			"worker_concurrency": 2,
		},
	}
	testClusterK3SConfigConf = &managementClient.K3sConfig{
		ClusterUpgradeStrategy: testClusterK3SUpgradeStrategyConfigConf,
		Version:                "version",
	}
	testClusterK3SConfigInterface = []interface{}{
		map[string]interface{}{
			"upgrade_strategy": testClusterK3SUpgradeStrategyConfigInterface,
			"version":          "version",
		},
	}
}

func TestFlattenClusterK3SUpgradeStrategyConfig(t *testing.T) {

	cases := []struct {
		Input          *managementClient.ClusterUpgradeStrategy
		ExpectedOutput []interface{}
	}{
		{
			testClusterK3SUpgradeStrategyConfigConf,
			testClusterK3SUpgradeStrategyConfigInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterK3SUpgradeStrategyConfig(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterK3SConfig(t *testing.T) {

	cases := []struct {
		Input          *managementClient.K3sConfig
		ExpectedOutput []interface{}
	}{
		{
			testClusterK3SConfigConf,
			testClusterK3SConfigInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterK3SConfig(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterK3SUpgradeStrategyConfig(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.ClusterUpgradeStrategy
	}{
		{
			testClusterK3SUpgradeStrategyConfigInterface,
			testClusterK3SUpgradeStrategyConfigConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterK3SUpgradeStrategyConfig(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterK3SConfig(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.K3sConfig
	}{
		{
			testClusterK3SConfigInterface,
			testClusterK3SConfigConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterK3SConfig(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
