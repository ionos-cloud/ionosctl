package convbytes

import (
	"regexp"
	"strconv"
	"strings"
)

const (
	B  int64 = 1
	KB       = B * 1024
	MB       = KB * 1024
	GB       = MB * 1024
	TB       = GB * 1024
	PB       = TB * 1024
)

// ToBytes converts a given value in a source unit to bytes.
func ToBytes(value int64, unit int64) int64 {
	return value * unit
}

// FromBytes converts a given byte value to a target unit.
func FromBytes(bytes int64, unit int64) int64 {
	if unit == 0 {
		return 0
	}
	return bytes / unit
}

// Convert converts a value from a source unit to a target unit.
func Convert(value int64, fromUnit int64, toUnit int64) int64 {
	bytes := ToBytes(value, fromUnit)
	return FromBytes(bytes, toUnit)
}

// StrToBytesOk converts your bytes string to bytes, for use with other convbytes pkg funcs
func StrToBytesOk(s string) (int64, bool) {
	s = strings.TrimSpace(s)

	// If the string is just a number, return it as is
	if num, err := strconv.ParseInt(s, 10, 64); err == nil {
		return num, true
	}

	// Match numbers followed by optional spaces and a unit
	r := regexp.MustCompile(`^(\d+)\s*([a-zA-Z]*)$`)
	matches := r.FindStringSubmatch(s)
	if matches == nil {
		return 0, false
	}

	value, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return 0, false
	}

	if matches[2] == "" {
		return value, true
	}

	unit, ok := map[string]int64{
		"B":  B,
		"KB": KB,
		"MB": MB,
		"GB": GB,
		"TB": TB,
		"PB": PB,
	}[strings.ToUpper(matches[2])]

	if !ok {
		return 0, false
	}

	return value * unit, true
}

func StrToBytes(s string) int64 {
	f, _ := StrToBytesOk(s)
	return f
}

// StrToUnitOk converts a size string to a target unit and returns the converted value and a boolean indicating if the conversion was successful.
func StrToUnitOk(s string, targetUnit int64) (int64, bool) {
	s = strings.TrimSpace(s)

	// If the string is just a number, assume it's in the target format
	if num, err := strconv.ParseInt(s, 10, 64); err == nil {
		return num, true
	}

	bytes, ok := StrToBytesOk(s)
	if !ok {
		return 0, false
	}
	return FromBytes(bytes, targetUnit), true
}

// StrToUnit converts a size string to a target unit and returns the converted value.
func StrToUnit(s string, targetUnit int64) int64 {
	val, _ := StrToUnitOk(s, targetUnit)
	return val
}
