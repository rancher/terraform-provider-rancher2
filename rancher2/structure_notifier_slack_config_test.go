package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
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
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
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
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
