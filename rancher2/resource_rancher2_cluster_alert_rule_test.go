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
	testAccRancher2ClusterAlertRuleType = "rancher2_cluster_alert_rule"
)

var (
	testAccRancher2ClusterAlertRuleGroup        string
	testAccRancher2ClusterAlertRule             string
	testAccRancher2ClusterAlertRuleUpdate       string
	testAccRancher2ClusterAlertRuleConfig       string
	testAccRancher2ClusterAlertRuleUpdateConfig string
)

func init() {
	testAccRancher2ClusterAlertRule = `
resource "` + testAccRancher2ClusterAlertRuleType + `" "foo" {
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  group_id = rancher2_cluster_alert_group.foo.id
  name = "foo"
  group_interval_seconds = 300
  repeat_interval_seconds = 3600
}
`
	testAccRancher2ClusterAlertRuleUpdate = `
resource "` + testAccRancher2ClusterAlertRuleType + `" "foo" {
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  group_id = rancher2_cluster_alert_group.foo.id
  name = "foo"
  group_interval_seconds = 600
  repeat_interval_seconds = 6000
}
`
	testAccRancher2ClusterAlertRuleConfig = testAccRancher2ClusterAlertGroupConfig + testAccRancher2ClusterAlertRule
	testAccRancher2ClusterAlertRuleUpdateConfig = testAccRancher2ClusterAlertGroupConfig + testAccRancher2ClusterAlertRuleUpdate
}

func TestAccRancher2ClusterAlertRule_basic(t *testing.T) {
	var ar *managementClient.ClusterAlertRule

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ClusterAlertRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2ClusterAlertRuleConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterAlertRuleExists(testAccRancher2ClusterAlertRuleType+".foo", ar),
					resource.TestCheckResourceAttr(testAccRancher2ClusterAlertRuleType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterAlertRuleType+".foo", "group_interval_seconds", "300"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterAlertRuleType+".foo", "severity", alertRuleSeverityCritical),
					resource.TestCheckResourceAttr(testAccRancher2ClusterAlertRuleType+".foo", "repeat_interval_seconds", "3600"),
				),
			},
			{
				Config: testAccRancher2ClusterAlertRuleUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterAlertRuleExists(testAccRancher2ClusterAlertRuleType+".foo", ar),
					resource.TestCheckResourceAttr(testAccRancher2ClusterAlertRuleType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterAlertRuleType+".foo", "group_interval_seconds", "600"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterAlertRuleType+".foo", "severity", alertRuleSeverityCritical),
					resource.TestCheckResourceAttr(testAccRancher2ClusterAlertRuleType+".foo", "repeat_interval_seconds", "6000"),
				),
			},
			{
				Config: testAccRancher2ClusterAlertRuleConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterAlertRuleExists(testAccRancher2ClusterAlertRuleType+".foo", ar),
					resource.TestCheckResourceAttr(testAccRancher2ClusterAlertRuleType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterAlertRuleType+".foo", "group_interval_seconds", "300"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterAlertRuleType+".foo", "severity", alertRuleSeverityCritical),
					resource.TestCheckResourceAttr(testAccRancher2ClusterAlertRuleType+".foo", "repeat_interval_seconds", "3600"),
				),
			},
		},
	})
}

func TestAccRancher2ClusterAlertRule_disappears(t *testing.T) {
	var ar *managementClient.ClusterAlertRule

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ClusterAlertRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2ClusterAlertRuleConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterAlertRuleExists(testAccRancher2ClusterAlertRuleType+".foo", ar),
					testAccRancher2ClusterAlertRuleDisappears(ar),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2ClusterAlertRuleDisappears(ar *managementClient.ClusterAlertRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2ClusterAlertRuleType {
				continue
			}
			client, err := testAccProvider.Meta().(*Config).ManagementClient()
			if err != nil {
				return err
			}

			ar, err = client.ClusterAlertRule.ByID(rs.Primary.ID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = client.ClusterAlertRule.Delete(ar)
			if err != nil {
				return fmt.Errorf("Error removing Cluster Alert Rule: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Pending:    []string{"removing"},
				Target:     []string{"removed"},
				Refresh:    clusterAlertRuleStateRefreshFunc(client, rs.Primary.ID),
				Timeout:    10 * time.Minute,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf(
					"[ERROR] waiting for cluster alert rule (%s) to be removed: %s", rs.Primary.ID, waitErr)
			}
		}
		return nil

	}
}

func testAccCheckRancher2ClusterAlertRuleExists(n string, ar *managementClient.ClusterAlertRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No cluster alert rule ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		foundAr, err := client.ClusterAlertRule.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("Cluster Alert Rule not found")
			}
			return err
		}

		ar = foundAr

		return nil
	}
}

func testAccCheckRancher2ClusterAlertRuleDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2ClusterAlertRuleType {
			continue
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		_, err = client.ClusterAlertRule.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("Cluster Alert Rule still exists")
	}
	return nil
}
