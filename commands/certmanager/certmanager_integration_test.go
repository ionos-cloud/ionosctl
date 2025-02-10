//go:build integration
// +build integration

package certmanager

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/commands/certmanager/certificate"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// test values
var (
	ca = &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{"Company, INC."},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{"Golden Gate Bridge"},
			PostalCode:    []string{"94016"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
)

func TestCertificateManagerServiceCmd(t *testing.T) {

	var err error
	t.Run(
		"cert list", func(t *testing.T) {
			viper.Reset()
			// os.Clearenv()

			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(constants.ArgQuiet, true)
			viper.Set(constants.ArgVerbose, false)

			err = certificate.CertListCmd().Command.Execute()
			assert.NoError(t, err)
		},
	)

	t.Run(
		"cert create from flags", func(t *testing.T) {
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
			viper.Set(constants.ArgForce, true)

			c := certificate.CertCreateCmd()
			c.Command.Flags().Set(constants.FlagCertName, "___test___certificate___test___")
			c.Command.Flags().Set(constants.FlagCert, caPEM.String())
			c.Command.Flags().Set(constants.FlagCertChain, caPEM.String())
			c.Command.Flags().Set(constants.FlagPrivateKey, caPrivKeyPEM.String())

			err = c.Command.Execute()
			assert.NoError(t, err)

			// var id string
			svc, err := client.Get()
			certs, _, err := svc.CertManagerClient.CertificateApi.CertificatesGet(context.Background()).Execute()
			assert.NoError(t, err)

			var id string
			for _, dto := range certs.GetItems() {
				if dto.GetProperties().Name == "___test___certificate___test___" {
					id = dto.GetId()
				}
			}

			g := certificate.CertGetCmd()
			g.Command.Flags().Set(constants.FlagCertId, id)
			assert.NoError(t, err)

			err = g.Command.Execute()
			assert.NoError(t, err)

			d := certificate.CertDeleteCmd()
			d.Command.Flags().Set(constants.FlagCertId, id)
			err = d.Command.Execute()
			assert.NoError(t, err)
		},
	)

	t.Run(
		"cert create from files", func(t *testing.T) {
			viper.Reset()

			// os.Mkdir("./testPaths", 0777)
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
			certPath := filepath.Join(".", "cert.pem")
			// os.Create(certPath)
			err = os.WriteFile(certPath, caPEM.Bytes(), 0777)
			assert.NoError(t, err)

			keyPath := filepath.Join(".", "key.pem")
			// os.Create(keyPath)
			os.WriteFile(keyPath, caPrivKeyPEM.Bytes(), 0777)
			assert.NoError(t, err)

			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(constants.ArgQuiet, true)
			viper.Set(constants.ArgVerbose, false)
			viper.Set(constants.ArgForce, true)

			c := certificate.CertCreateCmd()
			c.Command.Flags().Set(constants.FlagCertName, "test_certificate-files_test")
			c.Command.Flags().Set(constants.FlagCertPath, certPath)
			c.Command.Flags().Set(constants.FlagCertChainPath, certPath)
			c.Command.Flags().Set(constants.FlagPrivateKeyPath, keyPath)

			err = c.Command.Execute()
			assert.NoError(t, err)

			// var id string
			svc, err := client.Get()
			certs, _, err := svc.CertManagerClient.CertificateApi.CertificatesGet(context.Background()).
				Execute()
			assert.NoError(t, err)

			var id string
			for _, dto := range certs.GetItems() {
				if dto.GetProperties().Name == "test_certificate-files_test" {
					id = dto.GetId()
				}
			}

			p := certificate.CertUpdateCmd()
			p.Command.Flags().Set(constants.FlagCertId, id)
			p.Command.Flags().Set(constants.FlagCertName, "test_certificate-files-updated_test")
			err = p.Command.Execute()

			cert, _, err := svc.CertManagerClient.CertificateApi.CertificatesFindById(context.Background(), id).Execute()
			assert.NoError(t, err)
			assert.Equal(t, "test_certificate-files-updated_test", cert.GetProperties().Name)

			d := certificate.CertDeleteCmd()
			d.Command.Flags().Set(constants.FlagCertId, id)
			err = d.Command.Execute()
			assert.NoError(t, err)

			err = os.Remove(certPath)
			assert.NoError(t, err)

			err = os.Remove(keyPath)
			assert.NoError(t, err)
		},
	)
}
