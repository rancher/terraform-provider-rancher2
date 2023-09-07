package rancher2

import (
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
)

var (
	testNotifierWechatConfigConf      *managementClient.WechatConfig
	testNotifierWechatConfigInterface []interface{}
)

func init() {
	testNotifierWechatConfigConf = &managementClient.WechatConfig{
		Agent:            "agent",
		Corp:             "corp",
		DefaultRecipient: "default_recipient",
		Secret:           "secret",
		ProxyURL:         "proxy_url",
		RecipientType:    "recipient_type",
	}
	testNotifierWechatConfigInterface = []interface{}{
		map[string]interface{}{
			"agent":             "agent",
			"corp":              "corp",
			"default_recipient": "default_recipient",
			"secret":            "secret",
			"proxy_url":         "proxy_url",
			"recipient_type":    "recipient_type",
		},
	}
}

func TestFlattenNotifierWechatConfig(t *testing.T) {

	cases := []struct {
		Input          *managementClient.WechatConfig
		ExpectedOutput []interface{}
	}{
		{
			testNotifierWechatConfigConf,
			testNotifierWechatConfigInterface,
		},
	}

	for _, tc := range cases {
		output := flattenNotifierWechatConfig(tc.Input, testNotifierWechatConfigInterface)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestExpandNotifierWechatConfig(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.WechatConfig
	}{
		{
			testNotifierWechatConfigInterface,
			testNotifierWechatConfigConf,
		},
	}

	for _, tc := range cases {
		output := expandNotifierWechatConfig(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
