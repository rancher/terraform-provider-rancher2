package rancher2

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

const testAccRancher2PodSecurityAdmissionConfigurationTemplateType = "rancher2_pod_security_admission_configuration_template"

var (
	testAccCheckRancher2PodSecurityAdmissionConfigurationTemplate = `
resource "` + testAccRancher2PodSecurityAdmissionConfigurationTemplateType + `" "foo" {
  name = "foo"
  description = "Terraform PodSecurityAdmissionConfigurationTemplate acceptance test"
  defaults {
    audit = "privileged"
    audit_version = "latest"
    enforce = "privileged"
    enforce_version = "latest"
    warn = "privileged"
    warn_version = "latest"
  }
}
`
	testAccCheckRancher2PodSecurityAdmissionConfigurationTemplateUpdate = `
resource "` + testAccRancher2PodSecurityAdmissionConfigurationTemplateType + `" "foo" {
  name = "foo"
  description = "Terraform PodSecurityAdmissionConfigurationTemplate acceptance test - updated"
  defaults {
    audit = "restricted"
    audit_version = "latest"
    enforce = "restricted"
    enforce_version = "latest"
    warn = "restricted"
    warn_version = "latest"
  }
  exemptions {
    usernames = ["testuser"]
    runtime_classes = ["testclass"]
    namespaces = ["ingress-nginx","kube-system"]
  }
}
`
)

func init() {}

func TestAccRancher2PodSecurityAdmissionConfigurationTemplate_Basic(t *testing.T) {
	var psacTemplate *managementClient.PodSecurityAdmissionConfigurationTemplate

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2PodSecurityAdmissionConfigurationTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2PodSecurityAdmissionConfigurationTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2PodSecurityAdmissionConfigurationTemplateExists(testAccRancher2PodSecurityAdmissionConfigurationTemplateType+".foo", psacTemplate),
					resource.TestCheckResourceAttr(testAccRancher2PodSecurityAdmissionConfigurationTemplateType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2PodSecurityAdmissionConfigurationTemplateType+".foo", "description", "Terraform PodSecurityAdmissionConfigurationTemplate acceptance test"),
				),
			},
			{
				Config: testAccCheckRancher2PodSecurityAdmissionConfigurationTemplateUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2PodSecurityAdmissionConfigurationTemplateExists(testAccRancher2PodSecurityAdmissionConfigurationTemplateType+".foo", psacTemplate),
					resource.TestCheckResourceAttr(testAccRancher2PodSecurityAdmissionConfigurationTemplateType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2PodSecurityAdmissionConfigurationTemplateType+".foo", "description", "Terraform PodSecurityAdmissionConfigurationTemplate acceptance test - updated"),
				),
			},
		},
	})
}

func testAccCheckRancher2PodSecurityAdmissionConfigurationTemplateExists(n string, psacTemplate *managementClient.PodSecurityAdmissionConfigurationTemplate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no PodSecurityAdmissionConfigurationTemplate ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		foundPSAC, err := client.PodSecurityAdmissionConfigurationTemplate.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("PodSecurityAdmissionConfigurationTemplate not found")
			}
			return err
		}

		psacTemplate = foundPSAC

		return nil
	}
}

func testAccCheckRancher2PodSecurityAdmissionConfigurationTemplateDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rancher2_pod_security_admission_configuration_template" {
			continue
		}
		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		_, err = client.PodSecurityAdmissionConfigurationTemplate.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}

		return fmt.Errorf("PodSecurityAdmissionConfigurationTemplate still exists")
	}
	return nil
}
