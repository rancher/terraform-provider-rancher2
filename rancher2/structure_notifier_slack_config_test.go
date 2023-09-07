package rancher2

import (
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
)

var (
	testNotifierSlackConfigConf      *managementClient.SlackConfig
	testNotifierSlackConfigInterface []interface{}
)

func init() {
	testNotifierSlackConfigConf = &managementClient.SlackConfig{
		DefaultRecipient: "default_recipient",
		URL:              "url",
		ProxyURL:         "proxy_url",
	}
	testNotifierSlackConfigInterface = []interface{}{
		map[string]interface{}{
			"default_recipient": "default_recipient",
			"url":               "url",
			"proxy_url":         "proxy_url",
		},
	}
}

func TestFlattenNotifierSlackConfig(t *testing.T) {

	cases := []struct {
		Input          *managementClient.SlackConfig
		ExpectedOutput []interface{}
	}{
		{
			testNotifierSlackConfigConf,
			testNotifierSlackConfigInterface,
		},
	}

	for _, tc := range cases {
		output := flattenNotifierSlackConfig(tc.Input, testNotifierSlackConfigInterface)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestExpandNotifierSlackConfig(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.SlackConfig
	}{
		{
			testNotifierSlackConfigInterface,
			testNotifierSlackConfigConf,
		},
	}

	for _, tc := range cases {
		output := expandNotifierSlackConfig(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
