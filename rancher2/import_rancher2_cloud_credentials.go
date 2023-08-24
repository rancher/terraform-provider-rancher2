package rancher2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2CloudCredentialsImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	cloudCredentialID, driver := splitID(d.Id())
	d.Set("driver", driver)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	cloudCredential := &CloudCredential{}
	err = client.APIBaseClient.ByID(managementClient.CloudCredentialType, cloudCredentialID, cloudCredential)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	err = flattenCloudCredential(d, cloudCredential)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
