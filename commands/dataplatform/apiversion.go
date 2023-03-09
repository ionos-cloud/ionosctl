package dataplatform

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"

	"github.com/ionos-cloud/ionosctl/v6/pkg/config"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
)

func ApiVersionCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dataplatform",
		Verb:      "api-versions",
		Aliases:   []string{"versions", "api-version", "api", "meta"},
		ShortDesc: "Get API version",
		Example:   "ionosctl dataplatform meta",
		PreCmdRun: core.NoPreRun,
		CmdRun: func(c *core.CommandConfig) error {
			client, err := config.GetClient()
			list, _, err := client.DataplatformClient.DataPlatformMetaDataApi.VersionsGet(c.Context).Execute()
			if err != nil {
				return fmt.Errorf("fail VersionsGet: %w", err)
			}
			return c.Printer.Print(getApiVersionPrint(c, list))
		},
		InitClient: true,
	})

	cmd.AddBoolFlag(constants.ArgNoHeaders, "", false, "When using text output, don't print headers")
	cmd.AddStringSliceFlag(constants.ArgCols, "", nil, printer.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

type ApiVersionPrint struct {
	Version string `json:"Version,omitempty"`
}

var allCols = structs.Names(ApiVersionPrint{})

func MakeApiVersionPrintObject(objs []string) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(objs))

	for _, o := range objs {
		var printObj ApiVersionPrint
		printObj.Version = o
		o := structs.Map(printObj)

		out = append(out, o)
	}

	return out
}

func getApiVersionPrint(c *core.CommandConfig, dcs []string) printer.Result {
	r := printer.Result{}
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	if c != nil && dcs != nil {
		r.OutputJSON = dcs
		r.KeyValue = MakeApiVersionPrintObject(dcs)             // map header -> rows
		r.Columns = printer.GetHeadersAllDefault(allCols, cols) // headers
	}
	return r
}
