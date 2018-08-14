package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccRancher2ClusterDataSource_accessLog(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2ClusterDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.rancher2_cluster.foo", "name", "local"),
				),
			},
		},
	})
}

// Testing owner parameter
const testAccCheckRancher2ClusterDataSourceConfig = `
data "rancher2_cluster" "foo" {
	name = "local"
}
`
