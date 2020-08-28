package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenAuthConfigLdap(d *schema.ResourceData, in *managementClient.LdapConfig) error {
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
	err = d.Set("servers", toArrayInterface(in.Servers))
	if err != nil {
		return err
	}

	d.Set("service_account_distinguished_name", in.ServiceAccountDistinguishedName)
	d.Set("user_search_base", in.UserSearchBase)
	d.Set("certificate", Base64Encode(in.Certificate))
	d.Set("connection_timeout", int(in.ConnectionTimeout))
	d.Set("group_dn_attribute", in.GroupDNAttribute)
	d.Set("group_member_mapping_attribute", in.GroupMemberMappingAttribute)
	d.Set("group_member_user_attribute", in.GroupMemberUserAttribute)
	d.Set("group_name_attribute", in.GroupNameAttribute)
	d.Set("group_object_class", in.GroupObjectClass)
	d.Set("group_search_attribute", in.GroupSearchAttribute)
	d.Set("group_search_base", in.GroupSearchBase)
	d.Set("nested_group_membership_enabled", in.NestedGroupMembershipEnabled)
	d.Set("port", int(in.Port))
	d.Set("tls", in.TLS)
	d.Set("user_disabled_bit_mask", int(in.UserDisabledBitMask))
	d.Set("user_enabled_attribute", in.UserEnabledAttribute)
	d.Set("user_login_attribute", in.UserLoginAttribute)
	d.Set("user_member_attribute", in.UserMemberAttribute)
	d.Set("user_name_attribute", in.UserNameAttribute)
	d.Set("user_object_class", in.UserObjectClass)
	d.Set("user_search_attribute", in.UserSearchAttribute)

	return nil
}

// Expanders

func expandAuthConfigLdap(in *schema.ResourceData) (*managementClient.LdapConfig, error) {
	obj := &managementClient.LdapConfig{}
	if in == nil {
		return nil, fmt.Errorf("expanding ldap Auth Config: Input ResourceData is nil")
	}

	if v, ok := in.Get("access_mode").(string); ok && len(v) > 0 {
		obj.AccessMode = v
	}

	if v, ok := in.Get("allowed_principal_ids").([]interface{}); ok && len(v) > 0 {
		obj.AllowedPrincipalIDs = toArrayString(v)
	}

	if (obj.AccessMode == "required" || obj.AccessMode == "restricted") && len(obj.AllowedPrincipalIDs) == 0 {
		return nil, fmt.Errorf("expanding ldap Auth Config: allowed_principal_ids is required on access_mode %s", obj.AccessMode)
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

	if v, ok := in.Get("servers").([]interface{}); ok && len(v) > 0 {
		obj.Servers = toArrayString(v)
	}

	if v, ok := in.Get("service_account_distinguished_name").(string); ok && len(v) > 0 {
		obj.ServiceAccountDistinguishedName = v
	}

	if v, ok := in.Get("service_account_password").(string); ok && len(v) > 0 {
		obj.ServiceAccountPassword = v
	}

	if v, ok := in.Get("user_search_base").(string); ok && len(v) > 0 {
		obj.UserSearchBase = v
	}

	if v, ok := in.Get("certificate").(string); ok && len(v) > 0 {
		cert, err := Base64Decode(v)
		if err != nil {
			return nil, fmt.Errorf("expanding ldap Auth Config: certificate is not base64 encoded: %s", v)
		}
		obj.Certificate = cert
	}

	if v, ok := in.Get("connection_timeout").(int); ok && v > 0 {
		obj.ConnectionTimeout = int64(v)
	}

	if v, ok := in.Get("group_dn_attribute").(string); ok && len(v) > 0 {
		obj.GroupDNAttribute = v
	}

	if v, ok := in.Get("group_member_mapping_attribute").(string); ok && len(v) > 0 {
		obj.GroupMemberMappingAttribute = v
	}

	if v, ok := in.Get("group_member_user_attribute").(string); ok && len(v) > 0 {
		obj.GroupMemberUserAttribute = v
	}

	if v, ok := in.Get("group_name_attribute").(string); ok && len(v) > 0 {
		obj.GroupNameAttribute = v
	}

	if v, ok := in.Get("group_object_class").(string); ok && len(v) > 0 {
		obj.GroupObjectClass = v
	}

	if v, ok := in.Get("group_search_attribute").(string); ok && len(v) > 0 {
		obj.GroupSearchAttribute = v
	}

	if v, ok := in.Get("group_search_base").(string); ok && len(v) > 0 {
		obj.GroupSearchBase = v
	}

	if v, ok := in.Get("nested_group_membership_enabled").(bool); ok {
		obj.NestedGroupMembershipEnabled = v
	}

	if v, ok := in.Get("port").(int); ok && v > 0 {
		obj.Port = int64(v)
	}

	if v, ok := in.Get("tls").(bool); ok {
		obj.TLS = v
	}

	if v, ok := in.Get("user_disabled_bit_mask").(int); ok && v > 0 {
		obj.UserDisabledBitMask = int64(v)
	}

	if v, ok := in.Get("user_enabled_attribute").(string); ok && len(v) > 0 {
		obj.UserEnabledAttribute = v
	}

	if v, ok := in.Get("user_login_attribute").(string); ok && len(v) > 0 {
		obj.UserLoginAttribute = v
	}

	if v, ok := in.Get("user_member_attribute").(string); ok && len(v) > 0 {
		obj.UserMemberAttribute = v
	}

	if v, ok := in.Get("user_name_attribute").(string); ok && len(v) > 0 {
		obj.UserNameAttribute = v
	}

	if v, ok := in.Get("user_object_class").(string); ok && len(v) > 0 {
		obj.UserObjectClass = v
	}

	if v, ok := in.Get("user_search_attribute").(string); ok && len(v) > 0 {
		obj.UserSearchAttribute = v
	}

	return obj, nil
}
