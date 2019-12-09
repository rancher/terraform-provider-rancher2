package rancher2

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	managementClient "github.com/rancher/types/client/management/v3"
)

const (
	testAccRancher2NotifierType = "rancher2_notifier"
)

var (
	testAccRancher2NotifierPagerdutyConfig         string
	testAccRancher2NotifierPagerdutyUpdateConfig   string
	testAccRancher2NotifierPagerdutyRecreateConfig string
	testAccRancher2NotifierSlackConfig             string
	testAccRancher2NotifierSlackUpdateConfig       string
	testAccRancher2NotifierSlackRecreateConfig     string
	testAccRancher2NotifierSMTPConfig              string
	testAccRancher2NotifierSMTPUpdateConfig        string
	testAccRancher2NotifierSMTPRecreateConfig      string
	testAccRancher2NotifierWebhookConfig           string
	testAccRancher2NotifierWebhookUpdateConfig     string
	testAccRancher2NotifierWebhookRecreateConfig   string
	testAccRancher2NotifierWechatConfig            string
	testAccRancher2NotifierWechatUpdateConfig      string
	testAccRancher2NotifierWechatRecreateConfig    string
)

func init() {
	testAccRancher2NotifierPagerdutyConfig = `
resource "rancher2_notifier" "foo" {
  name = "foo"
  cluster_id = "` + testAccRancher2ClusterID + `"
  send_resolved = "true"
  description = "Terraform notifier acceptance test"
  pagerduty_config {
    service_key = "XXXXXXXX"
    proxy_url = "http://proxy.test.io"
  }
}
`

	testAccRancher2NotifierPagerdutyUpdateConfig = `
resource "rancher2_notifier" "foo" {
  name = "foo"
  cluster_id = "` + testAccRancher2ClusterID + `"
  send_resolved = "false"
  description = "Terraform notifier acceptance test - updated"
  pagerduty_config {
    service_key = "YYYYYYYY"
    proxy_url = "http://proxy2.test.io"
  }
}
 `

	testAccRancher2NotifierPagerdutyRecreateConfig = `
resource "rancher2_notifier" "foo" {
  name = "foo"
  cluster_id = "` + testAccRancher2ClusterID + `"
  send_resolved = "true"
  description = "Terraform notifier acceptance test"
  pagerduty_config {
    service_key = "XXXXXXXX"
    proxy_url = "http://proxy.test.io"
  }
}
 `

	testAccRancher2NotifierSlackConfig = `
resource "rancher2_notifier" "foo" {
  name = "foo"
  cluster_id = "` + testAccRancher2ClusterID + `"
  send_resolved = "true"
  description = "Terraform notifier acceptance test"
  slack_config {
    default_recipient = "XXXXXXXX"
    url = "http://url.test.io"
    proxy_url = "http://proxy.test.io"
  }
}
`

	testAccRancher2NotifierSlackUpdateConfig = `
resource "rancher2_notifier" "foo" {
  name = "foo"
  cluster_id = "` + testAccRancher2ClusterID + `"
  send_resolved = "false"
  description = "Terraform notifier acceptance test - updated"
  slack_config {
    default_recipient = "YYYYYYYY"
    url = "http://url2.test.io"
    proxy_url = "http://proxy2.test.io"
  }
}
 `

	testAccRancher2NotifierSlackRecreateConfig = `
resource "rancher2_notifier" "foo" {
  name = "foo"
  cluster_id = "` + testAccRancher2ClusterID + `"
  send_resolved = "true"
  description = "Terraform notifier acceptance test"
  slack_config {
    default_recipient = "XXXXXXXX"
    url = "http://url.test.io"
    proxy_url = "http://proxy.test.io"
  }
}
 `

	testAccRancher2NotifierSMTPConfig = `
resource "rancher2_notifier" "foo" {
  name = "foo"
  cluster_id = "` + testAccRancher2ClusterID + `"
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

	testAccRancher2NotifierSMTPUpdateConfig = `
resource "rancher2_notifier" "foo" {
  name = "foo"
  cluster_id = "` + testAccRancher2ClusterID + `"
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

	testAccRancher2NotifierSMTPRecreateConfig = `
resource "rancher2_notifier" "foo" {
  name = "foo"
  cluster_id = "` + testAccRancher2ClusterID + `"
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

	testAccRancher2NotifierWebhookConfig = `
resource "rancher2_notifier" "foo" {
  name = "foo"
  cluster_id = "` + testAccRancher2ClusterID + `"
  send_resolved = "true"
  description = "Terraform notifier acceptance test"
  webhook_config {
    url = "http://url.test.io"
    proxy_url = "http://proxy.test.io"
  }
}
`

	testAccRancher2NotifierWebhookUpdateConfig = `
resource "rancher2_notifier" "foo" {
  name = "foo"
  cluster_id = "` + testAccRancher2ClusterID + `"
  send_resolved = "false"
  description = "Terraform notifier acceptance test - updated"
  webhook_config {
    url = "http://url2.test.io"
    proxy_url = "http://proxy2.test.io"
  }
}
 `

	testAccRancher2NotifierWebhookRecreateConfig = `
resource "rancher2_notifier" "foo" {
  name = "foo"
  cluster_id = "` + testAccRancher2ClusterID + `"
  send_resolved = "true"
  description = "Terraform notifier acceptance test"
  webhook_config {
    url = "http://url.test.io"
    proxy_url = "http://proxy.test.io"
  }
}
 `

	testAccRancher2NotifierWechatConfig = `
resource "rancher2_notifier" "foo" {
  name = "foo"
  cluster_id = "` + testAccRancher2ClusterID + `"
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

	testAccRancher2NotifierWechatUpdateConfig = `
resource "rancher2_notifier" "foo" {
  name = "foo"
  cluster_id = "` + testAccRancher2ClusterID + `"
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

	testAccRancher2NotifierWechatRecreateConfig = `
resource "rancher2_notifier" "foo" {
  name = "foo"
  cluster_id = "` + testAccRancher2ClusterID + `"
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

}

func TestAccRancher2Notifier_basic_Pagerduty(t *testing.T) {
	var notifier *managementClient.Notifier

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NotifierDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2NotifierPagerdutyConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "description", "Terraform notifier acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "send_resolved", "true"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "pagerduty_config.0.service_key", "XXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "pagerduty_config.0.proxy_url", "http://proxy.test.io"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NotifierPagerdutyUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "description", "Terraform notifier acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "send_resolved", "false"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "pagerduty_config.0.proxy_url", "http://proxy2.test.io"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NotifierPagerdutyRecreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "description", "Terraform notifier acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "send_resolved", "true"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "pagerduty_config.0.service_key", "XXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "pagerduty_config.0.proxy_url", "http://proxy.test.io"),
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
			resource.TestStep{
				Config: testAccRancher2NotifierPagerdutyConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo", notifier),
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
			resource.TestStep{
				Config: testAccRancher2NotifierSlackConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "description", "Terraform notifier acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "send_resolved", "true"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "slack_config.0.default_recipient", "XXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "slack_config.0.url", "http://url.test.io"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "slack_config.0.proxy_url", "http://proxy.test.io"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NotifierSlackUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "slack_config.0.default_recipient", "YYYYYYYY"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "slack_config.0.url", "http://url2.test.io"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "slack_config.0.proxy_url", "http://proxy2.test.io"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NotifierSlackRecreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "description", "Terraform notifier acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "send_resolved", "true"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "slack_config.0.default_recipient", "XXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "slack_config.0.url", "http://url.test.io"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "slack_config.0.proxy_url", "http://proxy.test.io"),
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
			resource.TestStep{
				Config: testAccRancher2NotifierSlackConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo", notifier),
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
			resource.TestStep{
				Config: testAccRancher2NotifierSMTPConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "description", "Terraform notifier acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "send_resolved", "true"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "smtp_config.0.default_recipient", "XXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "smtp_config.0.host", "host.test.io"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "smtp_config.0.sender", "sender@test.io"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NotifierSMTPUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "smtp_config.0.default_recipient", "YYYYYYYY"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "smtp_config.0.host", "host2.test.io"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "smtp_config.0.sender", "sender2@test.io"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NotifierSMTPRecreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "description", "Terraform notifier acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "send_resolved", "true"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "smtp_config.0.default_recipient", "XXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "smtp_config.0.host", "host.test.io"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "smtp_config.0.sender", "sender@test.io"),
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
			resource.TestStep{
				Config: testAccRancher2NotifierSMTPConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo", notifier),
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
			resource.TestStep{
				Config: testAccRancher2NotifierWebhookConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "description", "Terraform notifier acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "send_resolved", "true"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "webhook_config.0.url", "http://url.test.io"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "webhook_config.0.proxy_url", "http://proxy.test.io"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NotifierWebhookUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "webhook_config.0.url", "http://url2.test.io"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "webhook_config.0.proxy_url", "http://proxy2.test.io"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NotifierWebhookRecreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "description", "Terraform notifier acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "send_resolved", "true"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "webhook_config.0.url", "http://url.test.io"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "webhook_config.0.proxy_url", "http://proxy.test.io"),
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
			resource.TestStep{
				Config: testAccRancher2NotifierWebhookConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo", notifier),
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
			resource.TestStep{
				Config: testAccRancher2NotifierWechatConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "description", "Terraform notifier acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "send_resolved", "true"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "wechat_config.0.default_recipient", "XXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "wechat_config.0.secret", "XXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "wechat_config.0.proxy_url", "http://proxy.test.io"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NotifierWechatUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "wechat_config.0.default_recipient", "YYYYYYYY"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "wechat_config.0.secret", "YYYYYYYY"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "wechat_config.0.proxy_url", "http://proxy2.test.io"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NotifierWechatRecreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "description", "Terraform notifier acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "send_resolved", "true"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "wechat_config.0.default_recipient", "XXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "wechat_config.0.secret", "XXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "wechat_config.0.proxy_url", "http://proxy.test.io"),
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
			resource.TestStep{
				Config: testAccRancher2NotifierWechatConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo", notifier),
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

		obj, err := client.Notifier.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}

		if obj.Removed != "" {
			return nil
		}

		return fmt.Errorf("Notifier still exists")
	}
	return nil
}
