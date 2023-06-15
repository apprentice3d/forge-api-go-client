package md

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
)

// TranslationParams is used when specifying the translation jobs
// See: https://forge.autodesk.com/en/docs/model-derivative/v2/reference/http/job-POST/
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
	Region string `json:"region"` // Region in which to store outputs. Possible values: US, EMEA. By default, it is set to US.
}

// FormatSpec is used within OutputSpecs and should be used when specifying the expected format and views (2d or/and 3d)
type FormatSpec struct {
	Type     string        `json:"type"`               // The requested output types.
	Views    []string      `json:"views"`              // An Array of the requested views.
	Advanced *AdvancedSpec `json:"advanced,omitempty"` // A set of special options, which you must specify only if the input file type is IFC, Revit, or Navisworks.
}

// AdvancedSpec is a set of extra translation options.
// You can specify them if the input file type is IFC, Revit, or Navisworks and the output is SVF/SVF2.
// You must specify them if the output is OBJ.
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
func translate(path string, params TranslationParams, xAdsHeaders *XAdsHeaders, token string) (result TranslationResult, err error) {

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
