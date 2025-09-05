package commands

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultTokenCols = []string{"Token"}
)

func ServerTokenCmd() *core.Command {
	ctx := context.TODO()
	tokenCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "token",
			Aliases:          []string{"t"},
			Short:            "Server Token Operations",
			Long:             "The sub-command of `ionosctl server token` allows you to get Token for specific Server.",
			TraverseChildren: true,
		},
	}

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, tokenCmd, core.CommandBuilder{
		Namespace:  "server",
		Resource:   "token",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a Token from a Server",
		LongDesc:   "Use this command to get the Server's jwToken.\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id",
		Example:    getTokenServerExample,
		PreCmdRun:  PreRunDcServerIds,
		CmdRun:     RunServerTokenGet,
		InitClient: true,
	})
	get.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddUUIDFlag(cloudapiv6.ArgServerId, cloudapiv6.ArgIdShort, "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.ArgDepthDescription)
	return core.WithConfigOverride(tokenCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}

func RunServerTokenGet(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"ServerToken with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))))

	t, _, err := c.CloudApiV6Services.Servers().GetToken(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
	)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.Token, t.Token, tabheaders.GetHeadersAllDefault(defaultTokenCols, nil))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}
