package rancher2

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

const (
	testAccRancher2TokenType = "rancher2_token"
)

var (
	testAccRancher2Token                    string
	testAccRancher2TokenUpdate              string
	testAccRancher2TokenCluster             string
	testAccRancher2TokenClusterUpdate       string
	testAccRancher2TokenClusterConfig       string
	testAccRancher2TokenClusterUpdateConfig string
)

func init() {
	testAccRancher2Token = `
resource "` + testAccRancher2TokenType + `" "foo" {
  description = "Terraform token acceptance test"
  ttl = 120
}
`
	testAccRancher2TokenUpdate = `
resource "` + testAccRancher2TokenType + `" "foo" {
  description = "Terraform token acceptance test - Updated"
  ttl = 120
}
 `
	testAccRancher2TokenCluster = `
resource "` + testAccRancher2TokenType + `" "foo-cluster" {
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  description = "Terraform token acceptance test"
  ttl = 120
}
`
	testAccRancher2TokenClusterUpdate = `
resource "` + testAccRancher2TokenType + `" "foo-cluster" {
  cluster_id = rancher2_cluster_sync.testacc.cluster_id
  description = "Terraform token acceptance test - Updated"
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
			{
				Config: testAccRancher2Token,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2TokenExists(testAccRancher2TokenType+".foo", token),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo", "description", "Terraform token acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo", "ttl", "120"),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo", "enabled", "true"),
				),
			},
			{
				Config: testAccRancher2TokenUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2TokenExists(testAccRancher2TokenType+".foo", token),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo", "description", "Terraform token acceptance test - Updated"),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo", "ttl", "120"),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo", "enabled", "true"),
				),
			},
			{
				Config: testAccRancher2Token,
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
			{
				Config: testAccRancher2Token,
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

	testAccRancher2TokenClusterConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2TokenCluster
	testAccRancher2TokenClusterUpdateConfig = testAccCheckRancher2ClusterSyncTestacc + testAccRancher2TokenClusterUpdate

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2TokenDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2TokenClusterConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2TokenExists(testAccRancher2TokenType+".foo-cluster", token),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo-cluster", "description", "Terraform token acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo-cluster", "ttl", "120"),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo-cluster", "enabled", "true"),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo-cluster", "cluster_id", testAccRancher2ClusterID),
				),
			},
			{
				Config: testAccRancher2TokenClusterUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2TokenExists(testAccRancher2TokenType+".foo-cluster", token),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo-cluster", "description", "Terraform token acceptance test - Updated"),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo-cluster", "ttl", "120"),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo-cluster", "enabled", "true"),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo-cluster", "cluster_id", testAccRancher2ClusterID),
				),
			},
			{
				Config: testAccRancher2TokenClusterConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2TokenExists(testAccRancher2TokenType+".foo-cluster", token),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo-cluster", "description", "Terraform token acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo-cluster", "ttl", "120"),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo-cluster", "enabled", "true"),
					resource.TestCheckResourceAttr(testAccRancher2TokenType+".foo-cluster", "cluster_id", testAccRancher2ClusterID),
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
			{
				Config: testAccRancher2TokenClusterConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2TokenExists(testAccRancher2TokenType+".foo-cluster", token),
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
