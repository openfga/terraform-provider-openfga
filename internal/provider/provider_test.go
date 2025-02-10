package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	providerConfig = `
provider "openfga" {
  api_url = "http://localhost:8080"
}
`
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"openfga": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {}
