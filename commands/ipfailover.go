package commands

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ipfailover() *core.Command {
	ctx := context.TODO()
	ipfailoverCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "ipfailover",
			Short:            "ipfailover Operations",
			Long:             `The sub-command of ` + "`" + `ionosctl ipfailover` + "`" + ` allows you to see information about ipfailovers available to create objects.`,
			TraverseChildren: true,
		},
	}
	globalFlags := ipfailoverCmd.GlobalFlags()
	globalFlags.StringSlice(config.ArgCols, defaultIpFailoverCols, "Columns to be printed in the standard output")
	_ = viper.BindPFlag(core.GetGlobalFlagName(ipfailoverCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = ipfailoverCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultIpFailoverCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Add Command
	*/
	addCmd := core.NewCommand(ctx, ipfailoverCmd, core.CommandBuilder{
		Namespace:  "ipfailover",
		Resource:   "ipfailover",
		Verb:       "add",
		ShortDesc:  "Add IPFailover to a LAN",
		LongDesc:   "Use this command to get a list of available ipfailovers to create objects on.",
		Example:    "",
		PreCmdRun:  noPreRun,
		CmdRun:     RunIpFailoverAdd,
		InitClient: true,
	})
	addCmd.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagLanId)
	_ = addCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	addCmd.AddStringFlag(config.ArgLanId, "", "", config.RequiredFlagLanId)
	_ = addCmd.Command.RegisterFlagCompletionFunc(config.ArgLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLansIds(os.Stderr, viper.GetString(core.GetFlagName(addCmd.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	addCmd.AddStringFlag(config.ArgServerId, "", "", config.RequiredFlagNicId)
	_ = addCmd.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetFlagName(addCmd.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	addCmd.AddStringFlag(config.ArgNicId, "", "", config.RequiredFlagNicId)
	_ = addCmd.Command.RegisterFlagCompletionFunc(config.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNicsIds(os.Stderr, viper.GetString(core.GetFlagName(addCmd.NS, config.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(addCmd.NS, config.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})
	addCmd.AddStringFlag(config.ArgIp, "", "", config.RequiredFlagNicId)

	/*
		Remove Command
	*/
	removeCmd := core.NewCommand(ctx, ipfailoverCmd, core.CommandBuilder{
		Namespace:  "ipfailover",
		Resource:   "ipfailover",
		Verb:       "remove",
		ShortDesc:  "Add IPFailover to a LAN",
		LongDesc:   "Use this command to get a list of available ipfailovers to create objects on.",
		Example:    "",
		PreCmdRun:  noPreRun,
		CmdRun:     RunIpFailoverRemove,
		InitClient: true,
	})
	removeCmd.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagLanId)
	_ = removeCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddStringFlag(config.ArgLanId, "", "", config.RequiredFlagLanId)
	_ = removeCmd.Command.RegisterFlagCompletionFunc(config.ArgLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLansIds(os.Stderr, viper.GetString(core.GetGlobalFlagName(removeCmd.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddStringFlag(config.ArgServerId, "", "", config.RequiredFlagNicId)
	_ = removeCmd.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetGlobalFlagName(removeCmd.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddStringFlag(config.ArgNicId, "", "", config.RequiredFlagNicId)
	_ = removeCmd.Command.RegisterFlagCompletionFunc(config.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNicsIds(os.Stderr, viper.GetString(core.GetGlobalFlagName(removeCmd.NS, config.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(removeCmd.NS, config.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddStringFlag(config.ArgIp, "", "", config.RequiredFlagNicId)

	return ipfailoverCmd
}

func RunIpFailoverAdd(c *core.CommandConfig) error {
	lanUpdated, resp, err := c.Lans().Update(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgLanId)),
		getIpFailoverInfo(c),
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	if properties, ok := lanUpdated.GetPropertiesOk(); ok && properties != nil {
		if ipfailovers, ok := properties.GetIpFailoverOk(); ok && ipfailovers != nil {
			return c.Printer.Print(getIpFailoverPrint(nil, c, getIpFailovers(ipfailovers)))
		}
	}
	return nil
}

func RunIpFailoverRemove(c *core.CommandConfig) error {
	oldLan, _, err := c.Lans().Get(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgLanId)),
	)
	if err != nil {
		return err
	}
	if properties, ok := oldLan.GetPropertiesOk(); ok && properties != nil {
		if ipfailovers, ok := properties.GetIpFailoverOk(); ok && ipfailovers != nil {
			_, resp, err := c.Lans().Update(
				viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
				viper.GetString(core.GetFlagName(c.NS, config.ArgLanId)),
				removeIpFailoverInfo(c, ipfailovers),
			)
			if err != nil {
				return err
			}

			if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
				return err
			}
		}
	}
	return c.Printer.Print(getIpFailoverPrint(nil, c, nil))
}

func getIpFailoverInfo(c *core.CommandConfig) resources.LanProperties {
	ip := viper.GetString(core.GetFlagName(c.NS, config.ArgIp))
	nicId := viper.GetString(core.GetFlagName(c.NS, config.ArgNicId))
	return resources.LanProperties{
		LanProperties: ionoscloud.LanProperties{
			IpFailover: &[]ionoscloud.IPFailover{
				{
					Ip:      &ip,
					NicUuid: &nicId,
				},
			},
		},
	}
}

func removeIpFailoverInfo(c *core.CommandConfig, failovers *[]ionoscloud.IPFailover) resources.LanProperties {
	removeIp := viper.GetString(core.GetFlagName(c.NS, config.ArgIp))
	removeNicId := viper.GetString(core.GetFlagName(c.NS, config.ArgNicId))

	newIpFailover := make([]ionoscloud.IPFailover, 0)
	for _, failover := range *failovers {
		if ip, ok := failover.GetIpOk(); ok && ip != nil && *ip != removeIp {
			if nicId, ok := failover.GetNicUuidOk(); ok && nicId != nil && *nicId != removeNicId {
				newIpFailover = append(newIpFailover, failover)
			}
		}
	}
	return resources.LanProperties{
		LanProperties: ionoscloud.LanProperties{
			IpFailover: &newIpFailover,
		},
	}
}

// Output Printing

var defaultIpFailoverCols = []string{"NicId", "IP"}

type IpFailoverPrint struct {
	NicId string `json:"NicId,omitempty"`
	IP    string `json:"IP,omitempty"`
}

func getIpFailoverPrint(resp *resources.Response, c *core.CommandConfig, ips []resources.IpFailover) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForRequest = viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForRequest))
		}
		if ips != nil {
			r.OutputJSON = ips
			r.KeyValue = getIpFailoverKVMaps(ips)
			r.Columns = getIpFailoverCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getIpFailovers(ips *[]ionoscloud.IPFailover) []resources.IpFailover {
	ls := make([]resources.IpFailover, 0)
	if ips != nil {
		for _, s := range *ips {
			ls = append(ls, resources.IpFailover{IPFailover: s})
		}
	}
	return ls
}

func getIpFailoverCols(flagName string, outErr io.Writer) []string {
	var cols []string
	if viper.IsSet(flagName) {
		cols = viper.GetStringSlice(flagName)
	} else {
		return defaultIpFailoverCols
	}

	columnsMap := map[string]string{
		"NicId": "NicId",
		"IP":    "IP",
	}
	var lanCols []string
	for _, k := range cols {
		col := columnsMap[k]
		if col != "" {
			lanCols = append(lanCols, col)
		} else {
			clierror.CheckError(errors.New("unknown column "+k), outErr)
		}
	}
	return lanCols
}

func getIpFailoverKVMaps(ls []resources.IpFailover) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(ls))
	for _, l := range ls {
		var ipPrint IpFailoverPrint
		if nicId, ok := l.GetNicUuidOk(); ok && nicId != nil {
			ipPrint.NicId = *nicId
		}
		if ip, ok := l.GetIpOk(); ok && ip != nil {
			ipPrint.IP = *ip
		}
		o := structs.Map(ipPrint)
		out = append(out, o)
	}
	return out
}
