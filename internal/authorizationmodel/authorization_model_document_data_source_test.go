package authorizationmodel_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/mauriceackel/terraform-provider-openfga/internal/acceptance"
)

func TestAccAuthorizationModelDocumentDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test DSL
			{
				Config: testAccAuthorizationModelDocumentDataSourceConfigDsl(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.openfga_authorization_model_document.test",
						tfjsonpath.New("result"),
						knownvalue.StringExact(expectedAuthorizationModelDocumentDataSourcResult),
					),
				},
			},
			// Test JSON
			{
				Config: testAccAuthorizationModelDocumentDataSourceConfigJson(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.openfga_authorization_model_document.test",
						tfjsonpath.New("result"),
						knownvalue.StringExact(expectedAuthorizationModelDocumentDataSourcResult),
					),
				},
			},
		},
	})
}

const expectedAuthorizationModelDocumentDataSourcResult = `{"schema_version":"1.1","type_definitions":[{"type":"user"}]}`

func testAccAuthorizationModelDocumentDataSourceConfigDsl() string {
	return fmt.Sprintf(`
%[1]s

data "openfga_authorization_model_document" "test" {
  dsl = <<EOT
model
  schema 1.1

type user
  EOT
}
`, acceptance.ProviderConfig)
}

func testAccAuthorizationModelDocumentDataSourceConfigJson() string {
	return fmt.Sprintf(`
%[1]s

data "openfga_authorization_model_document" "test" {
  json = <<EOT
{
  "type_definitions": [
    { "type": "user" }
  ],
  "schema_version": "1.1"
}
  EOT
}
`, acceptance.ProviderConfig)
}
