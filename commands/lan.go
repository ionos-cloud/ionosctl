package commands

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources/v5"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func lan() *core.Command {
	ctx := context.TODO()
	lanCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "lan",
			Aliases:          []string{"l"},
			Short:            "LAN Operations",
			Long:             "The sub-commands of `ionosctl lan` allow you to create, list, get, update, delete LANs.",
			TraverseChildren: true,
		},
	}
	globalFlags := lanCmd.GlobalFlags()
	globalFlags.StringSliceP(config.ArgCols, "", defaultLanCols, utils.ColsMessage(defaultLanCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(lanCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = lanCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultLanCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, lanCmd, core.CommandBuilder{
		Namespace:  "lan",
		Resource:   "lan",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List LANs",
		LongDesc:   "Use this command to retrieve a list of LANs provisioned in a specific Virtual Data Center.\n\nRequired values to run command:\n\n* Data Center Id",
		Example:    listLanExample,
		PreCmdRun:  PreRunDcLanIds,
		CmdRun:     RunLanList,
		InitClient: true,
	})
	list.AddStringFlag(config.ArgDataCenterId, "", "", config.DatacenterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, lanCmd, core.CommandBuilder{
		Namespace:  "lan",
		Resource:   "lan",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a LAN",
		LongDesc:   "Use this command to retrieve information of a given LAN.\n\nRequired values to run command:\n\n* Data Center Id\n* LAN Id",
		Example:    getLanExample,
		PreCmdRun:  PreRunDcLanIds,
		CmdRun:     RunLanGet,
		InitClient: true,
	})
	get.AddStringFlag(config.ArgDataCenterId, "", "", config.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(config.ArgLanId, config.ArgIdShort, "", config.LanId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLansIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, lanCmd, core.CommandBuilder{
		Namespace: "lan",
		Resource:  "lan",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a LAN",
		LongDesc: `Use this command to create a new LAN within a Virtual Data Center on your account. The name, the public option and the Private Cross-Connect Id can be set.

NOTE: IP Failover is configured after LAN creation using an update command.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id`,
		Example:    createLanExample,
		PreCmdRun:  PreRunDataCenterId,
		CmdRun:     RunLanCreate,
		InitClient: true,
	})
	create.AddStringFlag(config.ArgDataCenterId, "", "", config.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgName, config.ArgNameShort, "Unnamed LAN", "The name of the LAN")
	create.AddBoolFlag(config.ArgPublic, "", config.DefaultPublic, "Indicates if the LAN faces the public Internet (true) or not (false)")
	create.AddStringFlag(config.ArgPccId, "", "", "The unique Id of the Private Cross-Connect the LAN will connect to")
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgPccId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getPccsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for LAN creation to be executed")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for LAN creation [seconds]")

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, lanCmd, core.CommandBuilder{
		Namespace: "lan",
		Resource:  "lan",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a LAN",
		LongDesc: `Use this command to update a specified LAN. You can update the name, the public option for LAN and the Pcc Id to connect the LAN to a Private Cross-Connect.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* LAN Id`,
		Example:    updateLanExample,
		PreCmdRun:  PreRunDcLanIds,
		CmdRun:     RunLanUpdate,
		InitClient: true,
	})
	update.AddStringFlag(config.ArgDataCenterId, "", "", config.DatacenterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgLanId, config.ArgIdShort, "", config.LanId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLansIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgName, config.ArgNameShort, "", "The name of the LAN")
	update.AddStringFlag(config.ArgPccId, "", "", "The unique Id of the Private Cross-Connect the LAN will connect to")
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgPccId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getPccsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(config.ArgPublic, "", config.DefaultPublic, "Public option for LAN")
	update.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for LAN update to be executed")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for LAN update [seconds]")

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, lanCmd, core.CommandBuilder{
		Namespace: "lan",
		Resource:  "lan",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a LAN",
		LongDesc: `Use this command to delete a specified LAN from a Virtual Data Center.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* LAN Id`,
		Example:    deleteLanExample,
		PreCmdRun:  PreRunDcLanIds,
		CmdRun:     RunLanDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(config.ArgDataCenterId, "", "", config.DatacenterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(config.ArgLanId, config.ArgIdShort, "", config.LanId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLansIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for Request for LAN deletion to be executed")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for LAN deletion [seconds]")

	return lanCmd
}

func RunLanList(c *core.CommandConfig) error {
	lans, _, err := c.Lans().List(viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getLanPrint(nil, c, getLans(lans)))
}

func RunLanGet(c *core.CommandConfig) error {
	c.Printer.Verbose("Lan with id: %v from Datacenter with id: %v is getting...",
		viper.GetString(core.GetFlagName(c.NS, config.ArgLanId)), viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)))
	l, _, err := c.Lans().Get(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgLanId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLanPrint(nil, c, []v5.Lan{*l}))
}

