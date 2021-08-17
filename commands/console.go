package commands

import (
	"context"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources/v6"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func serverConsole() *core.Command {
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
	get.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(config.ArgServerId, config.ArgIdShort, "", config.RequiredFlagServerId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})

	return consoleCmd
}

func RunServerConsoleGet(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, config.ArgServerId))
	c.Printer.Verbose("ServerConsole with id: %v from Datacenter with id: %v is getting...", serverId, dcId)
	t, _, err := c.Servers().GetRemoteConsoleUrl(
		dcId,
		serverId,
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getConsolePrint(c, []v6.RemoteConsoleUrl{t}))
}

// Output Printing

var defaultConsoleCols = []string{"RemoteConsoleUrl"}

type RemoteConsolePrint struct {
	RemoteConsoleUrl string `json:"RemoteConsoleUrl,omitempty"`
}

func getConsolePrint(c *core.CommandConfig, ss []v6.RemoteConsoleUrl) printer.Result {
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

func getConsoleKVMaps(ss []v6.RemoteConsoleUrl) []map[string]interface{} {
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
