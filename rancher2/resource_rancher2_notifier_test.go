package rancher2

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

const (
	testAccRancher2NotifierType = "rancher2_notifier"
)

var (
	testAccRancher2NotifierDingtalk              string
	testAccRancher2NotifierDingtalkUpdate        string
	testAccRancher2NotifierDingtalkConfig        string
	testAccRancher2NotifierDingtalkUpdateConfig  string
	testAccRancher2NotifierMSTeams               string
	testAccRancher2NotifierMSTeamsUpdate         string
	testAccRancher2NotifierMSTeamsConfig         string
	testAccRancher2NotifierMSTeamsUpdateConfig   string
	testAccRancher2NotifierPagerduty             string
	testAccRancher2NotifierPagerdutyUpdate       string
	testAccRancher2NotifierPagerdutyConfig       string
	testAccRancher2NotifierPagerdutyUpdateConfig string
	testAccRancher2NotifierSlack                 string
	testAccRancher2NotifierSlackUpdate           string
	testAccRancher2NotifierSlackConfig           string
	testAccRancher2NotifierSlackUpdateConfig     string
	testAccRancher2NotifierSMTP                  string
	testAccRancher2NotifierSMTPUpdate            string
	testAccRancher2NotifierSMTPConfig            string
	testAccRancher2NotifierSMTPUpdateConfig      string
	testAccRancher2NotifierWebhook               string
	testAccRancher2NotifierWebhookUpdate         string
	testAccRancher2NotifierWebhookConfig         string
	testAccRancher2NotifierWebhookUpdateConfig   string
	testAccRancher2NotifierWechat                string
	testAccRancher2NotifierWechatUpdate          string
	testAccRancher2NotifierWechatConfig          string
	testAccRancher2NotifierWechatUpdateConfig    string
)

