package rancher2

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const (
	testAccRancher2CatalogType   = "rancher2_catalog"
	testAccRancher2CatalogGlobal = `
resource "` + testAccRancher2CatalogType + `" "foo-global" {
  name = "foo-global"
  url = "http://foo.com:8080"
  description= "Terraform catalog acceptance test"
  version = "helm_v3"
  annotations = {
    "testacc.terraform.io/test" = "true"
  }
  labels = {
    "testacc.terraform.io/test" = "true"
  }
}
`
	testAccRancher2CatalogGlobalUpdate = `
resource "` + testAccRancher2CatalogType + `" "foo-global" {
  name = "foo-global"
  url = "http://foo.updated.com:8080"
  description= "Terraform catalog acceptance test - updated"
  version = "helm_v3"
  annotations = {
    "testacc.terraform.io/test" = "false"
  }
  labels = {
    "testacc.terraform.io/test" = "false"
  }
}
 `
)

var (
	testAccRancher2CatalogCluster             string
	testAccRancher2CatalogClusterUpdate       string
	testAccRancher2CatalogClusterConfig       string
	testAccRancher2CatalogClusterUpdateConfig string
	testAccRancher2CatalogProject             string
	testAccRancher2CatalogProjectUpdate       string
	testAccRancher2CatalogProjectConfig       string
	testAccRancher2CatalogProjectUpdateConfig string
)

func init() {
	testAccRancher2CatalogCluster = `
resource "` + testAccRancher2CatalogType + `" "foo-cluster" {
  name = "foo-cluster"
  url = "http://foo.com:8080"
  description= "Terraform catalog acceptance test"
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  scope = "cluster"
  version = "helm_v2"
}
`
	testAccRancher2CatalogClusterUpdate = `
resource "` + testAccRancher2CatalogType + `" "foo-cluster" {
  name = "foo-cluster"
  url = "http://foo.updated.com:8080"
  description= "Terraform catalog acceptance test - updated"
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  scope = "cluster"
  version = "helm_v2"
}
 `
	testAccRancher2CatalogClusterConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2CatalogCluster
	testAccRancher2CatalogClusterUpdateConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2CatalogClusterUpdate
	testAccRancher2CatalogProject = `
resource "` + testAccRancher2CatalogType + `" "foo-project" {
  name = "foo-project"
  url = "http://foo.com:8080"
  description= "Terraform catalog acceptance test"
  project_id = rancher2_cluster_sync.testacc.default_project_id
  scope = "project"
  version = "helm_v2"
}
`
	testAccRancher2CatalogProjectUpdate = `
resource "` + testAccRancher2CatalogType + `" "foo-project" {
  name = "foo-project"
  url = "http://foo.updated.com:8080"
  description= "Terraform catalog acceptance test - updated"
  project_id = rancher2_cluster_sync.testacc.default_project_id
  scope = "project"
  version = "helm_v2"
}
`
	testAccRancher2CatalogProjectConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2CatalogProject
	testAccRancher2CatalogProjectUpdateConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2CatalogProjectUpdate
}

func TestAccRancher2Catalog_basic_Global(t *testing.T) {
	var catalog interface{}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2CatalogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2CatalogGlobal,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CatalogExists(testAccRancher2CatalogType+".foo-global", catalog),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-global", "name", "foo-global"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-global", "description", "Terraform catalog acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-global", "url", "http://foo.com:8080"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-global", "scope", "global"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-global", "version", "helm_v3"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-global", "annotations.testacc.terraform.io/test", "true"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-global", "labels.testacc.terraform.io/test", "true"),
				),
			},
			{
				Config: testAccRancher2CatalogGlobalUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CatalogExists(testAccRancher2CatalogType+".foo-global", catalog),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-global", "name", "foo-global"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-global", "description", "Terraform catalog acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-global", "url", "http://foo.updated.com:8080"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-global", "scope", "global"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-global", "version", "helm_v3"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-global", "annotations.testacc.terraform.io/test", "false"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-global", "labels.testacc.terraform.io/test", "false"),
				),
			},
			{
				Config: testAccRancher2CatalogGlobal,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CatalogExists(testAccRancher2CatalogType+".foo-global", catalog),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-global", "name", "foo-global"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-global", "description", "Terraform catalog acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-global", "url", "http://foo.com:8080"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-global", "scope", "global"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-global", "version", "helm_v3"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-global", "annotations.testacc.terraform.io/test", "true"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-global", "labels.testacc.terraform.io/test", "true"),
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
			{
				Config: testAccRancher2CatalogGlobal,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CatalogExists(testAccRancher2CatalogType+".foo-global", catalog),
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
			{
				Config: testAccRancher2CatalogClusterConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CatalogExists(testAccRancher2CatalogType+".foo-cluster", catalog),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-cluster", "name", "foo-cluster"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-cluster", "description", "Terraform catalog acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-cluster", "url", "http://foo.com:8080"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-cluster", "scope", "cluster"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-cluster", "cluster_id", testAccRancher2ClusterID),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-cluster", "version", "helm_v2"),
				),
			},
			{
				Config: testAccRancher2CatalogClusterUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CatalogExists(testAccRancher2CatalogType+".foo-cluster", catalog),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-cluster", "name", "foo-cluster"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-cluster", "description", "Terraform catalog acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-cluster", "url", "http://foo.updated.com:8080"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-cluster", "scope", "cluster"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-cluster", "cluster_id", testAccRancher2ClusterID),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-cluster", "version", "helm_v2"),
				),
			},
			{
				Config: testAccRancher2CatalogClusterConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CatalogExists(testAccRancher2CatalogType+".foo-cluster", catalog),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-cluster", "name", "foo-cluster"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-cluster", "description", "Terraform catalog acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-cluster", "url", "http://foo.com:8080"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-cluster", "scope", "cluster"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-cluster", "cluster_id", testAccRancher2ClusterID),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-cluster", "version", "helm_v2"),
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
			{
				Config: testAccRancher2CatalogClusterConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CatalogExists(testAccRancher2CatalogType+".foo-cluster", catalog),
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
			{
				Config: testAccRancher2CatalogProjectConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CatalogExists(testAccRancher2CatalogType+".foo-project", catalog),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-project", "name", "foo-project"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-project", "description", "Terraform catalog acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-project", "url", "http://foo.com:8080"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-project", "scope", "project"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-project", "version", "helm_v2"),
				),
			},
			{
				Config: testAccRancher2CatalogProjectUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CatalogExists(testAccRancher2CatalogType+".foo-project", catalog),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-project", "name", "foo-project"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-project", "description", "Terraform catalog acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-project", "url", "http://foo.updated.com:8080"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-project", "scope", "project"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-project", "version", "helm_v2"),
				),
			},
			{
				Config: testAccRancher2CatalogProjectConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CatalogExists(testAccRancher2CatalogType+".foo-project", catalog),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-project", "name", "foo-project"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-project", "description", "Terraform catalog acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-project", "url", "http://foo.com:8080"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-project", "scope", "project"),
					resource.TestCheckResourceAttr(testAccRancher2CatalogType+".foo-project", "version", "helm_v2"),
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
			{
				Config: testAccRancher2CatalogProjectConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CatalogExists(testAccRancher2CatalogType+".foo-project", catalog),
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
