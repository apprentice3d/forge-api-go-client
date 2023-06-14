/*
Package oauth provides the Golang implementation for the APS Authentication (OAuth) V2 REST API.
https://aps.autodesk.com/en/docs/oauth/v2/developers_guide/overview/

The API supports the following features:
- Two-legged authentication
- Three-legged authentication
- Refreshing tokens (only for three-legged)

To-do:
- Update APIs:
  - get user profile (replaces informational)

- Add missing APIs:
  - get OIDC specs
  - get JWKS
  - logout
  - introspect token
  - revoke token

Example of two-legged authentication:

func ExampleTwoLeggedAuth_Authenticate() {

	// acquire Forge secrets from environment
	clientID := os.Getenv("FORGE_CLIENT_ID")
	clientSecret := os.Getenv("FORGE_CLIENT_SECRET")

	if len(clientID) == 0 || len(clientSecret) == 0 {
		log.Fatalf("Could not get from env the Forge secrets")
	}

	// create oauth client
	authenticator := oauth.NewTwoLegged(clientID, clientSecret)

	// request a token with needed scopes, separated by spaces
	bearer, err := authenticator.GetToken("data:read data:write")

	if err != nil || len(bearer.AccessToken) == 0 {
		log.Fatalf("Could not get from env the Forge secrets")
	}

	// at this point, the bearer should contain the needed data. Check Bearer struct for more info
	fmt.Printf("Bearer now contains:\n"+
		"AccessToken: %s\n"+
		"TokenType: %s\n"+
		"Expires in: %d\n",
		bearer.AccessToken,
		bearer.TokenType,
		bearer.ExpiresIn)

}
*/
package oauth