func init() {
	testAccRancher2NotifierDingtalk = `
resource "` + testAccRancher2NotifierType + `" "foo-dingtalk" {
  name = "foo-dingtalk"
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  send_resolved = "true"
  description = "Terraform notifier acceptance test"
  dingtalk_config {
    url = "http://url.test.io"
    proxy_url = "http://proxy.test.io"
  }
}
`
	testAccRancher2NotifierDingtalkUpdate = `
resource "` + testAccRancher2NotifierType + `" "foo-dingtalk" {
  name = "foo-dingtalk"
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  send_resolved = "false"
  description = "Terraform notifier acceptance test - updated"
  dingtalk_config {
    url = "http://url2.test.io"
    proxy_url = "http://proxy2.test.io"
  }
}
`
	testAccRancher2NotifierDingtalkConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2NotifierDingtalk
	testAccRancher2NotifierDingtalkUpdateConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2NotifierDingtalkUpdate
	testAccRancher2NotifierMSTeams = `
resource "` + testAccRancher2NotifierType + `" "foo-msteams" {
  name = "foo-msteams"
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  send_resolved = "true"
  description = "Terraform notifier acceptance test"
  msteams_config {
    url = "http://url.test.io"
    proxy_url = "http://proxy.test.io"
  }
}
`
	testAccRancher2NotifierMSTeamsUpdate = `
resource "` + testAccRancher2NotifierType + `" "foo-msteams" {
  name = "foo-msteams"
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  send_resolved = "false"
  description = "Terraform notifier acceptance test - updated"
  msteams_config {
    url = "http://url2.test.io"
    proxy_url = "http://proxy2.test.io"
  }
}
`
	testAccRancher2NotifierMSTeamsConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2NotifierMSTeams
	testAccRancher2NotifierMSTeamsUpdateConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2NotifierMSTeamsUpdate
	testAccRancher2NotifierPagerduty = `
resource "` + testAccRancher2NotifierType + `" "foo-pagerduty" {
  name = "foo-pagerduty"
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  send_resolved = "true"
  description = "Terraform notifier acceptance test"
  pagerduty_config {
    service_key = "XXXXXXXX"
    proxy_url = "http://proxy.test.io"
  }
}
`
	testAccRancher2NotifierPagerdutyUpdate = `
resource "` + testAccRancher2NotifierType + `" "foo-pagerduty" {
  name = "foo-pagerduty"
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  send_resolved = "false"
  description = "Terraform notifier acceptance test - updated"
  pagerduty_config {
    service_key = "YYYYYYYY"
    proxy_url = "http://proxy2.test.io"
  }
}
`
	testAccRancher2NotifierPagerdutyConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2NotifierPagerduty
	testAccRancher2NotifierPagerdutyUpdateConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2NotifierPagerdutyUpdate
	testAccRancher2NotifierSlack = `
resource "` + testAccRancher2NotifierType + `" "foo-slack" {
  name = "foo-slack"
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  send_resolved = "true"
  description = "Terraform notifier acceptance test"
  slack_config {
    default_recipient = "XXXXXXXX"
    url = "http://url.test.io"
    proxy_url = "http://proxy.test.io"
  }
}
`
	testAccRancher2NotifierSlackUpdate = `
resource "` + testAccRancher2NotifierType + `" "foo-slack" {
  name = "foo-slack"
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  send_resolved = "false"
  description = "Terraform notifier acceptance test - updated"
  slack_config {
    default_recipient = "YYYYYYYY"
    url = "http://url2.test.io"
    proxy_url = "http://proxy2.test.io"
  }
}
`
	testAccRancher2NotifierSlackConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2NotifierSlack
	testAccRancher2NotifierSlackUpdateConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2NotifierSlackUpdate
	testAccRancher2NotifierSMTP = `
resource "` + testAccRancher2NotifierType + `" "foo-smtp" {
  name = "foo-smtp"
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  send_resolved = "true"
  description = "Terraform notifier acceptance test"
  smtp_config {
    default_recipient = "XXXXXXXX"
    host = "host.test.io"
    port = 25
    sender = "sender@test.io"
    tls = "true"
  }
}
`
	testAccRancher2NotifierSMTPUpdate = `
resource "` + testAccRancher2NotifierType + `" "foo-smtp" {
  name = "foo-smtp"
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  send_resolved = "false"
  description = "Terraform notifier acceptance test - updated"
  smtp_config {
    default_recipient = "YYYYYYYY"
    host = "host2.test.io"
    port = 25
    sender = "sender2@test.io"
    tls = "true"
  }
}
`
	testAccRancher2NotifierSMTPConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2NotifierSMTP
	testAccRancher2NotifierSMTPUpdateConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2NotifierSMTPUpdate
	testAccRancher2NotifierWebhook = `
resource "` + testAccRancher2NotifierType + `" "foo-webhook" {
  name = "foo-webhook"
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  send_resolved = "true"
  description = "Terraform notifier acceptance test"
  webhook_config {
    url = "http://url.test.io"
    proxy_url = "http://proxy.test.io"
  }
}
`
	testAccRancher2NotifierWebhookUpdate = `
resource "` + testAccRancher2NotifierType + `" "foo-webhook" {
  name = "foo-webhook"
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  send_resolved = "false"
  description = "Terraform notifier acceptance test - updated"
  webhook_config {
    url = "http://url2.test.io"
    proxy_url = "http://proxy2.test.io"
  }
}
`
	testAccRancher2NotifierWebhookConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2NotifierWebhook
	testAccRancher2NotifierWebhookUpdateConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2NotifierWebhookUpdate
	testAccRancher2NotifierWechat = `
resource "` + testAccRancher2NotifierType + `" "foo-wechat" {
  name = "foo-wechat"
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  send_resolved = "true"
  description = "Terraform notifier acceptance test"
  wechat_config {
    agent = "agent_id"
    corp = "corp_id"
    default_recipient = "XXXXXXXX"
    secret = "XXXXXXXX"
    proxy_url = "http://proxy.test.io"
  }
}
`
	testAccRancher2NotifierWechatUpdate = `
resource "` + testAccRancher2NotifierType + `" "foo-wechat" {
  name = "foo-wechat"
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  send_resolved = "false"
  description = "Terraform notifier acceptance test - updated"
  wechat_config {
    agent = "agent_id"
    corp = "corp_id"
    default_recipient = "YYYYYYYY"
    secret = "YYYYYYYY"
    proxy_url = "http://proxy2.test.io"
  }
}
`
	testAccRancher2NotifierWechatConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2NotifierWechat
	testAccRancher2NotifierWechatUpdateConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2NotifierWechatUpdate
}

