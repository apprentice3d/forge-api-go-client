package md

import (
	"encoding/base64"
	"io"
	"strings"

	"github.com/woweh/forge-api-go-client"
	"github.com/woweh/forge-api-go-client/oauth"
)

// ModelDerivativeAPI struct holds all paths necessary to access Model Derivative API
type ModelDerivativeAPI struct {
	// Forge authenticator, used to get access token, either 2-legged or 3-legged
	Authenticator oauth.ForgeAuthenticator
	// The relativePath depends on the region => either usPath or euPath
	relativePath string
	// The region where data resides, either US or EU (EMEA)
	region forge.Region
}

const (
	usPath = "/modelderivative/v2/designdata"
	euPath = "/modelderivative/v2/regions/eu/designdata"
)

// NewMdApi returns a Model Derivative API client for a specific region.
func NewMdApi(authenticator oauth.ForgeAuthenticator, region forge.Region) ModelDerivativeAPI {
	// default to US region
	path := usPath
	if region == forge.EU {
		path = euPath
	}

	return ModelDerivativeAPI{
		Authenticator: authenticator,
		relativePath:  path,
		region:        region,
	}
}

// Region of the ModelDerivativeAPI.
func (a *ModelDerivativeAPI) Region() forge.Region {
	return a.region
}

// SetRegion sets the Region _AND_ RelativePath of the ModelDerivativeAPI.
//   - If the region is US, the relativePath will be usPath (/modelderivative/v2/designdata)
//   - If the region is EU (== EMEA), the relativePath will be euPath (/modelderivative/v2/regions/eu/designdata)
//
// References:
// - https://aps.autodesk.com/en/docs/model-derivative/v2/reference/http/
func (a *ModelDerivativeAPI) SetRegion(region forge.Region) {
	a.region = region
	if region == forge.US {
		a.relativePath = usPath
	} else if region == forge.EU {
		a.relativePath = euPath
	}
}

// RelativePath of the ModelDerivativeAPI.
// Please note that the relativePath depends on the region.
//
// References:
// - https://aps.autodesk.com/en/docs/model-derivative/v2/reference/http/
func (a *ModelDerivativeAPI) RelativePath() string {
	return a.relativePath
}

// BaseUrl of the ModelDerivativeAPI.
func (a *ModelDerivativeAPI) BaseUrl() string {
	return a.Authenticator.HostPath() + a.relativePath
}

// StartTranslation starts a translation job with the given TranslationParams and XAdsHeaders.
//
// References:
//   - https://aps.autodesk.com/en/docs/model-derivative/v2/reference/http/jobs/job-POST/
func (a *ModelDerivativeAPI) StartTranslation(params TranslationParams, xHeaders XAdsHeaders) (
	result TranslationJob, err error,
) {
	bearer, err := a.Authenticator.GetToken("data:write data:read")
	if err != nil {
		return
	}

	return startTranslation(a.BaseUrl(), params, &xHeaders, bearer.AccessToken)
}

