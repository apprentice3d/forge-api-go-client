package dm

import (
	"encoding/json"
	"net/http"
	"github.com/outer-labs/forge-api-go-client/oauth"
)

// FolderAPI holds the necessary data for making calls to Forge Data Management service
type FolderAPI struct {
	oauth.TwoLeggedAuth
	FolderAPIPath string
}


// NewFolderAPIWithCredentials returns a Folder API client with default configurations
func NewFolderAPIWithCredentials(ClientID string, ClientSecret string) FolderAPI {
	return FolderAPI{
		oauth.NewTwoLeggedClient(ClientID, ClientSecret),
		"/data/v1/projects",
	}
}

// ListBuckets returns a list of all buckets created or associated with Forge secrets used for token creation
func (api FolderAPI) GetFolderDetails(projectKey, folderKey string) (result ItemDetails, err error) {
	
	// TO DO: take in optional header arguments
	// https://forge.autodesk.com/en/docs/data/v2/reference/http/projects-project_id-folders-folder_id-GET/
	bearer, err := api.Authenticate("data:read")
	if err != nil {
		return
	}

	path := api.Host + api.FolderAPIPath

	return getFolderDetails(path, projectKey, folderKey, bearer.AccessToken)
}

func (api FolderAPI) GetFolderContents(projectKey, folderKey string) (result FolderContents, err error) {
	bearer, err := api.Authenticate("data:read")
	if err != nil {
		return
	}
	path := api.Host + api.FolderAPIPath

	return getFolderContents(path, projectKey, folderKey, bearer.AccessToken)
}


/*
 *	SUPPORT FUNCTIONS
 */
func getFolderDetails(path, projectKey, folderKey, token string) (result ItemDetails, err error) {
	task := http.Client{}

	req, err := http.NewRequest("GET",
		path+"/"+projectKey+"/folders/"+folderKey,
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

func getFolderContents(path, projectKey, folderKey, token string) (result FolderContents, err error) {
	task := http.Client{}

	req, err := http.NewRequest("GET",
		path+"/"+projectKey+"/folders/"+folderKey+"/contents",
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