package apiversion

import (
	"context"
	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go-dbaas-mongo"
)

func ApiVersionCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dbaas-mongo",
		Verb:      "api-versions",
		Aliases:   []string{"ls"},
		ShortDesc: "Get Mongo API swagger files",
		Example:   "ionosctl dbaas mongo api-versions",
		PreCmdRun: core.NoPreRun,
		CmdRun: func(c *core.CommandConfig) error {
			list, _, err := c.DbaasMongoServices.ApiMetadata().List()
			if err != nil {
				return err
			}
			return c.Printer.Print(getApiVersionPrint(c, list))
		},
		InitClient: true,
	})

	return cmd
}

type ApiVersionPrint struct {
	Version string `json:"Name,omitempty"`
	Href    string `json:"Href,omitempty"`
}

var allCols = structs.Names(ApiVersionPrint{})

func MakeApiVersionPrintObject(objs []ionoscloud.APIVersion) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(objs))

	for _, o := range objs {
		var printObj ApiVersionPrint
		printObj.Version = *o.GetName()
		printObj.Href = *o.GetSwaggerUrl()
	}

	return out
}

func getApiVersionPrint(c *core.CommandConfig, dcs []ionoscloud.APIVersion) printer.Result {
	r := printer.Result{}
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	if c != nil && dcs != nil {
		r.OutputJSON = dcs
		r.KeyValue = MakeApiVersionPrintObject(dcs)             // map header -> rows
		r.Columns = printer.GetHeadersAllDefault(allCols, cols) // headers
	}
	return r
}
