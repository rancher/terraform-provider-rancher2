package rancher2

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const (
	testAccRancher2RegistryType = "rancher2_registry"
)

var (
	testAccRancher2Registry               string
	testAccRancher2RegistryUpdate         string
	testAccRancher2RegistryConfig         string
	testAccRancher2RegistryUpdateConfig   string
	testAccRancher2RegistryNs             string
	testAccRancher2RegistryNsUpdate       string
	testAccRancher2RegistryNsConfig       string
	testAccRancher2RegistryNsUpdateConfig string
)

func init() {
	testAccRancher2Registry = `
resource "` + testAccRancher2RegistryType + `" "foo" {
  name = "foo"
  description = "Terraform registry acceptance test"
  project_id = rancher2_cluster_sync.testacc.default_project_id
  registries {
    address = "test.io"
    username = "user"
    password = "pass"
  }
}
`
	testAccRancher2RegistryUpdate = `
resource "` + testAccRancher2RegistryType + `" "foo" {
  name = "foo"
  description = "Terraform registry acceptance test - updated"
  project_id = rancher2_cluster_sync.testacc.default_project_id
  registries {
    address = "test.io"
    username = "user2"
    password = "pass"
  }
}
`
	testAccRancher2RegistryNs = `
resource "` + testAccRancher2RegistryType + `" "foo-ns" {
  name = "foo-ns"
  description = "Terraform registry acceptance test"
  project_id = rancher2_cluster_sync.testacc.default_project_id
  namespace_id = rancher2_namespace.testacc.id
  registries {
    address = "test.io"
    username = "user"
    password = "pass"
  }
}
`
	testAccRancher2RegistryNsUpdate = `
resource "` + testAccRancher2RegistryType + `" "foo-ns" {
  name = "foo-ns"
  description = "Terraform registry acceptance test - updated"
  project_id = rancher2_cluster_sync.testacc.default_project_id
  namespace_id = rancher2_namespace.testacc.id
  registries {
    address = "test.io"
    username = "user2"
    password = "pass"
  }
}
 `
}

func TestAccRancher2Registry_basic_Project(t *testing.T) {
	var reg interface{}

	testAccRancher2RegistryConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2Registry
	testAccRancher2RegistryUpdateConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2RegistryUpdate

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2RegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2RegistryConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2RegistryExists(testAccRancher2RegistryType+".foo", reg),
					resource.TestCheckResourceAttr(testAccRancher2RegistryType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2RegistryType+".foo", "description", "Terraform registry acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2RegistryType+".foo", "registries.0.username", "user"),
				),
			},
			{
				Config: testAccRancher2RegistryUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2RegistryExists(testAccRancher2RegistryType+".foo", reg),
					resource.TestCheckResourceAttr(testAccRancher2RegistryType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2RegistryType+".foo", "description", "Terraform registry acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2RegistryType+".foo", "registries.0.username", "user2"),
				),
			},
			{
				Config: testAccRancher2RegistryConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2RegistryExists(testAccRancher2RegistryType+".foo", reg),
					resource.TestCheckResourceAttr(testAccRancher2RegistryType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2RegistryType+".foo", "description", "Terraform registry acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2RegistryType+".foo", "registries.0.username", "user"),
				),
			},
		},
	})
}

func TestAccRancher2Registry_disappears_Project(t *testing.T) {
	var reg interface{}

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2RegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2RegistryConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2RegistryExists(testAccRancher2RegistryType+".foo", reg),
					testAccRancher2RegistryDisappears(reg),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRancher2Registry_basic_Namespaced(t *testing.T) {
	var reg interface{}

	testAccRancher2RegistryNsConfig = testAccCheckRancher2ClusterSyncTestacc + testAccCheckRancher2NamespaceTestacc + testAccRancher2RegistryNs
	testAccRancher2RegistryNsUpdateConfig = testAccCheckRancher2ClusterSyncTestacc + testAccCheckRancher2NamespaceTestacc + testAccRancher2RegistryNsUpdate

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2RegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2RegistryNsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2RegistryExists(testAccRancher2RegistryType+".foo-ns", reg),
					resource.TestCheckResourceAttr(testAccRancher2RegistryType+".foo-ns", "name", "foo-ns"),
					resource.TestCheckResourceAttr(testAccRancher2RegistryType+".foo-ns", "description", "Terraform registry acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2RegistryType+".foo-ns", "registries.0.username", "user"),
				),
			},
			{
				Config: testAccRancher2RegistryNsUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2RegistryExists(testAccRancher2RegistryType+".foo-ns", reg),
					resource.TestCheckResourceAttr(testAccRancher2RegistryType+".foo-ns", "name", "foo-ns"),
					resource.TestCheckResourceAttr(testAccRancher2RegistryType+".foo-ns", "description", "Terraform registry acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2RegistryType+".foo-ns", "registries.0.username", "user2"),
				),
			},
			{
				Config: testAccRancher2RegistryNsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2RegistryExists(testAccRancher2RegistryType+".foo-ns", reg),
					resource.TestCheckResourceAttr(testAccRancher2RegistryType+".foo-ns", "name", "foo-ns"),
					resource.TestCheckResourceAttr(testAccRancher2RegistryType+".foo-ns", "description", "Terraform registry acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2RegistryType+".foo-ns", "registries.0.username", "user"),
				),
			},
		},
	})
}

func TestAccRancher2Registry_disappears_Namespaced(t *testing.T) {
	var reg interface{}

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2RegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2RegistryNsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2RegistryExists(testAccRancher2RegistryType+".foo-ns", reg),
					testAccRancher2RegistryDisappears(reg),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2RegistryDisappears(reg interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2RegistryType {
				continue
			}

			_, projectID := splitProjectID(rs.Primary.Attributes["project_id"])
			namespaceID := rs.Primary.Attributes["namespace_id"]

			reg, err := testAccProvider.Meta().(*Config).GetRegistry(rs.Primary.ID, projectID, namespaceID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = testAccProvider.Meta().(*Config).DeleteRegistry(reg)
			if err != nil {
				return fmt.Errorf("Error removing Registry: %s", err)
			}
		}
		return nil

	}
}

func testAccCheckRancher2RegistryExists(n string, reg interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No registry ID is set")
		}

		_, projectID := splitProjectID(rs.Primary.Attributes["project_id"])
		namespaceID := rs.Primary.Attributes["namespace_id"]

		foundReg, err := testAccProvider.Meta().(*Config).GetRegistry(rs.Primary.ID, projectID, namespaceID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("Registry not found")
			}
			return err
		}

		reg = foundReg

		return nil
	}
}

func testAccCheckRancher2RegistryDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2RegistryType {
			continue
		}

		_, projectID := splitProjectID(rs.Primary.Attributes["project_id"])
		namespaceID := rs.Primary.Attributes["namespace_id"]

		_, err := testAccProvider.Meta().(*Config).GetRegistry(rs.Primary.ID, projectID, namespaceID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
	}
	return nil
}
