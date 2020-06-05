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
	testAccRancher2Secret               string
	testAccRancher2SecretUpdate         string
	testAccRancher2SecretConfig         string
	testAccRancher2SecretUpdateConfig   string
	testAccRancher2SecretNs             string
	testAccRancher2SecretNsUpdate       string
	testAccRancher2SecretNsConfig       string
	testAccRancher2SecretNsUpdateConfig string
)

func init() {
	testAccRancher2Secret = `
resource "` + testAccRancher2SecretType + `" "foo" {
  name = "foo"
  description = "Terraform secret acceptance test"
  project_id = rancher2_cluster_sync.testacc.default_project_id
  data = {
    address = "dGVzdC5pbw=="
    password = "cGFzcw=="
    username = "dXNlcg=="
  }
}
`
	testAccRancher2SecretUpdate = `
resource "` + testAccRancher2SecretType + `" "foo" {
  name = "foo"
  description = "Terraform secret acceptance test - updated"
  project_id = rancher2_cluster_sync.testacc.default_project_id
  data = {
    address = "dGVzdC5pbw=="
    password = "cGFzcw=="
    username = "dXNlcjI="
  }
}
`
	testAccRancher2SecretNs = `
resource "` + testAccRancher2SecretType + `" "foo-ns" {
  name = "foo-ns"
  description = "Terraform secret acceptance test"
  project_id = rancher2_cluster_sync.testacc.default_project_id
  namespace_id = rancher2_namespace.testacc.id
  data = {
    address = "dGVzdC5pbw=="
    password = "cGFzcw=="
    username = "dXNlcg=="
  }
}
`
	testAccRancher2SecretNsUpdate = `
resource "` + testAccRancher2SecretType + `" "foo-ns" {
  name = "foo-ns"
  description = "Terraform secret acceptance test - updated"
  project_id = rancher2_cluster_sync.testacc.default_project_id
  namespace_id = rancher2_namespace.testacc.id
  data = {
    address = "dGVzdC5pbw=="
    password = "cGFzcw=="
    username = "dXNlcjI="
  }
}
 `
}

func TestAccRancher2Secret_basic_Project(t *testing.T) {
	var reg interface{}

	testAccRancher2SecretConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2Secret
	testAccRancher2SecretUpdateConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2SecretUpdate

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2SecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2SecretConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2SecretExists(testAccRancher2SecretType+".foo", reg),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo", "description", "Terraform secret acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo", "data.address", "dGVzdC5pbw=="),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo", "data.username", "dXNlcg=="),
				),
			},
			{
				Config: testAccRancher2SecretUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2SecretExists(testAccRancher2SecretType+".foo", reg),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo", "description", "Terraform secret acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo", "data.address", "dGVzdC5pbw=="),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo", "data.username", "dXNlcjI="),
				),
			},
			{
				Config: testAccRancher2SecretConfig,
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
			{
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

	testAccRancher2SecretNsConfig = testAccCheckRancher2ClusterSyncTestacc + testAccCheckRancher2NamespaceTestacc + testAccRancher2SecretNs
	testAccRancher2SecretNsUpdateConfig = testAccCheckRancher2ClusterSyncTestacc + testAccCheckRancher2NamespaceTestacc + testAccRancher2SecretNsUpdate

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2SecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2SecretNsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2SecretExists(testAccRancher2SecretType+".foo-ns", reg),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo-ns", "name", "foo-ns"),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo-ns", "description", "Terraform secret acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo-ns", "data.address", "dGVzdC5pbw=="),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo-ns", "data.username", "dXNlcg=="),
				),
			},
			{
				Config: testAccRancher2SecretNsUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2SecretExists(testAccRancher2SecretType+".foo-ns", reg),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo-ns", "name", "foo-ns"),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo-ns", "description", "Terraform secret acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo-ns", "data.address", "dGVzdC5pbw=="),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo-ns", "data.username", "dXNlcjI="),
				),
			},
			{
				Config: testAccRancher2SecretNsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2SecretExists(testAccRancher2SecretType+".foo-ns", reg),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo-ns", "name", "foo-ns"),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo-ns", "description", "Terraform secret acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo-ns", "data.address", "dGVzdC5pbw=="),
					resource.TestCheckResourceAttr(testAccRancher2SecretType+".foo-ns", "data.username", "dXNlcg=="),
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
			{
				Config: testAccRancher2SecretNsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2SecretExists(testAccRancher2SecretType+".foo-ns", reg),
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
	}
	return nil
}
