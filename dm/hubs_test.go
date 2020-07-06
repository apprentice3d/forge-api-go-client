package dm_test

import (
	"os"
	"testing"

	"github.com/outer-labs/forge-api-go-client/dm"
)

func TestHubAPI_GetHubDetails(t *testing.T) {

	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		t.Skipf("No Forge credentials present; skipping test")
	}

	hubAPI := dm.NewHubAPIWithCredentials(clientID, clientSecret)

	// testHubKey := "my_test_hub_key_for_go"
	testHubKey := os.Getenv("BIM_360_TEST_ACCOUNT_HUBKEY")

	t.Run("Get hub details", func(t *testing.T) {
		_, err := hubAPI.GetHubDetails(testHubKey)

		if err != nil {
			t.Fatalf("Failed to get hub details: %s\n", err.Error())
		}
	})

	t.Run("Get nonexistent hub", func(t *testing.T) {
		_, err := hubAPI.GetHubDetails(testHubKey + "30091981")

		if err == nil {
			t.Fatalf("Should fail getting getting details for non-existing hub\n")
		}
	})
}
