package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRancher2PodSecurityAdmissionConfigurationTemplateDataSource(t *testing.T) {
	testAccCheckRancher2PodSecurityAdmissionConfigurationTemplateDataSourceConfig := testAccCheckRancher2PodSecurityAdmissionConfigurationTemplate + `
data "` + testAccRancher2PodSecurityAdmissionConfigurationTemplateType + `" "foo" {
  name = rancher2_pod_security_admission_configuration_template.foo.name
}
`
	name := "data." + testAccRancher2PodSecurityAdmissionConfigurationTemplateType + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2PodSecurityAdmissionConfigurationTemplateDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "description", "Terraform PodSecurityAdmissionConfigurationTemplate acceptance test"),
				),
			},
		},
	})
}
