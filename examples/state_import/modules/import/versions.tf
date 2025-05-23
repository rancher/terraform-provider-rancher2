terraform {
  required_version = ">= 1.5.0"
  required_providers {
    rancher2 = {
      source  = "rancher/rancher2"
      version = ">= 6.0.0"
    }
  }
}
