package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRancher2SecretV2DataSource_Cluster(t *testing.T) {
	testAccCheckRancher2SecretV2ClusterDataSourceConfig := testAccRancher2SecretV2Config + `
data "` + testAccRancher2SecretV2Type + `" "foo" {
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  name = rancher2_secret_v2.foo.name
  namespace = rancher2_secret_v2.foo.namespace
}
`
	name := "data." + testAccRancher2SecretV2Type + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2SecretV2ClusterDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "namespace", "default"),
					resource.TestCheckResourceAttr(name, "data.username", "test"),
					resource.TestCheckResourceAttr(name, "data.password", "mypass"),
				),
			},
		},
	})
}
