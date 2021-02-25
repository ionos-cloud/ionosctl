package utils

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
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
func AskForConfirm(reader io.Reader, writer printer.PrintService, message string) error {
	if viper.GetBool(config.ArgIgnoreStdin) || viper.GetBool(config.ArgQuiet) {
		return nil
	}
	answer, err := getUserInput(reader, writer, message)
	if err != nil {
		return errors.New("unable to parse users input " + err.Error())
	}

	if answer != "y" && answer != "ye" && answer != "yes" {
		return errors.New("invalid user input")
	}
	return nil
}
