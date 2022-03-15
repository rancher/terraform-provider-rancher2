package rancher2

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	norman "github.com/rancher/norman/types"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

const (
	testAccRancher2CloudCredentialType            = "rancher2_cloud_credential"
	testAccRancher2CloudCredentialConfigAmazonec2 = `
resource "` + testAccRancher2CloudCredentialType + `" "foo-aws" {
  name = "foo-aws"
  description= "Terraform cloudCredential acceptance test"
  amazonec2_credential_config {
    access_key = "XXXXXXXXXXXXXXXXXXXX"
    secret_key = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  }
}
`
	testAccRancher2CloudCredentialUpdateConfigAmazonec2 = `
resource "` + testAccRancher2CloudCredentialType + `" "foo-aws" {
  name = "foo-aws"
  description= "Terraform cloudCredential acceptance test - updated"
  amazonec2_credential_config {
    access_key = "YYYYYYYYYYYYYYYYYYYY"
    secret_key = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  }
}
 `
	testAccRancher2CloudCredentialConfigAzure = `
resource "` + testAccRancher2CloudCredentialType + `" "foo-azure" {
  name = "foo-azure"
  description= "Terraform cloudCredential acceptance test"
  azure_credential_config {
	client_id = "XXXXXXXXXXXXXXXXXXXX"
    client_secret = "XXXXXXXXXXXXXXXXXXXX"
    subscription_id = "XXXXXXXXXXXXXXXXXXXX"
  }
}
`
	testAccRancher2CloudCredentialUpdateConfigAzure = `
resource "` + testAccRancher2CloudCredentialType + `" "foo-azure" {
  name = "foo-azure"
  description= "Terraform cloudCredential acceptance test - updated"
  azure_credential_config {
	client_id = "YYYYYYYYYYYYYYYYYYYY"
    client_secret = "XXXXXXXXXXXXXXXXXXXX"
    subscription_id = "XXXXXXXXXXXXXXXXXXXX"
  }
}
 `
	testAccRancher2CloudCredentialConfigDigitalocean = `
resource "` + testAccRancher2CloudCredentialType + `" "foo-do" {
  name = "foo-do"
  description= "Terraform cloudCredential acceptance test"
  digitalocean_credential_config {
    access_token = "XXXXXXXXXXXXXXXXXXXX"
  }
}
`
	testAccRancher2CloudCredentialUpdateConfigDigitalocean = `
resource "` + testAccRancher2CloudCredentialType + `" "foo-do" {
  name = "foo-do"
  description= "Terraform cloudCredential acceptance test - updated"
  digitalocean_credential_config {
    access_token = "YYYYYYYYYYYYYYYYYYYY"
  }
}
 `
	testAccRancher2CloudCredentialConfigGoogle = `
resource "` + testAccRancher2CloudCredentialType + `" "foo-google" {
  name = "foo-google"
  description= "Terraform cloudCredential acceptance test"
  google_credential_config {
    auth_encoded_json = "{\"auth_encoded_json\": true}"
  }
}
`
	testAccRancher2CloudCredentialUpdateConfigGoogle = `
resource "` + testAccRancher2CloudCredentialType + `" "foo-google" {
  name = "foo-google"
  description= "Terraform cloudCredential acceptance test - updated"
  google_credential_config {
    auth_encoded_json = "{\"auth_encoded_json\": false}"
  }
}
`
	testAccRancher2CloudCredentialConfigHarvester = `
resource "` + testAccRancher2CloudCredentialType + `" "foo-harvester" {
  name = "foo-harvester"
  description= "Terraform cloudCredential acceptance test"
  harvester_credential_config {
    cluster_id = "c-m-xb5qzqdk"
    cluster_type = "imported"
    kubeconfig_content = "XXXXXXXXXXXXXXXXXXXX"
  }
}
`
	testAccRancher2CloudCredentialUpdateConfigHarvester = `
resource "` + testAccRancher2CloudCredentialType + `" "foo-harvester" {
  name = "foo-harvester"
  description= "Terraform cloudCredential acceptance test - updated"
  harvester_credential_config {
    cluster_id = "c-m-xb5qzqdk"
    cluster_type = "imported"
    kubeconfig_content = "YYYYYYYYYYYYYYYYYYYY"
  }
}
`
	testAccRancher2CloudCredentialConfigOpenstack = `
resource "` + testAccRancher2CloudCredentialType + `" "foo-openstack" {
  name = "foo-openstack"
  description= "Terraform cloudCredential acceptance test"
  openstack_credential_config {
    password = "XXXXXXXXXXXXXXXXXXXX"
  }
}
`
	testAccRancher2CloudCredentialUpdateConfigOpenstack = `
resource "` + testAccRancher2CloudCredentialType + `" "foo-openstack" {
  name = "foo-openstack"
  description= "Terraform cloudCredential acceptance test - updated"
  openstack_credential_config {
    password = "YYYYYYYYYYYYYYYYYYYY"
  }
}
`
	testAccRancher2CloudCredentialConfigS3 = `
resource "` + testAccRancher2CloudCredentialType + `" "foo-s3" {
  name = "foo-s3"
  description= "Terraform cloudCredential acceptance test"
  s3_credential_config {
    access_key = "XXXXXXXXXXXXXXXXXXXX"
    secret_key = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
    default_bucket = "default_bucket"
    default_folder = "default_folder"
  }
}
`
	testAccRancher2CloudCredentialUpdateConfigS3 = `
resource "` + testAccRancher2CloudCredentialType + `" "foo-s3" {
  name = "foo-s3"
  description= "Terraform cloudCredential acceptance test - updated"
  s3_credential_config {
    access_key = "YYYYYYYYYYYYYYYYYYYY"
    secret_key = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
    default_bucket = "default_bucket_updated"
    default_folder = "default_folder_updated"
  }
}
`

	testAccRancher2CloudCredentialConfigVsphere = `
resource "` + testAccRancher2CloudCredentialType + `" "foo-vsphere" {
  name = "foo-vsphere"
  description= "Terraform cloudCredential acceptance test"
  vsphere_credential_config {
    password = "XXXXXXXXXXXXXXXXXXXX"
    username = "user"
    vcenter = "vcenter"
  }
}
`
	testAccRancher2CloudCredentialUpdateConfigVsphere = `
resource "` + testAccRancher2CloudCredentialType + `" "foo-vsphere" {
  name = "foo-vsphere"
  description= "Terraform cloudCredential acceptance test - updated"
  vsphere_credential_config {
    password = "YYYYYYYYYYYYYYYYYYYY"
    username = "user"
    vcenter = "vcenter2"
  }
}
 `
)

