package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccRancher2SecretDataSourceType = "rancher2_secret"
)

var (
	testAccCheckRancher2SecretProjectDataSourceConfig string
	testAccCheckRancher2SecretNsDataSourceConfig      string
)

func init() {
	testAccCheckRancher2SecretProjectDataSourceConfig = `
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
resource "rancher2_secret" "foo" {
  name = "foo"
  description = "Terraform registry acceptance test"
  project_id = "${rancher2_project.foo.id}"
  data = {
    address = "dGVzdC5pbw=="
    password = "cGFzcw=="
    username = "dXNlcg=="
  }
}
data "` + testAccRancher2SecretDataSourceType + `" "foo" {
  name = "${rancher2_secret.foo.name}"
  project_id = "${rancher2_project.foo.id}"
}
`
	testAccCheckRancher2SecretNsDataSourceConfig = `
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
resource "rancher2_secret" "foo" {
  name = "foo"
  description = "Terraform registry acceptance test"
  project_id = "${rancher2_project.foo.id}"
  namespace_id = "${rancher2_namespace.foo.id}"
  data = {
    address = "dGVzdC5pbw=="
    password = "cGFzcw=="
    username = "dXNlcg=="
  }
}
data "` + testAccRancher2SecretDataSourceType + `" "foo" {
  name = "${rancher2_secret.foo.name}"
  project_id = "${rancher2_project.foo.id}"
  namespace_id = "${rancher2_namespace.foo.id}"
}
`
}

func TestAccRancher2SecretDataSource_Project(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2SecretProjectDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data."+testAccRancher2SecretDataSourceType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr("data."+testAccRancher2SecretDataSourceType+".foo", "description", "Terraform registry acceptance test"),
					resource.TestCheckResourceAttr("data."+testAccRancher2SecretDataSourceType+".foo", "data.address", "dGVzdC5pbw=="),
					resource.TestCheckResourceAttr("data."+testAccRancher2SecretDataSourceType+".foo", "data.username", "dXNlcg=="),
				),
			},
		},
	})
}

func TestAccRancher2SecretDataSource_Namespaced(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2SecretNsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data."+testAccRancher2SecretDataSourceType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr("data."+testAccRancher2SecretDataSourceType+".foo", "description", "Terraform registry acceptance test"),
					resource.TestCheckResourceAttr("data."+testAccRancher2SecretDataSourceType+".foo", "data.address", "dGVzdC5pbw=="),
					resource.TestCheckResourceAttr("data."+testAccRancher2SecretDataSourceType+".foo", "data.username", "dXNlcg=="),
				),
			},
		},
	})
}
