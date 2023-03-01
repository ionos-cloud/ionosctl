package utils

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	"github.com/spf13/viper"
)

var getUserInput = func(reader io.Reader, writer printer.PrintService, message string) (string, error) {
	return readUserInput(reader, writer, message)
}

func readUserInput(reader io.Reader, writer printer.PrintService, message string) (string, error) {
	in := bufio.NewReader(reader)
	err := writer.Print("Warning: Are you sure you want to " + message + " (y/N) ? ")
	if err != nil {
		return "", err
	}
	answer, err := in.ReadString('\n')
	if err != nil {
		return "", err
	}

	answer = strings.TrimRight(answer, "\r\n")
	return strings.ToLower(answer), nil
}

// AskForConfirm parses and verifies user input for confirmation.
// Checks "--force" or "--quiet" WITHOUT core.GetFlagName. WARNING: if your Force/Quiet flags are bound to viper, they are most likely going to be ignored!
// DEPRECATED: Use confirm.Ask instead
func AskForConfirm(reader io.Reader, writer printer.PrintService, message string) error {
	if viper.GetBool(constants.ArgForce) || viper.GetBool(constants.ArgQuiet) {
		return nil
	}
	answer, err := getUserInput(reader, writer, message)
	if err != nil {
		return errors.New("unable to parse users input " + err.Error())
	}

	//TODO: when responding no should not return error with invalid user input

	if answer != "y" && answer != "ye" && answer != "yes" {
		return errors.New("invalid user input")
	}
	return nil
}
