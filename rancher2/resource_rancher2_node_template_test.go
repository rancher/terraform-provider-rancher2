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
	testAccRancher2NodeTemplateType      = "rancher2_node_template"
	testAccRancher2NodeTemplateAmazonec2 = `
resource "` + testAccRancher2NodeTemplateType + `" "foo-aws" {
  name = "foo-aws"
  description = "Terraform node driver amazonec2 acceptance test"
  cloud_credential_id = rancher2_cloud_credential.foo-aws.id
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
	testAccRancher2NodeTemplateAmazonec2Update = `
resource "` + testAccRancher2NodeTemplateType + `" "foo-aws" {
  name = "foo-aws2"
  description = "Terraform node driver amazonec2 acceptance test - updated"
  cloud_credential_id = rancher2_cloud_credential.foo-aws.id
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
	testAccRancher2NodeTemplateAmazonec2Config       = testAccRancher2CloudCredentialConfigAmazonec2 + testAccRancher2NodeTemplateAmazonec2
	testAccRancher2NodeTemplateAmazonec2UpdateConfig = testAccRancher2CloudCredentialConfigAmazonec2 + testAccRancher2NodeTemplateAmazonec2Update
	testAccRancher2NodeTemplateAzure                 = `
resource "` + testAccRancher2NodeTemplateType + `" "foo-azure" {
  name = "foo-azure"
  description = "Terraform node driver azure acceptance test"
  cloud_credential_id = rancher2_cloud_credential.foo-azure.id
  azure_config {
	image =  "image-XXXXXXXX"
	location =  "location-XXXXXXXX"
  }
}
`
	testAccRancher2NodeTemplateAzureUpdate = `
resource "` + testAccRancher2NodeTemplateType + `" "foo-azure" {
  name = "foo-azure2"
  description = "Terraform node driver azure acceptance test - updated"
  cloud_credential_id = rancher2_cloud_credential.foo-azure.id
  azure_config {
	image =  "image-YYYYYYYY"
	location =  "location-YYYYYYYY"
  }
}
`
	testAccRancher2NodeTemplateAzureConfig       = testAccRancher2CloudCredentialConfigAzure + testAccRancher2NodeTemplateAzure
	testAccRancher2NodeTemplateAzureUpdateConfig = testAccRancher2CloudCredentialConfigAzure + testAccRancher2NodeTemplateAzureUpdate
	testAccRancher2NodeTemplateDigitalocean      = `
resource "` + testAccRancher2NodeTemplateType + `" "foo-do" {
  name = "foo-do"
  description = "Terraform node driver digitalocean acceptance test"
  cloud_credential_id = rancher2_cloud_credential.foo-do.id
  digitalocean_config {
	image =  "image-XXXXXXXX"
	region =  "region-XXXXXXXX"
  }
}
`
	testAccRancher2NodeTemplateDigitaloceanUpdate = `
resource "` + testAccRancher2NodeTemplateType + `" "foo-do" {
  name = "foo-do2"
  description = "Terraform node driver digitalocean acceptance test - updated"
  cloud_credential_id = rancher2_cloud_credential.foo-do.id
  digitalocean_config {
	image =  "image-YYYYYYYY"
	region =  "region-YYYYYYYY"
  }
}
`
	testAccRancher2NodeTemplateDigitaloceanConfig       = testAccRancher2CloudCredentialConfigDigitalocean + testAccRancher2NodeTemplateDigitalocean
	testAccRancher2NodeTemplateDigitaloceanUpdateConfig = testAccRancher2CloudCredentialConfigDigitalocean + testAccRancher2NodeTemplateDigitaloceanUpdate
	testAccRancher2NodeTemplateOpennebulaDriver         = `
resource "rancher2_node_driver" "foo-opennebula" {
    active = true
    builtin = false
    name = "opennebula"
    url = "https://github.com/OpenNebula/docker-machine-opennebula/releases/download/release-0.2.0/docker-machine-driver-opennebula.tgz"
}
`
	testAccRancher2NodeTemplateOpennebula = `
resource "` + testAccRancher2NodeTemplateType + `" "foo-opennebula" {
  name = "foo-opennebula"
  description = "Terraform node template opennebula acceptance test"
  driver_id = rancher2_node_driver.foo-opennebula.id
  opennebula_config {
	user = "apiuser"
	password =  "password123"
	ssh_user = "rancher"
	template_name = "template-YYYYYYYY"
	xml_rpc_url = "http://XXXXXXXX/RPC2"
  }
}
`
	testAccRancher2NodeTemplateOpennebulaUpdate = `
resource "` + testAccRancher2NodeTemplateType + `" "foo-opennebula" {
  name = "foo-opennebula2"
  description = "Terraform node template opennebula acceptance test - updated"
  driver_id = rancher2_node_driver.foo-opennebula.id
  opennebula_config {
  	user = "apiuser"
	password =  "password123"
	ssh_user = "rancher2"
	template_name = "template-XXXXXXXX"
	xml_rpc_url = "http://XXXXXXXX/RPC2"
  }
}
`
	testAccRancher2NodeTemplateOpennebulaConfig       = testAccRancher2NodeTemplateOpennebulaDriver + testAccRancher2NodeTemplateOpennebula
	testAccRancher2NodeTemplateOpennebulaUpdateConfig = testAccRancher2NodeTemplateOpennebulaDriver + testAccRancher2NodeTemplateOpennebulaUpdate
	testAccRancher2NodeTemplateOpenstack              = `
resource "` + testAccRancher2NodeTemplateType + `" "foo-openstack" {
  name = "foo-openstack"
  description = "Terraform node driver openstack acceptance test"
  cloud_credential_id = rancher2_cloud_credential.foo-openstack.id
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
	testAccRancher2NodeTemplateOpenstackUpdate = `
resource "` + testAccRancher2NodeTemplateType + `" "foo-openstack" {
  name = "foo-openstack2"
  description = "Terraform node driver openstack acceptance test - updated"
  cloud_credential_id = "${rancher2_cloud_credential.foo-openstack.id}"
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
	testAccRancher2NodeTemplateOpenstackConfig       = testAccRancher2CloudCredentialConfigOpenstack + testAccRancher2NodeTemplateOpenstack
	testAccRancher2NodeTemplateOpenstackUpdateConfig = testAccRancher2CloudCredentialConfigOpenstack + testAccRancher2NodeTemplateOpenstackUpdate
	testAccRancher2NodeTemplateVsphere               = `
resource "` + testAccRancher2NodeTemplateType + `" "foo-vsphere" {
  name = "foo-vsphere"
  description = "Terraform node driver vsphere acceptance test"
  cloud_credential_id = rancher2_cloud_credential.foo-vsphere.id
  vsphere_config {
    cpu_count = "4"
	disk_size = "10240"
	pool =  "pool-XXXXXXXX"
  }
}
`
	testAccRancher2NodeTemplateVsphereUpdate = `
resource "` + testAccRancher2NodeTemplateType + `" "foo-vsphere" {
  name = "foo-vsphere2"
  description = "Terraform node driver vsphere acceptance test - updated"
  cloud_credential_id = rancher2_cloud_credential.foo-vsphere.id
  vsphere_config {
	cpu_count = "8"
	disk_size = "20480"
	pool =  "pool-YYYYYYYY"
  }
}
`
	testAccRancher2NodeTemplateVsphereConfig       = testAccRancher2CloudCredentialConfigVsphere + testAccRancher2NodeTemplateVsphere
	testAccRancher2NodeTemplateVsphereUpdateConfig = testAccRancher2CloudCredentialConfigVsphere + testAccRancher2NodeTemplateVsphereUpdate
)

func TestAccRancher2NodeTemplate_basic_Amazonec2(t *testing.T) {
	var nodeTemplate *NodeTemplate

	name := testAccRancher2NodeTemplateType + ".foo-aws"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NodeTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2NodeTemplateAmazonec2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(name, nodeTemplate),
					resource.TestCheckResourceAttr(name, "name", "foo-aws"),
					resource.TestCheckResourceAttr(name, "description", "Terraform node driver amazonec2 acceptance test"),
					resource.TestCheckResourceAttr(name, "driver", amazonec2ConfigDriver),
					resource.TestCheckResourceAttr(name, "amazonec2_config.0.ami", "ami-XXXXXXXXXXXXXXX"),
					resource.TestCheckResourceAttr(name, "amazonec2_config.0.subnet_id", "subnet-XXXXXXXX"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NodeTemplateAmazonec2UpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(name, nodeTemplate),
					resource.TestCheckResourceAttr(name, "name", "foo-aws2"),
					resource.TestCheckResourceAttr(name, "description", "Terraform node driver amazonec2 acceptance test - updated"),
					resource.TestCheckResourceAttr(name, "driver", amazonec2ConfigDriver),
					resource.TestCheckResourceAttr(name, "amazonec2_config.0.ami", "ami-YYYYYYYYYYYYYYY"),
					resource.TestCheckResourceAttr(name, "amazonec2_config.0.subnet_id", "subnet-YYYYYYYY"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NodeTemplateAmazonec2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(name, nodeTemplate),
					resource.TestCheckResourceAttr(name, "name", "foo-aws"),
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
				Config: testAccRancher2NodeTemplateAmazonec2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(testAccRancher2NodeTemplateType+".foo-aws", nodeTemplate),
					testAccRancher2NodeTemplateDisappears(nodeTemplate),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRancher2NodeTemplate_basic_Azure(t *testing.T) {
	var nodeTemplate *NodeTemplate

	name := testAccRancher2NodeTemplateType + ".foo-azure"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NodeTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2NodeTemplateAzureConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(name, nodeTemplate),
					resource.TestCheckResourceAttr(name, "name", "foo-azure"),
					resource.TestCheckResourceAttr(name, "description", "Terraform node driver azure acceptance test"),
					resource.TestCheckResourceAttr(name, "driver", azureConfigDriver),
					resource.TestCheckResourceAttr(name, "azure_config.0.image", "image-XXXXXXXX"),
					resource.TestCheckResourceAttr(name, "azure_config.0.location", "location-XXXXXXXX"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NodeTemplateAzureUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(name, nodeTemplate),
					resource.TestCheckResourceAttr(name, "name", "foo-azure2"),
					resource.TestCheckResourceAttr(name, "description", "Terraform node driver azure acceptance test - updated"),
					resource.TestCheckResourceAttr(name, "driver", azureConfigDriver),
					resource.TestCheckResourceAttr(name, "azure_config.0.image", "image-YYYYYYYY"),
					resource.TestCheckResourceAttr(name, "azure_config.0.location", "location-YYYYYYYY"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NodeTemplateAzureConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(name, nodeTemplate),
					resource.TestCheckResourceAttr(name, "name", "foo-azure"),
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
				Config: testAccRancher2NodeTemplateAzureConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(testAccRancher2NodeTemplateType+".foo-azure", nodeTemplate),
					testAccRancher2NodeTemplateDisappears(nodeTemplate),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRancher2NodeTemplate_basic_Digitalocean(t *testing.T) {
	var nodeTemplate *NodeTemplate

	name := testAccRancher2NodeTemplateType + ".foo-do"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NodeTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2NodeTemplateDigitaloceanConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(name, nodeTemplate),
					resource.TestCheckResourceAttr(name, "name", "foo-do"),
					resource.TestCheckResourceAttr(name, "description", "Terraform node driver digitalocean acceptance test"),
					resource.TestCheckResourceAttr(name, "driver", digitaloceanConfigDriver),
					resource.TestCheckResourceAttr(name, "digitalocean_config.0.image", "image-XXXXXXXX"),
					resource.TestCheckResourceAttr(name, "digitalocean_config.0.region", "region-XXXXXXXX"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NodeTemplateDigitaloceanUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(name, nodeTemplate),
					resource.TestCheckResourceAttr(name, "name", "foo-do2"),
					resource.TestCheckResourceAttr(name, "description", "Terraform node driver digitalocean acceptance test - updated"),
					resource.TestCheckResourceAttr(name, "driver", digitaloceanConfigDriver),
					resource.TestCheckResourceAttr(name, "digitalocean_config.0.image", "image-YYYYYYYY"),
					resource.TestCheckResourceAttr(name, "digitalocean_config.0.region", "region-YYYYYYYY"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NodeTemplateDigitaloceanConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(name, nodeTemplate),
					resource.TestCheckResourceAttr(name, "name", "foo-do"),
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
				Config: testAccRancher2NodeTemplateDigitaloceanConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(testAccRancher2NodeTemplateType+".foo-do", nodeTemplate),
					testAccRancher2NodeTemplateDisappears(nodeTemplate),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRancher2NodeTemplate_basic_Opennebula(t *testing.T) {
	var nodeTemplate *NodeTemplate

	name := testAccRancher2NodeTemplateType + ".foo-opennebula"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NodeTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2NodeTemplateOpennebulaConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(name, nodeTemplate),
					resource.TestCheckResourceAttr(name, "name", "foo-opennebula"),
					resource.TestCheckResourceAttr(name, "description", "Terraform node template opennebula acceptance test"),
					resource.TestCheckResourceAttr(name, "driver", opennebulaConfigDriver),
					resource.TestCheckResourceAttr(name, "opennebula_config.0.template_name", "template-YYYYYYYY"),
					resource.TestCheckResourceAttr(name, "opennebula_config.0.ssh_user", "rancher"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NodeTemplateOpennebulaUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(name, nodeTemplate),
					resource.TestCheckResourceAttr(name, "name", "foo-opennebula2"),
					resource.TestCheckResourceAttr(name, "description", "Terraform node template opennebula acceptance test - updated"),
					resource.TestCheckResourceAttr(name, "driver", opennebulaConfigDriver),
					resource.TestCheckResourceAttr(name, "opennebula_config.0.template_name", "template-XXXXXXXX"),
					resource.TestCheckResourceAttr(name, "opennebula_config.0.ssh_user", "rancher2"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NodeTemplateOpennebulaConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(name, nodeTemplate),
					resource.TestCheckResourceAttr(name, "name", "foo-opennebula"),
					resource.TestCheckResourceAttr(name, "description", "Terraform node template opennebula acceptance test"),
					resource.TestCheckResourceAttr(name, "driver", opennebulaConfigDriver),
					resource.TestCheckResourceAttr(name, "opennebula_config.0.template_name", "template-YYYYYYYY"),
					resource.TestCheckResourceAttr(name, "opennebula_config.0.ssh_user", "rancher"),
				),
			},
		},
	})
}

func TestAccRancher2NodeTemplate_disappears_Opennebula(t *testing.T) {
	var nodeTemplate *NodeTemplate

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NodeTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2NodeTemplateOpennebulaConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(testAccRancher2NodeTemplateType+".foo-opennebula", nodeTemplate),
					testAccRancher2NodeTemplateDisappears(nodeTemplate),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRancher2NodeTemplate_basic_Openstack(t *testing.T) {
	var nodeTemplate *NodeTemplate

	name := testAccRancher2NodeTemplateType + ".foo-openstack"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NodeTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2NodeTemplateOpenstackConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(name, nodeTemplate),
					resource.TestCheckResourceAttr(name, "name", "foo-openstack"),
					resource.TestCheckResourceAttr(name, "description", "Terraform node driver openstack acceptance test"),
					resource.TestCheckResourceAttr(name, "driver", openstackConfigDriver),
					resource.TestCheckResourceAttr(name, "openstack_config.0.image_name", "image-XXXXXXXX"),
					resource.TestCheckResourceAttr(name, "openstack_config.0.flavor_name", "flavor-XXXXXXXX"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NodeTemplateOpenstackUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(name, nodeTemplate),
					resource.TestCheckResourceAttr(name, "name", "foo-openstack2"),
					resource.TestCheckResourceAttr(name, "description", "Terraform node driver openstack acceptance test - updated"),
					resource.TestCheckResourceAttr(name, "driver", openstackConfigDriver),
					resource.TestCheckResourceAttr(name, "openstack_config.0.image_name", "image-YYYYYYYY"),
					resource.TestCheckResourceAttr(name, "openstack_config.0.flavor_name", "flavor-YYYYYYYY"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NodeTemplateOpenstackConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(name, nodeTemplate),
					resource.TestCheckResourceAttr(name, "name", "foo-openstack"),
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
				Config: testAccRancher2NodeTemplateOpenstackConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(testAccRancher2NodeTemplateType+".foo-openstack", nodeTemplate),
					testAccRancher2NodeTemplateDisappears(nodeTemplate),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRancher2NodeTemplate_basic_Vsphere(t *testing.T) {
	var nodeTemplate *NodeTemplate

	name := testAccRancher2NodeTemplateType + ".foo-vsphere"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NodeTemplateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2NodeTemplateVsphereConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(name, nodeTemplate),
					resource.TestCheckResourceAttr(name, "name", "foo-vsphere"),
					resource.TestCheckResourceAttr(name, "description", "Terraform node driver vsphere acceptance test"),
					resource.TestCheckResourceAttr(name, "driver", vmwarevsphereConfigDriver),
					resource.TestCheckResourceAttr(name, "vsphere_config.0.cpu_count", "4"),
					resource.TestCheckResourceAttr(name, "vsphere_config.0.disk_size", "10240"),
					resource.TestCheckResourceAttr(name, "vsphere_config.0.pool", "pool-XXXXXXXX"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NodeTemplateVsphereUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(name, nodeTemplate),
					resource.TestCheckResourceAttr(name, "name", "foo-vsphere2"),
					resource.TestCheckResourceAttr(name, "description", "Terraform node driver vsphere acceptance test - updated"),
					resource.TestCheckResourceAttr(name, "driver", vmwarevsphereConfigDriver),
					resource.TestCheckResourceAttr(name, "vsphere_config.0.cpu_count", "8"),
					resource.TestCheckResourceAttr(name, "vsphere_config.0.disk_size", "20480"),
					resource.TestCheckResourceAttr(name, "vsphere_config.0.pool", "pool-YYYYYYYY"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NodeTemplateVsphereConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(name, nodeTemplate),
					resource.TestCheckResourceAttr(name, "name", "foo-vsphere"),
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
				Config: testAccRancher2NodeTemplateVsphereConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeTemplateExists(testAccRancher2NodeTemplateType+".foo-vsphere", nodeTemplate),
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
