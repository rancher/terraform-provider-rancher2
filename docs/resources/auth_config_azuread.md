---
page_title: "rancher2_auth_config_azuread Resource"
---

# rancher2\_auth\_config\_azuread Resource

Provides a Rancher v2 Auth Config AzureAD resource. This can be used to configure and enable Auth Config AzureAD for Rancher v2 RKE clusters and retrieve their information.

In addition to the built-in local auth, only one external auth config provider can be enabled at a time.

**Important:** for existing Azure AD setups created in versions of Rancher before v2.6.7,
admins need to perform an update to the AuthConfig resource if they upgrade to Rancher v2.6.7+.
They need to first follow the documentation to set the proper [Application type permissions](https://rancher.com/docs/rancher/v2.6/en/admin-settings/authentication/azure-ad/#3-set-required-permissions-for-rancher) on the App registration on the Azure portal.

Then they need to go to the UI, which will [prompt](https://rancher.com/docs/rancher/v2.6/en/admin-settings/authentication/azure-ad/#migrating-from-azure-ad-graph-api-to-microsoft-graph-api) them to update their configuration on every screen. If they had set the permissions
properly in the previous step, then Rancher will perform the configuration migration successfully.

Finally, to avoid state drift in the Terraform configuration, admins need to look up the new values of the updated AuthConfig and
update the corresponding values in their `rancher2_auth_config_azuread` resource. They can then run `terraform plan` and ensure that
Terraform has nothing to change.

## Example Usage

```hcl
# Create a new rancher2 Auth Config AzureAD
resource "rancher2_auth_config_azuread" "azuread" {
  application_id = "<AZUREAD_APP_ID>"
  application_secret = "<AZUREAD_APP_SECRET>"
  auth_endpoint = "<AZUREAD_AUTH_ENDPOINT>"
  graph_endpoint = "<AZUREAD_GRAPH_ENDPOINT>"
  rancher_url = "<RANCHER_URL>"
  tenant_id = "<AZUREAD_TENANT_ID>"
  token_endpoint = "<AZUREAD_TOKEN_ENDPOINT>"
}
```

## Argument Reference

The following arguments are supported:

* `application_id` - (Required/Sensitive) AzureAD auth application ID (string)
* `application_secret` - (Required/Sensitive) AzureAD auth application secret (string)
* `auth_endpoint` - (Required) AzureAD auth endpoint (string)
* `graph_endpoint` - (Required) AzureAD graph endpoint (string)
* `rancher_url` - (Required) Rancher URL (string). "<rancher_url>/verify-auth-azure"
* `tenant_id` - (Required) AzureAD tenant ID (string)
* `token_endpoint` - (Required) AzureAD token endpoint (string)
* `endpoint` - (Optional) AzureAD endpoint. Default `https://login.microsoftonline.com/` (string)
* `access_mode` - (Optional) Access mode for auth. `required`, `restricted`, `unrestricted` are supported. Default `unrestricted` (string)
* `allowed_principal_ids` - (Optional) Allowed principal ids for auth. Required if `access_mode` is `required` or `restricted`. Ex: `azuread_user://<USER_ID>`  `azuread_group://<GROUP_ID>` (list)
* `enabled` - (Optional) Enable auth config provider. Default `true` (bool)
* `tls` - (Optional) Enable TLS connection. Default `true` (bool)
* `annotations` - (Optional/Computed) Annotations of the resource (map)
* `labels` - (Optional/Computed) Labels of the resource (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `name` - (Computed) The name of the resource (string)
* `type` - (Computed) The type of the resource (string)
