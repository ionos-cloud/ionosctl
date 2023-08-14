package logs

import (
	"time"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func LogsCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "logs",
			Short:            "Mongo Logs Operations",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(LogsListCmd())

	return cmd
}

type LogsPrint struct {
	Instance      int       `json:"Instance,omitempty"`
	Name          string    `json:"Name,omitempty"`
	MessageNumber int       `json:"MessageNumber,omitempty"`
	Message       string    `json:"Message,omitempty"`
	Time          time.Time `json:"Time,omitempty"`
}

var allCols = structs.Names(LogsPrint{})
var defaultCols = []string{"Instance", "Name", "MessageNumber", "Time"}

func MakeLogsPrintObject(logs *[]ionoscloud.ClusterLogsInstances) []map[string]interface{} {
	if logs == nil {
		return nil
	}

	out := make([]map[string]interface{}, 0, len(*logs))
	for idx, instance := range *logs {
		if instance.GetMessages() == nil {
			continue
		}
		for msgIdx, msg := range *instance.GetMessages() {
			var logsPrint LogsPrint

			logsPrint.Instance = idx
			logsPrint.MessageNumber = msgIdx
			if instance.GetName() != nil {
				logsPrint.Name = *instance.GetName()
			}
			if msg.GetMessage() != nil {
				logsPrint.Message = *msg.GetMessage()
			}
			if msg.GetTime() != nil {
				logsPrint.Time = *msg.GetTime()
			}

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
		r.KeyValue = MakeLogsPrintObject(dcs)                                                                                 // map header -> rows
		r.Columns = printer.GetHeaders(allCols, defaultCols, viper.GetStringSlice(core.GetFlagName(c.NS, constants.ArgCols))) // headers
	}
	return r
}
