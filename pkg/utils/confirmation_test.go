package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	mockprinter "github.com/ionos-cloud/ionosctl/pkg/printer/mocks"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testVar = "test"
	testErr = errors.New("error occurred")
)

func TestAskForConfirmYesReader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	p := mockprinter.NewMockPrintService(ctrl)
	p.EXPECT().Print(getTestMessage(testVar)).Return(nil)

	input := getUserInput
	defer func() { getUserInput = input }()

	viper.Set(config.ArgForce, false)
	reader := bytes.NewReader([]byte("YES\n"))
	getUserInput = func(io.Reader, printer.PrintService, string) (string, error) {
		return readUserInput(reader, p, testVar)
	}

	err := AskForConfirm(reader, p, testVar)
	assert.NoError(t, err)
}

func TestAskForConfirmPrinterErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	p := mockprinter.NewMockPrintService(ctrl)
	p.EXPECT().Print(getTestMessage(testVar)).Return(testErr)

	input := getUserInput
	defer func() { getUserInput = input }()

	viper.Set(config.ArgForce, false)
	reader := bytes.NewReader([]byte("YES\n"))
	getUserInput = func(io.Reader, printer.PrintService, string) (string, error) {
		return readUserInput(reader, p, testVar)
	}

	err := AskForConfirm(reader, p, testVar)
	assert.Error(t, err)
}

func TestAskForConfirmEOF(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	p := mockprinter.NewMockPrintService(ctrl)
	p.EXPECT().Print(getTestMessage(testVar)).Return(nil)

	input := getUserInput
	defer func() { getUserInput = input }()

	viper.Set(config.ArgForce, false)
	getUserInput = func(io.Reader, printer.PrintService, string) (string, error) {
		return readUserInput(os.Stdin, p, testVar)
	}

	err := AskForConfirm(os.Stdin, p, testVar)
	assert.Error(t, err)
	assert.True(t, err.Error() == `unable to parse users input EOF`)
}

func TestAskForConfirmYes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	p := mockprinter.NewMockPrintService(ctrl)

	input := getUserInput
	defer func() { getUserInput = input }()

	viper.Set(config.ArgForce, false)
	getUserInput = func(io.Reader, printer.PrintService, string) (string, error) {
		return "y", nil
	}

	err := AskForConfirm(os.Stdin, p, testVar)
	assert.NoError(t, err)
}

func TestAskForConfirmNo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	p := mockprinter.NewMockPrintService(ctrl)

	input := getUserInput
	defer func() { getUserInput = input }()

	viper.Set(config.ArgForce, false)
	getUserInput = func(io.Reader, printer.PrintService, string) (string, error) {
		return "no", nil
	}

	err := AskForConfirm(os.Stdin, p, testVar)
	assert.Error(t, err)
}

func TestAskForConfirmAny(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	p := mockprinter.NewMockPrintService(ctrl)

	input := getUserInput
	defer func() { getUserInput = input }()

	viper.Set(config.ArgForce, false)
	getUserInput = func(io.Reader, printer.PrintService, string) (string, error) {
		return "dummy", nil
	}

	err := AskForConfirm(os.Stdin, p, testVar)
	assert.Error(t, err)
}

func TestAskForConfirmError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	p := mockprinter.NewMockPrintService(ctrl)

	input := getUserInput
	defer func() { getUserInput = input }()

	viper.Set(config.ArgForce, false)
	getUserInput = func(io.Reader, printer.PrintService, string) (string, error) {
		return "", testErr
	}

	err := AskForConfirm(os.Stdin, p, testVar)
	assert.Error(t, err)
}

func TestAskForConfirmIgnoreStdin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	p := mockprinter.NewMockPrintService(ctrl)
	viper.Set(config.ArgForce, true)
	err := AskForConfirm(os.Stdin, p, testVar)
	assert.NoError(t, err)
}

func TestGetUserInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	p := mockprinter.NewMockPrintService(ctrl)
	p.EXPECT().Print(getTestMessage(testVar)).Return(nil)
	reader := bytes.NewReader([]byte("YES\n"))
	_, err := getUserInput(reader, p, testVar)
	assert.NoError(t, err)
}

func getTestMessage(msg string) string {
	return fmt.Sprintf("Warning: Are you sure you want to %s (y/N) ? ", msg)
}
