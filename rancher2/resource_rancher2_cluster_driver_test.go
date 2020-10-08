package rancher2

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

const (
	testAccRancher2ClusterDriverType = "rancher2_cluster_driver"
	testAccRancher2ClusterDriver     = `
resource "` + testAccRancher2ClusterDriverType + `" "foo" {
    active = false
    builtin = false
    checksum = "0x0"
    name = "foo"
    ui_url = "local://ui"
    url = "local://"
	whitelist_domains = ["*.foo.com"]
}
`
	testAccRancher2ClusterDriverUpdate = `
resource "` + testAccRancher2ClusterDriverType + `" "foo" {
    active = false
    builtin = false
    checksum = "0x1"
    name = "foo"
    ui_url = "local://ui/updated"
    url = "local://updated"
    whitelist_domains = ["*.foo.com", "updated.foo.com"]
}
 `
)

func TestAccRancher2ClusterDriver_basic(t *testing.T) {
	var clusterDriver *managementClient.KontainerDriver
	name := testAccRancher2ClusterDriverType + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ClusterDriverDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2ClusterDriver,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterDriverExists(name, clusterDriver),
					resource.TestCheckResourceAttr(name, "active", "false"),
					resource.TestCheckResourceAttr(name, "builtin", "false"),
					resource.TestCheckResourceAttr(name, "checksum", "0x0"),
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "ui_url", "local://ui"),
					resource.TestCheckResourceAttr(name, "url", "local://"),
					resource.TestCheckResourceAttr(name, "whitelist_domains.0", "*.foo.com"),
				),
			},
			{
				Config: testAccRancher2ClusterDriverUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterDriverExists(name, clusterDriver),
					resource.TestCheckResourceAttr(name, "active", "false"),
					resource.TestCheckResourceAttr(name, "builtin", "false"),
					resource.TestCheckResourceAttr(name, "checksum", "0x1"),
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "ui_url", "local://ui/updated"),
					resource.TestCheckResourceAttr(name, "url", "local://updated"),
					resource.TestCheckResourceAttr(name, "whitelist_domains.0", "*.foo.com"),
					resource.TestCheckResourceAttr(name, "whitelist_domains.1", "updated.foo.com"),
				),
			},
			{
				Config: testAccRancher2ClusterDriver,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterDriverExists(name, clusterDriver),
					resource.TestCheckResourceAttr(name, "active", "false"),
					resource.TestCheckResourceAttr(name, "builtin", "false"),
					resource.TestCheckResourceAttr(name, "checksum", "0x0"),
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "ui_url", "local://ui"),
					resource.TestCheckResourceAttr(name, "url", "local://"),
					resource.TestCheckResourceAttr(name, "whitelist_domains.0", "*.foo.com"),
				),
			},
		},
	})
}

func TestAccRancher2ClusterDriver_disappears(t *testing.T) {
	var clusterDriver *managementClient.KontainerDriver

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ClusterDriverDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2ClusterDriver,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterDriverExists(testAccRancher2ClusterDriverType+".foo", clusterDriver),
					testAccRancher2ClusterDriverDisappears(clusterDriver),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2ClusterDriverDisappears(clusterDriver *managementClient.KontainerDriver) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2ClusterDriverType {
				continue
			}
			client, err := testAccProvider.Meta().(*Config).ManagementClient()
			if err != nil {
				return err
			}

			clusterDriver, err := client.KontainerDriver.ByID(rs.Primary.ID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = client.KontainerDriver.Delete(clusterDriver)
			if err != nil {
				return fmt.Errorf("Error removing Cluster Driver: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Pending:    []string{"removing"},
				Target:     []string{"removed"},
				Refresh:    clusterDriverStateRefreshFunc(client, clusterDriver.ID),
				Timeout:    10 * time.Minute,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf(
					"[ERROR] waiting for cluster driver (%s) to be removed: %s", clusterDriver.ID, waitErr)
			}
		}
		return nil
	}
}

func testAccCheckRancher2ClusterDriverExists(n string, clusterDriver *managementClient.KontainerDriver) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cluster Driver ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		foundClusterDriver, err := client.KontainerDriver.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("Cluster Driver not found")
			}
			return err
		}

		clusterDriver = foundClusterDriver

		return nil
	}
}

func testAccCheckRancher2ClusterDriverDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2ClusterDriverType {
			continue
		}
		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		_, err = client.KontainerDriver.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("Cluster Driver still exists")
	}
	return nil
}
