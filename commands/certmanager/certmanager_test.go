package certmanager

import (
	"errors"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestCertificateManagerServiceCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(CertCmd())
	if ok := CertCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)

	if ok := CertCreateCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)

	if ok := CertDeleteCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)

	if ok := CertGetCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)

	if ok := CertListCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)

	if ok := CertGetApiVersionCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)

	t.Run("cert list", func(t *testing.T) {
		core.RootCmdTest.Command.SetArgs([]string{"list"})
		err = core.RootCmdTest.Command.Execute()
	})
}
