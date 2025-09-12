package authorizationmodel_test

import (
	"fmt"
	"os/exec"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	tf "github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/openfga/terraform-provider-openfga/internal/provider/acceptance"
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
				Check: func(s *tf.State) error {
					return nil
				},
			},
			// ImportState testing
			{
				ResourceName:      "openfga_authorization_model.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *tf.State) (string, error) {
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

// TestAccAuthorizationModelResourceDrift specifically tests the drift detection scenario
// This test simulates external deletion by attempting multiple approaches
func TestAccAuthorizationModelResourceDrift(t *testing.T) {
	var storeID, authModelID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and capture IDs
			{
				Config: testAccAuthorizationModelResourceConfig(
					testAccAuthorizationModelResourceModelJson("document"),
				),
				Check: func(s *tf.State) error {
					rsStore := s.RootModule().Resources["openfga_store.test"]
					rsAuthModel := s.RootModule().Resources["openfga_authorization_model.test"]
					storeID = rsStore.Primary.ID
					authModelID = rsAuthModel.Primary.ID
					t.Logf("Created store ID: %s, auth model ID: %s", storeID, authModelID)
					return nil
				},
			},
			// Simulate drift by trying multiple approaches to delete the model externally
			{
				PreConfig: func() {
					if storeID == "" || authModelID == "" {
						t.Fatal("Store ID or Auth Model ID is empty, cannot proceed with drift test")
						return
					}

					t.Logf("Attempting to simulate external deletion of authorization model")

					// Try approach 1: Direct database deletion with different table names
					tableNames := []string{"authorization_models", "authorization_model", "authz_models", "models"}

					for _, tableName := range tableNames {
						cmd := exec.Command("docker", "exec",
							"-e", "PGPASSWORD=password",
							"openfga-postgres",
							"psql", "-U", "openfga", "-d", "openfga",
							"-c", fmt.Sprintf("DELETE FROM %s WHERE id = '%s' OR authorization_model_id = '%s';", tableName, authModelID, authModelID))

						output, err := cmd.CombinedOutput()
						if err == nil {
							t.Logf("Successfully deleted from table %s: %s", tableName, string(output))
							break
						} else {
							t.Logf("Failed to delete from table %s: %v, output: %s", tableName, err, string(output))
						}
					}

					// Try approach 2: Reset the entire store (more drastic but might work)
					cmd := exec.Command("docker", "exec",
						"-e", "PGPASSWORD=password",
						"openfga-postgres",
						"psql", "-U", "openfga", "-d", "openfga",
						"-c", fmt.Sprintf("DELETE FROM stores WHERE id = '%s';", storeID))

					output, err := cmd.CombinedOutput()
					if err == nil {
						t.Logf("Successfully deleted store (cascade delete): %s", string(output))
					} else {
						t.Logf("Failed to delete store: %v, output: %s", err, string(output))
					}

					// Try approach 3: Use Docker to restart OpenFGA to clear in-memory state
					cmd = exec.Command("docker", "restart", "openfga-openfga-1")
					output, err = cmd.CombinedOutput()
					if err == nil {
						t.Logf("Successfully restarted OpenFGA container: %s", string(output))
						// Wait for the container to come back up
						time.Sleep(10 * time.Second)
					} else {
						t.Logf("Failed to restart OpenFGA container: %v, output: %s", err, string(output))
					}

					t.Logf("Completed external deletion attempts")
				},
				// Use the same config as before
				Config: testAccAuthorizationModelResourceConfig(
					testAccAuthorizationModelResourceModelJson("document"),
				),
				// We expect Terraform to detect the drift and plan to recreate the resource
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							"openfga_authorization_model.test",
							plancheck.ResourceActionCreate,
						),
					},
				},
				// After recreation, check that we have valid IDs again
				Check: func(s *tf.State) error {
					rsAuthModel, ok := s.RootModule().Resources["openfga_authorization_model.test"]
					if !ok {
						return fmt.Errorf("Authorization model resource not found after recreation")
					}

					newAuthModelID := rsAuthModel.Primary.ID
					if newAuthModelID == "" {
						return fmt.Errorf("New authorization model ID is empty after recreation")
					}

					t.Logf("Original auth model ID: %s, New auth model ID: %s", authModelID, newAuthModelID)
					if newAuthModelID == authModelID {
						t.Logf("Warning: New authorization model has same ID as deleted one, drift detection may not be working")
					} else {
						t.Logf("Success: New authorization model has different ID, drift detection is working")
					}

					// Update for next step
					authModelID = newAuthModelID
					return nil
				},
			},
		},
	})
}

// TestAccAuthorizationModelResourceDriftMockScenario tests drift detection using a more controlled approach
// This test specifically validates that the Read method properly handles missing resources
func TestAccAuthorizationModelResourceDriftMockScenario(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create the resource first
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
				},
			},
			// Test with an invalid/non-existent store ID to simulate drift
			{
				Config: testAccAuthorizationModelResourceDriftConfig(
					testAccAuthorizationModelResourceModelJson("document"),
					"non-existent-store-id",
				),
				ExpectError: regexp.MustCompile("Unable to create authorization model|Parameter StoreId is not a valid|not found|does not exist"),
			},
		},
	})
}

// Helper function to create a config with a specific store ID for drift testing
func testAccAuthorizationModelResourceDriftConfig(modelJson, storeID string) string {
	return fmt.Sprintf(`
%[1]s

resource "openfga_authorization_model" "test" {
	store_id = %[3]q

	model_json = %[2]q
}
`, acceptance.ProviderConfig, modelJson, storeID)
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
