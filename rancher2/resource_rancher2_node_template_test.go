package rancher2

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	norman "github.com/rancher/norman/types"
	managementClient "github.com/rancher/types/client/management/v3"
)

const (
	testAccRancher2NodeTemplateType            = "rancher2_node_template"
	testAccRancher2NodeTemplateConfigAmazonec2 = testAccRancher2CloudCredentialConfigAmazonec2 + `
resource "rancher2_node_template" "foo" {
  name = "foo"
  description = "Terraform node driver amazonec2 acceptance test"
  cloud_credential_id = "${rancher2_cloud_credential.foo.id}"
  amazonec2_config {
	ami =  "ami-XXXXXXXXXXXXXXX"
	region = "XX-west-1"
	security_group = ["XXXXXXXX"]
	subnet_id = "subnet-XXXXXXXX"
	vpc_id = "vpc-XXXXXXXX"
	zone = "a"
  }
}
`
	testAccRancher2NodeTemplateUpdateConfigAmazonec2 = testAccRancher2CloudCredentialConfigAmazonec2 + `
resource "rancher2_node_template" "foo" {
  name = "foo2"
  description = "Terraform node driver amazonec2 acceptance test - updated"
  cloud_credential_id = "${rancher2_cloud_credential.foo.id}"
  amazonec2_config {
	ami =  "ami-YYYYYYYYYYYYYYY"
	region = "XX-west-1"
	security_group = ["XXXXXXXX"]
	subnet_id = "subnet-YYYYYYYY"
	vpc_id = "vpc-XXXXXXXX"
	zone = "a"
  }
}
 `
	testAccRancher2NodeTemplateRecreateConfigAmazonec2 = testAccRancher2CloudCredentialConfigAmazonec2 + `
resource "rancher2_node_template" "foo" {
  name = "foo"
  description = "Terraform node driver amazonec2 acceptance test"
  cloud_credential_id = "${rancher2_cloud_credential.foo.id}"
  amazonec2_config {
	ami =  "ami-XXXXXXXXXXXXXXX"
	region = "XX-west-1"
	security_group = ["XXXXXXXX"]
	subnet_id = "subnet-XXXXXXXX"
	vpc_id = "vpc-XXXXXXXX"
	zone = "a"
  }
}
`
	testAccRancher2NodeTemplateConfigAzure = testAccRancher2CloudCredentialConfigAzure + `
resource "rancher2_node_template" "foo" {
  name = "foo"
  description = "Terraform node driver azure acceptance test"
  cloud_credential_id = "${rancher2_cloud_credential.foo.id}"
  azure_config {
	image =  "image-XXXXXXXX"
	location =  "location-XXXXXXXX"
  }
}
`
	testAccRancher2NodeTemplateUpdateConfigAzure = testAccRancher2CloudCredentialConfigAzure + `
resource "rancher2_node_template" "foo" {
  name = "foo2"
  description = "Terraform node driver azure acceptance test - updated"
  cloud_credential_id = "${rancher2_cloud_credential.foo.id}"
  azure_config {
	image =  "image-YYYYYYYY"
	location =  "location-YYYYYYYY"
  }
}
 `
	testAccRancher2NodeTemplateRecreateConfigAzure = testAccRancher2CloudCredentialConfigAzure + `
resource "rancher2_node_template" "foo" {
  name = "foo"
  description = "Terraform node driver azure acceptance test"
  cloud_credential_id = "${rancher2_cloud_credential.foo.id}"
  azure_config {
	image =  "image-XXXXXXXX"
	location =  "location-XXXXXXXX"
  }
}
`
	testAccRancher2NodeTemplateConfigDigitalocean = testAccRancher2CloudCredentialConfigDigitalocean + `
resource "rancher2_node_template" "foo" {
  name = "foo"
  description = "Terraform node driver digitalocean acceptance test"
  cloud_credential_id = "${rancher2_cloud_credential.foo.id}"
  digitalocean_config {
	image =  "image-XXXXXXXX"
	region =  "region-XXXXXXXX"
  }
}
`
	testAccRancher2NodeTemplateUpdateConfigDigitalocean = testAccRancher2CloudCredentialConfigDigitalocean + `
resource "rancher2_node_template" "foo" {
  name = "foo2"
  description = "Terraform node driver digitalocean acceptance test - updated"
  cloud_credential_id = "${rancher2_cloud_credential.foo.id}"
  digitalocean_config {
	image =  "image-YYYYYYYY"
	region =  "region-YYYYYYYY"
  }
}
 `
	testAccRancher2NodeTemplateRecreateConfigDigitalocean = testAccRancher2CloudCredentialConfigDigitalocean + `
resource "rancher2_node_template" "foo" {
  name = "foo"
  description = "Terraform node driver digitalocean acceptance test"
  cloud_credential_id = "${rancher2_cloud_credential.foo.id}"
  digitalocean_config {
	image =  "image-XXXXXXXX"
	region =  "region-XXXXXXXX"
  }
}
`
	testAccRancher2NodeTemplateConfigOpenstack = testAccRancher2CloudCredentialConfigOpenstack + `
resource "rancher2_node_template" "foo" {
  name = "foo"
  description = "Terraform node driver openstack acceptance test"
  cloud_credential_id = "${rancher2_cloud_credential.foo.id}"
  openstack_config {
  	username = "user"
    image_name =  "image-XXXXXXXX"
    region = "XX-west-1"
    flavor_name = "flavor-XXXXXXXX"
    auth_url = "http://XXXXXXXX"
    availability_zone = "zone-XXXXXXXX"
  }
}
`
	testAccRancher2NodeTemplateUpdateConfigOpenstack = testAccRancher2CloudCredentialConfigOpenstack + `
resource "rancher2_node_template" "foo" {
  name = "foo2"
  description = "Terraform node driver openstack acceptance test - updated"
  cloud_credential_id = "${rancher2_cloud_credential.foo.id}"
  openstack_config {
  	username = "user"
	image_name =  "image-YYYYYYYY"
	region = "XX-west-1"
	flavor_name = "flavor-YYYYYYYY"
	auth_url = "http://XXXXXXXX"
	availability_zone = "zone-XXXXXXXX"
  }
}
`
	testAccRancher2NodeTemplateRecreateConfigOpenstack = testAccRancher2CloudCredentialConfigOpenstack + `
resource "rancher2_node_template" "foo" {
  name = "foo"
  description = "Terraform node driver openstack acceptance test"
  cloud_credential_id = "${rancher2_cloud_credential.foo.id}"
  openstack_config {
  	username = "user"
	image_name =  "image-XXXXXXXX"
	region = "XX-west-1"
	flavor_name = "flavor-XXXXXXXX"
	auth_url = "http://XXXXXXXX"
	availability_zone = "zone-XXXXXXXX"
  }
}
`
	testAccRancher2NodeTemplateConfigVsphere = testAccRancher2CloudCredentialConfigVsphere + `
resource "rancher2_node_template" "foo" {
  name = "foo"
  description = "Terraform node driver vsphere acceptance test"
  cloud_credential_id = "${rancher2_cloud_credential.foo.id}"
  vsphere_config {
    cpu_count = "4"
	disk_size = "10240"
	pool =  "pool-XXXXXXXX"
  }
}
`
	testAccRancher2NodeTemplateUpdateConfigVsphere = testAccRancher2CloudCredentialConfigVsphere + `
resource "rancher2_node_template" "foo" {
  name = "foo2"
  description = "Terraform node driver vsphere acceptance test - updated"
  cloud_credential_id = "${rancher2_cloud_credential.foo.id}"
  vsphere_config {
	cpu_count = "8"
	disk_size = "20480"
	pool =  "pool-YYYYYYYY"
  }
}
`
	testAccRancher2NodeTemplateRecreateConfigVsphere = testAccRancher2CloudCredentialConfigVsphere + `
resource "rancher2_node_template" "foo" {
  name = "foo"
  description = "Terraform node driver vsphere acceptance test"
  cloud_credential_id = "${rancher2_cloud_credential.foo.id}"
  vsphere_config {
	cpu_count = "4"
	disk_size = "10240"
	pool =  "pool-XXXXXXXX"
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
					resource.TestCheckResourceAttr(name, "amazonec2_config.0.subnet_id", "subnet-XXXXXXXX"),
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
					resource.TestCheckResourceAttr(name, "amazonec2_config.0.subnet_id", "subnet-YYYYYYYY"),
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
					resource.TestCheckResourceAttr(name, "amazonec2_config.0.subnet_id", "subnet-XXXXXXXX"),
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
					resource.TestCheckResourceAttr(name, "azure_config.0.image", "image-XXXXXXXX"),
					resource.TestCheckResourceAttr(name, "azure_config.0.location", "location-XXXXXXXX"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NodeTemplateUpdateConfigAzure,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(name, nodeTemplate),
					resource.TestCheckResourceAttr(name, "name", "foo2"),
					resource.TestCheckResourceAttr(name, "description", "Terraform node driver azure acceptance test - updated"),
					resource.TestCheckResourceAttr(name, "driver", azureConfigDriver),
					resource.TestCheckResourceAttr(name, "azure_config.0.image", "image-YYYYYYYY"),
					resource.TestCheckResourceAttr(name, "azure_config.0.location", "location-YYYYYYYY"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NodeTemplateRecreateConfigAzure,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(name, nodeTemplate),
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "description", "Terraform node driver azure acceptance test"),
					resource.TestCheckResourceAttr(name, "driver", azureConfigDriver),
					resource.TestCheckResourceAttr(name, "azure_config.0.image", "image-XXXXXXXX"),
					resource.TestCheckResourceAttr(name, "azure_config.0.location", "location-XXXXXXXX"),
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
					resource.TestCheckResourceAttr(name, "digitalocean_config.0.image", "image-XXXXXXXX"),
					resource.TestCheckResourceAttr(name, "digitalocean_config.0.region", "region-XXXXXXXX"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NodeTemplateUpdateConfigDigitalocean,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(name, nodeTemplate),
					resource.TestCheckResourceAttr(name, "name", "foo2"),
					resource.TestCheckResourceAttr(name, "description", "Terraform node driver digitalocean acceptance test - updated"),
					resource.TestCheckResourceAttr(name, "driver", digitaloceanConfigDriver),
					resource.TestCheckResourceAttr(name, "digitalocean_config.0.image", "image-YYYYYYYY"),
					resource.TestCheckResourceAttr(name, "digitalocean_config.0.region", "region-YYYYYYYY"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NodeTemplateRecreateConfigDigitalocean,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(name, nodeTemplate),
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "description", "Terraform node driver digitalocean acceptance test"),
					resource.TestCheckResourceAttr(name, "driver", digitaloceanConfigDriver),
					resource.TestCheckResourceAttr(name, "digitalocean_config.0.image", "image-XXXXXXXX"),
					resource.TestCheckResourceAttr(name, "digitalocean_config.0.region", "region-XXXXXXXX"),
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

func TestAccRancher2NodeTemplate_basic_Openstack(t *testing.T) {
	var nodeTemplate *NodeTemplate

	name := testAccRancher2NodeTemplateType + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NodeTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2NodeTemplateConfigOpenstack,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(name, nodeTemplate),
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "description", "Terraform node driver openstack acceptance test"),
					resource.TestCheckResourceAttr(name, "driver", openstackConfigDriver),
					resource.TestCheckResourceAttr(name, "openstack_config.0.image_name", "image-XXXXXXXX"),
					resource.TestCheckResourceAttr(name, "openstack_config.0.flavor_name", "flavor-XXXXXXXX"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NodeTemplateUpdateConfigOpenstack,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(name, nodeTemplate),
					resource.TestCheckResourceAttr(name, "name", "foo2"),
					resource.TestCheckResourceAttr(name, "description", "Terraform node driver openstack acceptance test - updated"),
					resource.TestCheckResourceAttr(name, "driver", openstackConfigDriver),
					resource.TestCheckResourceAttr(name, "openstack_config.0.image_name", "image-YYYYYYYY"),
					resource.TestCheckResourceAttr(name, "openstack_config.0.flavor_name", "flavor-YYYYYYYY"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NodeTemplateRecreateConfigOpenstack,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(name, nodeTemplate),
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "description", "Terraform node driver openstack acceptance test"),
					resource.TestCheckResourceAttr(name, "driver", openstackConfigDriver),
					resource.TestCheckResourceAttr(name, "openstack_config.0.image_name", "image-XXXXXXXX"),
					resource.TestCheckResourceAttr(name, "openstack_config.0.flavor_name", "flavor-XXXXXXXX"),
				),
			},
		},
	})
}

func TestAccRancher2NodeTemplate_disappears_Openstack(t *testing.T) {
	var nodeTemplate *NodeTemplate

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NodeTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2NodeTemplateConfigOpenstack,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(testAccRancher2NodeTemplateType+".foo", nodeTemplate),
					testAccRancher2NodeTemplateDisappears(nodeTemplate),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRancher2NodeTemplate_basic_Vsphere(t *testing.T) {
	var nodeTemplate *NodeTemplate

	name := testAccRancher2NodeTemplateType + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NodeTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2NodeTemplateConfigVsphere,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(name, nodeTemplate),
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "description", "Terraform node driver vsphere acceptance test"),
					resource.TestCheckResourceAttr(name, "driver", vmwarevsphereConfigDriver),
					resource.TestCheckResourceAttr(name, "vsphere_config.0.cpu_count", "4"),
					resource.TestCheckResourceAttr(name, "vsphere_config.0.disk_size", "10240"),
					resource.TestCheckResourceAttr(name, "vsphere_config.0.pool", "pool-XXXXXXXX"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NodeTemplateUpdateConfigVsphere,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(name, nodeTemplate),
					resource.TestCheckResourceAttr(name, "name", "foo2"),
					resource.TestCheckResourceAttr(name, "description", "Terraform node driver vsphere acceptance test - updated"),
					resource.TestCheckResourceAttr(name, "driver", vmwarevsphereConfigDriver),
					resource.TestCheckResourceAttr(name, "vsphere_config.0.cpu_count", "8"),
					resource.TestCheckResourceAttr(name, "vsphere_config.0.disk_size", "20480"),
					resource.TestCheckResourceAttr(name, "vsphere_config.0.pool", "pool-YYYYYYYY"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NodeTemplateRecreateConfigVsphere,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(name, nodeTemplate),
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "description", "Terraform node driver vsphere acceptance test"),
					resource.TestCheckResourceAttr(name, "driver", vmwarevsphereConfigDriver),
					resource.TestCheckResourceAttr(name, "vsphere_config.0.cpu_count", "4"),
					resource.TestCheckResourceAttr(name, "vsphere_config.0.disk_size", "10240"),
					resource.TestCheckResourceAttr(name, "vsphere_config.0.pool", "pool-XXXXXXXX"),
				),
			},
		},
	})
}

func TestAccRancher2NodeTemplate_disappears_Vsphere(t *testing.T) {
	var nodeTemplate *NodeTemplate

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NodeTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2NodeTemplateConfigVsphere,
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
