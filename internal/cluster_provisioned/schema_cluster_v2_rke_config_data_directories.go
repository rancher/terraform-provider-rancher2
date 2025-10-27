package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Types

func clusterV2RKEConfigDataDirectoriesFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"system_agent": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Desired System Agent data directory.",
		},
		"provisioning": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Desired provisioning data directory.",
		},
		"k8s_distro": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Desired k8s distro data directory.",
		},
	}

	return s
}
