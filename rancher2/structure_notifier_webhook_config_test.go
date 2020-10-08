package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
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
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
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
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
