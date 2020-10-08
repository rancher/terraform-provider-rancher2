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
	testAccRancher2SettingType   = "rancher2_setting"
	testAccRancher2SettingConfig = `
resource "` + testAccRancher2SettingType + `" "foo" {
	name = "foo"
	value = "Terraform setting acceptance test"
}
`
	testAccRancher2SettingUpdateConfig = `
resource "` + testAccRancher2SettingType + `" "foo" {
	name = "foo"
	value = "Terraform setting acceptance test - updated"
}
 `
)

func TestAccRancher2Setting_basic(t *testing.T) {
	var setting *managementClient.Setting

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2SettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2SettingConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2SettingExists(testAccRancher2SettingType+".foo", setting),
					resource.TestCheckResourceAttr(testAccRancher2SettingType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2SettingType+".foo", "value", "Terraform setting acceptance test"),
				),
			},
			{
				Config: testAccRancher2SettingUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2SettingExists(testAccRancher2SettingType+".foo", setting),
					resource.TestCheckResourceAttr(testAccRancher2SettingType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2SettingType+".foo", "value", "Terraform setting acceptance test - updated"),
				),
			},
			{
				Config: testAccRancher2SettingConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2SettingExists(testAccRancher2SettingType+".foo", setting),
					resource.TestCheckResourceAttr(testAccRancher2SettingType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2SettingType+".foo", "value", "Terraform setting acceptance test"),
				),
			},
		},
	})
}

func TestAccRancher2Setting_disappears(t *testing.T) {
	var setting *managementClient.Setting

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2SettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2SettingConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2SettingExists(testAccRancher2SettingType+".foo", setting),
					testAccRancher2SettingDisappears(setting),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2SettingDisappears(setting *managementClient.Setting) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2SettingType {
				continue
			}
			client, err := testAccProvider.Meta().(*Config).ManagementClient()
			if err != nil {
				return err
			}

			setting, err = client.Setting.ByID(rs.Primary.ID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = client.Setting.Delete(setting)
			if err != nil {
				return fmt.Errorf("Error removing Setting: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Pending:    []string{"active"},
				Target:     []string{"removed"},
				Refresh:    settingStateRefreshFunc(client, setting.ID),
				Timeout:    10 * time.Minute,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf(
					"[ERROR] waiting for setting (%s) to be removed: %s", setting.ID, waitErr)
			}
		}
		return nil

	}
}

func testAccCheckRancher2SettingExists(n string, setting *managementClient.Setting) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No setting ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		foundSetting, err := client.Setting.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("Setting not found")
			}
			return err
		}

		setting = foundSetting

		return nil
	}
}

func testAccCheckRancher2SettingDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2SettingType {
			continue
		}
		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		_, err = client.Setting.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("Setting still exists")
	}
	return nil
}
