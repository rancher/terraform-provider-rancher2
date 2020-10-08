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
	testAccRancher2ProjectLoggingType = "rancher2_project_logging"
)

var (
	testAccRancher2ProjectLoggingSyslog             string
	testAccRancher2ProjectLoggingSyslogUpdate       string
	testAccRancher2ProjectLoggingConfigSyslog       string
	testAccRancher2ProjectLoggingUpdateConfigSyslog string
)

func init() {
	testAccRancher2ProjectLoggingSyslog = `
resource "` + testAccRancher2ProjectLoggingType + `" "foo" {
  name = "foo"
  project_id = rancher2_cluster_sync.testacc.default_project_id
  kind = "syslog"
  syslog_config {
    endpoint = "192.168.1.1:514"
    protocol = "udp"
    severity = "notice"
    ssl_verify = false
  }
}
`
	testAccRancher2ProjectLoggingSyslogUpdate = `
resource "` + testAccRancher2ProjectLoggingType + `" "foo" {
  name = "foo-updated"
  project_id = rancher2_cluster_sync.testacc.default_project_id
  kind = "syslog"
  syslog_config {
    endpoint = "192.168.1.1:514"
    protocol = "udp"
    severity = "notice"
    ssl_verify = false
  }
}
`
	testAccRancher2ProjectLoggingConfigSyslog = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2ProjectLoggingSyslog
	testAccRancher2ProjectLoggingUpdateConfigSyslog = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2ProjectLoggingSyslogUpdate

}

func TestAccRancher2ProjectLogging_basic_syslog(t *testing.T) {
	var project *managementClient.ProjectLogging

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ProjectLoggingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2ProjectLoggingConfigSyslog,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ProjectLoggingExists(testAccRancher2ProjectLoggingType+".foo", project),
					resource.TestCheckResourceAttr(testAccRancher2ProjectLoggingType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectLoggingType+".foo", "kind", "syslog"),
				),
			},
			{
				Config: testAccRancher2ProjectLoggingUpdateConfigSyslog,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ProjectLoggingExists(testAccRancher2ProjectLoggingType+".foo", project),
					resource.TestCheckResourceAttr(testAccRancher2ProjectLoggingType+".foo", "name", "foo-updated"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectLoggingType+".foo", "kind", "syslog"),
				),
			},
			{
				Config: testAccRancher2ProjectLoggingConfigSyslog,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ProjectLoggingExists(testAccRancher2ProjectLoggingType+".foo", project),
					resource.TestCheckResourceAttr(testAccRancher2ProjectLoggingType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectLoggingType+".foo", "kind", "syslog"),
				),
			},
		},
	})
}

func TestAccRancher2ProjectLogging_disappears_syslog(t *testing.T) {
	var project *managementClient.ProjectLogging

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ProjectLoggingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2ProjectLoggingConfigSyslog,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ProjectLoggingExists(testAccRancher2ProjectLoggingType+".foo", project),
					testAccRancher2ProjectLoggingDisappears(project),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2ProjectLoggingDisappears(pro *managementClient.ProjectLogging) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2ProjectLoggingType {
				continue
			}
			client, err := testAccProvider.Meta().(*Config).ManagementClient()
			if err != nil {
				return err
			}

			pro, err = client.ProjectLogging.ByID(rs.Primary.ID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = client.ProjectLogging.Delete(pro)
			if err != nil {
				return fmt.Errorf("Error removing Project Logging: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Pending:    []string{"active"},
				Target:     []string{"removed"},
				Refresh:    clusterLoggingStateRefreshFunc(client, pro.ID),
				Timeout:    10 * time.Minute,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf(
					"[ERROR] waiting for project logging (%s) to be removed: %s", pro.ID, waitErr)
			}
		}
		return nil

	}
}

func testAccCheckRancher2ProjectLoggingExists(n string, pro *managementClient.ProjectLogging) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Project Logging ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		foundPro, err := client.ProjectLogging.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("Project Logging not found")
			}
			return err
		}

		pro = foundPro

		return nil
	}
}

func testAccCheckRancher2ProjectLoggingDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2ProjectLoggingType {
			continue
		}
		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		_, err = client.ProjectLogging.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("Project Logging still exists")
	}
	return nil
}
