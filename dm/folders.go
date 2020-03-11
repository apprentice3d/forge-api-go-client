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

type FolderAPI3L struct {
	oauth.ThreeLeggedAuth
	FolderAPIPath string
}

// NewFolderAPIWithCredentials returns a Folder API client with default configurations
func NewFolderAPIWithCredentials(ClientID string, ClientSecret string) FolderAPI {
	return FolderAPI{
		oauth.NewTwoLeggedClient(ClientID, ClientSecret),
		"/data/v1/projects",
	}
}

func NewFolderAPI3LWithCredentials(threeLeggedAuth oauth.ThreeLeggedAuth) FolderAPI3L {
	return FolderAPI3L{
		threeLeggedAuth,
		"/data/v1/projects",
	}
}

// ListBuckets returns a list of all buckets created or associated with Forge secrets used for token creation
func (api FolderAPI) GetFolderDetails(projectKey, folderKey string) (result ForgeResponseObject, err error) {
	
	// TO DO: take in optional header arguments
	// https://forge.autodesk.com/en/docs/data/v2/reference/http/projects-project_id-folders-folder_id-GET/
	bearer, err := api.Authenticate("data:read")
	if err != nil {
		return
	}

	path := api.Host + api.FolderAPIPath

	return getFolderDetails(path, projectKey, folderKey, bearer.AccessToken)
}

func (api FolderAPI) GetFolderContents(projectKey, folderKey string) (result ForgeResponseArray, err error) {
	bearer, err := api.Authenticate("data:read")
	if err != nil {
		return
	}
	path := api.Host + api.FolderAPIPath

	return getFolderContents(path, projectKey, folderKey, bearer.AccessToken)
}

// Three legged Folder api calls
func (api FolderAPI3L) GetFolderDetailsThreeLegged(bearer oauth.Bearer, projectKey, folderKey string) (result ForgeResponseObject, err error) {
	
	// TO DO: take in optional header arguments
	// https://forge.autodesk.com/en/docs/data/v2/reference/http/projects-project_id-folders-folder_id-GET/
	refreshedBearer, err := api.RefreshToken(bearer.RefreshToken, "data:read")
	if err != nil {
		return
	}

	path := api.Host + api.FolderAPIPath

	return getFolderDetails(path, projectKey, folderKey, refreshedBearer.AccessToken)
}

func (api FolderAPI3L) GetFolderContentsThreeLegged(bearer oauth.Bearer, projectKey, folderKey string) (result ForgeResponseArray, err error) {
	refreshedBearer, err := api.RefreshToken(bearer.RefreshToken, "data:read")
	if err != nil {
		return
	}
	path := api.Host + api.FolderAPIPath

	return getFolderContents(path, projectKey, folderKey, refreshedBearer.AccessToken)
}


/*
 *	SUPPORT FUNCTIONS
 */
func getFolderDetails(path, projectKey, folderKey, token string) (result ForgeResponseObject, err error) {
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

func getFolderContents(path, projectKey, folderKey, token string) (result ForgeResponseArray, err error) {
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