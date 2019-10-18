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
	testAccRancher2RegistryProject          string
	testAccRancher2RegistryNamespace        string
	testAccRancher2RegistryConfig           string
	testAccRancher2RegistryUpdateConfig     string
	testAccRancher2RegistryRecreateConfig   string
	testAccRancher2RegistryNsConfig         string
	testAccRancher2RegistryNsUpdateConfig   string
	testAccRancher2RegistryNsRecreateConfig string
)

func init() {
	testAccRancher2RegistryProject = `
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
`

	testAccRancher2RegistryNamespace = `
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
`

	testAccRancher2RegistryConfig = testAccRancher2RegistryProject + `
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
`

	testAccRancher2RegistryUpdateConfig = testAccRancher2RegistryProject + `
resource "rancher2_registry" "foo" {
  name = "foo"
  description = "Terraform registry acceptance test - updated"
  project_id = "${rancher2_project.foo.id}"
  registries {
    address = "test.io"
    username = "user2"
    password = "pass"
  }
}
`

	testAccRancher2RegistryRecreateConfig = testAccRancher2RegistryProject + `
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
`

	testAccRancher2RegistryNsConfig = testAccRancher2RegistryProject + testAccRancher2RegistryNamespace + `
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
`

	testAccRancher2RegistryNsUpdateConfig = testAccRancher2RegistryProject + testAccRancher2RegistryNamespace + `
resource "rancher2_registry" "foo" {
  name = "foo"
  description = "Terraform registry acceptance test - updated"
  project_id = "${rancher2_project.foo.id}"
  namespace_id = "${rancher2_namespace.foo.id}"
  registries {
    address = "test.io"
    username = "user2"
    password = "pass"
  }
}
 `

	testAccRancher2RegistryNsRecreateConfig = testAccRancher2RegistryProject + testAccRancher2RegistryNamespace + `
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
 `
}

func TestAccRancher2Registry_basic_Project(t *testing.T) {
	var reg interface{}

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2RegistryDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2RegistryConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2RegistryExists(testAccRancher2RegistryType+".foo", reg),
					resource.TestCheckResourceAttr(testAccRancher2RegistryType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2RegistryType+".foo", "description", "Terraform registry acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2RegistryType+".foo", "registries.0.username", "user"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2RegistryUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2RegistryExists(testAccRancher2RegistryType+".foo", reg),
					resource.TestCheckResourceAttr(testAccRancher2RegistryType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2RegistryType+".foo", "description", "Terraform registry acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2RegistryType+".foo", "registries.0.username", "user2"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2RegistryRecreateConfig,
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
			resource.TestStep{
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

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2RegistryDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2RegistryNsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2RegistryExists(testAccRancher2RegistryType+".foo", reg),
					resource.TestCheckResourceAttr(testAccRancher2RegistryType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2RegistryType+".foo", "description", "Terraform registry acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2RegistryType+".foo", "registries.0.username", "user"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2RegistryNsUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2RegistryExists(testAccRancher2RegistryType+".foo", reg),
					resource.TestCheckResourceAttr(testAccRancher2RegistryType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2RegistryType+".foo", "description", "Terraform registry acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2RegistryType+".foo", "registries.0.username", "user2"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2RegistryNsRecreateConfig,
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

func TestAccRancher2Registry_disappears_Namespaced(t *testing.T) {
	var reg interface{}

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2RegistryDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2RegistryNsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2RegistryExists(testAccRancher2RegistryType+".foo", reg),
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
		return fmt.Errorf("Registry still exists")
	}
	return nil
}
