package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/completer"
	cloudapi_v6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
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
	globalFlags.StringSliceP(config.ArgCols, "", defaultTemplateCols, printer.ColsMessage(defaultTemplateCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(templateCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = templateCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultTemplateCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	core.NewCommand(ctx, templateCmd, core.CommandBuilder{
		Namespace:  "template",
		Resource:   "template",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Templates",
		LongDesc:   "Use this command to get a list of available public Templates.",
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
	get.AddStringFlag(cloudapi_v6.ArgTemplateId, cloudapi_v6.ArgIdShort, "", cloudapi_v6.TemplateId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgTemplateId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TemplatesIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	return templateCmd
}

func PreRunTemplateId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapi_v6.ArgTemplateId)
}

func RunTemplateList(c *core.CommandConfig) error {
	templates, _, err := c.CloudApiV6Services.Templates().List()
	if err != nil {
		return err
	}
	return c.Printer.Print(getTemplatePrint(c, getTemplates(templates)))
}

func RunTemplateGet(c *core.CommandConfig) error {
	c.Printer.Verbose("Template with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgTemplateId)))
	tpl, _, err := c.CloudApiV6Services.Templates().Get(viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgTemplateId)))
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
			r.Columns = getTemplateCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getTemplateCols(flagName string, outErr io.Writer) []string {
	var cols []string
	if viper.IsSet(flagName) {
		cols = viper.GetStringSlice(flagName)
	} else {
		return defaultTemplateCols
	}

	columnsMap := map[string]string{
		"TemplateId":  "TemplateId",
		"Name":        "Name",
		"Cores":       "Cores",
		"Ram":         "Ram",
		"StorageSize": "StorageSize",
	}
	var datacenterCols []string
	for _, k := range cols {
		col := columnsMap[k]
		if col != "" {
			datacenterCols = append(datacenterCols, col)
		} else {
			clierror.CheckError(errors.New("unknown column "+k), outErr)
		}
	}
	return datacenterCols
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
