package rancher2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2NodeTemplateImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return []*schema.ResourceData{}, err
	}
	nodeTemplate := &NodeTemplate{}
	err = client.APIBaseClient.ByID(managementClient.NodeTemplateType, d.Id(), nodeTemplate)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	err = flattenNodeTemplate(d, nodeTemplate)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
