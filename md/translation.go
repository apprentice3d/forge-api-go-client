package md

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/woweh/forge-api-go-client"
)

// TranslationParams are used when specifying translation jobs.
// See: https://aps.autodesk.com/en/docs/model-derivative/v2/reference/http/jobs/job-POST/#body-structure
type TranslationParams struct {
	Input  InputSpec  `json:"input"`  // InputSpec is used when specifying the source design
	Output OutputSpec `json:"output"` // OutputSpec is used when specifying the expected format and views (2d and/or 3d)
	// TODO: Add misc option
}

type InputSpec struct {
	// The URN of the source design. This is typically returned as `ObjectId` when you upload the source design to APS.
	// The URN needs to be Base64 (URL Safe) encoded.
	URN string `json:"urn"`
	// Set this to `true` if the source design is compressed as a zip file.
	// The design can consist of a single file or as in the case of Autodesk Inventor, multiple files.
	// If set to `true`, you must specify the rootFilename attribute.
	CompressedURN *bool `json:"compressedUrn,omitempty"`
	// The name of the top-level design file in the compressed file.
	// Mandatory if the compressedUrn is set to true.
	RootFileName *string `json:"rootFileName,omitempty"`
	//   - true - Instructs the server to check for references and download all referenced files.
	// If the design consists of multiple files (as with Autodesk Inventor assemblies) the translation fails if this attribute is not set to true.
	//   - false - (Default) Does not check for any references.
	CheckReferences *bool `json:"checkReferences,omitempty"`
}

// TranslationJob reflects data received upon successful creation of translation job
type TranslationJob struct {
	Result       string `json:"result"`
	URN          string `json:"urn"`
	AcceptedJobs struct {
		Output OutputSpec `json:"output"`
	}
}

// OutputSpec reflects data found upon creation translation job and receiving translation job status
type OutputSpec struct {
	Destination DestSpec     `json:"destination,omitempty"`
	Formats     []FormatSpec `json:"formats"`
}

// DestSpec is used within OutputSpecs and is useful when specifying the region for translation results
type DestSpec struct {
	Region forge.Region `json:"region"` // Region in which to store outputs. Possible values: US, EMEA. By default, it is set to US.
}

// OutputType is the requested output type.
// For a list of supported types, call the [GET formats endpoint].
// Note that Advanced Options are not supported for all output types.
// Make sure you specify the correct options for the requested output type.
// The API has only been tested with the following output types: svf, svf2 and obj
//
// [GET formats endpoint]: https://aps.autodesk.com/en/docs/model-derivative/v2/reference/http/informational/formats-GET/
type OutputType string

const (
	DWG       OutputType = "dwg"
	FBX       OutputType = "fbx"
	IFC       OutputType = "ifc"
	IGES      OutputType = "iges"
	OBJ       OutputType = "obj"
	STEP      OutputType = "step"
	STL       OutputType = "stl"
	SVF       OutputType = "svf"
	SVF2      OutputType = "svf2"
	Thumbnail OutputType = "thumbnail"
)

// ViewType - possible values: 2d, 3d
// Note that some output types have only one possible view.
// Make sure you specify the correct view for the requested output type.
type ViewType string

const (
	View2D ViewType = "2d"
	View3D ViewType = "3d"
)

// ViewTypes2DAnd3D returns a slice of ViewTypes containing View2D and View3D.
func ViewTypes2DAnd3D() []ViewType {
	return []ViewType{View2D, View3D}
}

// ViewType2D returns a slice of ViewTypes containing only View2D.
func ViewType2D() []ViewType {
	return []ViewType{View2D}
}

// ViewType3D returns a slice of ViewTypes containing only View3D.
func ViewType3D() []ViewType {
	return []ViewType{View3D}
}

// FormatSpec is used within OutputSpecs and should be used when specifying the expected format and views (2d or/and 3d)
type FormatSpec struct {
	Type     OutputType    `json:"type"`               // The requested output types.
	Views    []ViewType    `json:"views"`              // An Array of the requested views.
	Advanced *AdvancedSpec `json:"advanced,omitempty"` // A set of special options, which you must specify only if the input file type is IFC, Revit, or Navisworks.
}

