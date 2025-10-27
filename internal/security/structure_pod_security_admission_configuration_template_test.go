package rancher2

import (
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
)

var (
	defaultAll = &managementClient.PodSecurityAdmissionConfigurationTemplateDefaults{
		Audit:          "restricted",
		AuditVersion:   "latest",
		Enforce:        "restricted",
		EnforceVersion: "latest",
		Warn:           "restricted",
		WarnVersion:    "latest",
	}
	defaultAllRaw = []interface{}{
		map[string]interface{}{
			"audit":           "restricted",
			"audit_version":   "latest",
			"enforce":         "restricted",
			"enforce_version": "latest",
			"warn":            "restricted",
			"warn_version":    "latest",
		},
	}
	defaultPartly = &managementClient.PodSecurityAdmissionConfigurationTemplateDefaults{
		Audit:          "restricted",
		Enforce:        "restricted",
		EnforceVersion: "latest",
		Warn:           "restricted",
		WarnVersion:    "latest",
	}
	defaultPartlyWithDefaultValues = &managementClient.PodSecurityAdmissionConfigurationTemplateDefaults{
		Audit:          "restricted",
		AuditVersion:   "latest",
		Enforce:        "restricted",
		EnforceVersion: "latest",
		Warn:           "restricted",
		WarnVersion:    "latest",
	}
	defaultPartlyRaw = []interface{}{
		map[string]interface{}{
			"audit":           "restricted",
			"enforce":         "restricted",
			"enforce_version": "latest",
			"warn":            "restricted",
			"warn_version":    "latest",
		},
	}
	defaultPartlyWithEmptyValuesRaw = []interface{}{
		map[string]interface{}{
			"audit":           "restricted",
			"audit_version":   "",
			"enforce":         "restricted",
			"enforce_version": "latest",
			"warn":            "restricted",
			"warn_version":    "latest",
		},
	}
	defaultNone    = &managementClient.PodSecurityAdmissionConfigurationTemplateDefaults{}
	defaultNoneRaw = []interface{}{
		map[string]interface{}{},
	}
	exemptionsAll = &managementClient.PodSecurityAdmissionConfigurationTemplateExemptions{
		Usernames: []string{
			"admin1",
			"test-user-2",
		},
		Namespaces: []string{
			"cattle-system",
			"team-domain-1",
		},
		RuntimeClasses: []string{
			"runtime-cluster-1",
			"runtime-cluster-2",
		},
	}
	exemptionsAllRaw = []interface{}{
		map[string]interface{}{
			"usernames": []interface{}{
				"admin1",
				"test-user-2",
			},
			"runtime_classes": []interface{}{
				"runtime-cluster-1",
				"runtime-cluster-2",
			},
			"namespaces": []interface{}{
				"cattle-system",
				"team-domain-1",
			},
		},
	}
	exemptionsPartly = &managementClient.PodSecurityAdmissionConfigurationTemplateExemptions{
		Usernames: []string{
			"admin1",
			"test-user-2",
		},
		Namespaces: []string{
			"cattle-system",
			"team-domain-1",
		},
	}
	exemptionsPartlyWithDefaultValues = &managementClient.PodSecurityAdmissionConfigurationTemplateExemptions{
		Usernames: []string{
			"admin1",
			"test-user-2",
		},
		Namespaces: []string{
			"cattle-system",
			"team-domain-1",
		},
		RuntimeClasses: nil,
	}
	exemptionsPartlyRaw = []interface{}{
		map[string]interface{}{
			"usernames": []interface{}{
				"admin1",
				"test-user-2",
			},
			"namespaces": []interface{}{
				"cattle-system",
				"team-domain-1",
			},
		},
	}
	exemptionsPartlyWithEmptyValuesRaw = []interface{}{
		map[string]interface{}{
			"usernames": []interface{}{
				"admin1",
				"test-user-2",
			},
			"namespaces": []interface{}{
				"cattle-system",
				"team-domain-1",
			},
			"runtime_classes": []interface{}{},
		},
	}
	exemptionsNone    = &managementClient.PodSecurityAdmissionConfigurationTemplateExemptions{}
	exemptionsNoneRaw = []interface{}{
		map[string]interface{}{},
	}
	templateAll = &managementClient.PodSecurityAdmissionConfigurationTemplate{
		Name:        "template-all",
		Description: "this is the template for testing",
		Configuration: &managementClient.PodSecurityAdmissionConfigurationTemplateSpec{
			Defaults:   defaultAll,
			Exemptions: exemptionsAll,
		},
		Annotations: map[string]string{
			"purpose": "testing",
			"team":    "dev",
		},
		Labels: map[string]string{
			"purpose": "testing",
			"team":    "dev",
		},
	}
	templateAllRaw = map[string]interface{}{
		"name":        "template-all",
		"description": "this is the template for testing",
		"defaults":    defaultAllRaw,
		"exemptions":  exemptionsAllRaw,
		"annotations": map[string]interface{}{
			"purpose": "testing",
			"team":    "dev",
		},
		"labels": map[string]interface{}{
			"purpose": "testing",
			"team":    "dev",
		},
	}
	templatePartly = &managementClient.PodSecurityAdmissionConfigurationTemplate{
		Name:        "template-partly",
		Description: "this is the template for testing",
		Configuration: &managementClient.PodSecurityAdmissionConfigurationTemplateSpec{
			Defaults:   defaultPartly,
			Exemptions: exemptionsPartly,
		},
		Annotations: map[string]string{
			"purpose": "testing",
			"team":    "dev",
		},
		Labels: map[string]string{
			"purpose": "testing",
			"team":    "dev",
		},
	}
	templatePartlyWithDefaultValues = &managementClient.PodSecurityAdmissionConfigurationTemplate{
		Name:        "template-partly",
		Description: "this is the template for testing",
		Configuration: &managementClient.PodSecurityAdmissionConfigurationTemplateSpec{
			Defaults:   defaultPartlyWithDefaultValues,
			Exemptions: exemptionsPartlyWithDefaultValues,
		},
		Annotations: map[string]string{
			"purpose": "testing",
			"team":    "dev",
		},
		Labels: map[string]string{
			"purpose": "testing",
			"team":    "dev",
		},
	}
	templatePartlyRaw = map[string]interface{}{
		"name":        "template-partly",
		"description": "this is the template for testing",
		"defaults":    defaultPartlyRaw,
		"exemptions":  exemptionsPartlyRaw,
		"annotations": map[string]interface{}{
			"purpose": "testing",
			"team":    "dev",
		},
		"labels": map[string]interface{}{
			"purpose": "testing",
			"team":    "dev",
		},
	}
	templatePartlyWithEmptyValuesRaw = map[string]interface{}{
		"name":        "template-partly",
		"description": "this is the template for testing",
		"defaults":    defaultPartlyWithEmptyValuesRaw,
		"exemptions":  exemptionsPartlyWithEmptyValuesRaw,
		"annotations": map[string]interface{}{
			"purpose": "testing",
			"team":    "dev",
		},
		"labels": map[string]interface{}{
			"purpose": "testing",
			"team":    "dev",
		},
	}

	templateNone    = &managementClient.PodSecurityAdmissionConfigurationTemplate{}
	templateNoneRaw = map[string]interface{}{}
)