func TestAccRancher2Notifier_basic_Dingtalk(t *testing.T) {
	var notifier *managementClient.Notifier

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NotifierDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2NotifierDingtalkConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo-dingtalk", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-dingtalk", "name", "foo-dingtalk"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-dingtalk", "description", "Terraform notifier acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-dingtalk", "send_resolved", "true"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-dingtalk", "dingtalk_config.0.url", "http://url.test.io"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-dingtalk", "dingtalk_config.0.proxy_url", "http://proxy.test.io"),
				),
			},
			{
				Config: testAccRancher2NotifierDingtalkUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo-dingtalk", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-dingtalk", "name", "foo-dingtalk"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-dingtalk", "description", "Terraform notifier acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-dingtalk", "send_resolved", "false"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-dingtalk", "dingtalk_config.0.url", "http://url2.test.io"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-dingtalk", "dingtalk_config.0.proxy_url", "http://proxy2.test.io"),
				),
			},
			{
				Config: testAccRancher2NotifierDingtalkConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo-dingtalk", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-dingtalk", "name", "foo-dingtalk"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-dingtalk", "description", "Terraform notifier acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-dingtalk", "send_resolved", "true"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-dingtalk", "dingtalk_config.0.url", "http://url.test.io"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-dingtalk", "dingtalk_config.0.proxy_url", "http://proxy.test.io"),
				),
			},
		},
	})
}

func TestAccRancher2Notifier_disappears_Dingtalk(t *testing.T) {
	var notifier *managementClient.Notifier

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NotifierDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2NotifierDingtalkConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo-dingtalk", notifier),
					testAccRancher2NotifierDisappears(notifier),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRancher2Notifier_basic_MSTeams(t *testing.T) {
	var notifier *managementClient.Notifier

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NotifierDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2NotifierMSTeamsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo-msteams", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-msteams", "name", "foo-msteams"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-msteams", "description", "Terraform notifier acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-msteams", "send_resolved", "true"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-msteams", "msteams_config.0.url", "http://url.test.io"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-msteams", "msteams_config.0.proxy_url", "http://proxy.test.io"),
				),
			},
			{
				Config: testAccRancher2NotifierMSTeamsUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo-msteams", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-msteams", "name", "foo-msteams"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-msteams", "description", "Terraform notifier acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-msteams", "send_resolved", "false"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-msteams", "msteams_config.0.url", "http://url2.test.io"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-msteams", "msteams_config.0.proxy_url", "http://proxy2.test.io"),
				),
			},
			{
				Config: testAccRancher2NotifierMSTeamsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo-msteams", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-msteams", "name", "foo-msteams"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-msteams", "description", "Terraform notifier acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-msteams", "send_resolved", "true"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-msteams", "msteams_config.0.url", "http://url.test.io"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-msteams", "msteams_config.0.proxy_url", "http://proxy.test.io"),
				),
			},
		},
	})
}

func TestAccRancher2Notifier_disappears_MSTeams(t *testing.T) {
	var notifier *managementClient.Notifier

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NotifierDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2NotifierMSTeamsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo-msteams", notifier),
					testAccRancher2NotifierDisappears(notifier),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRancher2Notifier_basic_Pagerduty(t *testing.T) {
	var notifier *managementClient.Notifier

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NotifierDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2NotifierPagerdutyConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo-pagerduty", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-pagerduty", "name", "foo-pagerduty"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-pagerduty", "description", "Terraform notifier acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-pagerduty", "send_resolved", "true"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-pagerduty", "pagerduty_config.0.service_key", "XXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-pagerduty", "pagerduty_config.0.proxy_url", "http://proxy.test.io"),
				),
			},
			{
				Config: testAccRancher2NotifierPagerdutyUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo-pagerduty", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-pagerduty", "name", "foo-pagerduty"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-pagerduty", "description", "Terraform notifier acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-pagerduty", "send_resolved", "false"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-pagerduty", "pagerduty_config.0.proxy_url", "http://proxy2.test.io"),
				),
			},
			{
				Config: testAccRancher2NotifierPagerdutyConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo-pagerduty", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-pagerduty", "name", "foo-pagerduty"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-pagerduty", "description", "Terraform notifier acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-pagerduty", "send_resolved", "true"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-pagerduty", "pagerduty_config.0.service_key", "XXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-pagerduty", "pagerduty_config.0.proxy_url", "http://proxy.test.io"),
				),
			},
		},
	})
}

