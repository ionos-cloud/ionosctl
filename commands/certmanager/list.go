package certmanager

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
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

	cmd.Command.Flags().StringSlice(constants.ArgCols, nil, printer.ColsMessage(allCertificateCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCertificateCols, cobra.ShellCompDirectiveNoFileComp
	})
	return cmd
}

func CmdList(c *core.CommandConfig) error {
	fmt.Fprintf(c.Stderr, jsontabwriter.GenerateVerboseOutput("Getting Certificates..."))

	certs, _, err := c.CertificateManagerServices.Certs().List()
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("items", allCertificateJSONPaths, certs, allCertificateCols)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Stdout, out)

	return nil
}