func TestAccRancher2CloudCredential_basic_Amazonec2(t *testing.T) {
	var cloudCredential *CloudCredential

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2CloudCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2CloudCredentialConfigAmazonec2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CloudCredentialExists(testAccRancher2CloudCredentialType+".foo-aws", cloudCredential),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-aws", "name", "foo-aws"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-aws", "description", "Terraform cloudCredential acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-aws", "driver", amazonec2ConfigDriver),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-aws", "amazonec2_credential_config.0.access_key", "XXXXXXXXXXXXXXXXXXXX"),
				),
			},
			{
				Config: testAccRancher2CloudCredentialUpdateConfigAmazonec2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CloudCredentialExists(testAccRancher2CloudCredentialType+".foo-aws", cloudCredential),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-aws", "name", "foo-aws"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-aws", "description", "Terraform cloudCredential acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-aws", "driver", amazonec2ConfigDriver),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-aws", "amazonec2_credential_config.0.access_key", "YYYYYYYYYYYYYYYYYYYY"),
				),
			},
			{
				Config: testAccRancher2CloudCredentialConfigAmazonec2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CloudCredentialExists(testAccRancher2CloudCredentialType+".foo-aws", cloudCredential),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-aws", "name", "foo-aws"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-aws", "description", "Terraform cloudCredential acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-aws", "driver", amazonec2ConfigDriver),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-aws", "amazonec2_credential_config.0.access_key", "XXXXXXXXXXXXXXXXXXXX"),
				),
			},
		},
	})
}

