package rancher2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRancher2GlobalDNSProvider() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRancher2GlobalDNSProviderRead,

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
				Computed: true,
				Elem: &schema.Resource{
					Schema: globalDNSProviderAliConfigSchema(),
				},
			},
			"cloudflare_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: globalDNSProviderCloudFareConfigSchema(),
				},
			},
			"route53_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: globalDNSProviderRoute53ConfigSchema(),
				},
			},
		},
	}
}

func dataSourceRancher2GlobalDNSProviderRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}
	name := d.Get("name").(string)

	filters := map[string]interface{}{
		"name": name,
	}

	listOpts := NewListOpts(filters)

	globalDNSProvider, err := client.GlobalDnsProvider.List(listOpts)
	if err != nil {
		return diag.FromErr(err)
	}

	count := len(globalDNSProvider.Data)
	if count <= 0 {
		return diag.Errorf("[ERROR] global dns provider with name \"%s\" not found", name)
	}
	if count > 1 {
		return diag.Errorf("[ERROR] found %d global dns provider with name \"%s\"", count, name)
	}

	// TODO - Validate if this is required there IMO i se no reason to call it twice.
	_ = flattenGlobalDNSProvider(d, &globalDNSProvider.Data[0])

	//return fmt.Errorf("[ERROR] %#v\n%#v", d.Get("route53_config"), globalDNSProvider.Data[0].Route53ProviderConfig)

	return diag.FromErr(flattenGlobalDNSProvider(d, &globalDNSProvider.Data[0]))
}
