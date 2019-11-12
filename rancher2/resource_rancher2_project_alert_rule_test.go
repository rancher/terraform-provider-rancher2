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
	testAccRancher2ProjectAlertRuleType = "rancher2_project_alert_rule"
)

var (
	testAccRancher2ProjectAlertRuleGroup          string
	testAccRancher2ProjectAlertRuleConfig         string
	testAccRancher2ProjectAlertRuleUpdateConfig   string
	testAccRancher2ProjectAlertRuleRecreateConfig string
)

func init() {
	testAccRancher2ProjectAlertRuleGroup = `
resource "rancher2_project" "foo" {
  name = "foo"
  cluster_id = "` + testAccRancher2ClusterID + `"
  description = "Terraform project alert rule acceptance test"
  resource_quota {
    project_limit {
      limits_cpu = "2000m"
      limits_memory = "2000Mi"
      requests_storage = "2Gi"
    }
    namespace_default_limit {
      limits_cpu = "500m"
      limits_memory = "500Mi"
      requests_storage = "1Gi"
    }
  }
  container_resource_limit {
    limits_cpu = "20m"
    limits_memory = "20Mi"
    requests_cpu = "1m"
    requests_memory = "1Mi"
  }
}

resource "rancher2_project_alert_group" "foo" {
  name = "foo"
  description = "Terraform project alert rule acceptance test"
  project_id = "${rancher2_project.foo.id}"
  group_interval_seconds = 300
  repeat_interval_seconds = 3600
}
`

	testAccRancher2ProjectAlertRuleConfig = testAccRancher2ProjectAlertRuleGroup + `
resource "rancher2_project_alert_rule" "foo" {
  project_id = "${rancher2_project_alert_group.foo.project_id}"
  group_id = "${rancher2_project_alert_group.foo.id}"
  name = "foo"
  group_interval_seconds = 300
  repeat_interval_seconds = 3600
}
`

	testAccRancher2ProjectAlertRuleUpdateConfig = testAccRancher2ProjectAlertRuleGroup + `
resource "rancher2_project_alert_rule" "foo" {
  project_id = "${rancher2_project_alert_group.foo.project_id}"
  group_id = "${rancher2_project_alert_group.foo.id}"
  name = "foo"
  group_interval_seconds = 600
  repeat_interval_seconds = 6000
}
 `

	testAccRancher2ProjectAlertRuleRecreateConfig = testAccRancher2ProjectAlertRuleGroup + `
resource "rancher2_project_alert_rule" "foo" {
  project_id = "${rancher2_project_alert_group.foo.project_id}"
  group_id = "${rancher2_project_alert_group.foo.id}"
  name = "foo"
  group_interval_seconds = 300
  repeat_interval_seconds = 3600
}
 `
}

func TestAccRancher2ProjectAlertRule_basic(t *testing.T) {
	var ar *managementClient.ProjectAlertRule

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ProjectAlertRuleDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2ProjectAlertRuleConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ProjectAlertRuleExists(testAccRancher2ProjectAlertRuleType+".foo", ar),
					resource.TestCheckResourceAttr(testAccRancher2ProjectAlertRuleType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectAlertRuleType+".foo", "group_interval_seconds", "300"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectAlertRuleType+".foo", "severity", alertRuleSeverityCritical),
					resource.TestCheckResourceAttr(testAccRancher2ProjectAlertRuleType+".foo", "repeat_interval_seconds", "3600"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2ProjectAlertRuleUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ProjectAlertRuleExists(testAccRancher2ProjectAlertRuleType+".foo", ar),
					resource.TestCheckResourceAttr(testAccRancher2ProjectAlertRuleType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectAlertRuleType+".foo", "group_interval_seconds", "600"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectAlertRuleType+".foo", "severity", alertRuleSeverityCritical),
					resource.TestCheckResourceAttr(testAccRancher2ProjectAlertRuleType+".foo", "repeat_interval_seconds", "6000"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2ProjectAlertRuleRecreateConfig,
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
			resource.TestStep{
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

		obj, err := client.ProjectAlertRule.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
		if obj.Removed != "" {
			return nil
		}
		return fmt.Errorf("Project Alert Rule still exists")
	}
	return nil
}
