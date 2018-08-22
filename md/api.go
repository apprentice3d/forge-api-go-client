package md

import (
	"github.com/outer-labs/forge-api-go-client/oauth"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"bytes"
	"io/ioutil"
	"errors"
	"strconv"
	"io"
)

var (
	// TranslationSVFPreset specifies the minimum necessary for translating a generic (single file, uncompressed)
	// model into svf.
	TranslationSVFPreset = TranslationParams{
		Output: OutputSpec{
			Destination:DestSpec{"us"},
			Formats:[]FormatSpec{
				FormatSpec{
					"svf",
					[]string{"2d","3d"},
				},
			},
		},
	}
)

// API struct holds all paths necessary to access Model Derivative API
type ModelDerivativeAPI struct {
	oauth.TwoLeggedAuth
	ModelDerivativePath string
}

// NewAPIWithCredentials returns a Model Derivative API client with default configurations
func NewAPIWithCredentials(ClientID string, ClientSecret string) ModelDerivativeAPI {
	return ModelDerivativeAPI{
		oauth.NewTwoLeggedClient(ClientID, ClientSecret),
		"/modelderivative/v2/designdata",
	}
}

//TranslationParams is used when specifying the translation jobs
type TranslationParams struct {
	Input struct {
		URN           string  `json:"urn"`
		CompressedURN *bool   `json:"compressedUrn,omitempty"`
		RootFileName  *string `json:"rootFileName,omitempty"`
	} `json:"input"`
	Output OutputSpec `json:"output"`
}

// TranslationResult reflects data received upon successful creation of translation job
type TranslationResult struct {
	Result string `json:"result"`
	URN    string `json:"urn"`
	AcceptedJobs struct {
		Output OutputSpec `json:"output"`
	}
}

type ManifestResult struct{
	Type string `json:"type,omitempty"`
	HasThumbnail bool `json:"hasThumbnail,string,omitempty"`
	Status string `json:"status,omitempty"`
	Progress string `json:"progress,omitempty"`
	Region string `json:"region,omitempty"`
	URN string `json:"urn,omitempty"`
	Derivatives []DerivativeSpec `json:"derivatives,omitempty"`
}

type DerivativeSpec struct{
	Name string `json:"name,omitempty"`
	HasThumbnail bool `json:"hasThumbnail,string,omitempty"`
	Role string `json:"role,omitempty"`
	Status string `json:"status,omitempty"`
	Progress string `json:"progress,omitempty"`
	Children []ChildrenSpec `json:"children,omitempty"`
}

type ChildrenSpec struct{
	GUID string `json:"guid,omitempty"`
	Role string `json:"role,omitempty"`
	MIME string `json:"mime,omitempty"`
	URN string `json:"urn,omitempty"`
	Progress string `json:"progress,omitempty"`
	Status string `json:"status,omitempty"`
}

// OutputSpec reflects data found upon creation translation job and receiving translation job status
type OutputSpec struct {
	Destination DestSpec     `json:"destination,omitempty"`
	Formats     []FormatSpec `json:"formats"`
}

// DestSpec is used within OutputSpecs and is useful when specifying the region for translation results
type DestSpec struct {
	Region string `json:"region"`
}

// FormatSpec is used within OutputSpecs and should be used when specifying the expected format and views (2d or/and 3d)
type FormatSpec struct{
	Type  string   `json:"type"`
	Views []string `json:"views"`
}

type MetadataResult struct{
	Data MetadataSpec `json:"data",omitempty`
}

type MetadataSpec struct{
	Type string `json:"type",omitempty`
	Metadata []ViewSpec `json:"metadata",omitempty`
}

type ViewSpec struct{
	Name string `json:"name",omitempty`
	Role string `json:"role",omitempty`
	Guid string `json:"guid",omitempty`
}

type PropertiesResult struct{
	Data PropertiesSpec `json:"data",omitempty`
	Result string `json:"result",omitempty`
}

type PropertiesSpec struct{
	Type string    `json:"type"`
	Collection []ObjectSpec `json:"collection"`
}

type ObjectSpec struct{
	ObjectID int64 `json:"objectid"`
	Name string     `json:"name"`
	ExternalID string `json:"externalId"`
	Properties json.RawMessage
}

