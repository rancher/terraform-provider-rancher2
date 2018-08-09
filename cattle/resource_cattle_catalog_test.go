package cattle

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	managementClient "github.com/rancher/types/client/management/v3"
)

const (
	testAccCattleCatalogType   = "cattle_catalog"
	testAccCattleCatalogConfig = `
resource "cattle_catalog" "foo" {
  name = "foo"
  url = "http://foo.com:8080"
  description= "Foo catalog test"
}
`

	testAccCattleCatalogUpdateConfig = `
resource "cattle_catalog" "foo" {
  name = "foo"
  url = "http://foo.updated.com:8080"
  description= "Foo catalog test - updated"
}
 `

	testAccCattleCatalogRecreateConfig = `
resource "cattle_catalog" "foo" {
  name = "foo"
  url = "http://foo.com:8080"
  description= "Foo catalog test"
}
 `
)

func TestAccCattleCatalog_basic(t *testing.T) {
	var catalog *managementClient.Catalog

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCattleCatalogDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCattleCatalogConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCattleCatalogExists(testAccCattleCatalogType+".foo", catalog),
					resource.TestCheckResourceAttr(testAccCattleCatalogType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccCattleCatalogType+".foo", "description", "Foo catalog test"),
					resource.TestCheckResourceAttr(testAccCattleCatalogType+".foo", "url", "http://foo.com:8080"),
				),
			},
			resource.TestStep{
				Config: testAccCattleCatalogUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCattleCatalogExists(testAccCattleCatalogType+".foo", catalog),
					resource.TestCheckResourceAttr(testAccCattleCatalogType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccCattleCatalogType+".foo", "description", "Foo catalog test - updated"),
					resource.TestCheckResourceAttr(testAccCattleCatalogType+".foo", "url", "http://foo.updated.com:8080"),
				),
			},
			resource.TestStep{
				Config: testAccCattleCatalogRecreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCattleCatalogExists(testAccCattleCatalogType+".foo", catalog),
					resource.TestCheckResourceAttr(testAccCattleCatalogType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccCattleCatalogType+".foo", "description", "Foo catalog test"),
					resource.TestCheckResourceAttr(testAccCattleCatalogType+".foo", "url", "http://foo.com:8080"),
				),
			},
		},
	})
}

func TestAccCattleCatalog_disappears(t *testing.T) {
	var catalog *managementClient.Catalog

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCattleCatalogDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCattleCatalogConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCattleCatalogExists(testAccCattleCatalogType+".foo", catalog),
					testAccCattleCatalogDisappears(catalog),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCattleCatalogDisappears(cat *managementClient.Catalog) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccCattleCatalogType {
				continue
			}
			client, err := testAccProvider.Meta().(*Config).ManagementClient()
			if err != nil {
				return err
			}

			cat, err = client.Catalog.ByID(rs.Primary.ID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = client.Catalog.Delete(cat)
			if err != nil {
				return fmt.Errorf("Error removing Catalog: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Pending:    []string{"active"},
				Target:     []string{"removed"},
				Refresh:    CatalogStateRefreshFunc(client, cat.ID),
				Timeout:    10 * time.Minute,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf(
					"[ERROR] waiting for catalog (%s) to be removed: %s", cat.ID, waitErr)
			}
		}
		return nil

	}
}

func testAccCheckCattleCatalogExists(n string, cat *managementClient.Catalog) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No catalog ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		foundReg, err := client.Catalog.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("Catalog not found")
			}
			return err
		}

		cat = foundReg

		return nil
	}
}

func testAccCheckCattleCatalogDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccCattleCatalogType {
			continue
		}
		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		_, err = client.Catalog.ByID(rs.Primary.ID)
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
