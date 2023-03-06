package commands

import (
	"context"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddUUIDFlag(cloudapiv6.ArgServerId, cloudapiv6.ArgIdShort, "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)
	get.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.ArgDepthDescription)
	return tokenCmd
}

func RunServerTokenGet(c *core.CommandConfig) error {
	c.Printer.Verbose("ServerToken with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)))
	t, _, err := c.CloudApiV6Services.Servers().GetToken(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getTokenPrint(c, []resources.Token{t}))
}

// Output Printing

var defaultTokenCols = []string{"Token"}

type TokenPrint struct {
	Token string `json:"Token,omitempty"`
}

func getTokenPrint(c *core.CommandConfig, ss []resources.Token) printer.Result {
	r := printer.Result{}
	if c != nil {
		if ss != nil {
			r.OutputJSON = ss
			r.KeyValue = getTokenKVMaps(ss)
			r.Columns = defaultTokenCols
		}
	}
	return r
}

func getTokenKVMaps(ss []resources.Token) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(ss))
	for _, s := range ss {
		var tokenPrint TokenPrint
		if t, ok := s.GetTokenOk(); ok && t != nil {
			tokenPrint.Token = *t
		}
		o := structs.Map(tokenPrint)
		out = append(out, o)
	}
	return out
}
