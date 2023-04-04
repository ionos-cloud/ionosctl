package registry

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go-container-registry"
	"github.com/spf13/cobra"
)

func RegGetCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace:  "container-registry",
			Resource:   "registry",
			Verb:       "get",
			Aliases:    []string{"g"},
			ShortDesc:  "Get Properties of a Registry",
			LongDesc:   "Get Properties of a single Registry",
			Example:    "ionosctl container-registry registry get --id [REGISTRY_ID]",
			PreCmdRun:  PreCmdGet,
			CmdRun:     CmdGet,
			InitClient: true,
		},
	)

	cmd.AddStringFlag(FlagRegId, "i", "", "Registry ID", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(
		FlagRegId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return RegsIds(), cobra.ShellCompDirectiveNoFileComp
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

func CmdGet(c *core.CommandConfig) error {
	id, err := c.Command.Command.Flags().GetString(FlagRegId)
	if err != nil {
		return err
	}
	reg, _, err := c.ContainerRegistryServices.Registry().Get(id)
	if err != nil {
		return err
	}
	return c.Printer.Print(getRegistryPrint(nil, c, &[]ionoscloud.RegistryResponse{reg}, false))
}

func PreCmdGet(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(c.Command, c.NS, FlagRegId)
	if err != nil {
		return err
	}

	return nil
}
