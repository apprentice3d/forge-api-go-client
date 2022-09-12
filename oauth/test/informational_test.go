package oauth_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/woweh/forge-api-go-client/oauth"
)

//TODO: set up a pipeline for auto-creating a 3-legged oauth token
func TestInformation_AboutMe(t *testing.T) {

	//prepare the credentials
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	redirectURI := "http://localhost:3009/callback"
	refreshToken := ""

	authenticator := oauth.NewThreeLegged(clientID, clientSecret, redirectURI, refreshToken)

	info := oauth.NewInformationQuerier(authenticator)

	//aThreeLeggedToken := os.Getenv("THREE_LEGGED_TOKEN")

	profile, err := info.AboutMe()

	if err != nil {
		t.Error(err.Error())
		return
	}

	t.Logf("Received profile:\n"+
		"UserId: %s\n"+
		"UserName: %s\n"+
		"EmailId: %s\n"+
		"FirstName: %s\n"+
		"LastName: %s\n"+
		"EmailVerified: %t\n"+
		"Var2FaEnabled: %t\n"+
		"ProfileImages: %v",
		profile.UserID,
		profile.UserName,
		profile.EmailID,
		profile.FirstName,
		profile.LastName,
		profile.EmailVerified,
		profile.Var2FaEnabled,
		profile.ProfileImages)
}

func ExampleInformation_AboutMe() {

	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	authenticator := oauth.NewTwoLegged(clientID, clientSecret)
	//authenticator := oauth.NewThreeLegged("","", "")

	info := oauth.NewInformationQuerier(authenticator)

	profile, err := info.AboutMe()

	if err != nil {
		fmt.Printf("[ERROR] Could not retrieve profile, got %s\n", err.Error())
		return
	}

	fmt.Printf("Received profile:\n"+
		"UserId: %s\n"+
		"UserName: %s\n"+
		"EmailId: %s\n"+
		"FirstName: %s\n"+
		"LastName: %s\n"+
		"EmailVerified: %t\n"+
		"Var2FaEnabled: %t\n"+
		"ProfileImages: %v",
		profile.UserID,
		profile.UserName,
		profile.EmailID,
		profile.FirstName,
		profile.LastName,
		profile.EmailVerified,
		profile.Var2FaEnabled,
		profile.ProfileImages)
}
