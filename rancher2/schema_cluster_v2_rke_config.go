package rancher2

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Types

func clusterV2RKEConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"additional_manifest": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Cluster V2 additional manifest",
		},
		"local_auth_endpoint": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Deprecated:  "Use rancher2_cluster_v2.local_auth_endpoint instead",
			Description: "Cluster V2 local auth endpoint",
			Elem: &schema.Resource{
				Schema: clusterV2LocalAuthEndpointFields(),
			},
		},
		"upgrade_strategy": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Cluster V2 upgrade strategy",
			Elem: &schema.Resource{
				Schema: clusterV2RKEConfigUpgradeStrategyFields(),
			},
		},
		"chart_values": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Cluster V2 chart values. It should be in YAML format",
			ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
				v, ok := val.(string)
				if !ok || len(v) == 0 {
					return
				}
				_, err := ghodssyamlToMapInterface(v)
				if err != nil {
					errs = append(errs, fmt.Errorf("%q must be in yaml format, error: %v", key, err))
					return
				}
				return
			},
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				if old == "" || new == "" {
					return false
				}
				oldMap, _ := ghodssyamlToMapInterface(old)
				newMap, _ := ghodssyamlToMapInterface(new)
				return reflect.DeepEqual(oldMap, newMap)
			},
		},
		"machine_global_config": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Cluster V2 machine global config",
			ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
				v, ok := val.(string)
				if !ok || len(v) == 0 {
					return
				}
				_, err := ghodssyamlToMapInterface(v)
				if err != nil {
					errs = append(errs, fmt.Errorf("%q must be in yaml format, error: %v", key, err))
					return
				}
				return
			},
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				if old == "" || new == "" {
					return false
				}
				oldMap, _ := ghodssyamlToMapInterface(old)
				newMap, _ := ghodssyamlToMapInterface(new)
				return reflect.DeepEqual(oldMap, newMap)
			},
		},
		"machine_pools": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "Cluster V2 machine pools",
			Elem: &schema.Resource{
				Schema: clusterV2RKEConfigMachinePoolFields(),
			},
		},
		"machine_selector_config": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "Cluster V2 machine selector config",
			Elem: &schema.Resource{
				Schema: clusterV2RKEConfigSystemConfigFields(),
			},
		},
		"registries": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Cluster V2 registries",
			Elem: &schema.Resource{
				Schema: clusterV2RKEConfigRegistryFields(),
			},
		},
		"etcd": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Description: "Cluster V2 etcd",
			Elem: &schema.Resource{
				Schema: clusterV2RKEConfigETCDFields(),
			},
		},
		"rotate_certificates": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Cluster V2 certificate rotation",
			Elem: &schema.Resource{
				Schema: clusterV2RKEConfigRotateCertificatesFields(),
			},
		},
		"etcd_snapshot_create": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Cluster V2 etcd snapshot create",
			Elem: &schema.Resource{
				Schema: clusterV2RKEConfigETCDSnapshotCreateFields(),
			},
		},
		"etcd_snapshot_restore": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Cluster V2 etcd snapshot restore",
			Elem: &schema.Resource{
				Schema: clusterV2RKEConfigETCDSnapshotRestoreFields(),
			},
		},
	}

	return s
}
