package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	allConsoleJSONPaths = map[string]string{
		"RemoteConsoleUrl": "url",
	}

	defaultConsoleCols = []string{"RemoteConsoleUrl"}
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
	get.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddUUIDFlag(cloudapiv6.ArgServerId, cloudapiv6.ArgIdShort, "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)
	get.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.ArgDepthDescription)
	return consoleCmd
}

func RunServerConsoleGet(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Getting Consoler URL for Server with ID: %v from Datacenter with ID: %v...", serverId, dcId))

	t, resp, err := c.CloudApiV6Services.Servers().GetRemoteConsoleUrl(dcId, serverId)
	if err != nil {
		return err
	}
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}

	cols, err := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("", allConsoleJSONPaths, t.RemoteConsoleUrl,
		tabheaders.GetHeadersAllDefault(defaultImageCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}
