package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRancher2ClusterLoggingDataSource(t *testing.T) {
	testAccCheckRancher2ClusterLoggingDataSourceConfig := testAccRancher2ClusterLoggingSyslogConfig + `
data "` + testAccRancher2ClusterLoggingType + `" "foo" {
  cluster_id = rancher2_cluster_logging.foo.cluster_id
}
`
	name := "data." + testAccRancher2ClusterLoggingType + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2ClusterLoggingDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "kind", "syslog"),
					resource.TestCheckResourceAttr(name, "syslog_config.0.endpoint", "192.168.1.1:514"),
				),
			},
		},
	})
}
