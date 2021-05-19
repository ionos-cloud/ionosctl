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
			Short:            "IP Failover Operations",
			Long:             `The sub-command of ` + "`" + `ionosctl ipfailover` + "`" + ` allows you to see information about IP Failovers groups available on a LAN, to add/remove IP Failover group from a LAN.`,
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
		List Command
	*/
	listCmd := core.NewCommand(ctx, ipfailoverCmd, core.CommandBuilder{
		Namespace:  "ipfailover",
		Resource:   "ipfailover",
		Verb:       "list",
		ShortDesc:  "List IP Failovers groups from a LAN",
		LongDesc:   "Use this command to get a list of IP Failovers groups from a LAN.\n\nRequired values to run command:\n\n* Data Center Id\n* Lan Id",
		Example:    listIpFailoverExample,
		PreCmdRun:  PreRunDcLanIds,
		CmdRun:     RunIpFailoverList,
		InitClient: true,
	})
	listCmd.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = listCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	listCmd.AddStringFlag(config.ArgLanId, "", "", config.RequiredFlagLanId)
	_ = listCmd.Command.RegisterFlagCompletionFunc(config.ArgLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLansIds(os.Stderr, viper.GetString(core.GetFlagName(listCmd.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Add Command
	*/
	addCmd := core.NewCommand(ctx, ipfailoverCmd, core.CommandBuilder{
		Namespace: "ipfailover",
		Resource:  "ipfailover",
		Verb:      "add",
		ShortDesc: "Add IP Failover group to a LAN",
		LongDesc: `Use this command to add an IP Failover group to a LAN. 

Successfully setting up an IP Failover group requires three steps:

* Add a reserved IP address to a NIC that will become the IP Failover master.
* Use ` + "`" + `ionosctl ipfailover add` + "`" + ` command to enable IP Failover by providing the relevant IP and NIC Id values.
* Add the same reserved IP address to any other NICs that are a member of the same LAN. Those NICs will become IP Failover members.

Required values to run command:

* Data Center Id
* Lan Id
* Server Id
* Nic Id
* IP address`,
		Example:    addIpFailoverExample,
		PreCmdRun:  PreRunDcLanServerNicIdsIp,
		CmdRun:     RunIpFailoverAdd,
		InitClient: true,
	})
	addCmd.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = addCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	addCmd.AddStringFlag(config.ArgLanId, "", "", config.RequiredFlagLanId)
	_ = addCmd.Command.RegisterFlagCompletionFunc(config.ArgLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLansIds(os.Stderr, viper.GetString(core.GetFlagName(addCmd.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	addCmd.AddStringFlag(config.ArgServerId, "", "", config.RequiredFlagServerId)
	_ = addCmd.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetFlagName(addCmd.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	addCmd.AddStringFlag(config.ArgNicId, "", "", config.RequiredFlagNicId)
	_ = addCmd.Command.RegisterFlagCompletionFunc(config.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNicsIds(os.Stderr, viper.GetString(core.GetFlagName(addCmd.NS, config.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(addCmd.NS, config.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})
	addCmd.AddStringFlag(config.ArgIp, "", "", "IP address to be added to IP Failover Group "+config.RequiredFlag)
	addCmd.AddBoolFlag(config.ArgWaitForRequest, "", config.DefaultWait, "Wait for the Request for IpBlock creation to be executed")
	addCmd.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Request for IpBlock creation [seconds]")

	/*
		Remove Command
	*/
	removeCmd := core.NewCommand(ctx, ipfailoverCmd, core.CommandBuilder{
		Namespace: "ipfailover",
		Resource:  "ipfailover",
		Verb:      "remove",
		ShortDesc: "Remove IP Failover group from a LAN",
		LongDesc: `Use this command to remove an IP Failover group from a LAN.

Required values to run command:

* Data Center Id
* Lan Id
* Server Id
* Nic Id
* IP address`,
		Example:    removeIpFailoverExample,
		PreCmdRun:  PreRunDcLanServerNicIdsIp,
		CmdRun:     RunIpFailoverRemove,
		InitClient: true,
	})
	removeCmd.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = removeCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddStringFlag(config.ArgLanId, "", "", config.RequiredFlagLanId)
	_ = removeCmd.Command.RegisterFlagCompletionFunc(config.ArgLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLansIds(os.Stderr, viper.GetString(core.GetGlobalFlagName(removeCmd.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddStringFlag(config.ArgServerId, "", "", config.RequiredFlagServerId)
	_ = removeCmd.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetGlobalFlagName(removeCmd.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddStringFlag(config.ArgNicId, "", "", config.RequiredFlagNicId)
	_ = removeCmd.Command.RegisterFlagCompletionFunc(config.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNicsIds(os.Stderr, viper.GetString(core.GetGlobalFlagName(removeCmd.NS, config.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(removeCmd.NS, config.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddStringFlag(config.ArgIp, "", "", "Allocated IP "+config.RequiredFlag)
	removeCmd.AddBoolFlag(config.ArgWaitForRequest, "", config.DefaultWait, "Wait for the Request for IpBlock creation to be executed")
	removeCmd.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Request for IpBlock creation [seconds]")

	return ipfailoverCmd
}

func PreRunDcLanIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.NS, config.ArgDataCenterId, config.ArgLanId)
}

func PreRunDcLanServerNicIdsIp(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.NS, config.ArgDataCenterId, config.ArgLanId, config.ArgServerId, config.ArgNicId, config.ArgIp)
}

func RunIpFailoverList(c *core.CommandConfig) error {
	ipsFailovers := make([]resources.IpFailover, 0)
	obj, _, err := c.Lans().Get(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgLanId)),
	)
	if err != nil {
		return err
	}

	if properties, ok := obj.GetPropertiesOk(); ok && properties != nil {
		if ipFailovers, ok := properties.GetIpFailoverOk(); ok && ipFailovers != nil {
			for _, ip := range *ipFailovers {
				ipsFailovers = append(ipsFailovers, resources.IpFailover{IPFailover: ip})
			}
			return c.Printer.Print(getIpFailoverPrint(nil, c, ipsFailovers))
		} else {
			return errors.New("error getting ipfailovers")
		}
	} else {
		return errors.New("error getting lan properties")
	}
}

func RunIpFailoverAdd(c *core.CommandConfig) error {
	ipsFailovers := make([]resources.IpFailover, 0)
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
		if ipFailovers, ok := properties.GetIpFailoverOk(); ok && ipFailovers != nil {
			for _, ip := range *ipFailovers {
				ipsFailovers = append(ipsFailovers, resources.IpFailover{IPFailover: ip})
			}
			return c.Printer.Print(getIpFailoverPrint(nil, c, ipsFailovers))
		} else {
			return errors.New("error getting ipfailovers")
		}
	} else {
		return errors.New("error getting lan properties")
	}
}

func RunIpFailoverRemove(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "remove ip failover group from lan"); err != nil {
		return err
	}
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
			return c.Printer.Print(getIpFailoverPrint(resp, c, nil))
		} else {
			return errors.New("error getting ipfailovers")
		}
	} else {
		return errors.New("error getting lan properties")
	}
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

var defaultIpFailoverCols = []string{"NicId", "Ip"}

type IpFailoverPrint struct {
	NicId string `json:"NicId,omitempty"`
	Ip    string `json:"Ip,omitempty"`
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

func getIpFailoverCols(flagName string, outErr io.Writer) []string {
	if viper.IsSet(flagName) {
		var groupCols []string
		columnsMap := map[string]string{
			"NicId": "NicId",
			"Ip":    "Ip",
		}
		for _, k := range viper.GetStringSlice(flagName) {
			col := columnsMap[k]
			if col != "" {
				groupCols = append(groupCols, col)
			} else {
				clierror.CheckError(errors.New("unknown column "+k), outErr)
			}
		}
		return groupCols
	} else {
		return defaultIpFailoverCols
	}
}

func getIpFailoverKVMaps(ls []resources.IpFailover) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(ls))
	for _, l := range ls {
		var ipPrint IpFailoverPrint
		if nicId, ok := l.GetNicUuidOk(); ok && nicId != nil {
			ipPrint.NicId = *nicId
		}
		if ip, ok := l.GetIpOk(); ok && ip != nil {
			ipPrint.Ip = *ip
		}
		o := structs.Map(ipPrint)
		out = append(out, o)
	}
	return out
}
