package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"go.uber.org/multierr"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/query"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/waiter"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func LanCmd() *core.Command {
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
	globalFlags.StringSliceP(config.ArgCols, "", defaultLanCols, printer.ColsMessage(defaultLanCols))
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
		LongDesc:   "Use this command to retrieve a list of LANs provisioned in a specific Virtual Data Center.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.LANsFiltersUsage() + "\n\nRequired values to run command:\n\n* Data Center Id",
		Example:    listLanExample,
		PreCmdRun:  PreRunLansList,
		CmdRun:     RunLanList,
		InitClient: true,
	})
	list.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddIntFlag(cloudapiv6.ArgMaxResults, cloudapiv6.ArgMaxResultsShort, 0, cloudapiv6.ArgMaxResultsDescription)
	list.AddIntFlag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, config.DefaultListDepth, cloudapiv6.ArgDepthDescription)
	list.AddStringFlag(cloudapiv6.ArgOrderBy, "", "", cloudapiv6.ArgOrderByDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LANsFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv6.ArgFilters, cloudapiv6.ArgFiltersShort, []string{""}, cloudapiv6.ArgFiltersDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LANsFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddBoolFlag(config.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)

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
	get.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapiv6.ArgLanId, cloudapiv6.ArgIdShort, "", cloudapiv6.LanId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LansIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(config.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)
	get.AddIntFlag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, config.DefaultGetDepth, cloudapiv6.ArgDepthDescription)

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
	create.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed LAN", "The name of the LAN")
	create.AddBoolFlag(cloudapiv6.ArgPublic, cloudapiv6.ArgPublicShort, cloudapiv6.DefaultPublic, "Indicates if the LAN faces the public Internet (true) or not (false). E.g.: --public=true, --public=false")
	create.AddStringFlag(cloudapiv6.ArgPccId, "", "", "The unique Id of the Private Cross-Connect the LAN will connect to")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgPccId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.PccsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for LAN creation to be executed")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for LAN creation [seconds]")
	create.AddIntFlag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, config.DefaultCreateDepth, cloudapiv6.ArgDepthDescription)

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
	update.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgLanId, cloudapiv6.ArgIdShort, "", cloudapiv6.LanId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LansIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "The name of the LAN")
	update.AddStringFlag(cloudapiv6.ArgPccId, "", "", "The unique Id of the Private Cross-Connect the LAN will connect to")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgPccId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.PccsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(cloudapiv6.ArgPublic, "", cloudapiv6.DefaultPublic, "Public option for LAN. E.g.: --public=true, --public=false")
	update.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for LAN update to be executed")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for LAN update [seconds]")
	update.AddIntFlag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, config.DefaultUpdateDepth, cloudapiv6.ArgDepthDescription)

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
		PreCmdRun:  PreRunLanDelete,
		CmdRun:     RunLanDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(cloudapiv6.ArgLanId, cloudapiv6.ArgIdShort, "", cloudapiv6.LanId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LansIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for Request for LAN deletion to be executed")
	deleteCmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete all Lans from a Virtual Data Center.")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for LAN deletion [seconds]")
	deleteCmd.AddIntFlag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, config.DefaultDeleteDepth, cloudapiv6.ArgDepthDescription)

	return lanCmd
}

func PreRunLansList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId},
		[]string{cloudapiv6.ArgAll},
	); err != nil {
		return err
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgFilters)) {
		return query.ValidateFilters(c, completer.LANsFilters(), completer.LANsFiltersUsage())
	}
	return nil
}

func PreRunLanDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgLanId},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgAll},
	)
}

func RunLanList(c *core.CommandConfig) error {
	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	if !structs.IsZero(listQueryParams) {
		c.Printer.Verbose("List Query Parameters set: %v", utils.GetPropertiesKVSet(listQueryParams))
		if !structs.IsZero(listQueryParams.QueryParams) {
			c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(listQueryParams.QueryParams))
		}
	}
	lans, resp, err := c.CloudApiV6Services.Lans().List(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)), listQueryParams)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLanPrint(nil, c, getLans(lans)))
}

func RunLanGet(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	if !structs.IsZero(queryParams) {
		c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(queryParams))
	}
	c.Printer.Verbose("Lan with id: %v from Datacenter with id: %v is getting...",
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLanId)), viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	l, resp, err := c.CloudApiV6Services.Lans().Get(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLanId)),
		queryParams,
	)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLanPrint(nil, c, []resources.Lan{*l}))
}

func RunLanCreate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	if !structs.IsZero(queryParams) {
		c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(queryParams))
	}
	name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
	public := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgPublic))
	properties := ionoscloud.LanPropertiesPost{
		Name:   &name,
		Public: &public,
	}
	c.Printer.Verbose("Properties set for creating the Lan: Name: %v, Public: %v", name, public)
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPccId)) {
		pcc := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPccId))
		properties.SetPcc(pcc)
		c.Printer.Verbose("Property Pcc set: %v", pcc)
	}
	input := resources.LanPost{
		LanPost: ionoscloud.LanPost{
			Properties: &properties,
		},
	}
	c.Printer.Verbose("Creating LAN in Datacenter with ID: %v...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	l, resp, err := c.CloudApiV6Services.Lans().Create(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)), input, queryParams)
	if resp != nil && printer.GetId(resp) != "" {
		c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		OutputJSON:     l,
		KeyValue:       getLanPostsKVMaps([]resources.LanPost{*l}),
		Columns:        getLansCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr()),
		ApiResponse:    resp,
		Resource:       "lan",
		Verb:           "create",
		WaitForRequest: viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForRequest)),
	})
}

