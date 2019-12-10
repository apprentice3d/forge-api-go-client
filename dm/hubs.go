package dm

import (
	// "bytes"
	"encoding/json"
	"net/http"
	"github.com/outer-labs/forge-api-go-client/oauth"
)

// HubAPI holds the necessary data for making calls to Forge Data Management service
type HubAPI struct {
	oauth.TwoLeggedAuth
	HubAPIPath string
}

type HubDetails struct {
	Data    []Content `json:"data, omitempty"`
	JsonApi JsonAPI   `json:"jsonapi, omitempty"`
	Links   Link      `json:"links, omitempty"`
}

// CreateHubRequest contains the data necessary to be passed upon hub creation
// type CreateHubRequest struct {
// 	HubKey string `json:"bucketKey"`
// 	PolicyKey string `json:"policyKey"`
// }


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


// func (api HubAPI) CreateHub(hubKey, policyKey string) (result HubDetails, err error) {
// 	bearer, err := api.Authenticate("hub:create")
// 	if err != nil {
// 		return
// 	}
// 	path := api.Host + api.HubAPIPath
// 	result, err = createHub(path, hubKey, policyKey, bearer.AccessToken)

// 	return

// }

// func (api HubAPI) DeleteHub(hubKey string) error {
// 	bearer, err := api.Authenticate("hub:delete")
// 	if err != nil {
// 		return err
// 	}
// 	path := api.Host + api.HubAPIPath

// 	return deleteHub(path, hubKey, bearer.AccessToken)
// }


/*
 *	SUPPORT FUNCTIONS
 */
func getHubDetails(path, hubKey, token string) (result HubDetails, err error) {
	task := http.Client{}

	req, err := http.NewRequest("GET",
		path+"/"+hubKey+"/details",
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

// func createHub(path, hubKey, policyKey, token string) (result HubDetails, err error) {

// 	task := http.Client{}

// 	body, err := json.Marshal(
// 		CreateHubRequest{
// 			hubKey,
// 			policyKey,
// 		})
// 	if err != nil {
// 		return
// 	}

// 	req, err := http.NewRequest("POST",
// 		path,
// 		bytes.NewReader(body),
// 	)

// 	if err != nil {
// 		return
// 	}
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Authorization", "Bearer "+token)
// 	response, err := task.Do(req)
// 	if err != nil {
// 		return
// 	}
// 	defer response.Body.Close()

// 	decoder := json.NewDecoder(response.Body)
// 	if response.StatusCode != http.StatusOK {
//     err = &ErrorResult{StatusCode:response.StatusCode}
//     decoder.Decode(err)
// 		return
// 	}

// 	err = decoder.Decode(&result)

// 	return

// }

// func deleteHub(path, hubKey, token string) (err error) {
// 	task := http.Client{}

// 	req, err := http.NewRequest("DELETE",
// 		path+"/"+hubKey,
// 		nil,
// 	)

// 	if err != nil {
// 		return
// 	}

// 	req.Header.Set("Authorization", "Bearer "+token)
// 	response, err := task.Do(req)
// 	if err != nil {
// 		return
// 	}
// 	defer response.Body.Close()
//   decoder := json.NewDecoder(response.Body)
// 	if response.StatusCode != http.StatusOK {
//     err = &ErrorResult{StatusCode:response.StatusCode}
//     decoder.Decode(err)
// 		return
// 	}

// 	return
// }