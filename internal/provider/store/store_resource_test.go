package store_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	tf "github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/openfga/terraform-provider-openfga/internal/provider/acceptance"
	"os/exec"
	"testing"
)

func TestAccStoreResource(t *testing.T) {
	var storeID string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
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
				Check: func(s *tf.State) error {
					// Capture the store ID for later use in drift testing
					rs := s.RootModule().Resources["openfga_store.test"]
					storeID = rs.Primary.ID
					return nil
				},
			},
			// ImportState testing
			{
				ResourceName:      "openfga_store.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Drift testing: delete externally, then plan and apply recreate
			{
				PreConfig: func() {
					if storeID != "" {
						cmd := exec.Command("curl", "-X", "DELETE", "http://localhost:8080/stores/"+storeID)
						err := cmd.Run()
						if err != nil {
							t.Fatal(err)
						}

					}
				},
				Config: testAccStoreResourceConfig("store-1"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							"openfga_store.test",
							plancheck.ResourceActionCreate,
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
						knownvalue.StringExact("store-1"),
					),
				},
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
