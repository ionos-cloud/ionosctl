package apiversion

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
)

var allCols = []table.Column{
	{Name: "Version", JSONPath: "name", Default: true},
	{Name: "Href", JSONPath: "swaggerUrl", Default: true},
}

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

			return c.Out(table.Sprint(allCols, list, nil))
		},
		InitClient: true,
	})

	return cmd
}
