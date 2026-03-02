provider "aws" {
  default_tags {
    tags = {
      Id    = local.identifier
      Owner = local.owner
    }
  }
}

locals {
  identifier = var.identifier
  owner      = var.owner
}

module "s3_bucket" {
  source        = "terraform-aws-modules/s3-bucket/aws"
  version       = "5.0.2"
  bucket        = lower(local.identifier)
  force_destroy = true
  versioning = {
    status     = true
    mfa_delete = false
  }
}
