package rancher2

import (
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
)

var (
	testNotifierWebhookConfigConf      *managementClient.WebhookConfig
	testNotifierWebhookConfigInterface []interface{}
)

func init() {
	testNotifierWebhookConfigConf = &managementClient.WebhookConfig{
		URL:      "url",
		ProxyURL: "proxy_url",
	}
	testNotifierWebhookConfigInterface = []interface{}{
		map[string]interface{}{
			"url":       "url",
			"proxy_url": "proxy_url",
		},
	}
}

func TestFlattenNotifierWebhookConfig(t *testing.T) {

	cases := []struct {
		Input          *managementClient.WebhookConfig
		ExpectedOutput []interface{}
	}{
		{
			testNotifierWebhookConfigConf,
			testNotifierWebhookConfigInterface,
		},
	}

	for _, tc := range cases {
		output := flattenNotifierWebhookConfig(tc.Input, testNotifierWebhookConfigInterface)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestExpandNotifierWebhookConfig(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.WebhookConfig
	}{
		{
			testNotifierWebhookConfigInterface,
			testNotifierWebhookConfigConf,
		},
	}

	for _, tc := range cases {
		output := expandNotifierWebhookConfig(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
