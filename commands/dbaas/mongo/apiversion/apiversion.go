package apiversion

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/tabheaders"
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

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
			return nil
		},
		InitClient: true,
	})

	return cmd
}
