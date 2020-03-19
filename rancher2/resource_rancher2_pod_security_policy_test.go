package rancher2

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	managementClient "github.com/rancher/types/client/management/v3"
	"testing"
)


var (
	testPspTemplate = `
resource "rancher2_pod_security_policy_template" "foo" {
  name = "foo"
  description = "Terraform PodSecurityPolicyTemplate acceptance test"
  annotations {
    "something.annotation.io/field" = "another_field"
  }
  labels {
    app = "something"
  }
  allow_privilege_escalation = true
  allowed_capabilities = ["something"]
  allowed_csi_driver {
    name = "something"
  }
  allowed_csi_driver {
    name = "something_else"
  }
  allowed_flex_volume {
    driver = "something"
  }
  allowed_flex_volume {
    driver = "something_else"
  }
  allowed_host_path {
    path_prefix = "/"
    read_only = true
  }
  allowed_host_path {
    path_prefix = "//"
    read_only = false
  }
  allowed_proc_mount_types = ["something"]
  allowed_unsafe_sysctls = ["something"]
  default_add_capabilities = ["something"]
  default_allow_privilege_escalation = false
  forbidden_sysctls = ["something"]
  fs_group {
    rule = "deny"
    range {
      min = 0
      max = 100
    }
    range {
      min = 0
      max = 100
    }
  }
  fs_group {
    rule = "allow"
    ange {
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
    rule = "something_else"
    range {
      min = 1
      max = 100
    }
    range {
      min = 2
      max = 1024
    }
  }
  run_as_user {
    rule = "something"
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
    rule = "something_else"
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
    rule = "something"
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
  runtime_class {
    default_runtime_class_name = "something else"
    allowed_runtime_class_names  = ["something"]
  }
  se_linux {
    rule = "ShouldRunAs"
    se_linux_option {
      user = "me"
      role = "something"
      type = "something"
      level = "something"
    }
    se_linux_option {
      user = "you"
      role = "something else"
      type = "something else"
      level = "something else"
    }
  }
  se_linux {
    rule = "MaybeRanAs"
    se_linux_option {
      user = "me"
      role = "something"
      type = "something"
      level = "something"
    }
    se_linux_option {
      user = "you"
      role = "something else"
      type = "something else"
      level = "something else"
    }
  }
  supplemental_group {
    rule = "blah"
    range {
      min = 0
      max = 0
    }
    range {
      min = 0
      max = 0
    }
  }
  supplemental_group {
    rule = "gah"
    range {
      min = 0
      max = 0
    }
    range {
      min = 0
      max = 0
    }
  }
  volumes = ["something"]
}
`

)

func init() {}

func TestPspTemplateCreate(t *testing.T) {
	var notifier *managementClient.Notifier

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NotifierDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2NotifierPagerdutyConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "description", "Terraform notifier acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "send_resolved", "true"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "pagerduty_config.0.service_key", "XXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "pagerduty_config.0.proxy_url", "http://proxy.test.io"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NotifierPagerdutyUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "description", "Terraform notifier acceptance test - updated"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "send_resolved", "false"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "pagerduty_config.0.proxy_url", "http://proxy2.test.io"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NotifierPagerdutyRecreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo", notifier),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "description", "Terraform notifier acceptance test"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "send_resolved", "true"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "pagerduty_config.0.service_key", "XXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2NotifierType+".foo", "pagerduty_config.0.proxy_url", "http://proxy.test.io"),
				),
			},
		},
	})
}

func TestAccRancher2Notifier_disappears_Pagerduty(t *testing.T) {
	var notifier *managementClient.Notifier

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NotifierDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2NotifierPagerdutyConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NotifierExists(testAccRancher2NotifierType+".foo", notifier),
					testAccRancher2NotifierDisappears(notifier),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}


func testAccCheckRancher2NotifierDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2NotifierType {
			continue
		}
		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		obj, err := client.Notifier.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}

		if obj.Removed != "" {
			return nil
		}

		return fmt.Errorf("Notifier still exists")
	}
	return nil
}
