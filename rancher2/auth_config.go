package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	managementClient "github.com/rancher/types/client/management/v3"
)

const authDefaultAccessMode = "unrestricted"

var (
	authAccessModes = []string{"required", "restricted", "unrestricted"}
)

//Schemas

func authConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"access_mode": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(authAccessModes, true),
			Default:      authDefaultAccessMode,
		},
		"allowed_principal_ids": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"annotations": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"labels": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
	}
	return s
}

func getAuthConfigObject(kind string) (interface{}, error) {
	switch kind {
	case managementClient.ActiveDirectoryConfigType:
		return &managementClient.ActiveDirectoryConfig{}, nil
	case managementClient.ADFSConfigType:
		return &managementClient.ADFSConfig{}, nil
	case managementClient.AzureADConfigType:
		return &managementClient.AzureADConfig{}, nil
	case managementClient.FreeIpaConfigType:
		return &managementClient.LdapConfig{}, nil
	case managementClient.GithubConfigType:
		return &managementClient.GithubConfig{}, nil
	case managementClient.OpenLdapConfigType:
		return &managementClient.LdapConfig{}, nil
	case managementClient.PingConfigType:
		return &managementClient.PingConfig{}, nil
	default:
		return nil, fmt.Errorf("[ERROR] Auth config type %s not supported", kind)
	}
}
