---
page_title: "Rancher2 Provider"
---

# Rancher2 Provider

The Rancher2 provider is used to interact with the
resources supported by Rancher v2. 

The provider can be configured in 2 modes:
- Admin: this is the default mode, intended to manage rancher2 resources. It should be configured with the `api_url` of the Rancher server and API credentials, `token_key` or `access_key` and `secret_key`.
- Bootstrap: this mode is intended to bootstrap a rancher2 system. It is enabled if `bootstrap = true`. In this mode, `token_key` or `access_key` and `secret_key` can not be provided. More info at [rancher2_bootstrap resource](resources/bootstrap.html)

## Example Usage

```hcl
# Configure the Rancher2 provider to admin
provider "rancher2" {
  api_url    = "https://rancher.my-domain.com"
  access_key = var.rancher2_access_key
  secret_key = var.rancher2_secret_key
}
```

```hcl
# Configure the Rancher2 provider to bootstrap
provider "rancher2" {
  api_url   = "https://rancher.my-domain.com"
  bootstrap = true
}
```

```hcl
# Configure the Rancher2 provider to bootstrap and admin
# Provider config for bootstrap
provider "rancher2" {
  alias = "bootstrap"

  api_url   = "https://rancher.my-domain.com"
  bootstrap = true
}

# Create a new rancher2_bootstrap using bootstrap provider config
resource "rancher2_bootstrap" "admin" {
  provider = "rancher2.bootstrap"

  password = "blahblah"
}

# Provider config for admin
provider "rancher2" {
  alias = "admin"

  api_url = rancher2_bootstrap.admin.url
  token_key = rancher2_bootstrap.admin.token
  insecure = true
}

# Create a new rancher2 resource using admin provider config
resource "rancher2_catalog_v2" "foo" {
  provider = "rancher2.admin"

  cluster_id = "<CLUSTER_ID>"
  name = "test"
  url = "http://foo.com:8080"
}
```

## Argument Reference

The following arguments are supported:

* `api_url` - (Required) Rancher API url. It must be provided, but it can also be sourced from the `RANCHER_URL` environment variable.
* `access_key` - (Optional/Sensitive) Rancher API access key to connect to rancher. It can also be sourced from the `RANCHER_ACCESS_KEY` environment variable.
* `secret_key` - (Optional/Sensitive) Rancher API secret key to connect to rancher. It can also be sourced from the `RANCHER_SECRET_KEY` environment variable.
* `token_key` - (Optional/Sensitive) Rancher API token key to connect to rancher. It can also be sourced from the `RANCHER_TOKEN_KEY` environment variable. Could be used instead `access_key` and `secret_key`.
* `ca_certs` - CA certificates used to sign Rancher server tls certificates. Mandatory if self signed tls and insecure option false. It can also be sourced from the `RANCHER_CA_CERTS` environment variable.
* `insecure` - (Optional) Allow insecure connection to Rancher. Mandatory if self signed tls and not ca_certs provided. It can also be sourced from the `RANCHER_INSECURE` environment variable.
* `bootstrap` - (Optional) Enable bootstrap mode to manage `rancher2_bootstrap` resource. It can also be sourced from the `RANCHER_BOOTSTRAP` environment variable. Default: `false`
* `retries` - (Deprecated) Use timeout instead
* `timeout` - (Optional) Timeout duration to retry for Rancher connectivity and resource operations. Default: `"120s"`

## Security and Sensitive State

> **CRITICAL SECURITY NOTE:** Although some provider arguments and resource attributes are marked as sensitive to prevent them from being displayed in visual CLI outputs (such as `terraform plan` or `terraform apply`), **all sensitive values are stored in cleartext inside the Terraform state file.**
> 
> The sensitive flag in provider and resource schemas only prevents values from leaking into visual UI/CLI console logs. It does not encrypt, mask, or secure the values on disk or in the state storage backend. This represents the larger residual risk of managing secrets within Terraform configuration.
> 
> To mitigate this risk, it is highly recommended to:
> 1. Use a **secure remote state backend** (e.g., HashiCorp Consul, AWS S3 with KMS encryption, HashiCorp Cloud Platform, etc.) that natively supports **encryption at-rest** and **transit-level TLS encryption**.
> 2. Implement strict **Identity and Access Management (IAM)** and **access control policies** to restrict state file read permissions to authorized personnel and automated pipelines only.
> 3. Ensure that **live infrastructure repositories and generated Terraform plans are properly secured** and kept private. Live infrastructure plans and state information should never be public knowledge, as they expose sensitive architectural and operational details.
> 4. Be cautious when enabling Terraform debug logging (e.g., setting `TF_LOG=DEBUG` or `TRACE`). Debug logs will print sensitive fields, provider configurations, and API responses in plain text. Additionally, using the `--json` flag for machine-readable output will also include sensitive material in plain text. Always sanitize these logs and outputs before sharing them in bug reports or support tickets.