func TestAccRancher2CloudCredential_disappears_Amazonec2(t *testing.T) {
	var cloudCredential *CloudCredential

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2CloudCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2CloudCredentialConfigAmazonec2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CloudCredentialExists(testAccRancher2CloudCredentialType+".foo-aws", cloudCredential),
					testAccRancher2CloudCredentialDisappears(cloudCredential),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRancher2CloudCredential_basic_Azure(t *testing.T) {
	var cloudCredential *CloudCredential

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2CloudCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2CloudCredentialConfigAzure,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CloudCredentialExists(testAccRancher2CloudCredentialType+".foo-azure", cloudCredential),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-azure", "name", "foo-azure"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-azure", "description", "Terraform cloudCredential acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-azure", "driver", azureConfigDriver),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-azure", "azure_credential_config.0.client_id", "XXXXXXXXXXXXXXXXXXXX"),
				),
			},
			{
				Config: testAccRancher2CloudCredentialUpdateConfigAzure,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CloudCredentialExists(testAccRancher2CloudCredentialType+".foo-azure", cloudCredential),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-azure", "name", "foo-azure"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-azure", "description", "Terraform cloudCredential acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-azure", "driver", azureConfigDriver),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-azure", "azure_credential_config.0.client_id", "YYYYYYYYYYYYYYYYYYYY"),
				),
			},
			{
				Config: testAccRancher2CloudCredentialConfigAzure,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CloudCredentialExists(testAccRancher2CloudCredentialType+".foo-azure", cloudCredential),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-azure", "name", "foo-azure"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-azure", "description", "Terraform cloudCredential acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-azure", "driver", azureConfigDriver),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-azure", "azure_credential_config.0.client_id", "XXXXXXXXXXXXXXXXXXXX"),
				),
			},
		},
	})
}

func TestAccRancher2CloudCredential_disappears_Azure(t *testing.T) {
	var cloudCredential *CloudCredential

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2CloudCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2CloudCredentialConfigAzure,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CloudCredentialExists(testAccRancher2CloudCredentialType+".foo-azure", cloudCredential),
					testAccRancher2CloudCredentialDisappears(cloudCredential),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRancher2CloudCredential_basic_Digitalocean(t *testing.T) {
	var cloudCredential *CloudCredential

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2CloudCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2CloudCredentialConfigDigitalocean,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CloudCredentialExists(testAccRancher2CloudCredentialType+".foo-do", cloudCredential),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-do", "name", "foo-do"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-do", "description", "Terraform cloudCredential acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-do", "driver", digitaloceanConfigDriver),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-do", "digitalocean_credential_config.0.access_token", "XXXXXXXXXXXXXXXXXXXX"),
				),
			},
			{
				Config: testAccRancher2CloudCredentialUpdateConfigDigitalocean,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CloudCredentialExists(testAccRancher2CloudCredentialType+".foo-do", cloudCredential),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-do", "name", "foo-do"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-do", "description", "Terraform cloudCredential acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-do", "driver", digitaloceanConfigDriver),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-do", "digitalocean_credential_config.0.access_token", "YYYYYYYYYYYYYYYYYYYY"),
				),
			},
			{
				Config: testAccRancher2CloudCredentialConfigDigitalocean,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CloudCredentialExists(testAccRancher2CloudCredentialType+".foo-do", cloudCredential),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-do", "name", "foo-do"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-do", "description", "Terraform cloudCredential acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-do", "driver", digitaloceanConfigDriver),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-do", "digitalocean_credential_config.0.access_token", "XXXXXXXXXXXXXXXXXXXX"),
				),
			},
		},
	})
}

