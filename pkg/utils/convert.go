package utils

// TODO: This file should not exist anymore. Replace me with convbytes pkg calls

import (
	"errors"

	"github.com/ionos-cloud/ionosctl/v6/pkg/convbytes"
)

const (
	MegaBytes  = "MB"
	GigaBytes  = "GB"
	TerraBytes = "TB"
	PetaBytes  = "PB"
)

// ConvertSize converts the specified size to the unit specified
// Right now, it has support for MB, GB
//
// DEPRECATED: This func now simply calls `convbytes` pkg
func ConvertSize(sizeToConvert, unitToConvertTo string) (int, error) {
	// TODO: Replace all calls with calls to convbytes
	val, ok := convbytes.FromStringOk(sizeToConvert)
	if !ok {
		return 0, errors.New("invalid size string format")
	}

	switch unitToConvertTo {
	case MegaBytes:
		return int(convbytes.FromBytes(val, convbytes.MB)), nil
	case GigaBytes:
		return int(convbytes.FromBytes(val, convbytes.GB)), nil
	default:
		return 0, errors.New("error converting to the specified unit")
	}
}

// DEPRECATED: This func now simply calls `convbytes` pkg
func ConvertToMB(size, unit string) (int, error) {
	// TODO: Replace all calls with calls to convbytes

	val, ok := convbytes.FromStringOk(size)
	if !ok {
		return 0, errors.New("invalid size string format")
	}
	return int(convbytes.FromBytes(val, convbytes.MB)), nil
}

// DEPRECATED: This func now simply calls `convbytes` pkg
func ConvertToGB(size, unit string) (int, error) {
	// TODO: Replace all calls with calls to convbytes

	val, ok := convbytes.FromStringOk(size)
	if !ok {
		return 0, errors.New("invalid size string format")
	}
	return int(convbytes.FromBytes(val, convbytes.GB)), nil
}
