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

	if v, ok := in.Get("access_mode").(string); ok && v != "" {
		obj.AccessMode = v
	}

	if v, ok := in.Get("allowed_principal_ids").([]any); ok && len(v) > 0 {
		obj.AllowedPrincipalIDs = toArrayString(v)
	}

	if (obj.AccessMode == "required" || obj.AccessMode == "restricted") && len(obj.AllowedPrincipalIDs) == 0 {
		return nil, fmt.Errorf("expanding %s Auth Config: allowed_principal_ids is required on access_mode %s", AuthConfigCognitoName, obj.AccessMode)
	}

	if v, ok := in.Get("enabled").(bool); ok {
		obj.Enabled = v
	}

	if v, ok := in.Get("annotations").(map[string]any); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]any); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	if v, ok := in.Get("client_id").(string); ok && v != "" {
		obj.ClientID = v
	}

	if v, ok := in.Get("client_secret").(string); ok && v != "" {
		obj.ClientSecret = v
	}

	if v, ok := in.Get("issuer").(string); ok && v != "" {
		obj.Issuer = v
	}

	if v, ok := in.Get("rancher_url").(string); ok && v != "" {
		obj.RancherURL = v
	}

	if v, ok := in.Get("auth_endpoint").(string); ok && v != "" {
		obj.AuthEndpoint = v
	}

	if v, ok := in.Get("token_endpoint").(string); ok && v != "" {
		obj.TokenEndpoint = v
	}

	if v, ok := in.Get("userinfo_endpoint").(string); ok && v != "" {
		obj.UserInfoEndpoint = v
	}

	if v, ok := in.Get("jwks_url").(string); ok && v != "" {
		obj.JWKSUrl = v
	}

	if v, ok := in.Get("scopes").(string); ok && v != "" {
		obj.Scopes = v
	}

	if v, ok := in.Get("group_search_enabled").(bool); ok {
		obj.GroupSearchEnabled = &v
	}

	if v, ok := in.Get("groups_field").(string); ok && v != "" {
		obj.GroupsClaim = v
	}

	if v, ok := in.Get("certificate").(string); ok && v != "" {
		obj.Certificate = v
	}

	if v, ok := in.Get("private_key").(string); ok && v != "" {
		obj.PrivateKey = v
	}

	if v, ok := in.Get("name_claim").(string); ok && v != "" {
		obj.NameClaim = v
	}

	if v, ok := in.Get("email_claim").(string); ok && v != "" {
		obj.EmailClaim = v
	}

	if v, ok := in.Get("logout_all_enabled").(bool); ok {
		obj.LogoutAllEnabled = v
	}

	if v, ok := in.Get("logout_all_forced").(bool); ok {
		obj.LogoutAllForced = v
	}

	if v, ok := in.Get("end_session_endpoint").(string); ok && v != "" {
		obj.EndSessionEndpoint = v
	}

	return obj, nil
}
