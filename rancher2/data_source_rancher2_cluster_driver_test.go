package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRancher2ClusterDriverDataSource(t *testing.T) {
	testAccCheckRancher2ClusterDriverDataSourceConfig := `
data "` + testAccRancher2ClusterDriverType + `" "foo" {
  name = "amazonElasticContainerService"
}
`
	name := "data." + testAccRancher2ClusterDriverType + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2ClusterDriverDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "amazonElasticContainerService"),
					resource.TestCheckResourceAttr(name, "id", "amazonelasticcontainerservice"),
					resource.TestCheckResourceAttr(name, "labels.cattle.io/creator", "norman"),
				),
			},
		},
	})
}
