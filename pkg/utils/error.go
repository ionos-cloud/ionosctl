package utils

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/spf13/viper"
)

var (
	ErrAction = func() {
		os.Exit(1)
	}
)

type CliError struct {
	Err    error  `json:"Error,omitempty"`
	Detail string `json:"Detail,omitempty"`
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
		writeJSON(&cliErr, outErr)
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