func TestFlattenPodSecurityAdmissionConfigurationTemplateDefaults(t *testing.T) {
	type args struct {
		in *managementClient.PodSecurityAdmissionConfigurationTemplateDefaults
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{
			name: "with all fields",
			args: args{defaultAll},
			want: defaultAllRaw,
		},
		{
			name: "with some fields",
			args: args{defaultPartly},
			want: defaultPartlyRaw,
		},
		{
			name: "with no fields",
			args: args{defaultNone},
			want: defaultNoneRaw,
		},
		{
			name: "with nil value",
			args: args{nil},
			want: []interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, flattenPodSecurityAdmissionConfigurationTemplateDefaults(tt.args.in), "flattenPodSecurityAdmissionConfigurationTemplateDefaults(%v)", tt.args.in)
		})
	}
}

func TestFlattenPodSecurityAdmissionConfigurationTemplateExemptions(t *testing.T) {
	type args struct {
		in *managementClient.PodSecurityAdmissionConfigurationTemplateExemptions
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{
			name: "with all fields",
			args: args{exemptionsAll},
			want: exemptionsAllRaw,
		},
		{
			name: "with some fields",
			args: args{exemptionsPartly},
			want: exemptionsPartlyRaw,
		},
		{
			name: "with no fields",
			args: args{exemptionsNone},
			want: exemptionsNoneRaw,
		},
		{
			name: "with nil value",
			args: args{nil},
			want: []interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, flattenPodSecurityAdmissionConfigurationTemplateExemptions(tt.args.in), "flattenPodSecurityAdmissionConfigurationTemplateExemptions(%v)", tt.args.in)
		})
	}
}

