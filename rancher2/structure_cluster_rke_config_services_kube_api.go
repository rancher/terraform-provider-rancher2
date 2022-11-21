package rancher2

import (
	"fmt"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenClusterRKEConfigServicesKubeAPIAuditLogConfig(in *managementClient.AuditLogConfig) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	obj["format"] = in.Format
	obj["max_age"] = int(in.MaxAge)
	obj["max_backup"] = int(in.MaxBackup)
	obj["max_size"] = int(in.MaxSize)
	obj["path"] = in.Path

	if len(in.Policy) > 0 {
		policyStr, err := mapInterfaceToYAML(in.Policy)
		if err != nil {
			return nil, err
		}
		obj["policy"] = policyStr
	}

	return []interface{}{obj}, nil
}

func flattenClusterRKEConfigServicesKubeAPIAuditLog(in *managementClient.AuditLog) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	obj["enabled"] = in.Enabled

	config, err := flattenClusterRKEConfigServicesKubeAPIAuditLogConfig(in.Configuration)
	if err != nil {
		return []interface{}{}, err
	}
	obj["configuration"] = config

	return []interface{}{obj}, nil
}

func flattenClusterRKEConfigServicesKubeAPIEventRateLimit(in *managementClient.EventRateLimit) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	obj["enabled"] = in.Enabled

	if len(in.Configuration) > 0 {
		configStr, err := mapInterfaceToYAML(in.Configuration)
		if err != nil {
			return []interface{}{obj}
		}
		obj["configuration"] = configStr
	}

	return []interface{}{obj}
}

func flattenClusterRKEConfigServicesKubeAPISecretsEncryptionConfig(in *managementClient.SecretsEncryptionConfig) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	obj["enabled"] = in.Enabled

	if in.CustomConfig != nil {
		customConfigStr, err := interfaceToGhodssyaml(in.CustomConfig)
		if err != nil {
			return []interface{}{}, fmt.Errorf("Mashalling custom_config yaml: %v", err)
		}
		obj["custom_config"] = customConfigStr
	}

	return []interface{}{obj}, nil
}

