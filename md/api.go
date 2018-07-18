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
					[]string{"2d", "3d"},
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
type FormatSpec struct {
	Type  string   `json:"type"`
	Views []string `json:"views"`
}

// TranslateWithParams triggers translation job with settings specified in given TranslationParams
func (a ModelDerivativeAPI) TranslateWithParams(client * http.Client, params TranslationParams) (result TranslationResult, err error) {
	bearer, err := a.Authenticate(client, "data:write data:read")
	if err != nil {
		return
	}
	path := a.Host + a.ModelDerivativePath
	result, err = translate(client, path, params, bearer.AccessToken)

	return
}

// TranslateToSVF is a helper function that will use the TranslationSVFPreset for translating into svf a given ObjectID.
// It will also take care of converting objectID into Base64 (URL Safe) encoded URN.
func (a ModelDerivativeAPI) TranslateToSVF(client * http.Client, objectID string) (result TranslationResult, err error) {
	bearer, err := a.Authenticate(client, "data:write data:read")
	if err != nil {
		return
	}
	path := a.Host + a.ModelDerivativePath
	params := TranslationSVFPreset
	params.Input.URN = base64.RawStdEncoding.EncodeToString([]byte(objectID))

	result, err = translate(client, path, params, bearer.AccessToken)

	return
}


func (a ModelDerivativeAPI) GetManifest(client * http.Client, urn string) (result ManifestResult, err error) {
	bearer, err := a.Authenticate(client, "data:read")
	if err != nil {
		return
	}

	path := a.Host + a.ModelDerivativePath
	result, err = getManifest(client, path, urn, bearer.AccessToken)

	return
}

/*
 *	SUPPORT FUNCTIONS
 */
func translate(client * http.Client, path string, params TranslationParams, token string) (result TranslationResult, err error) {

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

func getManifest(client * http.Client, path string, urn string, token string) (result ManifestResult, err error) {
/*
	byteParams, err := json.Marshal(params)
	if err != nil {
		log.Println("Could not marshal the translation parameters")
		return
	}*/

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
