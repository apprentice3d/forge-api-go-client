package dm

import (
	"encoding/json"
	"errors"
	"fmt"
	"forge-api-go-client/oauth"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type FolderContents struct {
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
			DisplayName          string    `json:"displayName,omitempty"`
			CreateTime           time.Time `json:"createTime,omitempty"`
			CreateUserID         string    `json:"createUserId,omitempty"`
			CreateUserName       string    `json:"createUserName,omitempty"`
			LastModifiedTime     time.Time `json:"lastModifiedTime,omitempty"`
			LastModifiedUserID   string    `json:"lastModifiedUserId,omitempty"`
			LastModifiedUserName string    `json:"lastModifiedUserName,omitempty"`
			Extension            struct {
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
			Tip struct {
				Data struct {
					Type string `json:"type,omitempty"`
					ID   string `json:"id,omitempty"`
				} `json:"data,omitempty"`
				Links struct {
					Related struct {
						Href string `json:"href,omitempty"`
					} `json:"related,omitempty"`
				} `json:"links,omitempty"`
			} `json:"tip,omitempty"`
			Versions struct {
				Links struct {
					Related struct {
						Href string `json:"href,omitempty"`
					} `json:"related,omitempty"`
				} `json:"links,omitempty"`
			} `json:"versions,omitempty"`
			Parent struct {
				Data struct {
					Type string `json:"type,omitempty"`
					ID   string `json:"id,omitempty"`
				} `json:"data,omitempty"`
				Links struct {
					Related struct {
						Href string `json:"href,omitempty"`
					} `json:"related,omitempty"`
				} `json:"links,omitempty"`
			} `json:"parent,omitempty"`
			Refs struct {
				Links struct {
					Self struct {
						Href string `json:"href,omitempty"`
					} `json:"self,omitempty"`
					Related struct {
						Href string `json:"href,omitempty"`
					} `json:"related,omitempty"`
				} `json:"links,omitempty"`
			} `json:"refs,omitempty"`
		} `json:"relationships,omitempty"`
	} `json:"data,omitempty"`
	Included []struct {
		Type       string `json:"type,omitempty"`
		ID         string `json:"id,omitempty"`
		Attributes struct {
			Name                 string    `json:"name,omitempty"`
			DisplayName          string    `json:"displayName,omitempty"`
			CreateTime           time.Time `json:"createTime,omitempty"`
			CreateUserID         string    `json:"createUserId,omitempty"`
			CreateUserName       string    `json:"createUserName,omitempty"`
			LastModifiedTime     time.Time `json:"lastModifiedTime,omitempty"`
			LastModifiedUserID   string    `json:"lastModifiedUserId,omitempty"`
			LastModifiedUserName string    `json:"lastModifiedUserName,omitempty"`
			VersionNumber        int       `json:"versionNumber,omitempty"`
			MimeType             string    `json:"mimeType,omitempty"`
			FileType             string    `json:"fileType,omitempty"`
			StorageSize          uint64       `json:"storageSize,omitempty"`
			Extension            struct {
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
			Item struct {
				Data struct {
					Type string `json:"type,omitempty"`
					ID   string `json:"id,omitempty"`
				} `json:"data,omitempty"`
				Links struct {
					Related struct {
						Href string `json:"href,omitempty"`
					} `json:"related,omitempty"`
				} `json:"links,omitempty"`
			} `json:"item,omitempty"`
			Refs struct {
				Links struct {
					Self struct {
						Href string `json:"href,omitempty"`
					} `json:"self,omitempty"`
					Related struct {
						Href string `json:"href,omitempty"`
					} `json:"related,omitempty"`
				} `json:"links,omitempty"`
			} `json:"refs,omitempty"`
			Derivatives struct {
				Data struct {
					Type string `json:"type,omitempty"`
					ID   string `json:"id,omitempty"`
				} `json:"data,omitempty"`
				Meta struct {
					Link struct {
						Href string `json:"href,omitempty"`
					} `json:"link,omitempty"`
				} `json:"meta,omitempty"`
			} `json:"derivatives,omitempty"`
			Thumbnails struct {
				Data struct {
					Type string `json:"type,omitempty"`
					ID   string `json:"id,omitempty"`
				} `json:"data,omitempty"`
				Meta struct {
					Link struct {
						Href string `json:"href,omitempty"`
					} `json:"link,omitempty"`
				} `json:"meta,omitempty"`
			} `json:"thumbnails,omitempty"`
			Storage struct {
				Data struct {
					Type string `json:"type,omitempty"`
					ID   string `json:"id,omitempty"`
				} `json:"data,omitempty"`
				Meta struct {
					Link struct {
						Href string `json:"href,omitempty"`
					} `json:"link,omitempty"`
				} `json:"meta,omitempty"`
			} `json:"storage,omitempty"`
		} `json:"relationships,omitempty"`
	} `json:"included,omitempty"`
}

type ItemData struct {
	Jsonapi struct {
		Version string `json:"version,omitempty"`
	} `json:"jsonapi,omitempty"`
	Links struct {
		Self struct {
			Href string `json:"href,omitempty"`
		} `json:"self,omitempty"`
	} `json:"links,omitempty"`
	Data struct {
		Type       string `json:"type,omitempty"`
		ID         string `json:"id,omitempty"`
		Attributes struct {
			DisplayName          string    `json:"displayName,omitempty"`
			CreateTime           time.Time `json:"createTime,omitempty"`
			CreateUserID         string    `json:"createUserId,omitempty"`
			CreateUserName       string    `json:"createUserName,omitempty"`
			LastModifiedTime     time.Time `json:"lastModifiedTime,omitempty"`
			LastModifiedUserID   string    `json:"lastModifiedUserId,omitempty"`
			LastModifiedUserName string    `json:"lastModifiedUserName,omitempty"`
			Extension            struct {
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
			Tip struct {
				Data struct {
					Type string `json:"type,omitempty"`
					ID   string `json:"id,omitempty"`
				} `json:"data,omitempty"`
				Links struct {
					Related struct {
						Href string `json:"href,omitempty"`
					} `json:"related,omitempty"`
				} `json:"links,omitempty"`
			} `json:"tip,omitempty"`
			Versions struct {
				Links struct {
					Related struct {
						Href string `json:"href,omitempty"`
					} `json:"related,omitempty"`
				} `json:"links,omitempty"`
			} `json:"versions,omitempty"`
			Parent struct {
				Data struct {
					Type string `json:"type,omitempty"`
					ID   string `json:"id,omitempty"`
				} `json:"data,omitempty"`
				Links struct {
					Related struct {
						Href string `json:"href,omitempty"`
					} `json:"related,omitempty"`
				} `json:"links,omitempty"`
			} `json:"parent,omitempty"`
			Refs struct {
				Links struct {
					Self struct {
						Href string `json:"href,omitempty"`
					} `json:"self,omitempty"`
					Related struct {
						Href string `json:"href,omitempty"`
					} `json:"related,omitempty"`
				} `json:"links,omitempty"`
			} `json:"refs,omitempty"`
		} `json:"relationships,omitempty"`
	} `json:"data,omitempty"`
	Included []struct {
		Type       string `json:"type,omitempty"`
		ID         string `json:"id,omitempty"`
		Attributes struct {
			Name                 string    `json:"name,omitempty"`
			DisplayName          string    `json:"displayName,omitempty"`
			CreateTime           time.Time `json:"createTime,omitempty"`
			CreateUserID         string    `json:"createUserId,omitempty"`
			CreateUserName       string    `json:"createUserName,omitempty"`
			LastModifiedTime     time.Time `json:"lastModifiedTime,omitempty"`
			LastModifiedUserID   string    `json:"lastModifiedUserId,omitempty"`
			LastModifiedUserName string    `json:"lastModifiedUserName,omitempty"`
			VersionNumber        int       `json:"versionNumber,omitempty"`
			MimeType             string    `json:"mimeType,omitempty"`
			FileType             string    `json:"fileType,omitempty"`
			StorageSize          uint64       `json:"storageSize,omitempty"`
			Extension            struct {
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
			Item struct {
				Data struct {
					Type string `json:"type,omitempty"`
					ID   string `json:"id,omitempty"`
				} `json:"data,omitempty"`
				Links struct {
					Related struct {
						Href string `json:"href,omitempty"`
					} `json:"related,omitempty"`
				} `json:"links,omitempty"`
			} `json:"item,omitempty"`
			Refs struct {
				Links struct {
					Self struct {
						Href string `json:"href,omitempty"`
					} `json:"self,omitempty"`
					Related struct {
						Href string `json:"href,omitempty"`
					} `json:"related,omitempty"`
				} `json:"links,omitempty"`
			} `json:"refs,omitempty"`
			Derivatives struct {
				Data struct {
					Type string `json:"type,omitempty"`
					ID   string `json:"id,omitempty"`
				} `json:"data,omitempty"`
				Meta struct {
					Link struct {
						Href string `json:"href,omitempty"`
					} `json:"link,omitempty"`
				} `json:"meta,omitempty"`
			} `json:"derivatives,omitempty"`
			Thumbnails struct {
				Data struct {
					Type string `json:"type,omitempty"`
					ID   string `json:"id,omitempty"`
				} `json:"data,omitempty"`
				Meta struct {
					Link struct {
						Href string `json:"href,omitempty"`
					} `json:"link,omitempty"`
				} `json:"meta,omitempty"`
			} `json:"thumbnails,omitempty"`
			Storage struct {
				Data struct {
					Type string `json:"type,omitempty"`
					ID   string `json:"id,omitempty"`
				} `json:"data,omitempty"`
				Meta struct {
					Link struct {
						Href string `json:"href,omitempty"`
					} `json:"link,omitempty"`
				} `json:"meta,omitempty"`
			} `json:"storage,omitempty"`
		} `json:"relationships,omitempty"`
	} `json:"included,omitempty"`
}

type ProjectsAPI struct {
	oauth.TwoLeggedAuth
	ProjectsAPIPath string
}

func NewProjectsAPIWithCredentials(ClientID, ClientSecret, ProjectId string) ProjectsAPI {
	return ProjectsAPI{
		TwoLeggedAuth: oauth.NewTwoLeggedClient(ClientID, ClientSecret),
		ProjectsAPIPath:   fmt.Sprintf("/data/v1/projects/%s/", ProjectId), // projects/:project_id/folders/:folder_id/
	}
}

func (api *ProjectsAPI) GetFolderContents(FolderId string) (result *FolderContents, err error){
	//bearer, err := api.Authenticate("data:read")
	bearer, err := api.AuthenticateIfNecessary("data:read")
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("%s%sfolders/%s", api.Host, api.ProjectsAPIPath, FolderId)
	result, err = getFolderContents(path, bearer.AccessToken)
	return result, err
}

func getFolderContents(path string, token string) (result *FolderContents, err error) {
	task := http.Client{}

	req, err := http.NewRequest("GET", path+"/contents", nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	response, err := task.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = errors.New(strconv.Itoa(response.StatusCode))
		return nil, err
	}

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&result)
	return result, nil
}

func (api *ProjectsAPI) GetItemData(itemId string) (result *ItemData, err error){
	//bearer, err := api.Authenticate("data:read")
	bearer, err := api.AuthenticateIfNecessary("data:read")
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("%s%sitems/%s", api.Host, api.ProjectsAPIPath, itemId)
	result, err = getItemData(path, bearer.AccessToken)
	return result, err
}

func getItemData(path string, token string) (result *ItemData, err error) {
	task := http.Client{}

	req, err := http.NewRequest("GET", path, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	response, err := task.Do(req)
	if err != nil {
		log.Printf("Error at request: %v", err.Error())
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		err = errors.New(strconv.Itoa(response.StatusCode))
		log.Printf("Error at request: %v", err.Error())
		return nil, err
	}

	defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&result)
	return result, nil
}

func (api *ProjectsAPI) GetItemReader(itemStorageLink string) (result *io.ReadCloser, err error) {
	bearer, err := api.AuthenticateIfNecessary("data:read")
	if err != nil {
		return nil, err
	}

	result, err = getItemReader(itemStorageLink, bearer.AccessToken)
	return result, err
}

func getItemReader(link string, token string) (*io.ReadCloser, error) {
	task := http.Client{}

	req, err := http.NewRequest("GET", link, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	response, err := task.Do(req)
	if err != nil {
		log.Printf("Error at request: %v", err.Error())
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		err = errors.New(strconv.Itoa(response.StatusCode))
		log.Printf("Error at request: %v", err.Error())
		return nil, err
	}

	return &response.Body, nil
}