// TranslateWithParams triggers translation job with settings specified in given TranslationParams
func (a ModelDerivativeAPI) TranslateWithParams(params TranslationParams) (result TranslationResult, err error) {
	bearer, err := a.Authenticate("data:write data:read")
	if err != nil {
		return
	}
	path := a.Host + a.ModelDerivativePath
	result, err = translate(path, params, bearer.AccessToken)

	return
}

// TranslateToSVF is a helper function that will use the TranslationSVFPreset for translating into svf a given ObjectID.
// It will also take care of converting objectID into Base64 (URL Safe) encoded URN.
func (a ModelDerivativeAPI) TranslateToSVF(objectID string) (result TranslationResult, err error) {
	bearer, err := a.Authenticate("data:write data:read")
	if err != nil {
		return
	}
	path := a.Host + a.ModelDerivativePath
	params := TranslationSVFPreset
	params.Input.URN = base64.RawStdEncoding.EncodeToString([]byte(objectID))

	result, err = translate(path, params, bearer.AccessToken)

	return
}


func (a ModelDerivativeAPI) GetManifest(urn string) (result ManifestResult, err error) {
	bearer, err := a.Authenticate("data:read")
	if err != nil {
		return
	}

	path := a.Host + a.ModelDerivativePath
	result, err = getManifest(path, urn, bearer.AccessToken)

	return
}

func (a ModelDerivativeAPI) GetMetadata(urn string) (result MetadataResult, err error) {
	bearer, err := a.Authenticate("data:read")
	if err != nil {
		return
	}

	path := a.Host + a.ModelDerivativePath
	result, err = getMetadata(path, urn, bearer.AccessToken)

	return
}

func (a ModelDerivativeAPI) GetProperties(urn string, viewId string) (result PropertiesResult, err error) {
	bearer, err := a.Authenticate("data:read")
	if err != nil {
		return
	}

	path := a.Host + a.ModelDerivativePath
	result, err = getProperties(path, urn, viewId, bearer.AccessToken)

	return
}

func (a ModelDerivativeAPI) GetThumbnail(urn string) (reader io.ReadCloser, err error) {
	bearer, err := a.Authenticate("data:read")
	if err != nil {
		return
	}

	path := a.Host + a.ModelDerivativePath
	reader, err = getThumbnail(path, urn, bearer.AccessToken)

	return
}

/*
 *	SUPPORT FUNCTIONS
 */
func translate(path string, params TranslationParams, token string) (result TranslationResult, err error) {
	client := http.Client{}
	byteParams, err := json.Marshal(params)
	if err != nil {
		log.Println("Could not marshal the translation parameters")
		return
	}

	req, err := http.NewRequest("POST",
		path+"/job",
		bytes.NewBuffer(byteParams))

	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	response, err := client.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusOK {
		content, _ := ioutil.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	decoder := json.NewDecoder(response.Body)

	err = decoder.Decode(&result)

	return
}

func getManifest(path string, urn string, token string) (result ManifestResult, err error) {
	client := http.Client{}

	req, err := http.NewRequest("GET",
		path + "/" + urn + "/manifest",
		nil)

	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	response, err := client.Do(req)
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

func getThumbnail(path string, urn string, token string) (reader io.ReadCloser, err error) {
	client := http.Client{}

	req, err := http.NewRequest("GET",
		path + "/" + urn + "/thumbnail",
		nil)

	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	response, err := client.Do(req)
	if err != nil {
		return
	}

	if response.StatusCode != http.StatusOK {
		content, _ := ioutil.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	reader = response.Body
	return
}

func getProperties(path string, urn string, viewId string, token string) (
	result PropertiesResult, err error) {
	client := http.Client{}

	req, err := http.NewRequest("GET",
		path + "/" + urn + "/metadata/" + viewId + "/properties",
		nil)

	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	response, err := client.Do(req)
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

func getMetadata(path string, urn string, token string) (
	result MetadataResult, err error) {
	client := http.Client{}

	req, err := http.NewRequest("GET",
		path + "/" + urn + "/metadata",
		nil)

	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	response, err := client.Do(req)
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
