package rancher2

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceRancher2SecretV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2SecretV2Read,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
			},
			"immutable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"resource_version": {
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

func dataSourceRancher2SecretV2Read(d *schema.ResourceData, meta interface{}) error {
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	namespace := d.Get("namespace").(string)
	rancherID := namespace + "/" + name
	d.SetId(clusterID + secretV2ClusterIDsep + rancherID)

	secret, err := getSecretV2ByID(meta.(*Config), clusterID, rancherID)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Secret V2 %s not found at cluster %s", rancherID, clusterID)
			d.SetId("")
			return nil
		}
		return err
	}

	return flattenSecretV2(d, secret)
}
