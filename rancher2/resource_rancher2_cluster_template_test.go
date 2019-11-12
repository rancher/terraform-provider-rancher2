package rancher2

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	managementClient "github.com/rancher/types/client/management/v3"
)

const (
	testAccRancher2ClusterTemplateType   = "rancher2_cluster_template"
	testAccRancher2ClusterTemplateConfig = `
resource "rancher2_cluster_template" "foo" {
  name = "foo"
  members {
    access_type = "owner"
    user_principal_id = "local://user-XXXXX"
  }
  template_revisions {
    name = "V1"
    cluster_config {
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
    default = true
  }
  description = "Terraform cluster template acceptance test"
}
`
	testAccRancher2ClusterTemplateUpdateConfig = `
resource "rancher2_cluster_template" "foo" {
  name = "foo"
  members {
    access_type = "read-only"
    user_principal_id = "local://user-XXXXX"
  }
  template_revisions {
    name = "V2"
    cluster_config {
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
    default = true
  }
  description = "Terraform cluster template acceptance test - updated"
}
 `
	testAccRancher2ClusterTemplateRecreateConfig = `
resource "rancher2_cluster_template" "foo" {
  name = "foo"
  members {
    access_type = "owner"
    user_principal_id = "local://user-XXXXX"
  }
  template_revisions {
    name = "V1"
    cluster_config {
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
    default = true
  }
  description = "Terraform cluster template acceptance test"
}
`
)

func TestAccRancher2ClusterTemplate_basic(t *testing.T) {
	var clusterTemplate *managementClient.ClusterTemplate

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ClusterTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2ClusterTemplateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterTemplateExists(testAccRancher2ClusterTemplateType+".foo", clusterTemplate),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "description", "Terraform cluster template acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "template_revisions.0.name", "V1"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "template_revisions.0.default", "true"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "members.0.access_type", "owner"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2ClusterTemplateUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterTemplateExists(testAccRancher2ClusterTemplateType+".foo", clusterTemplate),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "description", "Terraform cluster template acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "template_revisions.0.name", "V2"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "template_revisions.0.default", "true"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "members.0.access_type", "read-only"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2ClusterTemplateRecreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterTemplateExists(testAccRancher2ClusterTemplateType+".foo", clusterTemplate),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "description", "Terraform cluster template acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "template_revisions.0.name", "V1"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "template_revisions.0.default", "true"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "members.0.access_type", "owner"),
				),
			},
		},
	})
}

func TestAccRancher2ClusterTemplate_disappears(t *testing.T) {
	var clusterTemplate *managementClient.ClusterTemplate

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ClusterTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2ClusterTemplateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterTemplateExists(testAccRancher2ClusterTemplateType+".foo", clusterTemplate),
					testAccRancher2ClusterTemplateDisappears(clusterTemplate),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2ClusterTemplateDisappears(clusterTemplate *managementClient.ClusterTemplate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2ClusterTemplateType {
				continue
			}
			client, err := testAccProvider.Meta().(*Config).ManagementClient()
			if err != nil {
				return err
			}

			clusterTemplate, err := client.ClusterTemplate.ByID(rs.Primary.ID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = client.ClusterTemplate.Delete(clusterTemplate)
			if err != nil {
				return fmt.Errorf("Error removing Cluster Template: %s", err)
			}
		}
		return nil
	}
}

func testAccCheckRancher2ClusterTemplateExists(n string, clusterTemplate *managementClient.ClusterTemplate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cluster Template ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		foundClusterTemplate, err := client.ClusterTemplate.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("Cluster Template not found")
			}
			return err
		}

		clusterTemplate = foundClusterTemplate

		return nil
	}
}

func testAccCheckRancher2ClusterTemplateDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2ClusterTemplateType {
			continue
		}
		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		_, err = client.ClusterTemplate.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("Cluster Template still exists")
	}
	return nil
}
