package rancher2

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	norman "github.com/rancher/norman/types"
	managementClient "github.com/rancher/types/client/management/v3"
)

const (
	testAccRancher2NodeTemplateType            = "rancher2_node_template"
	testAccRancher2NodeTemplateConfigAmazonec2 = `
resource "rancher2_node_template" "foo" {
  name = "foo"
  description = "Terraform node driver amazonec2 acceptance test"
  amazonec2_config {
	access_key = "XXXXXXXXXXXXXXXXXXXX"
	secret_key = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	ami =  "ami-XXXXXXXXXXXXXXX"
	region = "XX-west-1"
	security_group = ["XXXXXXXX"]
	subnet_id = "subnet-XXXXXXXX"
	vpc_id = "vpc-XXXXXXXX"
	zone = "a"
  }
}
`
	testAccRancher2NodeTemplateUpdateConfigAmazonec2 = `
resource "rancher2_node_template" "foo" {
  name = "foo2"
  description = "Terraform node driver amazonec2 acceptance test - updated"
  amazonec2_config {
	access_key = "XXXXXXXXXXXXXXXXXXXX"
	secret_key = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	ami =  "ami-YYYYYYYYYYYYYYY"
	region = "XX-west-1"
	security_group = ["XXXXXXXX"]
	subnet_id = "subnet-XXXXXXXX"
	vpc_id = "vpc-XXXXXXXX"
	zone = "a"
  }
}
 `
	testAccRancher2NodeTemplateRecreateConfigAmazonec2 = `
resource "rancher2_node_template" "foo" {
  name = "foo"
  description = "Terraform node driver amazonec2 acceptance test"
  amazonec2_config {
	access_key = "XXXXXXXXXXXXXXXXXXXX"
	secret_key = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	ami =  "ami-XXXXXXXXXXXXXXX"
	region = "XX-west-1"
	security_group = ["XXXXXXXX"]
	subnet_id = "subnet-XXXXXXXX"
	vpc_id = "vpc-XXXXXXXX"
	zone = "a"
  }
}
`
	testAccRancher2NodeTemplateConfigAzure = `
resource "rancher2_node_template" "foo" {
  name = "foo"
  description = "Terraform node driver azure acceptance test"
  azure_config {
	client_id = "XXXXXXXXXXXXXXXXXXXX"
    client_secret = "XXXXXXXXXXXXXXXXXXXX"
    subscription_id = "XXXXXXXXXXXXXXXXXXXX"
  }
}
`
	testAccRancher2NodeTemplateUpdateConfigAzure = `
resource "rancher2_node_template" "foo" {
  name = "foo2"
  description = "Terraform node driver azure acceptance test - updated"
  azure_config {
	client_id =  "YYYYYYYYYYYYYYYYYYYY"
    client_secret = "XXXXXXXXXXXXXXXXXXXX"
    subscription_id = "XXXXXXXXXXXXXXXXXXXX"
  }
}
 `
	testAccRancher2NodeTemplateRecreateConfigAzure = `
resource "rancher2_node_template" "foo" {
  name = "foo"
  description = "Terraform node driver azure acceptance test"
  azure_config {
	client_id = "XXXXXXXXXXXXXXXXXXXX"
    client_secret = "XXXXXXXXXXXXXXXXXXXX"
    subscription_id = "XXXXXXXXXXXXXXXXXXXX"
  }
}
`
	testAccRancher2NodeTemplateConfigDigitalocean = `
resource "rancher2_node_template" "foo" {
  name = "foo"
  description = "Terraform node driver digitalocean acceptance test"
  digitalocean_config {
	access_token = "XXXXXXXXXXXXXXXXXXXX"
  }
}
`
	testAccRancher2NodeTemplateUpdateConfigDigitalocean = `
resource "rancher2_node_template" "foo" {
  name = "foo2"
  description = "Terraform node driver digitalocean acceptance test - updated"
  digitalocean_config {
	access_token =  "YYYYYYYYYYYYYYYYYYYY"
  }
}
 `
	testAccRancher2NodeTemplateRecreateConfigDigitalocean = `
resource "rancher2_node_template" "foo" {
  name = "foo"
  description = "Terraform node driver digitalocean acceptance test"
  digitalocean_config {
	access_token = "XXXXXXXXXXXXXXXXXXXX"
  }
}
`
)

func TestAccRancher2NodeTemplate_basic_Amazonec2(t *testing.T) {
	var nodeTemplate *NodeTemplate

	name := testAccRancher2NodeTemplateType + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NodeTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2NodeTemplateConfigAmazonec2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(name, nodeTemplate),
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "description", "Terraform node driver amazonec2 acceptance test"),
					resource.TestCheckResourceAttr(name, "driver", amazonec2ConfigDriver),
					resource.TestCheckResourceAttr(name, "amazonec2_config.0.ami", "ami-XXXXXXXXXXXXXXX"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NodeTemplateUpdateConfigAmazonec2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(name, nodeTemplate),
					resource.TestCheckResourceAttr(name, "name", "foo2"),
					resource.TestCheckResourceAttr(name, "description", "Terraform node driver amazonec2 acceptance test - updated"),
					resource.TestCheckResourceAttr(name, "driver", amazonec2ConfigDriver),
					resource.TestCheckResourceAttr(name, "amazonec2_config.0.ami", "ami-YYYYYYYYYYYYYYY"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NodeTemplateRecreateConfigAmazonec2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(name, nodeTemplate),
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "description", "Terraform node driver amazonec2 acceptance test"),
					resource.TestCheckResourceAttr(name, "driver", amazonec2ConfigDriver),
					resource.TestCheckResourceAttr(name, "amazonec2_config.0.ami", "ami-XXXXXXXXXXXXXXX"),
				),
			},
		},
	})
}

