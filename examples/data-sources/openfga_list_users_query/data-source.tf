data "openfga_list_users_query" "basic" {
  store_id = "example_store_id"

  type     = "user"
  relation = "viewer"
  object   = "document:document-1"
}

data "openfga_list_users_query" "advanced" {
  store_id = "example_store_id"

  type     = "user"
  relation = "viewer"
  object   = "document:document-1"

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
