package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenAuthConfigActiveDirectory(d *schema.ResourceData, in *managementClient.ActiveDirectoryConfig) error {
	d.SetId(AuthConfigActiveDirectoryName)
	d.Set("name", AuthConfigActiveDirectoryName)
	d.Set("type", managementClient.ActiveDirectoryConfigType)
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
	d.Set("service_account_username", in.ServiceAccountUsername)
	d.Set("user_search_base", in.UserSearchBase)
	d.Set("certificate", in.Certificate)
	d.Set("connection_timeout", int(in.ConnectionTimeout))
	d.Set("default_login_domain", in.DefaultLoginDomain)
	d.Set("group_dn_attribute", in.GroupDNAttribute)
	d.Set("group_member_mapping_attribute", in.GroupMemberMappingAttribute)
	d.Set("group_member_user_attribute", in.GroupMemberUserAttribute)
	d.Set("group_name_attribute", in.GroupNameAttribute)
	d.Set("group_object_class", in.GroupObjectClass)
	d.Set("group_search_attribute", in.GroupSearchAttribute)
	d.Set("group_search_base", in.GroupSearchBase)
	d.Set("group_search_filter", in.GroupSearchFilter)
	d.Set("nested_group_membership_enabled", *in.NestedGroupMembershipEnabled)
	d.Set("port", int(in.Port))
	d.Set("start_tls", in.StartTLS)
	d.Set("tls", in.TLS)
	d.Set("user_disabled_bit_mask", int(in.UserDisabledBitMask))
	d.Set("user_enabled_attribute", in.UserEnabledAttribute)
	d.Set("user_login_attribute", in.UserLoginAttribute)
	d.Set("user_name_attribute", in.UserNameAttribute)
	d.Set("user_object_class", in.UserObjectClass)
	d.Set("user_search_attribute", in.UserSearchAttribute)
	d.Set("user_search_filter", in.UserSearchFilter)

	return nil
}

// Expanders

func expandAuthConfigActiveDirectory(in *schema.ResourceData) (*managementClient.ActiveDirectoryConfig, error) {
	obj := &managementClient.ActiveDirectoryConfig{}
	if in == nil {
		return nil, fmt.Errorf("expanding %s Auth Config: Input ResourceData is nil", AuthConfigActiveDirectoryName)
	}

	obj.Name = AuthConfigActiveDirectoryName
	obj.Type = managementClient.ActiveDirectoryConfigType

	if v, ok := in.Get("access_mode").(string); ok && len(v) > 0 {
		obj.AccessMode = v
	}

	if v, ok := in.Get("allowed_principal_ids").([]interface{}); ok && len(v) > 0 {
		obj.AllowedPrincipalIDs = toArrayString(v)
	}

	if (obj.AccessMode == "required" || obj.AccessMode == "restricted") && len(obj.AllowedPrincipalIDs) == 0 {
		return nil, fmt.Errorf("expanding %s Auth Config: allowed_principal_ids is required on access_mode %s", AuthConfigActiveDirectoryName, obj.AccessMode)
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

	if v, ok := in.Get("service_account_username").(string); ok && len(v) > 0 {
		obj.ServiceAccountUsername = v
	}

	if v, ok := in.Get("service_account_password").(string); ok && len(v) > 0 {
		obj.ServiceAccountPassword = v
	}

	if v, ok := in.Get("user_search_base").(string); ok && len(v) > 0 {
		obj.UserSearchBase = v
	}

	if v, ok := in.Get("certificate").(string); ok && len(v) > 0 {
		obj.Certificate = v
	}

	if v, ok := in.Get("connection_timeout").(int); ok && v > 0 {
		obj.ConnectionTimeout = int64(v)
	}

	if v, ok := in.Get("default_login_domain").(string); ok && len(v) > 0 {
		obj.DefaultLoginDomain = v
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

	if v, ok := in.Get("group_search_filter").(string); ok && len(v) > 0 {
		obj.GroupSearchFilter = v
	}

	if v, ok := in.Get("nested_group_membership_enabled").(bool); ok {
		obj.NestedGroupMembershipEnabled = &v
	}

	if v, ok := in.Get("port").(int); ok && v > 0 {
		obj.Port = int64(v)
	}

	if v, ok := in.Get("start_tls").(bool); ok {
		obj.StartTLS = v
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

	if v, ok := in.Get("user_name_attribute").(string); ok && len(v) > 0 {
		obj.UserNameAttribute = v
	}

	if v, ok := in.Get("user_object_class").(string); ok && len(v) > 0 {
		obj.UserObjectClass = v
	}

	if v, ok := in.Get("user_search_attribute").(string); ok && len(v) > 0 {
		obj.UserSearchAttribute = v
	}

	if v, ok := in.Get("user_search_filter").(string); ok && len(v) > 0 {
		obj.UserSearchFilter = v
	}

	return obj, nil
}
