package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/spf13/viper"
	"io"
	"os"
)

var (
	ErrAction = func() {
		os.Exit(1)
	}
)

type CliError struct {
	Err    error  `json:"error,omitempty"`
	Detail string `json:"detail,omitempty"`
}

// Standard error checking
func CheckError(err error, outErr io.Writer) {
	if err == nil {
		return
	}

	cliErr := CliError{
		Err:    err,
		Detail: err.Error(),
	}

	switch viper.GetString(config.ArgOutput) {
	case PrinterTypeJSON.String():
		b, _ := json.Marshal(&cliErr)
		fmt.Fprintf(outErr, string(b))
	case PrinterTypeText.String():
		errorConfirm(outErr, cliErr.Err.Error())
	default:
		err := errors.New(fmt.Sprintf(unknownTypeFormatErr, viper.GetString(config.ArgOutput)))
		errorConfirm(outErr, err.Error())
	}

	ErrAction()
}

func errorConfirm(writer io.Writer, msg string, args ...interface{}) {
	colorWarn := color.RedString("Error")
	fmt.Fprintf(writer, "\u2716 %s: %s\n", colorWarn, fmt.Sprintf(msg, args...))
}

/*
	Common errors
*/

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
	return fmt.Sprintf("%s required flag is not set\n", e.FlagName)
}
