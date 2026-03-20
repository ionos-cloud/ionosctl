package certificate

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
)

func CertGetCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace:  "certmanager",
		Resource:   "certificates",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get Certificate by ID",
		LongDesc:   "Use this command to retrieve a Certificate by ID.",
		Example:    "ionosctl certificate-manager get --certificate-id 47c5d9cc-b613-4b76-b0cc-dc531787a422",
		PreCmdRun:  PreCmdGet,
		CmdRun:     CmdGet,
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagCertId, constants.FlagIdShort, "", "Provide the specified Certificate", core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return CertificatesIds()
		}, constants.CertApiRegionalURL, constants.CertLocations),
	)
	cmd.AddBoolFlag(constants.FlagCert, "", false, "Print the certificate")
	cmd.AddBoolFlag(constants.FlagCertChain, "", false, "Print the certificate chain")

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

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	return c.Out(table.Sprint(allCols, cert, cols))
}

func PreCmdGet(c *core.PreCommandConfig) error {
	err := c.Command.Command.MarkFlagRequired(constants.FlagCertId)
	if err != nil {
		return err
	}

	c.Command.Command.MarkFlagsMutuallyExclusive(constants.FlagCert, constants.FlagCertChain)
	return nil
}
