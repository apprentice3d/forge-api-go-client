package forge

import "strings"

// Region is the region where the data resides.
type Region string

const (
	US   Region = "us"   // US Region - us in lowercase!
	EMEA Region = "emea" // EMEA (=> EU) Region - emea in lowercase!
)

// IsUS returns true if the region is US.
func (r Region) IsUS() bool {
	// case insensitive comparison!
	return strings.EqualFold(string(r), string(US))
}

// IsEMEA returns true if the region is EMEA (=> EU).
func (r Region) IsEMEA() bool {
	// case insensitive comparison!
	return strings.EqualFold(string(r), string(EMEA))
}
