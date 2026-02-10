package utils

import (
	"errors"
	"fmt"
	"regexp"
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

const (
	Seconds = "s"
	Minutes = "m"
	Hours   = "h"
	Days    = "D"
	Months  = "M"
	Years   = "Y"
)

func ConvertTime(timeToConvert, unitToConvertTo string) (int, error) {
	reg, err := regexp.Compile(`[0-9]+[smhDMY]`)
	if err != nil {
		return 0, fmt.Errorf("failed to compile time format regex: %v", err)
	}

	timeToConvert = strings.ReplaceAll(timeToConvert, " ", "")

	// Validate that the entire input matches the expected format
	validationReg, err := regexp.Compile(`^([0-9]+[smhDMY])+$`)
	if err != nil {
		return 0, fmt.Errorf("failed to compile validation regex: %v", err)
	}

	if !validationReg.MatchString(timeToConvert) {
		return 0, fmt.Errorf("invalid time format: '%s'. Accepted formats: Y, M, D, h, m, s (e.g., 1Y, 30D, 1m30s)", timeToConvert)
	}

	splitTime := reg.FindAllString(timeToConvert, -1)
	totalTime := 0
	for _, t := range splitTime {
		res := 0

		for _, unit := range []string{Seconds, Minutes, Hours, Days, Months, Years} {
			if !strings.HasSuffix(t, unit) {
				continue
			}

			switch unitToConvertTo {
			case Seconds:
				res, err = convertToSeconds(t, unit)
				if err != nil {
					return 0, err
				}
			case Minutes:
				res, err = convertToMinutes(t, unit)
				if err != nil {
					return 0, err
				}
			case Hours:
				res, err = convertToHours(t, unit)
				if err != nil {
					return 0, err
				}
			case Days:
				res, err = convertToDays(t, unit)
				if err != nil {
					return 0, err
				}
			case Months:
				res, err = convertToMonths(t, unit)
				if err != nil {
					return 0, err
				}
			case Years:
				res, err = convertToYears(t, unit)
				if err != nil {
					return 0, err
				}
			default:
				return 0, errors.New("error converting to the specified unit")
			}
		}

		totalTime += res
	}

	return totalTime, nil
}

func convertToSeconds(time, unit string) (int, error) {
	switch unit {
	case Seconds:
		s := strings.ReplaceAll(time, unit, "")
		return strconv.Atoi(s)
	case Minutes:
		s := strings.ReplaceAll(time, unit, "")
		m, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return m * 60, nil
	case Hours:
		s := strings.ReplaceAll(time, unit, "")
		h, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return h * 60 * 60, nil
	case Days:
		s := strings.ReplaceAll(time, unit, "")
		d, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return d * 60 * 60 * 24, nil
	case Months:
		s := strings.ReplaceAll(time, unit, "")
		m, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return m * 60 * 60 * 24 * 30, nil
	case Years:
		s := strings.ReplaceAll(time, unit, "")
		y, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return y * 60 * 60 * 24 * 365, nil
	}
	return 0, errors.New("error converting to seconds, no suffix: s, m, h, D, M, Y matched")
}

func convertToMinutes(time, unit string) (int, error) {
	switch unit {
	case Seconds:
		s := strings.ReplaceAll(time, unit, "")
		sec, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return sec / 60, nil
	case Minutes:
		s := strings.ReplaceAll(time, unit, "")
		return strconv.Atoi(s)
	case Hours:
		s := strings.ReplaceAll(time, unit, "")
		h, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return h * 60, nil
	case Days:
		s := strings.ReplaceAll(time, unit, "")
		d, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return d * 60 * 24, nil
	case Months:
		s := strings.ReplaceAll(time, unit, "")
		m, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return m * 60 * 24 * 30, nil
	case Years:
		s := strings.ReplaceAll(time, unit, "")
		y, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return y * 60 * 24 * 365, nil
	}
	return 0, errors.New("error converting to minutes, no suffix: s, m, h, D, M, Y matched")
}

func convertToHours(time, unit string) (int, error) {
	switch unit {
	case Seconds:
		s := strings.ReplaceAll(time, unit, "")
		sec, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return sec / (60 * 60), nil
	case Minutes:
		s := strings.ReplaceAll(time, unit, "")
		m, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return m / 60, nil
	case Hours:
		s := strings.ReplaceAll(time, unit, "")
		return strconv.Atoi(s)
	case Days:
		s := strings.ReplaceAll(time, unit, "")
		d, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return d * 24, nil
	case Months:
		s := strings.ReplaceAll(time, unit, "")
		m, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return m * 24 * 30, nil
	case Years:
		s := strings.ReplaceAll(time, unit, "")
		y, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return y * 24 * 365, nil
	}
	return 0, errors.New("error converting to hours, no suffix: s, m, h, D, M, Y matched")
}

func convertToDays(time, unit string) (int, error) {
	switch unit {
	case Seconds:
		s := strings.ReplaceAll(time, unit, "")
		sec, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return sec / (60 * 60 * 24), nil
	case Minutes:
		s := strings.ReplaceAll(time, unit, "")
		m, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return m / (60 * 24), nil
	case Hours:
		s := strings.ReplaceAll(time, unit, "")
		h, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return h / 24, nil
	case Days:
		s := strings.ReplaceAll(time, unit, "")
		return strconv.Atoi(s)
	case Months:
		s := strings.ReplaceAll(time, unit, "")
		m, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return m * 30, nil
	case Years:
		s := strings.ReplaceAll(time, unit, "")
		y, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return y * 365, nil
	}
	return 0, errors.New("error converting to days, no suffix: s, m, h, D, M, Y matched")
}

func convertToMonths(time, unit string) (int, error) {
	switch unit {
	case Seconds:
		s := strings.ReplaceAll(time, unit, "")
		sec, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return sec / (60 * 60 * 24 * 30), nil
	case Minutes:
		s := strings.ReplaceAll(time, unit, "")
		m, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return m / (60 * 24 * 30), nil
	case Hours:
		s := strings.ReplaceAll(time, unit, "")
		h, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return h / (24 * 30), nil
	case Days:
		s := strings.ReplaceAll(time, unit, "")
		d, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return d / 30, nil
	case Months:
		s := strings.ReplaceAll(time, unit, "")
		return strconv.Atoi(s)
	case Years:
		s := strings.ReplaceAll(time, unit, "")
		y, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return y * 12, nil
	}
	return 0, errors.New("error converting to months, no suffix: s, m, h, D, M, Y matched")
}

func convertToYears(time, unit string) (int, error) {
	switch unit {
	case Seconds:
		s := strings.ReplaceAll(time, unit, "")
		sec, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return sec / (60 * 60 * 24 * 365), nil
	case Minutes:
		s := strings.ReplaceAll(time, unit, "")
		m, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return m / (60 * 24 * 365), nil
	case Hours:
		s := strings.ReplaceAll(time, unit, "")
		h, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return h / (24 * 365), nil
	case Days:
		s := strings.ReplaceAll(time, unit, "")
		d, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return d / 365, nil
	case Months:
		s := strings.ReplaceAll(time, unit, "")
		m, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return m / 12, nil
	case Years:
		s := strings.ReplaceAll(time, unit, "")
		return strconv.Atoi(s)
	}
	return 0, errors.New("error converting to years, no suffix: s, m, h, D, M, Y matched")
}
