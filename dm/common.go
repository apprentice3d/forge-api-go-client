package dm

type Attribute struct {
	Name      string `json:"name, omitempty"`
	Region	  string `json: region, omitmepty"`
	Extension struct {
		Data    map[string]interface{} `json:"data, omitempty"`
		Version string                 `json:"version, omitempty"`
		Type    string                 `json:"type, omitempty"`
		Schema  struct {
			Href string `json:"href, omitempty"`
		} `json:"schema, omitempty"`
	} `json:"extension, omitempty"`
}

type Content struct {
	Relationships struct {
		Projects struct {
			Links Link `json:"links, omitempty"`
		} `json:"projects, omitempty"`
	} `json:"relationships, omitempty"`
	Attributes Attribute `json:"attributes, omitempty"`
	Type       string    `json:"type, omitempty"`
	Id         string    `json:"id, omitempty"`
	Links      Link      `json:"links, omitempty"`
}

type JsonAPI struct {
	Version string `json:"version, omitempty"`
}

type Link struct {
	Self struct {
		Href 	string `json:"href, omitempty"`
	} `json:"self, omitempty"`
	First struct {
		Href string `json:"href, omitempty"`
	} `json:"first, omitempty"`
	Prev struct {
		Href string `json:"href, omitempty"`
	} `json:"prev, omitempty"`
	Next struct {
		Href string `json:"href, omitempty"`
	} `json:"next, omitempty"`
	Related struct {
		Href string `json:"href, omitempty"`
	} `json:"related, omitempty"`
}
