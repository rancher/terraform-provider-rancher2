package rancher2

import managementClient "github.com/rancher/types/client/management/v3"

// Expanders

type nodeDriverFinder interface {
	ByID(id string) (*managementClient.NodeDriver, error)
}

func expandGenericNodeTemplateConfig(p []interface{}, client nodeDriverFinder) (*genericNodeTemplateConfig, error) {
	obj := &genericNodeTemplateConfig{}
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
