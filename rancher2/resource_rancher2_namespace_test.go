package rancher2

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	clusterClient "github.com/rancher/rancher/pkg/client/generated/cluster/v3"
)

const (
	testAccRancher2NamespaceType = "rancher2_namespace"
)

var (
	testAccRancher2Namespace             string
	testAccRancher2NamespaceUpdate       string
	testAccRancher2NamespaceConfig       string
	testAccRancher2NamespaceUpdateConfig string
)

func init() {
	testAccRancher2Namespace = `
resource "` + testAccRancher2NamespaceType + `" "foo" {
  name = "foo"
  description = "Terraform namespace acceptance test"
  project_id = rancher2_cluster_sync.testacc.default_project_id
  resource_quota {
    limit {
      limits_cpu = "100m"
      limits_memory = "100Mi"
      requests_storage = "1Gi"
    }
  }
  container_resource_limit {
    limits_cpu = "20m"
    limits_memory = "20Mi"
    requests_cpu = "1m"
    requests_memory = "1Mi"
  }
}
`
	testAccRancher2NamespaceUpdate = `
resource "` + testAccRancher2NamespaceType + `" "foo" {
  name = "foo"
  description = "Terraform namespace acceptance test - updated"
  project_id = rancher2_cluster_sync.testacc.default_project_id
  resource_quota {
    limit {
      limits_cpu = "100m"
      limits_memory = "100Mi"
      requests_storage = "2Gi"
    }
  }
  container_resource_limit {
    limits_cpu = "20m"
    limits_memory = "20Mi"
    requests_cpu = "1m"
    requests_memory = "1Mi"
  }
}
`
	testAccRancher2NamespaceConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2Namespace
	testAccRancher2NamespaceUpdateConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2NamespaceUpdate
}

func TestAccRancher2Namespace_basic(t *testing.T) {
	var ns *clusterClient.Namespace

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2NamespaceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NamespaceExists(testAccRancher2NamespaceType+".foo", ns),
					resource.TestCheckResourceAttr(testAccRancher2NamespaceType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2NamespaceType+".foo", "description", "Terraform namespace acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2NamespaceType+".foo", "wait_for_cluster", "false"),
				),
			},
			{
				Config: testAccRancher2NamespaceUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NamespaceExists(testAccRancher2NamespaceType+".foo", ns),
					resource.TestCheckResourceAttr(testAccRancher2NamespaceType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2NamespaceType+".foo", "description", "Terraform namespace acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2NamespaceType+".foo", "wait_for_cluster", "false"),
				),
			},
			{
				Config: testAccRancher2NamespaceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NamespaceExists(testAccRancher2NamespaceType+".foo", ns),
					resource.TestCheckResourceAttr(testAccRancher2NamespaceType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2NamespaceType+".foo", "description", "Terraform namespace acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2NamespaceType+".foo", "wait_for_cluster", "false"),
				),
			},
		},
	})
}

func TestAccRancher2Namespace_disappears(t *testing.T) {
	var ns *clusterClient.Namespace

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NamespaceDestroy,
		Steps: []resource.TestStep{
			{
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
			clusterID, err := clusterIDFromProjectID(rs.Primary.Attributes["project_id"])
			if err != nil {
				return err
			}
			client, err := testAccProvider.Meta().(*Config).ClusterClient(clusterID)
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
				Target:     []string{"removed", "forbidden"},
				Refresh:    namespaceStateRefreshFunc(client, ns.ID),
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

		clusterID, err := clusterIDFromProjectID(rs.Primary.Attributes["project_id"])
		if err != nil {
			return err
		}

		client, err := testAccProvider.Meta().(*Config).ClusterClient(clusterID)
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

		clusterID, err := clusterIDFromProjectID(rs.Primary.Attributes["project_id"])
		if err != nil {
			return err
		}
		client, err := testAccProvider.Meta().(*Config).ClusterClient(clusterID)
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
