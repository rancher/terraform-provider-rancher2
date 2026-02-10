terraform {
  required_version = ">= 1.5.0"
  required_providers {
    rancher2 = {
      source  = "rancher/rancher2"
      version = ">= 0.1.0"
    }
    random = {
      source  = "hashicorp/random"
      version = ">= 3.5.1"
    }
  }
}
