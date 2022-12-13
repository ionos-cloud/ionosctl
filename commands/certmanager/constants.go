package certmanager

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"
)

// flags
const (
	FlagCertId         = "certificate-id"
	FlagCertName       = "certificate-name"
	FlagCert           = "certificate"
	FlagCertPath       = "certificate-path"
	FlagCertChain      = "certificate-chain"
	FlagCertChainPath  = "certificate-chain-path"
	FlagPrivateKey     = "private-key"
	FlagPrivateKeyPath = "private-key-path"

	UsageCert = "ionosctl certificate-manager create --name, --cert/--cert-path, --cert-chain/--cert-chain-path, --private-key/--private-key-path"
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
