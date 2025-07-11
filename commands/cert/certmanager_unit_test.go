//go:build unit
// +build unit

package cert

import (
	"errors"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/commands/cert/certificate"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/stretchr/testify/assert"
)

var (
	nonAvailableCmdErr = errors.New("non-available cmd")
)

func TestCertificateManagerServiceCmdUnit(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(certificate.CertCmd())
	if ok := certificate.CertCmd().IsAvailableCommand(); !ok {
		err = nonAvailableCmdErr
	}
	assert.NoError(t, err)

	if ok := certificate.CertCreateCmd().IsAvailableCommand(); !ok {
		err = nonAvailableCmdErr
	}
	assert.NoError(t, err)

	if ok := certificate.CertDeleteCmd().IsAvailableCommand(); !ok {
		err = nonAvailableCmdErr
	}
	assert.NoError(t, err)

	if ok := certificate.CertGetCmd().IsAvailableCommand(); !ok {
		err = nonAvailableCmdErr
	}
	assert.NoError(t, err)

	if ok := certificate.CertListCmd().IsAvailableCommand(); !ok {
		err = nonAvailableCmdErr
	}
	assert.NoError(t, err)

	if ok := CertGetApiVersionCmd().IsAvailableCommand(); !ok {
		err = nonAvailableCmdErr
	}
	assert.NoError(t, err)
}
