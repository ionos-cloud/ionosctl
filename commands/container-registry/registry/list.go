package registry

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/viper"
)

func RegListCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace:  "container-registry",
			Resource:   "registry",
			Verb:       "list",
			Aliases:    []string{"l", "ls"},
			ShortDesc:  "List all registries in the contract",
			LongDesc:   "List all container registries in your contract, showing each registry's hostname, location, garbage-collection schedule, whether vulnerability scanning is enabled, and its state. Use --name to filter by a substring of the display name.",
			Example:    "ionosctl container-registry registry list",
			PreCmdRun:  core.NoPreRun,
			CmdRun:     CmdList,
			InitClient: true,
		},
	)

	cmd.AddStringFlag(
		constants.FlagName, constants.FlagNameShort, "",
		"Filter: list only registries whose DisplayName contains this substring (case-insensitive)",
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

	return c.Printer(allCols).Prefix("items").Print(regs)
}
