package clierror

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	"github.com/spf13/viper"
)

var (
	unknownTypeFormatErr = "unknown type format %s. Hint: use --output json|text"
	ErrAction            = func() {
		os.Exit(1)
	}
)

type CliError struct {
	Err    error  `json:"Error,omitempty"`
	Detail string `json:"Detail,omitempty"`
}

// CheckError Standard error checking
func CheckError(err error, outErr io.Writer) {
	if err == nil {
		return
	}
	cliErr := CliError{
		Err:    err,
		Detail: err.Error(),
	}

	switch viper.GetString(config.ArgOutput) {
	case printer.TypeJSON.String():
		printer.WriteJSON(&cliErr, outErr)
	case printer.TypeText.String():
		errorConfirm(outErr, cliErr.Err.Error())
	default:
		err := errors.New(fmt.Sprintf(unknownTypeFormatErr, viper.GetString(config.ArgOutput)))
		errorConfirm(outErr, err.Error())
	}

	ErrAction()
}

func errorConfirm(writer io.Writer, msg string, args ...interface{}) {
	if strings.HasSuffix(msg, "\n") {
		fmt.Fprintf(writer, "Error: %s", fmt.Sprintf(msg, args...))
	} else {
		fmt.Fprintf(writer, "Error: %s\n", fmt.Sprintf(msg, args...))
	}
}
