# RKE2 Cluster on Proxmox VE

This example provisions an RKE2 cluster in Rancher using Proxmox VE (PVE) as the infrastructure provider via the `docker-machine-driver-pve` node driver.

It assumes you already have:

- A Rancher instance with the Proxmox VE node driver installed and active
- A PVE cloud credential already created in Rancher (URL + API token stored there)
- A Proxmox VE VM template accessible from Rancher (cloud-init capable)

## Prerequisites

- OpenTofu >= 1.5.0 (or Terraform >= 1.5.0)
- Rancher2 provider >= 5.0.0
- PVE node driver: [Stellatarum/docker-machine-driver-pve](https://github.com/Stellatarum/docker-machine-driver-pve)

## Variables

| Name | Description | Required |
|------|-------------|----------|
| `rancher_url` | Rancher API URL | yes |
| `rancher_credentials` | Object with `access_key` and `secret_key` | yes |
| `rancher_insecure` | Skip TLS verification (self-signed certs) | no (default: `false`) |
| `cluster_name` | Name for the new cluster in Rancher | yes |
| `kubernetes_version` | RKE2 version (e.g. `v1.34.3+rke2r1`) | yes |
| `pve_cloud_credential_name` | Name of the existing PVE cloud credential in Rancher | yes |
| `pve_resource_pool` | Proxmox VE resource pool name | no |
| `pve_network_interface` | Network interface bus/device (e.g. `net0`) | yes |
| `pve_iso_device` | CD/DVD drive bus/device for cloud-init ISO (e.g. `ide2`) | yes |
| `pve_ssh_user` | SSH user created by cloud-init | no (default: `service`) |
| `pve_full_clone` | Force full disk clone instead of linked clone | no (default: `true`) |
| `pve_tags` | Comma-separated VM tags | no |
| `pve_server_template_id` | Proxmox VM template ID for server nodes | yes |
| `pve_worker_template_id` | Proxmox VM template ID for worker nodes | yes |
| `server_quantity` | Number of server (control-plane + etcd) nodes | yes |
| `worker_quantity` | Number of worker nodes | yes |

### terraform.tfvars example

```hcl
rancher_url               = "https://rancher.example.com"
rancher_insecure          = false
cluster_name              = "pve-cluster"
kubernetes_version        = "v1.34.3+rke2r1"
pve_cloud_credential_name = "my-pve-credential"
pve_resource_pool         = "my-pool"
pve_network_interface     = "net0"
pve_iso_device            = "ide2"
pve_ssh_user              = "service"
pve_full_clone            = true
pve_server_template_id    = "100"
pve_worker_template_id    = "100"
server_quantity           = 3
worker_quantity           = 2
```


