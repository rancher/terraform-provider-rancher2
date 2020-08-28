package rancher2

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

const testAccRancher2PodSecurityPolicyTemplateType = "rancher2_pod_security_policy_template"

var (
	testAccCheckRancher2PodSecurityPolicyTemplate = `
resource "` + testAccRancher2PodSecurityPolicyTemplateType + `" "foo" {
  name = "foo"
  description = "Terraform PodSecurityPolicyTemplate acceptance test"
  allow_privilege_escalation = false
  allowed_csi_driver {
    name = "something"
  }
  allowed_csi_driver {
    name = "something-else"
  }
  allowed_flex_volume {
    driver = "something"
  }
  allowed_flex_volume {
    driver = "something-else"
  }
  allowed_host_path {
    path_prefix = "/"
    read_only = true
  }
  allowed_host_path {
    path_prefix = "//"
    read_only = false
  }
  allowed_proc_mount_types = ["Default"]
  default_allow_privilege_escalation = false
  fs_group {
    rule = "MustRunAs"
    range {
      min = 0
      max = 100
    }
    range {
      min = 0
      max = 100
    }
  }
  host_ipc = false
  host_network = false
  host_pid = false
  host_port {
    min = 0
    max = 65535
  }
  host_port {
    min = 1024
    max = 8080
  }
  privileged = false
  read_only_root_filesystem = false
  required_drop_capabilities = ["something"]

  run_as_user {
    rule = "MustRunAs"
    range {
      min = 1
      max = 100
    }
    range {
      min = 2
      max = 1024
    }
  }
  run_as_group {
    rule = "MustRunAs"
    range {
      min = 1
      max = 100
    }
    range {
      min = 2
      max = 1024
    }
  }
  runtime_class {
    default_runtime_class_name = "something"
    allowed_runtime_class_names  = ["something"]
  }
  se_linux {
    rule = "RunAsAny"
  }
  supplemental_group {
    rule = "RunAsAny"
  }
  volumes = ["azureFile"]
}
`
	testAccCheckRancher2PodSecurityPolicyTemplateUpdate = `
resource "` + testAccRancher2PodSecurityPolicyTemplateType + `" "foo" {
  name = "foo"
  description = "Terraform PodSecurityPolicyTemplate acceptance test - updated"
  allow_privilege_escalation = false
  allowed_csi_driver {
    name = "something"
  }
  allowed_csi_driver {
    name = "something-else"
  }
  allowed_flex_volume {
    driver = "something"
  }
  allowed_flex_volume {
    driver = "something-else"
  }
  allowed_host_path {
    path_prefix = "/"
    read_only = true
  }
  allowed_host_path {
    path_prefix = "//"
    read_only = false
  }
  allowed_proc_mount_types = ["Default"]
  default_allow_privilege_escalation = false
  fs_group {
    rule = "MustRunAs"
    range {
      min = 0
      max = 100
    }
    range {
      min = 0
      max = 100
    }
  }
  host_ipc = false
  host_network = false
  host_pid = false
  host_port {
    min = 0
    max = 65535
  }
  host_port {
    min = 1024
    max = 8080
  }
  privileged = false
  read_only_root_filesystem = false
  required_drop_capabilities = ["something"]

  run_as_user {
    rule = "MustRunAs"
    range {
      min = 1
      max = 100
    }
    range {
      min = 2
      max = 1024
    }
  }
  run_as_group {
    rule = "MustRunAs"
    range {
      min = 1
      max = 100
    }
    range {
      min = 2
      max = 1024
    }
  }
  runtime_class {
    default_runtime_class_name = "something"
    allowed_runtime_class_names  = ["something"]
  }
  se_linux {
    rule = "RunAsAny"
  }
  supplemental_group {
    rule = "RunAsAny"
  }
  volumes = ["azureFile"]
}
`
)

func init() {}

func TestAccRancher2PodSecurityPolicyTemplate_Basic(t *testing.T) {
	var pspTemplate *managementClient.PodSecurityPolicyTemplate

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2PodSecurityPolicyTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancher2PodSecurityPolicyTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NPodSecurityPolicyTemplateExists(testAccRancher2PodSecurityPolicyTemplateType+".foo", pspTemplate),
					resource.TestCheckResourceAttr(testAccRancher2PodSecurityPolicyTemplateType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2PodSecurityPolicyTemplateType+".foo", "description", "Terraform PodSecurityPolicyTemplate acceptance test"),
				),
			},
			{
				Config: testAccCheckRancher2PodSecurityPolicyTemplateUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NPodSecurityPolicyTemplateExists(testAccRancher2PodSecurityPolicyTemplateType+".foo", pspTemplate),
					resource.TestCheckResourceAttr(testAccRancher2PodSecurityPolicyTemplateType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2PodSecurityPolicyTemplateType+".foo", "description", "Terraform PodSecurityPolicyTemplate acceptance test - updated"),
				),
			},
		},
	})
}

func testAccCheckRancher2NPodSecurityPolicyTemplateExists(n string, pspTemplate *managementClient.PodSecurityPolicyTemplate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No PodSecurityPolicyTemplate ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		foundPSP, err := client.PodSecurityPolicyTemplate.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("PodSecurityPolicyTemplate not found")
			}
			return err
		}

		pspTemplate = foundPSP

		return nil
	}
}

func testAccCheckRancher2PodSecurityPolicyTemplateDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rancher2_pod_security_policy_template" {
			continue
		}
		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		_, err = client.PodSecurityPolicyTemplate.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}

		return fmt.Errorf("PodSecurityPolicyTemplate still exists")
	}
	return nil
}