func TestAccRancher2NodeTemplate_disappears_Amazonec2(t *testing.T) {
	var nodeTemplate *NodeTemplate

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NodeTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2NodeTemplateConfigAmazonec2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(testAccRancher2NodeTemplateType+".foo", nodeTemplate),
					testAccRancher2NodeTemplateDisappears(nodeTemplate),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRancher2NodeTemplate_basic_Azure(t *testing.T) {
	var nodeTemplate *NodeTemplate

	name := testAccRancher2NodeTemplateType + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NodeTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2NodeTemplateConfigAzure,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(name, nodeTemplate),
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "description", "Terraform node driver azure acceptance test"),
					resource.TestCheckResourceAttr(name, "driver", azureConfigDriver),
					resource.TestCheckResourceAttr(name, "azure_config.0.client_id", "XXXXXXXXXXXXXXXXXXXX"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NodeTemplateUpdateConfigAzure,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(name, nodeTemplate),
					resource.TestCheckResourceAttr(name, "name", "foo2"),
					resource.TestCheckResourceAttr(name, "description", "Terraform node driver azure acceptance test - updated"),
					resource.TestCheckResourceAttr(name, "driver", azureConfigDriver),
					resource.TestCheckResourceAttr(name, "azure_config.0.client_id", "YYYYYYYYYYYYYYYYYYYY"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NodeTemplateRecreateConfigAzure,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(name, nodeTemplate),
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "description", "Terraform node driver azure acceptance test"),
					resource.TestCheckResourceAttr(name, "driver", azureConfigDriver),
					resource.TestCheckResourceAttr(name, "azure_config.0.client_id", "XXXXXXXXXXXXXXXXXXXX"),
				),
			},
		},
	})
}

func TestAccRancher2NodeTemplate_disappears_Azure(t *testing.T) {
	var nodeTemplate *NodeTemplate

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NodeTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2NodeTemplateConfigAzure,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(testAccRancher2NodeTemplateType+".foo", nodeTemplate),
					testAccRancher2NodeTemplateDisappears(nodeTemplate),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRancher2NodeTemplate_basic_Digitalocean(t *testing.T) {
	var nodeTemplate *NodeTemplate

	name := testAccRancher2NodeTemplateType + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NodeTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2NodeTemplateConfigDigitalocean,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(name, nodeTemplate),
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "description", "Terraform node driver digitalocean acceptance test"),
					resource.TestCheckResourceAttr(name, "driver", digitaloceanConfigDriver),
					resource.TestCheckResourceAttr(name, "digitalocean_config.0.access_token", "XXXXXXXXXXXXXXXXXXXX"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NodeTemplateUpdateConfigDigitalocean,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(name, nodeTemplate),
					resource.TestCheckResourceAttr(name, "name", "foo2"),
					resource.TestCheckResourceAttr(name, "description", "Terraform node driver digitalocean acceptance test - updated"),
					resource.TestCheckResourceAttr(name, "driver", digitaloceanConfigDriver),
					resource.TestCheckResourceAttr(name, "digitalocean_config.0.access_token", "YYYYYYYYYYYYYYYYYYYY"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NodeTemplateRecreateConfigDigitalocean,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(name, nodeTemplate),
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "description", "Terraform node driver digitalocean acceptance test"),
					resource.TestCheckResourceAttr(name, "driver", digitaloceanConfigDriver),
					resource.TestCheckResourceAttr(name, "digitalocean_config.0.access_token", "XXXXXXXXXXXXXXXXXXXX"),
				),
			},
		},
	})
}

func TestAccRancher2NodeTemplate_disappears_Digitalocean(t *testing.T) {
	var nodeTemplate *NodeTemplate

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NodeTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2NodeTemplateConfigDigitalocean,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(testAccRancher2NodeTemplateType+".foo", nodeTemplate),
					testAccRancher2NodeTemplateDisappears(nodeTemplate),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2NodeTemplateDisappears(nodeTemplate *NodeTemplate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2NodeTemplateType {
				continue
			}
			client, err := testAccProvider.Meta().(*Config).ManagementClient()
			if err != nil {
				return err
			}

			nodeTemplate := &norman.Resource{}
			err = client.APIBaseClient.ByID(managementClient.NodeTemplateType, rs.Primary.ID, nodeTemplate)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = client.APIBaseClient.Delete(nodeTemplate)
			if err != nil {
				return fmt.Errorf("Error removing Node Template: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Pending:    []string{"removing"},
				Target:     []string{"removed"},
				Refresh:    nodeTemplateStateRefreshFunc(client, nodeTemplate.ID),
				Timeout:    10 * time.Minute,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf("[ERROR] waiting for node template (%s) to be removed: %s", nodeTemplate.ID, waitErr)
			}
		}
		return nil
	}
}

func testAccCheckRancher2NodeTemplateExists(n string, nodeTemplate *NodeTemplate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Node Template ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		foundNodeTemplate := &NodeTemplate{}
		err = client.APIBaseClient.ByID(managementClient.NodeTemplateType, rs.Primary.ID, foundNodeTemplate)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("Node Template not found")
			}
			return err
		}

		nodeTemplate = foundNodeTemplate

		return nil
	}
}

func testAccCheckRancher2NodeTemplateDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2NodeTemplateType {
			continue
		}
		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		foundNodeTemplate := &NodeTemplate{}
		err = client.APIBaseClient.ByID(managementClient.NodeTemplateType, rs.Primary.ID, foundNodeTemplate)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("Node Template still exists")
	}
	return nil
}
