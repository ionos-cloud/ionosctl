package template

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo-v2/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/viper"
)

func TemplateGetCmd() *core.Command {
	ctx := context.TODO()
	get := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace:  "dbaas-mongo-v2",
		Resource:   "template",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a MongoDB Instance Template",
		Example:    "ionosctl dbaas mongo-v2 template get --template-id <template-id>",
		LongDesc:   "Use this command to retrieve details about a MongoDB Instance Template by using its ID.\n\nRequired values to run command:\n\n* Template Id",
		PreCmdRun:  PreRunTemplateId,
		CmdRun:     RunTemplateGet,
		InitClient: true,
	})

	get.AddUUIDFlag(constants.FlagTemplateId, constants.FlagIdShort, "", "The unique ID of the Instance Template", core.RequiredFlagOption(),
		core.WithCompletion(func() []string { return completer.TemplateIds("") }, constants.MongoDBApiRegionalURL, constants.MongoDBLocations),
	)
	return get
}

func PreRunTemplateId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagTemplateId)
}

func RunTemplateGet(c *core.CommandConfig) error {
	templateId := viper.GetString(core.GetFlagName(c.NS, constants.FlagTemplateId))
	c.Verbose("Getting Instance Template %v...", templateId)

	template, _, err := client.Must().MongoClientV2.InstanceTemplatesApi.InstancetemplatesFindById(
		context.Background(), templateId).Execute()
	if err != nil {
		return fmt.Errorf("could not get instance template: %w", err)
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	return c.Out(table.Sprint(templateCols, template, cols))
}
