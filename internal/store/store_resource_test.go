// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package store_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/mauriceackel/terraform-provider-openfga/internal/acceptance"
)

func TestAccStoreResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccStoreResourceConfig("store-1"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"openfga_store.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"openfga_store.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("store-1"),
					),
				},
			},
			// ImportState testing
			{
				ResourceName:      "openfga_store.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccStoreResourceConfig("store-2"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							"openfga_store.test",
							plancheck.ResourceActionDestroyBeforeCreate,
						),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"openfga_store.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"openfga_store.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("store-2"),
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccStoreResourceConfig(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "openfga_store" "test" {
  name = %[2]q
}
`, acceptance.ProviderConfig, name)
}
