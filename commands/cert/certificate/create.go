package certificate

import (
	"context"
	"fmt"
	"os"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/cert/v2"
)

func CertCreateCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "certmanager",
		Resource:  "certificate",
		Verb:      "create",
		Aliases:   []string{"add", "a", "c", "post"},
		ShortDesc: "Upload a new certificate",
		LongDesc: `Upload a TLS/SSL certificate so IONOS products (ALB, CDN) can serve HTTPS under your domain.

You must supply three PEM values, each either inline or from a file (the inline and *-path variants are mutually exclusive, and exactly one of each pair is required):
  * the certificate body (--certificate / --certificate-path)
  * the certificate chain of intermediate CAs (--certificate-chain / --certificate-chain-path)
  * the private key that pairs with the certificate (--private-key / --private-key-path)

The private key is write-only; it is never returned by get/list. This command does not renew certificates - upload a fresh one before expiry. For automatic issuance and renewal, use ` + "`autocertificate create`" + ` instead.`,
		Example: `# Upload from PEM files (recommended)
ionosctl certmanager certificate create --certificate-name my-cert --certificate-path cert.pem --certificate-chain-path chain.pem --private-key-path key.pem

# Upload with inline PEM values read from files by the shell
ionosctl certmanager certificate create --certificate-name my-cert --certificate "$(cat cert.pem)" --certificate-chain "$(cat chain.pem)" --private-key "$(cat key.pem)"`,
		PreCmdRun: func(c *core.PreCommandConfig) error {
			err := c.Command.Command.MarkFlagRequired(constants.FlagCertName)
			if err != nil {
				return fmt.Errorf("unable to mark flag %s as required: %w. Please report this at https://github.com/ionos-cloud/ionosctl/issues/new", constants.FlagCertName, err)
			}

			c.Command.Command.MarkFlagsMutuallyExclusive(constants.FlagCert, constants.FlagCertPath)
			c.Command.Command.MarkFlagsMutuallyExclusive(constants.FlagCertChain, constants.FlagCertChainPath)
			c.Command.Command.MarkFlagsMutuallyExclusive(constants.FlagPrivateKey, constants.FlagPrivateKeyPath)

			c.Command.Command.MarkFlagsOneRequired(constants.FlagCert, constants.FlagCertPath)
			c.Command.Command.MarkFlagsOneRequired(constants.FlagCertChain, constants.FlagCertChainPath)
			c.Command.Command.MarkFlagsOneRequired(constants.FlagPrivateKey, constants.FlagPrivateKeyPath)

			return nil
		},
		CmdRun:     CmdPost,
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagCertName, "n", "", "A friendly name for the certificate, used to identify it in listings and when other products reference it", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagCert, "", "", "The certificate body in PEM format, inline. Provide this or --certificate-path, not both")
	cmd.AddStringFlag(constants.FlagCertChain, "", "", "The chain of intermediate CA certificates in PEM format, inline. Provide this or --certificate-chain-path, not both")
	cmd.AddStringFlag(constants.FlagPrivateKey, "", "", "The private key paired with the certificate, in PEM format, inline. Write-only: never returned by get/list. Provide this or --private-key-path, not both")

	cmd.AddStringFlag(constants.FlagCertPath, "", "", "Path to a file containing the certificate body in PEM format. Provide this or --certificate, not both")
	cmd.AddStringFlag(constants.FlagCertChainPath, "", "", "Path to a file containing the intermediate CA certificate chain in PEM format. Provide this or --certificate-chain, not both")
	cmd.AddStringFlag(constants.FlagPrivateKeyPath, "", "", "Path to a file containing the private key in PEM format. Provide this or --private-key, not both")

	return cmd
}

func CmdPost(c *core.CommandConfig) error {
	input := cert.Certificate{}

	if c.Command.Command.Flags().Changed(constants.FlagCertName) {
		name, _ := c.Command.Command.Flags().GetString(constants.FlagCertName)

		input.Name = name
	}

	if c.Command.Command.Flags().Changed(constants.FlagCert) {
		certValue, _ := c.Command.Command.Flags().GetString(constants.FlagCert)

		input.Certificate = certValue
	}

	if c.Command.Command.Flags().Changed(constants.FlagCertChain) {
		certChain, _ := c.Command.Command.Flags().GetString(constants.FlagCertChain)

		input.CertificateChain = certChain
	}

	if c.Command.Command.Flags().Changed(constants.FlagPrivateKey) {
		privateKey, _ := c.Command.Command.Flags().GetString(constants.FlagPrivateKey)

		input.PrivateKey = privateKey
	}

	if c.Command.Command.Flags().Changed(constants.FlagCertPath) {
		certPath, _ := c.Command.Command.Flags().GetString(constants.FlagCertPath)

		bytes, err := os.ReadFile(certPath)
		if err != nil {
			return fmt.Errorf("error reading certificate file: %w", err)
		}
		input.Certificate = string(bytes)
	}

	if c.Command.Command.Flags().Changed(constants.FlagCertChainPath) {
		certChainPath, _ := c.Command.Command.Flags().GetString(constants.FlagCertChainPath)

		bytes, err := os.ReadFile(certChainPath)
		if err != nil {
			return fmt.Errorf("error reading certificate chain file: %w", err)
		}
		input.CertificateChain = string(bytes)
	}

	if c.Command.Command.Flags().Changed(constants.FlagPrivateKeyPath) {
		privateKeyPath, _ := c.Command.Command.Flags().GetString(constants.FlagPrivateKeyPath)

		bytes, err := os.ReadFile(privateKeyPath)
		if err != nil {
			return fmt.Errorf("error reading private key file: %w", err)
		}
		input.PrivateKey = string(bytes)
	}

	cert, _, err := client.Must().CertManagerClient.CertificateApi.CertificatesPost(context.Background()).
		CertificateCreate(cert.CertificateCreate{Properties: input}).Execute()
	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(cert)
}
