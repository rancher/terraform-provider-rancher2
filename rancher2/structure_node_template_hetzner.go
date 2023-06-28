package rancher2

import "strings"

// Flatteners

func flattenHetznerConfig(in *hetznerConfig) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	if len(in.APIToken) > 0 {
		obj["api_token"] = in.APIToken
	}

	if len(in.Image) > 0 {
		obj["image"] = in.Image
	}

	if len(in.ImageID) > 0 {
		obj["image_id"] = in.ImageID
	}

	if len(in.ServerLabels) > 0 {
		obj["server_labels"] = toMapInterface(in.ServerLabels)
	}

	if len(in.ServerLocation) > 0 {
		obj["server_location"] = in.ServerLocation
	}

	if len(in.ServerType) > 0 {
		obj["server_type"] = in.ServerType
	}

	if len(in.Networks) > 0 {
		obj["networks"] = strings.Join(in.Networks, ",")
	}

	obj["use_private_network"] = in.UsePrivateNetwork

	if len(in.UserData) > 0 {
		obj["userdata"] = in.UserData
	}

	if len(in.Volumes) > 0 {
		obj["volumes"] = strings.Join(in.Volumes, ",")
	}

	if len(in.Firewalls) > 0 {
		obj["firewalls"] = toArrayInterface(in.Firewalls)
	}

	if len(in.AdditionalKeys) > 0 {
		obj["additional_keys"] = toArrayInterface(in.AdditionalKeys)
	}

	if len(in.PlacementGroup) > 0 {
		obj["placement_group"] = in.PlacementGroup
	}

	return []interface{}{obj}
}

// Expanders

func expandHetznercloudConfig(p []interface{}) *hetznerConfig {
	obj := &hetznerConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["api_token"].(string); ok && len(v) > 0 {
		obj.APIToken = v
	}

	if v, ok := in["image"].(string); ok && len(v) > 0 {
		obj.Image = v
	}

	if v, ok := in["image_id"].(string); ok && len(v) > 0 {
		obj.ImageID = v
	}

	if v, ok := in["server_labels"].(map[string]interface{}); ok && len(v) > 0 {
		obj.ServerLabels = toMapString(v)
	}

	if v, ok := in["server_location"].(string); ok && len(v) > 0 {
		obj.ServerLocation = v
	}

	if v, ok := in["server_type"].(string); ok && len(v) > 0 {
		obj.ServerType = v
	}

	if v, ok := in["networks"].(string); ok && len(v) > 0 {
		obj.Networks = strings.Split(v, ",")
	}

	if v, ok := in["use_private_network"].(bool); ok {
		obj.UsePrivateNetwork = v
	}

	if v, ok := in["userdata"].(string); ok && len(v) > 0 {
		obj.UserData = v
	}

	if v, ok := in["volumes"].(string); ok && len(v) > 0 {
		obj.Volumes = strings.Split(v, ",")
	}

	if v, ok := in["firewalls"].([]interface{}); ok && len(v) > 0 {
		obj.Firewalls = toArrayString(v)
	}

	if v, ok := in["additional_keys"].([]interface{}); ok && len(v) > 0 {
		obj.AdditionalKeys = toArrayString(v)
	}

	if v, ok := in["placement_group"].(string); ok && len(v) > 0 {
		obj.PlacementGroup = v
	}

	return obj
}
