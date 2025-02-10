// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccStoreDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccStoreDataSourceConfig("store-1"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.openfga_store.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.openfga_store.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("store-1"),
					),
				},
			},
		},
	})
}

func testAccStoreDataSourceConfig(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "openfga_store" "test" {
  name = %[2]q
}

data "openfga_store" "test" {
  id = openfga_store.test.id
}
`, providerConfig, name)
}
