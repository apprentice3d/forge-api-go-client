package md

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

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

// getModelViewProperties returns a list of model views (Viewables) in the source design specified by the urn URI parameter.
// It also returns the ID that uniquely identifies the model view.
// You can use this ID with other metadata endpoints to obtain information about the objects within model view.
//
// Most design applications like Fusion 360 and Inventor contain only one model view per design. However, some applications like Revit allow multiple model views (e.g., HVAC, architecture, perspective) per design.
//
// Note You can retrieve metadata only from an input file that has been translated to SVF or SVF2.
//
// Reference:
// - https://aps.autodesk.com/en/docs/model-derivative/v2/reference/http/metadata/urn-metadata-GET/
func getModelViewProperties(path, urn, token string, xHeaders XAdsHeaders) (
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
