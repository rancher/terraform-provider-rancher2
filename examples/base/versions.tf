terraform {
  required_version = ">= 1.5.0"
  required_providers {
    local = {
      source  = "hashicorp/local"
      version = ">= 2.4"
    }
    random = {
      source  = "hashicorp/random"
      version = ">= 3.5.1"
    }
    github = {
      source  = "integrations/github"
      version = ">= 5.44"
    }
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.11"
    }
    http = {
      source  = "hashicorp/http"
      version = ">= 3.4"
    }
    null = {
      source  = "hashicorp/null"
      version = ">= 3"
    }
    tls = {
      source  = "hashicorp/tls"
      version = ">= 4.0"
    }
    acme = {
      source  = "vancluever/acme"
      version = ">= 2.0"
    }
    cloudinit = {
      source  = "hashicorp/cloudinit"
      version = ">= 2.3.3"
    }
    helm = {
      source  = "hashicorp/helm"
      version = ">= 2.14"
    }
    rancher2 = {
      source  = "rancher/rancher2"
      version = ">= 5.0.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.31.0"
    }
  }
}
