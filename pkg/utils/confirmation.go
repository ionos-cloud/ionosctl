package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/spf13/viper"
)

var getUserInput = func(reader io.Reader, writer io.Writer, message string) (string, error) {
	return readUserInput(reader, writer, message)
}

func warnConfirm(writer io.Writer, msg string, args ...interface{}) {
	fmt.Fprintf(writer, "Warning: %s", fmt.Sprintf(msg, args...))
}

func readUserInput(reader io.Reader, writer io.Writer, message string) (string, error) {
	in := bufio.NewReader(reader)
	warnConfirm(writer, "Are you sure you want to "+message+" (y/N) ? ")
	answer, err := in.ReadString('\n')
	if err != nil {
		return "", err
	}

	answer = strings.TrimRight(answer, "\r\n")
	return strings.ToLower(answer), nil
}

// AskForConfirm parses and verifies user input for confirmation.
func AskForConfirm(reader io.Reader, writer io.Writer, message string) error {
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
