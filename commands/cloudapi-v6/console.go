package commands

import (
	"context"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/completer"
	cloudapi_v6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ServerConsoleCmd() *core.Command {
	ctx := context.TODO()
	consoleCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "console",
			Aliases:          []string{"url"},
			Short:            "Server Remote Console URL Operations",
			Long:             `The sub-command of ` + "`" + `ionosctl server console` + "`" + ` allows you to get the URL for Remote Console of a specific Server.`,
			TraverseChildren: true,
		},
	}

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, consoleCmd, core.CommandBuilder{
		Namespace:  "server",
		Resource:   "console",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get the Remote Console URL to access a Server",
		LongDesc:   "Use this command to get the Server Remote Console link.\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id",
		Example:    getConsoleServerExample,
		PreCmdRun:  PreRunDcServerIds,
		CmdRun:     RunServerConsoleGet,
		InitClient: true,
	})
	get.AddStringFlag(cloudapi_v6.ArgDataCenterId, "", "", cloudapi_v6.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapi_v6.ArgServerId, cloudapi_v6.ArgIdShort, "", cloudapi_v6.ServerId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})

	return consoleCmd
}

func RunServerConsoleGet(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgServerId))
	c.Printer.Verbose("Getting Consoler URL for Server with ID: %v from Datacenter with ID: %v...", serverId, dcId)
	t, _, err := c.CloudApiV6Services.Servers().GetRemoteConsoleUrl(dcId, serverId)
	if err != nil {
		return err
	}
	return c.Printer.Print(getConsolePrint(c, []resources.RemoteConsoleUrl{t}))
}

// Output Printing

var defaultConsoleCols = []string{"RemoteConsoleUrl"}

type RemoteConsolePrint struct {
	RemoteConsoleUrl string `json:"RemoteConsoleUrl,omitempty"`
}

func getConsolePrint(c *core.CommandConfig, ss []resources.RemoteConsoleUrl) printer.Result {
	r := printer.Result{}
	if c != nil {
		if ss != nil {
			r.OutputJSON = ss
			r.KeyValue = getConsoleKVMaps(ss)
			r.Columns = defaultConsoleCols
		}
	}
	return r
}

func getConsoleKVMaps(ss []resources.RemoteConsoleUrl) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(ss))
	for _, s := range ss {
		var consolePrint RemoteConsolePrint
		if t, ok := s.GetUrlOk(); ok && t != nil {
			consolePrint.RemoteConsoleUrl = *t
		}
		o := structs.Map(consolePrint)
		out = append(out, o)
	}
	return out
}
