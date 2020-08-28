package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testNotifierPagerdutyConfigConf      *managementClient.PagerdutyConfig
	testNotifierPagerdutyConfigInterface []interface{}
)

func init() {
	testNotifierPagerdutyConfigConf = &managementClient.PagerdutyConfig{
		ServiceKey: "service_key",
		ProxyURL:   "proxy_url",
	}
	testNotifierPagerdutyConfigInterface = []interface{}{
		map[string]interface{}{
			"service_key": "service_key",
			"proxy_url":   "proxy_url",
		},
	}
}

func TestFlattenNotifierPagerdutyConfig(t *testing.T) {

	cases := []struct {
		Input          *managementClient.PagerdutyConfig
		ExpectedOutput []interface{}
	}{
		{
			testNotifierPagerdutyConfigConf,
			testNotifierPagerdutyConfigInterface,
		},
	}

	for _, tc := range cases {
		output := flattenNotifierPagerdutyConfig(tc.Input, testNotifierPagerdutyConfigInterface)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandNotifierPagerdutyConfig(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.PagerdutyConfig
	}{
		{
			testNotifierPagerdutyConfigInterface,
			testNotifierPagerdutyConfigConf,
		},
	}

	for _, tc := range cases {
		output := expandNotifierPagerdutyConfig(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
