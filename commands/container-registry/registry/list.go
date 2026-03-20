package registry

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
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

	cmd.AddStringFlag(
		constants.FlagName, constants.FlagNameShort, "",
		"Response filter to list only the Registries that contain the specified name in the DisplayName field. The value is case insensitive",
	)

	return cmd
}

func CmdList(c *core.CommandConfig) error {
	if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
		c.Verbose("Filtering after Registry Name: %v", viper.GetString(fn))
	}

	filterName := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
	req := client.Must().RegistryClient.RegistriesApi.RegistriesGet(context.Background())
	if filterName != "" {
		req = req.FilterName(filterName)
	}
	regs, _, err := client.Must().RegistryClient.RegistriesApi.RegistriesGetExecute(req)
	if err != nil {
		return err
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	return c.Out(table.Sprint(allCols, regs, cols, table.WithPrefix("items")))
}
