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

// suppressFunc is a DiffSuppressFunc that prevents Terraform from removing Rancher-managed annotations or labels.
// It also ignores changes to a predefined set of Rancher-managed annotation/label keys.
//
// Thought it is not recommended, users can still add annotations or labels whose keys contain Rancher-managed keys,
// but they won't be able to remove them once added.
// Note: Terraform prefixes the key `k` with either "annotations." or "labels."
var suppressFunc = func(k, old, new string, d *schema.ResourceData) bool {
	for _, exception := range exceptions {
		// Explicitly check if the key ends with the exception
		if strings.HasSuffix(k, exception) {
			return false
		}
	}

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
			DiffSuppressFunc: suppressFunc,
		},
		"labels": {
			Type:             schema.TypeMap,
			Optional:         true,
			Computed:         true,
			Description:      "Labels of the resource",
			DiffSuppressFunc: suppressFunc,
		},
	}
	return s
}
