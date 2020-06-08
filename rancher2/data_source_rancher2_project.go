package rancher2

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/rancher/norman/types"
)

func dataSourceRancher2Project() *schema.Resource {
	return &schema.Resource{
		Exists: dataSourceRancher2ProjectExists,
		Read:   dataSourceRancher2ProjectRead,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Description: "ID of the cluster to whom the project belongs",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "Name of the project",
				Type:        schema.TypeString,
				Required:    true,
			},
			"container_resource_limit": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: containerResourceLimitFields(),
				},
			},
			"description": {
				Description: "Description of the project",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"enable_project_monitoring": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable built-in project monitoring",
			},
			"pod_security_policy_template_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_quota": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: projectResourceQuotaFields(),
				},
			},
			"uuid": {
				Description: "UUID of the project",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: descriptions["annotations"],
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: descriptions["labels"],
			},
		},
	}
}

func dataSourceRancher2ProjectExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return false, err
	}

	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)

	projects, err := client.Project.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"cluster_id": clusterID,
			"name":       name,
		},
	})
	if err != nil {
		return false, err
	}

	cnt := len(projects.Data)
	if cnt <= 0 {
		return false, nil
	}
	if cnt > 1 {
		return false, fmt.Errorf("[ERROR] more than one project with specified name (\"%s\") found: %d", name, cnt)
	}

	// Only one project returned? Great...
	return true, nil
}

func dataSourceRancher2ProjectRead(d *schema.ResourceData, meta interface{}) error {
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	log.Printf("[INFO] Refreshing Rancher2 Project: %s (Cluster ID: %s)", name, clusterID)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	projects, err := client.Project.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"clusterId": clusterID,
			"name":      name,
		},
	})
	if err != nil {
		return err
	}

	cnt := len(projects.Data)
	if cnt <= 0 {
		return fmt.Errorf("[ERROR] project with name \"%s\" not found", name)
	}
	if cnt > 1 {
		return fmt.Errorf("[ERROR] more than one project with specified name (\"%s\") found: %d", name, cnt)
	}

	// Only one project returned? Great...
	project := projects.Data[0]
	d.SetId(project.ID)
	d.Set("cluster_id", project.ClusterID)
	d.Set("description", project.Description)
	d.Set("name", project.Name)
	d.Set("uuid", project.UUID)

	if project.ContainerDefaultResourceLimit != nil {
		containerLimit := flattenProjectContainerResourceLimit(project.ContainerDefaultResourceLimit)
		err := d.Set("container_resource_limit", containerLimit)
		if err != nil {
			return err
		}
	}

	d.Set("pod_security_policy_template_id", project.PodSecurityPolicyTemplateName)

	if project.ResourceQuota != nil && project.NamespaceDefaultResourceQuota != nil {
		resourceQuota := flattenProjectResourceQuota(project.ResourceQuota, project.NamespaceDefaultResourceQuota)
		err := d.Set("resource_quota", resourceQuota)
		if err != nil {
			return err
		}
	}

	err = d.Set("annotations", toMapInterface(project.Annotations))
	if err != nil {
		return err
	}
	err = d.Set("labels", toMapInterface(project.Labels))
	if err != nil {
		return err
	}

	return nil
}
