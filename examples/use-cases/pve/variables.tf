variable "rancher_credentials" {
  description = "Rancher API access and secret key."
  type = object({
    access_key = string
    secret_key = string
  })
  nullable  = false
  sensitive = true
}

variable "rancher_url" {
  description = "Rancher API URL (e.g. https://rancher.example.com)."
  type        = string
}

variable "cluster_name" {
  description = "Name for the new RKE2 cluster in Rancher."
  type        = string
}

variable "kubernetes_version" {
  description = "RKE2 Kubernetes version (e.g. v1.34.3+rke2r1)."
  type        = string
}

variable "registry_mirror" {
  description = "Internal registry mirror used for docker.io, ghcr.io and quay.io."
  type        = string
}

variable "rancher_insecure" {
  description = "Skip TLS verification for Rancher (use when CA cert is not in system trust store)."
  type        = bool
  default     = false
}

# Cloud credential name (must already exist in Rancher)
variable "pve_cloud_credential_name" {
  description = "Name of the Proxmox VE cloud credential already configured in Rancher."
  type        = string
}

# Proxmox VE — machine-level settings only; URL and credentials come from the cloud credential
variable "pve_resource_pool" {
  description = "Proxmox VE Resource Pool name."
  type        = string
  default     = ""
}

variable "pve_network_interface" {
  description = "Bus/device of the network interface to read the machine IP from (e.g. net0)."
  type        = string
  default     = "net0"
}

variable "pve_iso_device" {
  description = "Bus/device of the CD/DVD drive to mount the cloud-init ISO to (e.g. scsi1)."
  type        = string
  default     = "scsi1"
}

variable "pve_ssh_user" {
  description = "SSH user created via cloud-init."
  type        = string
  default     = "service"
}

variable "pve_full_clone" {
  description = "Forces a full copy of all disks, even if the storage supports linked clones."
  type        = bool
  default     = false
}

variable "pve_tags" {
  description = "Comma-separated list of tags to assign to the VMs."
  type        = string
  default     = ""
}

# Server (control-plane + etcd) pool
variable "pve_server_template_id" {
  description = "ID of the Proxmox VE VM template to use for server nodes."
  type        = string
}

variable "pve_server_sockets" {
  description = "Number of CPU sockets for server nodes."
  type        = string
  default     = "1"
}

variable "pve_server_cores" {
  description = "Number of CPU cores for server nodes."
  type        = string
  default     = "2"
}

variable "pve_server_memory" {
  description = "Memory in MiB for server nodes."
  type        = string
  default     = "4096"
}

variable "server_quantity" {
  description = "Number of server (control-plane + etcd) nodes."
  type        = number
  default     = 1
}

# Worker pool
variable "pve_worker_template_id" {
  description = "ID of the Proxmox VE VM template to use for worker nodes."
  type        = string
}

variable "pve_worker_sockets" {
  description = "Number of CPU sockets for worker nodes."
  type        = string
  default     = "1"
}

variable "pve_worker_cores" {
  description = "Number of CPU cores for worker nodes."
  type        = string
  default     = "2"
}

variable "pve_worker_memory" {
  description = "Memory in MiB for worker nodes."
  type        = string
  default     = "4096"
}

variable "worker_quantity" {
  description = "Number of worker nodes."
  type        = number
  default     = 1
}
