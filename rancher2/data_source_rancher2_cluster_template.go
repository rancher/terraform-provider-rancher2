package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceRancher2ClusterTemplate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2ClusterTemplateRead,

		Schema: map[string]*schema.Schema{
			"default_revision_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default cluster template revision ID",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Cluster template description",
			},
			"members": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Cluster template members",
				Elem: &schema.Resource{
					Schema: memberFields(),
				},
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cluster template name",
			},
			"template_revisions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Cluster template revisions",
				Elem: &schema.Resource{
					Schema: clusterTemplateRevisionFieldsData(),
				},
			},
			"annotations": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceRancher2ClusterTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	description := d.Get("description").(string)

	filters := map[string]interface{}{
		"name": name,
	}
	if len(description) > 0 {
		filters["description"] = description
	}
	listOpts := NewListOpts(filters)

	clusterTemplates, err := client.ClusterTemplate.List(listOpts)
	if err != nil {
		return err
	}

	count := len(clusterTemplates.Data)
	if count <= 0 {
		return fmt.Errorf("[ERROR] cluster template with name \"%s\" not found", name)
	}
	if count > 1 {
		return fmt.Errorf("[ERROR] found %d cluster template with name \"%s\"", count, name)
	}

	d.SetId(clusterTemplates.Data[0].ID)

	return resourceRancher2ClusterTemplateRead(d, meta)
}
