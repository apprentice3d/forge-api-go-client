package dm_test

import (
	"fmt"
	"forge-api-go-client/dm"
	"os"
	"testing"
)

func TestHubsAPI_ListHubs(t *testing.T) {
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	fmt.Printf("Using envs: %s\n%s\n", clientID, clientSecret)

	hubsAPI := dm.NewHubsAPIWithCredentials(clientID, clientSecret)

	t.Run("List Hubs", func(t *testing.T) {
		hubs, err := hubsAPI.ListHubs()
		if err!=nil{
			t.Fatalf("Failed to list hubs: %s\n", err.Error())
		}

		if len(hubs.Data) == 0 {
			t.Fatalf("Failed to list hubs. No hubs retreived.")
		}

		projects, err := hubsAPI.GetHubProjects(hubs.Data[0].Id)
		if err!=nil{
			t.Fatalf("Failed to list hub '%s' projects: %s\n", hubs.Data[0].Id, err.Error())
		}

		if len(projects.Data) == 0 {
			t.Fatalf("Failed to list hub projects. No projects retreived or failed to unmarshal.")
		}
	})
}