package rancher2

import (
	"errors"
	"fmt"
)

func flattenClusterBaseNodePool(in BaseNodePool) map[string]interface{} {
	obj := make(map[string]interface{})
	obj["add_default_label"] = in.AddDefaultLabel
	obj["add_default_taint"] = in.AddDefaultTaint

	if len(in.AdditionalLabels) > 0 {
		obj["additional_labels"] = flattenClusterBaseNodePoolAdditionalLabels(in)
	}

	if len(in.AdditionalTaints) > 0 {
		obj["additional_taints"] = flattenClusterBaseNodePoolAdditionalTaints(in)
	}

	obj["name"] = in.Name

	return obj
}

func flattenClusterBaseNodePoolAdditionalLabels(in BaseNodePool) map[string]interface{} {
	additionalLabelsObj := make(map[string]interface{})
	for key, value := range in.AdditionalLabels {
		additionalLabelsObj[key] = value
	}

	return additionalLabelsObj
}

func flattenClusterBaseNodePoolAdditionalTaints(in BaseNodePool) []interface{} {
	additionalTaintObjs := make([]interface{}, 0, len(in.AdditionalTaints))
	for _, taint := range in.AdditionalTaints {
		additionalTaintObj := make(map[string]interface{})

		if len(taint.Effect) > 0 {
			additionalTaintObj["effect"] = taint.Effect
		}

		if len(taint.Key) > 0 {
			additionalTaintObj["key"] = taint.Key
		}

		if len(taint.Value) > 0 {
			additionalTaintObj["value"] = taint.Value
		}

		additionalTaintObjs = append(additionalTaintObjs, additionalTaintObj)
	}

	return additionalTaintObjs
}

func expandClusterBaseNodePool(in map[string]interface{}) (BaseNodePool, error) {
	var obj BaseNodePool
	if v, ok := in["name"].(string); ok && len(v) > 0 {
		obj.Name = v
	} else {
		return obj, errors.New("'name' field must be provided for all pools")
	}

	if v, ok := in["add_default_label"].(bool); ok {
		obj.AddDefaultLabel = v
	}

	if v, ok := in["add_default_taint"].(bool); ok {
		obj.AddDefaultTaint = v
	}

	if v, ok := in["additional_labels"].(map[string]interface{}); ok && len(v) > 0 {
		obj.AdditionalLabels = toMapString(v)
	}

	if v, ok := in["additional_taints"].([]interface{}); ok && len(v) > 0 {
		additionalTaintsObjs, err := expandClusterBaseNodePoolAdditionalTaints(v, obj.Name)
		if err != nil {
			return obj, err
		}

		obj.AdditionalTaints = additionalTaintsObjs
	}

	return obj, nil
}

func expandClusterBaseNodePoolAdditionalTaints(additionalTaintsIn []interface{}, poolName string) ([]K8sTaint, error) {
	additionalTaintsObjs := make([]K8sTaint, 0, len(additionalTaintsIn))
	for index, additionalTaintIn := range additionalTaintsIn {
		if t, ok := additionalTaintIn.(map[string]interface{}); ok {
			taint := toMapString(t)
			additionalTaintsObj := K8sTaint{}

			if effect, ok := taint["effect"]; ok && len(effect) > 0 {
				additionalTaintsObj.Effect = effect
			}

			if key, ok := taint["key"]; ok && len(key) > 0 {
				additionalTaintsObj.Key = key
			}

			if value, ok := taint["value"]; ok && len(value) > 0 {
				additionalTaintsObj.Value = value
			}

			additionalTaintsObjs = append(additionalTaintsObjs, additionalTaintsObj)
		} else {
			return nil, fmt.Errorf("taint in index %d for worker pool %s contains unexpected content", index, poolName)
		}
	}
	return additionalTaintsObjs, nil
}
