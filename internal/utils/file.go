package utils

import (
	"errors"
	"fmt"
	"os"
)

func DeleteFile(filename string) error {
	return os.Remove(filename)
}

func SaveToFile(filename string, arg interface{}) (*os.File, error) {
	var file *os.File
	var err error
	if !FileExists(filename) {
		file, err = os.Create(filename)
		if err != nil {
			return nil, err
		}
		if arg != nil {
			byteKey := []byte(fmt.Sprintf("%v", arg.(interface{})))
			_, err = file.Write(byteKey)
			if err != nil {
				return nil, err
			}
		}
	} else {
		return nil, errors.New(fmt.Sprintf("file %v already exists. Please delete it before retrying", filename))
	}
	return file, nil
}

// FileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
