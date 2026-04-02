# provider "rancher2" {}

resource "rancher2_dev2" "full" {
  # id                  = "string" # this is read only, computed by the provider, and filtered from the API
  api_version = "string" # required
  kind        = "string" # required
  metadata = {
    namespace = "string" # optional, defaults to "default"
    name      = "string" # optional, mutually exclusive with generate_name, either name or generate_name must be set
    # generate_name = "string" # optional, mutually exclusive with name, either name or generate_name must be set
    # Annotations, labels, finalizers, and owner_references have special provider logic
    # The provider will enforce anything set by the user, but won't delete anything set by the API.
    # For instance if the user sets the label string = "string", and Kubernetes removes it, then Terraform will attempt to add it back.
    # If Kubernetes adds the label "node: A", it will be available as a read-only attribute to other resources and Terraform won't try to remove it.
    # If the user adds a label that Kubernetes changes, like say "node = A" and Kubernetes attempts to change it to "node: B", Terraform will detect drift and reconcile.
    annotations = { string = "string" } # optional, computed
    labels      = { string = "string" } # optional, computed
    # finalizers has a special mode, it normally works like labels, but when the user sends an empty list it sends a patch to Kubernetes removing all finalizers.
    finalizers = ["string"] # advanced use cases only, computed, optional with warning
    owner_references = [    # advanced use cases only, computed, optional with warning
      {
        api_version          = "string"
        kind                 = "string"
        name                 = "string"
        uid                  = "string"
        controller           = true
        block_owner_deletion = true
      }
    ]
    # "uid": "test",                      # computed, read-only
    # "generation": "test",               # computed, read-only
    # "creation_timestamp": "test",       # computed, read-only
    # "deletion_grace_period_seconds": 1, # computed, read-only
    # "deletion_timestamp": "test",       # computed, read-only
    # "managed_fields": "json blob"       # computed, read-only
    # "resource_version": "test",         # computed, read-only
    # "self_link": "test"                 # computed, read-only
  }
  spec = { # required
    # spec data below here in the same structure as the API
    string  = "test"
    bool    = false
    number  = 1
    int32   = 1
    int64   = 1
    float32 = 1.0
    float64 = 1.0
    map     = { "test" = "test" }
    list    = ["test"]
    object = {
      string_attribute = "test"
    }
    object_list = [
      {
        string_attribute = "test"
      }
    ]
    object_map = {
      "test" = {
        string_attribute = "test"
      }
    }
  }
  # status = "json blob" # computed, read only

  ## IGNORE api_response when using this as a template to generate a new resource, this is DEV resource only.
  api_responses = {
    "create" = {
      headers = {
        "Content-Type" = ["application/json"]
      }
      body        = jsonencode(local.createResponseBody)
      status_code = 200
    }
    "read" = {
      headers = {
        "Content-Type" = ["application/json"]
      }
      body        = jsonencode(local.createResponseBody)
      status_code = 200
    }
    "update" = {
      headers = {
        "Content-Type" = ["application/json"]
      }
      body        = jsonencode(local.createResponseBody)
      status_code = 200
    }
    "delete" = {
      headers     = {}
      body        = ""
      status_code = 200
    }
  }
}

# resource "rancher2_dev2" "minimal" {
#   api_version = "string"
#   kind        = "string"
#   metadata = {
#     name      = "string"
#     namespace = "string"
#   }
#   spec = {
#     string = "test"
#   }
# }

