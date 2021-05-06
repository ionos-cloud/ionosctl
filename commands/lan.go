package commands

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/hashicorp/go-multierror"
	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func lan() *builder.Command {
	ctx := context.TODO()
	lanCmd := &builder.Command{
		Command: &cobra.Command{
			Use:              "lan",
			Short:            "LAN Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl lan` + "`" + ` allow you to create, list, get, update, delete LANs.`,
			TraverseChildren: true,
		},
	}
	globalFlags := lanCmd.GlobalFlags()
	globalFlags.StringP(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = viper.BindPFlag(builder.GetGlobalFlagName(lanCmd.Name(), config.ArgDataCenterId), globalFlags.Lookup(config.ArgDataCenterId))
	_ = lanCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	globalFlags.StringSlice(config.ArgCols, defaultLanCols, "Columns to be printed in the standard output. Example: --cols \"ResourceId,Name\"")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(lanCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))

	/*
		List Command
	*/
	builder.NewCommand(ctx, lanCmd, PreRunGlobalDcId, RunLanList, "list", "List LANs",
		"Use this command to retrieve a list of LANs provisioned in a specific Virtual Data Center.\n\nRequired values to run command:\n\n* Data Center Id",
		listLanExample, true)

	/*
		Get Command
	*/
	get := builder.NewCommand(ctx, lanCmd, PreRunGlobalDcIdLanId, RunLanGet, "get", "Get a LAN",
		"Use this command to retrieve information of a given LAN.\n\nRequired values to run command:\n\n* Data Center Id\n* LAN Id",
		getLanExample, true)
	get.AddStringFlag(config.ArgLanId, "", "", config.RequiredFlagLanId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLansIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(lanCmd.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := builder.NewCommand(ctx, lanCmd, PreRunGlobalDcId, RunLanCreate, "create", "Create a LAN",
		`Use this command to create a new LAN within a Virtual Data Center on your account. The name, the public option and the Private Cross-Connect Id can be set.

NOTE: IP Failover is configured after LAN creation using an update command.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option.

Required values to run command:

* Data Center Id`, createLanExample, true)
	create.AddStringFlag(config.ArgLanName, "", "", "The name of the LAN")
	create.AddBoolFlag(config.ArgLanPublic, "", config.DefaultLanPublic, "Indicates if the LAN faces the public Internet (true) or not (false)")
	create.AddStringFlag(config.ArgPccId, "", "", "The unique Id of the Private Cross-Connect the LAN will connect to")
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgPccId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getPccsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(config.ArgWaitForRequest, "", config.DefaultWait, "Wait for LAN to be created")
	create.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for LAN to be created [seconds]")

	/*
		Update Command
	*/
	update := builder.NewCommand(ctx, lanCmd, PreRunGlobalDcIdLanId, RunLanUpdate, "update", "Update a LAN",
		`Use this command to update a specified LAN. You can update the name, the public option for LAN and the Pcc Id to connect the LAN to a Private Cross-Connect.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option.

Required values to run command:

* Data Center Id
* LAN Id`, updateLanExample, true)
	update.AddStringFlag(config.ArgLanId, "", "", config.RequiredFlagLanId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLansIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(lanCmd.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgLanName, "", "", "The name of the LAN")
	update.AddStringFlag(config.ArgPccId, "", "", "The unique Id of the Private Cross-Connect the LAN will connect to")
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgPccId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getPccsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(config.ArgLanPublic, "", config.DefaultLanPublic, "Public option for LAN")
	update.AddBoolFlag(config.ArgWaitForRequest, "", config.DefaultWait, "Wait for LAN to be updated")
	update.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for LAN to be updated [seconds]")

	/*
		Delete Command
	*/
	deleteCmd := builder.NewCommand(ctx, lanCmd, PreRunGlobalDcIdLanId, RunLanDelete, "delete", "Delete a LAN",
		`Use this command to delete a specified LAN from a Virtual Data Center.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option. You can force the command to execute without user input using `+"`"+`--force`+"`"+` option.

Required values to run command:

* Data Center Id
* LAN Id`, deleteLanExample, true)
	deleteCmd.AddStringFlag(config.ArgLanId, "", "", config.RequiredFlagLanId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLansIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(lanCmd.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, "", config.DefaultWait, "Wait for LAN to be deleted")
	deleteCmd.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for LAN to be deleted [seconds]")

	return lanCmd
}

func PreRunGlobalDcId(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredGlobalFlags(c.ParentName, config.ArgDataCenterId)
}

func PreRunGlobalDcIdLanId(c *builder.PreCommandConfig) error {
	var result *multierror.Error
	if err := builder.CheckRequiredGlobalFlags(c.ParentName, config.ArgDataCenterId); err != nil {
		result = multierror.Append(result, err)
	}
	if err := builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgLanId); err != nil {
		result = multierror.Append(result, err)
	}
	if result != nil {
		return result
	}
	return nil
}

func RunLanList(c *builder.CommandConfig) error {
	lans, _, err := c.Lans().List(viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)))
	if err != nil {
		return err
	}
	ss := getLans(lans)
	return c.Printer.Print(printer.Result{
		OutputJSON: lans,
		KeyValue:   getLansKVMaps(ss),
		Columns:    getLansCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
	})
}

func RunLanGet(c *builder.CommandConfig) error {
	lan, _, err := c.Lans().Get(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLanId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		OutputJSON: lan,
		KeyValue:   getLansKVMaps([]resources.Lan{*lan}),
		Columns:    getLansCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
	})
}

func RunLanCreate(c *builder.CommandConfig) error {
	name := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLanName))
	public := viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgLanPublic))
	properties := ionoscloud.LanPropertiesPost{
		Name:   &name,
		Public: &public,
	}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgPccId)) {
		properties.SetPcc(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgPccId)))
	}
	input := resources.LanPost{
		LanPost: ionoscloud.LanPost{
			Properties: &properties,
		},
	}
	l, resp, err := c.Lans().Create(viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)), input)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		OutputJSON:     l,
		KeyValue:       getLanPostsKVMaps([]resources.LanPost{*l}),
		Columns:        getLansCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
		ApiResponse:    resp,
		Resource:       "lan",
		Verb:           "create",
		WaitForRequest: viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWaitForRequest)),
	})
}

func RunLanUpdate(c *builder.CommandConfig) error {
	input := resources.LanProperties{}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgLanName)) {
		input.SetName(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLanName)))
	}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgLanPublic)) {
		input.SetPublic(viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgLanPublic)))
	}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgPccId)) {
		input.SetPcc(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgPccId)))
	}
	lan, resp, err := c.Lans().Update(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLanId)),
		input,
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		OutputJSON:     lan,
		KeyValue:       getLansKVMaps([]resources.Lan{*lan}),
		Columns:        getLansCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
		ApiResponse:    resp,
		Resource:       "lan",
		Verb:           "update",
		WaitForRequest: viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWaitForRequest)),
	})
}

func RunLanDelete(c *builder.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete lan"); err != nil {
		return err
	}
	resp, err := c.Lans().Delete(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLanId)),
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		ApiResponse:    resp,
		Resource:       "lan",
		Verb:           "delete",
		WaitForRequest: viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWaitForRequest)),
	})
}

// Output Printing

var defaultLanCols = []string{"LanId", "Name", "Public", "PccId"}

type LanPrint struct {
	LanId  string `json:"LanId,omitempty"`
	Name   string `json:"Name,omitempty"`
	Public bool   `json:"Public,omitempty"`
	PccId  string `json:"PccId,omitempty"`
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

func getLans(lans resources.Lans) []resources.Lan {
	ls := make([]resources.Lan, 0)
	for _, s := range *lans.Items {
		ls = append(ls, resources.Lan{Lan: s})
	}
	return ls
}

func getLansKVMaps(ls []resources.Lan) []map[string]interface{} {
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

func getLanPostsKVMaps(ls []resources.LanPost) []map[string]interface{} {
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
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)
	lanSvc := resources.NewLanService(clientSvc.Get(), context.TODO())
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
