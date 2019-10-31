package rancher2

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const (
	testAccRancher2SecretType = "rancher2_secret"
)

var (
	testAccRancher2SecretProject          string
	testAccRancher2SecretNamespace        string
	testAccRancher2SecretConfig           string
	testAccRancher2SecretUpdateConfig     string
	testAccRancher2SecretRecreateConfig   string
	testAccRancher2SecretNsConfig         string
	testAccRancher2SecretNsUpdateConfig   string
	testAccRancher2SecretNsRecreateConfig string
)

func init() {
	testAccRancher2SecretProject = `
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

	testAccRancher2SecretNamespace = `
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

	testAccRancher2SecretConfig = testAccRancher2SecretProject + `
resource "rancher2_secret" "foo" {
  name = "foo"
  description = "Terraform secret acceptance test"
  project_id = "${rancher2_project.foo.id}"
  data = {
    address = "dGVzdC5pbw=="
    password = "cGFzcw=="
    username = "dXNlcg=="
  }
}
`

	testAccRancher2SecretUpdateConfig = testAccRancher2SecretProject + `
resource "rancher2_secret" "foo" {
  name = "foo"
  description = "Terraform secret acceptance test - updated"
  project_id = "${rancher2_project.foo.id}"
  data = {
    address = "dGVzdC5pbw=="
    password = "cGFzcw=="
    username = "dXNlcjI="
  }
}
`

	testAccRancher2SecretRecreateConfig = testAccRancher2SecretProject + `
resource "rancher2_secret" "foo" {
  name = "foo"
  description = "Terraform secret acceptance test"
  project_id = "${rancher2_project.foo.id}"
  data = {
    address = "dGVzdC5pbw=="
    password = "cGFzcw=="
    username = "dXNlcg=="
  }
}
`

	testAccRancher2SecretNsConfig = testAccRancher2SecretProject + testAccRancher2SecretNamespace + `
resource "rancher2_secret" "foo" {
  name = "foo"
  description = "Terraform secret acceptance test"
  project_id = "${rancher2_project.foo.id}"
  namespace_id = "${rancher2_namespace.foo.id}"
  data = {
    address = "dGVzdC5pbw=="
    password = "cGFzcw=="
    username = "dXNlcg=="
  }
}
`

	testAccRancher2SecretNsUpdateConfig = testAccRancher2SecretProject + testAccRancher2SecretNamespace + `
resource "rancher2_secret" "foo" {
  name = "foo"
  description = "Terraform secret acceptance test - updated"
  project_id = "${rancher2_project.foo.id}"
  namespace_id = "${rancher2_namespace.foo.id}"
  data = {
    address = "dGVzdC5pbw=="
    password = "cGFzcw=="
    username = "dXNlcjI="
  }
}
 `

	testAccRancher2SecretNsRecreateConfig = testAccRancher2SecretProject + testAccRancher2SecretNamespace + `
resource "rancher2_secret" "foo" {
  name = "foo"
  description = "Terraform secret acceptance test"
  project_id = "${rancher2_project.foo.id}"
  namespace_id = "${rancher2_namespace.foo.id}"
  data = {
    address = "dGVzdC5pbw=="
    password = "cGFzcw=="
    username = "dXNlcg=="
  }
}
 `
}

func TestAccRancher2Secret_basic_Project(t *testing.T) {
	var reg interface{}

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2SecretDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2SecretConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2SecretExists(testAccRancher2SecretType+".foo", reg),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo", "description", "Terraform secret acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo", "data.address", "dGVzdC5pbw=="),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo", "data.username", "dXNlcg=="),
				),
			},
			resource.TestStep{
				Config: testAccRancher2SecretUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2SecretExists(testAccRancher2SecretType+".foo", reg),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo", "description", "Terraform secret acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo", "data.address", "dGVzdC5pbw=="),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo", "data.username", "dXNlcjI="),
				),
			},
			resource.TestStep{
				Config: testAccRancher2SecretRecreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2SecretExists(testAccRancher2SecretType+".foo", reg),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo", "description", "Terraform secret acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo", "data.address", "dGVzdC5pbw=="),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo", "data.username", "dXNlcg=="),
				),
			},
		},
	})
}

func TestAccRancher2Secret_disappears_Project(t *testing.T) {
	var reg interface{}

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2SecretDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2SecretConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2SecretExists(testAccRancher2SecretType+".foo", reg),
					testAccRancher2SecretDisappears(reg),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRancher2Secret_basic_Namespaced(t *testing.T) {
	var reg interface{}

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2SecretDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2SecretNsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2SecretExists(testAccRancher2SecretType+".foo", reg),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo", "description", "Terraform secret acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo", "data.address", "dGVzdC5pbw=="),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo", "data.username", "dXNlcg=="),
				),
			},
			resource.TestStep{
				Config: testAccRancher2SecretNsUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2SecretExists(testAccRancher2SecretType+".foo", reg),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo", "description", "Terraform secret acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo", "data.address", "dGVzdC5pbw=="),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo", "data.username", "dXNlcjI="),
				),
			},
			resource.TestStep{
				Config: testAccRancher2SecretNsRecreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2SecretExists(testAccRancher2SecretType+".foo", reg),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo", "description", "Terraform secret acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo", "data.address", "dGVzdC5pbw=="),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo", "data.username", "dXNlcg=="),
				),
			},
		},
	})
}

func TestAccRancher2Secret_disappears_Namespaced(t *testing.T) {
	var reg interface{}

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2SecretDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2SecretNsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2SecretExists(testAccRancher2SecretType+".foo", reg),
					testAccRancher2SecretDisappears(reg),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2SecretDisappears(reg interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2SecretType {
				continue
			}

			_, projectID := splitProjectID(rs.Primary.Attributes["project_id"])
			namespaceID := rs.Primary.Attributes["namespace_id"]

			reg, err := testAccProvider.Meta().(*Config).GetSecret(rs.Primary.ID, projectID, namespaceID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = testAccProvider.Meta().(*Config).DeleteSecret(reg)
			if err != nil {
				return fmt.Errorf("Error removing Secret: %s", err)
			}
		}
		return nil

	}
}

func testAccCheckRancher2SecretExists(n string, reg interface{}) resource.TestCheckFunc {
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

		foundReg, err := testAccProvider.Meta().(*Config).GetSecret(rs.Primary.ID, projectID, namespaceID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("Secret not found")
			}
			return err
		}

		reg = foundReg

		return nil
	}
}

func testAccCheckRancher2SecretDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2SecretType {
			continue
		}

		_, projectID := splitProjectID(rs.Primary.Attributes["project_id"])
		namespaceID := rs.Primary.Attributes["namespace_id"]

		_, err := testAccProvider.Meta().(*Config).GetSecret(rs.Primary.ID, projectID, namespaceID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("Secret still exists")
	}
	return nil
}
