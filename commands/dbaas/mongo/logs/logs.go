package logs

import (
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/json2table"
	ionoscloud "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"github.com/spf13/cobra"
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

var (
	allCols     = []string{"Instance", "Name", "MessageNumber", "Message", "Time"}
	defaultCols = []string{"Instance", "Name", "MessageNumber", "Time"}
)

func convertLogsToTable(logs *[]ionoscloud.ClusterLogsInstances) ([]map[string]interface{}, error) {
	if logs == nil {
		return nil, fmt.Errorf("no logs to process")
	}

	out := make([]map[string]interface{}, 0, len(*logs))
	for idx, instance := range *logs {
		if instance.GetMessages() == nil {
			continue
		}
		for msgIdx, msg := range *instance.GetMessages() {
			o, err := json2table.ConvertJSONToTable("", jsonpaths.LogsMessage, msg)
			if err != nil {
				return nil, fmt.Errorf("could not convert from JSON to Table format: %w", err)
			}

			o[0]["Instance"] = idx
			o[0]["MessageNumber"] = msgIdx
			if instance.GetName() != nil {
				o[0]["Name"] = *instance.GetName()
			}

			out = append(out, o...)
		}
	}

	return out, nil
}
