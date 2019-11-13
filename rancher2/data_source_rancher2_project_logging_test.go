package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccRancher2ProjectLoggingDataSourceType = "rancher2_project_logging"
)

var (
	testAccCheckRancher2ProjectLoggingDataSourceConfig string
)

func init() {
	testAccCheckRancher2ProjectLoggingDataSourceConfig = `
resource "rancher2_project" "foo" {
  name = "foo"
  cluster_id = "` + testAccRancher2ClusterID + `"
  description = "Terraform Project Logging acceptance test"
}

resource "rancher2_project_logging" "foo" {
  name = "foo"
  project_id = "${rancher2_project.foo.id}"
  kind = "syslog"
  syslog_config {
    endpoint = "192.168.1.1:514"
    protocol = "udp"
    severity = "notice"
    ssl_verify = false
  }
}

data "` + testAccRancher2ProjectLoggingDataSourceType + `" "foo" {
  project_id = "${rancher2_project_logging.foo.project_id}"
}
`
}

func TestAccRancher2ProjectLoggingDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2ProjectLoggingDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data."+testAccRancher2ProjectLoggingDataSourceType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr("data."+testAccRancher2ProjectLoggingDataSourceType+".foo", "kind", "syslog"),
					resource.TestCheckResourceAttr("data."+testAccRancher2ProjectLoggingDataSourceType+".foo", "syslog_config.0.endpoint", "192.168.1.1:514"),
				),
			},
		},
	})
}
