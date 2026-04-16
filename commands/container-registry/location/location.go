package location

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
)

var allCols = []table.Column{
	{Name: "LocationId", JSONPath: "id", Default: true},
}

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

	cmd.AddColsFlag(allCols)
	return cmd
}

func CmdList(c *core.CommandConfig) error {
	locs, _, err := client.Must().RegistryClient.LocationsApi.LocationsGet(context.Background()).Execute()
	if err != nil {
		return err
	}

	return c.Printer(allCols).Prefix("items").Print(locs)
}
