data "openfga_relationship_tuple" "example" {
  store_id = "example_store_id"

  user     = "user:user-1"
  relation = "viewer"
  object   = "document:document-1"
}
