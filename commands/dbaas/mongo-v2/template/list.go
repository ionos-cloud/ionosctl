package template

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	mongo "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v3"
	"github.com/spf13/viper"
)

func TemplateListCmd() *core.Command {
	ctx := context.TODO()
	cmd := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace:  "dbaas-mongo-v2",
		Resource:   "template",
		Verb:       "list",
		Aliases:    []string{"ls", "l"},
		ShortDesc:  "List MongoDB Instance Templates",
		LongDesc:   "Use this command to retrieve a list of MongoDB Instance Templates. You can filter by edition using `--edition`.",
		Example:    "ionosctl dbaas mongo-v2 template list --edition BUSINESS",
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunTemplateList,
		InitClient: true,
	})
	cmd.AddSetFlag(constants.FlagEdition, "", "", []string{string(mongo.MONGOEDITION_ENTERPRISE), string(mongo.MONGOEDITION_BUSINESS)},
		"Response filter to list only the templates of this edition")
	cmd.AddInt32Flag(constants.FlagLimit, "", 100, "The maximum number of elements to return")
	cmd.AddInt32Flag(constants.FlagOffset, "", 0, "The first element to return")
	return cmd
}

func RunTemplateList(c *core.CommandConfig) error {
	req := client.Must().MongoClientV2.InstanceTemplatesApi.InstancetemplatesGet(context.Background())
	if fn := core.GetFlagName(c.NS, constants.FlagEdition); viper.IsSet(fn) {
		req = req.FilterEdition(mongo.MongoEdition(viper.GetString(fn)))
	}
	if fn := core.GetFlagName(c.NS, constants.FlagLimit); viper.IsSet(fn) {
		req = req.Limit(viper.GetInt32(fn))
	}
	if fn := core.GetFlagName(c.NS, constants.FlagOffset); viper.IsSet(fn) {
		req = req.Offset(viper.GetInt32(fn))
	}
	templates, _, err := req.Execute()
	if err != nil {
		return err
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	return c.Out(table.Sprint(templateCols, templates, cols, table.WithPrefix("items")))
}
