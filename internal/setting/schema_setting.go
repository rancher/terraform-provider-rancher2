package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	norman "github.com/rancher/norman/types"
	managementAPI "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
)

const (
	settingV2Kind       = "Setting"
	settingV2APIGroup   = "management.cattle.io"
	settingV2APIVersion = "v3"
	settingV2APIType    = rancher2ManagementV2TypePrefix + ".setting"
)

//Types

type SettingV2 struct {
	norman.Resource
	managementAPI.Setting
}

//Schemas

func settingFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"value": {
			Type:     schema.TypeString,
			Required: true,
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}
