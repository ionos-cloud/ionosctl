package apiversion

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
)

var (
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

			out, err := jsontabwriter.GenerateOutput("", jsonpaths.DbaasMongoAPIVersion, list,
				tabheaders.GetHeadersAllDefault(allCols, nil))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
			return nil
		},
		InitClient: true,
	})

	return cmd
}
