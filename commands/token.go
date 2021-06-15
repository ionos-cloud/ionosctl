package commands

import (
	"context"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func serverToken() *core.Command {
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
	get.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(config.ArgServerId, config.ArgIdShort, "", config.RequiredFlagServerId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})

	return tokenCmd
}

func RunServerTokenGet(c *core.CommandConfig) error {
	t, _, err := c.Servers().GetToken(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgServerId)),
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
