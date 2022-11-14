package certmanager

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

var RequiredFlagSets = [16][]string{{FlagCertName, FlagCert, FlagCertChain, FlagPrivateKey},
	{FlagCertName, FlagCertPath, FlagCertChain, FlagPrivateKey},
	{FlagCertName, FlagCert, FlagCertChainPath, FlagPrivateKey},
	{FlagCertName, FlagCertPath, FlagCertChainPath, FlagPrivateKey},
	{FlagCertName, FlagCert, FlagCertChain, FlagPrivateKeyPath},
	{FlagCertName, FlagCertPath, FlagCertChain, FlagPrivateKeyPath},
	{FlagCertName, FlagCert, FlagCertChainPath, FlagPrivateKeyPath},
	{FlagCertName, FlagCertPath, FlagCertChainPath, FlagPrivateKeyPath},
}
