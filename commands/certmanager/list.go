package certmanager

import (
	"context"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/spf13/cobra"
)

func CertListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "certmanager",
		Resource:   "certificates",
		Verb:       "list",
		Aliases:    []string{"l"},
		ShortDesc:  "List all Certificates",
		LongDesc:   "Use this command to retrieve all Certificates.",
		Example:    "ionosctl certificate-manager list",
		PreCmdRun:  core.NoPreRun,
		CmdRun:     CmdList,
		InitClient: true,
	})
	_ = cmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

func CmdList(c *core.CommandConfig) error {
	c.Printer.Verbose("Getting Certificates...")
	cert, _, err := c.CertificateManagerServices.Certs().List()
	if err != nil {
		return err
	}
	list := cert.GetItems()
	return c.Printer.Print(getCertPrint(nil, c, list))
}
