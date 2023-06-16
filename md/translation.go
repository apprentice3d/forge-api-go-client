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

// TranslationResult reflects data received upon successful creation of translation job
type TranslationResult struct {
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

// OutputView - possible values: 2d, 3d
// Note that some output types have only one possible view.
// Make sure you specify the correct view for the requested output type.
type OutputView string

const (
	View2D OutputView = "2d"
	View3D OutputView = "3d"
)

// FormatSpec is used within OutputSpecs and should be used when specifying the expected format and views (2d or/and 3d)
type FormatSpec struct {
	Type     OutputType    `json:"type"`               // The requested output types.
	Views    []OutputView  `json:"views"`              // An Array of the requested views.
	Advanced *AdvancedSpec `json:"advanced,omitempty"` // A set of special options, which you must specify only if the input file type is IFC, Revit, or Navisworks.
}

// AdvancedSpec is a set of extra translation options.
// You *can* specify them if the input file type is IFC, Revit, or Navisworks and the output is SVF/SVF2.
// You *must* specify them if the output is OBJ.
type AdvancedSpec struct {
	// SVF/SVF2 option to be specified when the input file type is _IFC_. Specifies what _IFC_ loader to use during translation.
	ConversionMethod ConversionMethod `json:"conversionMethod,omitempty"`
	/* SVF/SVF2 option to be specified when the input file type is _IFC_. Specifies how storeys are translated.
	   NOTE: These options are applicable only when conversionMethod is set to modern or v3. */
	BuildingStoreys Option `json:"buildingStoreys,omitempty"`
	/* SVF/SVF2 option to be specified when the input file type is _IFC_. Specifies how spaces are translated.
	   NOTE: These options are applicable only when conversionMethod is set to modern or v3. */
	Spaces Option `json:"spaces,omitempty"`
	/* SVF/SVF2 option to be specified when the input file type is _IFC_. Specifies how openings are translated.
	   NOTE: These options are applicable only when conversionMethod is set to modern or v3. */
	OpeningElements Option `json:"openingElements,omitempty"`
	/* SVF/SVF2 option to be specified when the input file type is _Revit_.
	Generates master views when translating from the _Revit_ input format to SVF/SVF2.
	This option is ignored for all other input formats. This attribute defaults to false. */
	GenerateMasterViews *bool `json:"generateMasterViews,omitempty"`
	/* SVF/SVF2 option to be specified when the input file type is _Revit_.
	Specifies the materials to apply to the generated SVF/SVF2 derivatives. */
	MaterialMode MaterialMode `json:"materialMode,omitempty"`
	// SVF/SVF2 option to be specified when the input file type is _Navisworks_.
	HiddenObjects *bool `json:"hiddenObjects,omitempty"`
	// SVF/SVF2 option to be specified when the input file type is _Navisworks_.
	BasicMaterialProperties *bool `json:"basicMaterialProperties,omitempty"`
	// SVF/SVF2 option to be specified when the input file type is _Navisworks_.
	AutodeskMaterialProperties *bool `json:"autodeskMaterialProperties,omitempty"`
	// SVF/SVF2 option to be specified when the input file type is _Navisworks_.
	TimeLinerProperties *bool `json:"timelinerProperties,omitempty"`
	// OBJ option for creating a single or multiple OBJ files.
	ExportFileStructure ExportFileStructure `json:"exportFileStructure,omitempty"`
	/* OBJ option for translating models into different units.
	This causes the values to change. For example, from millimeters (10, 123, 31) to centimeters (1.0, 12.3, 3.1).
	If the source unit or the unit you are translating into is not supported, the values remain unchanged. */
	Unit Unit `json:"unit,omitempty"`
	/* OBJ option required for geometry extraction.
	   The model view ID (guid). Currently, only valid for 3d views. */
	ModelGuid string `json:"modelGuid,omitempty"`
	/* OBJ option required for geometry extraction. List object ids to be translated.
	   -1 will extract the entire model. Currently, only valid for 3d views. */
	ObjectIds *[]int `json:"objectIds,omitempty"`
}

// translate triggers a translation job with the given TranslationParams and xAdsHeaders.XAdsHeaders.
func translate(path string, params TranslationParams, xAdsHeaders *XAdsHeaders, token string) (
	result TranslationResult, err error,
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

	return
}
