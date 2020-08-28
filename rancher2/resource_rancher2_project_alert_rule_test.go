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
	testAccRancher2ProjectAlertRuleType = "rancher2_project_alert_rule"
)

var (
	testAccRancher2ProjectAlertRule             string
	testAccRancher2ProjectAlertRuleUpdate       string
	testAccRancher2ProjectAlertRuleConfig       string
	testAccRancher2ProjectAlertRuleUpdateConfig string
)

func init() {
	testAccRancher2ProjectAlertRule = `
resource "` + testAccRancher2ProjectAlertRuleType + `" "foo" {
  project_id = rancher2_cluster_sync.testacc.default_project_id
  group_id = rancher2_project_alert_group.foo.id
  name = "foo"
  group_interval_seconds = 300
  repeat_interval_seconds = 3600
}
`
	testAccRancher2ProjectAlertRuleUpdate = `
resource "` + testAccRancher2ProjectAlertRuleType + `" "foo" {
  project_id = rancher2_cluster_sync.testacc.default_project_id
  group_id = rancher2_project_alert_group.foo.id
  name = "foo"
  group_interval_seconds = 600
  repeat_interval_seconds = 6000
}
`
	testAccRancher2ProjectAlertRuleConfig = testAccRancher2ProjectAlertGroupConfig + testAccRancher2ProjectAlertRule
	testAccRancher2ProjectAlertRuleUpdateConfig = testAccRancher2ProjectAlertGroupConfig + testAccRancher2ProjectAlertRuleUpdate
}

func TestAccRancher2ProjectAlertRule_basic(t *testing.T) {
	var ar *managementClient.ProjectAlertRule

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ProjectAlertRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2ProjectAlertRuleConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ProjectAlertRuleExists(testAccRancher2ProjectAlertRuleType+".foo", ar),
					resource.TestCheckResourceAttr(testAccRancher2ProjectAlertRuleType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectAlertRuleType+".foo", "group_interval_seconds", "300"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectAlertRuleType+".foo", "severity", alertRuleSeverityCritical),
					resource.TestCheckResourceAttr(testAccRancher2ProjectAlertRuleType+".foo", "repeat_interval_seconds", "3600"),
				),
			},
			{
				Config: testAccRancher2ProjectAlertRuleUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ProjectAlertRuleExists(testAccRancher2ProjectAlertRuleType+".foo", ar),
					resource.TestCheckResourceAttr(testAccRancher2ProjectAlertRuleType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectAlertRuleType+".foo", "group_interval_seconds", "600"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectAlertRuleType+".foo", "severity", alertRuleSeverityCritical),
					resource.TestCheckResourceAttr(testAccRancher2ProjectAlertRuleType+".foo", "repeat_interval_seconds", "6000"),
				),
			},
			{
				Config: testAccRancher2ProjectAlertRuleConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ProjectAlertRuleExists(testAccRancher2ProjectAlertRuleType+".foo", ar),
					resource.TestCheckResourceAttr(testAccRancher2ProjectAlertRuleType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectAlertRuleType+".foo", "group_interval_seconds", "300"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectAlertRuleType+".foo", "severity", alertRuleSeverityCritical),
					resource.TestCheckResourceAttr(testAccRancher2ProjectAlertRuleType+".foo", "repeat_interval_seconds", "3600"),
				),
			},
		},
	})
}

func TestAccRancher2ProjectAlertRule_disappears(t *testing.T) {
	var ar *managementClient.ProjectAlertRule

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ProjectAlertRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2ProjectAlertRuleConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ProjectAlertRuleExists(testAccRancher2ProjectAlertRuleType+".foo", ar),
					testAccRancher2ProjectAlertRuleDisappears(ar),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2ProjectAlertRuleDisappears(ar *managementClient.ProjectAlertRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2ProjectAlertRuleType {
				continue
			}
			client, err := testAccProvider.Meta().(*Config).ManagementClient()
			if err != nil {
				return err
			}

			ar, err = client.ProjectAlertRule.ByID(rs.Primary.ID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = client.ProjectAlertRule.Delete(ar)
			if err != nil {
				return fmt.Errorf("Error removing Project Alert Rule: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Pending:    []string{"removing"},
				Target:     []string{"removed"},
				Refresh:    projectAlertRuleStateRefreshFunc(client, rs.Primary.ID),
				Timeout:    10 * time.Minute,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf(
					"[ERROR] waiting for project alert rule (%s) to be removed: %s", rs.Primary.ID, waitErr)
			}
		}
		return nil

	}
}

func testAccCheckRancher2ProjectAlertRuleExists(n string, ar *managementClient.ProjectAlertRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No project alert rule ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		foundAr, err := client.ProjectAlertRule.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("Project Alert Rule not found")
			}
			return err
		}

		ar = foundAr

		return nil
	}
}

func testAccCheckRancher2ProjectAlertRuleDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2ProjectAlertRuleType {
			continue
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		_, err = client.ProjectAlertRule.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("Project Alert Rule still exists")
	}
	return nil
}
