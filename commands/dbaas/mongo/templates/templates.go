package templates

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/utils"
	ionoscloud "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strconv"
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
	Name        string `json:"Name,omitempty"`
	Edition     string `json:"Edition,omitempty"`
	Cores       int32  `json:"Cores,omitempty"`
	StorageSize string `json:"StorageSize,omitempty"`
	Ram         string `json:"Ram,omitempty"`
}

var allCols = structs.Names(TemplatePrint{})

func getClusterRows(ls *[]ionoscloud.TemplateResponse) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(*ls))
	for _, t := range *ls {
		var cols TemplatePrint
		cols.TemplateId = *t.Id
		cols.Cores = *t.Properties.Cores
		cols.StorageSize = fmt.Sprintf("%d GB", *t.Properties.StorageSize)
		cols.Name = *t.Properties.Name
		cols.Edition = *t.Properties.Edition
		ramGb, err := utils.ConvertToGB(strconv.Itoa(int(*t.Properties.Ram)), utils.MegaBytes)
		if err == nil {
			cols.Ram = fmt.Sprintf("%d GB", ramGb)
		}
		o := structs.Map(cols)
		out = append(out, o)
	}
	return out
}
