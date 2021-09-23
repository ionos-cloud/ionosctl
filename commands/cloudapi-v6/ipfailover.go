package commands

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/waiter"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/internal/utils"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func IpfailoverCmd() *core.Command {
	ctx := context.TODO()
	ipfailoverCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "ipfailover",
			Aliases:          []string{"ipf"},
			Short:            "IP Failover Operations",
			Long:             "The sub-command of `ionosctl ipfailover` allows you to see information about IP Failovers groups available on a LAN, to add/remove IP Failover group from a LAN.",
			TraverseChildren: true,
		},
	}
	globalFlags := ipfailoverCmd.GlobalFlags()
	globalFlags.StringSliceP(config.ArgCols, "", defaultIpFailoverCols, printer.ColsMessage(defaultIpFailoverCols))
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
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List IP Failovers groups from a LAN",
		LongDesc:   "Use this command to get a list of IP Failovers groups from a LAN.\n\nRequired values to run command:\n\n* Data Center Id\n* Lan Id",
		Example:    listIpFailoverExample,
		PreCmdRun:  PreRunDcLanIds,
		CmdRun:     RunIpFailoverList,
		InitClient: true,
	})
	listCmd.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = listCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	listCmd.AddStringFlag(cloudapiv6.ArgLanId, "", "", cloudapiv6.LanId, core.RequiredFlagOption())
	_ = listCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LansIds(os.Stderr, viper.GetString(core.GetFlagName(listCmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Add Command
	*/
	addCmd := core.NewCommand(ctx, ipfailoverCmd, core.CommandBuilder{
		Namespace: "ipfailover",
		Resource:  "ipfailover",
		Verb:      "add",
		Aliases:   []string{"a"},
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
	addCmd.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = addCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	addCmd.AddStringFlag(cloudapiv6.ArgLanId, "", "", cloudapiv6.LanId, core.RequiredFlagOption())
	_ = addCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LansIds(os.Stderr, viper.GetString(core.GetFlagName(addCmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	addCmd.AddStringFlag(cloudapiv6.ArgServerId, "", "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = addCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(addCmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	addCmd.AddStringFlag(cloudapiv6.ArgNicId, "", "", cloudapiv6.NicId, core.RequiredFlagOption())
	_ = addCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NicsIds(os.Stderr, viper.GetString(core.GetFlagName(addCmd.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(addCmd.NS, cloudapiv6.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})
	addCmd.AddStringFlag(cloudapiv6.ArgIp, "", "", "IP address to be added to IP Failover Group", core.RequiredFlagOption())
	addCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for IP Failover creation to be executed")
	addCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for IP Failover creation [seconds]")

	/*
		Remove Command
	*/
	removeCmd := core.NewCommand(ctx, ipfailoverCmd, core.CommandBuilder{
		Namespace: "ipfailover",
		Resource:  "ipfailover",
		Verb:      "remove",
		Aliases:   []string{"r"},
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
	removeCmd.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = removeCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddStringFlag(cloudapiv6.ArgLanId, "", "", cloudapiv6.LanId, core.RequiredFlagOption())
	_ = removeCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LansIds(os.Stderr, viper.GetString(core.GetGlobalFlagName(removeCmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddStringFlag(cloudapiv6.ArgServerId, "", "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = removeCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetGlobalFlagName(removeCmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddStringFlag(cloudapiv6.ArgNicId, "", "", cloudapiv6.NicId, core.RequiredFlagOption())
	_ = removeCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NicsIds(os.Stderr, viper.GetString(core.GetGlobalFlagName(removeCmd.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(removeCmd.NS, cloudapiv6.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddStringFlag(cloudapiv6.ArgIp, "", "", "Allocated IP", core.RequiredFlagOption())
	removeCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for IP Failover deletion to be executed")
	removeCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for IP Failover deletion [seconds]")

	return ipfailoverCmd
}

func PreRunDcLanIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgLanId)
}

func PreRunDcLanServerNicIdsIp(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgLanId, cloudapiv6.ArgServerId, cloudapiv6.ArgNicId, cloudapiv6.ArgIp)
}

func RunIpFailoverList(c *core.CommandConfig) error {
	ipsFailovers := make([]resources.IpFailover, 0)
	obj, resp, err := c.CloudApiV6Services.Lans().Get(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLanId)),
	)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
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
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	lanId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLanId))
	c.Printer.Verbose("Adding an IP Failover group to LAN with ID: %v from Datacenter with ID: %v...", lanId, dcId)
	ipsFailovers := make([]resources.IpFailover, 0)
	lanUpdated, resp, err := c.CloudApiV6Services.Lans().Update(dcId, lanId, getIpFailoverInfo(c))
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	if properties, ok := lanUpdated.GetPropertiesOk(); ok && properties != nil {
		if ipFailovers, ok := properties.GetIpFailoverOk(); ok && ipFailovers != nil {
			for _, ip := range *ipFailovers {
				ipsFailovers = append(ipsFailovers, resources.IpFailover{IPFailover: ip})
			}
			return c.Printer.Print(getIpFailoverPrint(nil, c, ipsFailovers))
		} else {
			return errors.New("error getting updated ipfailovers")
		}
	} else {
		return errors.New("error getting updated lan properties")
	}
}

func RunIpFailoverRemove(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	lanId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLanId))
	c.Printer.Verbose("Removing IP Failover group from LAN with ID: %v from Datacenter with ID: %v...", lanId, dcId)
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "remove ip failover group from lan"); err != nil {
		return err
	}
	oldLan, _, err := c.CloudApiV6Services.Lans().Get(dcId, lanId)
	if err != nil {
		return err
	}
	if properties, ok := oldLan.GetPropertiesOk(); ok && properties != nil {
		if ipfailovers, ok := properties.GetIpFailoverOk(); ok && ipfailovers != nil {
			_, resp, err := c.CloudApiV6Services.Lans().Update(
				viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
				viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLanId)),
				removeIpFailoverInfo(c, ipfailovers),
			)
			if resp != nil {
				c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
			}
			if err != nil {
				return err
			}

			if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
				return err
			}
			return c.Printer.Print(getIpFailoverPrint(resp, c, nil))
		} else {
			return errors.New("error getting ipfailovers to update")
		}
	} else {
		return errors.New("error getting lan properties to update")
	}
}

func getIpFailoverInfo(c *core.CommandConfig) resources.LanProperties {
	ip := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIp))
	nicId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNicId))
	c.Printer.Verbose("Adding IpFailover with Ip: %v and NicUuid: %v", ip, nicId)
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
	removeIp := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIp))
	removeNicId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNicId))
	c.Printer.Verbose("Removing IpFailover with Ip: %v and NicUuid: %v", removeIp, removeNicId)
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
