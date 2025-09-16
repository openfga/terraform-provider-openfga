package authorizationmodel_test

import (
	"fmt"
	"os/exec"
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

func TestAccAuthorizationModelResourceDriftMockScenario(t *testing.T) {
	var savedModelID, savedStoreID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAuthorizationModelResourceConfig(
					testAccAuthorizationModelResourceModelJson("document"),
				),
				Check: func(s *tf.State) error {
					store, ok := s.RootModule().Resources["openfga_store.test"]
					if !ok {
						return fmt.Errorf("resource openfga_store.test not found")
					}
					model, ok := s.RootModule().Resources["openfga_authorization_model.test"]
					if !ok {
						return fmt.Errorf("resource openfga_authorization_model.test not found")
					}

					savedStoreID = store.Primary.ID
					savedModelID = model.Primary.ID
					t.Logf("Saved store ID: %s, model ID: %s", savedStoreID, savedModelID)
					return nil
				},
			},
			// Simulate drift by deleting model in DB and restarting container.
			{
				PreConfig: func() {
					simulateDriftByDeletingModelAndRestartingContainer(t, savedStoreID, savedModelID)
				},
				Config: testAccAuthorizationModelResourceConfig(
					testAccAuthorizationModelResourceModelJson("document"),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							"openfga_authorization_model.test",
							plancheck.ResourceActionCreate,
						),
					},
				},
			},
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

// simulates resource drift by:
// 1. Deleting the authorization model directly from the database
// 2. Restarting the OpenFGA container to clear memory cache
// 3. Waiting for the API to be ready again.
func simulateDriftByDeletingModelAndRestartingContainer(t *testing.T, savedStoreID, savedModelID string) {
	// Delete model directly from database
	cmd := exec.Command(
		"docker", "exec", "openfga-postgres",
		"psql", "-U", "openfga", "-d", "openfga",
		"-c", fmt.Sprintf(
			"DELETE FROM authorization_model WHERE store = '%s' AND authorization_model_id = '%s';",
			savedStoreID, savedModelID,
		),
	)
	cmd.Env = append(cmd.Env, "PGPASSWORD=password")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Logf("Failed to delete model: %v, output: %s", err, string(output))
	} else {
		t.Logf("Deleted model: %s", string(output))
	}

	fmt.Printf("deleting model %s", savedModelID)

	// Restart OpenFGA container
	// We need to restart the container because OpenFGA caches authorization models in memory.
	// Even after deleting the model directly in the database, the API might still return the cached
	// model unless the service is restarted. Restarting ensures that the next Terraform apply will
	// see the resource as missing and trigger drift detection correctly.
	restartCmd := exec.Command("docker", "restart", "docker-openfga-1") // adjust container name
	if output, err := restartCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to restart container: %v, output: %s", err, string(output))
	} else {
		t.Logf("Container restarted: %s", string(output))
	}

	// Wait for API to be ready
	ready := false
	for i := 0; i < 5; i++ {
		curlCmd := exec.Command("curl", "-s", "-o", "/dev/null", "-w", "%{http_code}", "http://localhost:8080/healthz")
		status, err := curlCmd.CombinedOutput()
		if err == nil && string(status) == "200" {
			ready = true
			break
		}
		t.Logf("Waiting for API, status: %s, err: %v", string(status), err)
		time.Sleep(time.Second)
	}
	if !ready {
		t.Fatalf("OpenFGA API not ready after restart")
	}
}
