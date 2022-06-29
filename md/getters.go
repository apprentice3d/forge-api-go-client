package md

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

func getDerivative(path string, urn, derivativeUrn, token string) (result []byte, err error) {
	task := http.Client{}

	req, err := http.NewRequest("GET",
		path+"/"+urn+"/manifest/"+derivativeUrn,
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

	if response.StatusCode != http.StatusOK {
		content, _ := ioutil.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	result, err = ioutil.ReadAll(response.Body)

	return
}

type Manifest struct {
	Type         string       `json:"type"`
	HasThumbnail string       `json:"hasThumbnail"`
	Status       string       `json:"status"`
	Progress     string       `json:"progress"`
	Region       string       `json:"region"`
	URN          string       `json:"urn"`
	Derivatives  []Derivative `json:"derivatives"`
}

type Derivative struct {
	Name         string `json:"name"`
	HasThumbnail string `json:"hasThumbnail"`
	Status       string `json:"status"`
	Progress     string `json:"progress"`
	Messages     []struct {
		Type string `json:"type"`
		//Message string `json:"message"`
		Code string `json:"code"`
	} `json:"messages,omitempty"`
	OutputType string  `json:"outputType"`
	Children   []Child `json:"children"`
}

//BUG: When translating a non-Revit model, the
// Manifest will contain an array of strings as message,
// while in case of others it is just a string

type Child struct {
	GUID         string    `json:"guid"`
	Type         string    `json:"type"`
	Role         string    `json:"role"`
	Name         string    `json:"name,omitempty"`
	Status       string    `json:"status,omitempty"`
	Progress     string    `json:"progress,omitempty"`
	Mime         string    `json:"mime,omitempty"`
	HasThumbnail string    `json:"hasThumbnail,omitempty"`
	URN          string    `json:"urn,omitempty"`
	ViewableID   string    `json:"viewableID,omitempty"`
	PhaseNames   string    `json:"phaseNames,omitempty"`
	Resolution   []float32 `json:"resolution,omitempty"`
	Children     []Child   `json:"children,omitempty"`
	Camera       []float32 `json:"camera,omitempty"`
	Messages     []struct {
		Type    string   `json:"type"`
		Message []string `json:"message"`
		Code    string   `json:"code"`
	} `json:"messages,omitempty"`
}

func getManifest(path string, urn, token string) (result Manifest, err error) {
	task := http.Client{}

	req, err := http.NewRequest("GET",
		path+"/"+urn+"/manifest",
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

	if response.StatusCode != http.StatusOK {
		content, _ := ioutil.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&result)

	return
}

type MetadataResponse struct {
	Data struct {
		Type     string `json:"type,omitempty"`
		Metadata []struct {
			Name         string `json:"name,omitempty"`
			Role         string `json:"role,omitempty"`
			GUID         string `json:"guid,omitempty"`
			IsMasterView bool   `json:"isMasterView,omitempty"`
		} `json:"metadata,omitempty"`
	} `json:"data,omitempty"`
}

func getMetadata(path string, urn, token string, xHeaders XHeaders) (result MetadataResponse, err error) {
	task := http.Client{}

	req, err := http.NewRequest("GET",
		path+"/"+urn+"/metadata",
		nil,
	)

	if err != nil {
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Add("x-ads-derivative-format", string(xHeaders.Format))
	req.Header.Add("x-ads-force", strconv.FormatBool(xHeaders.Overwrite))

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

/****************************    **********************/

type LMVManifest struct {
	Name           string `json:"name"`
	ToolkitVersion string `json:"toolkitversion"`
	ADSKID         struct {
		SourceSystem string `json:"sourcesystem"`
		Type         string `json:"type"`
		ID           string `json:"id"`
		Version      string `json:"version"`
	} `json:"adskID"`
	Assets   []Asset   `json:"assets"`
	Typesets []Typeset `json:"typesets"`
}

type Asset struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	URI   string `json:"URI"`
	Size  uint64 `json:"size"`
	USize uint64 `json:"usize"`
}

type Typeset struct {
	ID    string `json:"id"`
	Types []Type `json:"types"`
}

type Type struct {
	Class   string `json:"class"`
	Type    string `json:"type"`
	Version uint   `json:"version"`
}
