package rancher2

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	projectClient "github.com/rancher/rancher/pkg/client/generated/project/v3"
)

const (
	testAccRancher2AppType = "rancher2_app"
)

var (
	testAccRancher2App             string
	testAccRancher2AppUpdate       string
	testAccRancher2AppConfig       string
	testAccRancher2AppUpdateConfig string
)

func init() {
	testAccRancher2App = `
resource "rancher2_app" "foo" {
  catalog_name = "library"
  name = "foo"
  description = "Terraform app acceptance test"
  project_id = rancher2_cluster_sync.testacc.default_project_id
  template_name = "docker-registry"
  template_version = "1.8.1"
  target_namespace = rancher2_namespace.testacc.name
  answers = {
    "ingress_host" = "test.xip.io"
  }
  annotations = {
    "testacc.terraform.io/test" = "true"
  }
  labels = {
    "testacc.terraform.io/test" = "true"
  }
}
`
	testAccRancher2AppUpdate = `
resource "rancher2_app" "foo" {
  catalog_name = "library"
  name = "foo"
  description = "Terraform app acceptance test - updated"
  project_id = rancher2_cluster_sync.testacc.default_project_id
  template_name = "docker-registry"
  template_version = "1.8.1"
  target_namespace = rancher2_namespace.testacc.name
  answers = {
    "ingress_host" = "test2.xip.io"
  }
  annotations = {
    "testacc.terraform.io/test" = "false"
  }
  labels = {
    "testacc.terraform.io/test" = "false"
  }
}
`
	testAccRancher2AppConfig = testAccCheckRancher2ClusterSyncTestacc + testAccCheckRancher2NamespaceTestacc + testAccRancher2App
	testAccRancher2AppUpdateConfig = testAccCheckRancher2ClusterSyncTestacc + testAccCheckRancher2NamespaceTestacc + testAccRancher2AppUpdate
}

func TestAccRancher2App_basic(t *testing.T) {
	var app *projectClient.App

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2AppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2AppConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AppExists(testAccRancher2AppType+".foo", app),
					resource.TestCheckResourceAttr(testAccRancher2AppType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2AppType+".foo", "description", "Terraform app acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2AppType+".foo", "target_namespace", "testacc"),
					resource.TestCheckResourceAttr(testAccRancher2AppType+".foo", "external_id", "catalog://?catalog=library&template=docker-registry&version=1.8.1"),
					resource.TestCheckResourceAttr(testAccRancher2AppType+".foo", "answers.ingress_host", "test.xip.io"),
					resource.TestCheckResourceAttr(testAccRancher2AppType+".foo", "annotations.testacc.terraform.io/test", "true"),
					resource.TestCheckResourceAttr(testAccRancher2AppType+".foo", "labels.testacc.terraform.io/test", "true"),
				),
			},
			{
				Config: testAccRancher2AppUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AppExists(testAccRancher2AppType+".foo", app),
					resource.TestCheckResourceAttr(testAccRancher2AppType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2AppType+".foo", "description", "Terraform app acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2AppType+".foo", "target_namespace", "testacc"),
					resource.TestCheckResourceAttr(testAccRancher2AppType+".foo", "external_id", "catalog://?catalog=library&template=docker-registry&version=1.8.1"),
					resource.TestCheckResourceAttr(testAccRancher2AppType+".foo", "answers.ingress_host", "test2.xip.io"),
					resource.TestCheckResourceAttr(testAccRancher2AppType+".foo", "annotations.testacc.terraform.io/test", "false"),
					resource.TestCheckResourceAttr(testAccRancher2AppType+".foo", "labels.testacc.terraform.io/test", "false"),
				),
			},
			{
				Config: testAccRancher2AppConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AppExists(testAccRancher2AppType+".foo", app),
					resource.TestCheckResourceAttr(testAccRancher2AppType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2AppType+".foo", "description", "Terraform app acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2AppType+".foo", "target_namespace", "testacc"),
					resource.TestCheckResourceAttr(testAccRancher2AppType+".foo", "external_id", "catalog://?catalog=library&template=docker-registry&version=1.8.1"),
					resource.TestCheckResourceAttr(testAccRancher2AppType+".foo", "answers.ingress_host", "test.xip.io"),
					resource.TestCheckResourceAttr(testAccRancher2AppType+".foo", "annotations.testacc.terraform.io/test", "true"),
					resource.TestCheckResourceAttr(testAccRancher2AppType+".foo", "labels.testacc.terraform.io/test", "true"),
				),
			},
		},
	})
}

func TestAccRancher2App_disappears(t *testing.T) {
	var app *projectClient.App

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2AppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2AppConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AppExists(testAccRancher2AppType+".foo", app),
					testAccRancher2AppDisappears(app),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2AppDisappears(app *projectClient.App) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2AppType {
				continue
			}

			client, err := testAccProvider.Meta().(*Config).ProjectClient(rs.Primary.Attributes["project_id"])
			if err != nil {
				return err
			}

			app, err := client.App.ByID(rs.Primary.ID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = client.App.Delete(app)
			if err != nil {
				return fmt.Errorf("Error removing App: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Pending:    []string{"removing"},
				Target:     []string{"removed"},
				Refresh:    appStateRefreshFunc(client, rs.Primary.ID),
				Timeout:    10 * time.Minute,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf(
					"[ERROR] waiting for App (%s) to be removed: %s", app.ID, waitErr)
			}
		}
		return nil
	}
}

func testAccCheckRancher2AppExists(n string, app *projectClient.App) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No App ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ProjectClient(rs.Primary.Attributes["project_id"])
		if err != nil {
			return err
		}

		foundApp, err := client.App.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("App not found")
			}
			return err
		}

		app = foundApp

		return nil
	}
}

func testAccCheckRancher2AppDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2AppType {
			continue
		}

		client, err := testAccProvider.Meta().(*Config).ProjectClient(rs.Primary.Attributes["project_id"])
		if err != nil {
			return err
		}

		_, err = client.App.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("App still exists")
	}
	return nil
}
