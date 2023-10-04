package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenPodSecurityAdmissionConfigurationTemplateDefaults(in *managementClient.PodSecurityAdmissionConfigurationTemplateDefaults) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	obj := make(map[string]interface{})

	if len(in.Audit) > 0 {
		obj["audit"] = in.Audit
	}

	if len(in.AuditVersion) > 0 {
		obj["audit_version"] = in.AuditVersion
	}

	if len(in.Enforce) > 0 {
		obj["enforce"] = in.Enforce
	}

	if len(in.EnforceVersion) > 0 {
		obj["enforce_version"] = in.EnforceVersion
	}

	if len(in.Warn) > 0 {
		obj["warn"] = in.Warn
	}

	if len(in.WarnVersion) > 0 {
		obj["warn_version"] = in.WarnVersion
	}

	return []interface{}{obj}
}

func flattenPodSecurityAdmissionConfigurationTemplateExemptions(in *managementClient.PodSecurityAdmissionConfigurationTemplateExemptions) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	obj := make(map[string]interface{})

	if len(in.Usernames) > 0 {
		obj["usernames"] = in.Usernames
	}

	if len(in.RuntimeClasses) > 0 {
		obj["runtime_classes"] = in.RuntimeClasses
	}

	if len(in.Namespaces) > 0 {
		obj["namespaces"] = in.Namespaces
	}

	return []interface{}{obj}
}

func flattenPodSecurityAdmissionConfigurationTemplate(d *schema.ResourceData, in *managementClient.PodSecurityAdmissionConfigurationTemplate) error {
	if in == nil {
		return fmt.Errorf("[ERROR] flattening pod security admission configuration template: Input setting is nil")
	}

	d.SetId(in.ID)
	d.Set("name", in.Name)

	if len(in.Description) > 0 {
		d.Set("description", in.Description)
	}

	if in.Configuration.Defaults != nil {
		d.Set("defaults", flattenPodSecurityAdmissionConfigurationTemplateDefaults(in.Configuration.Defaults))
	}

	if in.Configuration.Exemptions != nil {
		d.Set("exemptions", flattenPodSecurityAdmissionConfigurationTemplateExemptions(in.Configuration.Exemptions))
	}

	return nil
}

func expandPodSecurityAdmissionConfigurationTemplateDefaults(in []interface{}) *managementClient.PodSecurityAdmissionConfigurationTemplateDefaults {
	if len(in) == 0 || in[0] == nil {
		return nil
	}

	obj := &managementClient.PodSecurityAdmissionConfigurationTemplateDefaults{}

	m := in[0].(map[string]interface{})

	if v, ok := m["audit"].(string); ok {
		obj.Audit = v
	}

	if v, ok := m["audit_version"].(string); ok {
		obj.AuditVersion = v
	}

	if v, ok := m["enforce"].(string); ok {
		obj.Enforce = v
	}

	if v, ok := m["enforce_version"].(string); ok {
		obj.EnforceVersion = v
	}

	if v, ok := m["warn"].(string); ok {
		obj.Warn = v
	}

	if v, ok := m["warn_version"].(string); ok {
		obj.WarnVersion = v
	}

	return obj
}

func expandPodSecurityAdmissionConfigurationTemplateExemptions(in []interface{}) *managementClient.PodSecurityAdmissionConfigurationTemplateExemptions {
	if len(in) == 0 || in[0] == nil {
		return nil
	}

	obj := &managementClient.PodSecurityAdmissionConfigurationTemplateExemptions{}

	m := in[0].(map[string]interface{})

	if v, ok := m["usernames"].([]interface{}); ok && len(v) > 0 {
		obj.Usernames = toArrayString(v)
	}

	if v, ok := m["runtime_classes"].([]interface{}); ok && len(v) > 0 {
		obj.RuntimeClasses = toArrayString(v)
	}

	if v, ok := m["namespaces"].([]interface{}); ok && len(v) > 0 {
		obj.Namespaces = toArrayString(v)
	}

	return obj
}

func expandPodSecurityAdmissionConfigurationTemplate(in *schema.ResourceData) (*managementClient.PodSecurityAdmissionConfigurationTemplate, error) {
	if in == nil {
		return nil, fmt.Errorf("[ERROR] expanding pod security admission configuration template: Input setting is nil")
	}

	obj := &managementClient.PodSecurityAdmissionConfigurationTemplate{}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	if v, ok := in.Get("name").(string); ok && len(v) > 0 {
		obj.Name = v
	}

	if v, ok := in.Get("description").(string); ok && len(v) > 0 {
		obj.Description = v
	}

	if v, ok := in.Get("defaults").([]interface{}); ok && len(v) > 0 {
		if obj.Configuration == nil {
			obj.Configuration = &managementClient.PodSecurityAdmissionConfigurationTemplateSpec{}
		}
		obj.Configuration.Defaults = expandPodSecurityAdmissionConfigurationTemplateDefaults(v)
	}

	if v, ok := in.Get("exemptions").([]interface{}); ok && len(v) > 0 {
		if obj.Configuration == nil {
			obj.Configuration = &managementClient.PodSecurityAdmissionConfigurationTemplateSpec{}
		}
		obj.Configuration.Exemptions = expandPodSecurityAdmissionConfigurationTemplateExemptions(v)
	}

	return obj, nil
}
