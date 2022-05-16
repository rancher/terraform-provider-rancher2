package rancher2

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/rancher/norman/types"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

const (
	testAccRancher2UserType = "rancher2_user"
)

var (
	testAccRancher2User, testAccRancher2UserUpdate                   string
	testAccRancher2UserWithToken, testAccRancher2UserWithTokenUpdate string

	token *managementClient.Token
)

func init() {
	testAccRancher2User = `
resource "` + testAccRancher2UserType + `" "foo" {
  name = "Terraform user acceptance test"
  username = "foo"
  password = "TestACC123456"
  enabled = true
}
`
	testAccRancher2UserWithToken = `
resource "` + testAccRancher2UserType + `" "foo" {
  name = "Terraform user acceptance test"
  username = "foo"
  password = "TestACC123456"
  enabled = true
  token_config {
	description = "foo"
	ttl = 120000
  }
}
`
	testAccRancher2UserUpdate = `
resource "` + testAccRancher2UserType + `" "foo" {
  name = "Terraform user acceptance test - Updated"
  username = "foo"
  password = "TestACC1234567"
  enabled = false
}
`
	testAccRancher2UserWithTokenUpdate = `
 resource "` + testAccRancher2UserType + `" "foo" {
   name = "Terraform user acceptance test"
   username = "foo"
   password = "TestACC123456"
   enabled = true
   token_config {
	 description = "foo"
	 ttl = 240000
   }
 }
`
}

func TestAccRancher2User_basic(t *testing.T) {
	var user *managementClient.User

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2UserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2User,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2UserExists(testAccRancher2UserType+".foo", user),
					resource.TestCheckResourceAttr(testAccRancher2UserType+".foo", "name", "Terraform user acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2UserType+".foo", "password", "TestACC123456"),
					resource.TestCheckResourceAttr(testAccRancher2UserType+".foo", "enabled", "true"),
					resource.TestCheckNoResourceAttr(testAccRancher2UserType+".foo", "token_id"),
				),
			},
			{
				Config: testAccRancher2UserUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2UserExists(testAccRancher2UserType+".foo", user),
					resource.TestCheckResourceAttr(testAccRancher2UserType+".foo", "name", "Terraform user acceptance test - Updated"),
					resource.TestCheckResourceAttr(testAccRancher2UserType+".foo", "password", "TestACC1234567"),
					resource.TestCheckResourceAttr(testAccRancher2UserType+".foo", "enabled", "false"),
				),
			},
			{
				Config: testAccRancher2User,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2UserExists(testAccRancher2UserType+".foo", user),
					resource.TestCheckResourceAttr(testAccRancher2UserType+".foo", "name", "Terraform user acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2UserType+".foo", "password", "TestACC123456"),
					resource.TestCheckResourceAttr(testAccRancher2UserType+".foo", "enabled", "true"),
				),
			},
		},
	})
}

func TestAccRancher2User_token(t *testing.T) {
	var user *managementClient.User

	tokenNameCheck, _ := regexp.Compile("token-\\w+")

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2UserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2UserWithToken,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2UserExists(testAccRancher2UserType+".foo", user),
					resource.TestCheckResourceAttr(testAccRancher2UserType+".foo", "name", "Terraform user acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2UserType+".foo", "password", "TestACC123456"),
					resource.TestCheckResourceAttr(testAccRancher2UserType+".foo", "enabled", "true"),
					testAccCheckRancher2UserToken(testAccRancher2UserType+".foo"),
					resource.TestMatchResourceAttr(testAccRancher2UserType+".foo", "token_name", tokenNameCheck),
					resource.TestCheckResourceAttr(testAccRancher2UserType+".foo", "token_enabled", "true"),
					resource.TestCheckResourceAttr(testAccRancher2UserType+".foo", "token_expired", "false"),
					resource.TestCheckResourceAttrSet(testAccRancher2UserType+".foo", "auth_token"),
					resource.TestCheckResourceAttrSet(testAccRancher2UserType+".foo", "access_key"),
					resource.TestCheckResourceAttrSet(testAccRancher2UserType+".foo", "secret_key"),
				),
			},
			{
				Config: testAccRancher2UserWithTokenUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2UserExists(testAccRancher2UserType+".foo", user),
					testAccCheckRancher2UserTokenChanged(testAccRancher2UserType+".foo"),
				),
			},
			{
				Config: testAccRancher2User,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2UserExists(testAccRancher2UserType+".foo", user),
					resource.TestCheckResourceAttr(testAccRancher2UserType+".foo", "token_id", ""),
					resource.TestCheckResourceAttr(testAccRancher2UserType+".foo", "token_name", ""),
					resource.TestCheckResourceAttr(testAccRancher2UserType+".foo", "auth_token", ""),
					resource.TestCheckResourceAttr(testAccRancher2UserType+".foo", "access_key", ""),
					resource.TestCheckResourceAttr(testAccRancher2UserType+".foo", "secret_key", ""),
				),
			},
		},
	})
}

func TestAccRancher2User_disappears(t *testing.T) {
	var user *managementClient.User

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2UserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2User,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2UserExists(testAccRancher2UserType+".foo", user),
					testAccRancher2UserDisappears(user),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2UserDisappears(pro *managementClient.User) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2UserType {
				continue
			}
			client, err := testAccProvider.Meta().(*Config).ManagementClient()
			if err != nil {
				return err
			}

			pro, err = client.User.ByID(rs.Primary.ID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = client.User.Delete(pro)
			if err != nil {
				return fmt.Errorf("Error removing User: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Pending:    []string{"active"},
				Target:     []string{"removed"},
				Refresh:    userStateRefreshFunc(client, pro.ID),
				Timeout:    10 * time.Minute,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf(
					"[ERROR] waiting for User (%s) to be removed: %s", pro.ID, waitErr)
			}
		}
		return nil

	}
}

func testAccCheckRancher2UserExists(n string, pro *managementClient.User) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No User ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		foundPro, err := client.User.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("User not found")
			}
			return err
		}

		pro = foundPro

		return nil
	}
}

func testAccCheckRancher2UserToken(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.Attributes["token_id"] == "" {
			return fmt.Errorf("No Token ID is set")
		}

		token = &managementClient.Token{
			Resource: types.Resource{
				ID: rs.Primary.Attributes["token_id"],
			},
			Name:  rs.Primary.Attributes["token_name"],
			Token: rs.Primary.Attributes["auth_token"],
		}

		return nil
	}
}

func testAccCheckRancher2UserTokenChanged(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.Attributes["token_id"] == "" {
			return fmt.Errorf("No Token ID is set")
		}

		if rs.Primary.Attributes["token_id"] == token.ID {
			return fmt.Errorf("Token ID has not changed")
		}

		if rs.Primary.Attributes["token_name"] == token.Name {
			return fmt.Errorf("Token name has not changed")
		}

		if rs.Primary.Attributes["auth_token"] == token.Token {
			return fmt.Errorf("Token value has not changed")
		}

		key := strings.Split(token.Token, ":")
		if rs.Primary.Attributes["access_key"] == key[0] {
			return fmt.Errorf("Token access key has not changed")
		}
		if rs.Primary.Attributes["secret_key"] == key[1] {
			return fmt.Errorf("Token secret key has not changed")
		}

		return nil
	}
}

func testAccCheckRancher2UserDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2UserType {
			continue
		}
		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		_, err = client.User.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("User still exists")
	}
	return nil
}
