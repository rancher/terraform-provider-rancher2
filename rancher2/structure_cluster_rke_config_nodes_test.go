package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/types/client/management/v3"
)

var (
	testClusterRKEConfigNodesConf      []managementClient.RKEConfigNode
	testClusterRKEConfigNodesInterface []interface{}
)

func init() {
	testClusterRKEConfigNodesConf = []managementClient.RKEConfigNode{
		managementClient.RKEConfigNode{
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
