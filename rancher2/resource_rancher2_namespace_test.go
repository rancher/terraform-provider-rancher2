package rancher2

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	clusterClient "github.com/rancher/types/client/cluster/v3"
)

const (
	testAccCattleNamespaceType   = "rancher2_namespace"
	testAccCattleNamespaceConfig = `
resource "rancher2_namespace" "foo" {
  name = "foo"
  description = "Foo namespace test"
  cluster_id = "local"
  project_name = "Default"
}
`

	testAccCattleNamespaceUpdateConfig = `
resource "rancher2_namespace" "foo" {
  name = "foo"
  description = "Foo namespace test - updated"
  cluster_id = "local"
}
 `

	testAccCattleNamespaceRecreateConfig = `
resource "rancher2_namespace" "foo" {
  name = "foo"
  description = "Foo namespace test"
  cluster_id = "local"
  project_name = "Default"
}
 `
)

func TestAccCattleNamespace_basic(t *testing.T) {
	var ns *clusterClient.Namespace

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCattleNamespaceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCattleNamespaceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCattleNamespaceExists(testAccCattleNamespaceType+".foo", ns),
					resource.TestCheckResourceAttr(testAccCattleNamespaceType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccCattleNamespaceType+".foo", "description", "Foo namespace test"),
					resource.TestCheckResourceAttr(testAccCattleNamespaceType+".foo", "cluster_id", "local"),
					resource.TestCheckResourceAttr(testAccCattleNamespaceType+".foo", "project_name", "Default"),
				),
			},
			resource.TestStep{
				Config: testAccCattleNamespaceUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCattleNamespaceExists(testAccCattleNamespaceType+".foo", ns),
					resource.TestCheckResourceAttr(testAccCattleNamespaceType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccCattleNamespaceType+".foo", "description", "Foo namespace test - updated"),
					resource.TestCheckResourceAttr(testAccCattleNamespaceType+".foo", "cluster_id", "local"),
					resource.TestCheckResourceAttr(testAccCattleNamespaceType+".foo", "project_name", ""),
				),
			},
			resource.TestStep{
				Config: testAccCattleNamespaceRecreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCattleNamespaceExists(testAccCattleNamespaceType+".foo", ns),
					resource.TestCheckResourceAttr(testAccCattleNamespaceType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccCattleNamespaceType+".foo", "description", "Foo namespace test"),
					resource.TestCheckResourceAttr(testAccCattleNamespaceType+".foo", "cluster_id", "local"),
					resource.TestCheckResourceAttr(testAccCattleNamespaceType+".foo", "project_name", "Default"),
				),
			},
		},
	})
}

func TestAccCattleNamespace_disappears(t *testing.T) {
	var ns *clusterClient.Namespace

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCattleNamespaceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCattleNamespaceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCattleNamespaceExists(testAccCattleNamespaceType+".foo", ns),
					testAccCattleNamespaceDisappears(ns),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCattleNamespaceDisappears(ns *clusterClient.Namespace) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccCattleNamespaceType {
				continue
			}
			client, err := testAccProvider.Meta().(*Config).ClusterClient(rs.Primary.Attributes["cluster_id"])
			if err != nil {
				return err
			}

			ns, err = client.Namespace.ByID(rs.Primary.ID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = client.Namespace.Delete(ns)
			if err != nil {
				return fmt.Errorf("Error removing Namespace: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Pending:    []string{"removing"},
				Target:     []string{"removed"},
				Refresh:    NamespaceStateRefreshFunc(client, ns.ID),
				Timeout:    10 * time.Minute,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf(
					"[ERROR] waiting for namespace (%s) to be removed: %s", ns.ID, waitErr)
			}
		}
		return nil

	}
}

func testAccCheckCattleNamespaceExists(n string, ns *clusterClient.Namespace) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No namespace ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ClusterClient(rs.Primary.Attributes["cluster_id"])
		if err != nil {
			return err
		}

		foundNs, err := client.Namespace.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("Namespace not found")
			}
			return err
		}

		ns = foundNs

		return nil
	}
}

func testAccCheckCattleNamespaceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccCattleNamespaceType {
			continue
		}
		client, err := testAccProvider.Meta().(*Config).ClusterClient(rs.Primary.Attributes["cluster_id"])
		if err != nil {
			return err
		}

		_, err = client.Namespace.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("Namespace still exists")
	}
	return nil
}
