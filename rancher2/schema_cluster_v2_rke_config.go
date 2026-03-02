package rancher2

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// Types

var (
	// StackPreference is the networking stack used by the cluster
	StackPreference = []string{"dual", "ipv6", "ipv4"}
)

func clusterV2RKEConfigFieldsV0() map[string]*schema.Schema {
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
		"machine_pool_defaults": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "Default values for machine pool configurations if unset",
			Elem: &schema.Resource{
				Schema: clusterV2RKEConfigMachinePoolDefaultFields(),
			},
		},
		"machine_selector_config": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "Cluster V2 machine selector config",
			Elem: &schema.Resource{
				Schema: clusterV2RKEConfigSystemConfigFieldsV0(),
			},
		},
		"machine_selector_files": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "Cluster V2 machine selector files",
			Elem: &schema.Resource{
				Schema: clusterV2RKEConfigMachineSelectorFilesFields(),
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
		"networking": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Description: "Cluster V2 networking",
			Elem: &schema.Resource{
				Schema: clusterV2Networking(),
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

func clusterV2Networking() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"stack_preference": {
			Type:     schema.TypeString,
			Optional: true,
			Description: "Specify the networking stack used by the cluster. The selected value configures the address " +
				"used for health and readiness probes of calico, etcd, kube-apiserver, kube-scheduler, kube-controller-manager, and kubelet. " +
				"It also defines the server URL in the authentication-token-webhook-config-file for " +
				"the Authorized Cluster Endpoint and the advertise-client-urls for etcd during snapshot restore. " +
				"When set to dual, the cluster uses localhost; " +
				"when set to ipv6, it uses [::1]; " +
				"when set to ipv4, it uses 127.0.0.1",
			ValidateFunc: validation.StringInSlice(StackPreference, false),
		},
	}

	return s
}

func clusterV2RKEConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"additional_manifest": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Cluster V2 additional manifest",
		},
		"data_directories": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Cluster V2 data directories",
			Elem: &schema.Resource{
				Schema: clusterV2RKEConfigDataDirectoriesFields(),
			},
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
		"machine_pool_defaults": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "Default values for machine pool configurations if unset",
			Elem: &schema.Resource{
				Schema: clusterV2RKEConfigMachinePoolDefaultFields(),
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
		"machine_selector_files": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "Cluster V2 machine selector files",
			Elem: &schema.Resource{
				Schema: clusterV2RKEConfigMachineSelectorFilesFields(),
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
		"networking": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Description: "Cluster V2 networking",
			Elem: &schema.Resource{
				Schema: clusterV2Networking(),
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
