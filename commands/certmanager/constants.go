package certmanager

// flags
const (
	FlagCertId              = "certificate-id"
	FlagCertName            = "certificate-name"
	FlagCert                = "certificate"
	FlagCertPath            = "certificate-path"
	FlagCertChain           = "certificate-chain"
	FlagCertChainPath       = "certificate-chain-path"
	FlagPrivateKey          = "private-key"
	FlagPrivateKeyPath      = "private-key-path"
	PostErrorFormatFlag = "%q requires at least 4 options by providing the required flag values or path values.\n\nUsage:\n%s\n%s\nFor more details, see '%s --help'."
	PostErrorExample1   = "ionosctl certificate-manager create --certificate-name CERTIFICATE_NAME --certificate CERTIFICATE --certificate-chain CERTIFICATE_CHAIN --private-key PRIVATE_KEY"
	PostErrorExample2   = "ionosctl certificate-manager create --certificate-name-path CERTIFICATE_NAME_PATH --certificate CERTIFICATE_PATH --certificate-chain-path CERTIFICATE_CHAIN_PATH --private-key-path PRIVATE_KEY_PATH\n"
	FlagAllFlag             = "all"
	FlagArgCols             = "cols"
	FlagArgWaitForState     = "wait-for-state"
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
