package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccRancher2ClusterDriverDataSourceType = "rancher2_cluster_driver"
)

var (
	testAccCheckRancher2ClusterDriverDataSourceConfig string
)

func init() {
	testAccCheckRancher2ClusterDriverDataSourceConfig = `
data "` + testAccRancher2ClusterDriverDataSourceType + `" "foo" {
  name = "amazonElasticContainerService"
}
`
}

func TestAccRancher2ClusterDriverDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2ClusterDriverDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data."+testAccRancher2ClusterDriverDataSourceType+".foo", "name", "amazonElasticContainerService"),
					resource.TestCheckResourceAttr("data."+testAccRancher2ClusterDriverDataSourceType+".foo", "id", "amazonelasticcontainerservice"),
					resource.TestCheckResourceAttr("data."+testAccRancher2ClusterDriverDataSourceType+".foo", "labels.cattle.io/creator", "norman"),
				),
			},
		},
	})
}
