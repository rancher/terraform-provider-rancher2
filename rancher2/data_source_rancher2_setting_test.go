package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccCheckRancher2SettingDataSourceConfig = `
data "` + testAccRancher2SettingType + `" "server-image" {
	name = "server-image"
}
`
)

func TestAccRancher2SettingDataSource_accessLog(t *testing.T) {
	name := "data." + testAccRancher2SettingType + ".server-image"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2SettingDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "value", "rancher/rancher"),
				),
			},
		},
	})
}
