package authorizationmodel_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/mauriceackel/terraform-provider-openfga/internal/provider/acceptance"
)

func TestAccAuthorizationModelResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccAuthorizationModelResourceConfig(
					testAccAuthorizationModelResourceModelJson("document"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"openfga_authorization_model.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"openfga_authorization_model.test",
						tfjsonpath.New("model_json"),
						knownvalue.StringExact(
							testAccAuthorizationModelResourceModelJson("document"),
						),
					),
				},
			},
			// ImportState testing
			{
				ResourceName:      "openfga_authorization_model.test",
				ImportState:       true,
				ImportStateVerify: true,
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
						"%s/%s",
						store.Primary.Attributes["id"],
						authorizationModel.Primary.Attributes["id"],
					), nil
				},
			},
			// Update and Read testing
			{
				Config: testAccAuthorizationModelResourceConfig(
					testAccAuthorizationModelResourceModelJson("file"),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							"openfga_authorization_model.test",
							plancheck.ResourceActionDestroyBeforeCreate,
						),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"openfga_authorization_model.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"openfga_authorization_model.test",
						tfjsonpath.New("model_json"),
						knownvalue.StringExact(
							testAccAuthorizationModelResourceModelJson("file"),
						),
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccAuthorizationModelResourceModelJson(typeName string) string {
	return fmt.Sprintf(`{"conditions":{"non_expired_grant":{"expression":"current_time == grant_time + grant_duration","name":"non_expired_grant","parameters":{"current_time":{"generic_types":[],"type_name":"TYPE_NAME_TIMESTAMP"},"grant_duration":{"generic_types":[],"type_name":"TYPE_NAME_DURATION"},"grant_time":{"generic_types":[],"type_name":"TYPE_NAME_TIMESTAMP"}}}},"schema_version":"1.1","type_definitions":[{"relations":{},"type":"user"},{"metadata":{"module":"","relations":{"viewer":{"directly_related_user_types":[{"condition":"non_expired_grant","type":"user"}],"module":""}}},"relations":{"viewer":{"this":{}}},"type":%[1]q}]}`, typeName)
}

func testAccAuthorizationModelResourceConfig(modelJson string) string {
	return fmt.Sprintf(`
%[1]s

resource "openfga_store" "test" {
	name = "test"
}

resource "openfga_authorization_model" "test" {
	store_id = openfga_store.test.id

	model_json = %[2]q
}
`, acceptance.ProviderConfig, modelJson)
}
