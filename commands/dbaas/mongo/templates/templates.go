package templates

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	ionoscloud "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TemplatesCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "templates",
			Aliases:          []string{"t"},
			Short:            "PostgreSQL Cluster Operations",
			Long:             "The sub-commands of `ionosctl dbaas postgres cluster` allow you to manage the PostgreSQL Clusters under your account.",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(TemplatesListCmd())
	return cmd
}

// TODO: Why is this tightly coupled to resources.ClusterResponse? Should just take Headers and Columns as params. should also be moved to printer package, to reduce duplication
//
// this is a nightmare to maintain if it is tightly coupled to every single resource!!!!!!!!!!!!
func getTemplatesPrint(resp *ionoscloud.APIResponse, c *core.CommandConfig, ls *[]ionoscloud.TemplateResponse) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForState = viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForState)) // this boolean is duplicated everywhere just to do an append of `& wait` to a verbose message
		}
		if ls != nil {
			r.OutputJSON = ls
			r.KeyValue = getClusterRows(ls)                                                             // map header -> rows
			r.Columns = getClusterHeaders(viper.GetStringSlice(core.GetFlagName(c.NS, config.ArgCols))) // headers
		}
	}
	return r
}

var allCols = structs.Names(TemplatePrint{})

type TemplatePrint struct {
	TemplateId  string `json:"TemplateId,omitempty"`
	Cores       int32  `json:"Cores,omitempty"`
	StorageSize string `json:"StorageSize,omitempty"`
	Ram         string `json:"Ram,omitempty"`
}

func getClusterRows(ls *[]ionoscloud.TemplateResponse) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(*ls))
	for _, t := range *ls {
		var cols TemplatePrint
		if idOk, ok := t.GetIdOk(); ok && idOk != nil {
			cols.TemplateId = *idOk
		}
		if coresOk, ok := t.GetCoresOk(); ok && coresOk != nil {
			cols.Cores = *coresOk
		}
		if ramOk, ok := t.GetRamOk(); ok && ramOk != nil {
			gb, _ := utils.ConvertToGB(fmt.Sprintf("%d", *ramOk), utils.MegaBytes)
			cols.Ram = fmt.Sprintf("%d GB", gb)
		}
		if storageSizeOk, ok := t.GetStorageSizeOk(); ok && storageSizeOk != nil {
			cols.StorageSize = fmt.Sprintf("%d GB", *storageSizeOk)
		}
		o := structs.Map(cols)
		out = append(out, o)
	}
	return out
}

func getClusterHeaders(customColumns []string) []string {
	if customColumns == nil {
		return allCols[0:6]
	}
	//for _, c := customColumns {
	//	if slices.Contains(allCols, c) {}
	//}
	return customColumns
}