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
	return bytes / unit
}

// Convert converts a value from a source unit to a target unit.
func Convert(value int64, fromUnit int64, toUnit int64) int64 {
	bytes := ToBytes(value, fromUnit)
	return FromBytes(bytes, toUnit)
}

// FromStringOk converts your bytes string to bytes, for use with other convbytes pkg funcs
func FromStringOk(s string) (int64, bool) {
	s = strings.TrimSpace(s)
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

func FromString(s string) int64 {
	f, _ := FromStringOk(s)
	return f
}
