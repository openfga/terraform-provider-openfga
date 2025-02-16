# No authentication
provider "openfga" {
  api_url = "http://localhost:8080"
}

# Pre-shared key/ token authentication
provider "openfga" {
  api_url = "http://localhost:8080"

  api_token = var.openfga_api_token
}

# OIDC authentication
provider "openfga" {
  api_url = "http://localhost:8080"

  client_id          = var.openfga_client_id
  client_secret      = var.openfga_client_secret
  token_endpoint_url = var.openfga_token_endpoint_url
}
