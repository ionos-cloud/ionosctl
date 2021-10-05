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
	globalFlags.StringSliceP(config.ArgCols, "", defaultPccCols, printer.ColsMessage(defaultPccCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(pccCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = pccCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultPccCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	core.NewCommand(ctx, pccCmd, core.CommandBuilder{
		Namespace:  "pcc",
		Resource:   "pcc",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Private Cross-Connects",
		LongDesc:   "Use this command to get a list of existing Private Cross-Connects available on your account.",
		Example:    listPccsExample,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunPccList,
		InitClient: true,
	})

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
	get.AddStringFlag(cloudapiv6.ArgPccId, cloudapiv6.ArgIdShort, "", cloudapiv6.PccId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgPccId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.PccsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

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
	create.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Private Cross-Connect creation to be executed")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Private Cross-Connect creation [seconds]")

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
	update.AddStringFlag(cloudapiv6.ArgPccId, cloudapiv6.ArgIdShort, "", cloudapiv6.PccId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgPccId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.PccsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Private Cross-Connect update to be executed")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Private Cross-Connect update [seconds]")

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
	deleteCmd.AddStringFlag(cloudapiv6.ArgPccId, cloudapiv6.ArgIdShort, "", cloudapiv6.PccId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgPccId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getPccsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Private Cross-Connect deletion to be executed")
	deleteCmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete all Private Cross-Connects.")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Private Cross-Connect deletion [seconds]")

	pccCmd.AddCommand(PeersCmd())

	return pccCmd
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
	pccs, resp, err := c.CloudApiV6Services.Pccs().List()
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getPccPrint(nil, c, getPccs(pccs)))
}

func RunPccGet(c *core.CommandConfig) error {
	c.Printer.Verbose("Private cross connect with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPccId)))
	u, resp, err := c.CloudApiV6Services.Pccs().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPccId)))
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getPccPrint(nil, c, getPcc(u)))
}

func RunPccCreate(c *core.CommandConfig) error {
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
	c.Printer.Verbose("Properties set for creating the private cross connect: Name: %v, Description: %v", name, description)
	u, resp, err := c.CloudApiV6Services.Pccs().Create(newUser)
	if resp != nil {
		c.Printer.Verbose("Request href: %v ", resp.Header.Get("location"))
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getPccPrint(resp, c, getPcc(u)))
}

func RunPccUpdate(c *core.CommandConfig) error {
	oldPcc, resp, err := c.CloudApiV6Services.Pccs().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPccId)))
	if err != nil {
		return err
	}
	newProperties := getPccInfo(oldPcc, c)
	pccUpd, resp, err := c.CloudApiV6Services.Pccs().Update(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPccId)), *newProperties)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getPccPrint(resp, c, getPcc(pccUpd)))
}

func RunPccDelete(c *core.CommandConfig) error {
	var resp *resources.Response
	pccId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPccId))
	allFlag := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll))
	if allFlag {
		err := DeleteAllPccs(c)
		if err != nil {
			return err
		}
	} else {
		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete private cross-connect"); err != nil {
			return err
		}
		c.Printer.Verbose("Private cross connect with id: %v is deleting...", pccId)
		resp, err := c.CloudApiV6Services.Pccs().Delete(pccId)
		if resp != nil {
			c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
		}
		if err != nil {
			return err
		}
		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
			return err
		}
	}
	return c.Printer.Print(getPccPrint(resp, c, nil))
}

func getPccInfo(oldUser *resources.PrivateCrossConnect, c *core.CommandConfig) *resources.PrivateCrossConnectProperties {
	var namePcc, description string
	if properties, ok := oldUser.GetPropertiesOk(); ok && properties != nil {
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
			namePcc = viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
			c.Printer.Verbose("Property Name set: %v", namePcc)
		} else {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				namePcc = *name
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgDescription)) {
			description = viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDescription))
			c.Printer.Verbose("Property Description set: %v", description)
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
	_ = c.Printer.Print("PrivateCrossConnects to be deleted:")
	pccs, resp, err := c.CloudApiV6Services.Pccs().List()
	if err != nil {
		return err
	}
	if pccsItems, ok := pccs.GetItemsOk(); ok && pccsItems != nil {
		for _, pcc := range *pccsItems {
			if id, ok := pcc.GetIdOk(); ok && id != nil {
				_ = c.Printer.Print("PrivateCrossConnect Id: " + *id)
			}
			if properties, ok := pcc.GetPropertiesOk(); ok && properties != nil {
				if name, ok := properties.GetNameOk(); ok && name != nil {
					_ = c.Printer.Print("PrivateCrossConnect Name: " + *name)
				}
			}
		}
		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete all the PrivateCrossConnects"); err != nil {
			return err
		}
		c.Printer.Verbose("Deleting all the PrivateCrossConnects...")

		for _, pcc := range *pccsItems {
			if id, ok := pcc.GetIdOk(); ok && id != nil {
				c.Printer.Verbose("Deleting PrivateCrossConnect with id: %v...", *id)
				resp, err = c.CloudApiV6Services.Pccs().Delete(*id)
				if resp != nil {
					c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
				}
				if err != nil {
					return err
				}
				if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
					return err
				}
			}
		}
	}
	return c.Printer.Print(getPccPrint(resp, c, nil))
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
	globalFlags.StringSliceP(config.ArgCols, "", defaultPccPeersCols,
		printer.ColsMessage(defaultPccPeersCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(peerCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = peerCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
	listPeers.AddStringFlag(cloudapiv6.ArgPccId, "", "", cloudapiv6.PccId, core.RequiredFlagOption())
	_ = listPeers.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgPccId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getPccsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	return peerCmd
}

func RunPccPeersList(c *core.CommandConfig) error {
	c.Printer.Verbose("Getting Peers from Private Cross-Connect with ID: %v...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPccId)))
	u, resp, err := c.CloudApiV6Services.Pccs().GetPeers(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPccId)))
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getPccPeerPrint(c, *u))
}

