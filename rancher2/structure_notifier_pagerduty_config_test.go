package rancher2

import (
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
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
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
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
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
