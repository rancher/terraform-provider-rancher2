package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRancher2CatalogV2DataSource_Cluster(t *testing.T) {
	testAccCheckRancher2CatalogV2ClusterDataSourceConfig := testAccRancher2CatalogV2Config + `
data "` + testAccRancher2CatalogV2Type + `" "foo" {
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  name = rancher2_catalog_v2.foo.name
}
`
	name := "data." + testAccRancher2CatalogV2Type + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2CatalogV2ClusterDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "git_repo", "https://git.rancher.io/charts"),
					resource.TestCheckResourceAttr(name, "git_branch", "dev-v2.5"),
				),
			},
		},
	})
}