// Output Printing

var defaultPccCols = []string{"PccId", "Name", "Description", "State"}

type PccPrint struct {
	PccId       string `json:"PccId,omitempty"`
	Name        string `json:"Name,omitempty"`
	Description string `json:"Description,omitempty"`
	State       string `json:"State,omitempty"`
}

func getPccPrint(resp *resources.Response, c *core.CommandConfig, pccs []resources.PrivateCrossConnect) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForRequest = viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForRequest))
		}
		if pccs != nil {
			r.OutputJSON = pccs
			r.KeyValue = getPccsKVMaps(pccs)
			r.Columns = getPccCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), os.Stderr)
		}
	}
	return r
}

var defaultPccPeersCols = []string{"LanId", "LanName", "DatacenterId", "DatacenterName", "Location"}

type PccPeerPrint struct {
	LanId          string `json:"LanId,omitempty"`
	LanName        string `json:"LanName,omitempty"`
	DatacenterId   string `json:"DatacenterId,omitempty"`
	DatacenterName string `json:"DatacenterName,omitempty"`
	Location       string `json:"Location,omitempty"`
}

func getPccPeerPrint(c *core.CommandConfig, pccs []resources.Peer) printer.Result {
	r := printer.Result{}
	if c != nil {
		if pccs != nil {
			r.OutputJSON = pccs
			r.KeyValue = getPccPeersKVMaps(pccs)
			r.Columns = getPccPeersCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), os.Stderr)
		}
	}
	return r
}

func getPccPeersCols(flagName string, outErr io.Writer) []string {
	if viper.IsSet(flagName) {
		var pccCols []string
		columnsMap := map[string]string{
			"LanId":          "LanId",
			"LanName":        "LanName",
			"DatacenterId":   "DatacenterId",
			"DatacenterName": "DatacenterName",
			"Location":       "Location",
		}
		for _, k := range viper.GetStringSlice(flagName) {
			col := columnsMap[k]
			if col != "" {
				pccCols = append(pccCols, col)
			} else {
				clierror.CheckError(errors.New("unknown column "+k), outErr)
			}
		}
		return pccCols
	} else {
		return defaultPccPeersCols
	}
}

func getPccCols(flagName string, outErr io.Writer) []string {
	if viper.IsSet(flagName) {
		var pccCols []string
		columnsMap := map[string]string{
			"PccId":       "PccId",
			"Name":        "Name",
			"Description": "Description",
			"State":       "State",
		}
		for _, k := range viper.GetStringSlice(flagName) {
			col := columnsMap[k]
			if col != "" {
				pccCols = append(pccCols, col)
			} else {
				clierror.CheckError(errors.New("unknown column "+k), outErr)
			}
		}
		return pccCols
	} else {
		return defaultPccCols
	}
}

func getPccs(pccs resources.PrivateCrossConnects) []resources.PrivateCrossConnect {
	u := make([]resources.PrivateCrossConnect, 0)
	if items, ok := pccs.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			u = append(u, resources.PrivateCrossConnect{PrivateCrossConnect: item})
		}
	}
	return u
}

func getPcc(u *resources.PrivateCrossConnect) []resources.PrivateCrossConnect {
	pccs := make([]resources.PrivateCrossConnect, 0)
	if u != nil {
		pccs = append(pccs, resources.PrivateCrossConnect{PrivateCrossConnect: u.PrivateCrossConnect})
	}
	return pccs
}

func getPccsKVMaps(us []resources.PrivateCrossConnect) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(us))
	for _, u := range us {
		var uPrint PccPrint
		if id, ok := u.GetIdOk(); ok && id != nil {
			uPrint.PccId = *id
		}
		if properties, ok := u.GetPropertiesOk(); ok && properties != nil {
			if n, ok := properties.GetNameOk(); ok && n != nil {
				uPrint.Name = *n
			}
			if d, ok := properties.GetDescriptionOk(); ok && d != nil {
				uPrint.Description = *d
			}
		}
		if metadata, ok := u.GetMetadataOk(); ok && metadata != nil {
			if state, ok := metadata.GetStateOk(); ok && state != nil {
				uPrint.State = *state
			}
		}
		o := structs.Map(uPrint)
		out = append(out, o)
	}
	return out
}

func getPccPeersKVMaps(ps []resources.Peer) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(ps))
	for _, p := range ps {
		var uPrint PccPeerPrint
		if lanId, ok := p.GetIdOk(); ok && lanId != nil {
			uPrint.LanId = *lanId
		}
		if loc, ok := p.GetLocationOk(); ok && loc != nil {
			uPrint.Location = *loc
		}
		if n, ok := p.GetNameOk(); ok && n != nil {
			uPrint.LanName = *n
		}
		if dcId, ok := p.GetDatacenterIdOk(); ok && dcId != nil {
			uPrint.DatacenterId = *dcId
		}
		if dcName, ok := p.GetDatacenterNameOk(); ok && dcName != nil {
			uPrint.DatacenterName = *dcName
		}
		o := structs.Map(uPrint)
		out = append(out, o)
	}
	return out
}

func getPccsIds(outErr io.Writer) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		config.GetServerUrl(),
	)
	clierror.CheckError(err, outErr)
	pccSvc := resources.NewPrivateCrossConnectService(clientSvc.Get(), context.TODO())
	pccs, _, err := pccSvc.List()
	clierror.CheckError(err, outErr)
	pccsIds := make([]string, 0)
	if items, ok := pccs.PrivateCrossConnects.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				pccsIds = append(pccsIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return pccsIds
}
