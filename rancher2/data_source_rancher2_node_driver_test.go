package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRancher2NodeDriverDataSource(t *testing.T) {
	testAccCheckRancher2NodeDriverDataSourceConfig := `
data "` + testAccRancher2NodeDriverType + `" "foo" {
  name = "amazonec2"
}
`
	name := "data." + testAccRancher2NodeDriverType + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2NodeDriverDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "amazonec2"),
					resource.TestCheckResourceAttr(name, "id", "amazonec2"),
					resource.TestCheckResourceAttr(name, "labels.cattle.io/creator", "norman"),
				),
			},
		},
	})
}
