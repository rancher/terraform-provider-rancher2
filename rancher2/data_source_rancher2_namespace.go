package rancher2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRancher2Namespace() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRancher2NamespaceRead,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project ID where k8s namespace belongs",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the k8s namespace managed by rancher v2",
			},
			"container_resource_limit": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: containerResourceLimitFields(),
				},
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the k8s namespace managed by rancher v2",
			},
			"resource_quota": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: namespaceResourceQuotaFields(),
				},
			},
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Annotations of the k8s namespace managed by rancher v2",
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Labels of the k8s namespace managed by rancher v2",
			},
		},
	}
}

func dataSourceRancher2NamespaceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectID := d.Get("project_id").(string)
	clusterID, err := clusterIDFromProjectID(projectID)
	if err != nil {
		return diag.FromErr(err)
	}
	name := d.Get("name").(string)

	client, err := meta.(*Config).ClusterClient(clusterID)
	if err != nil {
		return diag.FromErr(err)
	}

	filters := map[string]interface{}{
		"projectId": projectID,
		"name":      name,
	}
	listOpts := NewListOpts(filters)

	namespaces, err := client.Namespace.List(listOpts)
	if err != nil {
		return diag.FromErr(err)
	}

	count := len(namespaces.Data)
	if count <= 0 {
		return diag.Errorf("[ERROR] namespace with name \"%s\" on cluster ID \"%s\" not found", name, clusterID)
	}
	if count > 1 {
		return diag.Errorf("[ERROR] found %d namespace with name \"%s\" on cluster ID \"%s\"", count, name, clusterID)
	}

	return diag.FromErr(flattenNamespace(d, &namespaces.Data[0]))
}
