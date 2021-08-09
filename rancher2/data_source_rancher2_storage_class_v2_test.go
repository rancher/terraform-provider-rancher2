package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRancher2StorageClassV2DataSource_Cluster(t *testing.T) {
	testAccCheckRancher2StorageClassV2ClusterDataSourceConfig := testAccRancher2StorageClassV2Config + `
data "` + testAccRancher2StorageClassV2Type + `" "foo" {
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  name = rancher2_storage_class_v2.foo.name
}
`
	name := "data." + testAccRancher2StorageClassV2Type + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2StorageClassV2ClusterDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "k8s_provisioner", "tfp.test.io/provisioner"),
					resource.TestCheckResourceAttr(name, "reclaim_policy", "Delete"),
				),
			},
		},
	})
}
