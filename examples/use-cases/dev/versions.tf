terraform {
  required_version = ">= 1.5.0"
  required_providers {
    file = {
      source  = "rancher/file"
      version = ">= 2.2"
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
      version = "2.14"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.31.0"
    }
    time = {
      source  = "hashicorp/time"
      version = ">= 0.12.0"
    }
  }
}
