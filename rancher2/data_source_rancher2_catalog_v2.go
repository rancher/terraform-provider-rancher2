package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/rancher/rancher/pkg/apis/catalog.cattle.io/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func dataSourceRancher2CatalogV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2CatalogV2Read,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ca_bundle": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"git_branch": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"git_repo": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"insecure": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"resource_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret_namespace": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_account": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_account_namespace": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"annotations": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceRancher2CatalogV2Read(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	clusterID := d.Get("cluster_id").(string)

	client, err := meta.(*Config).catalogV2Client(clusterID)
	if err != nil {
		return err
	}
	obj, err := client.Get(name, "", metaV1.GetOptions{ResourceVersion: d.Get("resource_version").(string)})
	if err != nil {
		if errors.IsNotFound(err) || errors.IsForbidden(err) {
			return fmt.Errorf("[ERROR] catalog v2 on cluster %s with name \"%s\" not found", clusterID, name)
		}
	}

	return flattenCatalogV2(d, obj.(*v1.ClusterRepo))
}
