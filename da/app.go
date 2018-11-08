package da

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/apprentice3d/forge-api-go-client/oauth"
	"io/ioutil"
	"net/http"
	"strconv"
)

type AppList struct {
	InfoList
}
type AppParameters struct {
	URL  string `json:"endpointURL"`
	Data struct {
		Key         string `json:"key"`
		ContentType string `json:"content-type"`
		Policy      string `json:"policy"`
		Status      string `json:"success_action_status"`
		Redirect    string `json:"success_action_redirect"`
		Signature   string `json:"x-amz-signature"`
		Credential  string `json:"x-amz-credential"`
		Algorithm   string `json:"x-amz-algorithm"`
		Date        string `json:"x-amz-date"`
		Encryption  string `json:"x-amz-server-side-encryption"`
		Token       string `json:"x-amz-security-token"`
	} `json:"formData"`
}

type AppData struct {
	Engine  string `json:"engine"`
	Version uint   `json:"version"`
	ID      string `json:"id"`
}

type AppBundle struct {
	Parameters AppParameters `json:"uploadParameters"`
	AppData

	authenticator *oauth.TwoLeggedAuth
	path          string
	name          string
}

type AppDetails struct {
	Package string `json:"package"`
	AppData
}

type CreateAppRequest struct {
	ID     string `json:"id"`
	Engine string `json:"engine"`
}

type AppAlias struct {
	ID      string `json:"id"`
	Version uint   `json:"version"`
}

type AliasesList struct {
	Pagination string     `json:"paginationToken"`
	Data       []AppAlias `json:"data"`
}

type VersionList struct {
	Pagination string     `json:"paginationToken"`
	Data       []uint `json:"data"`
}


type VersionDetails struct {

}


// Delete removes the AppBundle, including all versions and aliases.
func (app *AppBundle) Delete() (err error) {

	bearer, err := app.authenticator.Authenticate("code:all")
	if err != nil {
		return
	}

	err = deleteApp(app.path, app.name, bearer.AccessToken)

	// TODO: research for a more elegant way of self-removing
	app.Parameters = AppParameters{}
	app.Engine = ""
	app.name = ""
	app.ID = ""
	app.Version = 0
	app.authenticator = nil
	app.path = ""

	return
}

//Details gets the details of the specified AppBundle.
func (app *AppBundle) Details() (details AppDetails, err error) {
	bearer, err := app.authenticator.Authenticate("code:all")
	if err != nil {
		return
	}
	details, err = getAppDetails(app.path, app.ID, bearer.AccessToken)

	return
}

//Aliases lists all aliases for the specified AppBundle.
func (app AppBundle) Aliases() (list AliasesList, err error) {
	bearer, err := app.authenticator.Authenticate("code:all")
	if err != nil {
		return
	}
	list, err = listAppAliases(app.path, app.name, bearer.AccessToken)

	return
}

//CreateAlias creates a new alias for this AppBundle.
//	Limit: 1. Number of aliases (LimitAliases).
func (app AppBundle) CreateAlias(alias string, version uint) (result AppAlias, err error) {
	bearer, err := app.authenticator.Authenticate("code:all")
	if err != nil {
		return
	}
	result, err = createAppAlias(app.path, app.name, alias, version, bearer.AccessToken)

	return
}

//ModifyAlias will switch the given alias to another existing version
func (app AppBundle) ModifyAlias(alias string, version uint) (result AppAlias, err error) {
	bearer, err := app.authenticator.Authenticate("code:all")
	if err != nil {
		return
	}
	result, err = modifyAppAlias(app.path, app.name, alias, version, bearer.AccessToken)

	return
}

//AliasDetail gets the details on given alias
func (app *AppBundle) AliasDetail(alias string) (details AppAlias, err error) {
	bearer, err := app.authenticator.Authenticate("code:all")
	if err != nil {
		return
	}
	details, err = getAliasDetails(app.path, app.name, alias, bearer.AccessToken)

	return
}

//DeleteAlias the alias for this AppBundle.
func (app AppBundle) DeleteAlias(alias string) (err error) {
	bearer, err := app.authenticator.Authenticate("code:all")
	if err != nil {
		return
	}
	err = deleteAppAlias(app.path, app.name, alias, bearer.AccessToken)

	return
}

//Versions lists all aliases for the specified AppBundle.
func (app AppBundle) Versions() (list VersionList, err error) {
	bearer, err := app.authenticator.Authenticate("code:all")
	if err != nil {
		return
	}
	list, err = listAppVersions(app.path, app.name, bearer.AccessToken)

	return
}

func (app AppBundle) CreateVersion(engine string) (result AppBundle, err error) {
	bearer, err := app.authenticator.Authenticate("code:all")
	if err != nil {
		return
	}
	result, err = createAppVersion(app.path, app.name, engine, bearer.AccessToken)
	result.authenticator = app.authenticator
	result.name = app.name
	result.path = app.path

	return
}


func (app *AppBundle) VersionDetails(version uint) (details AppData, err error) {
	bearer, err := app.authenticator.Authenticate("code:all")
	if err != nil {
		return
	}
	details, err = getVersionDetails(app.path, app.name, version, bearer.AccessToken)

	return
}


func (app AppBundle) DeleteVersion(version uint) (err error) {
	bearer, err := app.authenticator.Authenticate("code:all")
	if err != nil {
		return
	}
	err = deleteAppVersion(app.path, app.name, version, bearer.AccessToken)

	return
}











