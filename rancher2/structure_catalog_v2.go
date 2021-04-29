package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/rancher/rancher/pkg/apis/catalog.cattle.io/v1"
)

// Flatteners

func flattenCatalogV2(d *schema.ResourceData, in *ClusterRepo) error {
	if in == nil {
		return nil
	}

	if len(in.ID) > 0 {
		d.SetId(d.Get("cluster_id").(string) + catalogV2ClusterIDsep + in.ID)
	}
	d.Set("name", in.ObjectMeta.Name)
	err := d.Set("annotations", toMapInterface(in.ObjectMeta.Annotations))
	if err != nil {
		return err
	}
	err = d.Set("labels", toMapInterface(in.ObjectMeta.Labels))
	if err != nil {
		return err
	}
	d.Set("resource_version", in.ObjectMeta.ResourceVersion)

	d.Set("ca_bundle", string(in.Spec.CABundle))
	if in.Spec.Enabled != nil {
		d.Set("enabled", *in.Spec.Enabled)
	}
	d.Set("git_branch", in.Spec.GitBranch)
	d.Set("git_repo", in.Spec.GitRepo)
	d.Set("insecure", in.Spec.InsecureSkipTLSverify)
	if in.Spec.ClientSecret != nil {
		d.Set("secret_name", in.Spec.ClientSecret.Name)
		d.Set("secret_namespace", in.Spec.ClientSecret.Namespace)
	}
	d.Set("service_account", in.Spec.ServiceAccount)
	d.Set("service_account_namespace", in.Spec.ServiceAccountNamespace)
	d.Set("url", in.Spec.URL)

	return nil
}

// Expanders

func expandCatalogV2(in *schema.ResourceData) *ClusterRepo {
	if in == nil {
		return nil
	}
	obj := &ClusterRepo{}

	if len(in.Id()) > 0 {
		_, obj.ID = splitID(in.Id())
	}
	obj.TypeMeta.Kind = catalogV2Kind
	obj.TypeMeta.APIVersion = catalogV2APIGroup + "/" + catalogV2APIVersion

	obj.ObjectMeta.Name = in.Get("name").(string)
	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.ObjectMeta.Annotations = toMapString(v)
	}
	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.ObjectMeta.Labels = toMapString(v)
	}
	if v, ok := in.Get("resource_version").(string); ok {
		obj.ObjectMeta.ResourceVersion = v
	}
	if v, ok := in.Get("ca_bundle").(string); ok {
		obj.Spec.CABundle = []byte(v)
	}
	if v, ok := in.Get("enabled").(bool); ok {
		obj.Spec.Enabled = &v
	}
	if v, ok := in.Get("git_repo").(string); ok && len(v) > 0 {
		obj.Spec.GitRepo = v
		obj.Spec.GitBranch = catalogV2DefaultGitBranch
	}
	if v, ok := in.Get("git_branch").(string); ok && len(v) > 0 {
		obj.Spec.GitBranch = v
	}
	if v, ok := in.Get("insecure").(bool); ok {
		obj.Spec.InsecureSkipTLSverify = v
	}
	sName, nok := in.Get("secret_name").(string)
	sNamespace, nsok := in.Get("secret_namespace").(string)
	if nok && nsok && len(sName) > 0 {
		obj.Spec.ClientSecret = &v1.SecretReference{
			Name:      sName,
			Namespace: sNamespace,
		}
	}
	if v, ok := in.Get("service_account").(string); ok {
		obj.Spec.ServiceAccount = v
	}
	if v, ok := in.Get("service_account_namespace").(string); ok {
		obj.Spec.ServiceAccountNamespace = v
	}
	if v, ok := in.Get("url").(string); ok {
		obj.Spec.URL = v
	}

	return obj
}
