package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const (
	testAccRancher2CatalogDataSourceType = "rancher2_catalog"
)

var (
	testAccCheckRancher2CatalogDataSourceConfig string
)

func init() {
	testAccCheckRancher2CatalogDataSourceConfig = `
data "` + testAccRancher2CatalogDataSourceType + `" "library" {
  name = "library"
}
`
}

func TestAccRancher2CatalogDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2CatalogDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data."+testAccRancher2CatalogDataSourceType+".library", "name", "library"),
					resource.TestCheckResourceAttr("data."+testAccRancher2CatalogDataSourceType+".library", "url", "https://git.rancher.io/charts"),
				),
			},
		},
	})
}
