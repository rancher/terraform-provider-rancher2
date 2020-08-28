package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func dataSourceRancher2Catalog() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2CatalogRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"scope": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      catalogScopeGlobal,
				ValidateFunc: validation.StringInSlice(catalogScopes, true),
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kind": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"branch": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"username": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"version": {
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

func dataSourceRancher2CatalogRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	scope := d.Get("scope").(string)

	catalogs, err := meta.(*Config).GetCatalogByName(name, scope)
	if err != nil {
		return err
	}

	switch scope {
	case catalogScopeCluster:
		err = dataSourceRancher2CatalogCheck(len(catalogs.(*managementClient.ClusterCatalogCollection).Data), scope, name)
		if err != nil {
			return err
		}
		return flattenCatalog(d, &catalogs.(*managementClient.ClusterCatalogCollection).Data[0])
	case catalogScopeGlobal:
		err = dataSourceRancher2CatalogCheck(len(catalogs.(*managementClient.CatalogCollection).Data), scope, name)
		if err != nil {
			return err
		}
		return flattenCatalog(d, &catalogs.(*managementClient.CatalogCollection).Data[0])
	case catalogScopeProject:
		err = dataSourceRancher2CatalogCheck(len(catalogs.(*managementClient.ProjectCatalogCollection).Data), scope, name)
		if err != nil {
			return err
		}
		return flattenCatalog(d, &catalogs.(*managementClient.ProjectCatalogCollection).Data[0])
	default:
		return fmt.Errorf("[ERROR] Unsupported scope on catalog: %s", scope)
	}

}

func dataSourceRancher2CatalogCheck(i int, scope, name string) error {
	if i <= 0 {
		return fmt.Errorf("[ERROR] %s catalog with name \"%s\" not found", scope, name)
	}
	if i > 1 {
		return fmt.Errorf("[ERROR] found %d %s catalogs with name \"%s\"", i, scope, name)
	}
	return nil
}
