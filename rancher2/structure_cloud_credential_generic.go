package rancher2

// Flatteners

func flattenCloudCredentialGeneric(in *genericCredentialConfig, p []interface{}) []interface{} {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if in == nil {
		return []interface{}{}
	}

	if len(in.driverID) > 0 {
		obj["driver"] = in.driverID
	}
	if len(in.config) > 0 {
		obj["config"] = in.config
	}

	return []interface{}{obj}
}

// Expanders

func expandCloudCredentialGeneric(p []interface{}, client nodeDriverFinder) (*genericCredentialConfig, error) {
	obj := &genericCredentialConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["driver"].(string); ok && len(v) > 0 {
		obj.driverID = v
		nodeDriver, err := client.ByID(obj.driverID)
		if err != nil {
			return nil, err
		}
		obj.driverName = nodeDriver.Name
	}

	if v, ok := in["config"].(map[string]interface{}); ok {
		obj.config = v
	}

	return obj, nil
}
