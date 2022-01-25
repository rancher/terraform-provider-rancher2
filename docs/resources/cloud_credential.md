---
page_title: "rancher2_cloud_credential Resource"
---

# rancher2\_cloud\_credential Resource

Provides a Rancher v2 Cloud Credential resource. This can be used to create Cloud Credential for Rancher v2.2.x and retrieve their information.

amazonec2, azure, digitalocean, linode, openstack and vsphere credentials config are supported for Cloud Credential.

## Example Usage

```hcl
# Create a new rancher2 Cloud Credential
resource "rancher2_cloud_credential" "foo" {
  name = "foo"
  description = "foo test"
  amazonec2_credential_config {
    access_key = "<AWS_ACCESS_KEY>"
    secret_key = "<AWS_SECRET_KEY>"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Cloud Credential (string)
* `amazonec2_credential_config` - (Optional) AWS config for the Cloud Credential (list maxitems:1)
* `azure_credential_config` - (Optional) Azure config for the Cloud Credential (list maxitems:1)
* `description` - (Optional) Description for the Cloud Credential (string)
* `digitalocean_credential_config` - (Optional) DigitalOcean config for the Cloud Credential (list maxitems:1)
* `google_credential_config` - (Optional) Google config for the Cloud Credential (list maxitems:1)
* `linode_credential_config` - (Optional) Linode config for the Cloud Credential (list maxitems:1)
* `openstack_credential_config` - (Optional) OpenStack config for the Cloud Credential (list maxitems:1)
* `s3_credential_config` - (Optional) S3 config for the Cloud Credential. Just for Rancher 2.6.0 and above (list maxitems:1)
* `vsphere_credential_config` - (Optional) vSphere config for the Cloud Credential (list maxitems:1)
* `annotations` - (Optional) Annotations for Cloud Credential object (map)
* `labels` - (Optional/Computed) Labels for Cloud Credential object (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `driver` - (Computed) The driver of the Cloud Credential (string)

## Nested blocks

### `amazonec2_credential_config`

#### Arguments

* `access_key` - (Required/Sensitive) AWS access key (string)
* `secret_key` - (Required/Sensitive) AWS secret key (string)
* `default_region` - (Optional) AWS default region (string)

### `azure_credential_config`

#### Arguments

* `client_id` - (Required/Sensitive) Azure Service Principal Account ID (string)
* `client_secret` - (Required/Sensitive) Azure Service Principal Account password (string)
* `subscription_id` - (Required/Sensitive) Azure Subscription ID (string)
* `environment` - (Optional/Computed) Azure environment (e.g. AzurePublicCloud, AzureChinaCloud) (string)
* `tenant_id` - (Optional/Computed) Azure Tenant ID (string)

### `digitalocean_credential_config`

#### Arguments

* `access_token` - (Required/Sensitive) DigitalOcean access token (string)

### `google_credential_config`

#### Arguments

* `auth_encoded_json` - (Required/Sensitive) Google auth encoded json (string)

### `linode_credential_config`

#### Arguments

* `token` - (Required/Sensitive) Linode API token (string)

### `openstack_credential_config`

#### Arguments

* `password` - (Required/Sensitive) OpenStack password (string)

### `s3_credential_config`

#### Arguments

* `access_key` - (Required/Sensitive) AWS access key (string)
* `secret_key` - (Required/Sensitive) AWS secret key (string)
* `default_bucket` - (Optional) AWS default bucket (string)
* `default_endpoint` - (Optional) AWS default endpoint (string)
* `default_endpoint_ca` - (Optional/Sensitive) AWS default endpoint CA (string)
* `default_folder` - (Optional) AWS default folder (string)
* `default_region` - (Optional) AWS default region (string)
* `default_skip_ssl_verify` - (Optional) AWS default skip ssl verify. Default: `false` (bool)

### `vsphere_credential_config`

#### Arguments

* `password` - (Required/Sensitive) vSphere password (string)
* `username` - (Required) vSphere username (string)
* `vcenter` - (Required) vSphere IP/hostname for vCenter (string)
* `vcenter_port` - (Optional) vSphere Port for vCenter. Default `443` (string)

## Timeouts

`rancher2_cloud_credential` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating cloud credentials.
- `update` - (Default `10 minutes`) Used for cloud credential modifications.
- `delete` - (Default `10 minutes`) Used for deleting cloud credentials.

## Import

Node Template can be imported using the Rancher Node Template ID

```bash
terraform import rancher2_cloud_credential.foo &lt;cloud_credential_id&gt;.&lt;driver&gt;
```

### Argument Reference

The following arguments are supported:

* `cloud_credential_id` - (Required) The ID of the Cloud Credential (string)
* `driver` - (Required) The driver of the Cloud Credential (string)

Supported drivers:

* amazonec2
* azure
* digitalocean
* googlekubernetesengine
* linode
* openstack
* s3
* vmwarevsphere
