package location

import (
	"context"
	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go-container-registry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	list := locs.GetItems()
	return c.Printer.Print(getLocPrint(nil, c, list))
}

type LocPrint struct {
	LocationId string `json:"RegistryId,omitempty"`
}

var allCols = structs.Names(LocPrint{})

func getLocPrint(
	resp *ionoscloud.APIResponse, c *core.CommandConfig, response *[]ionoscloud.Location,
) printer.Result {
	r := printer.Result{}
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	if c != nil {
		if resp != nil {
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForState = viper.GetBool(
				core.GetFlagName(
					c.NS, constants.ArgWaitForRequest,
				),
			) // this boolean is duplicated everywhere just to do an append of `& wait` to a verbose message
		}
		if response != nil {
			r.OutputJSON = response
			r.KeyValue = getLocRows(response)                       // map header -> rows
			r.Columns = printer.GetHeadersAllDefault(allCols, cols) // headers
		}
	}
	return r
}

func getLocRows(locs *[]ionoscloud.Location) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(*locs))
	for _, loc := range *locs {
		var locPrint LocPrint
		if idOk, ok := loc.GetIdOk(); ok && idOk != nil {
			locPrint.LocationId = *idOk
		}
		o := structs.Map(locPrint)
		out = append(out, o)
	}
	return out
}
