resource "openfga_store" "example" {
  name = "example_store_name"
}

data "openfga_authorization_model_document" "example" {
  dsl = <<EOT
model
  schema 1.1

type user

type document
  relations
    define viewer: [user]
  EOT
}

resource "openfga_authorization_model" "example" {
  store_id = openfga_store.example.id

  model_json = data.openfga_authorization_model_document.example.result
}
