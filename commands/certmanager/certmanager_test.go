package certmanager

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
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

			// os.Mkdir("testPaths", 777)
			caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
			assert.NoError(t, err)

			caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
			assert.NoError(t, err)

			caPEM := new(bytes.Buffer)
			pem.Encode(caPEM, &pem.Block{
				Type:  "CERTIFICATE",
				Bytes: caBytes,
			})

			caPrivKeyPEM := new(bytes.Buffer)
			pem.Encode(caPrivKeyPEM, &pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
			})

			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgVerbose, false)

			c := CertCreateCmd()
			c.Command.Flags().Set(FlagCertName, "certificate")
			c.Command.Flags().Set(FlagCert, caPEM.String())
			c.Command.Flags().Set(FlagCertChain, caPEM.String())
			c.Command.Flags().Set(FlagPrivateKey, caPrivKeyPEM.String())

			err = c.Command.Execute()
			assert.NoError(t, err)
		},
	)
}
