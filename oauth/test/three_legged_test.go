package oauth_test

import (
	"github.com/apprentice3d/forge-api-go-client/oauth"
	"os"
	"testing"
)

func TestThreeLeggedAuthentication(t *testing.T) {

	//prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	t.Log("ClientID used in test: ", clientID)

	if len(clientID) == 0 || len(clientSecret) == 0 {
		t.Fatal("Could not get Forge env vars")
	}

	client := oauth.NewThreeLegged(clientID,
		clientSecret,
		"http://localhost:3009/callback", "")

	authCode := ""

	t.Run("Get the authorisation link", func(t *testing.T) {

		authLink, err := client.Authorize("data:read data:write user-profile:read", "something that will be passed back")

		if err != nil {
			t.Fatalf("Could not create the authorization link, got: %s", err.Error())
		}

		if len(authLink) == 0 {
			t.Fatalf("The authorization link is empty")
		}
		t.Log(authLink)
	})

	if authCode != "" {
		t.Run("Exchange auth code", func(t *testing.T) {
			token, err := client.ExchangeCode(authCode)

			if err != nil {
				t.Fatal("Could not exchange auth code for token: ", err.Error())
			}
			client.SetRefreshToken(token.RefreshToken)

			t.Log("Latest refresh token: ", token.RefreshToken)
		})

		t.Run("Get new token", func(t *testing.T) {
			token, err := client.GetToken("data:read")

			if err != nil {
				t.Fatal("Could not get new token: ", err.Error())
			}
			client.SetRefreshToken(token.RefreshToken)
			t.Log("Latest refresh token: ", token.RefreshToken)
		})
	}



}

func TestThreeLeggedAuthWithRefreshToken(t *testing.T) {

	//prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	t.Log("ClientID used in test: ", clientID)

	if len(clientID) == 0 || len(clientSecret) == 0 {
		t.Fatal("Could not get Forge env vars")
	}

	refreshToken := ""

	client := oauth.NewThreeLegged(clientID, clientSecret,
		"http://localhost:3009/callback", refreshToken)

	if refreshToken != "" {
		t.Run("Get new token", func(t *testing.T) {
			token, err := client.GetToken("data:read")

			if err != nil {
				t.Fatal("Could not get new token: ", err.Error())
			}
			client.SetRefreshToken(token.RefreshToken)
			t.Log("Latest refresh token: ", token.RefreshToken)
		})

		t.Run("Get another token", func(t *testing.T) {
			token, err := client.GetToken("data:write")

			if err != nil {
				t.Fatal("Could not get new token: ", err.Error())
			}
			client.SetRefreshToken(token.RefreshToken)
			t.Log("Latest refresh token: ", token.RefreshToken)
		})

		t.Run("Get token with wrong scope", func(t *testing.T) {
			token, err := client.GetToken("account:read")

			if err == nil {
				t.Fatal("Getting a token with superset scope should have failed")
			}
			client.SetRefreshToken(token.RefreshToken)
			t.Log("Latest refresh token: ", token.RefreshToken)
		})
	}



}
