package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	clusterDriverRKE = "rancherKubernetesEngine"
)

func clusterLegacyAnyMapListElem() *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeMap,
	}
}

func clusterFieldsV2() map[string]*schema.Schema {
	clusterDriversV2 := append([]string{}, clusterDrivers...)
	clusterDriversV2 = append(clusterDriversV2, clusterDriverRKE)

	s := clusterFields()
	s["driver"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Computed:     true,
		ValidateFunc: validation.StringInSlice(clusterDriversV2, true),
	}
	s["rke_config"] = &schema.Schema{
		Type:          schema.TypeList,
		MaxItems:      1,
		Optional:      true,
		Computed:      true,
		ConflictsWith: []string{"aks_config_v2", "eks_config_v2", "gke_config_v2", "k3s_config", "oke_config", "rke2_config", "imported_config"},
		Elem:          clusterLegacyAnyMapListElem(),
	}
	s["cluster_template_answers"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Computed:    true,
		Description: "Cluster template answers",
		Elem:        clusterLegacyAnyMapListElem(),
	}
	s["cluster_template_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Cluster template ID",
	}
	s["cluster_template_questions"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    true,
		Description: "Cluster template questions",
		Elem:        clusterLegacyAnyMapListElem(),
	}
	s["cluster_template_revision_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Cluster template revision ID",
	}

	return s
}
