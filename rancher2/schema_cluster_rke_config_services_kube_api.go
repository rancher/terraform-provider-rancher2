package rancher2

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	//"github.com/rancher/types/apis/management.cattle.io/v3"
	auditv1 "k8s.io/apiserver/pkg/apis/audit/v1"
)

const (
	clusterRKEConfigServicesKubeAPIApiversionTag                                = "apiVersion"
	clusterRKEConfigServicesKubeAPIKindTag                                      = "kind"
	clusterRKEConfigServicesKubeAPIAuditLogConfigPolicyAPIDefault               = "audit.k8s.io/v1"
	clusterRKEConfigServicesKubeAPIEventRateLimitConfigAPIDefault               = "eventratelimit.admission.k8s.io/v1alpha1"
	clusterRKEConfigServicesKubeAPIEncryptionConfigAPIDefault                   = "apiserver.config.k8s.io/v1"
	clusterRKEConfigServicesKubeAPIAdmissionConfigurationAPIDefault             = "apiserver.config.k8s.io/v1"
	clusterRKEConfigServicesKubeAPIAdmissionConfigurationPodSecurityAPIDefault  = "pod-security.admission.config.k8s.io/v1"
	clusterRKEConfigServicesKubeAPIAuditLogConfigPolicyKindDefault              = "Policy"
	clusterRKEConfigServicesKubeAPIEventRateLimitConfigKindDefault              = "Configuration"
	clusterRKEConfigServicesKubeAPIEncryptionConfigKindDefault                  = "EncryptionConfiguration"
	clusterRKEConfigServicesKubeAPIAdmissionConfigurationKindDefault            = "AdmissionConfiguration"
	clusterRKEConfigServicesKubeAPIAdmissionConfigurationPodSecurityKindDefault = "PodSecurityConfiguration"
)

var (
	clusterRKEConfigServicesKubeAPIRequired = []string{
		clusterRKEConfigServicesKubeAPIApiversionTag,
		clusterRKEConfigServicesKubeAPIKindTag,
	}
)

//Schemas

func clusterRKEConfigServicesKubeAPIAdmissionConfigurationFieldsV0() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"api_version": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     clusterRKEConfigServicesKubeAPIAdmissionConfigurationAPIDefault,
			Description: "Admission configuration ApiVersion",
		},
		"kind": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     clusterRKEConfigServicesKubeAPIAdmissionConfigurationKindDefault,
			Description: "Admission configuration Kind",
		},
		"plugins": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigServicesKubeAPIAdmissionConfigurationPluginsFieldsV0(),
			},
			Description: "Admission configuration plugins",
		},
	}
	return s
}

func clusterRKEConfigServicesKubeAPIAdmissionConfigurationFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"api_version": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     clusterRKEConfigServicesKubeAPIAdmissionConfigurationAPIDefault,
			Description: "Admission configuration ApiVersion",
		},
		"kind": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     clusterRKEConfigServicesKubeAPIAdmissionConfigurationKindDefault,
			Description: "Admission configuration Kind",
		},
		"plugins": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigServicesKubeAPIAdmissionConfigurationPluginsFields(),
			},
			Description: "Admission configuration plugins",
		},
	}
	return s
}

func clusterRKEConfigServicesKubeAPIAdmissionConfigurationPluginsFieldsV0() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Plugin name",
		},
		"path": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "Plugin path",
		},
		"configuration": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Plugin configuration",
		},
	}
	return s
}
func clusterRKEConfigServicesKubeAPIAdmissionConfigurationPluginsFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Plugin name",
		},
		"path": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "Plugin path",
		},
		"configuration": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
				v, ok := val.(string)
				if !ok || len(v) == 0 {
					return
				}
				m, err := ghodssyamlToMapInterface(v)
				if err != nil {
					errs = append(errs, fmt.Errorf("%q must be in yaml format, error: %v", key, err))
					return
				}
				for _, k := range clusterRKEConfigServicesKubeAPIRequired {
					check, ok := m[k].(string)
					if !ok || len(check) == 0 {
						errs = append(errs, fmt.Errorf("%s is required on yaml", k))
					}
					if k == clusterRKEConfigServicesKubeAPIKindTag {
						if check != clusterRKEConfigServicesKubeAPIAdmissionConfigurationPodSecurityKindDefault {
							errs = append(errs, fmt.Errorf("%s value %s should be: %s", k, check, clusterRKEConfigServicesKubeAPIAdmissionConfigurationPodSecurityKindDefault))
						}
					}

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
			Description: "Plugin configuration",
		},
	}
	return s
}

