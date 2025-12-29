package commands

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultTemplateCols = []string{"TemplateId", "Name", "Cores", "RAM", "StorageSize", "GPUs"}
)

func TemplateCmd() *core.Command {
	ctx := context.TODO()
	templateCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "template",
			Aliases:          []string{"tpl"},
			Short:            "Template Operations",
			Long:             "The sub-commands of `ionosctl template` allow you to see information about the Templates available.",
			TraverseChildren: true,
		},
	}
	globalFlags := templateCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultTemplateCols, tabheaders.ColsMessage(defaultTemplateCols))
	_ = viper.BindPFlag(core.GetFlagName(templateCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = templateCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultTemplateCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	_ = core.NewCommand(ctx, templateCmd, core.CommandBuilder{
		Namespace:  "template",
		Resource:   "template",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Templates",
		LongDesc:   "Use this command to get a list of available public Templates.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.TemplatesFiltersUsage(),
		Example:    listTemplateExample,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunTemplateList,
		InitClient: true,
	})

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, templateCmd, core.CommandBuilder{
		Namespace:  "template",
		Resource:   "template",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a specified Template",
		LongDesc:   "Use this command to get information about a specified Template.\n\nRequired values to run command:\n\n* Template Id",
		Example:    getTemplateExample,
		PreCmdRun:  PreRunTemplateId,
		CmdRun:     RunTemplateGet,
		InitClient: true,
	})
	get.AddUUIDFlag(cloudapiv6.ArgTemplateId, cloudapiv6.ArgIdShort, "", cloudapiv6.TemplateId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgTemplateId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TemplatesIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return core.WithConfigOverride(templateCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}

func PreRunTemplateId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgTemplateId)
}

func RunTemplateList(c *core.CommandConfig) error {

	templates, resp, err := c.CloudApiV6Services.Templates().List()
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	templatesConverted, err := resource2table.ConvertTemplatesToTable(templates.Templates)
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutputPreconverted(templates, templatesConverted,
		tabheaders.GetHeadersAllDefault(defaultTemplateCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}

func RunTemplateGet(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Template with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTemplateId))))

	tpl, resp, err := c.CloudApiV6Services.Templates().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTemplateId)))
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	templateConverted, err := resource2table.ConvertTemplateToTable(tpl.Template)
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutputPreconverted(tpl, templateConverted,
		tabheaders.GetHeadersAllDefault(defaultTemplateCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}
