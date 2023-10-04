package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	podSecurityAdmissionConfigurationDefaultMode    = "privileged"
	podSecurityAdmissionConfigurationDefaultVersion = "latest"
)

var (
	podSecurityAdmissionConfigurationModeTypes = []string{"privileged", "baseline", "restricted"}
)

//Schemas

func podSecurityAdmissionConfigurationTemplateDefaultFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"audit": {
			Type:         schema.TypeString,
			Default:      podSecurityAdmissionConfigurationDefaultMode,
			Description:  "Pod Security Admission Configuration audit. This audits a pod in violation of privileged, baseline, or restricted policy (default: privileged)",
			Optional:     true,
			ValidateFunc: validation.StringInSlice(podSecurityAdmissionConfigurationModeTypes, true),
		},
		"audit_version": {
			Type:        schema.TypeString,
			Default:     podSecurityAdmissionConfigurationDefaultVersion,
			Description: "Pod Security Admission Configuration audit version (default: latest)",
			Optional:    true,
		},
		"enforce": {
			Type:         schema.TypeString,
			Default:      podSecurityAdmissionConfigurationDefaultMode,
			Description:  "Pod Security Admission Configuration enforce. This rejects a pod in violation of privileged, baseline, or restricted policy (default: privileged)",
			Optional:     true,
			ValidateFunc: validation.StringInSlice(podSecurityAdmissionConfigurationModeTypes, true),
		},
		"enforce_version": {
			Type:        schema.TypeString,
			Default:     podSecurityAdmissionConfigurationDefaultVersion,
			Description: "Pod Security Admission Configuration enforce version (default: latest)",
			Optional:    true,
		},
		"warn": {
			Type:         schema.TypeString,
			Default:      podSecurityAdmissionConfigurationDefaultMode,
			Description:  "Pod Security Admission Configuration warn. This warns the user about a pod in violation of privileged, baseline, or restricted policy (default: privileged)",
			Optional:     true,
			ValidateFunc: validation.StringInSlice(podSecurityAdmissionConfigurationModeTypes, true),
		},
		"warn_version": {
			Type:        schema.TypeString,
			Default:     podSecurityAdmissionConfigurationDefaultVersion,
			Description: "Pod Security Admission Configuration warn version (default: latest)",
			Optional:    true,
		},
	}

	return s
}

func podSecurityAdmissionConfigurationTemplateExemptionFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"usernames": {
			Type:        schema.TypeList,
			Description: "Pod Security Admission Configuration username exemptions",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			ForceNew:    true,
		},
		"runtime_classes": {
			Type:        schema.TypeList,
			Description: "Pod Security Admission Configuration runtime class exemptions",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			ForceNew:    true,
		},
		"namespaces": {
			Type:        schema.TypeList,
			Description: "Pod Security Admission Configuration namespace exemptions",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			ForceNew:    true,
		},
	}

	return s
}

func podSecurityAdmissionConfigurationTemplateFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Pod Security Admission Configuration template name",
			ForceNew:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Pod Security Admission Configuration template description",
		},
		"defaults": {
			Type:        schema.TypeList,
			Description: "defaults allows the user to define admission control mode for Pod Security",
			Required:    true,
			Elem: &schema.Resource{
				Schema: podSecurityAdmissionConfigurationTemplateDefaultFields(),
			},
		},
		"exemptions": {
			Type:        schema.TypeList,
			Description: "exemptions allows the creation of pods for specific Usernames, RuntimeClassNames, and Namespaces that would otherwise be prohibited",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: podSecurityAdmissionConfigurationTemplateExemptionFields(),
			},
			ForceNew: true,
		},
	}

	return s
}
