package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenAuthConfigAzureAD(d *schema.ResourceData, in *managementClient.AzureADConfig) error {
	d.SetId(AuthConfigAzureADName)
	d.Set("name", AuthConfigAzureADName)
	d.Set("type", managementClient.AzureADConfigType)
	d.Set("access_mode", in.AccessMode)

	err := d.Set("allowed_principal_ids", toArrayInterface(in.AllowedPrincipalIDs))
	if err != nil {
		return err
	}

	d.Set("enabled", in.Enabled)

	err = d.Set("annotations", toMapInterface(in.Annotations))
	if err != nil {
		return err
	}
	err = d.Set("labels", toMapInterface(in.Labels))
	if err != nil {
		return err
	}

	d.Set("application_id", in.ApplicationID)
	d.Set("auth_endpoint", in.AuthEndpoint)
	d.Set("endpoint", in.Endpoint)
	d.Set("graph_endpoint", in.GraphEndpoint)
	d.Set("rancher_url", in.RancherURL)
	d.Set("tenant_id", in.TenantID)
	d.Set("token_endpoint", in.TokenEndpoint)

	return nil
}

// Expanders

func expandAuthConfigAzureAD(in *schema.ResourceData) (*managementClient.AzureADConfig, error) {
	obj := &managementClient.AzureADConfig{}
	if in == nil {
		return nil, fmt.Errorf("expanding %s Auth Config: Input ResourceData is nil", AuthConfigAzureADName)
	}

	obj.Name = AuthConfigAzureADName
	obj.Type = managementClient.AzureADConfigType

	if v, ok := in.Get("access_mode").(string); ok && len(v) > 0 {
		obj.AccessMode = v
	}

	if v, ok := in.Get("allowed_principal_ids").([]interface{}); ok && len(v) > 0 {
		obj.AllowedPrincipalIDs = toArrayString(v)
	}

	if (obj.AccessMode == "required" || obj.AccessMode == "restricted") && len(obj.AllowedPrincipalIDs) == 0 {
		return nil, fmt.Errorf("expanding %s Auth Config: allowed_principal_ids is required on access_mode %s", AuthConfigAzureADName, obj.AccessMode)
	}

	if v, ok := in.Get("enabled").(bool); ok {
		obj.Enabled = v
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	if v, ok := in.Get("application_id").(string); ok && len(v) > 0 {
		obj.ApplicationID = v
	}

	if v, ok := in.Get("application_secret").(string); ok && len(v) > 0 {
		obj.ApplicationSecret = v
	}

	if v, ok := in.Get("auth_endpoint").(string); ok && len(v) > 0 {
		obj.AuthEndpoint = v
	}

	if v, ok := in.Get("endpoint").(string); ok {
		obj.Endpoint = v
	}

	if v, ok := in.Get("graph_endpoint").(string); ok && len(v) > 0 {
		obj.GraphEndpoint = v
	}

	if v, ok := in.Get("rancher_url").(string); ok && len(v) > 0 {
		obj.RancherURL = v
	}

	if v, ok := in.Get("tenant_id").(string); ok && len(v) > 0 {
		obj.TenantID = v
	}

	if v, ok := in.Get("token_endpoint").(string); ok {
		obj.TokenEndpoint = v
	}

	return obj, nil
}
