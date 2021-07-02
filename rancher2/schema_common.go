package rancher2

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	commonAnnotationLabelCattle  = "cattle.io/"
	commonAnnotationLabelRancher = "rancher.io/"
)

//Schemas

func commonAnnotationLabelFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"annotations": {
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: "Annotations of the resource",
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				// Supressing diff for annotations containing cattle.io/
				if (strings.Contains(k, commonAnnotationLabelCattle) || strings.Contains(k, commonAnnotationLabelRancher)) && new == "" {
					return true
				}
				return false
			},
		},
		"labels": {
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: "Labels of the resource",
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				// Supressing diff for labels containing cattle.io/
				if (strings.Contains(k, commonAnnotationLabelCattle) || strings.Contains(k, commonAnnotationLabelRancher)) && new == "" {
					return true
				}
				return false
			},
		},
	}
	return s
}
