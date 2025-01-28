package registry

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RegDeleteCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace:  "container-registry",
			Resource:   "registry",
			Verb:       "delete",
			Aliases:    []string{"d"},
			ShortDesc:  "Delete a Registry",
			LongDesc:   "Delete a Registry.",
			Example:    "ionosctl container-registry registry delete --id [REGISTRY_ID]",
			PreCmdRun:  PreCmdDelete,
			CmdRun:     CmdDelete,
			InitClient: true,
		},
	)

	cmd.AddStringFlag(constants.FlagRegistryId, "i", "", "Specify the Registry ID", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagRegistryId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return RegsIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "Response delete all registries")

	cmd.Command.Flags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allCols, cobra.ShellCompDirectiveNoFileComp
		},
	)
	return cmd
}

func CmdDelete(c *core.CommandConfig) error {
	allFlag, err := c.Command.Command.Flags().GetBool(constants.ArgAll)
	if err != nil {
		return err
	}

	if allFlag {
		fmt.Fprintf(
			c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Deleting all Container Registries..."),
		)

		regs, _, err := c.ContainerRegistryServices.Registry().List("")
		if err != nil {
			return err
		}

		for _, reg := range regs.Items {
			msg := fmt.Sprintf("delete Container Registry: %s", *reg.Id)

			if !confirm.FAsk(c.Command.Command.InOrStdin(), msg, viper.GetBool(constants.ArgForce)) {
				return fmt.Errorf(confirm.UserDenied)
			}

			_, err = c.ContainerRegistryServices.Registry().Delete(*reg.Id)
			if err != nil {
				return err
			}
		}
	} else {
		id, err := c.Command.Command.Flags().GetString(constants.FlagRegistryId)
		if err != nil {
			return err
		}

		msg := fmt.Sprintf("delete Container Registry: %s", id)

		if !confirm.FAsk(c.Command.Command.InOrStdin(), msg, viper.GetBool(constants.ArgForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		_, err = c.ContainerRegistryServices.Registry().Delete(id)
		if err != nil {
			return err
		}
	}

	return nil
}

func PreCmdDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(
		c.Command, c.NS,
		[]string{constants.FlagRegistryId},
		[]string{constants.ArgAll},
	)
}
