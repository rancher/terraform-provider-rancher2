package rancher2

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const testAccRancher2ClusterV2Type = "rancher2_cluster_v2"

var (
	testAccRancher2ClusterV2             string
	testAccRancher2ClusterV2Update       string
	testAccRancher2ClusterV2Config       string
	testAccRancher2ClusterV2UpdateConfig string
)

func init() {
	testAccRancher2ClusterV2 = `
resource "` + testAccRancher2ClusterV2Type + `" "foo" {
  name = "foo"
  kubernetes_version = "v1.21.4+k3s1"
  enable_network_policy = true
  default_cluster_role_for_project_members = "user"
  rke_config {
    registries {
    	configs {
        hostname = "zmy-domain.test"
        insecure = true
      }
      configs {
        hostname = "my-domain.test"
      }
      mirrors {
        endpoints = ["https://amy-domain.com"]
        hostname = "docker.io"
      }
      mirrors {
        endpoints = ["https://xmy-domain.com"]
        hostname = "bdocker.io"
      }
    }
  }
}
`
	testAccRancher2ClusterV2Update = `
resource "` + testAccRancher2ClusterV2Type + `" "foo" {
  name = "foo"
  kubernetes_version = "v1.21.4+k3s1"
  enable_network_policy = false
  default_cluster_role_for_project_members = "user2"
  rke_config {
    registries {
      mirrors {
        endpoints = ["https://amy-domain.com"]
        hostname = "docker.io"
      }
      mirrors {
        endpoints = ["https://xmy-domain.com"]
        hostname = "bdocker.io"
      }
    }
  }
}
 `
	testAccRancher2ClusterV2Config = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2ClusterV2
	testAccRancher2ClusterV2UpdateConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2ClusterV2Update
}

func TestAccRancher2ClusterV2_basic(t *testing.T) {
	var cluster *ClusterV2

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ClusterV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2ClusterV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterV2Exists(testAccRancher2ClusterV2Type+".foo", cluster),
					resource.TestCheckResourceAttr(testAccRancher2ClusterV2Type+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterV2Type+".foo", "fleet_namespace", "fleet-default"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterV2Type+".foo", "kubernetes_version", "v1.21.4+k3s1"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterV2Type+".foo", "enable_network_policy", "true"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterV2Type+".foo", "default_cluster_role_for_project_members", "user"),
				),
			},
			{
				Config: testAccRancher2ClusterV2UpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterV2Exists(testAccRancher2ClusterV2Type+".foo", cluster),
					resource.TestCheckResourceAttr(testAccRancher2ClusterV2Type+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterV2Type+".foo", "fleet_namespace", "fleet-default"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterV2Type+".foo", "kubernetes_version", "v1.21.4+k3s1"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterV2Type+".foo", "enable_network_policy", "false"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterV2Type+".foo", "default_cluster_role_for_project_members", "user2"),
				),
			},
			{
				Config: testAccRancher2ClusterV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterV2Exists(testAccRancher2ClusterV2Type+".foo", cluster),
					resource.TestCheckResourceAttr(testAccRancher2ClusterV2Type+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterV2Type+".foo", "fleet_namespace", "fleet-default"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterV2Type+".foo", "kubernetes_version", "v1.21.4+k3s1"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterV2Type+".foo", "enable_network_policy", "true"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterV2Type+".foo", "default_cluster_role_for_project_members", "user"),
				),
			},
		},
	})
}

func TestAccRancher2ClusterV2_disappears(t *testing.T) {
	var cluster *ClusterV2

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ClusterV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2ClusterV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterV2Exists(testAccRancher2ClusterV2Type+".foo", cluster),
					testAccRancher2ClusterV2Disappears(cluster),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2ClusterV2Disappears(cat *ClusterV2) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2ClusterV2Type {
				continue
			}
			cluster, err := getClusterV2ByID(testAccProvider.Meta().(*Config), rs.Primary.ID)
			if err != nil {
				if IsNotFound(err) || IsForbidden(err) {
					return nil
				}
				return fmt.Errorf("testAccRancher2ClusterV2Disappears-get: %v", err)
			}
			err = deleteClusterV2(testAccProvider.Meta().(*Config), cluster)
			if err != nil {
				return fmt.Errorf("testAccRancher2ClusterV2Disappears-delete: %v", err)
			}
			stateConf := &resource.StateChangeConf{
				Pending:    []string{},
				Target:     []string{"removed"},
				Refresh:    clusterV2StateRefreshFunc(testAccProvider.Meta(), cluster.ID),
				Timeout:    120 * time.Second,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}
			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf("[ERROR] waiting for cluster (%s) to be deleted: %s", cluster.ID, waitErr)
			}
		}
		return nil

	}
}

func testAccCheckRancher2ClusterV2Exists(n string, cat *ClusterV2) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No cluster ID is set")
		}

		foundReg, err := getClusterV2ByID(testAccProvider.Meta().(*Config), rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return nil
			}
			return fmt.Errorf("testAccCheckRancher2ClusterV2Exists: %v", err)
		}

		cat = foundReg

		return nil
	}
}

func testAccCheckRancher2ClusterV2Destroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2ClusterV2Type {
			continue
		}
		_, err := getClusterV2ByID(testAccProvider.Meta().(*Config), rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return fmt.Errorf("testAccCheckRancher2ClusterV2Destroy: %v", err)
		}
		return fmt.Errorf("ClusterV2 still exists")
	}
	return nil
}
