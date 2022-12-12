package certmanager

import (
	"errors"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/ionos-cloud/ionosctl/pkg/constants"
)

func TestCertificateManagerServiceCmd(t *testing.T) {
	var err error
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

	t.Run(
		"cert list", func(t *testing.T) {
			viper.Reset()
			// os.Clearenv()

			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(constants.ArgQuiet, true)
			viper.Set(constants.ArgVerbose, false)

			err = CertListCmd().Command.Execute()
			assert.NoError(t, err)
		},
	)

	t.Run(
		"cert create", func(t *testing.T) {
			viper.Reset()

			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgVerbose, false)


			CertCreateCmd().Command.SetArgs([]string{"certificate-name", "cert_test"} )
			CertCreateCmd().Command.Execute()
			assert.NoError(t, err)
		},
	)
}
