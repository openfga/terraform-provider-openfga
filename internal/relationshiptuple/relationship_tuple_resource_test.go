package relationshiptuple_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/mauriceackel/terraform-provider-openfga/internal/acceptance"
)

func TestAccRelationshipTupleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccRelationshipTupleResourceConfig("user-1"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"openfga_relationship_tuple.test",
						tfjsonpath.New("store_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"openfga_relationship_tuple.test",
						tfjsonpath.New("user"),
						knownvalue.StringExact("user:user-1"),
					),
					statecheck.ExpectKnownValue(
						"openfga_relationship_tuple.test",
						tfjsonpath.New("relation"),
						knownvalue.StringExact("viewer"),
					),
					statecheck.ExpectKnownValue(
						"openfga_relationship_tuple.test",
						tfjsonpath.New("object"),
						knownvalue.StringExact("document:dummy"),
					),
					statecheck.ExpectKnownValue(
						"openfga_relationship_tuple.test",
						tfjsonpath.New("condition"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"name":         knownvalue.StringExact("non_expired_grant"),
							"context_json": knownvalue.StringExact(`{"grant_duration":"10m","grant_time":"2023-01-01T00:00:00Z"}`),
						}),
					),
				},
			},
			// ImportState testing
			{
				ResourceName:                         "openfga_relationship_tuple.test",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "user",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					store, ok := s.RootModule().Resources["openfga_store.test"]
					if !ok {
						return "", fmt.Errorf("Unable to find resource openfga_store.test")
					}

					return fmt.Sprintf(
						"%s/%s/%s/%s",
						store.Primary.Attributes["id"],
						"user:user-1",
						"viewer",
						"document:dummy",
					), nil
				},
			},
			// Update and Read testing
			{
				Config: testAccRelationshipTupleResourceConfig("user-2"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							"openfga_relationship_tuple.test",
							plancheck.ResourceActionDestroyBeforeCreate,
						),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"openfga_relationship_tuple.test",
						tfjsonpath.New("store_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"openfga_relationship_tuple.test",
						tfjsonpath.New("user"),
						knownvalue.StringExact("user:user-2"),
					),
					statecheck.ExpectKnownValue(
						"openfga_relationship_tuple.test",
						tfjsonpath.New("relation"),
						knownvalue.StringExact("viewer"),
					),
					statecheck.ExpectKnownValue(
						"openfga_relationship_tuple.test",
						tfjsonpath.New("object"),
						knownvalue.StringExact("document:dummy"),
					),
					statecheck.ExpectKnownValue(
						"openfga_relationship_tuple.test",
						tfjsonpath.New("condition"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"name":         knownvalue.StringExact("non_expired_grant"),
							"context_json": knownvalue.StringExact(`{"grant_duration":"10m","grant_time":"2023-01-01T00:00:00Z"}`),
						}),
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccRelationshipTupleResourceConfig(userName string) string {
	return fmt.Sprintf(`
%[1]s

resource "openfga_store" "test" {
  name = "test"
}

data "openfga_authorization_model_document" "test" {
  dsl = <<EOT
model
  schema 1.1

type user

type document
  relations
    define viewer: [user with non_expired_grant]

condition non_expired_grant(current_time: timestamp, grant_time: timestamp, grant_duration: duration) {
  current_time < grant_time + grant_duration
}
  EOT
}

resource "openfga_authorization_model" "test" {
  store_id = openfga_store.test.id

  model_json = data.openfga_authorization_model_document.test.result
}

resource "openfga_relationship_tuple" "test" {
  store_id = openfga_store.test.id

  user      = "user:%[2]s"
  relation  = "viewer"
  object    = "document:dummy"
  condition = {
    name         = "non_expired_grant"
	context_json = jsonencode({
      grant_time     = "2023-01-01T00:00:00Z"
	  grant_duration = "10m"
    })
  }

  depends_on = [openfga_authorization_model.test]
}
`,
		acceptance.ProviderConfig,
		userName,
	)
}
