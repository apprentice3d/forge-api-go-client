package oauth_test

import (
	"os"
	"testing"

	"github.com/woweh/forge-api-go-client/oauth"
)

func TestTwoLeggedAuthentication(t *testing.T) {

	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	if len(clientID) == 0 || len(clientSecret) == 0 {
		t.Fatalf("Could not get from env the Forge secrets")
	}

	t.Run("Valid Forge Secrets", func(t *testing.T) {
		authenticator := oauth.NewTwoLegged(clientID, clientSecret)

		bearer, err := authenticator.GetToken("data:read")

		if err != nil {
			t.Error(err.Error())
		}

		if len(bearer.AccessToken) == 0 {
			t.Errorf("Wrong bearer content: %v", bearer)
		}
	})

	t.Run("Invalid Forge Secrets", func(t *testing.T) {
		authenticator := oauth.NewTwoLegged("", clientSecret)

		bearer, err := authenticator.GetToken("data:read")

		if err == nil {
			t.Errorf("Expected to fail due to wrong credentials, but got %v", bearer)
		}

		if len(bearer.AccessToken) != 0 {
			t.Errorf("expected to not receive a token, but received: %s", bearer.AccessToken)
		}
	})

	t.Run("Invalid scope", func(t *testing.T) {
		authenticator := oauth.NewTwoLegged(clientID, clientSecret)

		bearer, err := authenticator.GetToken("data:invalidScopeValue")

		if err == nil {
			t.Errorf("Expected to fail due to wrong scope, but got %v\n", bearer)
		}

		if len(bearer.AccessToken) != 0 {
			t.Errorf("expected to not receive a token, but received: %s", bearer.AccessToken)
		}
	})

	t.Run("Invalid or unreachable host", func(t *testing.T) {
		authenticator := oauth.NewTwoLegged(clientID, clientSecret)
		authenticator.Host = "http://localhost"

		bearer, err := authenticator.GetToken("data:read")

		if err == nil {
			t.Errorf("Expected to fail due to wrong host, but got %v\n", bearer)
		}

		if len(bearer.AccessToken) != 0 {
			t.Errorf("expected to not receive a token, but received: %s", bearer.AccessToken)
		}
	})
}
