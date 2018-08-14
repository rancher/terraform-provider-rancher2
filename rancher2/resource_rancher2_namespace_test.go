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
	testAccRancher2NamespaceType   = "rancher2_namespace"
	testAccRancher2NamespaceConfig = `
resource "rancher2_namespace" "foo" {
  name = "foo"
  description = "Foo namespace test"
  cluster_id = "local"
  project_name = "Default"
}
`

	testAccRancher2NamespaceUpdateConfig = `
resource "rancher2_namespace" "foo" {
  name = "foo"
  description = "Foo namespace test - updated"
  cluster_id = "local"
}
 `

	testAccRancher2NamespaceRecreateConfig = `
resource "rancher2_namespace" "foo" {
  name = "foo"
  description = "Foo namespace test"
  cluster_id = "local"
  project_name = "Default"
}
 `
)

func TestAccRancher2Namespace_basic(t *testing.T) {
	var ns *clusterClient.Namespace

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NamespaceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2NamespaceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NamespaceExists(testAccRancher2NamespaceType+".foo", ns),
					resource.TestCheckResourceAttr(testAccRancher2NamespaceType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2NamespaceType+".foo", "description", "Foo namespace test"),
					resource.TestCheckResourceAttr(testAccRancher2NamespaceType+".foo", "cluster_id", "local"),
					resource.TestCheckResourceAttr(testAccRancher2NamespaceType+".foo", "project_name", "Default"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NamespaceUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NamespaceExists(testAccRancher2NamespaceType+".foo", ns),
					resource.TestCheckResourceAttr(testAccRancher2NamespaceType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2NamespaceType+".foo", "description", "Foo namespace test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2NamespaceType+".foo", "cluster_id", "local"),
					resource.TestCheckResourceAttr(testAccRancher2NamespaceType+".foo", "project_name", ""),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NamespaceRecreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NamespaceExists(testAccRancher2NamespaceType+".foo", ns),
					resource.TestCheckResourceAttr(testAccRancher2NamespaceType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2NamespaceType+".foo", "description", "Foo namespace test"),
					resource.TestCheckResourceAttr(testAccRancher2NamespaceType+".foo", "cluster_id", "local"),
					resource.TestCheckResourceAttr(testAccRancher2NamespaceType+".foo", "project_name", "Default"),
				),
			},
		},
	})
}

func TestAccRancher2Namespace_disappears(t *testing.T) {
	var ns *clusterClient.Namespace

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NamespaceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2NamespaceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NamespaceExists(testAccRancher2NamespaceType+".foo", ns),
					testAccRancher2NamespaceDisappears(ns),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2NamespaceDisappears(ns *clusterClient.Namespace) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2NamespaceType {
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

func testAccCheckRancher2NamespaceExists(n string, ns *clusterClient.Namespace) resource.TestCheckFunc {
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

func testAccCheckRancher2NamespaceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2NamespaceType {
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
