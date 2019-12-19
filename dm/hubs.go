package dm

import (
	// "fmt"
	"encoding/json"
	"net/http"
	// "github.com/outer-labs/forge-api-go-client/oauth"
	"../oauth"
)

// HubAPI holds the necessary data for making calls to Forge Data Management service
type HubAPI struct {
	oauth.TwoLeggedAuth
	HubAPIPath string
}

type HubDetails struct {
	Details DataDetails `json:"details, omitempty"`
}

var api HubAPI

// NewHubAPIWithCredentials returns a Hub API client with default configurations
func NewHubAPIWithCredentials(ClientID string, ClientSecret string) HubAPI {
	return HubAPI{
		oauth.NewTwoLeggedClient(ClientID, ClientSecret),
		"/project/v1/hubs",
	}
}

func (api HubAPI) GetHubDetails(hubKey string) (result HubDetails, err error) {
	bearer, err := api.Authenticate("data:read")
	if err != nil {
		return
	}
	path := api.Host + api.HubAPIPath

	return getHubDetails(path, hubKey, bearer.AccessToken)
}

/*
 *	SUPPORT FUNCTIONS
 */
func getHubDetails(path, hubKey, token string) (result HubDetails, err error) {
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