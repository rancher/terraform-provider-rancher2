package rancher2

import (
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
)

var (
	testClusterRKEConfigBastionHostConf      *managementClient.BastionHost
	testClusterRKEConfigBastionHostInterface []interface{}
)

func init() {
	testClusterRKEConfigBastionHostConf = &managementClient.BastionHost{
		Address:      "bastion.terraform.test",
		Port:         "22",
		SSHAgentAuth: true,
		SSHKey:       "XXXXXXXX",
		SSHKeyPath:   "/home/user/.ssh",
		User:         "test",
	}
	testClusterRKEConfigBastionHostInterface = []interface{}{
		map[string]interface{}{
			"address":        "bastion.terraform.test",
			"port":           "22",
			"ssh_agent_auth": true,
			"ssh_key":        "XXXXXXXX",
			"ssh_key_path":   "/home/user/.ssh",
			"user":           "test",
		},
	}
}

func TestFlattenClusterRKEConfigBastionHost(t *testing.T) {

	cases := []struct {
		Input          *managementClient.BastionHost
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigBastionHostConf,
			testClusterRKEConfigBastionHostInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigBastionHost(tc.Input, tc.ExpectedOutput)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestExpandClusterRKEConfigBastionHost(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.BastionHost
	}{
		{
			testClusterRKEConfigBastionHostInterface,
			testClusterRKEConfigBastionHostConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigBastionHost(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
