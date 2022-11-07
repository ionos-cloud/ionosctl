package logs

import (
	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ApiVersionCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "apiversion",
			Short:            "Mongo API Metadata Operations",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(ApiVersionCmd())

	return cmd
}

type ApiVersionPrint struct {
	Name    string `json:"Name,omitempty"`
	Href    string `json:"Href,omitempty"`
	Version string `json:"Version,omitempty"`
}

var allCols = structs.Names(ApiVersionPrint{})

func MakeLogsPrintObject(logs *[]ionoscloud.ClusterLogsInstances) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(*logs))
	for idx, instance := range *logs {
		for msgIdx, msg := range *instance.GetMessages() {
			var logsPrint LogsPrint

			logsPrint.InstanceNumber = idx
			logsPrint.MessageNumber = msgIdx
			logsPrint.Name = *instance.GetName()
			logsPrint.Message = *msg.GetMessage()
			logsPrint.Time = *msg.GetTime()

			o := structs.Map(logsPrint)
			out = append(out, o)
		}
	}

	return out
}

func getLogsPrint(c *core.CommandConfig, dcs *[]ionoscloud.ClusterLogsInstances) printer.Result {
	r := printer.Result{}
	if c != nil && dcs != nil {
		r.OutputJSON = dcs
		r.KeyValue = MakeLogsPrintObject(dcs)                                                                                                 // map header -> rows
		r.Columns = printer.GetHeadersAllDefault(structs.Names(LogsPrint{}), viper.GetStringSlice(core.GetFlagName(c.NS, constants.ArgCols))) // headers
	}
	return r
}
