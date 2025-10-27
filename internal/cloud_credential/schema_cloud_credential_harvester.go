package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	harvesterClusterTypeImported = "imported"
	harvesterClusterTypeExternal = "external"
)

var (
	harvesterClusterType = []string{harvesterClusterTypeImported, harvesterClusterTypeExternal}
)

//Types

type harvesterCredentialConfig struct {
	ClusterID         string `json:"clusterId,omitempty" yaml:"clusterId,omitempty"`
	ClusterType       string `json:"clusterType,omitempty" yaml:"clusterType,omitempty"`
	KubeconfigContent string `json:"kubeconfigContent,omitempty" yaml:"kubeconfigContent,omitempty"`
}

//Schemas

func cloudCredentialHarvesterFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cluster_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The cluster id of imported Harvester cluster",
		},
		"cluster_type": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "Harvester cluster type. must be imported or external",
			ValidateFunc: validation.StringInSlice(harvesterClusterType, true),
		},
		"kubeconfig_content": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Harvester cluster kubeconfig content",
		},
	}

	return s
}