func RunLanCreate(c *core.CommandConfig) error {
	name := viper.GetString(core.GetFlagName(c.NS, config.ArgName))
	public := viper.GetBool(core.GetFlagName(c.NS, config.ArgPublic))
	properties := ionoscloud.LanPropertiesPost{
		Name:   &name,
		Public: &public,
	}
	c.Printer.Verbose("Properties set for creating the Lan: Name: %v, Public: %v", name, public)
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgPccId)) {
		pcc := viper.GetString(core.GetFlagName(c.NS, config.ArgPccId))
		properties.SetPcc(pcc)
		c.Printer.Verbose("Property Pcc set: %v", pcc)
	}
	input := v5.LanPost{
		LanPost: ionoscloud.LanPost{
			Properties: &properties,
		},
	}
	l, resp, err := c.Lans().Create(viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)), input)
	if resp != nil {
		c.Printer.Verbose("Request href: %v ", resp.Header.Get("location"))
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		OutputJSON:     l,
		KeyValue:       getLanPostsKVMaps([]v5.LanPost{*l}),
		Columns:        getLansCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr()),
		ApiResponse:    resp,
		Resource:       "lan",
		Verb:           "create",
		WaitForRequest: viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForRequest)),
	})
}

func RunLanUpdate(c *core.CommandConfig) error {
	input := v5.LanProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, config.ArgName))
		input.SetName(name)
		c.Printer.Verbose("Property Name set: %v", name)
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgPublic)) {
		public := viper.GetBool(core.GetFlagName(c.NS, config.ArgPublic))
		input.SetPublic(public)
		c.Printer.Verbose("Property Public set: %v", public)
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgPccId)) {
		pcc := viper.GetString(core.GetFlagName(c.NS, config.ArgPccId))
		input.SetPcc(pcc)
		c.Printer.Verbose("Property Pcc set: %v", pcc)
	}
	lanUpdated, resp, err := c.Lans().Update(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgLanId)),
		input,
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getLanPrint(resp, c, []v5.Lan{*lanUpdated}))
}

func RunLanDelete(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete lan"); err != nil {
		return err
	}
	c.Printer.Verbose("Lan with id: %v is deleting...", viper.GetString(core.GetFlagName(c.NS, config.ArgLanId)))
	resp, err := c.Lans().Delete(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgLanId)),
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getLanPrint(resp, c, nil))
}

// Output Printing

var defaultLanCols = []string{"LanId", "Name", "Public", "PccId", "State"}

type LanPrint struct {
	LanId  string `json:"LanId,omitempty"`
	Name   string `json:"Name,omitempty"`
	Public bool   `json:"Public,omitempty"`
	PccId  string `json:"PccId,omitempty"`
	State  string `json:"State,omitempty"`
}

func getLanPrint(resp *v5.Response, c *core.CommandConfig, lans []v5.Lan) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForRequest = viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForRequest))
		}
		if lans != nil {
			r.OutputJSON = lans
			r.KeyValue = getLansKVMaps(lans)
			r.Columns = getLansCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getLansCols(flagName string, outErr io.Writer) []string {
	var cols []string
	if viper.IsSet(flagName) {
		cols = viper.GetStringSlice(flagName)
	} else {
		return defaultLanCols
	}

	columnsMap := map[string]string{
		"LanId":  "LanId",
		"Name":   "Name",
		"Public": "Public",
		"PccId":  "PccId",
		"State":  "State",
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

func getLans(lans v5.Lans) []v5.Lan {
	ls := make([]v5.Lan, 0)
	for _, s := range *lans.Items {
		ls = append(ls, v5.Lan{Lan: s})
	}
	return ls
}

func getLansKVMaps(ls []v5.Lan) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(ls))
	for _, l := range ls {
		var lanprint LanPrint
		if id, ok := l.GetIdOk(); ok && id != nil {
			lanprint.LanId = *id
		}
		if properties, ok := l.GetPropertiesOk(); ok && properties != nil {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				lanprint.Name = *name
			}
			if public, ok := properties.GetPublicOk(); ok && public != nil {
				lanprint.Public = *public
			}
			if pccId, ok := properties.GetPccOk(); ok && pccId != nil {
				lanprint.PccId = *pccId
			}
		}
		if metadata, ok := l.GetMetadataOk(); ok && metadata != nil {
			if state, ok := metadata.GetStateOk(); ok && state != nil {
				lanprint.State = *state
			}
		}
		o := structs.Map(lanprint)
		out = append(out, o)
	}
	return out
}

func getLanPostsKVMaps(ls []v5.LanPost) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(ls))
	for _, l := range ls {
		properties := l.GetProperties()
		var lanprint LanPrint
		if id, ok := l.GetIdOk(); ok && id != nil {
			lanprint.LanId = *id
		}
		if name, ok := properties.GetNameOk(); ok && name != nil {
			lanprint.Name = *name
		}
		if public, ok := properties.GetPublicOk(); ok && public != nil {
			lanprint.Public = *public
		}
		if pccId, ok := properties.GetPccOk(); ok && pccId != nil {
			lanprint.PccId = *pccId
		}
		o := structs.Map(lanprint)
		out = append(out, o)
	}
	return out
}

func getLansIds(outErr io.Writer, datacenterId string) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := v5.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		config.GetServerUrl(),
	)
	clierror.CheckError(err, outErr)
	lanSvc := v5.NewLanService(clientSvc.Get(), context.TODO())
	lans, _, err := lanSvc.List(datacenterId)
	clierror.CheckError(err, outErr)
	lansIds := make([]string, 0)
	if items, ok := lans.Lans.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				lansIds = append(lansIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return lansIds
}
