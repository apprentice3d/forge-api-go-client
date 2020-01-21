package dm

type ForgeResponseObject struct {
	JsonApi 				JsonAPI 		`json:"jsonApi`
	Links 					Links 			`json:"links"`
	Data 					Data 			`json:"data"`
	Included 				*[]Data 		`json:"included, omitempty"`
}

type ForgeResponseArray struct {
	JsonApi 				JsonAPI 		`json:"jsonApi`
	Links 					Links 			`json:"links"`
	Data 					[]Data 			`json:"data"`
	Included 				*[]Data 		`json:"included, omitempty"`
}

type JsonAPI struct {
	Version 				string 			`json:"version"`
}

type Links struct {
	Self 					*Href 			`json:"self, omitempty"`
	Related 				*Href 			`json:"related, omitempty"`
	First 					*Href 			`json:"first, omitempty"`
	Prev 					*Href 			`json:"prev, omitempty"`
	Next 					*Href 			`json:"next, omitempty"`
}

type Data struct {
	Type 					string 			`json:"type"` 	 
	Id 					string 			`json:"id"`
	Attributes 				*Attributes 		`json:"attributes, omitempty"`
	Relationships 				*Relationships 		`json:"relationships, omitempty"`
	Links 					*Links 			`json:"links, omitempty"`
}

type Attributes struct {
	Name      				string 			`json:"name"`
	Extension 				Extension 		`json:"extension"`
	Region	  				*string 		`json:"region, omitempty"`
	Scopes 					*[]string  		`json:"scopes, omitempty"`
	DisplayName      			*string 		`json:"displayName, omitempty"`
	ObjectCount      			*int 			`json:"objectCount, omitempty"`
	CreateTime      			*string 		`json:"createTime, omitempty"`
	CreateUserId      			*string 		`json:"createUserId, omitempty"`
	CreateUserName      			*string 		`json:"createUserName, omitempty"`
	LastModifiedTime    			*string 		`json:"lastModifiedTime, omitempty"`
	LastModifiedUserId  			*string 		`json:"lastModifiedUserId, omitempty"`
	LastModifiedUserName    		*string 		`json:"lastModifiedUserName, omitempty"`
	Hidden      				*bool 			`json:"displayName, omitempty"`
	VersionNumber      			*int 			`json:"versionNumber, omitempty"`
	Mimetype      				*string 		`json:"mimeType, omitempty"`
	FileType      				*string 		`json:"fileType, omitempty"`
	StorageSize      			*int 			`json:"storageSize, omitempty"`
	Reserved 				*bool 			`json:"reserved, omitempty"`
	ReservedTime 				*string 		`json:"reservedTime, omitempty"`
	ReservedUserId 				*string 		`json:"reservedUserId, omitempty"`
	ReservedUserName			*string 		`json:"reservedUserName, omitempty"`
	PathInProject 				*string 		`json:"pathInProject, omitempty"`
}

type Relationships struct {
	Hub 					*RelatedLinks		`json:"hub, omitempty"`
	Projects 				*RelatedLinks 		`json:projects, omitempty"`
	RootFolder 				*RelatedLinks 		`json:"rootFolder, omitempty"`
	TopFolders 				*RelatedLinks 		`json:"topFolders, omitempty"`
	Parent  				*RelatedLinks 		`json:"parent, omitempty"`
	Tip  					*RelatedLinks 		`json:"tip, omitempty"`
	Versions  				*RelatedLinks 		`json:"versions, omitempty"`
	Contents  				*RelatedLinks 		`json:"contents, omitempty"`
	Refs  					*RelatedLinks 		`json:"refs, omitempty"`		
	Links 					*RelatedLinks 		`json:"links, omitempty"`
	Item 					*RelatedLinks		`json:"item, omitempty"`
	Storage 				*RelatedLinks		`json:"storage, omitempty"`
	Derivatives 				*RelatedLinks		`json:"derivatives, omitempty"`
	Thumbnails 				*RelatedLinks		`json:"thumbmails, omitempty"`
	DownloadFormats 			*RelatedLinks		`json:"downloadFormats, omitempty"`
}

type Extension struct {
	Type 					string 			`json:"type"`
	Version 				string 			`json:"version"`
	Schema 					Href			`json:"schema"`
	Data 					struct{} 		`json:"data"`
}

type RelatedLinks struct {
	Meta 					*Meta 			`json:"meta, omitempty"`
	Links 					*Links 			`json:"links, omitempty"`
	Data 					*Data			`json:"data, omitempty"`
}

type Meta struct {
	Link   					Href 			`json:"href"`	
}

type Href struct {
	Href 					string 			`json:"href"`
}

// Note on use of omitempty: https://www.sohamkamani.com/blog/golang/2018-07-19-golang-omitempty/
