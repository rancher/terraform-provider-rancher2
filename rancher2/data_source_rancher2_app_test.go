package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccRancher2AppDataSourceType = "rancher2_app"
)

var (
	testAccCheckRancher2AppDataSourceConfig string
)

func init() {
	testAccCheckRancher2AppDataSourceConfig = `
resource "rancher2_project" "foo" {
  name = "foo"
  cluster_id = "` + testAccRancher2ClusterID + `"
  description = "Terraform project acceptance test"
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
resource "rancher2_app" "foo" {
  catalog_name = "library"
  name = "foo"
  description = "Terraform app acceptance test"
  project_id = "${rancher2_project.foo.id}"
  template_name = "docker-registry"
  template_version = "1.8.1"
  target_namespace = "${rancher2_namespace.foo.name}"
  answers = {
    "ingress_host" = "test.xip.io"
  }
}
data "` + testAccRancher2AppDataSourceType + `" "foo" {
  name = "${rancher2_app.foo.name}"
  project_id = "${rancher2_project.foo.id}"
}
`
}

func TestAccRancher2AppDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2AppDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data."+testAccRancher2AppDataSourceType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr("data."+testAccRancher2AppDataSourceType+".foo", "description", "Terraform app acceptance test"),
					resource.TestCheckResourceAttr("data."+testAccRancher2AppDataSourceType+".foo", "target_namespace", "foo"),
					resource.TestCheckResourceAttr("data."+testAccRancher2AppDataSourceType+".foo", "external_id", "catalog://?catalog=library&template=docker-registry&version=1.8.1"),
					resource.TestCheckResourceAttr("data."+testAccRancher2AppDataSourceType+".foo", "answers.ingress_host", "test.xip.io"),
				),
			},
		},
	})
}
