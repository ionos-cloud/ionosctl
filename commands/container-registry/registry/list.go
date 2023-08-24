package registry

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RegListCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace:  "container-registry",
			Resource:   "registry",
			Verb:       "list",
			Aliases:    []string{"l", "ls"},
			ShortDesc:  "List all Registries",
			LongDesc:   "List all managed container registries for your account",
			Example:    "ionosctl container-registry registry list",
			PreCmdRun:  core.NoPreRun,
			CmdRun:     CmdList,
			InitClient: true,
		},
	)

	cmd.AddStringFlag(FlagName, "n", "",
		"Response filter to list only the Registries that contain the specified name in the DisplayName field. The value is case insensitive",
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
	if viper.IsSet(core.GetFlagName(c.NS, FlagName)) {
		fmt.Fprintf(c.Stderr, jsontabwriter.GenerateVerboseOutput(
			"Filtering after Registry Name: %v", viper.GetString(core.GetFlagName(c.NS, "name"))))
	}

	regs, _, err := c.ContainerRegistryServices.Registry().List(viper.GetString(core.GetFlagName(c.NS, "name")))
	if err != nil {
		return err
	}

	cols, err := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("items", allJSONPaths, regs, printer.GetHeadersAllDefault(allCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Stdout, out)
	return nil
}
