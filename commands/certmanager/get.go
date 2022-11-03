package certmanager

import (
	"context"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/spf13/cobra"

	sdkgo "github.com/ionos-cloud/sdk-go-cert-manager"
)

func CertGetCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "certmanager",
		Resource:  "certificates",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Get Certificate by ID",
		LongDesc:  "Use this command to retrieve a Certificate by ID.",
		Example:   "ionosctl certificate-manager get --certificate-id 47c5d9cc-b613-4b76-b0cc-dc531787a422",
		PreCmdRun: PreCmdGet,
		CmdRun: CmdGet,
		InitClient: true,
	})

	cmd.AddStringFlag(CertId, "i", "", "Response get a single certificate (required)")
	cmd.AddBoolFlag(Cert, "", false, "Print the certificate")
	cmd.AddBoolFlag(CertChain, "", false, "Print the certificate chain")
	_ = cmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

func CmdGet(c *core.CommandConfig) error {
	var certFlag, certChainFlag bool
	certFlag, err := c.Command.Command.Flags().GetBool(Cert)
	if err != nil {
		return err
	}
	certChainFlag, err = c.Command.Command.Flags().GetBool(CertChain)
	if err != nil {
		return err
	}
	c.Printer.Verbose("Getting Certificates...")
	id, err := c.Command.Command.Flags().GetString(CertId)
	if err != nil {
		return err
	}
	cert, _, err := c.CertificateManagerServices.Certs().Get(id)
	if err != nil {
		return err
	}

	flags := []bool{certFlag, certChainFlag}
	if certFlag || certChainFlag {
		return c.Printer.Print(printProperties(&cert, c, flags))
	}

	return c.Printer.Print(getCertPrint(nil, c, &[]sdkgo.CertificateDto{cert}))
}

func PreCmdGet(c *core.PreCommandConfig) error {
	err := c.Command.Command.MarkFlagRequired(CertId)
	if err != nil {
		return err
	}

	c.Command.Command.MarkFlagsMutuallyExclusive(Cert, CertChain)
	return nil
}
