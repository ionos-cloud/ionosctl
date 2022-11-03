package certmanager

// flags
const (
	CertId              = "certificate-id"
	CertName            = "certificate-name"
	CertNamePath        = "certificate-name-path"
	Cert                = "certificate"
	CertPath            = "certificate-path"
	CertChain           = "certificate-chain"
	CertChainPath       = "certificate-chain-path"
	PrivateKey          = "private-key"
	PrivateKeyPath      = "private-key-path"
	PostErrorFormatFlag = "%q requires at least 4 options by providing the required flag values or path values.\n\nUsage:\n%s\n%s\nFor more details, see '%s --help'."
	PostErrorExample1   = "ionosctl certificate-manager create --certificate-name CERTIFICATE_NAME --certificate CERTIFICATE --certificate-chain CERTIFICATE_CHAIN --private-key PRIVATE_KEY"
	PostErrorExample2   = "ionosctl certificate-manager create --certificate-name-path CERTIFICATE_NAME_PATH --certificate CERTIFICATE_PATH --certificate-chain-path CERTIFICATE_CHAIN_PATH --private-key-path PRIVATE_KEY_PATH\n"
	AllFlag             = "all"
)

var RequiredFlagSets = [16][]string{{CertName, Cert, CertChain, PrivateKey},
	{CertNamePath, Cert, CertChain, PrivateKey},
	{CertName, CertPath, CertChain, PrivateKey},
	{CertNamePath, CertPath, CertChain, PrivateKey},
	{CertName, Cert, CertChainPath, PrivateKey},
	{CertNamePath, Cert, CertChainPath, PrivateKey},
	{CertName, CertPath, CertChainPath, PrivateKey},
	{CertNamePath, CertPath, CertChainPath, PrivateKey},
	{CertName, Cert, CertChain, PrivateKeyPath},
	{CertNamePath, Cert, CertChain, PrivateKeyPath},
	{CertName, CertPath, CertChain, PrivateKeyPath},
	{CertNamePath, CertPath, CertChain, PrivateKeyPath},
	{CertName, Cert, CertChainPath, PrivateKeyPath},
	{CertNamePath, Cert, CertChainPath, PrivateKeyPath},
	{CertName, CertPath, CertChainPath, PrivateKeyPath},
	{CertNamePath, CertPath, CertChainPath, PrivateKeyPath}}
