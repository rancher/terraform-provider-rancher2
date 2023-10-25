package rancher2

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

const (
	testAccRancher2CustomUserTokenType = "rancher2_custom_user_token"
)

var (
	testAccRancher2CustomUserToken       string
	testAccRancher2CustomUserTokenUpdate string
)

func init() {
	testAccRancher2CustomUserToken = `
resource "` + testAccRancher2CustomUserTokenType + `" "foo" {
  username = "foo"
  password = "TestACC123456"
  description = "Terraform custom user token acceptance test"
  depends_on = [
    rancher2_global_role_binding.foo
  ]
}
`
	testAccRancher2CustomUserTokenUpdate = `
resource "` + testAccRancher2CustomUserTokenType + `" "foo" {
  username = "foo"
  password = "TestACC123456"
  description = "Terraform custom user token acceptance test - updated"
  ttl = 120
  depends_on = [
    rancher2_global_role_binding.foo
  ]
}
`
}

func TestAccRancher2CustomUserToken_basic(t *testing.T) {
	var user *managementClient.User
	var globalRole *managementClient.GlobalRoleBinding
	var token *managementClient.Token

	testAccRancher2CustomUserToken = testAccRancher2User + testAccRancher2GlobalRoleBinding + testAccRancher2CustomUserToken
	testAccRancher2CustomUserTokenUpdate = testAccRancher2User + testAccRancher2GlobalRoleBinding + testAccRancher2CustomUserTokenUpdate

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2UserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2CustomUserToken,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2UserExists(testAccRancher2UserType+".foo", user),
					testAccCheckRancher2GlobalRoleBindingExists(testAccRancher2GlobalRoleBindingType+".foo", globalRole),
					testAccCheckRancher2CustomUserTokenExists(testAccRancher2CustomUserTokenType+".foo", token),
					resource.TestCheckResourceAttr(testAccRancher2CustomUserTokenType+".foo", "description", "Terraform custom user token acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2CustomUserTokenType+".foo", "ttl", "0"),
					resource.TestCheckResourceAttr(testAccRancher2CustomUserTokenType+".foo", "enabled", "true"),
				),
			},
			{
				Config: testAccRancher2CustomUserTokenUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2UserExists(testAccRancher2UserType+".foo", user),
					testAccCheckRancher2GlobalRoleBindingExists(testAccRancher2GlobalRoleBindingType+".foo", globalRole),
					testAccCheckRancher2CustomUserTokenExists(testAccRancher2CustomUserTokenType+".foo", token),
					resource.TestCheckResourceAttr(testAccRancher2CustomUserTokenType+".foo", "description", "Terraform custom user token acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2CustomUserTokenType+".foo", "ttl", "120"),
					resource.TestCheckResourceAttr(testAccRancher2CustomUserTokenType+".foo", "enabled", "true"),
				),
			},
			{
				Config: testAccRancher2CustomUserToken,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CustomUserTokenExists(testAccRancher2CustomUserTokenType+".foo", token),
				),
			},
		},
	})
}

func TestAccRancher2CustomUserToken_disappears(t *testing.T) {
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

func testAccCheckRancher2CustomUserTokenExists(n string, token *managementClient.Token) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("[ERROR] Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("[ERROR] No Token ID is set")
		}

		if testAccProviderConfig.URL == "" {
			return fmt.Errorf("[ERROR] No config")
		}

		_, tempToken, err := DoUserLogin(testAccProviderConfig.URL, "", "foo", "TestACC123456", "0", "Temp Terraform API token for ACC tests", testAccProviderConfig.CACerts, testAccProviderConfig.Insecure)
		if err != nil {
			return fmt.Errorf("[ERROR] Login with %s user: %v", "foo", err)
		}

		options := testAccProviderConfig.CreateClientOpts()
		options.URL = options.URL + rancher2ClientAPIVersion
		options.TokenKey = tempToken
		client, err := managementClient.NewClient(options)
		if err != nil {
			return err
		}

		foundToken, err := client.Token.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("[ERROR] Token not found")
			}
			return err
		}

		token = foundToken

		return nil
	}
}
