package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/query"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	globalFlags.StringSliceP(constants.ArgCols, "", defaultTemplateCols, printer.ColsMessage(defaultTemplateCols))
	_ = viper.BindPFlag(core.GetFlagName(templateCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = templateCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultTemplateCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, templateCmd, core.CommandBuilder{
		Namespace:  "template",
		Resource:   "template",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Templates",
		LongDesc:   "Use this command to get a list of available public Templates.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.TemplatesFiltersUsage(),
		Example:    listTemplateExample,
		PreCmdRun:  PreRunTemplateList,
		CmdRun:     RunTemplateList,
		InitClient: true,
	})
	list.AddInt32Flag(cloudapiv6.ArgMaxResults, cloudapiv6.ArgMaxResultsShort, cloudapiv6.DefaultMaxResults, cloudapiv6.ArgMaxResultsDescription)
	list.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultListDepth, cloudapiv6.ArgDepthDescription)
	list.AddStringFlag(cloudapiv6.ArgOrderBy, "", "", cloudapiv6.ArgOrderByDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TemplatesFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv6.ArgFilters, cloudapiv6.ArgFiltersShort, []string{""}, cloudapiv6.ArgFiltersDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TemplatesFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)

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
		return completer.TemplatesIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)
	get.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.ArgDepthDescription)
	return templateCmd
}

func PreRunTemplateList(c *core.PreCommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgFilters)) {
		return query.ValidateFilters(c, completer.TemplatesFilters(), completer.TemplatesFiltersUsage())
	}
	return nil
}

func PreRunTemplateId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgTemplateId)
}

func RunTemplateList(c *core.CommandConfig) error {
	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	templates, resp, err := c.CloudApiV6Services.Templates().List(listQueryParams)
	if resp != nil {
		c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getTemplatePrint(c, getTemplates(templates)))
}

func RunTemplateGet(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	c.Printer.Verbose("Template with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTemplateId)))
	tpl, resp, err := c.CloudApiV6Services.Templates().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTemplateId)), queryParams)
	if resp != nil {
		c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getTemplatePrint(c, getTemplate(tpl)))
}

// Output Printing

var defaultTemplateCols = []string{"TemplateId", "Name", "Cores", "Ram", "StorageSize"}

type TemplatePrint struct {
	TemplateId  string  `json:"TemplateId,omitempty"`
	Name        string  `json:"Name,omitempty"`
	Cores       float32 `json:"Cores,omitempty"`
	Ram         string  `json:"Ram,omitempty"`
	StorageSize string  `json:"StorageSize,omitempty"`
}

func getTemplatePrint(c *core.CommandConfig, tpls []resources.Template) printer.Result {
	r := printer.Result{}
	if c != nil {
		if tpls != nil {
			r.OutputJSON = tpls
			r.KeyValue = getTemplatesKVMaps(tpls)
			r.Columns = printer.GetHeadersAllDefault(defaultTemplateCols, viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgCols)))
		}
	}
	return r
}

func getTemplates(templates resources.Templates) []resources.Template {
	tpls := make([]resources.Template, 0)
	if items, ok := templates.GetItemsOk(); ok && items != nil {
		for _, d := range *items {
			tpls = append(tpls, resources.Template{Template: d})
		}
	}
	return tpls
}

func getTemplate(template *resources.Template) []resources.Template {
	tpls := make([]resources.Template, 0)
	if template != nil {
		tpls = append(tpls, resources.Template{Template: template.Template})
	}
	return tpls
}

func getTemplatesKVMaps(tpls []resources.Template) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(tpls))
	for _, tpl := range tpls {
		o := getTemplateKVMap(tpl)
		out = append(out, o)
	}
	return out
}

func getTemplateKVMap(tpl resources.Template) map[string]interface{} {
	var tplPrint TemplatePrint
	if tplId, ok := tpl.GetIdOk(); ok && tplId != nil {
		tplPrint.TemplateId = *tplId
	}
	if properties, ok := tpl.GetPropertiesOk(); ok && properties != nil {
		if name, ok := properties.GetNameOk(); ok && name != nil {
			tplPrint.Name = *name
		}
		if c, ok := properties.GetCoresOk(); ok && c != nil {
			tplPrint.Cores = *c
		}
		if r, ok := properties.GetRamOk(); ok && r != nil {
			tplPrint.Ram = fmt.Sprintf("%vMB", *r)
		}
		if storageSize, ok := properties.GetStorageSizeOk(); ok && storageSize != nil {
			tplPrint.StorageSize = fmt.Sprintf("%vGB", *storageSize)
		}
	}
	return structs.Map(tplPrint)
}
