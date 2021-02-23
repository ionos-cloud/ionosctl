package clierror

import (
	"fmt"
)

type RequiredFlagErr struct {
	FlagName string
}

var _ error = &RequiredFlagErr{}

func NewRequiredFlagErr(flagName string) *RequiredFlagErr {
	return &RequiredFlagErr{
		FlagName: flagName,
	}
}

func (e *RequiredFlagErr) Error() string {
	return fmt.Sprintf("%s required flag is not set", e.FlagName)
}
