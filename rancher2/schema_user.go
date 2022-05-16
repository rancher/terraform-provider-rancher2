package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var (
	userTokenFieldsList = []string{"token_id", "token_name", "token_enabled", "token_expired", "auth_token", "access_key", "secret_key"}
)

//Schemas

func userTokenConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cluster_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Cluster ID for scoped token",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Token description",
		},
		"ttl": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     0,
			Description: "Token time to live in seconds",
		},
		"renew": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Renew expired or disabled token",
		},
	}
	return s
}

func userFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"password": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"username": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"principal_ids": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"temp_token": {
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},
		"temp_token_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"login_role_binding_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"token_config": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "",
			Elem: &schema.Resource{
				Schema: userTokenConfigFields(),
			},
		},
		"token_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Token ID",
		},
		"token_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Token name",
		},
		"token_enabled": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Token enabled",
		},
		"token_expired": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Token expired",
		},
		"auth_token": {
			Type:        schema.TypeString,
			Computed:    true,
			Sensitive:   true,
			Description: "Token value",
		},
		"access_key": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Token access key",
		},
		"secret_key": {
			Type:        schema.TypeString,
			Computed:    true,
			Sensitive:   true,
			Description: "Token secret key",
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}

// Diffs
func userTokenComputedIf(keys []string) schema.CustomizeDiffFunc {
	diffs := make([]schema.CustomizeDiffFunc, 0, 7)
	for _, k := range keys {
		diffs = append(diffs, customdiff.ComputedIf(k, userTokenComputedConditionFunc()))
	}
	return customdiff.All(diffs...)
}

func userTokenComputedConditionFunc() customdiff.ResourceConditionFunc {
	return func(d *schema.ResourceDiff, meta interface{}) bool {
		if v, ok := d.Get("token_config").([]interface{}); ok && len(v) == 0 {
			return false
		}

		if d.HasChange("token_config") {
			return true
		}

		if r, ok := d.Get("token_config.0.renew").(bool); ok && r {
			if v, ok := d.Get("token_expired").(bool); ok && v {
				return true
			}
		}

		if v, ok := d.Get("token_config").([]interface{}); ok && len(v) > 0 {
			if v, ok := d.Get("token_id").([]interface{}); ok && len(v) == 0 {
				return true
			}
		}

		return false
	}
}
