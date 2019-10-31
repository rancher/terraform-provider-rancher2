package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRancher2SettingDataSource_accessLog(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2SettingDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.rancher2_setting.server-image", "value", "rancher/rancher"),
				),
			},
		},
	})
}

// Testing owner parameter
const testAccCheckRancher2SettingDataSourceConfig = `
data "rancher2_setting" "server-image" {
	name = "server-image"
}
`
