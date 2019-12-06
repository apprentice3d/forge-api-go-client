package dm

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

type ListedHubs struct {
	JsonAPI []struct {
		Version   string `json:"version"`
	} `json:"jsonAPI"`
	Links []struct {
		Self []struct {	
			Href string `json:"href"`
		} `json:"self"`
	} `json:"links"`
	Data []struct {
		Type   string `json:"type"`
		Id string `json:"id"`
		Attributes []struct {	
			Name string `json:"name"`
			Extension []struct{
				Type string `json:"type"`
				Version string `json:"version"`
				Schema []struct {
					Href string `json:"href"`
				} `json:"schema"`
				Data string `json"data"`
			} `json:"extension"`
			Region string `json:"region"`
		} `json:"attributes"`
		Relationships []struct {
			Projects []struct {
				Links []struct {
					Related []struct {
						Href string `json:"href"`
					} `json:"related"`
				} `json:"links"`
			} `json:"projects"`
		} `json:"relationships"`
		Links []struct {
			Self []struct {
				Href string `json:"href"`
			}
		} `json:"links"`
	} `json:"data"`
}

// NewHubAPIWithCredentials returns a Hub API client with default configurations
func NewHubAPIWithCredentials(ClientID string, ClientSecret string) HubAPI {
	return HubAPI{
		oauth.NewTwoLeggedClient(ClientID, ClientSecret),
		"/project/v1/hubs",
	}
}

func (api HubAPI) ListHubs(id, name, extension string) (result ListedHubs, err error){

}

func 