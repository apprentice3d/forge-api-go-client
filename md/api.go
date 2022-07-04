package md

import (
	"encoding/base64"

	"github.com/apprentice3d/forge-api-go-client/md/xAdsHeaders"
	"github.com/apprentice3d/forge-api-go-client/oauth"
)

// ModelDerivativeAPI struct holds all paths necessary to access Model Derivative API
type ModelDerivativeAPI struct {
	Authenticator       oauth.ForgeAuthenticator
	ModelDerivativePath string
}

// NewMDAPI returns a Model Derivative API client with default configurations
func NewMDAPI(authenticator oauth.ForgeAuthenticator) ModelDerivativeAPI {
	return ModelDerivativeAPI{
		authenticator,
		"/modelderivative/v2/designdata",
	}
}

// TranslateWithParamsAndXHeaders triggers translation job with settings specified in given TranslationParams
func (a ModelDerivativeAPI) TranslateWithParamsAndXHeaders(params TranslationParams, xHeaders xAdsHeaders.Headers) (result TranslationResult, err error) {
	bearer, err := a.Authenticator.GetToken("data:write data:read")
	if err != nil {
		return
	}
	path := a.Authenticator.GetHostPath() + a.ModelDerivativePath
	result, err = translate(path, params, xHeaders, bearer.AccessToken)

	return
}

// TranslateWithParams triggers translation job with settings specified in given TranslationParams
func (a ModelDerivativeAPI) TranslateWithParams(params TranslationParams) (result TranslationResult, err error) {
	return a.TranslateWithParamsAndXHeaders(params, xAdsHeaders.Default())
}

// TranslationSVFPreset specifies the minimum necessary for translating a generic (single file, uncompressed)
// model into svf.
var TranslationSVFPreset = TranslationParams{
	Output: OutputSpec{
		Destination: DestSpec{"us"},
		Formats: []FormatSpec{
			{
				Type:  "svf",
				Views: []string{"2d", "3d"},
			},
		},
	},
}

// TranslateToSVF is a helper function that will use the TranslationSVFPreset for translating into svf a given ObjectID.
// It will also take care of converting objectID into Base64 (URL Safe) encoded URN.
func (a ModelDerivativeAPI) TranslateToSVF(objectID string) (result TranslationResult, err error) {
	bearer, err := a.Authenticator.GetToken("data:write data:read")
	if err != nil {
		return
	}
	path := a.Authenticator.GetHostPath() + a.ModelDerivativePath
	params := TranslationSVFPreset
	params.Input.URN = base64.RawStdEncoding.EncodeToString([]byte(objectID))

	result, err = translate(path, params, xAdsHeaders.Default(), bearer.AccessToken)

	return
}

// GetManifest returns information about derivatives that correspond to a specific source file,
// including derivative URNs and statuses.
func (a ModelDerivativeAPI) GetManifest(urn string) (result Manifest, err error) {
	bearer, err := a.Authenticator.GetToken("data:read")
	if err != nil {
		return
	}
	path := a.Authenticator.GetHostPath() + a.ModelDerivativePath
	result, err = getManifest(path, urn, bearer.AccessToken)

	return
}

// GetDerivative downloads a selected derivative. To download the file, you need to specify the file’s URN,
// which you retrieve by calling the GET :urn/manifest endpoint.
func (a ModelDerivativeAPI) GetDerivative(urn, derivativeUrn string) (data []byte, err error) {
	bearer, err := a.Authenticator.GetToken("data:read")
	if err != nil {
		return
	}
	path := a.Authenticator.GetHostPath() + a.ModelDerivativePath
	data, err = getDerivative(path, urn, derivativeUrn, bearer.AccessToken)

	return
}

// GetMetadata returns a list of model views (Viewables) in the source design specified by the urn URI parameter.
// It also returns the ID that uniquely identifies the model view.
// You can use this ID with other metadata endpoints to obtain information about the objects within model view.
//  NOTE:
// You can retrieve metadata only from an input file that has been translated to SVF or SVF2.
func (a ModelDerivativeAPI) GetMetadata(urn string, xHeaders xAdsHeaders.Headers) (result MetadataResponse, err error) {
	bearer, err := a.Authenticator.GetToken("data:read")
	if err != nil {
		return
	}
	path := a.Authenticator.GetHostPath() + a.ModelDerivativePath
	result, err = getMetadata(path, urn, bearer.AccessToken, xHeaders)

	return
}
