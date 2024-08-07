package rancher2

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			"exponential_backoff_max_wait": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Maximum amount of seconds to wait before retrying",
			},
			"exponential_backoff_min_wait": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Minimum amount of seconds to wait before retrying",
			},
			"exponential_backoff_max_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Maximum number of retries before returning error",
			},
			"git_branch": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"git_repo": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"insecure_plain_http": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Only valid for OCI URL's. Allows insecure connections to registries without enforcing TLS checks",
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
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)

	catalog, err := getCatalogV2ByID(meta.(*Config), clusterID, name)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Catalog V2 %s not found at cluster %s", name, clusterID)
			d.SetId("")
			return nil
		}
		return err
	}

	return flattenCatalogV2(d, catalog)
}
