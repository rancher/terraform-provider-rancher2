package rancher2

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	managementClient "github.com/rancher/types/client/management/v3"
)

const (
	testAccRancher2NodePoolType = "rancher2_node_pool"
	testAccRancher2Cluster      = `
resource "rancher2_cluster" "foo" {
  name = "foo-custom"
  description = "Terraform node pool cluster acceptance test"
  rke_config {
    network {
      plugin = "canal"
    }
  }
}
`
	testAccRancher2CloudCredential = `
resource "rancher2_cloud_credential" "foo" {
  name = "foo"
  description= "Terraform cloudCredential acceptance test"
  amazonec2_credential_config {
	access_key = "XXXXXXXXXXXXXXXXXXXX"
	secret_key = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  }
}
`
	testAccRancher2NodeTemplate = `
resource "rancher2_node_template" "foo" {
  name = "foo"
  description = "Terraform node pool acceptance test"
  cloud_credential_id = "${rancher2_cloud_credential.foo.id}"
  amazonec2_config {
	ami =  "ami-XXXXXXXXXXXXXXX"
	region = "XX-west-1"
	security_group = ["XXXXXXXX"]
	subnet_id = "subnet-XXXXXXXX"
	vpc_id = "vpc-XXXXXXXX"
	zone = "a"
  }
}
`
)

var (
	testAccRancher2NodePoolConfig         string
	testAccRancher2NodePoolUpdateConfig   string
	testAccRancher2NodePoolRecreateConfig string
)

func init() {
	testAccRancher2NodePoolConfig = testAccRancher2Cluster + testAccRancher2CloudCredential + testAccRancher2NodeTemplate + `
resource "rancher2_node_pool" "foo" {
  cluster_id =  "${rancher2_cluster.foo.id}"
  name = "foo"
  hostname_prefix =  "foo-cluster-0"
  delete_not_ready_after_secs = 120
  node_template_id = "${rancher2_node_template.foo.id}"
  quantity = 1
  control_plane = true
  etcd = true
  worker = true
}
`

	testAccRancher2NodePoolUpdateConfig = testAccRancher2Cluster + testAccRancher2CloudCredential + testAccRancher2NodeTemplate + `
resource "rancher2_node_pool" "foo" {
  cluster_id =  "${rancher2_cluster.foo.id}"
  name = "foo"
  hostname_prefix =  "foo-cluster-0"
  delete_not_ready_after_secs = 60
  node_template_id = "${rancher2_node_template.foo.id}"
  quantity = 3
  control_plane = false
  etcd = true
  worker = false
}
`

	testAccRancher2NodePoolRecreateConfig = testAccRancher2Cluster + testAccRancher2CloudCredential + testAccRancher2NodeTemplate + `
resource "rancher2_node_pool" "foo" {
  cluster_id =  "${rancher2_cluster.foo.id}"
  name = "foo"
  hostname_prefix =  "foo-cluster-0"
  delete_not_ready_after_secs = 120
  node_template_id = "${rancher2_node_template.foo.id}"
  quantity = 1
  control_plane = true
  etcd = true
  worker = true
}
`
}

func TestAccRancher2NodePool_basic(t *testing.T) {
	var nodePool *managementClient.NodePool

	name := testAccRancher2NodePoolType + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NodePoolDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2NodePoolConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodePoolExists(name, nodePool),
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "hostname_prefix", "foo-cluster-0"),
					resource.TestCheckResourceAttr(name, "delete_not_ready_after_secs", "120"),
					resource.TestCheckResourceAttr(name, "control_plane", "true"),
					resource.TestCheckResourceAttr(name, "etcd", "true"),
					resource.TestCheckResourceAttr(name, "worker", "true"),
					resource.TestCheckResourceAttr(name, "quantity", "1"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NodePoolUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodePoolExists(name, nodePool),
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "hostname_prefix", "foo-cluster-0"),
					resource.TestCheckResourceAttr(name, "delete_not_ready_after_secs", "60"),
					resource.TestCheckResourceAttr(name, "control_plane", "false"),
					resource.TestCheckResourceAttr(name, "etcd", "true"),
					resource.TestCheckResourceAttr(name, "worker", "false"),
					resource.TestCheckResourceAttr(name, "quantity", "3"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NodePoolRecreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodePoolExists(name, nodePool),
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "hostname_prefix", "foo-cluster-0"),
					resource.TestCheckResourceAttr(name, "delete_not_ready_after_secs", "120"),
					resource.TestCheckResourceAttr(name, "control_plane", "true"),
					resource.TestCheckResourceAttr(name, "etcd", "true"),
					resource.TestCheckResourceAttr(name, "worker", "true"),
					resource.TestCheckResourceAttr(name, "quantity", "1"),
				),
			},
		},
	})
}

func TestAccRancher2NodePool_disappears(t *testing.T) {
	var nodePool *managementClient.NodePool

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NodePoolDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2NodePoolConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodePoolExists(testAccRancher2NodePoolType+".foo", nodePool),
					testAccRancher2NodePoolDisappears(nodePool),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2NodePoolDisappears(nodePool *managementClient.NodePool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2NodePoolType {
				continue
			}
			client, err := testAccProvider.Meta().(*Config).ManagementClient()
			if err != nil {
				return err
			}

			nodePool, err := client.NodePool.ByID(rs.Primary.ID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = client.NodePool.Delete(nodePool)
			if err != nil {
				return fmt.Errorf("Error removing Node Pool: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Pending:    []string{"removing"},
				Target:     []string{"removed"},
				Refresh:    nodePoolStateRefreshFunc(client, nodePool.ID),
				Timeout:    10 * time.Minute,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf("[ERROR] waiting for node pool (%s) to be removed: %s", nodePool.ID, waitErr)
			}
		}
		return nil
	}
}

func testAccCheckRancher2NodePoolExists(n string, nodePool *managementClient.NodePool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Node Pool ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		foundNodePool, err := client.NodePool.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("Node Pool not found")
			}
			return err
		}

		nodePool = foundNodePool

		return nil
	}
}

func testAccCheckRancher2NodePoolDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2NodePoolType {
			continue
		}
		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		_, err = client.NodePool.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("Node Pool still exists")
	}
	return nil
}
