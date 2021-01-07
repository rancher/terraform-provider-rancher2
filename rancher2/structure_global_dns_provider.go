package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenAliDNSProvider(in *managementClient.AlidnsProviderConfig, p []interface{}) []interface{} {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}
	if in == nil {
		return []interface{}{}
	}

	if len(in.AccessKey) > 0 {
		obj["access_key"] = in.AccessKey
	}

	return []interface{}{obj}
}

func flattenCloudFlareDNSProvider(in *managementClient.CloudflareProviderConfig, p []interface{}) []interface{} {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}
	if in == nil {
		return []interface{}{}
	}

	if len(in.APIEmail) > 0 {
		obj["api_email"] = in.APIEmail
	}

	if in.ProxySetting != nil {
		obj["proxy_setting"] = *in.ProxySetting
	}

	return []interface{}{obj}
}

func flattenRoute53DNSProvider(in *managementClient.Route53ProviderConfig, p []interface{}) []interface{} {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}
	if in == nil {
		return []interface{}{}
	}

	if len(in.AccessKey) > 0 {
		obj["access_key"] = in.AccessKey
	}
	if len(in.CredentialsPath) > 0 {
		obj["credentials_path"] = in.CredentialsPath
	}
	if len(in.Region) > 0 {
		obj["region"] = in.Region
	}
	if len(in.RoleArn) > 0 {
		obj["role_arn"] = in.RoleArn
	}
	if len(in.ZoneType) > 0 {
		obj["zone_type"] = in.ZoneType
	}

	return []interface{}{obj}
}

func flattenGlobalDNSProvider(d *schema.ResourceData, in *managementClient.GlobalDnsProvider) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)
	d.Set("root_domain", in.RootDomain)
	d.Set("name", in.Name)

	if in.AlidnsProviderConfig != nil {
		err := d.Set("alidns_config", flattenAliDNSProvider(in.AlidnsProviderConfig, d.Get("alidns_config").([]interface{})))
		if err != nil {
			return err
		}
		d.Set("dns_provider", globalDNSProviderAlidnsKind)
	}
	if in.CloudflareProviderConfig != nil {
		err := d.Set("cloudflare_config", flattenCloudFlareDNSProvider(in.CloudflareProviderConfig, d.Get("cloudflare_config").([]interface{})))
		if err != nil {
			return err
		}
		d.Set("dns_provider", globalDNSProviderCloudflareKind)
	}
	if in.Route53ProviderConfig != nil {
		err := d.Set("route53_config", flattenRoute53DNSProvider(in.Route53ProviderConfig, d.Get("route53_config").([]interface{})))
		if err != nil {
			return err
		}
		d.Set("dns_provider", globalDNSProviderRoute53Kind)
	}

	err := d.Set("annotations", toMapInterface(in.Annotations))
	if err != nil {
		return err
	}

	err = d.Set("labels", toMapInterface(in.Labels))
	if err != nil {
		return err
	}

	return nil
}

// Expanders

func expandAliDNSConfig(p []interface{}) *managementClient.AlidnsProviderConfig {
	obj := &managementClient.AlidnsProviderConfig{}
	if len(p) == 0 || p[0] == nil {
		return nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["access_key"].(string); ok && len(v) > 0 {
		obj.AccessKey = v
	}

	if v, ok := in["secret_key"].(string); ok && len(v) > 0 {
		obj.SecretKey = v
	}

	return obj
}

func expandCloudFlareDNSConfig(p []interface{}) *managementClient.CloudflareProviderConfig {
	obj := &managementClient.CloudflareProviderConfig{}
	if len(p) == 0 || p[0] == nil {
		return nil
	}
	in := p[0].(map[string]interface{})

	obj.APIEmail = in["api_email"].(string)
	obj.APIKey = in["api_key"].(string)
	proxySetting := in["proxy_setting"].(bool)
	obj.ProxySetting = &proxySetting

	return obj
}

func expandRoute53DNSConfig(p []interface{}) *managementClient.Route53ProviderConfig {
	obj := &managementClient.Route53ProviderConfig{}
	if len(p) == 0 || p[0] == nil {
		return nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["access_key"].(string); ok && len(v) > 0 {
		obj.AccessKey = v
	}

	if v, ok := in["secret_key"].(string); ok && len(v) > 0 {
		obj.SecretKey = v
	}

	if v, ok := in["credentials_path"].(string); ok && len(v) > 0 {
		obj.CredentialsPath = v
	}

	if v, ok := in["zone_type"].(string); ok {
		obj.ZoneType = v
	}

	if v, ok := in["role_arn"].(string); ok {
		obj.RoleArn = v
	}

	if v, ok := in["region"].(string); ok {
		obj.Region = v
	}

	return obj

}

func expandGlobalDNSProvider(in *schema.ResourceData) *managementClient.GlobalDnsProvider {
	obj := &managementClient.GlobalDnsProvider{}
	if in == nil {
		return nil
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	obj.Name = in.Get("name").(string)
	obj.RootDomain = in.Get("root_domain").(string)

	if v, ok := in.Get("alidns_config").([]interface{}); ok && len(v) > 0 {
		obj.AlidnsProviderConfig = expandAliDNSConfig(v)
		in.Set("dns_provider", globalDNSProviderAlidnsKind)
	}
	if v, ok := in.Get("cloudflare_config").([]interface{}); ok && len(v) > 0 {
		obj.CloudflareProviderConfig = expandCloudFlareDNSConfig(v)
		in.Set("dns_provider", globalDNSProviderCloudflareKind)
	}
	if v, ok := in.Get("route53_config").([]interface{}); ok && len(v) > 0 {
		obj.Route53ProviderConfig = expandRoute53DNSConfig(v)
		in.Set("dns_provider", globalDNSProviderRoute53Kind)
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj
}