func clusterRKEConfigServicesKubeAPIAuditLogConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"format": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "json",
		},
		"max_age": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  30,
		},
		"max_backup": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  10,
		},
		"max_size": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  100,
		},
		"path": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "/var/log/kube-audit/audit-log.json",
		},
		"policy": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
				v, ok := val.(string)
				if !ok || len(v) == 0 {
					return
				}
				m, err := ghodssyamlToMapInterface(v)
				if err != nil {
					errs = append(errs, fmt.Errorf("%q must be in yaml format, error: %v", key, err))
					return
				}
				for _, k := range clusterRKEConfigServicesKubeAPIRequired {
					check, ok := m[k].(string)
					if !ok || len(check) == 0 {
						errs = append(errs, fmt.Errorf("%s is required on yaml", k))
					}
					if k == clusterRKEConfigServicesKubeAPIKindTag {
						if check != clusterRKEConfigServicesKubeAPIAuditLogConfigPolicyKindDefault {
							errs = append(errs, fmt.Errorf("%s value %s should be: %s", k, check, clusterRKEConfigServicesKubeAPIAuditLogConfigPolicyKindDefault))
						}
					}

				}
				return
			},
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				if old == "" || new == "" {
					return false
				}
				oldPolicy := &auditv1.Policy{}
				newPolicy := &auditv1.Policy{}
				oldMap, _ := ghodssyamlToMapInterface(old)
				newMap, _ := ghodssyamlToMapInterface(new)
				oldStr, _ := mapInterfaceToJSON(oldMap)
				newStr, _ := mapInterfaceToJSON(newMap)
				jsonToInterface(oldStr, oldPolicy)
				jsonToInterface(newStr, newPolicy)
				return reflect.DeepEqual(oldPolicy, newPolicy)
			},
		},
	}
	return s
}

func clusterRKEConfigServicesKubeAPIAuditLogFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"configuration": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigServicesKubeAPIAuditLogConfigFields(),
			},
		},
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
	return s
}

func clusterRKEConfigServicesKubeAPIAdmissionConfigPluginsFieldsV0() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Plugin name",
		},
		"path": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "Plugin path",
		},
		"configuration": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Plugin configuration",
		},
	}
	return s
}
func clusterRKEConfigServicesKubeAPIAdmissionConfigPluginsFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Plugin name",
		},
		"path": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "Plugin path",
		},
		"configuration": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
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
			Description: "Plugin configuration",
		},
	}
	return s
}

func clusterRKEConfigServicesKubeAPIEventRateLimitFieldsV0() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"configuration": {
			Type:     schema.TypeMap,
			Optional: true,
		},
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
	return s
}

func clusterRKEConfigServicesKubeAPIEventRateLimitFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"configuration": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
				v, ok := val.(string)
				if !ok || len(v) == 0 {
					return
				}
				m, err := ghodssyamlToMapInterface(v)
				if err != nil {
					errs = append(errs, fmt.Errorf("%q must be in yaml format, error: %v", key, err))
					return
				}
				for _, k := range clusterRKEConfigServicesKubeAPIRequired {
					check, ok := m[k].(string)
					if !ok || len(check) == 0 {
						errs = append(errs, fmt.Errorf("%s is required on yaml", k))
					}
					if k == clusterRKEConfigServicesKubeAPIKindTag {
						if check != clusterRKEConfigServicesKubeAPIEventRateLimitConfigKindDefault {
							errs = append(errs, fmt.Errorf("%s value %s should be: %s", k, check, clusterRKEConfigServicesKubeAPIEventRateLimitConfigKindDefault))
						}
					}

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
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
	return s
}

