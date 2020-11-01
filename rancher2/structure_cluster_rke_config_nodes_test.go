package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testClusterRKEConfigNodeDrainInputConf           *managementClient.NodeDrainInput
	testClusterRKEConfigNodeDrainInputInterface      []interface{}
	testClusterRKEConfigNodeUpgradeStrategyConf      *managementClient.NodeUpgradeStrategy
	testClusterRKEConfigNodeUpgradeStrategyInterface []interface{}
	testClusterRKEConfigNodesConf                    []managementClient.RKEConfigNode
	testClusterRKEConfigNodesInterface               []interface{}
)

func init() {
	testClusterRKEConfigNodeDrainInputConf = &managementClient.NodeDrainInput{
		DeleteLocalData:  false,
		Force:            false,
		GracePeriod:      -1,
		IgnoreDaemonSets: newTrue(),
		Timeout:          60,
	}
	testClusterRKEConfigNodeDrainInputInterface = []interface{}{
		map[string]interface{}{
			"delete_local_data":  false,
			"force":              false,
			"grace_period":       -1,
			"ignore_daemon_sets": true,
			"timeout":            60,
		},
	}
	testClusterRKEConfigNodeUpgradeStrategyConf = &managementClient.NodeUpgradeStrategy{
		Drain:                      newFalse(),
		DrainInput:                 testClusterRKEConfigNodeDrainInputConf,
		MaxUnavailableControlplane: "2",
		MaxUnavailableWorker:       "20%",
	}
	testClusterRKEConfigNodeUpgradeStrategyInterface = []interface{}{
		map[string]interface{}{
			"drain":                        false,
			"drain_input":                  testClusterRKEConfigNodeDrainInputInterface,
			"max_unavailable_controlplane": "2",
			"max_unavailable_worker":       "20%",
		},
	}
	testClusterRKEConfigNodesConf = []managementClient.RKEConfigNode{
		{
			Address:          "url.terraform.test",
			DockerSocket:     "docker.sock",
			HostnameOverride: "terra-test",
			InternalAddress:  "192.168.1.1",
			Labels: map[string]string{
				"label_one": "one",
				"label_two": "two",
			},
			NodeID:       "test1",
			Port:         "22",
			Role:         []string{"worker"},
			SSHAgentAuth: true,
			SSHKey:       "XXXXXXXX",
			SSHKeyPath:   "/home/user/.ssh",
			User:         "test",
		},
	}
	testClusterRKEConfigNodesInterface = []interface{}{
		map[string]interface{}{
			"address":           "url.terraform.test",
			"docker_socket":     "docker.sock",
			"hostname_override": "terra-test",
			"internal_address":  "192.168.1.1",
			"labels": map[string]interface{}{
				"label_one": "one",
				"label_two": "two",
			},
			"node_id":        "test1",
			"port":           "22",
			"role":           []interface{}{"worker"},
			"ssh_agent_auth": true,
			"ssh_key":        "XXXXXXXX",
			"ssh_key_path":   "/home/user/.ssh",
			"user":           "test",
		},
	}
}

func TestFlattenClusterRKEConfigNodeDrainInput(t *testing.T) {

	cases := []struct {
		Input          *managementClient.NodeDrainInput
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigNodeDrainInputConf,
			testClusterRKEConfigNodeDrainInputInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterRKEConfigNodeDrainInput(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterRKEConfigNodeUpgradeStrategy(t *testing.T) {

	cases := []struct {
		Input          *managementClient.NodeUpgradeStrategy
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigNodeUpgradeStrategyConf,
			testClusterRKEConfigNodeUpgradeStrategyInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterRKEConfigNodeUpgradeStrategy(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterRKEConfigNodes(t *testing.T) {

	cases := []struct {
		Input          []managementClient.RKEConfigNode
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigNodesConf,
			testClusterRKEConfigNodesInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigNodes(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigNodeDrainInput(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.NodeDrainInput
	}{
		{
			testClusterRKEConfigNodeDrainInputInterface,
			testClusterRKEConfigNodeDrainInputConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterRKEConfigNodeDrainInput(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigNodeUpgradeStrategy(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.NodeUpgradeStrategy
	}{
		{
			testClusterRKEConfigNodeUpgradeStrategyInterface,
			testClusterRKEConfigNodeUpgradeStrategyConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterRKEConfigNodeUpgradeStrategy(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigNodes(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput []managementClient.RKEConfigNode
	}{
		{
			testClusterRKEConfigNodesInterface,
			testClusterRKEConfigNodesConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigNodes(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
