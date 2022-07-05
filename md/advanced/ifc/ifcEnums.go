// Package ifc provides types and consts for IFC specific advanced translation options.
package ifc

type ConversionMethod string // An option to be specified when the input file type is IFC. Specifies what IFC loader to use during translation.

const (
	Legacy ConversionMethod = "legacy" // Use the old Navisworks IFC loader
	Modern ConversionMethod = "modern" // Use the revit IFC loader (recommended over the legacy option).
	V3     ConversionMethod = "v3"     // Use the newer revit IFC loader (recommended over the older modern option)
)

type Option string /* IFC specific option that specifies how elements (BuildingStoreys, Spaces or OpeningElements) are translated.
NOTE: These options are applicable only when conversionMethod is set to modern or v3. */

const (
	Hide Option = "hide" // (default) elements are extracted but not visible by default.
	Show Option = "show" // elements are extracted and are visible by default.
	Skip Option = "skip" // elements are not translated.
)
