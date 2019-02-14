package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

//Schemas

func etcdFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"ca_cert": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"cert": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"creation": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"external_urls": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
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
		"key": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"path": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"retention": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"snapshot": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
	}
	return s
}

// Flatteners

func flattenEtcd(in *managementClient.ETCDService) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.CACert) > 0 {
		obj["ca_cert"] = in.CACert
	}

	if len(in.Cert) > 0 {
		obj["cert"] = in.Cert
	}

	if len(in.Creation) > 0 {
		obj["creation"] = in.Creation
	}

	if len(in.ExternalURLs) > 0 {
		obj["external_urls"] = toArrayInterface(in.ExternalURLs)
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

	if len(in.Key) > 0 {
		obj["key"] = in.Key
	}

	if len(in.Path) > 0 {
		obj["path"] = in.Path
	}

	if len(in.Retention) > 0 {
		obj["retention"] = in.Retention
	}

	obj["snapshot"] = in.Snapshot

	return []interface{}{obj}, nil
}

// Expanders

func expandEtcd(p []interface{}) (*managementClient.ETCDService, error) {
	obj := &managementClient.ETCDService{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["ca_cert"].(string); ok && len(v) > 0 {
		obj.CACert = v
	}

	if v, ok := in["cert"].(string); ok && len(v) > 0 {
		obj.Cert = v
	}

	if v, ok := in["creation"].(string); ok && len(v) > 0 {
		obj.Creation = v
	}

	if v, ok := in["external_urls"].([]interface{}); ok && len(v) > 0 {
		obj.ExternalURLs = toArrayString(v)
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

	if v, ok := in["key"].(string); ok && len(v) > 0 {
		obj.Key = v
	}

	if v, ok := in["path"].(string); ok && len(v) > 0 {
		obj.Path = v
	}

	if v, ok := in["retention"].(string); ok && len(v) > 0 {
		obj.Retention = v
	}

	if v, ok := in["snapshot"].(bool); ok {
		obj.Snapshot = v
	}

	return obj, nil
}
