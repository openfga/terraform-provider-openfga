data "openfga_authorization_model" "latest" {
  store_id = "example_store_id"
}

data "openfga_authorization_model" "specific" {
  store_id = "example_store_id"

  id = "example_authorization_model_id"
}
