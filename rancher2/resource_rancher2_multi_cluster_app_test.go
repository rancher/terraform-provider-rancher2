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
	testAccRancher2MultiClusterAppType = "rancher2_multi_cluster_app"
)

var (
	testAccRancher2MultiClusterApp             string
	testAccRancher2MultiClusterAppUpdate       string
	testAccRancher2MultiClusterAppConfig       string
	testAccRancher2MultiClusterAppUpdateConfig string
)

func init() {
	testAccRancher2MultiClusterApp = `
resource "` + testAccRancher2MultiClusterAppType + `" "foo" {
  catalog_name = "library"
  name = "foo"
  targets {
    project_id = rancher2_cluster_sync.testacc.default_project_id
  }
  template_name = "docker-registry"
  template_version = "1.8.1"
  answers {
    values = {
      "ingress_host" = "test.xip.io"
    }
  }
  roles = ["project-member"]
}
`
	testAccRancher2MultiClusterAppUpdate = `
resource "` + testAccRancher2MultiClusterAppType + `" "foo" {
  catalog_name = "library"
  name = "foo"
  targets {
    project_id = rancher2_cluster_sync.testacc.default_project_id
  }
  template_name = "docker-registry"
  template_version = "1.8.1"
  answers {
    values = {
      "ingress_host" = "test2.xip.io"
    }
  }
  roles = ["cluster-admin"]
}
`
	testAccRancher2MultiClusterAppConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2MultiClusterApp
	testAccRancher2MultiClusterAppUpdateConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2MultiClusterAppUpdate
}

func TestAccRancher2MultiClusterApp_basic(t *testing.T) {
	var app *managementClient.MultiClusterApp

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2MultiClusterAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2MultiClusterAppConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2MultiClusterAppExists(testAccRancher2MultiClusterAppType+".foo", app),
					resource.TestCheckResourceAttr(testAccRancher2MultiClusterAppType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2MultiClusterAppType+".foo", "template_version_id", "cattle-global-data:library-docker-registry-1.8.1"),
					resource.TestCheckResourceAttr(testAccRancher2MultiClusterAppType+".foo", "answers.0.values.ingress_host", "test.xip.io"),
					resource.TestCheckResourceAttr(testAccRancher2MultiClusterAppType+".foo", "roles.0", "project-member"),
				),
			},
			{
				Config: testAccRancher2MultiClusterAppUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2MultiClusterAppExists(testAccRancher2MultiClusterAppType+".foo", app),
					resource.TestCheckResourceAttr(testAccRancher2MultiClusterAppType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2MultiClusterAppType+".foo", "template_version_id", "cattle-global-data:library-docker-registry-1.8.1"),
					resource.TestCheckResourceAttr(testAccRancher2MultiClusterAppType+".foo", "answers.0.values.ingress_host", "test2.xip.io"),
					resource.TestCheckResourceAttr(testAccRancher2MultiClusterAppType+".foo", "roles.0", "cluster-admin"),
				),
			},
			{
				Config: testAccRancher2MultiClusterAppConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2MultiClusterAppExists(testAccRancher2MultiClusterAppType+".foo", app),
					resource.TestCheckResourceAttr(testAccRancher2MultiClusterAppType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2MultiClusterAppType+".foo", "template_version_id", "cattle-global-data:library-docker-registry-1.8.1"),
					resource.TestCheckResourceAttr(testAccRancher2MultiClusterAppType+".foo", "answers.0.values.ingress_host", "test.xip.io"),
					resource.TestCheckResourceAttr(testAccRancher2MultiClusterAppType+".foo", "roles.0", "project-member"),
				),
			},
		},
	})
}

func TestAccRancher2MultiClusterApp_disappears(t *testing.T) {
	var app *managementClient.MultiClusterApp

	time.Sleep(5 * time.Second)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2MultiClusterAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2MultiClusterAppConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2MultiClusterAppExists(testAccRancher2MultiClusterAppType+".foo", app),
					testAccRancher2MultiClusterAppDisappears(app),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2MultiClusterAppDisappears(mca *managementClient.MultiClusterApp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2MultiClusterAppType {
				continue
			}

			client, err := testAccProvider.Meta().(*Config).ManagementClient()
			if err != nil {
				return err
			}

			mca, err := client.MultiClusterApp.ByID(rs.Primary.ID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = client.MultiClusterApp.Delete(mca)
			if err != nil {
				return fmt.Errorf("Error removing multi cluster app: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Pending:    []string{"removing"},
				Target:     []string{"removed"},
				Refresh:    multiClusterAppStateRefreshFunc(client, rs.Primary.ID),
				Timeout:    10 * time.Minute,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf(
					"[ERROR] waiting for multi cluster app (%s) to be removed: %s", rs.Primary.ID, waitErr)
			}
			time.Sleep(5 * time.Second)
		}
		return nil
	}
}

func testAccCheckRancher2MultiClusterAppExists(n string, mca *managementClient.MultiClusterApp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No multi cluster app ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		foundMultiClusterApp, err := client.MultiClusterApp.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("Multi cluster app not found")
			}
			return err
		}

		mca = foundMultiClusterApp

		return nil
	}
}

func testAccCheckRancher2MultiClusterAppDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2MultiClusterAppType {
			continue
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		_, err = client.MultiClusterApp.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("Multi cluster app still exists")
	}
	return nil
}
