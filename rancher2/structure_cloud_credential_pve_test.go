package rancher2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testCloudCredentialPveConf      map[string]interface{}
	testCloudCredentialPveInterface []interface{}
)

func init() {
	testCloudCredentialPveConf = map[string]interface{}{
		"pveUrl":         "https://pve.example.com:8006",
		"pveTokenId":     "root@pam!rancher",
		"pveTokenSecret": "secret-uuid",
		"pveInsecureTls": false,
	}
	testCloudCredentialPveInterface = []interface{}{
		map[string]interface{}{
			"pve_url":          "https://pve.example.com:8006",
			"pve_token_id":     "root@pam!rancher",
			"pve_token_secret": "secret-uuid",
			"pve_insecure_tls": false,
		},
	}
}

func TestFlattenCloudCredentialPve(t *testing.T) {
	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput []interface{}
	}{
		{
			testCloudCredentialPveConf,
			testCloudCredentialPveInterface,
		},
	}

	for _, tc := range cases {
		output := flattenCloudCredentialPve(tc.Input, []interface{}{})
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestExpandCloudCredentialPve(t *testing.T) {
	cases := []struct {
		Input          []interface{}
		ExpectedOutput map[string]interface{}
	}{
		{
			testCloudCredentialPveInterface,
			testCloudCredentialPveConf,
		},
	}

	for _, tc := range cases {
		output := expandCloudCredentialPve(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
