package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRancher2ProjectDataSource(t *testing.T) {
	testAccCheckRancher2ProjectDataSourceConfig := testAccCheckRancher2ClusterSyncTestacc + testAccRancher2Project + `
data "` + testAccRancher2ProjectType + `" "system" {
  name = "System"
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
}
`
	name := "data." + testAccRancher2ProjectType + ".system"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2ProjectDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "System"),
					resource.TestCheckResourceAttr(name, "cluster_id", testAccRancher2ClusterID),
				),
			},
		},
	})
}
