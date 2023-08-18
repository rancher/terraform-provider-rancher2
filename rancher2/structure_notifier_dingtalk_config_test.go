package rancher2

import (
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
)

var (
	testNotifierDingtalkConfigConf      *managementClient.DingtalkConfig
	testNotifierDingtalkConfigInterface []interface{}
)

func init() {
	testNotifierDingtalkConfigConf = &managementClient.DingtalkConfig{
		URL:      "url",
		ProxyURL: "proxy_url",
		Secret:   "secret",
	}
	testNotifierDingtalkConfigInterface = []interface{}{
		map[string]interface{}{
			"url":       "url",
			"proxy_url": "proxy_url",
			"secret":    "secret",
		},
	}
}

func TestFlattenNotifierDingtalkConfig(t *testing.T) {

	cases := []struct {
		Input          *managementClient.DingtalkConfig
		ExpectedOutput []interface{}
	}{
		{
			testNotifierDingtalkConfigConf,
			testNotifierDingtalkConfigInterface,
		},
	}
	for _, tc := range cases {
		output := flattenNotifierDingtalkConfig(tc.Input, tc.ExpectedOutput)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestExpandNotifierDingtalkConfig(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.DingtalkConfig
	}{
		{
			testNotifierDingtalkConfigInterface,
			testNotifierDingtalkConfigConf,
		},
	}

	for _, tc := range cases {
		output := expandNotifierDingtalkConfig(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")

	}
}
