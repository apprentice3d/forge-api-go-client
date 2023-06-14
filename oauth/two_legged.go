package oauth

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// NewTwoLegged returns a 2-legged authenticator with default host and authPath
func NewTwoLegged(clientID, clientSecret string) *TwoLeggedAuth {
	return &TwoLeggedAuth{
		AuthData{
			clientID,
			clientSecret,
			"https://developer.api.autodesk.com",
			"/authentication/v2",
		},
	}
}

// GetToken allows getting a token with a given scope
// References:
// - https://aps.autodesk.com/en/docs/oauth/v2/reference/http/gettoken-POST/ -
// - https://aps.autodesk.com/en/docs/oauth/v2/tutorials/get-2-legged-token/
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

// GetHostPath returns host path, usually different in case of prd stg and dev environments
func (a *AuthData) GetHostPath() string {
	return a.Host
}

// SetHostPath allows changing the host, usually useful for switching between prd stg and dev environments
func (a *AuthData) SetHostPath(host string) {
	a.Host = host
}
