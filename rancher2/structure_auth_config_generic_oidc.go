package rancher2

import (
	"fmt"

	"github.com/go-viper/mapstructure/v2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenAuthConfigGenericOIDC(d *schema.ResourceData, in *managementClient.GenericOIDCConfig) error {
	d.SetId(AuthConfigGenericOIDCName)
	d.Set("name", AuthConfigGenericOIDCName)
	d.Set("type", managementClient.GenericOIDCConfigType)

	if err := flattenOIDCConfig(d, in); err != nil {
		return fmt.Errorf("flattening AuthConfig for GenericOIDC: %s", err)
	}

	return nil
}

// Expanders

func expandAuthConfigGenericOIDC(in *schema.ResourceData) (*managementClient.GenericOIDCConfig, error) {
	obj := &managementClient.GenericOIDCConfig{}
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

	if v, ok := in.Get("name_claim").(string); ok && len(v) > 0 {
		obj.NameClaim = v
	}

	if v, ok := in.Get("email_claim").(string); ok && len(v) > 0 {
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

// flattenOIDCConfig is a generic OIDC flattener.
//
// It converts the provided input to a map and looks up known keys in the map.
func flattenOIDCConfig(d *schema.ResourceData, in any) error {
	var oidcData map[string]any
	if err := mapstructure.Decode(in, &oidcData); err != nil {
		return fmt.Errorf("decoding struct: %w", err)
	}

	d.Set("access_mode", oidcData["AccessMode"].(string))

	err := d.Set("allowed_principal_ids", oidcData["AllowedPrincipalIDs"])
	if err != nil {
		return err
	}

	d.Set("enabled", oidcData["Enabled"].(bool))

	err = d.Set("annotations", oidcData["Annotations"])
	if err != nil {
		return err
	}

	err = d.Set("labels", oidcData["Labels"])
	if err != nil {
		return err
	}

	d.Set("client_id", oidcData["ClientID"])
	d.Set("issuer", oidcData["Issuer"])
	d.Set("rancher_url", oidcData["RancherURL"])
	d.Set("auth_endpoint", oidcData["AuthEndpoint"])
	d.Set("token_endpoint", oidcData["TokenEndpoint"])
	d.Set("userinfo_endpoint", oidcData["UserInfoEndpoint"])
	d.Set("jwks_url", oidcData["JWKSUrl"])
	d.Set("scopes", oidcData["Scopes"])

	groupSearchEnabled := oidcData["GroupSearchEnabled"]
	if groupSearchEnabled != nil {
		d.Set("group_search_enabled", groupSearchEnabled)
	}
	d.Set("groups_field", oidcData["GroupsClaim"])
	d.Set("certificate", oidcData["Certificate"])
	d.Set("private_key", oidcData["PrivateKey"])

	d.Set("name_claim", oidcData["NameClaim"])
	d.Set("email_claim", oidcData["EmailClaim"])

	d.Set("logout_all_enabled", oidcData["LogoutAllEnabled"])
	d.Set("logout_all_forced", oidcData["LogoutAllForced"])

	if v, ok := oidcData["EndSessionEndpoint"]; ok {
		d.Set("end_session_endpoint", v)
	}

	return nil
}

// expandOIDCConfig is a generic OIDC flattener.
//
// This is a generic function for configuring AuthConfigs based on the
// GenericOIDC fields.
func expandOIDCConfig[T *Q, Q any](in *schema.ResourceData, obj T) (T, error) {
	var objData map[string]any
	if err := mapstructure.Decode(obj, &objData); err != nil {
		return obj, fmt.Errorf("decoding struct: %w", err)
	}

	if v, ok := in.Get("access_mode").(string); ok && v != "" {
		objData["AccessMode"] = v
	}

	if v, ok := in.Get("allowed_principal_ids").([]any); ok && len(v) > 0 {
		objData["AllowedPrincipalIDs"] = toArrayString(v)
	}

	if (objData["AccessMode"] == "required" || objData["AccessMode"] == "restricted") && len(objData["AllowedPrincipalIDs"].([]string)) == 0 {
		return nil, fmt.Errorf("expanding %s Auth Config: allowed_principal_ids is required on access_mode %s", objData["Name"], objData["AccessMode"])
	}

	if v, ok := in.Get("enabled").(bool); ok {
		objData["Enabled"] = v
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		objData["Annotations"] = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		objData["Labels"] = toMapString(v)
	}

	if v, ok := in.Get("client_id").(string); ok && v != "" {
		objData["ClientID"] = v
	}

	if v, ok := in.Get("client_secret").(string); ok && len(v) > 0 {
		objData["ClientSecret"] = v
	}

	if v, ok := in.Get("issuer").(string); ok && v != "" {
		objData["Issuer"] = v
	}

	if v, ok := in.Get("rancher_url").(string); ok && v != "" {
		objData["RancherURL"] = v
	}

	if v, ok := in.Get("auth_endpoint").(string); ok && v != "" {
		objData["AuthEndpoint"] = v
	}

	if v, ok := in.Get("token_endpoint").(string); ok && v != "" {
		objData["TokenEndpoint"] = v
	}

	if v, ok := in.Get("userinfo_endpoint").(string); ok && v != "" {
		objData["UserInfoEndpoint"] = v
	}

	if v, ok := in.Get("jwks_url").(string); ok && v != "" {
		objData["JWKSUrl"] = v
	}

	if v, ok := in.Get("scopes").(string); ok && v != "" {
		objData["Scopes"] = v
	}

	if v, ok := in.Get("group_search_enabled").(bool); ok {
		objData["GroupSearchEnabled"] = &v
	}

	if v, ok := in.Get("groups_field").(string); ok && v != "" {
		objData["GroupsClaim"] = v
	}

	if v, ok := in.Get("certificate").(string); ok && v != "" {
		objData["Certificate"] = v
	}

	if v, ok := in.Get("private_key").(string); ok && v != "" {
		objData["PrivateKey"] = v
	}

	if v, ok := in.Get("name_claim").(string); ok && v != "" {
		objData["NameClaim"] = v
	}

	if v, ok := in.Get("email_claim").(string); ok && v != "" {
		objData["EmailClaim"] = v
	}

	if v, ok := in.Get("logout_all_enabled").(bool); ok {
		objData["LogoutAllEnabled"] = v
	}

	if v, ok := in.Get("logout_all_forced").(bool); ok {
		objData["LogoutAllForced"] = v
	}

	if v, ok := in.Get("end_session_endpoint").(string); ok && v != "" {
		objData["EndSessionEndpoint"] = v
	}

	if err := mapstructure.Decode(objData, obj); err != nil {
		return obj, fmt.Errorf("decoding struct: %w", err)
	}

	return obj, nil
}
