package oauth_test

import (
	"os"
	"testing"

	"github.com/apprentice3d/forge-api-go-client/oauth"
)

func TestThreeLeggedAuth_Authorize(t *testing.T) {

	//prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		t.Skipf("No Forge credentials present; skipping test")
	}

	if len(clientID) == 0 || len(clientSecret) == 0 {
		t.Fatal("Could not get Forge env vars")
	}

	client := oauth.NewThreeLeggedClient(clientID,
		clientSecret,
		"http://localhost:3009/callback")

	authLink, err := client.Authorize("data:read data:write", "something that will be passed back")

	if err != nil {
		t.Errorf("Could not create the authorization link, got: %s", err.Error())
	}

	if len(authLink) == 0 {
		t.Errorf("The authorization link is empty")
	}

}
