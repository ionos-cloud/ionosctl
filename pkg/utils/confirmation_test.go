package utils

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testVar = "test"
	testErr = errors.New("error occurred")
)

func TestAskForConfirm_YesReader(t *testing.T) {
	input := getUserInput
	defer func() { getUserInput = input }()

	var b bytes.Buffer
	viper.Set(config.ArgIgnoreStdin, false)
	reader := bytes.NewReader([]byte("YES\n"))
	getUserInput = func(io.Reader, io.Writer, string) (string, error) {
		return readUserInput(reader, &b, "yes")
	}

	err := AskForConfirm(reader, &b, testVar)
	assert.NoError(t, err)
}

func TestAskForConfirm_EOF(t *testing.T) {
	input := getUserInput
	defer func() { getUserInput = input }()

	var b bytes.Buffer
	viper.Set(config.ArgIgnoreStdin, false)
	getUserInput = func(io.Reader, io.Writer, string) (string, error) {
		return readUserInput(os.Stdin, &b, "yes")
	}

	err := AskForConfirm(os.Stdin, &b, testVar)
	assert.Error(t, err)
	assert.True(t, err.Error() == `unable to parse users input EOF`)
}

func TestAskForConfirm_Yes(t *testing.T) {
	input := getUserInput
	defer func() { getUserInput = input }()

	var b bytes.Buffer
	viper.Set(config.ArgIgnoreStdin, false)
	getUserInput = func(io.Reader, io.Writer, string) (string, error) {
		return "y", nil
	}

	err := AskForConfirm(os.Stdin, &b, testVar)
	assert.NoError(t, err)
}

func TestAskForConfirm_No(t *testing.T) {
	input := getUserInput
	defer func() { getUserInput = input }()

	var b bytes.Buffer
	viper.Set(config.ArgIgnoreStdin, false)
	getUserInput = func(io.Reader, io.Writer, string) (string, error) {
		return "no", nil
	}

	err := AskForConfirm(os.Stdin, &b, testVar)
	assert.Error(t, err)
}

func TestAskForConfirm_Any(t *testing.T) {
	input := getUserInput
	defer func() { getUserInput = input }()

	var b bytes.Buffer
	viper.Set(config.ArgIgnoreStdin, false)
	getUserInput = func(io.Reader, io.Writer, string) (string, error) {
		return "dummy", nil
	}

	err := AskForConfirm(os.Stdin, &b, testVar)
	assert.Error(t, err)
}

func TestAskForConfirm_Error(t *testing.T) {
	input := getUserInput
	defer func() { getUserInput = input }()

	var b bytes.Buffer
	viper.Set(config.ArgIgnoreStdin, false)
	getUserInput = func(io.Reader, io.Writer, string) (string, error) {
		return "", testErr
	}

	err := AskForConfirm(os.Stdin, &b, testVar)
	assert.Error(t, err)
}

func TestAskForConfirm_IgnoreStdin(t *testing.T) {
	var b bytes.Buffer
	viper.Set(config.ArgIgnoreStdin, true)
	err := AskForConfirm(os.Stdin, &b, testVar)
	assert.NoError(t, err)
}
