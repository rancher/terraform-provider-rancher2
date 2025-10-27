package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenAuthConfigGenericOIDC(d *schema.ResourceData, in *managementClient.OIDCConfig) error {
	d.SetId(AuthConfigGenericOIDCName)
	d.Set("name", AuthConfigGenericOIDCName)
	d.Set("type", managementClient.GenericOIDCConfigType)
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
	d.Set("issuer", in.Issuer)
	d.Set("rancher_url", in.RancherURL)
	d.Set("auth_endpoint", in.AuthEndpoint)
	d.Set("token_endpoint", in.TokenEndpoint)
	d.Set("userinfo_endpoint", in.UserInfoEndpoint)
	d.Set("jwks_url", in.JWKSUrl)
	d.Set("scopes", in.Scopes)
	if in.GroupSearchEnabled != nil {
		d.Set("group_search_enabled", *in.GroupSearchEnabled)
	}
	d.Set("groups_field", in.GroupsClaim)
	d.Set("certificate", in.Certificate)
	d.Set("private_key", in.PrivateKey)

	return nil
}

// Expanders

func expandAuthConfigGenericOIDC(in *schema.ResourceData) (*managementClient.OIDCConfig, error) {
	obj := &managementClient.OIDCConfig{}
	if in == nil {
		return nil, fmt.Errorf("expanding %s Auth Config: Input ResourceData is nil", AuthConfigGenericOIDCName)
	}

	obj.Name = AuthConfigGenericOIDCName
	obj.Type = managementClient.GenericOIDCConfigType

	if v, ok := in.Get("access_mode").(string); ok && len(v) > 0 {
		obj.AccessMode = v
	}

	if v, ok := in.Get("allowed_principal_ids").([]interface{}); ok && len(v) > 0 {
		obj.AllowedPrincipalIDs = toArrayString(v)
	}

	if (obj.AccessMode == "required" || obj.AccessMode == "restricted") && len(obj.AllowedPrincipalIDs) == 0 {
		return nil, fmt.Errorf("expanding %s Auth Config: allowed_principal_ids is required on access_mode %s", AuthConfigGenericOIDCName, obj.AccessMode)
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

	if v, ok := in.Get("issuer").(string); ok && len(v) > 0 {
		obj.Issuer = v
	}

	if v, ok := in.Get("rancher_url").(string); ok && len(v) > 0 {
		obj.RancherURL = v
	}

	if v, ok := in.Get("auth_endpoint").(string); ok && len(v) > 0 {
		obj.AuthEndpoint = v
	}

	if v, ok := in.Get("token_endpoint").(string); ok && len(v) > 0 {
		obj.TokenEndpoint = v
	}

	if v, ok := in.Get("userinfo_endpoint").(string); ok && len(v) > 0 {
		obj.UserInfoEndpoint = v
	}

	if v, ok := in.Get("jwks_url").(string); ok && len(v) > 0 {
		obj.JWKSUrl = v
	}

	if v, ok := in.Get("scopes").(string); ok && len(v) > 0 {
		obj.Scopes = v
	}

	if v, ok := in.Get("group_search_enabled").(bool); ok {
		obj.GroupSearchEnabled = &v
	}

	if v, ok := in.Get("groups_field").(string); ok && len(v) > 0 {
		obj.GroupsClaim = v
	}

	if v, ok := in.Get("certificate").(string); ok && len(v) > 0 {
		obj.Certificate = v
	}

	if v, ok := in.Get("private_key").(string); ok && len(v) > 0 {
		obj.PrivateKey = v
	}

	return obj, nil
}
