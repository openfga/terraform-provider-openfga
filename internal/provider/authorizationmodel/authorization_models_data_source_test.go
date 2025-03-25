package authorizationmodel_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/mauriceackel/terraform-provider-openfga/internal/provider/acceptance"
)

func TestAccAuthorizationModelsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test Empty
			{
				Config: testAccAuthorizationModelsDataSourceConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.openfga_authorization_models.test",
						tfjsonpath.New("authorization_models"),
						knownvalue.ListSizeExact(0),
					),
				},
			},
			// Setup authorization models
			{
				Config: testAccAuthorizationModelsDataSourceConfig(
					testAccAuthorizationModelsDataSourceModelJson("document"),
					testAccAuthorizationModelsDataSourceModelJson("file"),
				),
			},
			// Read testing
			{
				Config: testAccAuthorizationModelsDataSourceConfig(
					testAccAuthorizationModelsDataSourceModelJson("document"),
					testAccAuthorizationModelsDataSourceModelJson("file"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.openfga_authorization_models.test",
						tfjsonpath.New("authorization_models"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"id": knownvalue.NotNull(),
								"model_json": knownvalue.StringExact(
									testAccAuthorizationModelsDataSourceModelJson("file"),
								),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"id": knownvalue.NotNull(),
								"model_json": knownvalue.StringExact(
									testAccAuthorizationModelsDataSourceModelJson("document"),
								),
							}),
						}),
					),
				},
			},
		},
	})
}

func testAccAuthorizationModelsDataSourceModelJson(typeName string) string {
	return fmt.Sprintf(`{"conditions":{},"schema_version":"1.1","type_definitions":[{"relations":{},"type":%[1]q}]}`, typeName)
}

func testAccAuthorizationModelsDataSourceConfig(modelJsons ...string) string {
	var resources = `
resource "openfga_store" "test" {
	name = "test"
}
	`

	for idx, modelJson := range modelJsons {
		var dependsOn string
		if idx > 0 {
			dependsOn = fmt.Sprintf(`depends_on = [openfga_authorization_model.model_%[1]d]`, idx-1)
		}

		resources += fmt.Sprintf(`
resource "openfga_authorization_model" "model_%[1]d" {
	store_id = openfga_store.test.id

	model_json = %[2]q
  %[3]s
}
`, idx, modelJson, dependsOn)
	}

	return fmt.Sprintf(`
%[1]s

%[2]s

data "openfga_authorization_models" "test" {
	store_id = openfga_store.test.id
}
`, acceptance.ProviderConfig, resources)
}
