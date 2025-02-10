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

func TestAccStoresDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test Empty
			{
				Config: testAccStoresDataSourceConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.openfga_stores.test",
						tfjsonpath.New("stores"),
						knownvalue.ListSizeExact(0),
					),
				},
			},
			// Setup stores
			{
				Config: testAccStoresDataSourceConfig("store-1", "store-2"),
			},
			// Read testing
			{
				Config: testAccStoresDataSourceConfig("store-1", "store-2"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.openfga_stores.test",
						tfjsonpath.New("stores"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"id":   knownvalue.NotNull(),
								"name": knownvalue.StringExact("store-1"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"id":   knownvalue.NotNull(),
								"name": knownvalue.StringExact("store-2"),
							}),
						}),
					),
				},
			},
		},
	})
}

func testAccStoresDataSourceConfig(names ...string) string {
	var resources string
	for idx, name := range names {
		var dependsOn string
		if idx > 0 {
			dependsOn = fmt.Sprintf(`depends_on = [openfga_store.store_%[1]d]`, idx-1)
		}

		resources += fmt.Sprintf(`
resource "openfga_store" "store_%[1]d" {
  name = %[2]q
  %[3]s
}
`, idx, name, dependsOn)
	}

	return fmt.Sprintf(`
%[1]s

%[2]s

data "openfga_stores" "test" {}
`, providerConfig, resources)
}
