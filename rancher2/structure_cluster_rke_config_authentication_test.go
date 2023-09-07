package rancher2

import (
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
)

var (
	testClusterRKEConfigAuthenticationConf      *managementClient.AuthnConfig
	testClusterRKEConfigAuthenticationInterface []interface{}
)

func init() {
	testClusterRKEConfigAuthenticationConf = &managementClient.AuthnConfig{
		SANs:     []string{"sans1", "sans2"},
		Strategy: "strategy",
	}
	testClusterRKEConfigAuthenticationInterface = []interface{}{
		map[string]interface{}{
			"sans":     []interface{}{"sans1", "sans2"},
			"strategy": "strategy",
		},
	}
}

func TestFlattenClusterRKEConfigAuthentication(t *testing.T) {

	cases := []struct {
		Input          *managementClient.AuthnConfig
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigAuthenticationConf,
			testClusterRKEConfigAuthenticationInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigAuthentication(tc.Input)
		if err != nil {
			assert.FailNow(t, "[ERROR] on flattener: %#v", err)
		}
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestExpandClusterRKEConfigAuthentication(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.AuthnConfig
	}{
		{
			testClusterRKEConfigAuthenticationInterface,
			testClusterRKEConfigAuthenticationConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigAuthentication(tc.Input)
		if err != nil {
			assert.FailNow(t, "[ERROR] on expander: %#v", err)
		}
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
