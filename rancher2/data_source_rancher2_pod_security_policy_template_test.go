package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testAccCheckRancher2PodSecurityPolicyTemplateDataSourceConfig string

func TestAccRancher2PodSecurityPolicyTemplateDataSource(t *testing.T) {
	testAccCheckRancher2PodSecurityPolicyTemplateDataSourceConfig := testAccCheckRancher2PodSecurityPolicyTemplate + `
data "` + testAccRancher2PodSecurityPolicyTemplateType + `" "foo" {
  name = rancher2_pod_security_policy_template.foo.name
}
`
	name := "data." + testAccRancher2PodSecurityPolicyTemplateType + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2PodSecurityPolicyTemplateDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "description", "Terraform PodSecurityPolicyTemplate acceptance test"),
				),
			},
		},
	})
}
