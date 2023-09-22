package clierror

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"

	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	"github.com/spf13/viper"
)

var (
	unknownTypeFormatErr = "unknown type format %s. Hint: use --output json|text"
)

type CliError struct {
	Err    error  `json:"Error,omitempty"`
	Detail string `json:"Detail,omitempty"`
}

// CheckErrorAndDie Standard error checking
//
// DEPRECATED: Use die.Die instead if you want to kill the execution, or return the errors in a go-ish fashion!
func CheckErrorAndDie(err error, outErr io.Writer) {
	if err == nil {
		return
	}
	cliErr := CliError{
		Err:    err,
		Detail: err.Error(),
	}

	switch viper.GetString(constants.ArgOutput) {
	case printer.TypeJSON.String():
		printer.WriteJSON(&cliErr, outErr)
	case printer.TypeText.String():
		errorConfirm(outErr, cliErr.Err.Error())
	default:
		err := errors.New(fmt.Sprintf(unknownTypeFormatErr, viper.GetString(constants.ArgOutput)))
		errorConfirm(outErr, err.Error())
	}

	os.Exit(1)
}

func errorConfirm(writer io.Writer, msg string, args ...interface{}) {
	if strings.HasSuffix(msg, "\n") {
		fmt.Fprintf(writer, "Error: %s", fmt.Sprintf(msg, args...))
	} else {
		fmt.Fprintf(writer, "Error: %s\n", fmt.Sprintf(msg, args...))
	}
}
