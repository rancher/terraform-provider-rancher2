package rancher2

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRancher2Setting() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRancher2SettingRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRancher2SettingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	log.Printf("[INFO] Refreshing Rancher2 Setting: %s", name)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	setting, err := client.Setting.ByID(name)
	if err != nil || setting == nil {
		return diag.FromErr(err)
	}

	d.SetId(name)
	d.Set("value", setting.Value)

	return nil
}
