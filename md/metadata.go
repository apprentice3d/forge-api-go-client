package md

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type MetaData struct {
	Data struct {
		Type  string `json:"type,omitempty"`
		Views []View `json:"metadata,omitempty"`
	} `json:"data,omitempty"`
}

type View struct {
	Name         string   `json:"name,omitempty"`
	Role         ViewType `json:"role,omitempty"`
	Guid         string   `json:"guid,omitempty"`
	IsMasterView bool     `json:"isMasterView,omitempty"`
}

type ObjectTree struct {
	Data struct {
		Type    string           `json:"type"`
		Objects []ObjectTreeNode `json:"objects"`
	} `json:"data"`
}

type ObjectTreeNode struct {
	ObjectId int              `json:"objectid"`
	Name     string           `json:"name"`
	Objects  []ObjectTreeNode `json:"objects"`
}

const (
	timeToWait = time.Duration(5) * time.Second
	maxRetries = 12 * 5 // => 5 minutes max
)

func getMetadata(path, urn, token string) (result MetaData, err error) {
	task := http.Client{}

	req, err := http.NewRequest("GET", path+"/"+urn+"/metadata", nil)
	if err != nil {
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)

	response, err := task.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		content, _ := io.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	err = json.NewDecoder(response.Body).Decode(&result)

	return
}

func getObjectTree(path, urn, modelGuid, token string, forceGet bool, xHeaders XAdsHeaders) (
	result ObjectTree, err error,
) {
	// retry logic, not very elegant but it works
	tries := 0
retry:
	tries++
	task := http.Client{}

	url := path + "/" + urn + "/metadata/" + modelGuid + "?forceget=" + strconv.FormatBool(forceGet)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Add("x-ads-force", strconv.FormatBool(xHeaders.Overwrite))
	req.Header.Add("x-ads-derivative-format", string(xHeaders.Format))

	response, err := task.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusAccepted {
		// 202 Accepted => the request has been accepted for processing, but the processing has not been completed.
		if tries < maxRetries {
			log.Println("=> retry number: ", tries)
			time.Sleep(timeToWait)
			goto retry
		}
	} else if response.StatusCode != http.StatusOK {
		content, _ := io.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	err = json.NewDecoder(response.Body).Decode(&result)

	return
}

func getModelViewProperties(path, urn, modelGuid, token string, xHeaders XAdsHeaders) (
	jsonData []byte, err error,
) {

	// retry logic, not very elegant but it works
	tries := 0
retry:
	tries++
	task := http.Client{}

	req, err := http.NewRequest("GET", path+"/"+urn+"/metadata/"+modelGuid+"/properties", nil)
	if err != nil {
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Add("x-ads-force", strconv.FormatBool(xHeaders.Overwrite))
	req.Header.Add("x-ads-derivative-format", string(xHeaders.Format))

	response, err := task.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusAccepted {
		// 202 Accepted => the request has been accepted for processing, but the processing has not been completed.
		if tries < maxRetries {
			log.Println("=> retry number: ", tries)
			time.Sleep(timeToWait)
			goto retry
		}
	} else if response.StatusCode != http.StatusOK {
		content, _ := io.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	return io.ReadAll(response.Body)
}
