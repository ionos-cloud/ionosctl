package apiversion

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
)

var (
	allJSONPaths = map[string]string{
		"Version": "name",
		"Href":    "swaggerUrl",
	}

	allCols = []string{"Version", "Href"}
)

func ApiVersionCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dbaas-mongo",
		Verb:      "api-versions",
		Aliases:   []string{"versions", "api-version"},
		ShortDesc: "Get Mongo API swagger files",
		Example:   "ionosctl dbaas mongo api-versions",
		PreCmdRun: core.NoPreRun,
		CmdRun: func(c *core.CommandConfig) error {
			list, _, err := client.Must().MongoClient.MetadataApi.InfosVersionGet(context.Background()).Execute()
			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			out, err := jsontabwriter.GenerateOutput("", allJSONPaths, list,
				printer.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Stdout, out)
			return nil
		},
		InitClient: true,
	})

	return cmd
}
