package md

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/woweh/forge-api-go-client"
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

// IsSuccess returns true if the status is success.
func (s *Status) IsSuccess() bool {
	// case insensitive comparison!
	return strings.EqualFold(string(*s), string(StatusSuccess))
}

// IsFailed returns true if the status is failed.
func (s *Status) IsFailed() bool {
	// case insensitive comparison!
	return strings.EqualFold(string(*s), string(StatusFailed))
}

// IsPending returns true if the status is pending.
func (s *Status) IsPending() bool {
	// case insensitive comparison!
	return strings.EqualFold(string(*s), string(StatusPending))
}

// IsInProgress returns true if the status is in progress.
func (s *Status) IsInProgress() bool {
	// case insensitive comparison!
	return strings.EqualFold(string(*s), string(StatusInProgress))
}

// IsTimeout returns true if the status is timeout.
func (s *Status) IsTimeout() bool {
	// case insensitive comparison!
	return strings.EqualFold(string(*s), string(StatusTimeout))
}

// IsEmpty returns true if this ProgressReport is empty.
func (pr *ProgressReport) IsEmpty() bool {
	return pr.Status == "" && pr.Progress == ""
}

type ProgressReport struct {
	Status   Status `json:"status"`
	Progress string `json:"progress"`
}

type Manifest struct {
	ProgressReport
	Type         string       `json:"type"`
	HasThumbnail string       `json:"hasThumbnail"`
	Region       forge.Region `json:"region"`
	URN          string       `json:"urn"`
	Version      string       `json:"version"`
	Derivatives  []Derivative `json:"derivatives"`
}

type Derivative struct {
	ProgressReport
	Name         string      `json:"name"`
	HasThumbnail string      `json:"hasThumbnail"`
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
	ProgressReport
	GUID         string    `json:"guid"`
	Type         string    `json:"type"`
	Role         string    `json:"role"`
	Name         string    `json:"name,omitempty"`
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
	log.Println("- Base64 encoded design URL: ", urn)
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
