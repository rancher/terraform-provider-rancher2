package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccRancher2ProjectDataSourceType = "rancher2_project"
)

var (
	testAccCheckRancher2ProjectDataSourceConfig string
)

func init() {
	testAccCheckRancher2ProjectDataSourceConfig = `
data "` + testAccRancher2ProjectDataSourceType + `" "system" {
  name = "System"
  cluster_id = "` + testAccRancher2ClusterID + `"
}
`
}

func TestAccRancher2ProjectDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2ProjectDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data."+testAccRancher2ProjectDataSourceType+".system", "name", "System"),
					resource.TestCheckResourceAttr("data."+testAccRancher2ProjectDataSourceType+".system", "cluster_id", testAccRancher2ClusterID),
				),
			},
		},
	})
}
