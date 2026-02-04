package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenAuthConfigCognito(d *schema.ResourceData, in *managementClient.CognitoConfig) error {
	d.SetId(AuthConfigCognitoName)
	d.Set("name", AuthConfigCognitoName)
	d.Set("type", managementClient.CognitoConfigType)

	if err := flattenOIDCConfig(d, in); err != nil {
		return fmt.Errorf("flattening AuthConfig for Cognito: %s", err)
	}

	return nil
}

// Expanders

func expandAuthConfigCognito(in *schema.ResourceData) (*managementClient.CognitoConfig, error) {
	obj := &managementClient.CognitoConfig{}
	if in == nil {
		return nil, fmt.Errorf("expanding :%s Auth Config: Input ResourceData is nil", AuthConfigCognitoName)
	}

	obj.Name = AuthConfigCognitoName
	obj.Type = managementClient.CognitoConfigType

	if v, ok := in.Get("access_mode").(string); ok && len(v) > 0 {
		obj.AccessMode = v
	}

	if v, ok := in.Get("allowed_principal_ids").([]interface{}); ok && len(v) > 0 {
		obj.AllowedPrincipalIDs = toArrayString(v)
	}

	if (obj.AccessMode == "required" || obj.AccessMode == "restricted") && len(obj.AllowedPrincipalIDs) == 0 {
		return nil, fmt.Errorf("expanding %s Auth Config: allowed_principal_ids is required on access_mode %s", AuthConfigCognitoName, obj.AccessMode)
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

	if v, ok := in.Get("name_claim").(string); ok && len(v) > 0 {
		obj.NameClaim = v
	}

	if v, ok := in.Get("email_claim").(string); ok && len(v) > 0 {
		obj.EmailClaim = v
	}

	return obj, nil
}