func TestAccRancher2Notifier_disappears_Pagerduty(t *testing.T) {
	var notifier *managementClient.Notifier

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NotifierDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2NotifierPagerdutyConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo-pagerduty", notifier),
					testAccRancher2NotifierDisappears(notifier),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRancher2Notifier_basic_Slack(t *testing.T) {
	var notifier *managementClient.Notifier

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NotifierDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2NotifierSlackConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo-slack", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-slack", "name", "foo-slack"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-slack", "description", "Terraform notifier acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-slack", "send_resolved", "true"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-slack", "slack_config.0.default_recipient", "XXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-slack", "slack_config.0.url", "http://url.test.io"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-slack", "slack_config.0.proxy_url", "http://proxy.test.io"),
				),
			},
			{
				Config: testAccRancher2NotifierSlackUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo-slack", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-slack", "name", "foo-slack"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-slack", "slack_config.0.default_recipient", "YYYYYYYY"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-slack", "slack_config.0.url", "http://url2.test.io"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-slack", "slack_config.0.proxy_url", "http://proxy2.test.io"),
				),
			},
			{
				Config: testAccRancher2NotifierSlackConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo-slack", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-slack", "name", "foo-slack"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-slack", "description", "Terraform notifier acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-slack", "send_resolved", "true"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-slack", "slack_config.0.default_recipient", "XXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-slack", "slack_config.0.url", "http://url.test.io"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-slack", "slack_config.0.proxy_url", "http://proxy.test.io"),
				),
			},
		},
	})
}

func TestAccRancher2Notifier_disappears_Slack(t *testing.T) {
	var notifier *managementClient.Notifier

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NotifierDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2NotifierSlackConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo-slack", notifier),
					testAccRancher2NotifierDisappears(notifier),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRancher2Notifier_basic_SMTP(t *testing.T) {
	var notifier *managementClient.Notifier

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NotifierDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2NotifierSMTPConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo-smtp", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-smtp", "name", "foo-smtp"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-smtp", "description", "Terraform notifier acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-smtp", "send_resolved", "true"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-smtp", "smtp_config.0.default_recipient", "XXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-smtp", "smtp_config.0.host", "host.test.io"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-smtp", "smtp_config.0.sender", "sender@test.io"),
				),
			},
			{
				Config: testAccRancher2NotifierSMTPUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo-smtp", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-smtp", "name", "foo-smtp"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-smtp", "smtp_config.0.default_recipient", "YYYYYYYY"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-smtp", "smtp_config.0.host", "host2.test.io"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-smtp", "smtp_config.0.sender", "sender2@test.io"),
				),
			},
			{
				Config: testAccRancher2NotifierSMTPConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo-smtp", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-smtp", "name", "foo-smtp"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-smtp", "description", "Terraform notifier acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-smtp", "send_resolved", "true"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-smtp", "smtp_config.0.default_recipient", "XXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-smtp", "smtp_config.0.host", "host.test.io"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-smtp", "smtp_config.0.sender", "sender@test.io"),
				),
			},
		},
	})
}

