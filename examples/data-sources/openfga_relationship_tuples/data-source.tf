data "openfga_relationship_tuples" "all" {
  store_id = "example_store_id"
}

data "openfga_relationship_tuples" "query" {
  store_id = "example_store_id"

  query = {
    user     = "user:user-1"
    relation = "viewer"
    object   = "document:"
  }
}
