package rancher2

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	managementClient "github.com/rancher/types/client/management/v3"
)

const (
	testAccRancher2TokenType = "rancher2_token"
)

var (
	testAccRancher2TokenUserConfig            string
	testAccRancher2TokenConfig                string
	testAccRancher2TokenUpdateConfig          string
	testAccRancher2TokenRecreateConfig        string
	testAccRancher2TokenClusterConfig         string
	testAccRancher2TokenClusterUpdateConfig   string
	testAccRancher2TokenClusterRecreateConfig string
)

func init() {
	testAccRancher2TokenConfig = `
resource "rancher2_token" "foo" {
  description = "Terraform token acceptance test"
  ttl = 120
}
`

	testAccRancher2TokenUpdateConfig = `
resource "rancher2_token" "foo" {
  description = "Terraform token acceptance test - Updated"
  ttl = 120
}
 `

	testAccRancher2TokenRecreateConfig = `
resource "rancher2_token" "foo" {
  description = "Terraform token acceptance test"
  ttl = 120
}
 `
	testAccRancher2TokenClusterConfig = `
resource "rancher2_token" "foo" {
  cluster_id = "` + testAccRancher2ClusterID + `"
  description = "Terraform token acceptance test"
  ttl = 120
}
`

	testAccRancher2TokenClusterUpdateConfig = `
resource "rancher2_token" "foo" {
  cluster_id = "` + testAccRancher2ClusterID + `"
  description = "Terraform token acceptance test - Updated"
  ttl = 120
}
 `

	testAccRancher2TokenClusterRecreateConfig = `
resource "rancher2_token" "foo" {
  cluster_id = "` + testAccRancher2ClusterID + `"
  description = "Terraform token acceptance test"
  ttl = 120
}
 `
}

func TestAccRancher2Token_basic(t *testing.T) {
	var token *managementClient.Token

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2TokenDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2TokenConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2TokenExists(testAccRancher2TokenType+".foo", token),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo", "description", "Terraform token acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo", "ttl", "120"),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo", "enabled", "true"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2TokenUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2TokenExists(testAccRancher2TokenType+".foo", token),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo", "description", "Terraform token acceptance test - Updated"),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo", "ttl", "120"),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo", "enabled", "true"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2TokenRecreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2TokenExists(testAccRancher2TokenType+".foo", token),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo", "description", "Terraform token acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo", "ttl", "120"),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo", "enabled", "true"),
				),
			},
		},
	})
}

func TestAccRancher2Token_disappears(t *testing.T) {
	var token *managementClient.Token

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2TokenDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2TokenConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2TokenExists(testAccRancher2TokenType+".foo", token),
					testAccRancher2TokenDisappears(token),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRancher2TokenScoped_basic(t *testing.T) {
	var token *managementClient.Token

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2TokenDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2TokenClusterConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2TokenExists(testAccRancher2TokenType+".foo", token),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo", "description", "Terraform token acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo", "ttl", "120"),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo", "enabled", "true"),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo", "cluster_id", testAccRancher2ClusterID),
				),
			},
			resource.TestStep{
				Config: testAccRancher2TokenClusterUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2TokenExists(testAccRancher2TokenType+".foo", token),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo", "description", "Terraform token acceptance test - Updated"),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo", "ttl", "120"),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo", "enabled", "true"),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo", "cluster_id", testAccRancher2ClusterID),
				),
			},
			resource.TestStep{
				Config: testAccRancher2TokenClusterRecreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2TokenExists(testAccRancher2TokenType+".foo", token),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo", "description", "Terraform token acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo", "ttl", "120"),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo", "enabled", "true"),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo", "cluster_id", testAccRancher2ClusterID),
				),
			},
		},
	})
}

func TestAccRancher2TokenScoped_disappears(t *testing.T) {
	var token *managementClient.Token

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2TokenDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2TokenClusterConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2TokenExists(testAccRancher2TokenType+".foo", token),
					testAccRancher2TokenDisappears(token),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2TokenDisappears(token *managementClient.Token) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2TokenType {
				continue
			}
			client, err := testAccProvider.Meta().(*Config).ManagementClient()
			if err != nil {
				return err
			}

			token, err = client.Token.ByID(rs.Primary.ID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = client.Token.Delete(token)
			if err != nil {
				return fmt.Errorf("Error removing Token: %s", err)
			}
		}
		return nil

	}
}

func testAccCheckRancher2TokenExists(n string, token *managementClient.Token) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Token ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		foundToken, err := client.Token.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("Token not found")
			}
			return err
		}

		token = foundToken

		return nil
	}
}

func testAccCheckRancher2TokenDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2TokenType {
			continue
		}
		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		_, err = client.Token.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("Token still exists")
	}
	return nil
}