// Used by datasource
func clusterRKEConfigServicesKubeAPIEventRateLimitFieldsData() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"configuration": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
	return s
}

func clusterRKEConfigServicesKubeAPISecretsEncryptionConfigFieldsV0() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"custom_config": {
			Type:     schema.TypeMap,
			Optional: true,
		},
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
	return s
}

func clusterRKEConfigServicesKubeAPISecretsEncryptionConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"custom_config": {
			Type:     schema.TypeString,
			Optional: true,
			ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
				v, ok := val.(string)
				if !ok || len(v) == 0 {
					return
				}
				m, err := ghodssyamlToMapInterface(v)
				if err != nil {
					errs = append(errs, fmt.Errorf("%q must be in yaml format, error: %v", key, err))
					return
				}
				for _, k := range clusterRKEConfigServicesKubeAPIRequired {
					check, ok := m[k].(string)
					if !ok || len(check) == 0 {
						errs = append(errs, fmt.Errorf("%s is required on yaml", k))
					}
					if k == clusterRKEConfigServicesKubeAPIKindTag {
						if check != clusterRKEConfigServicesKubeAPIEncryptionConfigKindDefault {
							errs = append(errs, fmt.Errorf("%s value %s should be: %s", k, check, clusterRKEConfigServicesKubeAPIEncryptionConfigKindDefault))
						}
					}

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
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
	return s
}

// Used by datasource
func clusterRKEConfigServicesKubeAPISecretsEncryptionConfigFieldsData() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"custom_config": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
	return s
}

func clusterRKEConfigServicesKubeAPIFieldsV0() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"admission_configuration": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigServicesKubeAPIAdmissionConfigurationFieldsV0(),
			},
			Description: "Cluster admission configuration",
		},
		"always_pull_images": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"audit_log": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigServicesKubeAPIAuditLogFields(),
			},
		},
		"event_rate_limit": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigServicesKubeAPIEventRateLimitFieldsV0(),
			},
		},
		"extra_args": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"extra_binds": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"extra_env": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"image": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"pod_security_policy": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"secrets_encryption_config": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigServicesKubeAPISecretsEncryptionConfigFieldsV0(),
			},
		},
		"service_cluster_ip_range": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"service_node_port_range": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
	return s
}

func clusterRKEConfigServicesKubeAPIFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"admission_configuration": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigServicesKubeAPIAdmissionConfigurationFields(),
			},
			Description: "Cluster admission configuration",
		},
		"always_pull_images": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"audit_log": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigServicesKubeAPIAuditLogFields(),
			},
		},
		"event_rate_limit": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigServicesKubeAPIEventRateLimitFields(),
			},
		},
		"extra_args": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"extra_binds": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"extra_env": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"image": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"pod_security_policy": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"secrets_encryption_config": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigServicesKubeAPISecretsEncryptionConfigFields(),
			},
		},
		"service_cluster_ip_range": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"service_node_port_range": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
	return s
}

// Used by datasource
func clusterRKEConfigServicesKubeAPIFieldsData() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"admission_configuration": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigServicesKubeAPIAdmissionConfigurationFields(),
			},
			Description: "Cluster admission configuration",
		},
		"always_pull_images": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"audit_log": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigServicesKubeAPIAuditLogFields(),
			},
		},
		"event_rate_limit": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigServicesKubeAPIEventRateLimitFieldsData(),
			},
		},
		"extra_args": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"extra_binds": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"extra_env": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"image": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"pod_security_policy": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"secrets_encryption_config": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigServicesKubeAPISecretsEncryptionConfigFieldsData(),
			},
		},
		"service_cluster_ip_range": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"service_node_port_range": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
	return s
}
