package dm

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/outer-labs/forge-api-go-client/oauth"
)

type ProjectDetails struct {
	Data    []Content `json:"data, omitempty"`
	JsonApi JsonAPI   `json:"jsonapi, omitempty"`
	Links   Link      `json:"links, omitempty"`
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

type Content struct {
	Relationships Relationships `json:"relationships, omitempty"`
	Attributes Attribute `json:"attributes, omitempty"`
	Type       string    `json:"type, omitempty"`
	Id         string    `json:"id, omitempty"`
	Links      Link      `json:"links, omitempty"`
}

type Relationships struct {
	Hub 		[]Hub 		`json:"hub, omitempty"`
	RootFolder 	RootFolder 	`json:"rootfolder, omitempty"`
	TopFolders TopFolders 	`json:"topfolders, omitempty"`
}

type Hub struct {
	Links 		Link 		`json:"links, omitempty"`
	Data 		[]Content 	`json:"data, omitempty"`
}

type RootFolder struct {
	Meta struct {
		Links Link `json:"links, omitempty"`
	} `json:"meta, omitempty"`
	Data    []Content `json:"data, omitempty"`
}

type TopFolders struct {
	Links Link `json:"links, omitempty"`
}

type Project struct {
	Links Link `json:"links, omitempty"`
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