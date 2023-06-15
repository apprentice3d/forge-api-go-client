package md

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type derivativeDownloadUrl struct {
	Etag        string `json:"etag"`
	Size        int    `json:"size"`
	Url         string `json:"url"`
	ContentType string `json:"content-type"`
	Expiration  int64  `json:"expiration"`
}

func getDerivative(path, urn, derivativeUrn, token string) (result []byte, err error) {

	task := http.Client{}

	req, err := http.NewRequest("GET", path+"/"+urn+"/manifest/"+derivativeUrn+"/signedcookies", nil)

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

	var getUrlResult derivativeDownloadUrl

	// deserialize the response
	err = json.NewDecoder(response.Body).Decode(&getUrlResult)
	if err != nil {
		return
	}

	signedCookieValue := strings.Join(response.Header.Values("Set-Cookie"), ";")

	return downloadDerivative(getUrlResult, signedCookieValue)
}

func downloadDerivative(downloadUrl derivativeDownloadUrl, cookieValue string) (result []byte, err error) {

	task := http.Client{}

	req, err := http.NewRequest("GET", downloadUrl.Url, nil)
	if err != nil {
		return
	}

	req.Header.Set("Cookie", cookieValue)
	req.Header.Set("Content-Type", downloadUrl.ContentType)
	req.Header.Set("Content-Length", strconv.Itoa(downloadUrl.Size))
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

	result, err = io.ReadAll(response.Body)
	if err != nil {
		return
	}
	if len(result) != downloadUrl.Size {
		err = errors.New("downloaded file size is different than the expected size")
		return
	}

	return result, nil
}

// BUG: When translating a non-Revit model, the
// Manifest will contain an array of strings as message,
// while in case of others it is just a string

type Manifest struct {
	Type         string       `json:"type"`
	HasThumbnail string       `json:"hasThumbnail"`
	Status       string       `json:"status"`
	Progress     string       `json:"progress"`
	Region       string       `json:"region"`
	URN          string       `json:"urn"`
	Version      string       `json:"version"`
	Derivatives  []Derivative `json:"derivatives"`
}

type Derivative struct {
	Name         string      `json:"name"`
	HasThumbnail string      `json:"hasThumbnail"`
	Status       string      `json:"status"`
	Progress     string      `json:"progress"`
	Messages     []Message   `json:"messages,omitempty"`
	OutputType   string      `json:"outputType"`
	Properties   *Properties `json:"properties,omitempty"`
	Children     []Child     `json:"children"`
}

type Message struct {
	Type string `json:"type"`
	Code string `json:"code"`
	// Message can either be a string, or an array of strings
	// Use reflection to handle this, for example:
	// reflect.TypeOf(result.Derivatives[0].Messages[0].Message).Kind() == reflect.String
	Message any `json:"message,omitempty"`
}

type Properties struct {
	DocumentInformation DocumentInformation `json:"Document Information"`
}

type DocumentInformation struct {
	NavisworksFileCreator string `json:"Navisworks File Creator"`
	IFCApplicationName    string `json:"IFC Application Name"`
	IFCApplicationVersion string `json:"IFC Application Version"`
	IFCSchema             string `json:"IFC Schema"`
	IFCLoader             string `json:"IFC Loader"`
}

type Child struct {
	GUID         string    `json:"guid"`
	Type         string    `json:"type"`
	Role         string    `json:"role"`
	Name         string    `json:"name,omitempty"`
	Status       string    `json:"status,omitempty"`
	Progress     string    `json:"progress,omitempty"`
	Mime         string    `json:"mime,omitempty"`
	UseAsDefault *bool     `json:"useAsDefault,omitempty"`
	HasThumbnail string    `json:"hasThumbnail,omitempty"`
	URN          string    `json:"urn,omitempty"`
	ViewableID   string    `json:"viewableID,omitempty"`
	PhaseNames   string    `json:"phaseNames,omitempty"`
	Resolution   []float32 `json:"resolution,omitempty"`
	Children     []Child   `json:"children,omitempty"`
	Camera       []float32 `json:"camera,omitempty"`
	ModelGUID    *string   `json:"modelGuid,omitempty"`
	ObjectIDs    []int     `json:"objectIds,omitempty"`
	Messages     []Message `json:"messages,omitempty"`
}

func getManifest(path, urn, token string) (result Manifest, err error) {
	task := http.Client{}

	req, err := http.NewRequest("GET", path+"/"+urn+"/manifest", nil)

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

func getMetadata(path, urn, token string, xHeaders XAdsHeaders) (result MetadataResponse, err error) {
	task := http.Client{}

	req, err := http.NewRequest("GET", path+"/"+urn+"/metadata", nil)

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

	if response.StatusCode != http.StatusOK {
		content, _ := io.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	err = json.NewDecoder(response.Body).Decode(&result)

	return
}

func getModelViewProperties(path, urn, guid, token string, xHeaders XAdsHeaders) (
	jsonData []byte, err error,
) {
	task := http.Client{}

	req, err := http.NewRequest("GET", path+"/"+urn+"/metadata", nil)
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

	if response.StatusCode != http.StatusOK {
		content, _ := io.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	return io.ReadAll(response.Body)
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

func getObjectTree(path, urn, guid, token string, forceGet bool, xHeaders XAdsHeaders) (
	result ObjectTree, err error,
) {
	task := http.Client{}

	url := path + "/" + urn + "/metadata/" + guid + "?forceget=" + strconv.FormatBool(forceGet)

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

	if response.StatusCode != http.StatusOK {
		content, _ := io.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	err = json.NewDecoder(response.Body).Decode(&result)

	return
}
