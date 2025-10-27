package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	clusterScanCisConfigProfilePermissive = "permissive"
	clusterScanCisConfigProfileHardened   = "hardened"
)

var (
	clusterScanCisConfigProfileKinds = []string{clusterScanCisConfigProfilePermissive, clusterScanCisConfigProfileHardened}
)

//Schemas

func clusterScanCisConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"debug_master": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Debug master",
		},
		"debug_worker": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Debug worker",
		},
		"override_benchmark_version": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Override Benchmark Version",
		},
		"override_skip": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"profile": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      clusterScanCisConfigProfilePermissive,
			Description:  "Profile",
			ValidateFunc: validation.StringInSlice(clusterScanCisConfigProfileKinds, true),
		},
	}

	return s
}

func clusterScanConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cis_scan_config": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Description: "The cluster cis scan config",
			Elem: &schema.Resource{
				Schema: clusterScanCisConfigFields(),
			},
		},
	}

	return s
}

func clusterScanFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cluster_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The cluster ID to scan",
		},
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The cluster scan name",
		},
		"run_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The cluster scan run type",
		},
		"scan_config": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterScanConfigFields(),
			},
			Description: "The cluster scan run type",
		},
		"scan_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The cluster scan type",
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The cluster scan status",
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}
