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
	testAccRancher2GlobalRoleBindingType = "rancher2_global_role_binding"
)

var (
	testAccRancher2GlobalRoleBinding             string
	testAccRancher2GlobalRoleBindingUpdate       string
	testAccRancher2GlobalRoleBindingConfig       string
	testAccRancher2GlobalRoleBindingUpdateConfig string
)

func init() {
	testAccRancher2GlobalRoleBinding = `
resource "` + testAccRancher2GlobalRoleBindingType + `" "foo" {
  name = "foo-test"
  global_role_id = "user-base"
  user_id = rancher2_user.foo.id
}
`
	testAccRancher2GlobalRoleBindingUpdate = `
resource "` + testAccRancher2GlobalRoleBindingType + `" "foo" {
  name = "foo-test-updated"
  global_role_id = "user-base"
  user_id = rancher2_user.foo.id
}
`
}

func TestAccRancher2GlobalRoleBinding_basic(t *testing.T) {
	var globalRole *managementClient.GlobalRoleBinding

	testAccRancher2GlobalRoleBindingConfig = testAccRancher2User + testAccRancher2GlobalRoleBinding
	testAccRancher2GlobalRoleBindingUpdateConfig = testAccRancher2User + testAccRancher2GlobalRoleBindingUpdate
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2GlobalRoleBindingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2GlobalRoleBindingConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2GlobalRoleBindingExists(testAccRancher2GlobalRoleBindingType+".foo", globalRole),
					resource.TestCheckResourceAttr(testAccRancher2GlobalRoleBindingType+".foo", "name", "foo-test"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalRoleBindingType+".foo", "global_role_id", "user-base"),
				),
			},
			{
				Config: testAccRancher2GlobalRoleBindingUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2GlobalRoleBindingExists(testAccRancher2GlobalRoleBindingType+".foo", globalRole),
					resource.TestCheckResourceAttr(testAccRancher2GlobalRoleBindingType+".foo", "name", "foo-test-updated"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalRoleBindingType+".foo", "global_role_id", "user-base"),
				),
			},
			{
				Config: testAccRancher2GlobalRoleBindingConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2GlobalRoleBindingExists(testAccRancher2GlobalRoleBindingType+".foo", globalRole),
					resource.TestCheckResourceAttr(testAccRancher2GlobalRoleBindingType+".foo", "name", "foo-test"),
					resource.TestCheckResourceAttr(testAccRancher2GlobalRoleBindingType+".foo", "global_role_id", "user-base"),
				),
			},
		},
	})
}

func TestAccRancher2GlobalRoleBinding_disappears(t *testing.T) {
	var globalRole *managementClient.GlobalRoleBinding

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2GlobalRoleBindingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2GlobalRoleBindingConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2GlobalRoleBindingExists(testAccRancher2GlobalRoleBindingType+".foo", globalRole),
					testAccRancher2GlobalRoleBindingDisappears(globalRole),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2GlobalRoleBindingDisappears(pro *managementClient.GlobalRoleBinding) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2GlobalRoleBindingType {
				continue
			}
			client, err := testAccProvider.Meta().(*Config).ManagementClient()
			if err != nil {
				return err
			}

			pro, err = client.GlobalRoleBinding.ByID(rs.Primary.ID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = client.GlobalRoleBinding.Delete(pro)
			if err != nil {
				return fmt.Errorf("Error removing Global Role Binding: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Pending:    []string{"active"},
				Target:     []string{"removed"},
				Refresh:    globalRoleBindingStateRefreshFunc(client, pro.ID),
				Timeout:    10 * time.Minute,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf(
					"[ERROR] waiting for Global Role Binding (%s) to be removed: %s", pro.ID, waitErr)
			}
		}
		return nil

	}
}

func testAccCheckRancher2GlobalRoleBindingExists(n string, pro *managementClient.GlobalRoleBinding) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Global Role Binding ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		foundPro, err := client.GlobalRoleBinding.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("Global Role Binding not found")
			}
			return err
		}

		pro = foundPro

		return nil
	}
}

func testAccCheckRancher2GlobalRoleBindingDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2GlobalRoleBindingType {
			continue
		}
		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		_, err = client.GlobalRoleBinding.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("Global Role Binding still exists")
	}
	return nil
}
