data "openfga_authorization_model_document" "dsl" {
  dsl = file("path/to/model.fga")
}

data "openfga_authorization_model_document" "mod" {
  mod_file_path = "path/to/fga.mod"
}

data "openfga_authorization_model_document" "json" {
  json = file("path/to/model.json")
}

data "openfga_authorization_model_document" "model" {
  model = {
    schema_version = "1.1"
    type_definitions = [{
      type = "user"
    }]
  }
}
