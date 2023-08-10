package utils

import (
	"errors"
	"strconv"
	"strings"
)

const (
	MegaBytes  = "MB"
	GigaBytes  = "GB"
	TerraBytes = "TB"
	PetaBytes  = "PB"
)

// ConvertSize converts the specified size to the unit specified
// Right now, it has support for MB, GB
func ConvertSize(sizeToConvert, unitToConvertTo string) (int, error) {
	for _, unit := range []string{MegaBytes, GigaBytes, PetaBytes, TerraBytes} {
		if !strings.HasSuffix(sizeToConvert, unit) {
			continue
		}
		sizeToConvert = strings.ReplaceAll(sizeToConvert, " ", "")
		switch unitToConvertTo {
		case MegaBytes:
			return convertToMB(sizeToConvert, unit)
		case GigaBytes:
			return convertToGB(sizeToConvert, unit)
		default:
			return 0, errors.New("error converting to the specified unit")
		}
	}
	return strconv.Atoi(sizeToConvert)
}

func ConvertToMB(size, unit string) (int, error) {
	return convertToMB(size, unit)
}

func convertToMB(size, unit string) (int, error) {
	switch unit {
	case MegaBytes:
		s := strings.ReplaceAll(size, unit, "")
		return strconv.Atoi(s)
	case GigaBytes:
		s := strings.ReplaceAll(size, unit, "")
		gb, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return gb * 1024, nil
	case TerraBytes:
		s := strings.ReplaceAll(size, unit, "")
		tb, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return tb * 1024 * 1024, nil
	case PetaBytes:
		s := strings.ReplaceAll(size, unit, "")
		pb, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return pb * 1024 * 1024 * 1024, nil
	default:
		return 0, errors.New("error converting in MB, no suffix: MB, GB, TB, PB matched")
	}
}

func ConvertToGB(size, unit string) (int, error) {
	return convertToGB(size, unit)
}

func convertToGB(size, unit string) (int, error) {
	switch unit {
	case MegaBytes:
		s := strings.ReplaceAll(size, unit, "")
		mb, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return mb / 1024, nil
	case GigaBytes:
		s := strings.ReplaceAll(size, unit, "")
		return strconv.Atoi(s)
	case TerraBytes:
		s := strings.ReplaceAll(size, unit, "")
		tb, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return tb * 1024, nil
	case PetaBytes:
		s := strings.ReplaceAll(size, unit, "")
		pb, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return pb * 1024 * 1024, nil
	default:
		return 0, errors.New("error converting in GB, no suffix: MB, GB, TB, PB matched")
	}
}
