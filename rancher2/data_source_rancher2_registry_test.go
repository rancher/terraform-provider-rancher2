package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccRancher2RegistryDataSourceType = "rancher2_registry"
)

var (
	testAccCheckRancher2RegistryProjectDataSourceConfig string
	testAccCheckRancher2RegistryNsDataSourceConfig      string
)

func init() {
	testAccCheckRancher2RegistryProjectDataSourceConfig = `
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
resource "rancher2_registry" "foo" {
  name = "foo"
  description = "Terraform registry acceptance test"
  project_id = "${rancher2_project.foo.id}"
  registries {
    address = "test.io"
    username = "user"
    password = "pass"
  }
}
data "` + testAccRancher2RegistryDataSourceType + `" "foo" {
  name = "${rancher2_registry.foo.name}"
  project_id = "${rancher2_project.foo.id}"
}
`
	testAccCheckRancher2RegistryNsDataSourceConfig = `
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
resource "rancher2_registry" "foo" {
  name = "foo"
  description = "Terraform registry acceptance test"
  project_id = "${rancher2_project.foo.id}"
  namespace_id = "${rancher2_namespace.foo.id}"
  registries {
    address = "test.io"
    username = "user"
    password = "pass"
  }
}
data "` + testAccRancher2RegistryDataSourceType + `" "foo" {
  name = "${rancher2_registry.foo.name}"
  project_id = "${rancher2_project.foo.id}"
  namespace_id = "${rancher2_namespace.foo.id}"
}
`
}

func TestAccRancher2RegistryDataSource_Project(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2RegistryProjectDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data."+testAccRancher2RegistryDataSourceType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr("data."+testAccRancher2RegistryDataSourceType+".foo", "description", "Terraform registry acceptance test"),
					resource.TestCheckResourceAttr("data."+testAccRancher2RegistryDataSourceType+".foo", "registries.0.address", "test.io"),
					resource.TestCheckResourceAttr("data."+testAccRancher2RegistryDataSourceType+".foo", "registries.0.username", "user"),
				),
			},
		},
	})
}

func TestAccRancher2RegistryDataSource_Namespaced(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2RegistryNsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data."+testAccRancher2RegistryDataSourceType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr("data."+testAccRancher2RegistryDataSourceType+".foo", "description", "Terraform registry acceptance test"),
					resource.TestCheckResourceAttr("data."+testAccRancher2RegistryDataSourceType+".foo", "registries.0.address", "test.io"),
					resource.TestCheckResourceAttr("data."+testAccRancher2RegistryDataSourceType+".foo", "registries.0.username", "user"),
				),
			},
		},
	})
}
