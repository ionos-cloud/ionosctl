//go:build unit
// +build unit

package certmanager

import (
	"errors"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/stretchr/testify/assert"
)

var (
	nonAvailableCmdErr = errors.New("non-available cmd")
)

func TestCertificateManagerServiceCmdUnit(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(CertCmd())
	if ok := CertCmd().IsAvailableCommand(); !ok {
		err = nonAvailableCmdErr
	}
	assert.NoError(t, err)

	if ok := CertCreateCmd().IsAvailableCommand(); !ok {
		err = nonAvailableCmdErr
	}
	assert.NoError(t, err)

	if ok := CertDeleteCmd().IsAvailableCommand(); !ok {
		err = nonAvailableCmdErr
	}
	assert.NoError(t, err)

	if ok := CertGetCmd().IsAvailableCommand(); !ok {
		err = nonAvailableCmdErr
	}
	assert.NoError(t, err)

	if ok := CertListCmd().IsAvailableCommand(); !ok {
		err = nonAvailableCmdErr
	}
	assert.NoError(t, err)

	if ok := CertGetApiVersionCmd().IsAvailableCommand(); !ok {
		err = nonAvailableCmdErr
	}
	assert.NoError(t, err)
}
