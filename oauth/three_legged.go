package oauth

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)


// NewThreeLegged returns a 3-legged authenticator with default host and authPath,
// giving client secrets, redirectURI and optionally with a starting refresh token (useful for CLI apps)
func NewThreeLegged(clientID, clientSecret, redirectURI, refreshToken string) *ThreeLeggedAuth {
	return &ThreeLeggedAuth{
		AuthData{
			clientID,
			clientSecret,
			"https://developer.api.autodesk.com",
			"/authentication/v1",
		},
		redirectURI,
		refreshToken,
	}
}

// Authorize method returns an URL to redirect an end user, where it will be asked to give his consent for app to
//access the specified resources.
//
// The resources for which the permission is asked are specified as a space-separated list of required scopes.
// State can be used to specify, as URL-encoded payload, some arbitrary data that the authentication flow will pass back
// verbatim in a state query parameter to the callback URL.
//	Note: You do not call this URL directly in your server code.
//	See the Get a 3-Legged Token tutorial for more information on how to use this endpoint.
func (a ThreeLeggedAuth) Authorize(scope string, state string) (string, error) {

	request, err := http.NewRequest("GET",
		a.Host+a.authPath+"/authorize",
		nil,
	)

	if err != nil {
		return "", err
	}

	query := request.URL.Query()
	query.Add("client_id", a.ClientID)
	query.Add("response_type", "code")
	query.Add("redirect_uri", a.RedirectURI)
	query.Add("scope", scope)
	query.Add("state", state)

	request.URL.RawQuery = query.Encode()

	return request.URL.String(), nil
}

func (a *ThreeLeggedAuth) SetRefreshToken(refreshtoken string) {
	a.RefreshToken = refreshtoken
}

//ExchangeCode is used to exchange the authorization code for a token and an exchange token
func (a *ThreeLeggedAuth) ExchangeCode(code string) (bearer Bearer, err error) {

	task := http.Client{}

	body := url.Values{}
	body.Add("client_id", a.ClientID)
	body.Add("client_secret", a.ClientSecret)
	body.Add("grant_type", "authorization_code")
	body.Add("code", code)
	body.Add("redirect_uri", a.RedirectURI)

	req, err := http.NewRequest("POST",
		a.Host+a.authPath+"/gettoken",
		bytes.NewBufferString(body.Encode()),
	)

	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := task.Do(req)

	if err != nil {
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		content, _ := ioutil.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&bearer)

	a.RefreshToken = bearer.RefreshToken

	return
}

func (a *ThreeLeggedAuth) GetToken(scope string) (token Bearer, err error) {
	token, err = a.GetNewRefreshToken(a.RefreshToken, scope)
	a.RefreshToken = token.RefreshToken
	return
}

// GetNewRefreshToken is used to get a new access token by using the refresh token provided by ExchangeCode
func (a ThreeLeggedAuth) GetNewRefreshToken(refreshToken string, scope string) (bearer Bearer, err error) {

	task := http.Client{}

	body := url.Values{}
	body.Add("client_id", a.ClientID)
	body.Add("client_secret", a.ClientSecret)
	body.Add("grant_type", "refresh_token")
	body.Add("refresh_token", refreshToken)
	body.Add("scope", scope)

	req, err := http.NewRequest("POST",
		a.Host+a.authPath+"/refreshtoken",
		bytes.NewBufferString(body.Encode()),
	)

	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := task.Do(req)

	if err != nil {
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		content, _ := ioutil.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&bearer)

	return
}


func (a ThreeLeggedAuth) GetRefreshToken() string {
	return a.RefreshToken
}