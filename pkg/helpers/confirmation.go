package helpers

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
)

func retrieveUserInput(message string) (string, error) {
	return readUserInput(os.Stdin, message)
}

func warnConfirm(msg string, args ...interface{}) {
	colorWarn := color.YellowString("Warning")
	fmt.Fprintf(color.Output, "%s: %s", colorWarn, fmt.Sprintf(msg, args...))
}

func readUserInput(in io.Reader, message string) (string, error) {
	reader := bufio.NewReader(in)
	warnConfirm("Are you sure you want to " + message + " (y/N) ? ")
	answer, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	answer = strings.TrimRight(answer, "\r\n")

	return strings.ToLower(answer), nil
}

// AskForConfirm parses and verifies user input for confirmation.
func AskForConfirm(message string) error {
	answer, err := retrieveUserInput(message)
	if err != nil {
		return errors.New("unable to parse users input " + err.Error())
	}

	if answer != "y" && answer != "ye" && answer != "yes" {
		return errors.New("invalid user input")
	}

	return nil
}
