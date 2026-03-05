
resource "rancher2_dev" "full" {
  id                = "full_test"
  user_token        = "test"
  bool_attribute    = false
  number_attribute  = 1.1
  int64_attribute   = 1
  int32_attribute   = 1
  float64_attribute = 1.2
  float32_attribute = 1.3
  string_attribute  = "dev-test"
  list_attribute    = ["this", "is", "a", "list"]
  set_attribute     = toset(["this", "is", "a", "list"])
  map_attribute = {
    "this" = "is"
    "a"    = "map"
  }
  nested_object = {
    string_attribute = "test"
    nested_nested_object = {
      string_attribute = "tst"
      bool_attribute   = false
    }
  }
  nested_object_list = [
    {
      string_attribute = "test"
      nested_nested_object = {
        string_attribute = "tst"
        bool_attribute   = false
      }
    },
  ]
  nested_object_map = {
    "first" = {
      string_attribute = "test"
      nested_nested_object = {
        string_attribute = "tst"
        bool_attribute   = false
      }
    }
  }
}

resource "rancher2_dev" "required" {
  id               = "required_test"
  number_attribute = 1.1
  string_attribute = "dev-test"
}

