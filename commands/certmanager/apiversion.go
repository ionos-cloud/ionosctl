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

var (
	APIVersionJSONPaths = map[string]string{
		"Href":    "href",
		"Name":    "name",
		"Version": "version",
	}

	allAPIVersionCols = []string{"Name", "Href", "Version"}
)

func CertGetApiVersionCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "certmanager",
		Resource:   "certificates",
		Verb:       "api-version",
		Aliases:    []string{"api", "info"},
		ShortDesc:  "Get Certificate Manager API Version",
		LongDesc:   "Use this command to retrieve API Version.",
		Example:    "ionosctl certificate-manager api-version",
		PreCmdRun:  core.NoPreRun,
		CmdRun:     CmdGetApiVersion,
		InitClient: true,
	})

	cmd.Command.Flags().StringSlice(constants.ArgCols, nil, printer.ColsMessage(allCertificateCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allAPIVersionCols, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(constants.ArgNoHeaders, "n", false, "Response delete all certificates")

	return cmd
}

func CmdGetApiVersion(c *core.CommandConfig) error {
	fmt.Fprintf(c.Stderr, jsontabwriter.GenerateVerboseOutput("Getting Api Version..."))

	cols, err := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if err != nil {
		return err
	}

	APIVersion, _, err := c.CertificateManagerServices.Certs().GetApiVersion()
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("", APIVersionJSONPaths, APIVersion, printer.GetHeadersAllDefault(allAPIVersionCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Stdout, out)

	return nil
}
