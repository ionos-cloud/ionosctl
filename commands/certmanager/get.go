package certmanager

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	sdkgo "github.com/ionos-cloud/sdk-go-cert-manager"
	"github.com/spf13/cobra"
)

func CertGetCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
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

	cmd.AddStringFlag(FlagCertId, "i", "", "Response get a single certificate", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(FlagCertId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return CertificatesIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(FlagCert, "", false, "Print the certificate")
	cmd.AddBoolFlag(FlagCertChain, "", false, "Print the certificate chain")

	cmd.Command.Flags().StringSlice(constants.ArgCols, nil, printer.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

func CmdGet(c *core.CommandConfig) error {
	var certFlag, certChainFlag, getCertOrChain bool
	certFlag, err := c.Command.Command.Flags().GetBool(FlagCert)
	if err != nil {
		return err
	}
	certChainFlag, err = c.Command.Command.Flags().GetBool(FlagCertChain)
	if err != nil {
		return err
	}
	if !certFlag && certChainFlag {
		getCertOrChain = true
	}
	c.Printer.Verbose("Getting Certificates...")
	id, err := c.Command.Command.Flags().GetString(FlagCertId)
	if err != nil {
		return err
	}
	cert, _, err := c.CertificateManagerServices.Certs().Get(id)
	if err != nil {
		return err
	}

	if certFlag || certChainFlag {
		return c.Printer.Print(printProperties(&cert, getCertOrChain))
	}

	return c.Printer.Print(getCertPrint(nil, c, &[]sdkgo.CertificateDto{cert}))
}

func PreCmdGet(c *core.PreCommandConfig) error {
	err := c.Command.Command.MarkFlagRequired(FlagCertId)
	if err != nil {
		return err
	}

	c.Command.Command.MarkFlagsMutuallyExclusive(FlagCert, FlagCertChain)
	return nil
}