func TestAccRancher2CloudCredential_disappears_Digitalocean(t *testing.T) {
	var cloudCredential *CloudCredential

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2CloudCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2CloudCredentialConfigDigitalocean,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CloudCredentialExists(testAccRancher2CloudCredentialType+".foo-do", cloudCredential),
					testAccRancher2CloudCredentialDisappears(cloudCredential),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRancher2CloudCredential_basic_Google(t *testing.T) {
	var cloudCredential *CloudCredential

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2CloudCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2CloudCredentialConfigGoogle,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CloudCredentialExists(testAccRancher2CloudCredentialType+".foo-google", cloudCredential),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-google", "name", "foo-google"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-google", "description", "Terraform cloudCredential acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-google", "driver", googleConfigDriver),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-google", "google_credential_config.0.auth_encoded_json", "{\"auth_encoded_json\": true}"),
				),
			},
			{
				Config: testAccRancher2CloudCredentialUpdateConfigGoogle,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CloudCredentialExists(testAccRancher2CloudCredentialType+".foo-google", cloudCredential),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-google", "name", "foo-google"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-google", "description", "Terraform cloudCredential acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-google", "driver", googleConfigDriver),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-google", "google_credential_config.0.auth_encoded_json", "{\"auth_encoded_json\": false}"),
				),
			},
			{
				Config: testAccRancher2CloudCredentialConfigGoogle,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CloudCredentialExists(testAccRancher2CloudCredentialType+".foo-google", cloudCredential),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-google", "name", "foo-google"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-google", "description", "Terraform cloudCredential acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-google", "driver", googleConfigDriver),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-google", "google_credential_config.0.auth_encoded_json", "{\"auth_encoded_json\": true}"),
				),
			},
		},
	})
}

func TestAccRancher2CloudCredential_disappears_Google(t *testing.T) {
	var cloudCredential *CloudCredential

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2CloudCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2CloudCredentialConfigGoogle,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CloudCredentialExists(testAccRancher2CloudCredentialType+".foo-google", cloudCredential),
					testAccRancher2CloudCredentialDisappears(cloudCredential),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRancher2CloudCredential_basic_Harvester(t *testing.T) {
	var cloudCredential *CloudCredential

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2CloudCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2CloudCredentialConfigHarvester,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CloudCredentialExists(testAccRancher2CloudCredentialType+".foo-harvester", cloudCredential),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-harvester", "name", "foo-harvester"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-harvester", "description", "Terraform cloudCredential acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-harvester", "driver", harvesterConfigDriver),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-harvester", "harvester_credential_config.0.kubeconfig_content", "XXXXXXXXXXXXXXXXXXXX"),
				),
			},
			{
				Config: testAccRancher2CloudCredentialUpdateConfigHarvester,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CloudCredentialExists(testAccRancher2CloudCredentialType+".foo-harvester", cloudCredential),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-harvester", "name", "foo-harvester"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-harvester", "description", "Terraform cloudCredential acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-harvester", "driver", harvesterConfigDriver),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-harvester", "harvester_credential_config.0.kubeconfig_content", "YYYYYYYYYYYYYYYYYYYY"),
				),
			},
			{
				Config: testAccRancher2CloudCredentialConfigHarvester,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CloudCredentialExists(testAccRancher2CloudCredentialType+".foo-harvester", cloudCredential),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-harvester", "name", "foo-harvester"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-harvester", "description", "Terraform cloudCredential acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-harvester", "driver", harvesterConfigDriver),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-harvester", "harvester_credential_config.0.kubeconfig_content", "XXXXXXXXXXXXXXXXXXXX"),
				),
			},
		},
	})
}

