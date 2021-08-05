package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources/v6"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func template() *core.Command {
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
	globalFlags.StringSliceP(config.ArgCols, "", defaultTemplateCols, utils.ColsMessage(defaultTemplateCols))
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
		PreCmdRun:  noPreRun,
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
	get.AddStringFlag(config.ArgTemplateId, config.ArgIdShort, "", config.RequiredFlagTemplateId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgTemplateId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getTemplatesIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	return templateCmd
}

func PreRunTemplateId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.NS, config.ArgTemplateId)
}

func RunTemplateList(c *core.CommandConfig) error {
	templates, _, err := c.Templates().List()
	if err != nil {
		return err
	}
	return c.Printer.Print(getTemplatePrint(c, getTemplates(templates)))
}

func RunTemplateGet(c *core.CommandConfig) error {
	tpl, _, err := c.Templates().Get(viper.GetString(core.GetFlagName(c.NS, config.ArgTemplateId)))
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

func getTemplatePrint(c *core.CommandConfig, tpls []v6.Template) printer.Result {
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

func getTemplates(templates v6.Templates) []v6.Template {
	tpls := make([]v6.Template, 0)
	if items, ok := templates.GetItemsOk(); ok && items != nil {
		for _, d := range *items {
			tpls = append(tpls, v6.Template{Template: d})
		}
	}
	return tpls
}

func getTemplate(template *v6.Template) []v6.Template {
	tpls := make([]v6.Template, 0)
	if template != nil {
		tpls = append(tpls, v6.Template{Template: template.Template})
	}
	return tpls
}

func getTemplatesKVMaps(tpls []v6.Template) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(tpls))
	for _, tpl := range tpls {
		o := getTemplateKVMap(tpl)
		out = append(out, o)
	}
	return out
}

func getTemplateKVMap(tpl v6.Template) map[string]interface{} {
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

func getTemplatesIds(outErr io.Writer) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := v6.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		config.GetServerUrl(),
	)
	clierror.CheckError(err, outErr)
	tplSvc := v6.NewTemplateService(clientSvc.Get(), context.TODO())
	tpls, _, err := tplSvc.List()
	clierror.CheckError(err, outErr)
	tplsIds := make([]string, 0)
	if items, ok := tpls.Templates.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				tplsIds = append(tplsIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return tplsIds
}
