package rancher2

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
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
    default = true
    enabled = true
    cluster_config {
      cluster_auth_endpoint {
        enabled = true
      }
      enable_cluster_alerting = false
      enable_cluster_monitoring = true
      enable_network_policy = false
      rke_config {
        ignore_docker_version = true
        addon_job_timeout = "30"
        ssh_agent_auth = "false"
        authentication {
          strategy = "x509|webhook"
        }
        private_registries {
          password   = "test"
          url        = "xxxxxxxxx.tes.local"
          user       = "test"
        }
        monitoring {
          provider = "metrics-server"
        }
        ingress {
          provider = "nginx"
          node_selector = {
            app = "ingress"
          }
        }
        network {
          plugin = "canal"
        }
        services {
          etcd {
            creation = "12h"
            retention = "72h"
            snapshot = false
            backup_config {
              enabled = true
              interval_hours = "12"
              retention = "6"
              safe_timestamp = false
            }
          }
          kube_api {
              service_node_port_range = "30000-32767"
              pod_security_policy = false
              always_pull_images = false
          }
        }
        upgrade_strategy {
          drain = true
          max_unavailable_worker = "10%"
          max_unavailable_controlplane = "1"
          drain_input {
	        delete_local_data = false
	        force = false
	        grace_period = "-1"
	        ignore_daemon_sets = true
	        timeout = "120"
          }
        }
      }
    }
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
    name = "V1"
    default = true
    enabled = true
    cluster_config {
      cluster_auth_endpoint {
        enabled = true
      }
      enable_cluster_alerting = false
      enable_cluster_monitoring = true
      enable_network_policy = false
      rke_config {
        ignore_docker_version = true
        addon_job_timeout = "30"
        ssh_agent_auth = "false"
        authentication {
          strategy = "x509|webhook"
        }
        private_registries {
          password   = "test"
          url        = "xxxxxxxxx.tes.local"
          user       = "test"
        }
        monitoring {
          provider = "metrics-server"
        }
        ingress {
          provider = "nginx"
          node_selector = {
            app = "ingress"
          }
        }
        network {
          plugin = "canal"
        }
        services {
          etcd {
            creation = "12h"
            retention = "72h"
            snapshot = false
            backup_config {
              enabled = true
              interval_hours = "12"
              retention = "6"
              safe_timestamp = false
            }
          }
          kube_api {
              service_node_port_range = "30000-32767"
              pod_security_policy = false
              always_pull_images = false
          }
        }
        upgrade_strategy {
          drain = true
          max_unavailable_worker = "10%"
          max_unavailable_controlplane = "1"
          drain_input {
	        delete_local_data = false
	        force = false
	        grace_period = "-1"
	        ignore_daemon_sets = true
	        timeout = "120"
          }
        }
      }
    }
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
    }
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
			{
				Config: testAccRancher2ClusterTemplateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterTemplateExists(testAccRancher2ClusterTemplateType+".foo", clusterTemplate),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "description", "Terraform cluster template acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "template_revisions.0.name", "V1"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "template_revisions.0.default", "true"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "members.0.access_type", "owner"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "template_revisions.0.cluster_config.0.rke_config.0.network.0.plugin", "canal"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "template_revisions.0.cluster_config.0.rke_config.0.services.0.etcd.0.retention", "72h"),
				),
			},
			{
				Config: testAccRancher2ClusterTemplateUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterTemplateExists(testAccRancher2ClusterTemplateType+".foo", clusterTemplate),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "description", "Terraform cluster template acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "template_revisions.1.name", "V2"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "template_revisions.1.default", "false"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "members.0.access_type", "read-only"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "template_revisions.1.cluster_config.0.rke_config.0.network.0.plugin", "canal"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "template_revisions.1.cluster_config.0.rke_config.0.services.0.etcd.0.retention", "48h"),
				),
			},
			{
				Config: testAccRancher2ClusterTemplateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterTemplateExists(testAccRancher2ClusterTemplateType+".foo", clusterTemplate),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "description", "Terraform cluster template acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "template_revisions.0.name", "V1"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "template_revisions.0.default", "true"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "members.0.access_type", "owner"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "template_revisions.0.cluster_config.0.rke_config.0.network.0.plugin", "canal"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterTemplateType+".foo", "template_revisions.0.cluster_config.0.rke_config.0.services.0.etcd.0.retention", "72h"),
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
			{
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
