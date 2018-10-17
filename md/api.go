package md

import (
	"github.com/apprentice3d/forge-api-go-client/oauth"
	"encoding/base64"
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

// TranslationSVFPreset specifies the minimum necessary for translating a generic (single file, uncompressed)
// model into svf.
var TranslationSVFPreset = TranslationParams{
	Output: OutputSpec{
		Destination: DestSpec{"us"},
		Formats: []FormatSpec{
			FormatSpec{
				"svf",
				[]string{"2d", "3d"},
			},
		},
	},
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


// GetManifest returns information about derivatives that correspond to a specific source file,
// including derivative URNs and statuses.
func (a ModelDerivativeAPI) GetManifest(urn string) (result Manifest, err error) {
	bearer, err := a.Authenticate("data:read")
	if err != nil {
		return
	}
	path := a.Host + a.ModelDerivativePath
	result, err = getManifest(path, urn, bearer.AccessToken)

	return
}


// GetDerivative downloads a selected derivative. To download the file, you need to specify the file’s URN,
// which you retrieve by calling the GET :urn/manifest endpoint.
func (a ModelDerivativeAPI) GetDerivative(urn, derivativeUrn string) (data []byte, err error) {
	bearer, err := a.Authenticate("data:read")
	if err != nil {
		return
	}
	path := a.Host + a.ModelDerivativePath
	data, err = getDerivative(path, urn, derivativeUrn, bearer.AccessToken)

	return
}