func TestFlattenPodSecurityAdmissionConfigurationTemplate(t *testing.T) {
	tests := []struct {
		Name                 string
		Input                *managementClient.PodSecurityAdmissionConfigurationTemplate
		ExpectedOutput       map[string]interface{}
		ExpectedErrorMessage string
	}{
		{
			Name:           "with values for all fields",
			Input:          templateAll,
			ExpectedOutput: templateAllRaw,
		},
		{
			Name:           "with values for some fields",
			Input:          templatePartly,
			ExpectedOutput: templatePartlyWithEmptyValuesRaw,
		},
		{
			Name:           "with no value for any field",
			Input:          templateNone,
			ExpectedOutput: templateNoneRaw,
		},
		{
			Name:                 "with nil value",
			Input:                nil,
			ExpectedErrorMessage: "[ERROR] flattening pod security admission configuration template: Input setting is nil",
		},
	}

	for _, tc := range tests {
		resourceData := schema.TestResourceDataRaw(t, podSecurityAdmissionConfigurationTemplateFields(), nil)
		err := flattenPodSecurityAdmissionConfigurationTemplate(resourceData, tc.Input)
		if err != nil {
			if tc.ExpectedErrorMessage == "" {
				assert.FailNow(t, "No error is expected")
			} else {
				if !strings.Contains(err.Error(), nilInputErrMsg) {
					assert.FailNow(t, "failed to get the expected error", "expected error should contain: %#v actual: %#v", nilInputErrMsg, err.Error())
				}
			}
			break
		}
		result := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			result[k] = resourceData.Get(k)
		}
		if !reflect.DeepEqual(tc.ExpectedOutput, result) {
			assert.FailNow(t, "Unexpected output from flattener", "test case: %v \nexpected: %#v \nactual: %#v", tc.Name, tc.ExpectedOutput, result)
		}
	}
}

func TestExpandPodSecurityAdmissionConfigurationTemplateDefaults(t *testing.T) {
	type args struct {
		in []interface{}
	}
	tests := []struct {
		name string
		args args
		want *managementClient.PodSecurityAdmissionConfigurationTemplateDefaults
	}{
		{
			name: "with values for all fields",
			args: args{defaultAllRaw},
			want: defaultAll,
		},
		{
			name: "with values for some fields",
			args: args{defaultPartlyRaw},
			want: defaultPartly,
		},
		{
			name: "with no value for any fields",
			args: args{defaultNoneRaw},
			want: defaultNone,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, expandPodSecurityAdmissionConfigurationTemplateDefaults(tt.args.in), "expandPodSecurityAdmissionConfigurationTemplateDefaults(%v)", tt.args.in)
		})
	}
}

func TestExpandPodSecurityAdmissionConfigurationTemplateExemptions(t *testing.T) {
	type args struct {
		in []interface{}
	}
	tests := []struct {
		name string
		args args
		want *managementClient.PodSecurityAdmissionConfigurationTemplateExemptions
	}{
		{
			name: "with values for all fields",
			args: args{exemptionsAllRaw},
			want: exemptionsAll,
		},
		{
			name: "with values for some fields",
			args: args{exemptionsPartlyRaw},
			want: exemptionsPartly,
		},
		{
			name: "with no value for any fields",
			args: args{exemptionsNoneRaw},
			want: exemptionsNone,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, expandPodSecurityAdmissionConfigurationTemplateExemptions(tt.args.in), "expandPodSecurityAdmissionConfigurationTemplateExemptions(%v)", tt.args.in)
		})
	}
}

func TestExpandPodSecurityAdmissionConfigurationTemplate(t *testing.T) {
	tests := []struct {
		Name           string
		Input          map[string]interface{}
		ExpectedOutput *managementClient.PodSecurityAdmissionConfigurationTemplate
	}{
		{
			Name:           "with values for all fields",
			Input:          templateAllRaw,
			ExpectedOutput: templateAll,
		},
		{
			Name:           "with values for some fields",
			Input:          templatePartlyRaw,
			ExpectedOutput: templatePartlyWithDefaultValues,
		},
		{
			Name:           "with no value for any field",
			Input:          templateNoneRaw,
			ExpectedOutput: templateNone,
		},
	}

	for _, tc := range tests {
		resourceData := schema.TestResourceDataRaw(t, podSecurityAdmissionConfigurationTemplateFields(), tc.Input)
		output, err := expandPodSecurityAdmissionConfigurationTemplate(resourceData)
		if err != nil {
			assert.FailNow(t, "No error is expected", "err: %s", err.Error())
		}
		if !reflect.DeepEqual(tc.ExpectedOutput, output) {
			assert.FailNow(t, "Unexpected output from expander", "test case: %v \nexpected: %#v \nactual: %#v",
				tc.Name, tc.ExpectedOutput, output)
		}
	}

	// nil case
	_, err := expandPodSecurityAdmissionConfigurationTemplate(nil)
	if err != nil && !strings.Contains(err.Error(), nilInputErrMsg) {
		assert.FailNow(t, "failed to get the expected error", "expected error should contain: %#v actual: %#v", nilInputErrMsg, err.Error())
	}
}