func TestAccRancher2CloudCredential_disappears_Harvester(t *testing.T) {
	var cloudCredential *CloudCredential

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2CloudCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2CloudCredentialConfigHarvester,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CloudCredentialExists(testAccRancher2CloudCredentialType+".foo-harvester", cloudCredential),
					testAccRancher2CloudCredentialDisappears(cloudCredential),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRancher2CloudCredential_basic_Openstack(t *testing.T) {
	var cloudCredential *CloudCredential

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2CloudCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2CloudCredentialConfigOpenstack,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CloudCredentialExists(testAccRancher2CloudCredentialType+".foo-openstack", cloudCredential),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-openstack", "name", "foo-openstack"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-openstack", "description", "Terraform cloudCredential acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-openstack", "driver", openstackConfigDriver),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-openstack", "openstack_credential_config.0.password", "XXXXXXXXXXXXXXXXXXXX"),
				),
			},
			{
				Config: testAccRancher2CloudCredentialUpdateConfigOpenstack,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CloudCredentialExists(testAccRancher2CloudCredentialType+".foo-openstack", cloudCredential),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-openstack", "name", "foo-openstack"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-openstack", "description", "Terraform cloudCredential acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-openstack", "driver", openstackConfigDriver),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-openstack", "openstack_credential_config.0.password", "YYYYYYYYYYYYYYYYYYYY"),
				),
			},
			{
				Config: testAccRancher2CloudCredentialConfigOpenstack,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CloudCredentialExists(testAccRancher2CloudCredentialType+".foo-openstack", cloudCredential),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-openstack", "name", "foo-openstack"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-openstack", "description", "Terraform cloudCredential acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-openstack", "driver", openstackConfigDriver),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-openstack", "openstack_credential_config.0.password", "XXXXXXXXXXXXXXXXXXXX"),
				),
			},
		},
	})
}

func TestAccRancher2CloudCredential_disappears_Openstack(t *testing.T) {
	var cloudCredential *CloudCredential

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2CloudCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2CloudCredentialConfigOpenstack,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CloudCredentialExists(testAccRancher2CloudCredentialType+".foo-openstack", cloudCredential),
					testAccRancher2CloudCredentialDisappears(cloudCredential),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRancher2CloudCredential_basic_S3(t *testing.T) {
	var cloudCredential *CloudCredential

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2CloudCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2CloudCredentialConfigS3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CloudCredentialExists(testAccRancher2CloudCredentialType+".foo-s3", cloudCredential),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-s3", "name", "foo-s3"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-s3", "description", "Terraform cloudCredential acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-s3", "driver", s3ConfigDriver),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-s3", "s3_credential_config.0.access_key", "XXXXXXXXXXXXXXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-s3", "s3_credential_config.0.default_bucket", "default_bucket"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-s3", "s3_credential_config.0.default_folder", "default_folder"),
				),
			},
			{
				Config: testAccRancher2CloudCredentialUpdateConfigS3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CloudCredentialExists(testAccRancher2CloudCredentialType+".foo-s3", cloudCredential),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-s3", "name", "foo-s3"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-s3", "description", "Terraform cloudCredential acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-s3", "driver", s3ConfigDriver),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-s3", "s3_credential_config.0.access_key", "YYYYYYYYYYYYYYYYYYYY"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-s3", "s3_credential_config.0.default_bucket", "default_bucket_updated"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-s3", "s3_credential_config.0.default_folder", "default_folder_updated"),
				),
			},
			{
				Config: testAccRancher2CloudCredentialConfigS3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CloudCredentialExists(testAccRancher2CloudCredentialType+".foo-s3", cloudCredential),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-s3", "name", "foo-s3"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-s3", "description", "Terraform cloudCredential acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-s3", "driver", s3ConfigDriver),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-s3", "s3_credential_config.0.access_key", "XXXXXXXXXXXXXXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-s3", "s3_credential_config.0.default_bucket", "default_bucket"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-s3", "s3_credential_config.0.default_folder", "default_folder"),
				),
			},
		},
	})
}