func TestAccRancher2Notifier_disappears_SMTP(t *testing.T) {
	var notifier *managementClient.Notifier

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NotifierDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2NotifierSMTPConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo-smtp", notifier),
					testAccRancher2NotifierDisappears(notifier),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRancher2Notifier_basic_Webhook(t *testing.T) {
	var notifier *managementClient.Notifier

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NotifierDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2NotifierWebhookConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo-webhook", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-webhook", "name", "foo-webhook"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-webhook", "description", "Terraform notifier acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-webhook", "send_resolved", "true"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-webhook", "webhook_config.0.url", "http://url.test.io"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-webhook", "webhook_config.0.proxy_url", "http://proxy.test.io"),
				),
			},
			{
				Config: testAccRancher2NotifierWebhookUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo-webhook", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-webhook", "name", "foo-webhook"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-webhook", "webhook_config.0.url", "http://url2.test.io"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-webhook", "webhook_config.0.proxy_url", "http://proxy2.test.io"),
				),
			},
			{
				Config: testAccRancher2NotifierWebhookConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo-webhook", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-webhook", "name", "foo-webhook"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-webhook", "description", "Terraform notifier acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-webhook", "send_resolved", "true"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-webhook", "webhook_config.0.url", "http://url.test.io"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-webhook", "webhook_config.0.proxy_url", "http://proxy.test.io"),
				),
			},
		},
	})
}

func TestAccRancher2Notifier_disappears_Webhook(t *testing.T) {
	var notifier *managementClient.Notifier

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NotifierDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2NotifierWebhookConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo-webhook", notifier),
					testAccRancher2NotifierDisappears(notifier),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRancher2Notifier_basic_Wechat(t *testing.T) {
	var notifier *managementClient.Notifier

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NotifierDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2NotifierWechatConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo-wechat", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-wechat", "name", "foo-wechat"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-wechat", "description", "Terraform notifier acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-wechat", "send_resolved", "true"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-wechat", "wechat_config.0.default_recipient", "XXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-wechat", "wechat_config.0.secret", "XXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-wechat", "wechat_config.0.proxy_url", "http://proxy.test.io"),
				),
			},
			{
				Config: testAccRancher2NotifierWechatUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo-wechat", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-wechat", "name", "foo-wechat"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-wechat", "wechat_config.0.default_recipient", "YYYYYYYY"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-wechat", "wechat_config.0.secret", "YYYYYYYY"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-wechat", "wechat_config.0.proxy_url", "http://proxy2.test.io"),
				),
			},
			{
				Config: testAccRancher2NotifierWechatConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo-wechat", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-wechat", "name", "foo-wechat"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-wechat", "description", "Terraform notifier acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-wechat", "send_resolved", "true"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-wechat", "wechat_config.0.default_recipient", "XXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-wechat", "wechat_config.0.secret", "XXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo-wechat", "wechat_config.0.proxy_url", "http://proxy.test.io"),
				),
			},
		},
	})
}

func TestAccRancher2Notifier_disappears_Wechat(t *testing.T) {
	var notifier *managementClient.Notifier

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NotifierDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2NotifierWechatConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo-wechat", notifier),
					testAccRancher2NotifierDisappears(notifier),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2NotifierDisappears(notifier *managementClient.Notifier) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2NotifierType {
				continue
			}

			client, err := testAccProvider.Meta().(*Config).ManagementClient()
			if err != nil {
				return err
			}

			notifier, err = client.Notifier.ByID(rs.Primary.ID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = client.Notifier.Delete(notifier)
			if err != nil {
				return fmt.Errorf("Error removing Notifier: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Pending:    []string{"active"},
				Target:     []string{"removed"},
				Refresh:    notifierStateRefreshFunc(client, rs.Primary.ID),
				Timeout:    10 * time.Minute,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf(
					"[ERROR] waiting for notifier (%s) to be removed: %s", rs.Primary.ID, waitErr)
			}
		}
		return nil

	}
}

func testAccCheckRancher2NotifierExists(n string, notifier *managementClient.Notifier) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No notifier ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		foundNot, err := client.Notifier.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("Notifier not found")
			}
			return err
		}

		notifier = foundNot

		return nil
	}
}

func testAccCheckRancher2NotifierDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2NotifierType {
			continue
		}
		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		_, err = client.Notifier.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}

		return fmt.Errorf("Notifier still exists")
	}
	return nil
}
