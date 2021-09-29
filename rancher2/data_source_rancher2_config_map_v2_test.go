package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRancher2ConfigMapV2DataSource_Cluster(t *testing.T) {
	testAccCheckRancher2ConfigMapV2ClusterDataSourceConfig := testAccRancher2ConfigMapV2Config + `
data "` + testAccRancher2ConfigMapV2Type + `" "foo" {
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  name = rancher2_config_map_v2.foo.name
  namespace = rancher2_config_map_v2.foo.namespace
}
`
	name := "data." + testAccRancher2ConfigMapV2Type + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2ConfigMapV2ClusterDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "namespace", "default"),
					resource.TestCheckResourceAttr(name, "data.param1", "true"),
					resource.TestCheckResourceAttr(name, "data.param2", "40000"),
				),
			},
		},
	})
}
