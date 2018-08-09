package cattle

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccCattleClusterDataSource_accessLog(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCattleClusterDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cattle_cluster.foo", "name", "default"),
				),
			},
		},
	})
}

// Testing owner parameter
const testAccCheckCattleClusterDataSourceConfig = `
data "cattle_cluster" "default" {
	name = "default"
}
`
