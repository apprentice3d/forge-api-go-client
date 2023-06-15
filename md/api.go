package md

import (
	"encoding/base64"

	"github.com/woweh/forge-api-go-client"
	"github.com/woweh/forge-api-go-client/oauth"
)

// ModelDerivativeAPI struct holds all paths necessary to access Model Derivative API
type ModelDerivativeAPI struct {
	Authenticator       oauth.ForgeAuthenticator
	ModelDerivativePath string
	Region              forge.Region
}

// NewMDAPI returns a Model Derivative API client.
//
// NOTE: this uses the US region.
//
// Deprecated: Use NewMdApi instead.
func NewMDAPI(authenticator oauth.ForgeAuthenticator) ModelDerivativeAPI {
	return ModelDerivativeAPI{
		authenticator,
		"/modelderivative/v2/designdata",
		forge.US,
	}
}

// NewMdApi returns a Model Derivative API client for a specific region.
func NewMdApi(authenticator oauth.ForgeAuthenticator, region forge.Region) ModelDerivativeAPI {
	// default to US region
	path := "/modelderivative/v2/designdata"
	if region == forge.EMEA {
		path = "/modelderivative/v2/regions/eu/designdata"
	}

	return ModelDerivativeAPI{
		authenticator,
		path,
		region,
	}
}

// TranslateWithParamsAndXHeaders starts a translation job with the given TranslationParams and XAdsHeaders.
//
// References:
//   - https://aps.autodesk.com/en/docs/model-derivative/v2/reference/http/jobs/job-POST/
func (a *ModelDerivativeAPI) TranslateWithParamsAndXHeaders(
	params TranslationParams, xHeaders XAdsHeaders,
) (result TranslationResult, err error) {
	bearer, err := a.Authenticator.GetToken("data:write data:read")
	if err != nil {
		return
	}
	path := a.Authenticator.GetHostPath() + a.ModelDerivativePath
	result, err = translate(path, params, &xHeaders, bearer.AccessToken)

	return
}

// TranslateWithParams triggers translation job with settings specified in given TranslationParams
//
// Deprecated: Use TranslateWithParamsAndXHeaders instead.
func (a *ModelDerivativeAPI) TranslateWithParams(params TranslationParams) (result TranslationResult, err error) {
	bearer, err := a.Authenticator.GetToken("data:write data:read")
	if err != nil {
		return
	}
	path := a.Authenticator.GetHostPath() + a.ModelDerivativePath
	result, err = translate(path, params, nil, bearer.AccessToken)

	return
}

// TranslationSVFPreset specifies default parameters for translating a generic model into svf.
//   - Output: svf
//   - Views: 2d, 3d
//   - Destination: US
var TranslationSVFPreset = TranslationParams{
	Output: OutputSpec{
		Destination: DestSpec{forge.US},
		Formats: []FormatSpec{
			{
				Type:  SVF,
				Views: []OutputView{View2D, View3D},
			},
		},
	},
}

// NewTranslationParams creates a TranslationParams struct with the given objectID, outputType, views, and advanced options.
//   - The region will be taken from the ModelDerivativeAPI.
//   - The advanced options can be nil.
//
// Make sure to use the correct views and advanced options for the given outputType.
// There are no checks for this.
func (a *ModelDerivativeAPI) NewTranslationParams(
	urn string, outputType OutputType, views []OutputView, advanced *AdvancedSpec,
) TranslationParams {
	return TranslationParams{
		Input: InputSpec{
			URN: urn,
		},
		Output: OutputSpec{
			Destination: DestSpec{a.Region},
			Formats: []FormatSpec{
				{
					Type:     outputType,
					Views:    views,
					Advanced: advanced,
				},
			},
		},
	}
}

// UrnFromObjectId creates a Base64 (URL Safe) encoded URN from the given objectID.
func UrnFromObjectId(objectID string) string {
	return base64.RawStdEncoding.EncodeToString([]byte(objectID))
}

// IfcAdvancedSpec returns an IFC specific AdvancedSpec.
//
//	NOTE:
//
// The storeys, spaces, and openings options are applicable only when conversionMethod is set to `modern` or `v3`.
func IfcAdvancedSpec(conversionMethod ConversionMethod, storeys, spaces, openings Option) *AdvancedSpec {
	if conversionMethod == Legacy {
		return &AdvancedSpec{ConversionMethod: conversionMethod}
	}
	return &AdvancedSpec{
		ConversionMethod: conversionMethod,
		BuildingStoreys:  storeys,
		Spaces:           spaces,
		OpeningElements:  openings,
	}
}

// RevitAdvancedSpec returns a Revit specific AdvancedSpec.
func RevitAdvancedSpec(generateMasterViews *bool, materialMode MaterialMode) *AdvancedSpec {
	return &AdvancedSpec{
		GenerateMasterViews: generateMasterViews,
		MaterialMode:        materialMode,
	}
}

