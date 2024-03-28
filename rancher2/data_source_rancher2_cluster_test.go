package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccRancher2ClusterDataSourceType = "rancher2_cluster"
)

var (
	testAccCheckRancher2ClusterDataSourceConfig string
)

func init() {
	testAccCheckRancher2ClusterDataSourceConfig = `
resource "rancher2_cluster" "foo" {
  name = "foo"
  description = "Terraform custom cluster acceptance test"
  rke_config {
    network {
      plugin = "canal"
    }
    services {
      etcd {
        creation = "6h"
        retention = "24h"
      }
    }
  }
}
data "` + testAccRancher2ClusterDataSourceType + `" "foo" {
  name = rancher2_cluster.foo.name
}
`
}

func TestAccRancher2ClusterDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2ClusterDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data."+testAccRancher2ClusterDataSourceType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr("data."+testAccRancher2ClusterDataSourceType+".foo", "description", "Terraform custom cluster acceptance test"),
					resource.TestCheckResourceAttr("data."+testAccRancher2ClusterDataSourceType+".foo", "rke_config.0.network.0.plugin", "canal"),
					resource.TestCheckResourceAttr("data."+testAccRancher2ClusterDataSourceType+".foo", "rke_config.0.services.0.etcd.0.creation", "6h"),
					resource.TestCheckResourceAttr("data."+testAccRancher2ClusterDataSourceType+".foo", "labels.cattle.io/creator", "norman"),
				),
			},
		},
	})
}
