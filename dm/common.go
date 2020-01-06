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
	Relationships 	Relationships 	`json:"relationships, omitempty"`
	Attributes 	  	Attribute 		`json:"attributes, omitempty"`
	Type       		string    		`json:"type, omitempty"`
	Id         		string    		`json:"id, omitempty"`
	Links      		Link      		`json:"links, omitempty"`
}

type DataDetails struct {
	Data    	[]Content 	`json:"data, omitempty"`
	JsonApi 	JsonAPI   	`json:"jsonapi, omitempty"`
	Links   	Link      	`json:"links, omitempty"`
}

type FolderContents struct {
	JsonApi 	JsonAPI   	`json:"jsonapi, omitempty"`
	Links   	Link      	`json:"links, omitempty"`
	Data    	[]Content 	`json:"data, omitempty"`
	Included 	[]Content 	`json:"included, omitempty"`
}

type ItemDetails struct {
	Data    	Content 	`json:"data, omitempty"`
	JsonApi 	JsonAPI   	`json:"jsonapi, omitempty"`
	Links   	Link      	`json:"links, omitempty"`
	Included 	[]Content 	`json:"included, omitempty"`
}

type Hub struct {
	Links 		Link 		`json:"links, omitempty"`
	Data 		[]Content 	`json:"data, omitempty"`
}

type JsonAPI struct {
	Version 	string 		`json:"version, omitempty"`
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

type Project struct {
	Links Link `json:"links, omitempty"`
}

type Relationships struct {
	Projects 	Project 	`json:"projects, omitempty"`
	Hub 		struct {
		Links 		Link 		`json:"links, omitempty"`
		Data 		struct {
			ID 		string 	`json:"id, omitempty"`
			Type 	string 	`json:"type, omitempty"`
		} 	`json:"data, omitempty"`
	} 		`json:"hub, omitempty"`
	RootFolder 	RootFolder 	`json:"rootfolder, omitempty"`
	TopFolders 	TopFolders 	`json:"topfolders, omitempty"`
	Storage Storage `json:"storage, omitempty"`
}

type Storage struct {
	Meta struct {
		Links Link `json:"links, omitempty"`
	} `json:"meta, omitempty"`
	Data struct {
		ID 		string 	`json:"id, omitempty"`
		Type 	string 	`json:"type, omitempty"`
	} 	`json:"data, omitempty"`
}

type RootFolder struct {
	Meta struct {
		Links Link `json:"links, omitempty"`
	} `json:"meta, omitempty"`
	Data struct {
		ID 		string 	`json:"id, omitempty"`
		Type 	string 	`json:"type, omitempty"`
	} 	`json:"data, omitempty"`
}

type TopFolders struct {
	Links 		Link 	`	json:"links, omitempty"`
}