# read only attributes
output "rancher2_dev2_internal_id" {
  value = rancher2_dev2.full.id
}
output "rancher2_dev2_uid" {
  value = rancher2_dev2.full.metadata.uid
}
output "rancher2_dev2_generation" {
  value = rancher2_dev2.full.metadata.generation
}
output "rancher2_dev2_creation_timestamp_simple" {
  # input timestamps must be in RFC3339 format, which is common in Kubernetes
  value = formatdate("DD-MMM-YY hh:mm:ss", rancher2_dev2.full.metadata.creation_timestamp)
}
output "rancher2_dev2_deletion_grace_period_seconds" {
  value = rancher2_dev2.full.metadata.deletion_grace_period_seconds
}
output "rancher2_dev2_deletion_timestamp_simple" {
  # input timestamps must be in RFC3339 format, which is common in Kubernetes
  value = formatdate("DD-MMM-YY hh:mm:ss", rancher2_dev2.full.metadata.deletion_timestamp)
}
output "rancher2_dev2_managed_fields" {
  value = jsondecode(rancher2_dev2.full.metadata.managed_fields)
}
output "rancher2_dev2_resource_version" {
  value = rancher2_dev2.full.metadata.resource_version
}
output "rancher2_dev2_self_link" {
  value = rancher2_dev2.full.metadata.self_link
}
output "rancher2_dev2_status" {
  value = jsondecode(rancher2_dev2.full.status)
}

# special attributes (writable attributes added by kubernetes, not user, read only)
output "rancher2_dev2_string_annotation" {
  value = rancher2_dev2.full.metadata.annotations.string
}
output "rancher2_dev2_string_label" {
  value = rancher2_dev2.full.metadata.labels.string
}
output "rancher2_dev2_finalizers" {
  value = rancher2_dev2.full.metadata.finalizers
}
output "rancher2_dev2_owner_references" {
  value = rancher2_dev2.full.metadata.owner_references
}

# attributes added by user
output "rancher2_dev2_name" {
  value = rancher2_dev2.full.metadata.name
}
# output "rancher2_dev2_name_prefix" { # mutually exclusive with name
#   value = rancher2_dev2.full.metadata.generate_name
# }
output "rancher2_dev2_namespace" {
  value = rancher2_dev2.full.metadata.namespace
}
output "rancher2_dev2_api_version" {
  value = rancher2_dev2.full.api_version
}
output "rancher2_dev2_kind" {
  value = rancher2_dev2.full.kind
}
output "rancher2_dev2_spec" {
  value = rancher2_dev2.full.spec
}
output "rancher2_dev2_annotation_string" {
  value = rancher2_dev2.full.metadata.annotations.string
}
output "rancher2_dev2_label_string" {
  value = rancher2_dev2.full.metadata.labels.string
}
output "rancher2_dev2_finalizer_string" {
  value = try(rancher2_dev2.full.metadata.finalizers[index(rancher2_dev2.full.metadata.finalizers, "string")], "")
}
output "rancher2_dev2_owner_reference_api_version" {
  value = try(rancher2_dev2.full.metadata.owner_references[index(rancher2_dev2.full.metadata.owner_references.*.uid, "string")], "")
}

locals {
  createResponseBody = {
    api_version = "string"
    kind        = "string"
    metadata = {
      uid           = "string"
      name          = "string"
      generate_name = ""
      namespace     = "string"
      annotations = {
        string = "string"
      }
      labels = {
        string = "string"
      }
      finalizers = ["string"]
      owner_references = [
        {
          api_version          = "string"
          kind                 = "string"
          name                 = "string"
          uid                  = "string"
          controller           = true
          block_owner_deletion = true
        }
      ]
      generation         = 1
      creation_timestamp = "test"
      deletion_timestamp = "test"
      managed_fields = jsonencode({
        "field" = "test_managed_fields"
      })
      resource_version = "test"
      self_link        = "test"
      # deletion_grace_period_seconds = 1 # this will always be null in the API
    }
    spec = {
      string  = "test"
      bool    = false
      number  = 1
      int32   = 1
      int64   = 1
      float32 = 1.0
      float64 = 1.0
      map     = { "test" = "test" }
      list    = ["test"]
      object = {
        string_attribute = "test"
      }
      object_list = [
        {
          string_attribute = "test"
        }
      ]
      object_map = {
        "test" = {
          string_attribute = "test"
        }
      }
    }
    status = {
      string = "test"
    }
  }
}
