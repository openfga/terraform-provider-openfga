package relationshiptuple_test

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/openfga/terraform-provider-openfga/internal/provider/acceptance"
)

func TestAccRelationshipTupleResource(t *testing.T) {
	var storeID string

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
						tfjsonpath.New("authorization_model_id"),
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
						knownvalue.StringExact("document:document-1"),
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
				Check: func(s *terraform.State) error {
					// Capture the store ID for later use in drift testing
					rs := s.RootModule().Resources["openfga_store.test"]
					storeID = rs.Primary.ID
					return nil
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

					authorizationModel, ok := s.RootModule().Resources["openfga_authorization_model.test"]
					if !ok {
						return "", fmt.Errorf("Unable to find resource openfga_authorization_model.test")
					}

					return fmt.Sprintf(
						"%s/%s/%s/%s/%s",
						store.Primary.Attributes["id"],
						authorizationModel.Primary.Attributes["id"],
						"user:user-1",
						"viewer",
						"document:document-1",
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
						knownvalue.StringExact("document:document-1"),
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
			// Drift testing: delete externally, then plan and apply recreate
			{
				PreConfig: func() {
					if storeID != "" {
						jsonBody := `{"deletes":{"tuple_keys":[{"user":"user:user-2","relation":"viewer","object":"document:document-1","condition":{"name":"non_expired_grant","context":{"grant_time":"2023-01-01T00:00:00Z","grant_duration":"10m"}}}]}}`
						cmd := exec.Command("curl", "-X", "POST", "-H", "Content-Type: application/json", "-d", jsonBody, "http://localhost:8080/stores/"+storeID+"/write")
						err := cmd.Run()
						if err != nil {
							t.Fatal(err)
						}
					}
				},
				Config: testAccRelationshipTupleResourceConfig("user-2"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							"openfga_relationship_tuple.test",
							plancheck.ResourceActionCreate,
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
						knownvalue.StringExact("document:document-1"),
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
	store_id               = openfga_store.test.id
	authorization_model_id = openfga_authorization_model.test.id

	user      = "user:%[2]s"
	relation  = "viewer"
	object    = "document:document-1"
	condition = {
		name         = "non_expired_grant"
		context_json = jsonencode({
			grant_time     = "2023-01-01T00:00:00Z"
			grant_duration = "10m"
		})
	}
}
`, acceptance.ProviderConfig, userName)
}
