package da

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

type EngineList struct {
	InfoList
}

type EngineDetails struct {
	ProductVersion string `json:"productVersion"`
	Description string `json:"description"`
	Version uint `json:"version"`
	Id string `json:"id"`
}



func listEngines(path string, token string) (list EngineList, err error) {

	task := http.Client{}
	req, err := http.NewRequest("GET",
		path+"/engines",
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


func getEngineDetails(path string, engineID string, token string) (details EngineDetails, err error) {

	task := http.Client{}
	req, err := http.NewRequest("GET",
		path+"/engines/"+engineID,
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
	err = decoder.Decode(&details)

	return
}