package rancher2

import (
	"encoding/base64"
	"fmt"

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

	encodedCABundle := base64.StdEncoding.EncodeToString(in.Spec.CABundle)
	d.Set("ca_bundle", encodedCABundle)
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

	if in.Spec.ExponentialBackOffValues != nil {
		d.Set("exponential_backoff_min_wait", in.Spec.ExponentialBackOffValues.MinWait)
		d.Set("exponential_backoff_max_wait", in.Spec.ExponentialBackOffValues.MaxWait)
		d.Set("exponential_backoff_max_retries", in.Spec.ExponentialBackOffValues.MaxRetries)
	}

	d.Set("insecure_plain_http", in.Spec.InsecurePlainHTTP)

	d.Set("service_account", in.Spec.ServiceAccount)
	d.Set("service_account_namespace", in.Spec.ServiceAccountNamespace)
	d.Set("url", in.Spec.URL)

	return nil
}

// Expanders

func expandCatalogV2(in *schema.ResourceData) (*ClusterRepo, error) {
	if in == nil {
		return nil, nil
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
		// The rancher API expects the CA bundle to be in base64-encoded DER
		// format. The caBundle field of ClusterRepo is of type []byte, so
		// json.Marshal base64-encodes this field for us. So internally, we
		// set the caBundle field to non-base64-encoded DER.
		decodedCABundle, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return nil, fmt.Errorf("failed to base64-decode ca_bundle: %w", err)
		}
		obj.Spec.CABundle = decodedCABundle
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
	if v, ok := in.Get("insecure_plain_http").(bool); ok {
		obj.Spec.InsecurePlainHTTP = v
	}
	if v, ok := in.Get("exponential_backoff_min_wait").(int); ok {
		if obj.Spec.ExponentialBackOffValues != nil {
			obj.Spec.ExponentialBackOffValues.MinWait = v
		} else {
			obj.Spec.ExponentialBackOffValues = &v1.ExponentialBackOffValues{
				MinWait: v,
			}
		}
	}

	if v, ok := in.Get("exponential_backoff_max_wait").(int); ok {
		if obj.Spec.ExponentialBackOffValues != nil {
			obj.Spec.ExponentialBackOffValues.MaxWait = v
		} else {
			obj.Spec.ExponentialBackOffValues = &v1.ExponentialBackOffValues{
				MaxWait: v,
			}
		}
	}

	if v, ok := in.Get("exponential_backoff_max_retries").(int); ok {
		if obj.Spec.ExponentialBackOffValues != nil {
			obj.Spec.ExponentialBackOffValues.MaxRetries = v
		} else {
			obj.Spec.ExponentialBackOffValues = &v1.ExponentialBackOffValues{
				MaxRetries: v,
			}
		}
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

	return obj, nil
}
