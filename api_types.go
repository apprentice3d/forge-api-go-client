package forge

// This file contains types and constants for the entire APS (aka Forge) API.

import "strings"

// HostName / domain of the Autodesk Forge API.
const HostName = "https://developer.api.autodesk.com"

// Region is the region where the data resides.
type Region string

const (
	US   Region = "us"   // US Region - us in lowercase!
	EU   Region = "emea" // EU (== EMEA) Region - emea in lowercase!
	EMEA Region = "emea" // EMEA (== EU) Region - emea in lowercase!
)

// IsUS returns true if the region is US.
func (r Region) IsUS() bool {
	// case insensitive comparison!
	return strings.EqualFold(string(r), string(US))
}

// IsEMEA returns true if the region is EMEA (== EU).
func (r Region) IsEMEA() bool {
	// case insensitive comparison!
	return strings.EqualFold(string(r), string(EMEA))
}

// IsEU returns true if the region is EU (== EMEA).
func (r Region) IsEU() bool {
	// case insensitive comparison!
	return strings.EqualFold(string(r), string(EU))
}
