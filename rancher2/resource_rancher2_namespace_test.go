package rancher2

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	clusterClient "github.com/rancher/types/client/cluster/v3"
)

const (
	testAccRancher2NamespaceType = "rancher2_namespace"
)

var (
	testAccRancher2NamespaceProject        string
	testAccRancher2NamespaceConfig         string
	testAccRancher2NamespaceUpdateConfig   string
	testAccRancher2NamespaceRecreateConfig string
)

func init() {
	testAccRancher2NamespaceProject = `
resource "rancher2_project" "foo" {
  name = "foo"
  cluster_id = "` + testAccRancher2ClusterID + `"
  description = "Terraform namespace acceptance test"
  resource_quota {
    project_limit {
      limits_cpu = "2000m"
      limits_memory = "2000Mi"
      requests_storage = "2Gi"
    }
    namespace_default_limit {
      limits_cpu = "500m"
      limits_memory = "500Mi"
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

	testAccRancher2NamespaceConfig = testAccRancher2NamespaceProject + `
resource "rancher2_namespace" "foo" {
  name = "foo"
  description = "Terraform namespace acceptance test"
  project_id = "${rancher2_project.foo.id}"
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

	testAccRancher2NamespaceUpdateConfig = testAccRancher2NamespaceProject + `
resource "rancher2_namespace" "foo" {
  name = "foo"
  description = "Terraform namespace acceptance test - updated"
  project_id = "${rancher2_project.foo.id}"
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

	testAccRancher2NamespaceRecreateConfig = testAccRancher2NamespaceProject + `
resource "rancher2_namespace" "foo" {
  name = "foo"
  description = "Terraform namespace acceptance test"
  project_id = "${rancher2_project.foo.id}"
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
}

func TestAccRancher2Namespace_basic(t *testing.T) {
	var ns *clusterClient.Namespace

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NamespaceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2NamespaceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NamespaceExists(testAccRancher2NamespaceType+".foo", ns),
					resource.TestCheckResourceAttr(testAccRancher2NamespaceType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2NamespaceType+".foo", "description", "Terraform namespace acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2NamespaceType+".foo", "wait_for_cluster", "false"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NamespaceUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NamespaceExists(testAccRancher2NamespaceType+".foo", ns),
					resource.TestCheckResourceAttr(testAccRancher2NamespaceType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2NamespaceType+".foo", "description", "Terraform namespace acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2NamespaceType+".foo", "wait_for_cluster", "false"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NamespaceRecreateConfig,
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
			resource.TestStep{
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
				Target:     []string{"removed"},
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

		obj, err := client.Namespace.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
		if obj.Removed != "" {
			return nil
		}
		return fmt.Errorf("Namespace still exists")
	}
	return nil
}