// NavisworksAdvancedSpec returns a Navisworks specific AdvancedSpec.
func NavisworksAdvancedSpec(hiddenObjects, basicMaterialProperties, autodeskMaterialProperties, timeLinerProperties *bool) *AdvancedSpec {
	return &AdvancedSpec{
		HiddenObjects:              hiddenObjects,
		BasicMaterialProperties:    basicMaterialProperties,
		AutodeskMaterialProperties: autodeskMaterialProperties,
		TimeLinerProperties:        timeLinerProperties,
	}
}

// ObjAdvancedSpec returns a OBJ specific AdvancedSpec.
func ObjAdvancedSpec(
	exportFileStructure ExportFileStructure, unit Unit, modelGuid string, objectIds *[]int,
) *AdvancedSpec {
	return &AdvancedSpec{
		ExportFileStructure: exportFileStructure,
		Unit:                unit,
		ModelGuid:           modelGuid,
		ObjectIds:           objectIds,
	}
}

// TranslateToSVF is a helper function for translating a file to SVF using the default TranslationSVFPreset.
// The objectID will be converted into a Base64 (URL Safe) encoded URN.
//
// Deprecated: Use TranslateWithParamsAndXHeaders instead.
func (a *ModelDerivativeAPI) TranslateToSVF(objectID string) (result TranslationResult, err error) {
	bearer, err := a.Authenticator.GetToken("data:write data:read")
	if err != nil {
		return
	}
	path := a.Authenticator.GetHostPath() + a.ModelDerivativePath
	params := TranslationSVFPreset
	params.Input.URN = UrnFromObjectId(objectID)

	result, err = translate(path, params, nil, bearer.AccessToken)

	return
}

// GetManifest returns information about derivatives that correspond to a specific source file,
// including derivative URNs and translation statuses.
func (a *ModelDerivativeAPI) GetManifest(urn string) (result Manifest, err error) {
	bearer, err := a.Authenticator.GetToken("data:read")
	if err != nil {
		return
	}
	path := a.Authenticator.GetHostPath() + a.ModelDerivativePath
	result, err = getManifest(path, urn, bearer.AccessToken)

	return
}

// GetDerivative downloads a selected derivative. To download the file, you need to specify the file’s URN, which you retrieve from the manifest.
// You can fetch the manifest using the GetManifest function.
func (a *ModelDerivativeAPI) GetDerivative(urn, derivativeUrn string) (jsonData []byte, err error) {
	bearer, err := a.Authenticator.GetToken("jsonData:read")
	if err != nil {
		return
	}
	path := a.Authenticator.GetHostPath() + a.ModelDerivativePath
	jsonData, err = getDerivative(path, urn, derivativeUrn, bearer.AccessToken)

	return
}

// GetMetadata returns a list of model views (Viewables) in the source design specified by the `urn` URI parameter.
// It also returns the ID that uniquely identifies the model view.
// You can use this ID with other metadata endpoints to obtain information about the objects within model view.
//
//	NOTE: You can retrieve metadata only from an input file that has been translated to SVF or SVF2.
//
// References:
//   - https://aps.autodesk.com/en/docs/model-derivative/v2/reference/http/metadata/urn-metadata-GET/
func (a *ModelDerivativeAPI) GetMetadata(urn string, xHeaders XAdsHeaders) (result MetadataResponse, err error) {
	bearer, err := a.Authenticator.GetToken("data:read")
	if err != nil {
		return
	}
	path := a.Authenticator.GetHostPath() + a.ModelDerivativePath
	result, err = getMetadata(path, urn, bearer.AccessToken, xHeaders)

	return
}

// GetModelViewProperties returns the properties of the objects in the model view as one json blob.
//   - You can get the guid (unique model view ID) by using the GetMetadata function.
//
// References:
//   - https://aps.autodesk.com/en/docs/model-derivative/v2/reference/http/metadata/urn-metadata-guid-properties-GET/
func (a *ModelDerivativeAPI) GetModelViewProperties(urn, guid string, xHeaders XAdsHeaders) (jsonData []byte, err error) {
	bearer, err := a.Authenticator.GetToken("data:read")
	if err != nil {
		return
	}

	path := a.Authenticator.GetHostPath() + a.ModelDerivativePath

	return getModelViewProperties(path, urn, guid, bearer.AccessToken, xHeaders)
}

// GetObjectTree returns the object tree of the model view.
//   - You can get the guid (unique model view ID) by using the GetMetadata function.
//   - Use forceGet = `true` to retrieve the object tree even if it exceeds the recommended maximum size (20 MB). The default for forceGet is `false`.
func (a *ModelDerivativeAPI) GetObjectTree(urn, guid string, forceGet bool, xHeaders XAdsHeaders) (result ObjectTree, err error) {
	bearer, err := a.Authenticator.GetToken("data:read")
	if err != nil {
		return
	}

	path := a.Authenticator.GetHostPath() + a.ModelDerivativePath

	return getObjectTree(path, urn, guid, bearer.AccessToken, forceGet, xHeaders)
}
