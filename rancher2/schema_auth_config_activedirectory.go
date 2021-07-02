package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const AuthConfigActiveDirectoryName = "activedirectory"

//Schemas

func authConfigActiveDirectoryFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"servers": {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"service_account_username": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"service_account_password": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"test_username": {
			Type:     schema.TypeString,
			Required: true,
		},
		"test_password": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"user_search_base": {
			Type:     schema.TypeString,
			Required: true,
		},
		"certificate": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
			StateFunc: TrimSpace,
		},
		"connection_timeout": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  5000,
		},
		"default_login_domain": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"group_dn_attribute": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"group_member_mapping_attribute": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"group_member_user_attribute": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"group_name_attribute": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"group_object_class": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"group_search_attribute": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"group_search_base": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"group_search_filter": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"nested_group_membership_enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"port": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  389,
		},
		"start_tls": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"tls": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"user_disabled_bit_mask": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  2,
		},
		"user_enabled_attribute": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"user_login_attribute": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"user_name_attribute": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"user_object_class": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"user_search_attribute": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"user_search_filter": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}

	for k, v := range authConfigFields() {
		s[k] = v
	}

	return s
}
