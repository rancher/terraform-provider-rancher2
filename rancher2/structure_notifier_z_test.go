package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

var (
	testNotifierPagerdutyConf      *managementClient.Notifier
	testNotifierPagerdutyInterface map[string]interface{}
	testNotifierSlackConf          *managementClient.Notifier
	testNotifierSlackInterface     map[string]interface{}
	testNotifierSMTPConf           *managementClient.Notifier
	testNotifierSMTPInterface      map[string]interface{}
	testNotifierWebhookConf        *managementClient.Notifier
	testNotifierWebhookInterface   map[string]interface{}
	testNotifierWechatConf         *managementClient.Notifier
	testNotifierWechatInterface    map[string]interface{}
)

func init() {
	testNotifierPagerdutyConf = &managementClient.Notifier{
		Name:            "name",
		ClusterID:       "cluster_id",
		Description:     "description",
		PagerdutyConfig: testNotifierPagerdutyConfigConf,
	}
	testNotifierPagerdutyInterface = map[string]interface{}{
		"name":             "name",
		"cluster_id":       "cluster_id",
		"description":      "description",
		"pagerduty_config": testNotifierPagerdutyConfigInterface,
	}
	testNotifierSlackConf = &managementClient.Notifier{
		Name:        "name",
		ClusterID:   "cluster_id",
		Description: "description",
		SlackConfig: testNotifierSlackConfigConf,
	}
	testNotifierSlackInterface = map[string]interface{}{
		"name":         "name",
		"cluster_id":   "cluster_id",
		"description":  "description",
		"slack_config": testNotifierSlackConfigInterface,
	}
	testNotifierSMTPConf = &managementClient.Notifier{
		Name:        "name",
		ClusterID:   "cluster_id",
		Description: "description",
		SMTPConfig:  testNotifierSMTPConfigConf,
	}
	testNotifierSMTPInterface = map[string]interface{}{
		"name":        "name",
		"cluster_id":  "cluster_id",
		"description": "description",
		"smtp_config": testNotifierSMTPConfigInterface,
	}
	testNotifierWebhookConf = &managementClient.Notifier{
		Name:          "name",
		ClusterID:     "cluster_id",
		Description:   "description",
		WebhookConfig: testNotifierWebhookConfigConf,
	}
	testNotifierWebhookInterface = map[string]interface{}{
		"name":           "name",
		"cluster_id":     "cluster_id",
		"description":    "description",
		"webhook_config": testNotifierWebhookConfigInterface,
	}
	testNotifierWechatConf = &managementClient.Notifier{
		Name:         "name",
		ClusterID:    "cluster_id",
		Description:  "description",
		WechatConfig: testNotifierWechatConfigConf,
	}
	testNotifierWechatInterface = map[string]interface{}{
		"name":          "name",
		"cluster_id":    "cluster_id",
		"description":   "description",
		"wechat_config": testNotifierWechatConfigInterface,
	}
}

func TestFlattenNotifier(t *testing.T) {

	cases := []struct {
		Input          *managementClient.Notifier
		ExpectedOutput map[string]interface{}
	}{
		{
			testNotifierPagerdutyConf,
			testNotifierPagerdutyInterface,
		},
		{
			testNotifierSlackConf,
			testNotifierSlackInterface,
		},
		{
			testNotifierSMTPConf,
			testNotifierSMTPInterface,
		},
		{
			testNotifierWebhookConf,
			testNotifierWebhookInterface,
		},
		{
			testNotifierWechatConf,
			testNotifierWechatInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, notifierFields(), map[string]interface{}{})
		err := flattenNotifier(output, tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			expectedOutput[k] = output.Get(k)
		}
		if !reflect.DeepEqual(expectedOutput, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, expectedOutput)
		}
	}
}

func TestExpandNotifier(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *managementClient.Notifier
	}{
		{
			testNotifierPagerdutyInterface,
			testNotifierPagerdutyConf,
		},
		{
			testNotifierSlackInterface,
			testNotifierSlackConf,
		},
		{
			testNotifierSMTPInterface,
			testNotifierSMTPConf,
		},
		{
			testNotifierWebhookInterface,
			testNotifierWebhookConf,
		},
		{
			testNotifierWechatInterface,
			testNotifierWechatConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, notifierFields(), tc.Input)
		output, err := expandNotifier(inputResourceData)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
