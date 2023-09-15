package location

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/tabheaders"
	"github.com/spf13/cobra"
)

var (
	allJSONPaths = map[string]string{
		"LocationId": "id",
	}

	allCols = []string{"LocationId"}
)

func RegLocationsListCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace:  "container-registry",
			Resource:   "registry",
			Verb:       "locations",
			Aliases:    []string{"location", "loc", "l", "locs"},
			ShortDesc:  "List all Registries Locations",
			LongDesc:   "List all managed container registries locations for your account",
			Example:    "ionosctl container-registry locations",
			PreCmdRun:  core.NoPreRun,
			CmdRun:     CmdList,
			InitClient: true,
		},
	)

	cmd.Command.Flags().StringSlice(constants.ArgCols, nil, printer.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allCols, cobra.ShellCompDirectiveNoFileComp
		},
	)
	return cmd
}

func CmdList(c *core.CommandConfig) error {
	locs, _, err := c.ContainerRegistryServices.Location().Get()
	if err != nil {
		return err
	}

	cols, err := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("items", allJSONPaths, locs, tabheaders.GetHeadersAllDefault(allCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}
