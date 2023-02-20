package token

import (
	"context"
	"github.com/ionos-cloud/ionosctl/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go-container-registry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TokenPostCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace:  "container-registry",
			Resource:   "token",
			Verb:       "create",
			Aliases:    []string{"c"},
			ShortDesc:  "Create a new token",
			LongDesc:   "Create a new token used to access a container registry",
			Example:    "ionosctl container-registry token create --registry-id [REGISTRY-ID] --name [TOKEN-NAME]",
			PreCmdRun:  PreCmdListToken,
			CmdRun:     CmdListToken,
			InitClient: true,
		},
	)

	cmd.AddStringFlag("name", "", "", "Name of the Token", core.RequiredFlagOption())
	cmd.AddStringFlag("expiry-date", "", "", "Expiry date of the Token")
	cmd.AddStringFlag("status", "", "", "Status of the Token")
	cmd.AddStringSliceFlag("scope-actions", "", []string{}, "Scope actions of the Token")
	cmd.AddStringFlag("scope-name", "", "", "Scope name of the Token")
	cmd.AddStringFlag("scope-type", "", "", "Scope type of the Token")

	cmd.AddStringFlag("registry-id", "r", "", "Registry ID", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(
		"registry-id", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return registry.RegsIds(), cobra.ShellCompDirectiveNoFileComp
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

func PreCmdPostToken(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(
		c.Command, c.NS,
		[]string{"registry-id", "name"},
		[]string{"registry-id", "name", "scope-actions", "scope-name", "scope-type"},
	)

	return nil
}

func CmdPostToken(c *core.CommandConfig) error {
	allFlag := viper.GetBool(core.GetFlagName(c.NS, "all"))
	if !allFlag {
		id := viper.GetString(core.GetFlagName(c.NS, "registry-id"))
		tokens, _, err := c.ContainerRegistryServices.Token().List(id)
		if err != nil {
			return err
		}
		list := tokens.GetItems()
		return c.Printer.Print(getTokenPrint(nil, c, list, false))
	}
	var list []ionoscloud.TokenResponse
	regs, _, err := c.ContainerRegistryServices.Registry().List("")
	if err != nil {
		return err
	}
	for _, reg := range *regs.GetItems() {
		tokens, _, err := c.ContainerRegistryServices.Token().List(*reg.Id)
		if err != nil {
			return err
		}
		list = append(list, *tokens.GetItems()...)
	}
	return c.Printer.Print(getTokenPrint(nil, c, &list, false))
}
