package oauth

import (
	"encoding/base64"
	"net/http"
)

// setBasicAuthHeader sets the Basic Authorization header for a given request
func setBasicAuthHeader(r *http.Request, a AuthData) {
	base64encodedClientIdAndSecret := base64.StdEncoding.EncodeToString([]byte(a.ClientID + ":" + a.ClientSecret))
	r.Header.Set("Authorization", "Basic "+base64encodedClientIdAndSecret)
}
