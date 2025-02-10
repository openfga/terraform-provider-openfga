// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	providerApiUrl = "http://localhost:8080"
)

var (
	providerConfig = fmt.Sprintf(`
provider "openfga" {
  api_url = %[1]q
}
`, providerApiUrl)
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"openfga": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {}
