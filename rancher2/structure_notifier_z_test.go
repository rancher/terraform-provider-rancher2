package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testNotifierDingtalkConf       *managementClient.Notifier
	testNotifierDingtalkInterface  map[string]interface{}
	testNotifierMSTeamsConf        *managementClient.Notifier
	testNotifierMSTeamsInterface   map[string]interface{}
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
	testNotifierDingtalkConf = &managementClient.Notifier{
		Name:           "name",
		ClusterID:      "cluster_id",
		Description:    "description",
		DingtalkConfig: testNotifierDingtalkConfigConf,
	}
	testNotifierDingtalkInterface = map[string]interface{}{
		"name":            "name",
		"cluster_id":      "cluster_id",
		"description":     "description",
		"dingtalk_config": testNotifierDingtalkConfigInterface,
	}
	testNotifierMSTeamsConf = &managementClient.Notifier{
		Name:          "name",
		ClusterID:     "cluster_id",
		Description:   "description",
		MSTeamsConfig: testNotifierMSTeamsConfigConf,
	}
	testNotifierMSTeamsInterface = map[string]interface{}{
		"name":           "name",
		"cluster_id":     "cluster_id",
		"description":    "description",
		"msteams_config": testNotifierMSTeamsConfigInterface,
	}
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
			testNotifierDingtalkConf,
			testNotifierDingtalkInterface,
		},
		{
			testNotifierMSTeamsConf,
			testNotifierMSTeamsInterface,
		},
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
			testNotifierDingtalkInterface,
			testNotifierDingtalkConf,
		},
		{
			testNotifierMSTeamsInterface,
			testNotifierMSTeamsConf,
		},
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
