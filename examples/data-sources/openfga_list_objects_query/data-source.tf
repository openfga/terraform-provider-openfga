data "openfga_list_objects_query" "basic" {
  store_id = "example_store_id"

  user     = "user:user-1"
  relation = "viewer"
  type     = "document"
}

data "openfga_list_objects_query" "advanced" {
  store_id = "example_store_id"

  user     = "user:user-1"
  relation = "viewer"
  type     = "document"

  contextual_tuples = [
    {
      user     = "user:user-1"
      relation = "viewer"
      object   = "document:document-1"
    }
  ]

  context_json = jsonencode({
    time = timestamp()
  })
}
