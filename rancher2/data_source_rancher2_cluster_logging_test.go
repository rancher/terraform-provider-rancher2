package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccRancher2ClusterLoggingDataSourceType = "rancher2_cluster_logging"
)

var (
	testAccCheckRancher2ClusterLoggingDataSourceConfig string
)

func init() {
	testAccCheckRancher2ClusterLoggingDataSourceConfig = `
resource "rancher2_cluster_logging" "foo" {
  name = "foo"
  cluster_id = "` + testAccRancher2ClusterID + `"
  kind = "syslog"
  syslog_config {
    endpoint = "192.168.1.1:514"
    protocol = "udp"
    severity = "notice"
    ssl_verify = false
  }
}

data "` + testAccRancher2ClusterLoggingDataSourceType + `" "foo" {
  cluster_id = "${rancher2_cluster_logging.foo.cluster_id}"
}
`
}

func TestAccRancher2ClusterLoggingDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2ClusterLoggingDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data."+testAccRancher2ClusterLoggingDataSourceType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr("data."+testAccRancher2ClusterLoggingDataSourceType+".foo", "kind", "syslog"),
					resource.TestCheckResourceAttr("data."+testAccRancher2ClusterLoggingDataSourceType+".foo", "syslog_config.0.endpoint", "192.168.1.1:514"),
				),
			},
		},
	})
}
