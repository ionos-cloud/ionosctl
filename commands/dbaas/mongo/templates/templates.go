package templates

import (
	"fmt"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
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
			Short:            "Mongo Templates Operations",
			Long:             "Templates can be used to create MongoDB clusters; they contain properties such as number of cores, RAM, and the storage size",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(TemplatesListCmd())
	return cmd
}

func getTemplatesPrint(c *core.CommandConfig, ls *[]ionoscloud.TemplateResponse) printer.Result {
	r := printer.Result{}
	if c != nil && ls != nil {
		r.OutputJSON = ls
		r.KeyValue = getClusterRows(ls)                                                                                    // map header -> rows
		r.Columns = printer.GetHeadersAllDefault(allCols, viper.GetStringSlice(core.GetFlagName(c.NS, constants.ArgCols))) // headers
	}
	return r
}

type TemplatePrint struct {
	TemplateId  string `json:"TemplateId,omitempty"`
	Cores       int32  `json:"Cores,omitempty"`
	StorageSize string `json:"StorageSize,omitempty"`
	Ram         string `json:"Ram,omitempty"`
}

var allCols = structs.Names(TemplatePrint{})

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
