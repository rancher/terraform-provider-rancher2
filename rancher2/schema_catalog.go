package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	catalogKindHelm     = "helm"
	catalogScopeCluster = "cluster"
	catalogScopeGlobal  = "global"
	catalogScopeProject = "project"
)

var (
	catalogKinds  = []string{catalogKindHelm}
	catalogScopes = []string{catalogScopeCluster, catalogScopeGlobal, catalogScopeProject}
)

// Shemas

func catalogFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"url": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"branch": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "master",
		},
		"cluster_id": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"kind": &schema.Schema{
			Type:         schema.TypeString,
			Optional:     true,
			Default:      catalogKindHelm,
			ValidateFunc: validation.StringInSlice(catalogKinds, true),
		},
		"password": &schema.Schema{
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"project_id": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"refresh": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"scope": &schema.Schema{
			Type:         schema.TypeString,
			Optional:     true,
			Default:      catalogScopeGlobal,
			ValidateFunc: validation.StringInSlice(catalogScopes, true),
		},
		"username": &schema.Schema{
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
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
