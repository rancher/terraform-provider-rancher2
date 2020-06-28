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
	catalogHelmV2       = "helm_v2"
	catalogHelmV3       = "helm_v3"
)

var (
	catalogKinds        = []string{catalogKindHelm}
	catalogScopes       = []string{catalogScopeCluster, catalogScopeGlobal, catalogScopeProject}
	catalogHelmVersions = []string{catalogHelmV2, catalogHelmV3}
)

// Shemas

func catalogFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"url": {
			Type:     schema.TypeString,
			Required: true,
		},
		"branch": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "master",
		},
		"cluster_id": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"kind": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      catalogKindHelm,
			ValidateFunc: validation.StringInSlice(catalogKinds, true),
		},
		"password": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"project_id": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"refresh": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"scope": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      catalogScopeGlobal,
			ValidateFunc: validation.StringInSlice(catalogScopes, true),
		},
		"username": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"version": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Computed:     true,
			ValidateFunc: validation.StringInSlice(catalogHelmVersions, true),
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}
