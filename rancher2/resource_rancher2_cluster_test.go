package rancher2

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	managementClient "github.com/rancher/types/client/management/v3"
)

const (
	testAccRancher2ClusterType      = "rancher2_cluster"
	testAccRancher2ClusterConfigRKE = `
resource "rancher2_cluster" "foo" {
  name = "foo"
  description = "Terraform custom cluster acceptance test"
  kind = "rke"
  rke_config {
    network {
      plugin = "canal"
    }
    services {
      etcd {
        creation = "6h"
        retention = "24h"
      }
	}
  }
}
`

	testAccRancher2ClusterUpdateConfigRKE = `
resource "rancher2_cluster" "foo" {
  name = "foo"
  description = "Terraform custom cluster acceptance test - updated"
  kind = "rke"
  rke_config {
    network {
      plugin = "canal"
    }
    services {
      etcd {
        creation = "6h"
        retention = "24h"
      }
	}
  }
}
 `

	testAccRancher2ClusterRecreateConfigRKE = `
resource "rancher2_cluster" "foo" {
  name = "foo"
  description = "Terraform custom cluster acceptance test"
  kind = "rke"
  rke_config {
    network {
      plugin = "canal"
    }
    services {
      etcd {
        creation = "6h"
        retention = "24h"
      }
	}
  }
}
 `

	testAccRancher2ClusterConfigImported = `
resource "rancher2_cluster" "foo" {
  name = "foo"
  description = "Terraform imported cluster acceptance test"
  kind = "imported"
}
`

	testAccRancher2ClusterUpdateConfigImported = `
resource "rancher2_cluster" "foo" {
  name = "foo"
  description = "Terraform imported cluster acceptance test - updated"
  kind = "imported"
}
 `

	testAccRancher2ClusterRecreateConfigImported = `
resource "rancher2_cluster" "foo" {
  name = "foo"
  description = "Terraform imported cluster acceptance test"
  kind = "imported"
}
 `
)

func TestAccRancher2Cluster_basic_RKE(t *testing.T) {
	var cluster *managementClient.Cluster

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ClusterDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2ClusterConfigRKE,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterExists(testAccRancher2ClusterType+".foo", cluster),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "description", "Terraform custom cluster acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "kind", "rke"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "rke_config.0.services.0.etcd.0.creation", "6h"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "rke_config.0.services.0.etcd.0.retention", "24h"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2ClusterUpdateConfigRKE,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterExists(testAccRancher2ClusterType+".foo", cluster),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "description", "Terraform custom cluster acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "kind", "rke"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "rke_config.0.services.0.etcd.0.creation", "6h"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "rke_config.0.services.0.etcd.0.retention", "24h"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2ClusterRecreateConfigRKE,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterExists(testAccRancher2ClusterType+".foo", cluster),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "description", "Terraform custom cluster acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "kind", "rke"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "rke_config.0.services.0.etcd.0.creation", "6h"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "rke_config.0.services.0.etcd.0.retention", "24h"),
				),
			},
		},
	})
}

func TestAccRancher2Cluster_disappears_RKE(t *testing.T) {
	var cluster *managementClient.Cluster

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ClusterDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2ClusterConfigRKE,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterExists(testAccRancher2ClusterType+".foo", cluster),
					testAccRancher2ClusterDisappears(cluster),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRancher2Cluster_basic_Imported(t *testing.T) {
	var cluster *managementClient.Cluster

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ClusterDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2ClusterConfigImported,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterExists(testAccRancher2ClusterType+".foo", cluster),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "description", "Terraform imported cluster acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "kind", "imported"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2ClusterUpdateConfigImported,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterExists(testAccRancher2ClusterType+".foo", cluster),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "description", "Terraform imported cluster acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "kind", "imported"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2ClusterRecreateConfigImported,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterExists(testAccRancher2ClusterType+".foo", cluster),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "description", "Terraform imported cluster acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "kind", "imported"),
				),
			},
		},
	})
}

func TestAccRancher2Cluster_disappears_Imported(t *testing.T) {
	var cluster *managementClient.Cluster

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ClusterDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2ClusterConfigImported,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterExists(testAccRancher2ClusterType+".foo", cluster),
					testAccRancher2ClusterDisappears(cluster),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2ClusterDisappears(pro *managementClient.Cluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2ClusterType {
				continue
			}
			client, err := testAccProvider.Meta().(*Config).ManagementClient()
			if err != nil {
				return err
			}

			pro, err = client.Cluster.ByID(rs.Primary.ID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = client.Cluster.Delete(pro)
			if err != nil {
				return fmt.Errorf("Error removing Cluster: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Pending:    []string{"active", "removing"},
				Target:     []string{"removed"},
				Refresh:    clusterRegistrationTokenStateRefreshFunc(client, pro.ID),
				Timeout:    10 * time.Minute,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf(
					"[ERROR] waiting for Cluster (%s) to be removed: %s", pro.ID, waitErr)
			}
		}
		return nil

	}
}

func testAccCheckRancher2ClusterExists(n string, pro *managementClient.Cluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cluster ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		foundPro, err := client.Cluster.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("Cluster not found")
			}
			return err
		}

		pro = foundPro

		return nil
	}
}

func testAccCheckRancher2ClusterDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2ClusterType {
			continue
		}
		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		_, err = client.Cluster.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("Cluster still exists")
	}
	return nil
}
