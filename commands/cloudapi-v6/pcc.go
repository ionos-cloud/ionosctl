package commands

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/query"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/waiter"
	"github.com/ionos-cloud/ionosctl/v6/internal/confirm"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/utils"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	allPccJSONPaths = map[string]string{
		"PccId":       "id",
		"Name":        "properties.name",
		"Description": "properties.description",
		"State":       "metadata.state",
	}

	allPccPeerJSONPaths = map[string]string{
		"LanId":          "id",
		"LanName":        "name",
		"DatacenterId":   "datacenterId",
		"DatacenterName": "datacenterName",
		"Location":       "location",
	}

	defaultPccCols      = []string{"PccId", "Name", "Description", "State"}
	defaultPccPeersCols = []string{"LanId", "LanName", "DatacenterId", "DatacenterName", "Location"}
)

func PccCmd() *core.Command {
	ctx := context.TODO()
	pccCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "pcc",
			Short:            "Private Cross-Connect Operations",
			Long:             "The sub-commands of `ionosctl pcc` allow you to list, get, create, update, delete Private Cross-Connect. To add Private Cross-Connect to a Lan, check the `ionosctl lan update` command.",
			TraverseChildren: true,
		},
	}
	globalFlags := pccCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultPccCols, printer.ColsMessage(defaultPccCols))
	_ = viper.BindPFlag(core.GetFlagName(pccCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = pccCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultPccCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, pccCmd, core.CommandBuilder{
		Namespace:  "pcc",
		Resource:   "pcc",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Private Cross-Connects",
		LongDesc:   "Use this command to get a list of existing Private Cross-Connects available on your account.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.PccsFiltersUsage(),
		Example:    listPccsExample,
		PreCmdRun:  PreRunPccList,
		CmdRun:     RunPccList,
		InitClient: true,
	})
	list.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, cloudapiv6.DefaultMaxResults, constants.DescMaxResults)
	list.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultListDepth, cloudapiv6.ArgDepthDescription)
	list.AddStringFlag(cloudapiv6.ArgOrderBy, "", "", cloudapiv6.ArgOrderByDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.PccsFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv6.ArgFilters, cloudapiv6.ArgFiltersShort, []string{""}, cloudapiv6.ArgFiltersDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.PccsFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, pccCmd, core.CommandBuilder{
		Namespace:  "pcc",
		Resource:   "pcc",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a Private Cross-Connect",
		LongDesc:   "Use this command to retrieve details about a specific Private Cross-Connect.\n\nRequired values to run command:\n\n* Pcc Id",
		Example:    getPccExample,
		PreCmdRun:  PreRunPccId,
		CmdRun:     RunPccGet,
		InitClient: true,
	})
	get.AddUUIDFlag(cloudapiv6.ArgPccId, cloudapiv6.ArgIdShort, "", cloudapiv6.PccId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgPccId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.PccsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)
	get.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.ArgDepthDescription)

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, pccCmd, core.CommandBuilder{
		Namespace:  "pcc",
		Resource:   "pcc",
		Verb:       "create",
		Aliases:    []string{"c"},
		ShortDesc:  "Create a Private Cross-Connect",
		LongDesc:   "Use this command to create a Private Cross-Connect. You can specify the name and the description for the Private Cross-Connect.",
		Example:    createPccExample,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunPccCreate,
		InitClient: true,
	})
	create.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed PrivateCrossConnect", "The name for the Private Cross-Connect")
	create.AddStringFlag(cloudapiv6.ArgDescription, cloudapiv6.ArgDescriptionShort, "", "The description for the Private Cross-Connect")
	create.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Private Cross-Connect creation to be executed")
	create.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Private Cross-Connect creation [seconds]")
	create.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultCreateDepth, cloudapiv6.ArgDepthDescription)

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, pccCmd, core.CommandBuilder{
		Namespace: "pcc",
		Resource:  "pcc",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a Private Cross-Connect",
		LongDesc: `Use this command to update details about a specific Private Cross-Connect. Name and description can be updated.

Required values to run command:

* Pcc Id`,
		Example:    updatePccExample,
		PreCmdRun:  PreRunPccId,
		CmdRun:     RunPccUpdate,
		InitClient: true,
	})
	update.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "The name for the Private Cross-Connect")
	update.AddStringFlag(cloudapiv6.ArgDescription, cloudapiv6.ArgDescriptionShort, "", "The description for the Private Cross-Connect")
	update.AddUUIDFlag(cloudapiv6.ArgPccId, cloudapiv6.ArgIdShort, "", cloudapiv6.PccId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgPccId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.PccsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Private Cross-Connect update to be executed")
	update.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Private Cross-Connect update [seconds]")
	update.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultUpdateDepth, cloudapiv6.ArgDepthDescription)

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, pccCmd, core.CommandBuilder{
		Namespace: "pcc",
		Resource:  "pcc",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a Private Cross-Connect",
		LongDesc: `Use this command to delete a Private Cross-Connect.

Required values to run command:

* Pcc Id`,
		Example:    deletePccExample,
		PreCmdRun:  PreRunPccDelete,
		CmdRun:     RunPccDelete,
		InitClient: true,
	})
	deleteCmd.AddUUIDFlag(cloudapiv6.ArgPccId, cloudapiv6.ArgIdShort, "", cloudapiv6.PccId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgPccId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.PccsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Private Cross-Connect deletion to be executed")
	deleteCmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete all Private Cross-Connects.")
	deleteCmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Private Cross-Connect deletion [seconds]")
	deleteCmd.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultDeleteDepth, cloudapiv6.ArgDepthDescription)

	pccCmd.AddCommand(PeersCmd())

	return pccCmd
}