func RunLanUpdate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	if !structs.IsZero(queryParams) {
		c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(queryParams))
	}
	input := resources.LanProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
		input.SetName(name)
		c.Printer.Verbose("Property Name set: %v", name)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPublic)) {
		public := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgPublic))
		input.SetPublic(public)
		c.Printer.Verbose("Property Public set: %v", public)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPccId)) {
		pcc := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPccId))
		input.SetPcc(pcc)
		c.Printer.Verbose("Property Pcc set: %v", pcc)
	}
	c.Printer.Verbose("Updating LAN with ID: %v from Datacenter with ID: %v...",
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLanId)), viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	lanUpdated, resp, err := c.CloudApiV6Services.Lans().Update(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLanId)),
		input,
		queryParams,
	)
	if resp != nil && printer.GetId(resp) != "" {
		c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getLanPrint(resp, c, []resources.Lan{*lanUpdated}))
}

func RunLanDelete(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	if !structs.IsZero(queryParams) {
		c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(queryParams))
	}
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	lanId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLanId))
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllLans(c); err != nil {
			return err
		}
		return c.Printer.Print(printer.Result{Resource: c.Resource, Verb: c.Verb})
	} else {
		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete lan"); err != nil {
			return err
		}
		c.Printer.Verbose("Starting deleting LAN with ID: %v from Datacenter with ID: %v...", lanId, dcId)
		resp, err := c.CloudApiV6Services.Lans().Delete(dcId, lanId, queryParams)
		if resp != nil && printer.GetId(resp) != "" {
			c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
		}
		if err != nil {
			return err
		}
		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
			return err
		}
		return c.Printer.Print(getLanPrint(resp, c, nil))
	}
}

func DeleteAllLans(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	if !structs.IsZero(queryParams) {
		c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(queryParams))
	}
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	c.Printer.Verbose("Datacenter ID: %v", dcId)
	c.Printer.Verbose("Getting Lans...")
	lans, resp, err := c.CloudApiV6Services.Lans().List(dcId, resources.ListQueryParams{})
	if err != nil {
		return err
	}
	if lansItems, ok := lans.GetItemsOk(); ok && lansItems != nil {
		if len(*lansItems) > 0 {
			_ = c.Printer.Print("Lans to be deleted:")
			for _, lan := range *lansItems {
				toPrint := ""
				if id, ok := lan.GetIdOk(); ok && id != nil {
					toPrint += "Lan Id: " + *id
				}
				if properties, ok := lan.GetPropertiesOk(); ok && properties != nil {
					if name, ok := properties.GetNameOk(); ok && name != nil {
						toPrint += " Lan Name: " + *name
					}
				}
				_ = c.Printer.Print(toPrint)
			}
			if err = utils.AskForConfirm(c.Stdin, c.Printer, "delete all the Lans"); err != nil {
				return err
			}
			c.Printer.Verbose("Deleting all the Lans...")
			var multiErr error
			for _, lan := range *lansItems {
				if id, ok := lan.GetIdOk(); ok && id != nil {
					c.Printer.Verbose("Starting deleting Lan with id: %v...", *id)
					resp, err = c.CloudApiV6Services.Lans().Delete(dcId, *id, queryParams)
					if resp != nil && printer.GetId(resp) != "" {
						c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
					}
					if err != nil {
						multiErr = multierr.Append(multiErr, fmt.Errorf(config.DeleteAllAppendErr, c.Resource, *id, err))
						continue
					} else {
						_ = c.Printer.Print(fmt.Sprintf(config.StatusDeletingAll, c.Resource, *id))
					}
					if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
						multiErr = multierr.Append(multiErr, fmt.Errorf(config.WaitDeleteAllAppendErr, c.Resource, *id, err))
						continue
					}
				}
			}
			if multiErr != nil {
				return multiErr
			}
			return nil
		} else {
			return errors.New("no Lans found")
		}
	} else {
		return errors.New("could not get items of Lans")
	}
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

func getLanPrint(resp *resources.Response, c *core.CommandConfig, lans []resources.Lan) printer.Result {
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

func getLans(lans resources.Lans) []resources.Lan {
	lanObjs := make([]resources.Lan, 0)
	if items, ok := lans.GetItemsOk(); ok && items != nil {
		for _, lan := range *items {
			lanObjs = append(lanObjs, resources.Lan{Lan: lan})
		}
	}
	return lanObjs
}

func getLansKVMaps(ls []resources.Lan) []map[string]interface{} {
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
