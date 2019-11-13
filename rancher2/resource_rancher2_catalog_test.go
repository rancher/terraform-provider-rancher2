package rancher2

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const (
	testAccRancher2CatalogType         = "rancher2_catalog"
	testAccRancher2CatalogGlobalConfig = `
resource "rancher2_catalog" "foo" {
  name = "foo"
  url = "http://foo.com:8080"
  description= "Terraform catalog acceptance test"
}
`

	testAccRancher2CatalogGlobalUpdateConfig = `
resource "rancher2_catalog" "foo" {
  name = "foo"
  url = "http://foo.updated.com:8080"
  description= "Terraform catalog acceptance test - updated"
}
 `

	testAccRancher2CatalogGlobalRecreateConfig = `
resource "rancher2_catalog" "foo" {
  name = "foo"
  url = "http://foo.com:8080"
  description= "Terraform catalog acceptance test"
}
 `
)

var (
	testAccRancher2CatalogClusterConfig         string
	testAccRancher2CatalogClusterUpdateConfig   string
	testAccRancher2CatalogClusterRecreateConfig string
	testAccRancher2CatalogProject               string
	testAccRancher2CatalogProjectConfig         string
	testAccRancher2CatalogProjectUpdateConfig   string
	testAccRancher2CatalogProjectRecreateConfig string
)

func init() {
	testAccRancher2CatalogClusterConfig = `
resource "rancher2_catalog" "foo" {
  name = "foo"
  url = "http://foo.com:8080"
  description= "Terraform catalog acceptance test"
  cluster_id = "` + testAccRancher2ClusterID + `"
  scope = "cluster"
}
`

	testAccRancher2CatalogClusterUpdateConfig = `
resource "rancher2_catalog" "foo" {
  name = "foo"
  url = "http://foo.updated.com:8080"
  description= "Terraform catalog acceptance test - updated"
  cluster_id = "` + testAccRancher2ClusterID + `"
  scope = "cluster"
}
 `

	testAccRancher2CatalogClusterRecreateConfig = `
resource "rancher2_catalog" "foo" {
  name = "foo"
  url = "http://foo.com:8080"
  description= "Terraform catalog acceptance test"
  cluster_id = "` + testAccRancher2ClusterID + `"
  scope = "cluster"
}
 `

	testAccRancher2CatalogProject = `
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
`

	testAccRancher2CatalogProjectConfig = testAccRancher2CatalogProject + `
resource "rancher2_catalog" "foo" {
  name = "foo"
  url = "http://foo.com:8080"
  description= "Terraform catalog acceptance test"
  project_id = "${rancher2_project.foo.id}"
  scope = "project"
}
`

	testAccRancher2CatalogProjectUpdateConfig = testAccRancher2CatalogProject + `
resource "rancher2_catalog" "foo" {
  name = "foo"
  url = "http://foo.updated.com:8080"
  description= "Terraform catalog acceptance test - updated"
  project_id = "${rancher2_project.foo.id}"
  scope = "project"
}
 `

	testAccRancher2CatalogProjectRecreateConfig = testAccRancher2CatalogProject + `
resource "rancher2_catalog" "foo" {
  name = "foo"
  url = "http://foo.com:8080"
  description= "Terraform catalog acceptance test"
  project_id = "${rancher2_project.foo.id}"
  scope = "project"
}
 `

}

func TestAccRancher2Catalog_basic_Global(t *testing.T) {
	var catalog interface{}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2CatalogDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2CatalogGlobalConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CatalogExists(testAccRancher2CatalogType+".foo", catalog),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "description", "Terraform catalog acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "url", "http://foo.com:8080"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "scope", "global"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2CatalogGlobalUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CatalogExists(testAccRancher2CatalogType+".foo", catalog),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "description", "Terraform catalog acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "url", "http://foo.updated.com:8080"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "scope", "global"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2CatalogGlobalRecreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CatalogExists(testAccRancher2CatalogType+".foo", catalog),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "description", "Terraform catalog acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "url", "http://foo.com:8080"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "scope", "global"),
				),
			},
		},
	})
}