// NewTranslationParams creates a TranslationParams struct with the given urn, outputType, views, and advanced options.
//   - The region will be taken from the ModelDerivativeAPI.
//   - The advanced options can be nil.
//
// Make sure to use the correct combination of views and advanced options for the given outputType.
// There are no checks for this.
func (a *ModelDerivativeAPI) NewTranslationParams(
	urn string, outputType OutputType, views []ViewType, advanced *AdvancedSpec,
) TranslationParams {
	return TranslationParams{
		Input: InputSpec{
			URN: urn,
		},
		Output: OutputSpec{
			Destination: DestSpec{a.region},
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

// DefaultTranslationParams creates a TranslationParams struct with the given urn.
//   - The region will be taken from the ModelDerivativeAPI.
//   - The outputType will be SVF.
//   - The views will be 2D and 3D.
//   - The advanced options will be nil.
func (a *ModelDerivativeAPI) DefaultTranslationParams(urn string) TranslationParams {
	return a.NewTranslationParams(urn, SVF, ViewTypes2DAnd3D(), nil)
}

// UrnFromObjectId creates a Base64 (URL Safe) encoded URN from the given objectID.
//
// OssApi.UploadObject will return an objectID that can be used here.
//
// The URN is required as input for translating the object (CAD file), see:
//   - NewTranslationParams
//   - DefaultTranslationParams
func UrnFromObjectId(objectID string) string {
	return base64.RawStdEncoding.EncodeToString([]byte(objectID))
}

// GetManifest returns information about derivatives that correspond to a specific source file, including derivative URNs and translation statuses.
//
// References:
//   - https://aps.autodesk.com/en/docs/model-derivative/v2/reference/http/manifest/urn-manifest-GET/
func (a *ModelDerivativeAPI) GetManifest(urn string) (result Manifest, err error) {
	bearer, err := a.Authenticator.GetToken("data:read")
	if err != nil {
		return
	}

	return getManifest(a.BaseUrl(), urn, bearer.AccessToken)
}

// GetPropertiesDatabaseUrn returns the URN of the SQLite properties database from the manifest.
// If the database URN is not found, an empty string is returned.
// The database URN is used to download the SQLite properties database using the GetDerivative function.
func (m *Manifest) GetPropertiesDatabaseUrn() string {
	for _, derivative := range m.Derivatives {
		for _, child := range derivative.Children {
			if child.Role == "Autodesk.CloudPlatform.PropertyDatabase" {
				return child.URN
			}
		}
	}
	return ""
}

// GetProgressReport returns the ProgressReport (status and progress) of a translation.
func (m *Manifest) GetProgressReport() ProgressReport {
	return m.ProgressReport
}

// GetProgressReportOfChild returns the ProgressReport of a translation of a given outputType for a specific child,
// identified by its model/view GUID string.
// If the child is not found, an empty ProgressReport is returned.
func (m *Manifest) GetProgressReportOfChild(derivativeOutputType, modelViewGuid string) ProgressReport {
	for _, derivative := range m.Derivatives {
		// strings.EqualFold => ignore casing
		if strings.EqualFold(derivative.OutputType, derivativeOutputType) {
			for _, child := range derivative.Children {
				if child.ModelGUID != nil && *child.ModelGUID == modelViewGuid {
					return child.ProgressReport
				}
			}
		}
	}
	return ProgressReport{}
}

// GetSourceFileName returns the source file name of the translation.
// If the source file name is not found, an empty string is returned.
func (m *Manifest) GetSourceFileName() string {
	// Is this always the name of the first derivative?
	for _, derivative := range m.Derivatives {
		if len(derivative.Name) > 0 {
			return derivative.Name
		}
	}
	return ""
}

// GetDerivative downloads a selected derivative.
// To download the file, you need to specify the file’s URN, which you retrieve from the manifest.
// You can fetch the manifest using the GetManifest function.
//
// References:
//   - https://aps.autodesk.com/en/docs/model-derivative/v2/reference/http/urn-manifest-derivativeUrn-signedcookies-GET/
//
// Example:
//
// urn := "..."
// derivativeUrn := "..."
// writer, err := os.Create("propertiesDb.sqlite")
//
//	if err != nil {
//	  log.Fatal(err)
//	}
//
// defer writer.Close()
// written, err := client.GetDerivative(urn, derivativeUrn, writer)
//
//	if err != nil {
//	  log.Fatal(err)
//	}
//
// log.Printf("Downloaded %d bytes", written)
func (a *ModelDerivativeAPI) GetDerivative(urn, derivativeUrn string, writer io.Writer) (written int64, err error) {
	bearer, err := a.Authenticator.GetToken("data:read")
	if err != nil {
		return
	}

	return getDerivative(a.BaseUrl(), urn, derivativeUrn, bearer.AccessToken, writer)
}

// GetMetadata returns a list of model views (Viewables) in the source design specified by the `urn` URI parameter.
// It also returns the ID that uniquely identifies the model view.
// You can use this ID with other metadata endpoints to obtain information about the objects within model view.
//
//	NOTE: You can retrieve metadata only from an input file that has been translated to SVF or SVF2.
//
// Reference:
//   - https://aps.autodesk.com/en/docs/model-derivative/v2/reference/http/metadata/urn-metadata-GET/
func (a *ModelDerivativeAPI) GetMetadata(urn string) (result MetaData, err error) {
	bearer, err := a.Authenticator.GetToken("data:read")
	if err != nil {
		return
	}

	return getMetadata(a.BaseUrl(), urn, bearer.AccessToken)
}

// GetMasterModelViewGuid returns the GUID of the master view.
// If no master view is found, the GUID of the first 3D view is returned.
// If no 3D view is found, an empty string is returned.
func (m *MetaData) GetMasterModelViewGuid() string {
	// 1st look for the master view
	for _, view := range m.Data.Views {
		if view.IsMasterView {
			return view.Guid
		}
	}
	// else return the first 3d view
	for _, view := range m.Data.Views {
		if view.Role == View3D {
			return view.Guid
		}
	}
	return ""
}

// GetModelViewProperties returns a list of all properties of all objects that are displayed in the model view specified by the modelGuid URI parameter.
//
// Properties are returned as a flat list ordered, by their `objectId`.
// The `objectId` is a non-persistent ID assigned to an object when a design file is translated to the SVF or SVF2 format.
// This means that:
//   - A design file must be translated to SVF or SVF2 before you can retrieve properties.
//   - The `objectId` of an object can change if the design is translated to SVF or SVF2 again. If you require a persistent ID to reference an object, use externalId.
//
// Note: Before you call this endpoint:
//   - Get the modelGuid (unique model view ID) by using the GetModelViewProperties function.
//
// Reference:
//   - https://aps.autodesk.com/en/docs/model-derivative/v2/reference/http/metadata/urn-metadata-guid-properties-GET/
func (a *ModelDerivativeAPI) GetModelViewProperties(urn, modelGuid string, xHeaders XAdsHeaders) (
	jsonData []byte, err error,
) {
	bearer, err := a.Authenticator.GetToken("data:read")
	if err != nil {
		return
	}

	return getModelViewProperties(a.BaseUrl(), urn, modelGuid, bearer.AccessToken, xHeaders)
}

// GetObjectTree returns a hierarchical list of objects (object tree) in the model view specified by the modelGuid URI parameter.
//
// Note: Before you call this endpoint:
//   - Get the modelGuid (unique model view ID) by using the GetModelViewProperties function.
//
// Use forceGet = `true` to retrieve the object tree even if it exceeds the recommended maximum size (20 MB). The default for forceGet is `false`.
//
// Reference:
//   - https://aps.autodesk.com/en/docs/model-derivative/v2/reference/http/metadata/urn-metadata-guid-GET/
func (a *ModelDerivativeAPI) GetObjectTree(urn, modelGuid string, forceGet bool, xHeaders XAdsHeaders) (
	result ObjectTree, err error,
) {
	bearer, err := a.Authenticator.GetToken("data:read")
	if err != nil {
		return
	}

	return getObjectTree(a.BaseUrl(), urn, modelGuid, bearer.AccessToken, forceGet, xHeaders)
}
