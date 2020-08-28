package rancher2

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

const (
	testAccRancher2NodePoolType = "rancher2_node_pool"
)

var (
	testAccRancher2NodePool             string
	testAccRancher2NodePoolUpdate       string
	testAccRancher2NodePoolConfig       string
	testAccRancher2NodePoolUpdateConfig string
)

func init() {
	testAccRancher2NodePool = `
resource "` + testAccRancher2NodePoolType + `" "foo" {
  cluster_id =  rancher2_cluster.foo.id
  name = "foo"
  hostname_prefix =  "foo-cluster-0"
  delete_not_ready_after_secs = 120
  node_template_id = rancher2_node_template.foo-aws.id
  quantity = 1
  control_plane = true
  etcd = true
  worker = true
}
`
	testAccRancher2NodePoolUpdate = `
resource "` + testAccRancher2NodePoolType + `" "foo" {
  cluster_id =  rancher2_cluster.foo.id
  name = "foo"
  hostname_prefix =  "foo-cluster-0"
  delete_not_ready_after_secs = 60
  node_template_id = rancher2_node_template.foo-aws.id
  quantity = 3
  control_plane = false
  etcd = true
  worker = false
}
`
}

func TestAccRancher2NodePool_basic(t *testing.T) {
	var nodePool *managementClient.NodePool

	testAccRancher2NodePoolConfig = testAccRancher2ClusterConfigRKE + testAccRancher2CloudCredentialConfigAmazonec2 + testAccRancher2NodeTemplateAmazonec2 + testAccRancher2NodePool
	testAccRancher2NodePoolUpdateConfig = testAccRancher2ClusterConfigRKE + testAccRancher2CloudCredentialConfigAmazonec2 + testAccRancher2NodeTemplateAmazonec2 + testAccRancher2NodePoolUpdate

	name := testAccRancher2NodePoolType + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NodePoolDestroy,
		Steps: []resource.TestStep{
			{
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
			{
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
			{
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
			{
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
