package md

// This file provides enums (types and consts) for advanced translation options.

/*
// Revit specific options
*/type (
	// MaterialMode is a Revit specific option that specifies the materials to apply to the generated SVF/SVF2 derivatives.
	MaterialMode string
	// Unit is an OBJ specific option for translating models into different units.
	Unit string
	// ConversionMethod is an IFC specific option that specifies what IFC loader to use during translation.
	ConversionMethod string // An
	// Option are IFC specific options that specify how elements (BuildingStoreys, Spaces or OpeningElements) are translated.
	//   NOTE: These options are applicable only when conversionMethod is set to modern or v3.
	Option string
	// ExportFileStructure is a OBJ specific option for creating a single or multiple OBJ files.
	ExportFileStructure string
)

/*
// Revit specific options
*/

const (
	Auto     MaterialMode = "auto"     // (Default) Use the current setting of the default view of the input file.
	Basic    MaterialMode = "basic"    // Use basic materials.
	AutoDesk MaterialMode = "autodesk" // Use Autodesk materials.
)

/*
// OBJ specific options
*/

const (
	Single   ExportFileStructure = "single"   // (default): creates one OBJ file for all the input files (assembly file).
	Multiple ExportFileStructure = "multiple" // creates a separate OBJ file for each object
)

const (
	Meter      Unit = "meter"
	Decimeter  Unit = "decimeter"
	Centimeter Unit = "centimeter"
	Millimeter Unit = "millimeter"
	Micrometer Unit = "micrometer"
	Nanometer  Unit = "nanometer"
	Yard       Unit = "yard"
	Foot       Unit = "foot"
	Inch       Unit = "inch"
	Mil        Unit = "mil"
	MicroInch  Unit = "microinch"
	None       Unit = ""
)

/*
// IFC specific options
*/

const (
	Legacy ConversionMethod = "legacy" // Use the old Navisworks IFC loader
	Modern ConversionMethod = "modern" // Use the revit IFC loader (recommended over the legacy option).
	V3     ConversionMethod = "v3"     // Use the newer revit IFC loader (recommended over the older modern option)
)

const (
	Hide Option = "hide" // (default) elements are extracted but not visible by default.
	Show Option = "show" // elements are extracted and are visible by default.
	Skip Option = "skip" // elements are not translated.
)
