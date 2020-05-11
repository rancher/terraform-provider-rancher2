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
resource "` + testAccRancher2ClusterTemplateType + `" "foo" {
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
      scheduled_cluster_scan {
	    enabled = true
	    scan_config {
	      cis_scan_config {
	        debug_master = true
	        debug_worker = true
	      }
	    }
	    schedule_config {
	      cron_schedule = "30 * * * *"
	      retention = 5
	    }
	  }
    }
    default = true
  }
  description = "Terraform cluster template acceptance test"
}
`
	testAccRancher2ClusterTemplateUpdateConfig = `
resource "` + testAccRancher2ClusterTemplateType + `" "foo" {
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
            retention = "48h"
          }
        }
      }
      scheduled_cluster_scan {
	    enabled = true
	    scan_config {
	      cis_scan_config {
	        debug_master = true
	        debug_worker = true
	      }
	    }
	    schedule_config {
	      cron_schedule = "30 10 * * *"
	      retention = 5
	    }
	  }
    }
    default = true
  }
  description = "Terraform cluster template acceptance test - updated"
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
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "template_revisions.0.cluster_config.0.rke_config.0.network.0.plugin", "canal"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "template_revisions.0.cluster_config.0.rke_config.0.services.0.etcd.0.retention", "24h"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "template_revisions.0.cluster_config.0.scheduled_cluster_scan.0.scan_config.0.cis_scan_config.0.debug_worker", "true"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "template_revisions.0.cluster_config.0.scheduled_cluster_scan.0.schedule_config.0.cron_schedule", "30 * * * *"),
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
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "template_revisions.0.cluster_config.0.rke_config.0.network.0.plugin", "canal"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "template_revisions.0.cluster_config.0.rke_config.0.services.0.etcd.0.retention", "48h"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "template_revisions.0.cluster_config.0.scheduled_cluster_scan.0.scan_config.0.cis_scan_config.0.debug_worker", "true"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "template_revisions.0.cluster_config.0.scheduled_cluster_scan.0.schedule_config.0.cron_schedule", "30 10 * * *"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2ClusterTemplateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterTemplateExists(testAccRancher2ClusterTemplateType+".foo", clusterTemplate),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "description", "Terraform cluster template acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "template_revisions.0.name", "V1"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "template_revisions.0.default", "true"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "members.0.access_type", "owner"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "template_revisions.0.cluster_config.0.rke_config.0.network.0.plugin", "canal"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "template_revisions.0.cluster_config.0.rke_config.0.services.0.etcd.0.retention", "24h"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "template_revisions.0.cluster_config.0.scheduled_cluster_scan.0.scan_config.0.cis_scan_config.0.debug_worker", "true"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "template_revisions.0.cluster_config.0.scheduled_cluster_scan.0.schedule_config.0.cron_schedule", "30 * * * *"),
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
