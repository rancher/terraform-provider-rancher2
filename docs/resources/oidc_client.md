---
page_title: "rancher2_oidc_client Resource"
---

# rancher2\_oidc_client Resource

Provides a Rancher OIDC Client. This can be used to configure the OIDC Clients
Provides a Rancher OIDC Client. This can be used to configure the OIDC Clients available for the Rancher OIDC Provider.

## Example Usage

```hcl
# Create a new rancher2 OIDC Client
resource "rancher2_oidc_client" "foo" {
  description = "Foo OIDC Client"
  token_expiration_seconds = 600 # expiration of the id_token and access_token
  refresh_token_expiration_seconds = 7200 # expiration of the refresh_token
  redirect_uris = [
    "http://127.0.0.1:5556/auth/rancher/callback",
    "http://127.0.0.1:33418/",
    "https://vscode.dev/redirect"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `redirect_uris` - (Required) List of allowed redirect URIs for this OIDC Client (list)
* `description` - (Optional) A human-readable description for the OIDC Client (string)
* `token_expiration_seconds` - (Optional/Computed) ID Token and Access Token will only be valid for this many seconds (int)
* `refresh_token_expiration_seconds` - (Optional/Computed) How long can the refresh token be used for (int)
* `annotations` - (Optional/Computed) Annotations for OIDC Client object (map)
* `labels` - (Optional/Computed) Labels for OIDC Client object (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `client_id` - (Computed) The ID to be used when authenticating as this OIDC Client.

## Import

OIDC Clients can be imported using the Client name in the format `<client_name>`

```
$ terraform import rancher2_oidc_client.foo <client_name>
```