func flattenClusterRKEConfigServicesKubeAPI(in *managementClient.KubeAPIService) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.AdmissionConfiguration) > 0 {
		obj["admission_configuration"] = in.AdmissionConfiguration
	}

	obj["always_pull_images"] = in.AlwaysPullImages

	if in.AuditLog != nil {
		auditLog, err := flattenClusterRKEConfigServicesKubeAPIAuditLog(in.AuditLog)
		if err != nil {
			return []interface{}{}, err
		}
		obj["audit_log"] = auditLog
	}

	if in.EventRateLimit != nil {
		obj["event_rate_limit"] = flattenClusterRKEConfigServicesKubeAPIEventRateLimit(in.EventRateLimit)
	}

	if len(in.ExtraArgs) > 0 {
		obj["extra_args"] = toMapInterface(in.ExtraArgs)
	}

	if len(in.ExtraBinds) > 0 {
		obj["extra_binds"] = toArrayInterface(in.ExtraBinds)
	}

	if len(in.ExtraEnv) > 0 {
		obj["extra_env"] = toArrayInterface(in.ExtraEnv)
	}

	if len(in.Image) > 0 {
		obj["image"] = in.Image
	}

	obj["pod_security_policy"] = in.PodSecurityPolicy

	if in.SecretsEncryptionConfig != nil {
		customConfig, err := flattenClusterRKEConfigServicesKubeAPISecretsEncryptionConfig(in.SecretsEncryptionConfig)
		if err != nil {
			return []interface{}{}, err
		}
		obj["secrets_encryption_config"] = customConfig
	}

	if len(in.ServiceClusterIPRange) > 0 {
		obj["service_cluster_ip_range"] = in.ServiceClusterIPRange
	}

	if len(in.ServiceNodePortRange) > 0 {
		obj["service_node_port_range"] = in.ServiceNodePortRange
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandClusterRKEConfigServicesKubeAPIAuditLogConfig(p []interface{}) (*managementClient.AuditLogConfig, error) {
	obj := &managementClient.AuditLogConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["format"].(string); ok && len(v) > 0 {
		obj.Format = v
	}

	if v, ok := in["max_age"].(int); ok && v > 0 {
		obj.MaxAge = int64(v)
	}

	if v, ok := in["max_backup"].(int); ok && v > 0 {
		obj.MaxBackup = int64(v)
	}

	if v, ok := in["max_size"].(int); ok && v > 0 {
		obj.MaxSize = int64(v)
	}

	if v, ok := in["path"].(string); ok && len(v) > 0 {
		obj.Path = v
	}

	if v, ok := in["policy"].(string); ok && len(v) > 0 {
		policyMap, err := ghodssyamlToMapInterface(v)
		if err != nil {
			return nil, err
		}

		obj.Policy = policyMap
	}

	return obj, nil
}

func expandClusterRKEConfigServicesKubeAPIAuditLog(p []interface{}) (*managementClient.AuditLog, error) {
	obj := &managementClient.AuditLog{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["enabled"].(bool); ok {
		obj.Enabled = v
	}

	if v, ok := in["configuration"].([]interface{}); ok && len(v) > 0 {
		config, err := expandClusterRKEConfigServicesKubeAPIAuditLogConfig(v)
		if err != nil {
			return nil, err
		}
		obj.Configuration = config
	}

	return obj, nil
}

func expandClusterRKEConfigServicesKubeAPIEventRateLimit(p []interface{}) *managementClient.EventRateLimit {
	obj := &managementClient.EventRateLimit{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["enabled"].(bool); ok {
		obj.Enabled = v
	}

	if v, ok := in["configuration"].(string); ok && len(v) > 0 {
		configMap, err := ghodssyamlToMapInterface(v)
		if err != nil {
			return obj
		}
		obj.Configuration = configMap
	}

	return obj
}

func expandClusterRKEConfigServicesKubeAPISecretsEncryptionConfig(p []interface{}) *managementClient.SecretsEncryptionConfig {
	obj := &managementClient.SecretsEncryptionConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["enabled"].(bool); ok {
		obj.Enabled = v
	}

	if v, ok := in["custom_config"].(string); ok && len(v) > 0 {
		customConfig := &managementClient.EncryptionConfiguration{}
		err := ghodssyamlToInterface(v, customConfig)
		if err != nil {
			return obj
		}
		obj.CustomConfig = customConfig
	}

	return obj
}

func expandClusterRKEConfigServicesKubeAPI(p []interface{}) (*managementClient.KubeAPIService, error) {
	obj := &managementClient.KubeAPIService{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["admission_configuration"].(map[string]interface{}); ok && len(v) > 0 {
		obj.AdmissionConfiguration = v
	}

	if v, ok := in["always_pull_images"].(bool); ok {
		obj.AlwaysPullImages = v
	}

	if v, ok := in["audit_log"].([]interface{}); ok && len(v) > 0 {
		auditLog, err := expandClusterRKEConfigServicesKubeAPIAuditLog(v)
		if err != nil {
			return nil, err
		}
		obj.AuditLog = auditLog
	}

	if v, ok := in["event_rate_limit"].([]interface{}); ok && len(v) > 0 {
		obj.EventRateLimit = expandClusterRKEConfigServicesKubeAPIEventRateLimit(v)
	}

	if v, ok := in["extra_args"].(map[string]interface{}); ok && len(v) > 0 {
		obj.ExtraArgs = toMapString(v)
	}

	if v, ok := in["extra_binds"].([]interface{}); ok && len(v) > 0 {
		obj.ExtraBinds = toArrayString(v)
	}

	if v, ok := in["extra_env"].([]interface{}); ok && len(v) > 0 {
		obj.ExtraEnv = toArrayString(v)
	}

	if v, ok := in["image"].(string); ok && len(v) > 0 {
		obj.Image = v
	}

	if v, ok := in["pod_security_policy"].(bool); ok {
		obj.PodSecurityPolicy = v
	}

	if v, ok := in["secrets_encryption_config"].([]interface{}); ok && len(v) > 0 {
		obj.SecretsEncryptionConfig = expandClusterRKEConfigServicesKubeAPISecretsEncryptionConfig(v)
	}

	if v, ok := in["service_cluster_ip_range"].(string); ok && len(v) > 0 {
		obj.ServiceClusterIPRange = v
	}

	if v, ok := in["service_node_port_range"].(string); ok && len(v) > 0 {
		obj.ServiceNodePortRange = v
	}

	return obj, nil
}
