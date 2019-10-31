package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccRancher2NodePoolDataSourceType = "rancher2_node_pool"
)

var (
	testAccCheckRancher2NodePoolDataSourceConfig string
)

func init() {
	testAccCheckRancher2NodePoolDataSourceConfig = `
resource "rancher2_cluster" "foo" {
  name = "foo-custom"
  description = "Terraform node pool cluster acceptance test"
  rke_config {
    network {
      plugin = "canal"
    }
  }
}

resource "rancher2_cloud_credential" "foo" {
  name = "foo"
  description= "Terraform cloudCredential acceptance test"
  amazonec2_credential_config {
	access_key = "XXXXXXXXXXXXXXXXXXXX"
	secret_key = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  }
}

resource "rancher2_node_template" "foo" {
  name = "foo"
  description = "Terraform node pool acceptance test"
  cloud_credential_id = "${rancher2_cloud_credential.foo.id}"
  amazonec2_config {
	ami =  "ami-XXXXXXXXXXXXXXX"
	region = "XX-west-1"
	security_group = ["XXXXXXXX"]
	subnet_id = "subnet-XXXXXXXX"
	vpc_id = "vpc-XXXXXXXX"
	zone = "a"
  }
}

resource "rancher2_node_pool" "foo" {
  cluster_id =  "${rancher2_cluster.foo.id}"
  name = "foo"
  hostname_prefix =  "foo-cluster-0"
  node_template_id = "${rancher2_node_template.foo.id}"
  quantity = 1
  control_plane = true
  etcd = true
  worker = true
}

data "` + testAccRancher2NodePoolDataSourceType + `" "foo" {
  name = "${rancher2_node_pool.foo.name}"
  cluster_id = "${rancher2_node_pool.foo.cluster_id}"
}
`
}

func TestAccRancher2NodePoolDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2NodePoolDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data."+testAccRancher2NodePoolDataSourceType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr("data."+testAccRancher2NodePoolDataSourceType+".foo", "hostname_prefix", "foo-cluster-0"),
					resource.TestCheckResourceAttr("data."+testAccRancher2NodePoolDataSourceType+".foo", "control_plane", "true"),
					resource.TestCheckResourceAttr("data."+testAccRancher2NodePoolDataSourceType+".foo", "labels.cattle.io/creator", "norman"),
				),
			},
		},
	})
}
