package rancher2

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	commonAnnotationLabelCattle  = "cattle.io/"
	commonAnnotationLabelRancher = "rancher.io/"
)

// exceptions is a list of annotation/label keys that should not be suppressed.
var exceptions = []string{
	"rancher.io/imported-cluster-version-management",
}

// supressFunc is a DiffSuppressFunc that prevents Terraform from trying to remove Rancher-managed annotations/labels;
// it ignores the annotation from a predefined list of annotation/label keys.
var supressFunc = func(k, old, new string, d *schema.ResourceData) bool {
	for _, exception := range exceptions {
		// The key `k` is prefixed with "annotations." or "labels."
		// suppress the diff if the key contains the exception
		if strings.Contains(k, exception) {
			return false
		}
	}

	// Suppress the diff for Rancher-managed annotations/labels
	if (strings.Contains(k, commonAnnotationLabelCattle) || strings.Contains(k, commonAnnotationLabelRancher)) && new == "" {
		return true
	}

	return false
}

// Schemas

func commonAnnotationLabelFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"annotations": {
			Type:             schema.TypeMap,
			Optional:         true,
			Computed:         true,
			Description:      "Annotations of the resource",
			DiffSuppressFunc: supressFunc,
		},
		"labels": {
			Type:             schema.TypeMap,
			Optional:         true,
			Computed:         true,
			Description:      "Labels of the resource",
			DiffSuppressFunc: supressFunc,
		},
	}
	return s
}
