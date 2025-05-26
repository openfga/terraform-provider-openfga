package query_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/openfga/terraform-provider-openfga/internal/provider/acceptance"
)

func TestAccCheckQueryDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCheckQueryDataSourceConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.openfga_check_query.allowed",
						tfjsonpath.New("result"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"data.openfga_check_query.forbidden",
						tfjsonpath.New("result"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"data.openfga_check_query.contextually_allowed",
						tfjsonpath.New("result"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"data.openfga_check_query.contextually_allowed_with_context",
						tfjsonpath.New("result"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"data.openfga_check_query.contextually_forbidden_with_context",
						tfjsonpath.New("result"),
						knownvalue.Bool(false),
					),
				},
			},
		},
	})
}

func testAccCheckQueryDataSourceConfig() string {
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
	store_id               = openfga_store.test.id
	authorization_model_id = openfga_authorization_model.test.id

	user      = "user:user-1"
	relation  = "viewer"
	object    = "document:document-1"
}

data "openfga_check_query" "allowed" {
	depends_on = [openfga_relationship_tuple.test]

	store_id = openfga_store.test.id
	authorization_model_id = openfga_authorization_model.test.id

	user     = "user:user-1"
	relation = "viewer"
	object   = "document:document-1"
}

data "openfga_check_query" "forbidden" {
	depends_on = [openfga_relationship_tuple.test]

	store_id = openfga_store.test.id
	authorization_model_id = openfga_authorization_model.test.id

	user     = "user:user-2"
	relation = "viewer"
	object   = "document:document-1"
}

data "openfga_check_query" "contextually_allowed" {
	depends_on = [openfga_relationship_tuple.test]

	store_id = openfga_store.test.id
	authorization_model_id = openfga_authorization_model.test.id

	user     = "user:user-2"
	relation = "viewer"
	object   = "document:document-1"

	contextual_tuples = [{
		user     = "user:user-2"
		relation = "viewer"
		object   = "document:document-1"
	}]
}

data "openfga_check_query" "contextually_allowed_with_context" {
	depends_on = [openfga_relationship_tuple.test]

	store_id = openfga_store.test.id
	authorization_model_id = openfga_authorization_model.test.id

	user     = "user:user-2"
	relation = "viewer"
	object   = "document:document-1"
  
	contextual_tuples = [{
		user      = "user:user-2"
		relation  = "viewer"
		object    = "document:document-1"
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

data "openfga_check_query" "contextually_forbidden_with_context" {
	depends_on = [openfga_relationship_tuple.test]

	store_id = openfga_store.test.id
	authorization_model_id = openfga_authorization_model.test.id

	user     = "user:user-2"
	relation = "viewer"
	object   = "document:document-1"
  
	contextual_tuples = [{
		user      = "user:user-2"
		relation  = "viewer"
		object    = "document:document-1"
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
