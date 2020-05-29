package dm

import (
	"encoding/json"
	"errors"
	"fmt"
	"forge-api-go-client/oauth"
	"io/ioutil"
	"net/http"
	"strconv"
)

type HubsAPI struct {
	oauth.TwoLeggedAuth
	HubsAPIPath string
}

func NewHubsAPIWithCredentials(ClientID, ClientSecret string) HubsAPI {
	return HubsAPI{
		TwoLeggedAuth: oauth.NewTwoLeggedClient(ClientID, ClientSecret),
		HubsAPIPath:   "/project/v1/hubs",
	}
}

type Hubs struct {
	Data    []Content `json:"data,omitempty"`
	JsonApi JsonAPI   `json:"jsonapi,omitempty"`
	Links   Link      `json:"links,omitempty"`
}

type JsonAPI struct {
	Version string `json:"version,omitempty"`
}

type Link struct {
	Self struct {
		Href string `json:"href,omitempty"`
	} `json:"self,omitempty"`
}

type Content struct {
	Relationships struct {
		Projects Project `json:"projects,omitempty"`
	} `json:"relationships,omitempty"`
	Attributes Attribute `json:"attributes,omitempty"`
	Type       string    `json:"type,omitempty"`
	Id         string    `json:"id,omitempty"`
	Links      Link      `json:"links,omitempty"`
}

type Project struct {
	Links struct {
		Related struct {
			Href string `json:"href,omitempty"`
		} `json:"related,omitempty"`
	} `json:"links,omitempty"`
}

type Attribute struct {
	Name      string `json:"name,omitempty"`
	Extension struct {
		Data    map[string]interface{} `json:"data,omitempty"`
		Version string                 `json:"version,omitempty"`
		Type    string                 `json:"type,omitempty"`
		Schema  struct {
			Href string `json:"href,omitempty"`
		} `json:"schema,omitempty"`
	} `json:"extension,omitempty"`
}

type Projects struct {
	Jsonapi struct {
		Version string `json:"version,omitempty"`
	} `json:"jsonapi,omitempty"`
	Links struct {
		Self struct {
			Href string `json:"href,omitempty"`
		} `json:"self,omitempty"`
	} `json:"links,omitempty"`
	Data []struct {
		Type       string `json:"type,omitempty"`
		ID         string `json:"id,omitempty"`
		Attributes struct {
			Name      string `json:"name,omitempty"`
			Extension struct {
				Type    string `json:"type,omitempty"`
				Version string `json:"version,omitempty"`
				Schema  struct {
					Href string `json:"href,omitempty"`
				} `json:"schema,omitempty"`
				Data struct {
				} `json:"data,omitempty"`
			} `json:"extension,omitempty"`
		} `json:"attributes,omitempty"`
		Links struct {
			Self struct {
				Href string `json:"href,omitempty"`
			} `json:"self,omitempty"`
		} `json:"links,omitempty"`
		Relationships struct {
			Hub struct {
				Data struct {
					Type string `json:"type,omitempty"`
					ID   string `json:"id,omitempty"`
				} `json:"data,omitempty"`
				Links struct {
					Related struct {
						Href string `json:"href,omitempty"`
					} `json:"related,omitempty"`
				} `json:"links,omitempty"`
			} `json:"hub,omitempty"`
			RootFolder struct {
				Data struct {
					Type string `json:"type,omitempty"`
					ID   string `json:"id,omitempty"`
				} `json:"data,omitempty"`
				Meta struct {
					Link struct {
						Href string `json:"href,omitempty"`
					} `json:"link,omitempty"`
				} `json:"meta,omitempty"`
			} `json:"rootFolder,omitempty"`
		} `json:"relationships,omitempty"`
	} `json:"data,omitempty"`
}


func (api *HubsAPI) ListHubs() (result Hubs, err error){
	bearer, err := api.Authenticate("data:read")
	if err != nil {
		return Hubs{}, err
	}

	path := api.Host + api.HubsAPIPath
	result, err = listHubs(path, bearer.AccessToken)
	return result, err
}

func listHubs(path string, token string) (result Hubs, err error) {
	task := http.Client{}

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return Hubs{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	response, err := task.Do(req)
	if err != nil {
		return Hubs{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		content, _ := ioutil.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return Hubs{}, err
	}

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&result)
	return result, nil
}

func (api *HubsAPI) GetHubProjects(hubId string) (result Projects, err error){
	bearer, err := api.Authenticate("data:read")
	if err != nil {
		return Projects{}, err
	}

	path := fmt.Sprintf("%s/%s/%s/projects", api.Host, api.HubsAPIPath, hubId)
	result, err = getHubProjects(path, bearer.AccessToken)
	return result, err
}

func getHubProjects(path string, token string) (result Projects, err error) {
	task := http.Client{}

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return Projects{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	response, err := task.Do(req)
	if err != nil {
		return Projects{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		content, _ := ioutil.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return Projects{}, err
	}

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&result)
	return result, nil
}
