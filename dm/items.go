package dm

import (
	"encoding/json"
	"net/http"
)

type ItemDetails struct {
	Details 	DataDetails 	`json:"details, omitempty"`
	Included 	[]Content 		`json:"included, omitempty"`
}

// ListBuckets returns a list of all buckets created or associated with Forge secrets used for token creation
func (api FolderAPI) GetItemDetails(projectKey, itemKey string) (result ItemDetails, err error) {
	
	// TO DO: take in optional header argument
	// https://forge.autodesk.com/en/docs/data/v2/reference/http/projects-project_id-items-item_id-GET/
	bearer, err := api.Authenticate("data:read")
	if err != nil {
		return
	}

	path := api.Host + api.FolderAPIPath

	return getItemDetails(path, projectKey, itemKey, bearer.AccessToken)
}

/*
 *	SUPPORT FUNCTIONS
 */
func getItemDetails(path, projectKey, itemKey, token string) (result ItemDetails, err error) {
	task := http.Client{}

	req, err := http.NewRequest("GET",
		path+"/"+projectKey+"/items/"+itemKey,
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
