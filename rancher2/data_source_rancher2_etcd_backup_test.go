package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRancher2EtcdBackupDataSource(t *testing.T) {
	testAccCheckRancher2EtcdBackupDataSourceConfig := testAccRancher2EtcdBackupConfig + `
data "` + testAccRancher2EtcdBackupType + `" "foo" {
  name = rancher2_etcd_backup.foo.name
  cluster_id = rancher2_etcd_backup.foo.cluster_id
}
`
	name := "data." + testAccRancher2EtcdBackupType + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2EtcdBackupDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "manual", "true"),
					resource.TestCheckResourceAttr(name, "labels.cattle.io/creator", "norman"),
				),
			},
		},
	})
}
