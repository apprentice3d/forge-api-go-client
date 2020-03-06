package dm

import (
	// "fmt"
	"encoding/json"
	"net/http"
	"github.com/outer-labs/forge-api-go-client/oauth"
)

// HubAPI holds the necessary data for making calls to Forge Data Management service
type HubAPI struct {
	oauth.TwoLeggedAuth
	HubAPIPath string
}

var api HubAPI

// NewHubAPIWithCredentials returns a Hub API client with default configurations
func NewHubAPIWithCredentials(ClientID string, ClientSecret string) HubAPI {
	return HubAPI{
		oauth.NewTwoLeggedClient(ClientID, ClientSecret),
		"/project/v1/hubs",
	}
}

func (api HubAPI) GetHubs() (result ForgeResponseArray, err error) {
	bearer, err := api.Authenticate("data:read")
	if err != nil {
		return
	}
	path := api.Host + api.HubAPIPath

	return getHubs(path, bearer.AccessToken)
}

func (api HubAPI) GetHubDetails(hubKey string) (result ForgeResponseObject, err error) {
	bearer, err := api.Authenticate("data:read")
	if err != nil {
		return
	}
	path := api.Host + api.HubAPIPath

	return getHubDetails(path, hubKey, bearer.AccessToken)
}


func GetHubsThreeLegged(bearer oauth.Bearer) (result ForgeResponseArray, err error) {
	// bearer, err := api.Authenticate("data:read")
	// if err != nil {
	// 	return
	// }
	// path := api.Host + api.HubAPIPath

	//To do? check if access token needs to be refreshed? if so, run bearer.RefreshToken?
	path := "https://developer.api.autodesk.com/project/v1/hubs"
	return getHubs(path, bearer.AccessToken)
}




/*
 *	SUPPORT FUNCTIONS
 */

func getHubs(path, token string) (result ForgeResponseArray, err error) {
	task := http.Client{}

	req, err := http.NewRequest("GET",
		path,
		nil,
	)

	if err != nil {
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)

	response, err := task.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()
	
	decoder := json.NewDecoder(response.Body)
		if response.StatusCode != http.StatusOK {
			err = &ErrorResult{StatusCode:response.StatusCode}
			decoder.Decode(err)
				return
		}
	
	err = decoder.Decode(&result)

	return
}

func getHubDetails(path, hubKey, token string) (result ForgeResponseObject, err error) {
	task := http.Client{}

	req, err := http.NewRequest("GET",
		path+"/"+hubKey,
		nil,
	)

	if err != nil {
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)

	response, err := task.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()
	
	decoder := json.NewDecoder(response.Body)
		if response.StatusCode != http.StatusOK {
			err = &ErrorResult{StatusCode:response.StatusCode}
			decoder.Decode(err)
				return
		}
	
	err = decoder.Decode(&result)

	return
}