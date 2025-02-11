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

func TestAccAuthorizationModelDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccAuthorizationModelDataSourceConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.openfga_authorization_model.specific",
						tfjsonpath.New("model_json"),
						knownvalue.StringExact(expectedFirstAuthorizationModelDataSourceModelJson),
					),
					statecheck.ExpectKnownValue(
						"data.openfga_authorization_model.latest",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.openfga_authorization_model.latest",
						tfjsonpath.New("model_json"),
						knownvalue.StringExact(expectedLatestAuthorizationModelDataSourceModelJson),
					),
				},
			},
		},
	})
}

const expectedFirstAuthorizationModelDataSourceModelJson = `{"schema_version":"1.1","type_definitions":[{"type":"document"}]}`
const expectedLatestAuthorizationModelDataSourceModelJson = `{"schema_version":"1.1","type_definitions":[{"type":"file"}]}`

func testAccAuthorizationModelDataSourceConfig() string {
	return fmt.Sprintf(`
%[1]s

resource "openfga_store" "test" {
  name = "test"
}

resource "openfga_authorization_model" "first" {
  store_id = openfga_store.test.id
 
  model_json = %[2]q
}

resource "openfga_authorization_model" "latest" {
  store_id = openfga_store.test.id
 
  model_json = %[3]q

  depends_on = [openfga_authorization_model.first]
}

data "openfga_authorization_model" "specific" {
  id       = openfga_authorization_model.first.id
  store_id = openfga_store.test.id
}

data "openfga_authorization_model" "latest" {
  store_id = openfga_store.test.id

  depends_on = [openfga_authorization_model.latest]
}

data "openfga_store" "test" {
  id = openfga_store.test.id
}
`, acceptance.ProviderConfig, expectedFirstAuthorizationModelDataSourceModelJson, expectedLatestAuthorizationModelDataSourceModelJson)
}
