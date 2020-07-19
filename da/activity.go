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

type Param struct {
	Zip         bool   `json:"zip"`
	Description string `json:"description"`
	OnDemand    bool   `json:"ondemand"`
	Required    bool   `json:"required"`
	Verb        string `json:"verb"`
	LocalName   string `json:"localName"`
}

type Setting struct {
	Script string `json:"script"`
}

type ActivityConfig struct {
	ID          string           `json:"id"`
	CommandLine []string         `json:"commandLine"`
	Description string           `json:"description"`
	AppBundles  []string         `json:"appbundles"`
	Engine      string           `json:"engine"`
	Parameters  map[string]Param `json:"paramaters"`
	Settings    Setting          `json:"settings"`
}

type Activity struct {
	ActivityConfig

	authenticator oauth.ForgeAuthenticator
	path          string
	name          string

}

func (activity *Activity) Delete() (err error) {
	bearer, err := activity.authenticator.GetToken("code:all")
	if err != nil {
		return
	}

	err = deleteActivity(activity.path, activity.ID, bearer.AccessToken)

	activity.Parameters = make(map[string]Param)
	activity.ID = ""
	activity.CommandLine = make([]string,0)
	activity.Description = ""
	activity.AppBundles = make([]string,0)
	activity.Engine = ""
	activity.Settings = Setting{}
	activity.authenticator = nil
	activity.path = ""
	activity.name = ""

	return
}


//CreateAlias creates a new alias for this Activity.
func (activity Activity) CreateAlias(alias string, version uint) (result Alias, err error) {
	bearer, err := activity.authenticator.GetToken("code:all")
	if err != nil {
		return
	}
	result, err = createActivityAlias(activity.path, activity.name, alias, version, bearer.AccessToken)

	return
}


/*
 *	SUPPORT FUNCTIONS
 */


/*
  ACTIVITY
*/

func createActivity(path string, activity ActivityConfig, token string) (result Activity, err error) {

	task := http.Client{}

	body, err := json.Marshal(
		activity)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST",
		path+"/activities",
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

func deleteActivity(path string, activityId string, token string) (err error) {

	task := http.Client{}
	req, err := http.NewRequest("DELETE",
		path+"/activities/"+activityId,
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

func listActivityAliases(path string, activityId, token string) (list AliasesList, err error) {

	task := http.Client{}
	req, err := http.NewRequest("GET",
		path+"/activities/"+activityId+"/aliases",
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

func createActivityAlias(path, activityId, alias string, version uint, token string) (result Alias, err error) {

	task := http.Client{}

	body, err := json.Marshal(
		Alias{
			alias,
			version,
		})
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST",
		path+"/activities/"+activityId+"/aliases",
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

func modifyActivityAlias(path, activityId, alias string, version uint, token string) (result Alias, err error) {

	task := http.Client{}

	body, err := json.Marshal(
		struct {
			Version uint `json:"version"`
		}{version})
	if err != nil {
		return
	}

	req, err := http.NewRequest("PATCH",
		path+"/activities/"+activityId+"/aliases/"+alias,
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

func getActivityAliasDetails(path, activityId, alias, token string) (result Alias, err error) {

	task := http.Client{}

	req, err := http.NewRequest("GET",
		path+"/activities/"+activityId+"/aliases/"+alias,
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



func deleteActivityAlias(path string, activityId, alias, token string) (err error) {

	task := http.Client{}
	req, err := http.NewRequest("DELETE",
		path+"/activities/"+activityId+"/aliases/"+alias,
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

func listActivityVersions(path string, activityId, token string) (list VersionList, err error) {

	task := http.Client{}
	req, err := http.NewRequest("GET",
		path+"/activities/"+activityId+"/versions",
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

func createActivityVersion(path, activityId, engine string, token string) (result ActivityConfig, err error) {

	task := http.Client{}

	body, err := json.Marshal(
		struct{
			Engine string `json:"engine"`
		}{engine})
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST",
		path+"/activities/"+activityId+"/versions",
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



func getActivityVersionDetails(path, activityId string, version uint, token string) (result ActivityConfig, err error) {

	task := http.Client{}

	req, err := http.NewRequest("GET",
		path+"/activities/"+activityId+"/versions/"+strconv.Itoa(int(version)),
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



func deleteActivityVersion(path, activityId string, version uint, token string) (err error) {

	task := http.Client{}
	req, err := http.NewRequest("DELETE",
		path+"/activities/"+activityId+"/versions/"+strconv.Itoa(int(version)),
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