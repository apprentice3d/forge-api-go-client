package dm_test

import (
	"fmt"
	"github.com/outer-labs/forge-api-go-client/dm"
	"log"
	"os"
	"testing"
	"net/http"
)

func TestHubAPI_GetHubDetails(t *testing.T) {

	// prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	hubAPI := dm.NewHubAPIWithCredentials(clientID, clientSecret)

	testHubKey := "my_test_hub_key_for_go"

	t.Run("Create a hub", func(t *testing.T) {
		_, err := hubAPI.CreateHub(testHubKey, "transient")

		if err != nil {
			t.Fatalf("Failed to create a hub: %s\n", err.Error())
		}
	})

	t.Run("Get hub details", func(t *testing.T) {
		_, err := hubAPI.GetHubDetails(testHubKey)

		if err != nil {
			t.Fatalf("Failed to get hub details: %s\n", err.Error())
		}
	})

	t.Run("Delete created hub", func(t *testing.T) {
		err := hubAPI.DeleteHub(testHubKey)

		if err != nil {
			t.Fatalf("Failed to delete hub: %s\n", err.Error())
		}
	})

	t.Run("Get nonexistent hub", func(t *testing.T) {
		_, err := hubAPI.GetHubDetails(testHubKey + "30091981")

		if err == nil {
			t.Fatalf("Should fail getting getting details for non-existing hub\n")
		}
	})
}


