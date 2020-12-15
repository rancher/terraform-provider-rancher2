package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceRancher2GlobalDNSProvider() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2GlobalDNSProviderRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"annotations": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"root_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dns_provider": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"alidns_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: globalDNSProviderAliConfigSchema(),
				},
			},
			"cloudflare_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: globalDNSProviderCloudFareConfigSchema(),
				},
			},
			"route53_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: globalDNSProviderRoute53ConfigSchema(),
				},
			},
		},
	}
}

func dataSourceRancher2GlobalDNSProviderRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}
	name := d.Get("name").(string)

	filters := map[string]interface{}{
		"name": name,
	}

	listOpts := NewListOpts(filters)

	globalDNSProvider, err := client.GlobalDnsProvider.List(listOpts)
	if err != nil {
		return err
	}

	count := len(globalDNSProvider.Data)
	if count <= 0 {
		return fmt.Errorf("[ERROR] global dns provider with name \"%s\" not found", name)
	}
	if count > 1 {
		return fmt.Errorf("[ERROR] found %d global dns provider with name \"%s\"", count, name)
	}

	flattenGlobalDNSProvider(d, &globalDNSProvider.Data[0])

	//return fmt.Errorf("[ERROR] %#v\n%#v", d.Get("route53_config"), globalDNSProvider.Data[0].Route53ProviderConfig)

	return flattenGlobalDNSProvider(d, &globalDNSProvider.Data[0])
}