// AdvancedSpec is a set of extra translation options.
//   - You *can* specify them if the input file type is IFC, Revit, or Navisworks and the output is SVF/SVF2.
//   - You *must* specify them if the output is OBJ.
type AdvancedSpec struct {
	// ConversionMethod specifies what _IFC_ loader to use during translation (_IFC_ => SVF/SVF2).
	ConversionMethod IfcConversionMethod `json:"conversionMethod,omitempty"`

	// BuildingStoreys specifies how storeys are translated (_IFC_ => SVF/SVF2).
	// NOTE: These options are applicable **only** when conversionMethod is set to modern or v3.
	BuildingStoreys IfcOption `json:"buildingStoreys,omitempty"`

	// Spaces specifies how spaces are translated (_IFC_ => SVF/SVF2).
	// NOTE: These options are applicable **only** when conversionMethod is set to modern or v3.
	Spaces IfcOption `json:"spaces,omitempty"`

	// OpeningElements specifies how openings are translated (_IFC_ => SVF/SVF2).
	// NOTE: These options are applicable **only** when conversionMethod is set to modern or v3.
	OpeningElements IfcOption `json:"openingElements,omitempty"`

	// TwoDViews specifies the format that 2D views must be rendered to (_Revit_ => SVF/SVF2).
	TwoDViews Rvt2dViews `json:"2dviews,omitempty"`

	// ExtractorVersion specifies what version of the Revit translator/extractor to use (_Revit_ => SVF/SVF2).
	ExtractorVersion RvtExtractorVersion `json:"extractorVersion,omitempty"`

	// GenerateMasterViews specifies if master views shall be created (_Revit_ => SVF/SVF2).
	// This attribute defaults to false.
	GenerateMasterViews *bool `json:"generateMasterViews,omitempty"`

	// MaterialMode specifies the materials to apply to the generated SVF/SVF2 derivatives (_Revit_ => SVF/SVF2).
	MaterialMode RvtMaterialMode `json:"materialMode,omitempty"`

	// HiddenObjects specifies whether hidden objects must be extracted or not (_Navisworks_ => SVF/SVF2).
	HiddenObjects *bool `json:"hiddenObjects,omitempty"`

	// BasicMaterialProperties specifies whether basic material properties must be extracted or not (_Navisworks_ => SVF/SVF2).
	BasicMaterialProperties *bool `json:"basicMaterialProperties,omitempty"`

	// AutodeskMaterialProperties specifies how to handle Autodesk material properties (_Navisworks_ => SVF/SVF2).
	AutodeskMaterialProperties *bool `json:"autodeskMaterialProperties,omitempty"`

	// TimeLinerProperties specifies whether timeliner properties must be extracted or not (_Navisworks_ => SVF/SVF2).
	TimeLinerProperties *bool `json:"timelinerProperties,omitempty"`

	// ExportFileStructure specifies if a single or multiple OBJ files shall be generated (SVF/SVF2 => _OBJ_).
	ExportFileStructure ObjExportFileStructure `json:"exportFileStructure,omitempty"`

	/* Unit specifies the unit for translating models (SVF/SVF2 => _OBJ_).
	This causes the values to change. For example, from millimeters (10, 123, 31) to centimeters (1.0, 12.3, 3.1).
	If the source unit or the unit you are translating into is not supported, the values remain unchanged. */
	Unit ObjUnit `json:"unit,omitempty"`

	// ModelGuid specifies the model view ID (guid) required for geometry extraction (SVF/SVF2 => _OBJ_).
	// Currently, only valid for 3d views.
	ModelGuid string `json:"modelGuid,omitempty"`

	// ObjectIds are required for geometry extraction (SVF/SVF2 => _OBJ_). List object ids to be translated.
	//   NOTE: -1 will extract the entire model. Currently, only valid for 3d views.
	ObjectIds *[]int `json:"objectIds,omitempty"`
}

// startTranslation triggers a translation job with the given TranslationParams and xAdsHeaders.XAdsHeaders.
func startTranslation(path string, params TranslationParams, xAdsHeaders *XAdsHeaders, token string) (
	result TranslationJob, err error,
) {

	byteParams, err := json.Marshal(params)
	if err != nil {
		log.Println("Could not marshal the translation parameters")
		return
	}

	req, err := http.NewRequest("POST", path+"/job", bytes.NewBuffer(byteParams))
	if err != nil {
		return
	}

	log.Println("Creating translation job...")
	log.Println("- Base64 encoded design URL: ", params.Input.URN)
	log.Println("- URL: ", req.URL.String())

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	if xAdsHeaders != nil {
		req.Header.Add("x-ads-derivative-format", string(xAdsHeaders.Format))
		req.Header.Add("x-ads-force", strconv.FormatBool(xAdsHeaders.Overwrite))
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusOK {
		content, _ := io.ReadAll(response.Body)
		err = errors.New("[" + strconv.Itoa(response.StatusCode) + "] " + string(content))
		return
	}

	decoder := json.NewDecoder(response.Body)

	err = decoder.Decode(&result)

	log.Println("Translation job created successfully")

	return
}
