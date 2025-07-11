package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	norman "github.com/rancher/norman/types"
	managementv3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
)

const (
	clusterProxyConfigV2ApiType    = "management.cattle.io.clusterproxyconfig"
	clusterProxyConfigV2Kind       = "ClusterProxyConfig"
	clusterProxyConfigV2APIVersion = "management.cattle.io/v3"
)

type ClusterProxyConfigV2 struct {
	norman.Resource
	managementv3.ClusterProxyConfig
}

// Flatteners

func flattenClusterProxyConfigV2(d *schema.ResourceData, in *ClusterProxyConfigV2) error {
	if in == nil {
		return fmt.Errorf("[ERROR] flattening cluster proxy config: Input unstructured is nil")
	}

	d.Set("cluster_id", (in.GetNamespace()))
	d.Set("enabled", in.Enabled)

	return nil
}

// Expanders

func expandClusterProxyConfigV2(in *schema.ResourceData) (*ClusterProxyConfigV2, error) {
	if in == nil {
		return nil, fmt.Errorf("[ERROR] expanding cluster proxy config: Input ResourceData is nil")
	}

	obj := &ClusterProxyConfigV2{}

	if len(in.Id()) > 0 {
		obj.ID = in.Id()
	}

	obj.TypeMeta.Kind = clusterProxyConfigV2Kind
	obj.TypeMeta.APIVersion = clusterProxyConfigV2APIVersion

	obj.ObjectMeta.Namespace = in.Get("cluster_id").(string)
	obj.ObjectMeta.Name = clusterProxyConfigV2Name
	obj.Enabled = in.Get("enabled").(bool)

	return obj, nil
}
