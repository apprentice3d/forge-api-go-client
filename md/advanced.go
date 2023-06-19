package md

// This file provides "enums" (= types and consts) and factory functions for advanced translation options.

type (
	// RvtMaterialMode is a Revit specific option that specifies the materials to apply to the generated SVF/SVF2 derivatives.
	RvtMaterialMode string

	// Rvt2dViews defines the format that 2D views must be rendered to.
	// Available options are:
	//   - legacy - (Default) Render to a model derivative specific file format.
	//   - pdf - Render to the PDF file format. This option applies only to Revit 2022 files and newer.
	Rvt2dViews string

	// RvtExtractorVersion specifies what version of the Revit translator/extractor to use.
	//   NOTE: If this attribute is not specified, the system uses the current official release version.
	// Possible values:
	//   - next - Makes the translation job use the newest available version of the translator/extractor.
	//  	This option is meant to be used by beta testers who wish to test a pre-release version of the translator.
	// 		If no pre-release version is available, this option makes the translation job use the current official release version.
	//   - previous - Makes the translation job use the previous official release version of the translator/extractor.
	//  	This option is meant to be used as a workaround in case of regressions caused by a new official release of the translator/extractor.
	RvtExtractorVersion string

	// ObjUnit is an OBJ specific option for translating models into different units.
	ObjUnit string

	// ObjExportFileStructure is a OBJ specific option for creating a single or multiple OBJ files.
	ObjExportFileStructure string

	// IfcConversionMethod is an IFC specific option that specifies what IFC loader to use during translation.
	IfcConversionMethod string // An

	// IfcOption are IFC specific options that specify how elements (BuildingStoreys, Spaces or OpeningElements) are translated.
	//   NOTE: These options are applicable only when conversionMethod is set to modern or v3.
	IfcOption string
)

/*
Revit specific options
*/

const (
	RvtAuto     RvtMaterialMode     = "auto"     // (Default) Use the current setting of the default view of the input file.
	RvtBasic    RvtMaterialMode     = "basic"    // Use basic materials.
	RvtAutoDesk RvtMaterialMode     = "autodesk" // Use Autodesk materials.
	RvtLegacy   Rvt2dViews          = "legacy"   // (Default) Render to a model derivative specific file format.
	RvtPdf      Rvt2dViews          = "pdf"      // Render to the PDF file format. This option applies only to Revit 2022 files and newer.
	RvtNext     RvtExtractorVersion = "next"     // Makes the translation job use the newest available version of the translator/extractor.
	RvtPrevious RvtExtractorVersion = "previous" // Makes the translation job use the previous official release version of the translator/extractor.
)

/*
OBJ specific options
*/

const (
	ObjSingle   ObjExportFileStructure = "single"   // (default): creates one OBJ file for all the input files (assembly file).
	ObjMultiple ObjExportFileStructure = "multiple" // creates a separate OBJ file for each object
)

const (
	ObjMeter      ObjUnit = "meter"
	ObjDecimeter  ObjUnit = "decimeter"
	ObjCentimeter ObjUnit = "centimeter"
	ObjMillimeter ObjUnit = "millimeter"
	ObjMicrometer ObjUnit = "micrometer"
	ObjNanometer  ObjUnit = "nanometer"
	ObjYard       ObjUnit = "yard"
	ObjFoot       ObjUnit = "foot"
	ObjInch       ObjUnit = "inch"
	ObjMil        ObjUnit = "mil"
	ObjMicroInch  ObjUnit = "microinch"
	ObjNone       ObjUnit = ""
)

/*
IFC specific options
*/

const (
	IfcLegacy IfcConversionMethod = "legacy" // Use the old Navisworks IFC loader
	IfcModern IfcConversionMethod = "modern" // Use the revit IFC loader (recommended over the legacy option).
	IfcV3     IfcConversionMethod = "v3"     // Use the newer revit IFC loader (recommended over the older modern option)
)

const (
	IfcHide IfcOption = "hide" // (default) elements are extracted but not visible by default.
	IfcShow IfcOption = "show" // elements are extracted and are visible by default.
	IfcSkip IfcOption = "skip" // elements are not translated.
)

// IfcAdvancedSpec returns an IFC specific AdvancedSpec.
//
// NOTE:
//   - The storeys, spaces, and openings options are applicable only when conversionMethod is set to `modern` or `v3`.
//   - If the conversionMethod is set to `legacy`, the storeys, spaces, and openings options are ignored.
//   - Use empty strings for the IfcOptions to use the default values.
//
// Reference:
//   - https://aps.autodesk.com/en/docs/model-derivative/v2/reference/http/jobs/job-POST/#body-structure
//     -> Case 4 - Input file type is IFC:
func IfcAdvancedSpec(conversionMethod IfcConversionMethod, storeys, spaces, openings IfcOption) *AdvancedSpec {
	if conversionMethod == IfcLegacy {
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
//   - Use empty strings for materialMode, twoDView and version to use the default values.
//
// Reference:
//   - https://aps.autodesk.com/en/docs/model-derivative/v2/reference/http/jobs/job-POST/#body-structure
//     -> Case 6 - Input file type is RVT:
func RevitAdvancedSpec(
	generateMasterViews bool, materialMode RvtMaterialMode, twoDView Rvt2dViews, version RvtExtractorVersion,
) *AdvancedSpec {
	return &AdvancedSpec{
		GenerateMasterViews: &generateMasterViews,
		MaterialMode:        materialMode,
		TwoDViews:           twoDView,
		ExtractorVersion:    version,
	}
}

// NavisworksAdvancedSpec returns a Navisworks specific AdvancedSpec.
//   - The Autodesk defaults for the 4 options are all `false`.
//
// Reference:
//   - https://aps.autodesk.com/en/docs/model-derivative/v2/reference/http/jobs/job-POST/#body-structure
//     -> Case 5 - Input file type is NWD:
func NavisworksAdvancedSpec(hiddenObjects, basicMaterialProperties, autodeskMaterialProperties, timeLinerProperties bool) *AdvancedSpec {
	return &AdvancedSpec{
		HiddenObjects:              &hiddenObjects,
		BasicMaterialProperties:    &basicMaterialProperties,
		AutodeskMaterialProperties: &autodeskMaterialProperties,
		TimeLinerProperties:        &timeLinerProperties,
	}
}

// ObjAdvancedSpec returns a OBJ specific AdvancedSpec.
//
// Reference:
//   - https://aps.autodesk.com/en/docs/model-derivative/v2/reference/http/jobs/job-POST/#body-structure
//     -> Attributes that Apply to OBJ Outputs
func ObjAdvancedSpec(
	exportFileStructure ObjExportFileStructure, unit ObjUnit, modelGuid string, objectIds *[]int,
) *AdvancedSpec {
	return &AdvancedSpec{
		ExportFileStructure: exportFileStructure,
		Unit:                unit,
		ModelGuid:           modelGuid,
		ObjectIds:           objectIds,
	}
}
