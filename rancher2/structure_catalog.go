package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenClusterCatalog(d *schema.ResourceData, in *managementClient.ClusterCatalog) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)
	d.Set("name", in.Name)
	d.Set("url", in.URL)
	d.Set("description", in.Description)
	d.Set("kind", in.Kind)
	d.Set("branch", in.Branch)

	if len(in.Password) > 0 {
		d.Set("password", in.Password)
	}

	d.Set("username", in.Username)
	d.Set("cluster_id", in.ClusterID)

	if len(in.HelmVersion) > 0 {
		d.Set("version", in.HelmVersion)
	}

	err := d.Set("annotations", toMapInterface(in.Annotations))
	if err != nil {
		return err
	}

	err = d.Set("labels", toMapInterface(in.Labels))
	if err != nil {
		return err
	}

	return nil
}

func flattenGlobalCatalog(d *schema.ResourceData, in *managementClient.Catalog) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)
	d.Set("name", in.Name)
	d.Set("url", in.URL)
	d.Set("description", in.Description)
	d.Set("kind", in.Kind)
	d.Set("branch", in.Branch)

	if len(in.Password) > 0 {
		d.Set("password", in.Password)
	}

	d.Set("username", in.Username)

	if len(in.HelmVersion) > 0 {
		d.Set("version", in.HelmVersion)
	}

	err := d.Set("annotations", toMapInterface(in.Annotations))
	if err != nil {
		return err
	}

	err = d.Set("labels", toMapInterface(in.Labels))
	if err != nil {
		return err
	}

	return nil
}

func flattenProjectCatalog(d *schema.ResourceData, in *managementClient.ProjectCatalog) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)
	d.Set("name", in.Name)
	d.Set("url", in.URL)
	d.Set("description", in.Description)
	d.Set("kind", in.Kind)
	d.Set("branch", in.Branch)

	if len(in.Password) > 0 {
		d.Set("password", in.Password)
	}

	d.Set("username", in.Username)
	d.Set("project_id", in.ProjectID)

	if len(in.HelmVersion) > 0 {
		d.Set("version", in.HelmVersion)
	}

	err := d.Set("annotations", toMapInterface(in.Annotations))
	if err != nil {
		return err
	}

	err = d.Set("labels", toMapInterface(in.Labels))
	if err != nil {
		return err
	}

	return nil
}

func flattenCatalog(d *schema.ResourceData, in interface{}) error {
	if in == nil {
		return nil
	}

	scope := d.Get("scope").(string)
	switch scope {
	case catalogScopeCluster:
		return flattenClusterCatalog(d, in.(*managementClient.ClusterCatalog))
	case catalogScopeGlobal:
		return flattenGlobalCatalog(d, in.(*managementClient.Catalog))
	case catalogScopeProject:
		return flattenProjectCatalog(d, in.(*managementClient.ProjectCatalog))
	default:
		return fmt.Errorf("[ERROR] Unsupported scope on catalog: %s", scope)
	}
}

// Expanders

func expandClusterCatalog(in *schema.ResourceData) *managementClient.ClusterCatalog {
	obj := &managementClient.ClusterCatalog{}

	if in == nil {
		return nil
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	obj.Name = in.Get("name").(string)
	obj.URL = in.Get("url").(string)
	obj.Description = in.Get("description").(string)
	obj.Kind = in.Get("kind").(string)
	obj.Branch = in.Get("branch").(string)
	obj.ClusterID = in.Get("cluster_id").(string)

	if v, ok := in.Get("password").(string); ok && len(v) > 0 {
		obj.Password = v
	}
	if v, ok := in.Get("username").(string); ok && len(v) > 0 {
		obj.Username = v
	}

	if v, ok := in.Get("version").(string); ok && len(v) > 0 {
		obj.HelmVersion = v
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj
}

func expandGlobalCatalog(in *schema.ResourceData) *managementClient.Catalog {
	obj := &managementClient.Catalog{}

	if in == nil {
		return nil
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	obj.Name = in.Get("name").(string)
	obj.URL = in.Get("url").(string)
	obj.Description = in.Get("description").(string)
	obj.Kind = in.Get("kind").(string)
	obj.Branch = in.Get("branch").(string)

	if v, ok := in.Get("password").(string); ok && len(v) > 0 {
		obj.Password = v
	}
	if v, ok := in.Get("username").(string); ok && len(v) > 0 {
		obj.Username = v
	}

	if v, ok := in.Get("version").(string); ok && len(v) > 0 {
		obj.HelmVersion = v
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj
}

func expandProjectCatalog(in *schema.ResourceData) *managementClient.ProjectCatalog {
	obj := &managementClient.ProjectCatalog{}

	if in == nil {
		return nil
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	obj.Name = in.Get("name").(string)
	obj.URL = in.Get("url").(string)
	obj.Description = in.Get("description").(string)
	obj.Kind = in.Get("kind").(string)
	obj.Branch = in.Get("branch").(string)
	obj.ProjectID = in.Get("project_id").(string)

	if v, ok := in.Get("password").(string); ok && len(v) > 0 {
		obj.Password = v
	}
	if v, ok := in.Get("username").(string); ok && len(v) > 0 {
		obj.Username = v
	}

	if v, ok := in.Get("version").(string); ok && len(v) > 0 {
		obj.HelmVersion = v
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj
}

func expandCatalog(in *schema.ResourceData) interface{} {
	if in == nil {
		return nil
	}

	scope := in.Get("scope").(string)
	switch scope {
	case catalogScopeCluster:
		return expandClusterCatalog(in)
	case catalogScopeGlobal:
		return expandGlobalCatalog(in)
	case catalogScopeProject:
		return expandProjectCatalog(in)
	default:
		return nil
	}
}
