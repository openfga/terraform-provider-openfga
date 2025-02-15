package query_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/mauriceackel/terraform-provider-openfga/internal/acceptance"
)

func TestAccListObjectsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccListObjectsDataSourceConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.openfga_list_objects_query.with_results",
						tfjsonpath.New("result"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("document:dummy"),
						}),
					),
					statecheck.ExpectKnownValue(
						"data.openfga_list_objects_query.without_results",
						tfjsonpath.New("result"),
						knownvalue.ListSizeExact(0),
					),
					statecheck.ExpectKnownValue(
						"data.openfga_list_objects_query.with_contextual_results",
						tfjsonpath.New("result"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("document:dummy"),
						}),
					),
					statecheck.ExpectKnownValue(
						"data.openfga_list_objects_query.with_contextual_context_results",
						tfjsonpath.New("result"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("document:dummy"),
						}),
					),
					statecheck.ExpectKnownValue(
						"data.openfga_list_objects_query.without_contextual_context_results",
						tfjsonpath.New("result"),
						knownvalue.ListSizeExact(0),
					),
				},
			},
		},
	})
}

func testAccListObjectsDataSourceConfig() string {
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
    define viewer: [user, user with larger_than]

condition larger_than(required: int, provided: int) {
  provided > required
}
  EOT
}

resource "openfga_authorization_model" "test" {
  store_id = openfga_store.test.id

  model_json = data.openfga_authorization_model_document.test.result
}

resource "openfga_relationship_tuple" "test" {
  store_id = openfga_store.test.id

  user      = "user:user-1"
  relation  = "viewer"
  object    = "document:dummy"

  depends_on = [openfga_authorization_model.test]
}

data "openfga_list_objects_query" "with_results" {
  depends_on = [openfga_relationship_tuple.test]

  store_id = openfga_store.test.id
  authorization_model_id = openfga_authorization_model.test.id

  user     = "user:user-1"
  relation = "viewer"
  type     = "document"
}

data "openfga_list_objects_query" "without_results" {
  depends_on = [openfga_relationship_tuple.test]

  store_id = openfga_store.test.id
  authorization_model_id = openfga_authorization_model.test.id

  user     = "user:user-2"
  relation = "viewer"
  type     = "document"
}

data "openfga_list_objects_query" "with_contextual_results" {
  depends_on = [openfga_relationship_tuple.test]

  store_id = openfga_store.test.id
  authorization_model_id = openfga_authorization_model.test.id

  user     = "user:user-2"
  relation = "viewer"
  type     = "document"
  
  contextual_tuples = [{
    user     = "user:user-2"
    relation = "viewer"
    object   = "document:dummy"
  }]
}

data "openfga_list_objects_query" "with_contextual_context_results" {
  depends_on = [openfga_relationship_tuple.test]

  store_id = openfga_store.test.id
  authorization_model_id = openfga_authorization_model.test.id

  user     = "user:user-2"
  relation = "viewer"
  type     = "document"
  
  contextual_tuples = [{
    user      = "user:user-2"
    relation  = "viewer"
    object    = "document:dummy"
	condition = {
      name = "larger_than"
      context_json = jsonencode({
	    provided = 100
      })
    }
  }]

  context_json = jsonencode({
    required = 50
  })
}

data "openfga_list_objects_query" "without_contextual_context_results" {
  depends_on = [openfga_relationship_tuple.test]

  store_id = openfga_store.test.id
  authorization_model_id = openfga_authorization_model.test.id

  user     = "user:user-2"
  relation = "viewer"
  type     = "document"
  
  contextual_tuples = [{
    user      = "user:user-2"
    relation  = "viewer"
    object    = "document:dummy"
	condition = {
      name = "larger_than"
      context_json = jsonencode({
	    provided = 100
      })
    }
  }]

  context_json = jsonencode({
    required = 100
  })
}
`, acceptance.ProviderConfig)
}
