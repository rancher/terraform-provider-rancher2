package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

const (
	questionTypeInt      = "int"
	questionTypeBool     = "boolean"
	questionTypeString   = "string"
	questionTypePassword = "password"
)

var (
	questionTypeKinds = []string{questionTypeInt, questionTypeBool, questionTypeString, questionTypePassword}
)

func getConflicts(fieldNames []string, fieldName string) []string {
	conflicts := make([]string, 0, len(fieldNames)-1)
	for _, name := range fieldNames {
		if name != fieldName {
			conflicts = append(conflicts, name)
		}
	}
	return conflicts
}

func questionFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"default": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Default variable value",
		},
		"required": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Required variable",
		},
		"type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      questionTypeString,
			Description:  "Variable type",
			ValidateFunc: validation.StringInSlice(questionTypeKinds, true),
		},
		"variable": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Variable name",
		},
	}

	return s
}

func expandQuestions(p []interface{}) []managementClient.Question {
	if len(p) == 0 || p[0] == nil {
		return nil
	}

	obj := make([]managementClient.Question, len(p))

	for i := range p {
		in := p[i].(map[string]interface{})

		if v, ok := in["default"].(string); ok && len(v) > 0 {
			obj[i].Default = v
		}

		if v, ok := in["required"].(bool); ok {
			obj[i].Required = v
		}

		if v, ok := in["type"].(string); ok && len(v) > 0 {
			obj[i].Type = v
		}

		if v, ok := in["variable"].(string); ok && len(v) > 0 {
			obj[i].Variable = v
		}
	}

	return obj
}

func flattenQuestions(p []managementClient.Question) []interface{} {
	if len(p) == 0 {
		return []interface{}{}
	}

	out := make([]interface{}, len(p))
	for i, in := range p {
		obj := make(map[string]interface{})

		obj["default"] = in.Default
		obj["required"] = in.Required
		obj["type"] = in.Type
		obj["variable"] = in.Variable

		out[i] = obj
	}

	return out
}
