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

// NewTwoLegged returns a 2-legged authenticator with default host and authPath
func NewTwoLegged(clientID, clientSecret string) *TwoLeggedAuth {
	return &TwoLeggedAuth {
		AuthData{
			clientID,
			clientSecret,
			"https://developer.api.autodesk.com",
			"/authentication/v1",
		},

	}
}

// GetToken allows getting a token with a given scope
func (a TwoLeggedAuth) GetToken(scope string) (bearer Bearer, err error) {

	task := http.Client{}

	body := url.Values{}
	body.Add("client_id", a.ClientID)
	body.Add("client_secret", a.ClientSecret)
	body.Add("grant_type", "client_credentials")
	body.Add("scope", scope)

	req, err := http.NewRequest("POST",
		a.Host+a.authPath+"/authenticate",
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


// GetHostPath returns host path, usually different in case of prd stg and dev environments
func (a AuthData) GetHostPath() string {
	return a.Host
}

// SetHostPath allows changing the host, usually useful for switching between prd stg and dev environments
func (a *AuthData) SetHostPath(host string) {
	a.Host = host
}