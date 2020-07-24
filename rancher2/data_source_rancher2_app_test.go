package rancher2

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccRancher2AppDataSource(t *testing.T) {
	testAccCheckRancher2AppDataSourceConfig := testAccRancher2AppConfig + `
data "` + testAccRancher2AppType + `" "foo" {
  name = rancher2_app.foo.name
  project_id = rancher2_app.foo.project_id
}
`
	name := "data." + testAccRancher2AppType + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2AppDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "description", "Terraform app acceptance test"),
					resource.TestCheckResourceAttr(name, "target_namespace", "testacc"),
					resource.TestCheckResourceAttr(name, "external_id", "catalog://?catalog=library&template=docker-registry&version=1.8.1"),
					resource.TestCheckResourceAttr(name, "answers.ingress_host", "test.xip.io"),
				),
			},
		},
	})
}

func testAccRancher2CheckClusterID() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return fmt.Errorf("testAccRancher2ClusterID %s", testAccRancher2ClusterID)

	}
}
