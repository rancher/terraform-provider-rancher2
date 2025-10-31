package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	clusterProxyConfigV2Name = "clusterproxyconfig"
)

func resourceRancher2ClusterProxyConfigV2Import(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceRancher2ClusterProxyConfigV2Read(d, meta)
	if err != nil || d.Id() == "" {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
