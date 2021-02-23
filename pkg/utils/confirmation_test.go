package utils

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	mockprinter "github.com/ionos-cloud/ionosctl/pkg/utils/printer/mocks"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testVar = "test"
	testErr = errors.New("error occurred")
)

func TestAskForConfirm_YesReader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	p := mockprinter.NewMockPrintService(ctrl)
	p.EXPECT().Print("Warning: Are you sure you want to " + testVar + " (y/N) ? ").Return(nil)

	input := getUserInput
	defer func() { getUserInput = input }()

	viper.Set(config.ArgIgnoreStdin, false)
	reader := bytes.NewReader([]byte("YES\n"))
	getUserInput = func(io.Reader, printer.PrintService, string) (string, error) {
		return readUserInput(reader, p, testVar)
	}

	err := AskForConfirm(reader, p, testVar)
	assert.NoError(t, err)
}

func TestAskForConfirm_PrinterErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	p := mockprinter.NewMockPrintService(ctrl)
	p.EXPECT().Print("Warning: Are you sure you want to " + testVar + " (y/N) ? ").Return(testErr)

	input := getUserInput
	defer func() { getUserInput = input }()

	viper.Set(config.ArgIgnoreStdin, false)
	reader := bytes.NewReader([]byte("YES\n"))
	getUserInput = func(io.Reader, printer.PrintService, string) (string, error) {
		return readUserInput(reader, p, testVar)
	}

	err := AskForConfirm(reader, p, testVar)
	assert.Error(t, err)
}

func TestAskForConfirm_EOF(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	p := mockprinter.NewMockPrintService(ctrl)
	p.EXPECT().Print("Warning: Are you sure you want to " + testVar + " (y/N) ? ").Return(nil)

	input := getUserInput
	defer func() { getUserInput = input }()

	viper.Set(config.ArgIgnoreStdin, false)
	getUserInput = func(io.Reader, printer.PrintService, string) (string, error) {
		return readUserInput(os.Stdin, p, testVar)
	}

	err := AskForConfirm(os.Stdin, p, testVar)
	assert.Error(t, err)
	assert.True(t, err.Error() == `unable to parse users input EOF`)
}

func TestAskForConfirm_Yes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	p := mockprinter.NewMockPrintService(ctrl)

	input := getUserInput
	defer func() { getUserInput = input }()

	viper.Set(config.ArgIgnoreStdin, false)
	getUserInput = func(io.Reader, printer.PrintService, string) (string, error) {
		return "y", nil
	}

	err := AskForConfirm(os.Stdin, p, testVar)
	assert.NoError(t, err)
}

func TestAskForConfirm_No(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	p := mockprinter.NewMockPrintService(ctrl)

	input := getUserInput
	defer func() { getUserInput = input }()

	viper.Set(config.ArgIgnoreStdin, false)
	getUserInput = func(io.Reader, printer.PrintService, string) (string, error) {
		return "no", nil
	}

	err := AskForConfirm(os.Stdin, p, testVar)
	assert.Error(t, err)
}

func TestAskForConfirm_Any(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	p := mockprinter.NewMockPrintService(ctrl)

	input := getUserInput
	defer func() { getUserInput = input }()

	viper.Set(config.ArgIgnoreStdin, false)
	getUserInput = func(io.Reader, printer.PrintService, string) (string, error) {
		return "dummy", nil
	}

	err := AskForConfirm(os.Stdin, p, testVar)
	assert.Error(t, err)
}

func TestAskForConfirm_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	p := mockprinter.NewMockPrintService(ctrl)

	input := getUserInput
	defer func() { getUserInput = input }()

	viper.Set(config.ArgIgnoreStdin, false)
	getUserInput = func(io.Reader, printer.PrintService, string) (string, error) {
		return "", testErr
	}

	err := AskForConfirm(os.Stdin, p, testVar)
	assert.Error(t, err)
}

func TestAskForConfirm_IgnoreStdin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	p := mockprinter.NewMockPrintService(ctrl)
	viper.Set(config.ArgIgnoreStdin, true)
	err := AskForConfirm(os.Stdin, p, testVar)
	assert.NoError(t, err)
}

func TestGetUserInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	p := mockprinter.NewMockPrintService(ctrl)
	p.EXPECT().Print("Warning: Are you sure you want to " + testVar + " (y/N) ? ").Return(nil)
	reader := bytes.NewReader([]byte("YES\n"))
	_, err := getUserInput(reader, p, testVar)
	assert.NoError(t, err)
}