func PreRunPccList(c *core.PreCommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgFilters)) {
		return query.ValidateFilters(c, completer.PccsFilters(), completer.PccsFiltersUsage())
	}
	return nil
}

func PreRunPccId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgPccId)
}

func PreRunPccDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgPccId},
		[]string{cloudapiv6.ArgAll},
	)
}

func RunPccList(c *core.CommandConfig) error {
	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	pccs, resp, err := c.CloudApiV6Services.Pccs().List(listQueryParams)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols, err := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("items", allPccJSONPaths, pccs.PrivateCrossConnects,
		tabheaders.GetHeadersAllDefault(defaultPccCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

func RunPccGet(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Private cross connect with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPccId))))

	u, resp, err := c.CloudApiV6Services.Pccs().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPccId)), queryParams)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols, err := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("", allPccJSONPaths, u.PrivateCrossConnect,
		tabheaders.GetHeadersAllDefault(defaultPccCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

func RunPccCreate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams

	name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
	description := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDescription))

	newUser := resources.PrivateCrossConnect{
		PrivateCrossConnect: ionoscloud.PrivateCrossConnect{
			Properties: &ionoscloud.PrivateCrossConnectProperties{
				Name:        &name,
				Description: &description,
			},
		},
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Properties set for creating the private cross connect: Name: %v, Description: %v", name, description))

	u, resp, err := c.CloudApiV6Services.Pccs().Create(newUser, queryParams)
	if resp != nil && printer.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}

	cols, err := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("", allPccJSONPaths, u.PrivateCrossConnect,
		tabheaders.GetHeadersAllDefault(defaultPccCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

func RunPccUpdate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	oldPcc, resp, err := c.CloudApiV6Services.Pccs().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPccId)), queryParams)
	if err != nil {
		return err
	}

	newProperties := getPccInfo(oldPcc, c)
	pccUpd, resp, err := c.CloudApiV6Services.Pccs().Update(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPccId)), *newProperties, queryParams)
	if resp != nil && printer.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}

	cols, err := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("", allPccJSONPaths, pccUpd.PrivateCrossConnect,
		tabheaders.GetHeadersAllDefault(defaultPccCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

func RunPccDelete(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	pccId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPccId))

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllPccs(c); err != nil {
			return err
		}

		fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Private Cross Connects successfully deleted"))
		return nil
	}

	if !confirm.Ask("delete private cross-connect") {
		return nil
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Starting deleting Private cross connect with id: %v...", pccId))

	resp, err := c.CloudApiV6Services.Pccs().Delete(pccId, queryParams)
	if resp != nil && printer.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Private Cross Connect successfully deleted"))
	return nil
}

