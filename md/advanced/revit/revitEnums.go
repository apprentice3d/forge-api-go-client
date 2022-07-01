package revit

type MaterialMode string // An option to be specified when the input file type is Revit. Specifies the materials to apply to the generated SVF/SVF2 derivatives.

const (
	Auto     MaterialMode = "auto"     // (Default) Use the current setting of the default view of the input file.
	Basic    MaterialMode = "basic"    // Use basic materials.
	AutoDesk MaterialMode = "autodesk" // Use Autodesk materials.
)
