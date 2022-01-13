package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2CloudCredentialsImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	cloudCredential := &CloudCredential{}
	err = client.APIBaseClient.ByID(managementClient.CloudCredentialType, d.Id(), cloudCredential)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	drivers := []string{
		amazonec2ConfigDriver,
		azureConfigDriver,
		digitaloceanConfigDriver,
		googleConfigDriver,
		s3ConfigDriver,
		vmwarevsphereConfigDriver,
	}

	// Missing "driver" field in api.
	for _, driver := range drivers {
		d.Set("driver", driver)
		err = flattenCloudCredential(d, cloudCredential)
		if err != nil {
			return []*schema.ResourceData{}, err
		}
	}

	return []*schema.ResourceData{d}, nil
}
