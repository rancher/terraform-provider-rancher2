package rancher2

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	norman "github.com/rancher/norman/types"
	"github.com/rancher/rancher/pkg/apis/catalog.cattle.io/v1"
)

const (
	appV2Kind              = "App"
	appV2APIGroup          = "catalog.cattle.io"
	appV2APIVersion        = "v1"
	appV2APIType           = rancher2CatalogTypePrefix + ".app"
	appV2OperationAPIType  = rancher2CatalogTypePrefix + ".operation"
	appV2ValueGlobal       = "global."
	appV2ClusterIDsep      = "."
	appV2DefaultRegistryID = "system-default-registry"
)

//Types

type AppV2 struct {
	norman.Resource
	v1.App
}

type AppV2Operation struct {
	norman.Resource
	v1.Operation
}

// Schemas

func appV2Fields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cluster_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "App v2 name",
		},
		"namespace": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "App v2 namespace",
		},
		"repo_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Repo name",
		},
		"chart_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Chart name",
		},
		"chart_version": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Chart version",
		},
		"cluster_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"project_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Deploy app within project ID",
		},
		"system_default_registry": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"values": {
			Type:             schema.TypeString,
			Optional:         true,
			Sensitive:        false,
			Description:      "App v2 custom values yaml",
			ValidateFunc:     validateAppSchema,
			DiffSuppressFunc: suppressAppDiff,
		},
		"deployment_values": {
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: false,
			Description: "Values YAML file including computed values. This field prevents incorrect discrepancies from " +
				"showing in the terraform plan output when files change but values stay the same, due to additional " +
				"computed values included by the provider itself.",
		},
		"cleanup_on_fail": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Cleanup app V2 on failed chart upgrade",
		},
		"disable_hooks": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Disable app V2 chart hooks",
		},
		"disable_open_api_validation": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Disable app V2 Open API Validation",
		},
		"force_upgrade": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Force app V2 chart upgrade",
		},
		"wait": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Wait until app is deployed",
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}

func validateAppSchema(val interface{}, key string) (warns []string, errs []error) {
	v, ok := val.(string)
	if !ok || len(v) == 0 {
		return
	}
	_, err := ghodssyamlToMapInterface(v)
	if err != nil {
		errs = append(errs, fmt.Errorf("[ERROR] %q must be in YAML format, error: %v", key, err))
		return
	}
	return
}

func suppressAppDiff(_, old, new string, d *schema.ResourceData) bool {
	oldMap, _ := ghodssyamlToMapInterface(old)
	newMap, _ := ghodssyamlToMapInterface(new)

	if newMap == nil {
		newMap = map[string]interface{}{}
	}

	return reflect.DeepEqual(oldMap, newMap)
}
