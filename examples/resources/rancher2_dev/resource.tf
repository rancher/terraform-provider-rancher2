
resource "rancher2_dev" "full" {
  # id                = "full_test" # this is read only
  # int32_attribute   = 1           # this is read only
  string_attribute  = "dev-test" # required
  number_attribute  = 1.1        # required
  bool_attribute    = false
  int64_attribute   = 1
  float64_attribute = 1.2
  float32_attribute = 1.3
  list_attribute    = ["this", "is", "a", "list"]
  set_attribute     = toset(["this", "is", "a", "set"])
  map_attribute = {
    "this" = "is"
    "a"    = "map"
  }
  nested_object = {
    string_attribute = "test"
    nested_nested_object = {
      string_attribute = "tst"
      # bool_attribute   = false # read only
    }
  }
  nested_object_list = [
    {
      string_attribute = "test"
      nested_nested_object = {
        string_attribute = "tst"
        # bool_attribute   = false # read only
      }
    },
  ]
  nested_object_map = {
    "first" = {
      string_attribute = "test"
      nested_nested_object = {
        string_attribute = "tst"
        # bool_attribute   = false # read only
      }
    }
  }
}

# resource "rancher2_dev" "required" {
#   number_attribute = 1.1
#   string_attribute = "dev-test-required"
# }
