package dm_test

import (
	"os"
	"testing"
	"../dm"
	// "github.com/outer-labs/forge-api-go-client/dm"
)

func TestProjectAPI_GetProjects(t *testing.T) {

	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	hubAPI := dm.NewHubAPIWithCredentials(clientID, clientSecret)

	// testHubKey := "my_test_hub_key_for_go"
	testHubKey := os.Getenv("BIM_360_TEST_ACCOUNT_HUBKEY")

	t.Run("Get project details", func(t *testing.T) {
		_, err := hubAPI.GetProjectDetails(testHubKey)

		if err != nil {
			t.Fatalf("Failed to get project details: %s\n", err.Error())
		}
	})

	t.Run("Get nonexistent project", func(t *testing.T) {
		_, err := hubAPI.GetProjectDetails(testHubKey + "30091981")

		if err == nil {
			t.Fatalf("Should fail getting getting details for non-existing project\n")
		}
	})
}


