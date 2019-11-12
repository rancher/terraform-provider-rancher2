package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccRancher2NodeDriverDataSourceType = "rancher2_node_driver"
)

var (
	testAccCheckRancher2NodeDriverDataSourceConfig string
)

func init() {
	testAccCheckRancher2NodeDriverDataSourceConfig = `
data "` + testAccRancher2NodeDriverDataSourceType + `" "foo" {
  name = "amazonec2"
}
`
}

func TestAccRancher2NodeDriverDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2NodeDriverDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data."+testAccRancher2NodeDriverDataSourceType+".foo", "name", "amazonec2"),
					resource.TestCheckResourceAttr("data."+testAccRancher2NodeDriverDataSourceType+".foo", "id", "amazonec2"),
					resource.TestCheckResourceAttr("data."+testAccRancher2NodeDriverDataSourceType+".foo", "labels.cattle.io/creator", "norman"),
				),
			},
		},
	})
}