func TestAccRancher2Catalog_disappears_Global(t *testing.T) {
	var catalog interface{}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2CatalogDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2CatalogGlobalConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CatalogExists(testAccRancher2CatalogType+".foo", catalog),
					testAccRancher2CatalogDisappears(catalog),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRancher2Catalog_basic_Cluster(t *testing.T) {
	var catalog interface{}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2CatalogDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2CatalogClusterConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CatalogExists(testAccRancher2CatalogType+".foo", catalog),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "description", "Terraform catalog acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "url", "http://foo.com:8080"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "scope", "cluster"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "cluster_id", testAccRancher2ClusterID),
				),
			},
			resource.TestStep{
				Config: testAccRancher2CatalogClusterUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CatalogExists(testAccRancher2CatalogType+".foo", catalog),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "description", "Terraform catalog acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "url", "http://foo.updated.com:8080"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "scope", "cluster"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "cluster_id", testAccRancher2ClusterID),
				),
			},
			resource.TestStep{
				Config: testAccRancher2CatalogClusterRecreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CatalogExists(testAccRancher2CatalogType+".foo", catalog),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "description", "Terraform catalog acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "url", "http://foo.com:8080"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "scope", "cluster"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "cluster_id", testAccRancher2ClusterID),
				),
			},
		},
	})
}

func TestAccRancher2Catalog_disappears_Cluster(t *testing.T) {
	var catalog interface{}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2CatalogDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2CatalogClusterConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CatalogExists(testAccRancher2CatalogType+".foo", catalog),
					testAccRancher2CatalogDisappears(catalog),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRancher2Catalog_basic_Project(t *testing.T) {
	var catalog interface{}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2CatalogDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2CatalogProjectConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CatalogExists(testAccRancher2CatalogType+".foo", catalog),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "description", "Terraform catalog acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "url", "http://foo.com:8080"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "scope", "project"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2CatalogProjectUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CatalogExists(testAccRancher2CatalogType+".foo", catalog),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "description", "Terraform catalog acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "url", "http://foo.updated.com:8080"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "scope", "project"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2CatalogProjectRecreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CatalogExists(testAccRancher2CatalogType+".foo", catalog),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "description", "Terraform catalog acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "url", "http://foo.com:8080"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo", "scope", "project"),
				),
			},
		},
	})
}

func TestAccRancher2Catalog_disappears_Project(t *testing.T) {
	var catalog interface{}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2CatalogDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2CatalogProjectConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CatalogExists(testAccRancher2CatalogType+".foo", catalog),
					testAccRancher2CatalogDisappears(catalog),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2CatalogDisappears(cat interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2CatalogType {
				continue
			}

			scope := rs.Primary.Attributes["scope"]
			id := rs.Primary.ID
			cat, err := testAccProvider.Meta().(*Config).GetCatalog(id, scope)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = testAccProvider.Meta().(*Config).DeleteCatalog(scope, cat)
			if err != nil {
				return fmt.Errorf("Error removing %s Catalog: %s", scope, err)
			}

			stateConf := &resource.StateChangeConf{
				Pending:    []string{"active"},
				Target:     []string{"removed"},
				Refresh:    catalogStateRefreshFunc(testAccProvider.Meta(), id, scope),
				Timeout:    10 * time.Minute,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf(
					"[ERROR] waiting for %s catalog (%s) to be removed: %s", scope, id, waitErr)
			}
		}
		return nil

	}
}

func testAccCheckRancher2CatalogExists(n string, cat interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No catalog ID is set")
		}

		scope := rs.Primary.Attributes["scope"]
		id := rs.Primary.ID
		foundReg, err := testAccProvider.Meta().(*Config).GetCatalog(id, scope)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}

		cat = foundReg

		return nil
	}
}

func testAccCheckRancher2CatalogDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2CatalogType {
			continue
		}
		scope := rs.Primary.Attributes["scope"]
		id := rs.Primary.ID
		_, err := testAccProvider.Meta().(*Config).GetCatalog(id, scope)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("Catalog still exists")
	}
	return nil
}
