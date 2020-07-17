package rancher2

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

// Flatteners
func flattenRoute53DNSProvider(in *managementClient.Route53ProviderConfig) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.Region) > 0 {
		obj["region"] = in.Region
	}

	if len(in.SecretKey) > 0 {
		obj["secret_key"] = in.SecretKey
	}

	if len(in.AccessKey) > 0 {
		obj["access_key"] = in.AccessKey
	}

	if len(in.RoleArn) > 0 {
		obj["role_arn"] = in.RoleArn
	}

	if len(in.CredentialsPath) > 0 {
		obj["credentials_path"] = in.CredentialsPath
	}

	return []interface{}{obj}, nil
}

func flattenAliDNSProvider(in *managementClient.AlidnsProviderConfig) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.SecretKey) > 0 {
		obj["secret_key"] = in.SecretKey
	}

	if len(in.AccessKey) > 0 {
		obj["access_key"] = in.AccessKey
	}

	return []interface{}{obj}, nil
}

func flattenCloudFlareDNSProvider(in *managementClient.CloudflareProviderConfig) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.APIKey) > 0 {
		obj["api_key"] = in.APIKey
	}

	if len(in.APIEmail) > 0 {
		obj["api_email"] = in.APIEmail
	}

	obj["proxy_setting"] = in.ProxySetting

	return []interface{}{obj}, nil
}

func flattenGlobalDNSProvider(d *schema.ResourceData, in *managementClient.GlobalDNSProvider) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)
	d.Set("root_domain", in.RootDomain)
	d.Set("name", in.Name)

	switch d.Get("dns_provider") {
	case "route53":
		route53Config, err := flattenRoute53DNSProvider(in.Route53ProviderConfig)
		if err != nil {
			return err
		}
		err = d.Set("route53_config", route53Config)
		if err != nil {
			return err
		}
	case "alidns":
		aliDNSConfig, err := flattenAliDNSProvider(in.AlidnsProviderConfig)
		if err != nil {
			return err
		}
		err = d.Set("alidns_config", aliDNSConfig)
		if err != nil {
			return err
		}
	case "cloudfare":
		cfDNS, err := flattenCloudFlareDNSProvider(in.CloudflareProviderConfig)
		if err != nil {
			return err
		}
		err = d.Set("cloudflare_config", cfDNS)
	}

	return nil

}

// Expanders
func expandGlobalDNSProvider(in *schema.ResourceData) (*managementClient.GlobalDNSProvider, error) {
	obj := &managementClient.GlobalDNSProvider{}
	if in == nil {
		return nil, fmt.Errorf("resource rancher2_global_dns_provider data cannot be nil")
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	obj.Name = in.Get("name").(string)
	obj.RootDomain = in.Get("root_domain").(string)

	switch in.Get("dns_provider").(string) {
	case "route53":
		if v, ok := in.Get("route53_config").([]interface{}); ok && len(v) > 0 {
			route53Config, err := expandRoute53DNSConfig(v)
			if err != nil {
				return nil, err
			}
			obj.Route53ProviderConfig = route53Config
		}
	case "alidns":
		if v, ok := in.Get("alidns_config").([]interface{}); ok && len(v) > 0 {
			aliDNSConfig, err := expandAliDNSConfig(v)
			if err != nil {
				return nil, err
			}
			obj.AlidnsProviderConfig = aliDNSConfig
		}
	case "cloudflare":
		if v, ok := in.Get("rcloudflare_config").([]interface{}); ok && len(v) > 0 {
			cloudFlareConfig, err := expandCloudFlareDNSConfig(v)
			if err != nil {
				return nil, err
			}
			obj.CloudflareProviderConfig = cloudFlareConfig
		}
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj, nil
}

func expandRoute53DNSConfig(p []interface{}) (*managementClient.Route53ProviderConfig, error) {
	obj := &managementClient.Route53ProviderConfig{}
	if len(p) == 0 || p[0] == nil {
		return nil, fmt.Errorf("route_53 cannot be empty")
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

	return obj, nil

}

func expandAliDNSConfig(p []interface{}) (*managementClient.AlidnsProviderConfig, error) {
	obj := &managementClient.AlidnsProviderConfig{}
	if len(p) == 0 || p[0] == nil {
		return nil, fmt.Errorf("alidns config cannot be empty")
	}

	in := p[0].(map[string]interface{})

	if v, ok := in["access_key"].(string); ok && len(v) > 0 {
		obj.AccessKey = v
	}

	if v, ok := in["secret_key"].(string); ok && len(v) > 0 {
		obj.SecretKey = v
	}

	return obj, nil
}

func expandCloudFlareDNSConfig(p []interface{}) (*managementClient.CloudflareProviderConfig, error) {
	obj := &managementClient.CloudflareProviderConfig{}
	if len(p) == 0 || p[0] == nil {
		return nil, fmt.Errorf("cloudflare_config connot be empty")
	}

	in := p[0].(map[string]interface{})
	obj.APIEmail = in["api_email"].(string)
	obj.APIKey = in["api_key"].(string)
	obj.ProxySetting = in["proxy_setting"].(*bool)

	return obj, nil
}
