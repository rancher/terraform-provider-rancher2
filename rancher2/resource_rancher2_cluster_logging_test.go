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
	testAccRancher2ClusterLoggingType = "rancher2_cluster_logging"
)

var (
	testAccRancher2ClusterLoggingSyslog             string
	testAccRancher2ClusterLoggingSyslogUpdate       string
	testAccRancher2ClusterLoggingSyslogConfig       string
	testAccRancher2ClusterLoggingSyslogUpdateConfig string
)

func init() {
	testAccRancher2ClusterLoggingSyslog = `
resource "` + testAccRancher2ClusterLoggingType + `" "foo" {
  name = "foo"
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  kind = "syslog"
  syslog_config {
    endpoint = "192.168.1.1:514"
    protocol = "udp"
    severity = "notice"
    ssl_verify = false
  }
}
`
	testAccRancher2ClusterLoggingSyslogUpdate = `
resource "` + testAccRancher2ClusterLoggingType + `" "foo" {
  name = "foo-updated"
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  kind = "syslog"
  syslog_config {
    endpoint = "192.168.1.1:514"
    protocol = "udp"
    severity = "notice"
    ssl_verify = false
  }
}
`
	testAccRancher2ClusterLoggingSyslogConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2ClusterLoggingSyslog
	testAccRancher2ClusterLoggingSyslogUpdateConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2ClusterLoggingSyslogUpdate
}

func TestAccRancher2ClusterLogging_basic_syslog(t *testing.T) {
	var cluster *managementClient.ClusterLogging

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ClusterLoggingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2ClusterLoggingSyslogConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterLoggingExists(testAccRancher2ClusterLoggingType+".foo", cluster),
					resource.TestCheckResourceAttr(testAccRancher2ClusterLoggingType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterLoggingType+".foo", "cluster_id", testAccRancher2ClusterID),
				),
			},
			{
				Config: testAccRancher2ClusterLoggingSyslogUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterLoggingExists(testAccRancher2ClusterLoggingType+".foo", cluster),
					resource.TestCheckResourceAttr(testAccRancher2ClusterLoggingType+".foo", "name", "foo-updated"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterLoggingType+".foo", "cluster_id", testAccRancher2ClusterID),
				),
			},
			{
				Config: testAccRancher2ClusterLoggingSyslogConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterLoggingExists(testAccRancher2ClusterLoggingType+".foo", cluster),
					resource.TestCheckResourceAttr(testAccRancher2ClusterLoggingType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterLoggingType+".foo", "cluster_id", testAccRancher2ClusterID),
				),
			},
		},
	})
}

func TestAccRancher2ClusterLogging_disappears_syslog(t *testing.T) {
	var cluster *managementClient.ClusterLogging

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ClusterLoggingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2ClusterLoggingSyslogConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterLoggingExists(testAccRancher2ClusterLoggingType+".foo", cluster),
					testAccRancher2ClusterLoggingDisappears(cluster),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2ClusterLoggingDisappears(clu *managementClient.ClusterLogging) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2ClusterLoggingType {
				continue
			}
			client, err := testAccProvider.Meta().(*Config).ManagementClient()
			if err != nil {
				return err
			}

			clu, err = client.ClusterLogging.ByID(rs.Primary.ID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = client.ClusterLogging.Delete(clu)
			if err != nil {
				return fmt.Errorf("Error removing Cluster Logging: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Pending:    []string{"active"},
				Target:     []string{"removed"},
				Refresh:    clusterLoggingStateRefreshFunc(client, clu.ID),
				Timeout:    10 * time.Minute,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf(
					"[ERROR] waiting for cluster logging (%s) to be removed: %s", clu.ID, waitErr)
			}
		}
		return nil

	}
}

func testAccCheckRancher2ClusterLoggingExists(n string, clu *managementClient.ClusterLogging) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cluster Logging ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		foundClu, err := client.ClusterLogging.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("Cluster Logging not found")
			}
			return err
		}

		clu = foundClu

		return nil
	}
}

func testAccCheckRancher2ClusterLoggingDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2ClusterLoggingType {
			continue
		}
		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		_, err = client.ClusterLogging.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("Cluster Logging still exists")
	}
	return nil
}
