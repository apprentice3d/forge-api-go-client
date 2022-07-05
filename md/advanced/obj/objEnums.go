// Package obj provides types and consts for OBJ specific advanced translation options.
package obj

type ExportFileStructure string // An option for creating a single or multiple OBJ files.

const (
	Single   ExportFileStructure = "single"   // (default): creates one OBJ file for all the input files (assembly file).
	Multiple ExportFileStructure = "multiple" // creates a separate OBJ file for each object
)

type Unit string // An option for translating models into different units.

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
)
