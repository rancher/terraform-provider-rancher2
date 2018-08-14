package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccCattleSettingDataSource_accessLog(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCattleSettingDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.rancher2_setting.server-image", "value", "rancher/rancher"),
				),
			},
		},
	})
}

// Testing owner parameter
const testAccCheckCattleSettingDataSourceConfig = `
data "rancher2_setting" "server-image" {
	name = "server-image"
}
`
