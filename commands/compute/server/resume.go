package server

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ServerResumeCmd() *core.Command {
	resume := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "server",
		Resource:  "server",
		Verb:      "resume",
		Aliases:   []string{"res"},
		ShortDesc: "Resume a Cube Server",
		LongDesc: `Use this command to resume a Cube Server. The operation can only be applied to suspended Cube Servers.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the resource to reach AVAILABLE state. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* Server Id`,
		Example:    "ionosctl compute server resume --datacenter-id DATACENTER_ID --server-id SERVER_ID",
		PreCmdRun:  PreRunDcServerIds,
		CmdRun:     RunServerResume,
		InitClient: true,
	})
	resume.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = resume.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	resume.AddUUIDFlag(cloudapiv6.ArgServerId, cloudapiv6.ArgIdShort, "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = resume.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIdsCustom(viper.GetString(core.GetFlagName(resume.NS, cloudapiv6.ArgDataCenterId)),
			resources.ListQueryParams{
				Filters: &map[string][]string{
					"type": {"CUBE"},
				},
			}), cobra.ShellCompDirectiveNoFileComp
	})

	return resume
}