func TestAccRancher2CloudCredential_disappears_S3(t *testing.T) {
	var cloudCredential *CloudCredential

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2CloudCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2CloudCredentialConfigS3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CloudCredentialExists(testAccRancher2CloudCredentialType+".foo-s3", cloudCredential),
					testAccRancher2CloudCredentialDisappears(cloudCredential),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRancher2CloudCredential_basic_Vsphere(t *testing.T) {
	var cloudCredential *CloudCredential

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2CloudCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2CloudCredentialConfigVsphere,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CloudCredentialExists(testAccRancher2CloudCredentialType+".foo-vsphere", cloudCredential),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-vsphere", "name", "foo-vsphere"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-vsphere", "description", "Terraform cloudCredential acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-vsphere", "driver", vmwarevsphereConfigDriver),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-vsphere", "vsphere_credential_config.0.password", "XXXXXXXXXXXXXXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-vsphere", "vsphere_credential_config.0.vcenter", "vcenter"),
				),
			},
			{
				Config: testAccRancher2CloudCredentialUpdateConfigVsphere,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CloudCredentialExists(testAccRancher2CloudCredentialType+".foo-vsphere", cloudCredential),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-vsphere", "name", "foo-vsphere"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-vsphere", "description", "Terraform cloudCredential acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-vsphere", "driver", vmwarevsphereConfigDriver),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-vsphere", "vsphere_credential_config.0.password", "YYYYYYYYYYYYYYYYYYYY"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-vsphere", "vsphere_credential_config.0.vcenter", "vcenter2"),
				),
			},
			{
				Config: testAccRancher2CloudCredentialConfigVsphere,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CloudCredentialExists(testAccRancher2CloudCredentialType+".foo-vsphere", cloudCredential),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-vsphere", "name", "foo-vsphere"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-vsphere", "description", "Terraform cloudCredential acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-vsphere", "driver", vmwarevsphereConfigDriver),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-vsphere", "vsphere_credential_config.0.password", "XXXXXXXXXXXXXXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2CloudCredentialType+".foo-vsphere", "vsphere_credential_config.0.vcenter", "vcenter"),
				),
			},
		},
	})
}

func TestAccRancher2CloudCredential_disappears_Vsphere(t *testing.T) {
	var cloudCredential *CloudCredential

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2CloudCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancher2CloudCredentialConfigVsphere,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2CloudCredentialExists(testAccRancher2CloudCredentialType+".foo-vsphere", cloudCredential),
					testAccRancher2CloudCredentialDisappears(cloudCredential),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2CloudCredentialDisappears(cloudCredential *CloudCredential) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2CloudCredentialType {
				continue
			}
			client, err := testAccProvider.Meta().(*Config).ManagementClient()
			if err != nil {
				return err
			}

			cloudCredential := &norman.Resource{}
			err = client.APIBaseClient.ByID(managementClient.CloudCredentialType, rs.Primary.ID, cloudCredential)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = client.APIBaseClient.Delete(cloudCredential)
			if err != nil {
				return fmt.Errorf("Error removing Cloud Credential: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Pending:    []string{"active"},
				Target:     []string{"removed"},
				Refresh:    cloudCredentialStateRefreshFunc(client, cloudCredential.ID),
				Timeout:    10 * time.Minute,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf(
					"[ERROR] waiting for Cloud Credential (%s) to be removed: %s", cloudCredential.ID, waitErr)
			}
		}
		return nil

	}
}

func testAccCheckRancher2CloudCredentialExists(n string, cloudCredential *CloudCredential) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cloud Credential ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		foundReg := &CloudCredential{}
		err = client.APIBaseClient.ByID(managementClient.CloudCredentialType, rs.Primary.ID, foundReg)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("Cloud Credential not found")
			}
			return err
		}

		cloudCredential = foundReg

		return nil
	}
}

func testAccCheckRancher2CloudCredentialDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2CloudCredentialType {
			continue
		}
		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		cloudCredential := &CloudCredential{}
		err = client.APIBaseClient.ByID(managementClient.CloudCredentialType, rs.Primary.ID, cloudCredential)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("Cloud Credential still exists")
	}
	return nil
}
