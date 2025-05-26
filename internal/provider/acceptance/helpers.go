package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/openfga/terraform-provider-openfga/internal/provider"
)

const (
	ProviderApiUrl = "http://localhost:8080"
)

var (
	ProviderConfig = fmt.Sprintf(`
provider "openfga" {
	api_url = %[1]q
}
`, ProviderApiUrl)
)

var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"openfga": providerserver.NewProtocol6WithError(provider.New("test")()),
}

func TestAccPreCheck(t *testing.T) {}
