package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenAuthConfigGithub(d *schema.ResourceData, in *managementClient.GithubConfig) error {
	d.SetId(AuthConfigGithubName)
	d.Set("name", AuthConfigGithubName)
	d.Set("type", managementClient.GithubConfigType)
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

	d.Set("client_id", in.ClientID)
	d.Set("hostname", in.Hostname)
	d.Set("tls", in.TLS)

	return nil
}

// Expanders

func expandAuthConfigGithub(in *schema.ResourceData) (*managementClient.GithubConfig, error) {
	obj := &managementClient.GithubConfig{}
	if in == nil {
		return nil, fmt.Errorf("expanding %s Auth Config: Input ResourceData is nil", AuthConfigGithubName)
	}

	obj.Name = AuthConfigGithubName
	obj.Type = managementClient.GithubConfigType

	if v, ok := in.Get("access_mode").(string); ok && len(v) > 0 {
		obj.AccessMode = v
	}

	if v, ok := in.Get("allowed_principal_ids").([]interface{}); ok && len(v) > 0 {
		obj.AllowedPrincipalIDs = toArrayString(v)
	}

	if (obj.AccessMode == "required" || obj.AccessMode == "restricted") && len(obj.AllowedPrincipalIDs) == 0 {
		return nil, fmt.Errorf("expanding %s Auth Config: allowed_principal_ids is required on access_mode %s", AuthConfigGithubName, obj.AccessMode)
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

	if v, ok := in.Get("client_id").(string); ok && len(v) > 0 {
		obj.ClientID = v
	}

	if v, ok := in.Get("client_secret").(string); ok && len(v) > 0 {
		obj.ClientSecret = v
	}

	if v, ok := in.Get("hostname").(string); ok && len(v) > 0 {
		obj.Hostname = v
	}

	if v, ok := in.Get("tls").(bool); ok {
		obj.TLS = v
	}

	return obj, nil
}