func getPccInfo(oldUser *resources.PrivateCrossConnect, c *core.CommandConfig) *resources.PrivateCrossConnectProperties {
	var namePcc, description string

	if properties, ok := oldUser.GetPropertiesOk(); ok && properties != nil {
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
			namePcc = viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Name set: %v", namePcc))
		} else {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				namePcc = *name
			}
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgDescription)) {
			description = viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDescription))

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Description set: %v", description))
		} else {
			if desc, ok := properties.GetDescriptionOk(); ok && desc != nil {
				description = *desc
			}
		}
	}

	return &resources.PrivateCrossConnectProperties{
		PrivateCrossConnectProperties: ionoscloud.PrivateCrossConnectProperties{
			Name:        &namePcc,
			Description: &description,
		},
	}
}

func DeleteAllPccs(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting PrivateCrossConnects..."))

	pccs, resp, err := c.CloudApiV6Services.Pccs().List(cloudapiv6.ParentResourceListQueryParams)
	if err != nil {
		return err
	}

	pccsItems, ok := pccs.GetItemsOk()
	if !ok || pccsItems == nil {
		return fmt.Errorf("could not get items of PrivateCrossConnects")
	}

	if len(*pccsItems) <= 0 {
		return fmt.Errorf("no PrivateCrossConnects found")
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("PrivateCrossConnects to be deleted:"))

	for _, pcc := range *pccsItems {
		delIdAndName := ""

		if id, ok := pcc.GetIdOk(); ok && id != nil {
			delIdAndName += "PrivateCrossConnect Id: " + *id
		}

		if properties, ok := pcc.GetPropertiesOk(); ok && properties != nil {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				delIdAndName += " PrivateCrossConnect Name: " + *name
			}
		}

		fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput(delIdAndName))
	}

	if !confirm.Ask("delete all the private cross-connects") {
		return nil
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Deleting all the PrivateCrossConnects..."))

	var multiErr error
	for _, pcc := range *pccsItems {
		id, ok := pcc.GetIdOk()
		if !ok || id == nil {
			continue
		}

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
			"Starting deleting PrivateCrossConnect with id: %v...", *id))

		resp, err = c.CloudApiV6Services.Pccs().Delete(*id, queryParams)
		if resp != nil && printer.GetId(resp) != "" {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime))
		}
		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
			continue
		}

		fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput(constants.MessageDeletingAll, c.Resource, *id))

		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrWaitDeleteAll, c.Resource, *id, err))
		}
	}

	if multiErr != nil {
		return multiErr
	}

	return nil
}

func PeersCmd() *core.Command {
	ctx := context.TODO()
	peerCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "peers",
			Short:            "Private Cross-Connect Peers Operations",
			Long:             "The sub-command of `ionosctl pcc peers` allows you to get a list of Peers from a Private Cross-Connect.",
			TraverseChildren: true,
		},
	}
	globalFlags := peerCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultPccPeersCols,
		printer.ColsMessage(defaultPccPeersCols))
	_ = viper.BindPFlag(core.GetFlagName(peerCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = peerCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultPccPeersCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	listPeers := core.NewCommand(ctx, peerCmd, core.CommandBuilder{
		Namespace:  "pcc",
		Resource:   "peers",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "Get a list of Peers from a Private Cross-Connect",
		LongDesc:   "Use this command to get a list of Peers from a Private Cross-Connect.\n\nRequired values to run command:\n\n* Pcc Id",
		Example:    listPccPeersExample,
		PreCmdRun:  PreRunPccId,
		CmdRun:     RunPccPeersList,
		InitClient: true,
	})
	listPeers.AddUUIDFlag(cloudapiv6.ArgPccId, "", "", cloudapiv6.PccId, core.RequiredFlagOption())
	_ = listPeers.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgPccId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.PccsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	listPeers.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)

	return peerCmd
}

func RunPccPeersList(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Getting Peers from Private Cross-Connect with ID: %v...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPccId))))

	u, resp, err := c.CloudApiV6Services.Pccs().GetPeers(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPccId)))
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	peers := make([]ionoscloud.Peer, 0)

	if u != nil {
		for _, p := range *u {
			peers = append(peers, p.Peer)
		}
	}

	cols, err := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("", allPccPeerJSONPaths, peers,
		tabheaders.GetHeadersAllDefault(defaultPccPeersCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}
