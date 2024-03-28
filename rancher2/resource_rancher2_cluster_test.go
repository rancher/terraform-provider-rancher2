package rancher2

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	norman "github.com/rancher/norman/types"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

const (
	testAccRancher2ClusterType      = "rancher2_cluster"
	testAccRancher2ClusterConfigRKE = `
resource "` + testAccRancher2ClusterType + `" "foo" {
  name = "foo"
  description = "Terraform custom cluster acceptance test"
  rke_config {
    network {
      plugin = "canal"
    }
    private_registries {
      is_default = false
      password = "test_pass"
      url= "test.local"
      user = "test_user"
      ecr_credential_plugin {
        aws_access_key_id = "test_key"
        aws_session_token = "test_secret"
      }
    }
    services {
      etcd {
        creation = "6h"
        retention = "24h"
        backup_config {
          enabled = true
          interval_hours = 20
          retention = 10
        }
      }
      kube_api {
        audit_log {
          enabled = true
          configuration {
            max_age = 5
            max_backup = 5
            max_size = 100
            path = "-"
            format = "json"
            policy = "apiVersion: audit.k8s.io/v1\nkind: Policy\nmetadata:\n  creationTimestamp: null\nomitStages:\n- RequestReceived\nrules:\n- level: RequestResponse\n  resources:\n  - resources:\n    - pods\n"
          }
        }
        event_rate_limit {
          configuration = "apiVersion: eventratelimit.admission.k8s.io/v1alpha1\nkind: Configuration\nlimits:\n- type: Server\n  burst: 30000\n  qps: 6000\n"
          enabled = false
        }
      }
    }
    upgrade_strategy {
      drain = true
      max_unavailable_worker = "20%"
    }
  }
  annotations = {
    "testacc.terraform.io/test" = "true"
  }
  labels = {
    "testacc.terraform.io/test" = "true"
  }
}
`
	testAccRancher2ClusterUpdateConfigRKE = `
resource "` + testAccRancher2ClusterType + `" "foo" {
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
      kube_api {
        audit_log {
          enabled = true
          configuration {
            max_age = 7
            max_backup = 5
            max_size = 100
            path = "-"
            format = "json"
            policy = "apiVersion: audit.k8s.io/v1\nkind: Policy\nmetadata:\n  creationTimestamp: null\nomitStages:\n- RequestReceived\nrules:\n- level: RequestResponse\n  resources:\n  - resources:\n    - pods\n"
          }
        }
        event_rate_limit {
          configuration = "apiVersion: eventratelimit.admission.k8s.io/v1alpha1\nkind: Configuration\nlimits:\n- type: Server\n  burst: 30000\n  qps: 6000\n"
          enabled = false
        }
      }
    }
    upgrade_strategy {
      drain = false
      max_unavailable_worker = "10%"
    }
  }
  annotations = {
    "testacc.terraform.io/test" = "false"
  }
  labels = {
    "testacc.terraform.io/test" = "false"
  }
}
 `
	testAccRancher2ClusterConfigImported = `
resource "` + testAccRancher2ClusterType + `" "foo" {
  name = "foo"
  description = "Terraform imported cluster acceptance test"
}
`

	testAccRancher2ClusterUpdateConfigImported = `
resource "` + testAccRancher2ClusterType + `" "foo" {
  name = "foo"
  description = "Terraform imported cluster acceptance test - updated"
}
 `
	testAccRancher2ClusterConfigK3S = `
resource "` + testAccRancher2ClusterType + `" "foo" {
  name = "foo"
  description = "Terraform k3s cluster acceptance test"
  k3s_config {
    upgrade_strategy {
      drain_server_nodes = false
      drain_worker_nodes = false
      server_concurrency = 1
      worker_concurrency = 2
    }
  }
}
`
	testAccRancher2ClusterUpdateConfigK3S = `
resource "` + testAccRancher2ClusterType + `" "foo" {
  name = "foo"
  description = "Terraform k3s cluster acceptance test - updated"
  k3s_config {
    upgrade_strategy {
      drain_server_nodes = false
      drain_worker_nodes = false
      server_concurrency = 1
      worker_concurrency = 2
    }
  }
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
			{
				Config: testAccRancher2ClusterConfigRKE,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterExists(testAccRancher2ClusterType+".foo", cluster),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "description", "Terraform custom cluster acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "rke_config.0.services.0.etcd.0.creation", "6h"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "rke_config.0.services.0.etcd.0.retention", "24h"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "rke_config.0.services.0.kube_api.0.audit_log.0.configuration.0.max_age", "5"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "rke_config.0.upgrade_strategy.0.drain", "true"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "rke_config.0.upgrade_strategy.0.max_unavailable_worker", "20%"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "annotations.testacc.terraform.io/test", "true"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "labels.testacc.terraform.io/test", "true"),
				),
			},
			{
				Config: testAccRancher2ClusterUpdateConfigRKE,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterExists(testAccRancher2ClusterType+".foo", cluster),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "description", "Terraform custom cluster acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "rke_config.0.services.0.etcd.0.creation", "12h"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "rke_config.0.services.0.etcd.0.retention", "72h"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "rke_config.0.services.0.kube_api.0.audit_log.0.configuration.0.max_age", "7"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "rke_config.0.upgrade_strategy.0.drain", "false"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "rke_config.0.upgrade_strategy.0.max_unavailable_worker", "10%"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "annotations.testacc.terraform.io/test", "false"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "labels.testacc.terraform.io/test", "false"),
				),
			},
			{
				Config: testAccRancher2ClusterConfigRKE,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterExists(testAccRancher2ClusterType+".foo", cluster),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "description", "Terraform custom cluster acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "rke_config.0.services.0.etcd.0.creation", "6h"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "rke_config.0.services.0.etcd.0.retention", "24h"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "rke_config.0.services.0.kube_api.0.audit_log.0.configuration.0.max_age", "5"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "rke_config.0.upgrade_strategy.0.drain", "true"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "rke_config.0.upgrade_strategy.0.max_unavailable_worker", "20%"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "annotations.testacc.terraform.io/test", "true"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "labels.testacc.terraform.io/test", "true"),
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
			{
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
			{
				Config: testAccRancher2ClusterConfigImported,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterExists(testAccRancher2ClusterType+".foo", cluster),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "description", "Terraform imported cluster acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "driver", ""),
				),
			},
			{
				Config: testAccRancher2ClusterUpdateConfigImported,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterExists(testAccRancher2ClusterType+".foo", cluster),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "description", "Terraform imported cluster acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "driver", ""),
				),
			},
			{
				Config: testAccRancher2ClusterConfigImported,
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
			{
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

func TestAccRancher2Cluster_basic_K3S(t *testing.T) {
	var cluster *Cluster

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2ClusterConfigK3S,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterExists(testAccRancher2ClusterType+".foo", cluster),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "description", "Terraform k3s cluster acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "driver", ""),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "k3s_config.0.upgrade_strategy.0.drain_server_nodes", "false"),
				),
			},
			{
				Config: testAccRancher2ClusterUpdateConfigK3S,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterExists(testAccRancher2ClusterType+".foo", cluster),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "description", "Terraform k3s cluster acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "driver", ""),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "k3s_config.0.upgrade_strategy.0.drain_server_nodes", "false"),
				),
			},
			{
				Config: testAccRancher2ClusterConfigK3S,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterExists(testAccRancher2ClusterType+".foo", cluster),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "description", "Terraform k3s cluster acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "driver", ""),
					resource.TestCheckResourceAttr(testAccRancher2ClusterType+".foo", "k3s_config.0.upgrade_strategy.0.drain_server_nodes", "false"),
				),
			},
		},
	})
}

func TestAccRancher2Cluster_disappears_K3S(t *testing.T) {
	var cluster *Cluster

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2ClusterConfigK3S,
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