/*
 *	SUPPORT FUNCTIONS
 */



/*
   APPBUNDLE
*/

func listApps(path string, token string) (list AppList, err error) {

	task := http.Client{}
	req, err := http.NewRequest("GET",
		path+"/appbundles",
		nil,
	)

	req.Header.Set("Authorization", "Bearer "+token)
	response, err := task.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		content, _ := ioutil.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&list)

	return
}

func createApp(path, name, engine, token string) (result AppBundle, err error) {

	task := http.Client{}

	body, err := json.Marshal(
		CreateAppRequest{
			name,
			engine,
		})
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST",
		path+"/appbundles",
		bytes.NewReader(body),
	)

	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	response, err := task.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		content, _ := ioutil.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	decoder := json.NewDecoder(response.Body)

	err = decoder.Decode(&result)

	return
}

func getAppDetails(path, appID, token string) (result AppDetails, err error) {

	task := http.Client{}

	req, err := http.NewRequest("GET",
		path+"/appbundles/"+appID,
		nil,
	)

	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	response, err := task.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		content, _ := ioutil.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	decoder := json.NewDecoder(response.Body)

	err = decoder.Decode(&result)

	return
}

func deleteApp(path string, id string, token string) (err error) {

	task := http.Client{}
	req, err := http.NewRequest("DELETE",
		path+"/appbundles/"+id,
		nil,
	)

	req.Header.Set("Authorization", "Bearer "+token)
	response, err := task.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		content, _ := ioutil.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	return
}


/*
	ALIASES
*/

func listAppAliases(path string, appName, token string) (list AliasesList, err error) {

	task := http.Client{}
	req, err := http.NewRequest("GET",
		path+"/appbundles/"+appName+"/aliases",
		nil,
	)

	req.Header.Set("Authorization", "Bearer "+token)
	response, err := task.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		content, _ := ioutil.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&list)

	return
}

func createAppAlias(path, appName, alias string, version uint, token string) (result AppAlias, err error) {

	task := http.Client{}

	body, err := json.Marshal(
		AppAlias{
			alias,
			version,
		})
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST",
		path+"/appbundles/"+appName+"/aliases",
		bytes.NewReader(body),
	)

	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	response, err := task.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		content, _ := ioutil.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	decoder := json.NewDecoder(response.Body)

	err = decoder.Decode(&result)

	return
}

func modifyAppAlias(path, appName, alias string, version uint, token string) (result AppAlias, err error) {

	task := http.Client{}

	body, err := json.Marshal(
		struct {
			Version uint `json:"version"`
		}{version})
	if err != nil {
		return
	}

	req, err := http.NewRequest("PATCH",
		path+"/appbundles/"+appName+"/aliases/"+alias,
		bytes.NewReader(body),
	)

	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	response, err := task.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		content, _ := ioutil.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	decoder := json.NewDecoder(response.Body)

	err = decoder.Decode(&result)

	return
}

func getAliasDetails(path, appName, alias, token string) (result AppAlias, err error) {

	task := http.Client{}

	req, err := http.NewRequest("GET",
		path+"/appbundles/"+appName+"/aliases/"+alias,
		nil,
	)

	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	response, err := task.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		content, _ := ioutil.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	decoder := json.NewDecoder(response.Body)

	err = decoder.Decode(&result)

	return
}



func deleteAppAlias(path string, appName, alias, token string) (err error) {

	task := http.Client{}
	req, err := http.NewRequest("DELETE",
		path+"/appbundles/"+appName+"/aliases/"+alias,
		nil,
	)

	req.Header.Set("Authorization", "Bearer "+token)
	response, err := task.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		content, _ := ioutil.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	return
}


/*
   VERSIONS
*/

func listAppVersions(path string, appName, token string) (list VersionList, err error) {

	task := http.Client{}
	req, err := http.NewRequest("GET",
		path+"/appbundles/"+appName+"/versions",
		nil,
	)

	req.Header.Set("Authorization", "Bearer "+token)
	response, err := task.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		content, _ := ioutil.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&list)

	return
}

func createAppVersion(path, appName, engine string, token string) (result AppBundle, err error) {

	task := http.Client{}

	body, err := json.Marshal(
		struct{
			Engine string `json:"engine"`
		}{engine})
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST",
		path+"/appbundles/"+appName+"/versions",
		bytes.NewReader(body),
	)

	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	response, err := task.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		content, _ := ioutil.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	decoder := json.NewDecoder(response.Body)

	err = decoder.Decode(&result)

	return
}



func getVersionDetails(path, appName string, version uint, token string) (result AppData, err error) {

	task := http.Client{}

	req, err := http.NewRequest("GET",
		path+"/appbundles/"+appName+"/versions/"+strconv.Itoa(int(version)),
		nil,
	)

	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	response, err := task.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		content, _ := ioutil.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	decoder := json.NewDecoder(response.Body)

	err = decoder.Decode(&result)

	return
}



func deleteAppVersion(path string, appName string, version uint, token string) (err error) {

	task := http.Client{}
	req, err := http.NewRequest("DELETE",
		path+"/appbundles/"+appName+"/versions/"+strconv.Itoa(int(version)),
		nil,
	)

	req.Header.Set("Authorization", "Bearer "+token)
	response, err := task.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		content, _ := ioutil.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	return
}