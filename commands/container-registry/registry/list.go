package registry

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/cobra"
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

	cmd.AddStringFlag(
		constants.FlagName, constants.FlagNameShort, "",
		"Response filter to list only the Registries that contain the specified name in the DisplayName field. The value is case insensitive",
	)

	cmd.Command.Flags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allCols, cobra.ShellCompDirectiveNoFileComp
		},
	)
	return cmd
}

func CmdList(c *core.CommandConfig) error {
	filterName, _ := c.Command.Command.Flags().GetString(constants.FlagName)
	if c.Command.Command.Flags().Changed(constants.FlagName) {
		fmt.Fprintf(
			c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
				"Filtering after Registry Name: %v", filterName,
			),
		)
	}
	req := client.Must().RegistryClient.RegistriesApi.RegistriesGet(context.Background())
	if filterName != "" {
		req = req.FilterName(filterName)
	}
	regs, _, err := client.Must().RegistryClient.RegistriesApi.RegistriesGetExecute(req)
	if err != nil {
		return err
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	out, err := jsontabwriter.GenerateOutput(
		"items", jsonpaths.ContainerRegistryRegistry, regs, tabheaders.GetHeadersAllDefault(allCols, cols),
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}
