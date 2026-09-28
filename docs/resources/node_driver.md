---
page_title: "rancher2_node_driver Resource"
---

# rancher2\_node\_driver Resource

Provides a Rancher v2 Node Driver resource. This can be used to register node
drivers and retrieve their information. RKE2 provisioning requires an
RKE2-capable driver and its generated machine configuration schema.

## Example Usage

```hcl
# Create a new rancher2 Node Driver
resource "rancher2_node_driver" "foo" {
    active = true
    builtin = false
    checksum = "0x0"
    description = "Foo description"
    external_id = "foo_external"
    name = "foo"
    ui_url = "local://ui"
    url = "local://"
    whitelist_domains = ["*.foo.com"]
}
```

### Registering the T-Cloud Public driver

The following configuration pins an existing external `opentelekomcloud` driver
to v2.2.1 and manages the credential fields used by
[`tcloud_public_credential_config`](cloud_credential.md#tcloud_public_credential_config).
It does not install the Dashboard extension or network controller, create a
cluster, or manage a shared network.

On Rancher v2.15.1, a fresh v3 API create assigns a generated `nd-*` Kubernetes
metadata name even when `name = "opentelekomcloud"`; the machine provisioner
then cannot resolve `nodedrivers.management.cattle.io/opentelekomcloud`.
Bootstrap the driver with the stable metadata name before importing it into
Terraform:

```yaml
apiVersion: management.cattle.io/v3
kind: NodeDriver
metadata:
  name: opentelekomcloud
  annotations:
    lifecycle.cattle.io/create.node-driver-controller: "true"
    passwordFields: password,secretKey
    privateCredentialFields: password,secretKey
    publicCredentialFields: username,domainName,projectName,projectId,region,authUrl,accessKey
spec:
  active: true
  checksum: 4c1b3dce8e595655460d3045abc2e75c113b41d0ecfc0667befec1edce49a03f
  description: T-Cloud Public
  displayName: opentelekomcloud
  url: https://github.com/opentelekomcloud/docker-machine-opentelekomcloud/releases/download/v2.2.1/docker-machine-driver-opentelekomcloud_2.2.1_linux_amd64.tar.gz
```

Apply the manifest to Rancher's local cluster, wait for activation, and import
it before applying the Terraform configuration below. Do not use this resource
to perform the initial v2.15.1 registration until Rancher preserves the required
metadata name.

```hcl
resource "rancher2_node_driver" "tcloud_public" {
  name        = "opentelekomcloud"
  description = "T-Cloud Public"
  active      = true
  builtin     = false
  url         = "https://github.com/opentelekomcloud/docker-machine-opentelekomcloud/releases/download/v2.2.1/docker-machine-driver-opentelekomcloud_2.2.1_linux_amd64.tar.gz"
  checksum    = "4c1b3dce8e595655460d3045abc2e75c113b41d0ecfc0667befec1edce49a03f"

  annotations = {
    passwordFields          = "password,secretKey"
    privateCredentialFields = "password,secretKey"
    publicCredentialFields  = "username,domainName,projectName,projectId,region,authUrl,accessKey"
  }

  lifecycle {
    prevent_destroy = true
  }
}
```

The checksum is the **SHA-256 of the downloaded archive**, expressed as 64
hexadecimal characters without a `sha256:` prefix. Rancher verifies the download
before extracting the driver binary. This archive is for Linux amd64 on the
Rancher execution host, regardless of the Terraform client's OS. Use a matching
artifact and checksum if your execution host has a different architecture.

Keep `name = "opentelekomcloud"`: Rancher's v3 API maps `name` to the
NodeDriver's `spec.displayName`, which the controller uses for driver/schema
identity and derives from the executable. It is not a branding-only label.
Use the Terraform resource label and `description` for T-Cloud Public branding;
no separate `display_name` argument is needed.

The annotations deliberately exclude `authMethod`, which is not a v2.2.1 driver
flag. They select the seven public and two private fields of the credential
schema. With `active = true`, Rancher exposes these credentials using the
existing node-driver resource arguments; no additional provider parameter is
required.

Register and import the driver before creating credentials and machine
configurations. Use `depends_on = [rancher2_node_driver.tcloud_public]` on both
resources when managing them in the same configuration. Driver activation is not
proof of shared-network readiness or a usable RKE2 cluster. Before workload
provisioning, verify the generated schemas in Rancher's **local** cluster:

```sh
kubectl get dynamicschemas.management.cattle.io opentelekomcloudcredentialconfig -o yaml
kubectl wait --for=condition=Established --timeout=120s \
  crd/opentelekomcloudconfigs.rke-machine-config.cattle.io
```

Confirm that the credential schema exposes the nine fields above and the machine
CRD serves `v1` with kind `OpentelekomcloudConfig`. The contract, checksum,
fresh driver registration, and end-to-end provisioning were verified against
Rancher v2.15.1 and driver v2.2.1. Earlier driver v2.2.0 releases accepted only
`OS_*` authentication variables and are not compatible with the
`OPENTELEKOMCLOUD_*` variables supplied by Rancher's provisioning jobs. Use
v2.2.1 or later.

Import the stable Kubernetes metadata name rather than creating a second
registration:

```sh
terraform import rancher2_node_driver.tcloud_public opentelekomcloud
```

Inspect the subsequent plan before applying: the existing URL, checksum,
annotations, or `builtin` status may differ. Do not automatically replace or
deactivate a legacy driver used by other clusters. Keep this shared registration
until all dependent clusters and machines have been removed; intentionally
remove `prevent_destroy` only when decommissioning an unused driver.

## Argument Reference

The following arguments are supported:

* `active` - (Required) Specify if the node driver state (bool)
* `builtin` - (Required) Specify wheter the node driver is an internal node driver or not (bool)
* `name` - (Required) Name of the node driver (string)
* `url` - (Required) The URL to download the machine driver binary for 64-bit Linux (string)
* `checksum` - (Optional) Verify that the downloaded driver matches the expected checksum (string)
* `description` - (Optional) Description of the node driver (string)
* `external_id` - (Optional) External ID (string)
* `ui_url` - (Optional) The URL to load for customized Add Nodes screen for this driver (string)
* `whitelist_domains` - (Optional) Domains to whitelist for the ui (list)
* `annotations` - (Optional/Computed) Annotations of the resource (map)
* `labels` - (Optional/Computed) Labels of the resource (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)

## Timeouts

`rancher2_node_driver` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating node drivers.
- `update` - (Default `10 minutes`) Used for node driver modifications.
- `delete` - (Default `10 minutes`) Used for deleting node drivers.

## Import

Node Driver can be imported using the Rancher Node Driver ID

```
$ terraform import rancher2_node_driver.foo &lt;node_driver_id&gt;
```
