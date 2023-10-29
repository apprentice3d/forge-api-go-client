package oauth

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/woweh/forge-api-go-client"
)

// NewTwoLegged returns a 2-legged authenticator with default host and authPath
func NewTwoLegged(clientID, clientSecret string) *TwoLeggedAuth {
	return &TwoLeggedAuth{
		AuthData{
			clientID,
			clientSecret,
			forge.HostName,
			"/authentication/v2",
		},
	}
}

// GetToken allows getting a token with a given scope
//
// Parameter:
// - scope: a space separated list, like "data:read data:search viewables:read".
//
// References:
// - https://aps.autodesk.com/en/docs/oauth/v2/reference/http/gettoken-POST/ -
// - https://aps.autodesk.com/en/docs/oauth/v2/tutorials/get-2-legged-token/
// - https://aps.autodesk.com/en/docs/oauth/v2/developers_guide/scopes/
func (a *TwoLeggedAuth) GetToken(scope string) (bearer Bearer, err error) {

	task := http.Client{}

	body := url.Values{}
	body.Add("grant_type", "client_credentials")
	body.Add("scope", scope)

	req, err := http.NewRequest(
		"POST",
		a.Host+a.authPath+"/token",
		bytes.NewBufferString(body.Encode()),
	)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	setBasicAuthHeader(req, a.AuthData)
	response, err := task.Do(req)

	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		content, _ := io.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&bearer)

	return
}

func (a *TwoLeggedAuth) GetRefreshToken() string {
	return ""
}

// HostPath returns host path, usually different in case of prd stg and dev environments
// Note:
//   - This might be useful for Autodesk internal use, but not for external developers.
func (a *AuthData) HostPath() string {
	return a.Host
}

// SetHostPath allows changing the host, usually useful for switching between prd stg and dev environments
// Note:
//   - This might be useful for Autodesk internal use, but not for external developers.
func (a *AuthData) SetHostPath(host string) {
	a.Host = host
}
