package rancher2

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	norman "github.com/rancher/norman/types"
	managementClient "github.com/rancher/types/client/management/v3"
)

const (
	testAccRancher2ClusterType      = "rancher2_cluster"
	testAccRancher2ClusterConfigRKE = `
resource "rancher2_cluster" "foo" {
  name = "foo"
  description = "Terraform custom cluster acceptance test"
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
  rke_config {
    network {
      plugin = "canal"
    }
    services {
      etcd {
        creation = "12h"
        retention = "72h"
      }
	}
  }
}
 `

	testAccRancher2ClusterRecreateConfigRKE = `
resource "rancher2_cluster" "foo" {
  name = "foo"
  description = "Terraform custom cluster acceptance test"
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
}
`

	testAccRancher2ClusterUpdateConfigImported = `
resource "rancher2_cluster" "foo" {
  name = "foo"
  description = "Terraform imported cluster acceptance test - updated"
}
 `

	testAccRancher2ClusterRecreateConfigImported = `
resource "rancher2_cluster" "foo" {
  name = "foo"
  description = "Terraform imported cluster acceptance test"
}
 `
)

func TestAccRancher2Cluster_basic_RKE(t *testing.T) {
	var cluster *Cluster

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
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "rke_config.0.services.0.etcd.0.creation", "12h"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "rke_config.0.services.0.etcd.0.retention", "72h"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2ClusterRecreateConfigRKE,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterExists(testAccRancher2ClusterType+".foo", cluster),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "description", "Terraform custom cluster acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "rke_config.0.services.0.etcd.0.creation", "6h"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "rke_config.0.services.0.etcd.0.retention", "24h"),
				),
			},
		},
	})
}

func TestAccRancher2Cluster_disappears_RKE(t *testing.T) {
	var cluster *Cluster

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
	var cluster *Cluster

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
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "driver", ""),
				),
			},
			resource.TestStep{
				Config: testAccRancher2ClusterUpdateConfigImported,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterExists(testAccRancher2ClusterType+".foo", cluster),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "description", "Terraform imported cluster acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "driver", ""),
				),
			},
			resource.TestStep{
				Config: testAccRancher2ClusterRecreateConfigImported,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterExists(testAccRancher2ClusterType+".foo", cluster),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "description", "Terraform imported cluster acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "driver", ""),
				),
			},
		},
	})
}

func TestAccRancher2Cluster_disappears_Imported(t *testing.T) {
	var cluster *Cluster

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

func testAccRancher2ClusterDisappears(pro *Cluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2ClusterType {
				continue
			}
			client, err := testAccProvider.Meta().(*Config).ManagementClient()
			if err != nil {
				return err
			}

			pro := &norman.Resource{}
			err = client.APIBaseClient.ByID(managementClient.ClusterType, rs.Primary.ID, pro)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = client.APIBaseClient.Delete(pro)
			if err != nil {
				return fmt.Errorf("Error removing Cluster: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Pending:    []string{"active", "removing"},
				Target:     []string{"removed"},
				Refresh:    clusterStateRefreshFunc(client, pro.ID),
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

func testAccCheckRancher2ClusterExists(n string, pro *Cluster) resource.TestCheckFunc {
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

		foundPro := &Cluster{}
		err = client.APIBaseClient.ByID(managementClient.ClusterType, rs.Primary.ID, foundPro)
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

		obj := &Cluster{}
		err = client.APIBaseClient.ByID(managementClient.ClusterType, rs.Primary.ID, obj)
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
