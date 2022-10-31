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
		LongDesc:  "Use this command to retrieve a  Certificate by ID.",
		Example:   "ionosctl certificate-manager get --certificate-id 12345678-1234-1234-1234-123456789012",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			err := c.Command.Command.MarkFlagRequired("certificate-id")
			if err != nil {
				return err
			}
			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			c.Printer.Verbose("Getting Certificates...")
			id, err := c.Command.Command.Flags().GetString("certificate-id")
			if err != nil {
				return err
			}
			cert, _, err := c.CertificateManagerServices.Certs().Get(id)
			if err != nil {
				return err
			}

			return c.Printer.Print(getCertPrint(nil, c, &[]sdkgo.CertificateDto{cert}))
		},
		InitClient: true,
	})

	//TODO: add certifcate and certificate-chain boolean flags

	cmd.AddStringFlag("certificate-id", "", "", "Response get a single certificate (required)")
	_ = cmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
