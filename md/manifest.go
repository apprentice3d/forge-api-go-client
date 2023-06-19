package md

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
)

// BUG: When translating a non-Revit model, the
// Manifest will contain an array of strings as message,
// while in case of others it is just a string

// Status is the status of the translation
type Status string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "inprogress"
	StatusSuccess    Status = "success"
	StatusFailed     Status = "failed"
	StatusTimeout    Status = "timeout"
)

type Manifest struct {
	Type         string       `json:"type"`
	HasThumbnail string       `json:"hasThumbnail"`
	Status       Status       `json:"status"`
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
	// Message can either be a string, or an array of strings.
	// This is a bug in the REST API.
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

	log.Println("Requesting manifest...")
	log.Println("- Base64  encoded design URN: ", urn)
	log.Println("- URL: ", req.URL.String())

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

	log.Println("Manifest received.")

	return
}
