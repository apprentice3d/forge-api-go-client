package dm

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/outer-labs/forge-api-go-client/oauth"
)

type Hubs struct {
	Data    []Content `json:"data, omitempty"`
	JsonApi JsonAPI   `json:"jsonapi, omitempty"`
	Links   Link      `json:"links, omitempty"`
}

type JsonAPI struct {
	Version string `json:"version, omitempty"`
}

type Link struct {
	Self struct {
		Href string `json:"href, omitempty"`
	} `json:"self, omitempty"`
}

type Content struct {
	Relationships struct {
		Projects Project `json:"projects, omitempty"`
	} `json:"relationships, omitempty"`
	Attributes Attribute `json:"attributes, omitempty"`
	Type       string    `json:"type, omitempty"`
	Id         string    `json:"id, omitempty"`
	Links      Link      `json:"links, omitempty"`
}

type Project struct {
	Links struct {
		Related struct {
			Href string `json:"href, omitempty"`
		} `json:"related, omitempty"`
	} `json:"links, omitempty"`
}

type Attribute struct {
	Name      string `json:"name, omitempty"`
	Extension struct {
		Data    map[string]interface{} `json:"data, omitempty"`
		Version string                 `json:"version, omitempty"`
		Type    string                 `json:"type, omitempty"`
		Schema  struct {
			Href string `json:"href, omitempty"`
		} `json:"schema, omitempty"`
	} `json:"extension, omitempty"`
}


// HubAPI holds the necessary data for making Bucket related calls to Forge Data Management service
type HubAPI struct {
	oauth.TwoLeggedAuth
	HubAPIPath string
}

// NewHubAPIWithCredentials returns a Hub API client with default configurations
func NewHubAPIWithCredentials(ClientID string, ClientSecret string) HubAPI {
	return HubAPI{
		oauth.NewTwoLeggedClient(ClientID, ClientSecret),
		"/project/v1/hubs",
	}
}

func (api HubAPI) GetHubDetails(path, hubKey, token string) (result HubDetails, err error) {
	bearer, err := api.Authenticate("hub:read")
	if err != nil {
		return
	}
	path := api.Host + api.BucketAPIPath

	return getHubDetails(path, hubKey, bearer.AccessToken)
}

/*
 *	SUPPORT FUNCTIONS
 */
func getHubDetails(path, hubKey, token string) (result ListedBuckets, err error) {
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