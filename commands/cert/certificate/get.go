package certificate

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func CertGetCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "certmanager",
		Resource:  "certificates",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Get a certificate by ID",
		LongDesc: `Retrieve a single uploaded certificate by ID, including metadata such as its serial number, validity dates (notBefore/notAfter), and whether it has expired.

By default the metadata is printed as a table. To print the PEM material instead, pass --certificate (the certificate body) or --certificate-chain (the chain); these two are mutually exclusive. The private key is write-only and is never returned.`,
		Example: `# Show certificate metadata
ionosctl certmanager certificate get --certificate-id 47c5d9cc-b613-4b76-b0cc-dc531787a422

# Print the certificate body PEM (e.g. to pipe into a file)
ionosctl certmanager certificate get --certificate-id 47c5d9cc-b613-4b76-b0cc-dc531787a422 --certificate`,
		PreCmdRun:  PreCmdGet,
		CmdRun:     CmdGet,
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagCertId, constants.FlagIdShort, "", "The ID (UUID) of the certificate to retrieve", core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return CertificatesIds()
		}, constants.CertApiRegionalURL, constants.CertLocations),
	)
	cmd.AddBoolFlag(constants.FlagCert, "", false, "Instead of the metadata table, print only the certificate body PEM. Mutually exclusive with --certificate-chain")
	cmd.AddBoolFlag(constants.FlagCertChain, "", false, "Instead of the metadata table, print only the certificate chain PEM. Mutually exclusive with --certificate")

	return cmd
}

func CmdGet(c *core.CommandConfig) error {
	var certFlag, certChainFlag, getCertOrChain bool

	certFlag, err := c.Command.Command.Flags().GetBool(constants.FlagCert)
	if err != nil {
		return err
	}

	certChainFlag, err = c.Command.Command.Flags().GetBool(constants.FlagCertChain)
	if err != nil {
		return err
	}

	if !certFlag && certChainFlag {
		getCertOrChain = true
	}

	c.Verbose("Getting Certificates...")

	id, err := c.Command.Command.Flags().GetString(constants.FlagCertId)
	if err != nil {
		return err
	}

	cert, _, err := client.Must().CertManagerClient.CertificateApi.CertificatesFindById(context.Background(), id).Execute()
	if err != nil {
		return err
	}

	if certFlag || certChainFlag {
		c.Msg("%s", printProperties(cert.Properties, getCertOrChain))

		return nil
	}

	return c.Printer(allCols).Print(cert)
}

func PreCmdGet(c *core.PreCommandConfig) error {
	err := c.Command.Command.MarkFlagRequired(constants.FlagCertId)
	if err != nil {
		return err
	}

	c.Command.Command.MarkFlagsMutuallyExclusive(constants.FlagCert, constants.FlagCertChain)
	return nil
}
