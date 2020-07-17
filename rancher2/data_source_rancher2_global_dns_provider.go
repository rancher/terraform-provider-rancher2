package rancher2

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

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
				Optional: true,
			},
			"dns_provider": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"route53", "cloudflare", "alidns"}, true),
			},
			"route53_config": {
				Type:          schema.TypeSet,
				MaxItems:      1,
				Optional:      true,
				Elem: &schema.Resource{
					Schema: globalDNSProviderRoute53ConfigSchema(),
				},
			},
			"alidns_config": {
				Type:          schema.TypeSet,
				MaxItems:      1,
				Optional:      true,
				Elem: &schema.Resource{
					Schema: globalDNSProviderAliConfigSchema(),
				},
			},
			"cloudflare_config": {
				Type:          schema.TypeSet,
				MaxItems:      1,
				Optional:      true,
				Elem: &schema.Resource{
					Schema: globalDNSProviderCloudFareConfigSchema(),
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

	globalDNSProvider, err := client.GlobalDNSProvider.List(listOpts)
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

	return flattenGlobalDNSProvider(d, &globalDNSProvider.Data[0])
}
