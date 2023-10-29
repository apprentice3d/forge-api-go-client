package oauth

// ForgeAuthenticator defines an interface that allows abstraction of 2-legged and a 3-legged context.
//
//	This provides useful when an API accepts both 2-legged and 3-legged context tokens
type ForgeAuthenticator interface {
	GetToken(scope string) (Bearer, error)
	HostPath() string
	GetRefreshToken() string
}

// AuthData reflects the data common to 2-legged and 3-legged api calls
type AuthData struct {
	ClientID     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
	Host         string `json:"host,omitempty"`
	authPath     string
}

// TwoLeggedAuth struct holds data necessary for making requests in 2-legged context
type TwoLeggedAuth struct {
	AuthData
}

// ThreeLeggedAuth struct holds data necessary for making requests in 3-legged context
type ThreeLeggedAuth struct {
	AuthData
	RedirectURI  string
	RefreshToken string
}

// Bearer reflects the response when acquiring a 2-legged token or in 3-legged context for exchanging the authorization
// code for a token + refresh token and when exchanging the refresh token for a new token
type Bearer struct {
	TokenType    string `json:"token_type"`              // Will always be Bearer
	ExpiresIn    int32  `json:"expires_in"`              // Access token expiration time (in seconds)
	AccessToken  string `json:"access_token"`            // The access token
	RefreshToken string `json:"refresh_token,omitempty"` // The refresh token used in 3-legged oauth
}
