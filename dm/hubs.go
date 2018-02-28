package dm
//
//type Hubs struct {
//	Data    []Content `json:"data,omitempty"`
//	JsonApi JsonAPI   `json:"jsonapi,omitempty"`
//	Links   Link      `json:"links,omitempty"`
//}
//
//type JsonAPI struct {
//	Version string `json:"version,omitempty"`
//}
//
//type Link struct {
//	Self struct {
//		Href string `json:"href,omitempty"`
//	} `json:"self,omitempty"`
//}
//
//type Content struct {
//	Relationships struct {
//		Projects Project `json:"projects,omitempty"`
//	} `json:"relationships,omitempty"`
//	Attributes Attribute `json:"attributes,omitempty"`
//	Type       string    `json:"type,omitempty"`
//	Id         string    `json:"id,omitempty"`
//	Links      Link      `json:"links,omitempty"`
//}
//
//type Project struct {
//	Links struct {
//		Related struct {
//			Href string `json:"href,omitempty"`
//		} `json:"related,omitempty"`
//	} `json:"links,omitempty"`
//}
//
//type Attribute struct {
//	Name      string `json:"name,omitempty"`
//	Extension struct {
//		Data    map[string]interface{} `json:"data,omitempty"`
//		Version string                 `json:"version,omitempty"`
//		Type    string                 `json:"type,omitempty"`
//		Schema  struct {
//			Href string `json:"href,omitempty"`
//		} `json:"schema,omitempty"`
//	} `json:"extension,omitempty"`
//}
