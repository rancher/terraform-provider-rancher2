package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccRancher2NamespaceDataSourceType = "rancher2_namespace"
)

var (
	testAccCheckRancher2NamespaceDataSourceConfig string
)

func init() {
	testAccCheckRancher2NamespaceDataSourceConfig = `
resource "rancher2_project" "foo" {
  name = "foo"
  cluster_id = "` + testAccRancher2ClusterID + `"
  description = "Terraform namespace acceptance test"
  resource_quota {
    project_limit {
      limits_cpu = "2000m"
      limits_memory = "2000Mi"
      requests_storage = "2Gi"
    }
    namespace_default_limit {
      limits_cpu = "500m"
      limits_memory = "500Mi"
      requests_storage = "1Gi"
    }
  }
}
resource "rancher2_namespace" "foo" {
  name = "foo"
  description = "Terraform namespace acceptance test"
  project_id = "${rancher2_project.foo.id}"
  resource_quota {
    limit {
      limits_cpu = "100m"
      limits_memory = "100Mi"
      requests_storage = "1Gi"
    }
  }
}
data "` + testAccRancher2NamespaceDataSourceType + `" "foo" {
  name = "${rancher2_namespace.foo.name}"
  project_id = "${rancher2_namespace.foo.project_id}"
}
`
}

func TestAccRancher2NamespaceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2NamespaceDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data."+testAccRancher2NamespaceDataSourceType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr("data."+testAccRancher2NamespaceDataSourceType+".foo", "description", "Terraform namespace acceptance test"),
					resource.TestCheckResourceAttr("data."+testAccRancher2NamespaceDataSourceType+".foo", "labels.cattle.io/creator", "norman"),
				),
			},
		},
	})
}
