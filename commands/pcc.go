package commands

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/fatih/structs"
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

func pcc() *builder.Command {
	ctx := context.TODO()
	pccCmd := &builder.Command{
		Command: &cobra.Command{
			Use:              "pcc",
			Short:            "Private Cross-Connect Operations",
			Long:             `The sub-command of ` + "`" + `ionosctl pcc` + "`" + ` allows you to list, get, create, update, delete Private Cross-Connect. To add Private Cross-Connect to a Lan, check the ` + "`" + `ionosctl lan` + "`" + ` commands.`,
			TraverseChildren: true,
		},
	}
	globalFlags := pccCmd.Command.PersistentFlags()
	globalFlags.StringSlice(config.ArgCols, defaultUserCols, "Columns to be printed in the standard output")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(pccCmd.Command.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))

	/*
		List Command
	*/
	builder.NewCommand(ctx, pccCmd, noPreRun, RunPccList, "list", "List Private Cross-Connects",
		"Use this command to get a list of existing Private Cross-Connects available on your account.", listPccsExample, true)

	/*
		Get Command
	*/
	get := builder.NewCommand(ctx, pccCmd, PreRunPccIdValidate, RunPccGet, "get", "Get a Private Cross-Connect",
		"Use this command to retrieve details about a specific Private Cross-Connect.\n\nRequired values to run command:\n\n* Pcc Id", getPccExample, true)
	get.AddStringFlag(config.ArgPccId, "", "", config.RequiredFlagPccId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgPccId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getPccsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get Peers Command
	*/
	getPeers := builder.NewCommand(ctx, pccCmd, PreRunPccIdValidate, RunPccGetPeers, "get-peers", "Get a Private Cross-Connect Peers",
		"Use this command to get a list of Peers from a Private Cross-Connect.\n\nRequired values to run command:\n\n* Pcc Id", getPccPeersExample, true)
	getPeers.AddStringFlag(config.ArgPccId, "", "", config.RequiredFlagPccId)
	_ = getPeers.Command.RegisterFlagCompletionFunc(config.ArgPccId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getPccsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := builder.NewCommand(ctx, pccCmd, noPreRun, RunPccCreate, "create", "Create a Private Cross-Connect",
		"Use this command to create a Private Cross-Connect. You can specify the name and the description for the Private Cross-Connect.",
		createPccExample, true)
	create.AddStringFlag(config.ArgPccName, "", "", "The name for the Private Cross-Connect")
	create.AddStringFlag(config.ArgPccDescription, "", "", "The description for the Private Cross-Connect")
	create.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for Private Cross-Connect to be created")
	create.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Private Cross-Connect to be created [seconds]")

	/*
		Update Command
	*/
	update := builder.NewCommand(ctx, pccCmd, PreRunPccIdValidate, RunPccUpdate, "update", "Update a Private Cross-Connect",
		`Use this command to update details about a specific Private Cross-Connect. Name and description can be updated.

Required values to run command:

* Pcc Id`, updatePccExample, true)
	update.AddStringFlag(config.ArgPccName, "", "", "The name for the Private Cross-Connect")
	update.AddStringFlag(config.ArgPccDescription, "", "", "The description for the Private Cross-Connect")
	update.AddStringFlag(config.ArgPccId, "", "", config.RequiredFlagPccId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgPccId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getPccsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for Private Cross-Connect to be updated")
	update.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Private Cross-Connect to be updated [seconds]")

	/*
		Delete Command
	*/
	deleteCmd := builder.NewCommand(ctx, pccCmd, PreRunPccIdValidate, RunPccDelete, "delete", "Delete a Private Cross-Connect",
		`Use this command to delete a Private Cross-Connect.

Required values to run command:

* Pcc Id`, deletePccExample, true)
	deleteCmd.AddStringFlag(config.ArgPccId, "", "", config.RequiredFlagPccId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgPccId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getPccsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for Private Cross-Connect to be deleted")
	deleteCmd.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Private Cross-Connect to be deleted [seconds]")

	return pccCmd
}

func PreRunPccIdValidate(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgPccId)
}

func RunPccList(c *builder.CommandConfig) error {
	pccs, _, err := c.Pccs().List()
	if err != nil {
		return err
	}
	return c.Printer.Print(getPccPrint(nil, c, getPccs(pccs)))
}

func RunPccGet(c *builder.CommandConfig) error {
	u, _, err := c.Pccs().Get(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgPccId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getPccPrint(nil, c, getPcc(u)))
}

func RunPccGetPeers(c *builder.CommandConfig) error {
	u, _, err := c.Pccs().GetPeers(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgPccId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getPccPeerPrint(c, *u))
}

func RunPccCreate(c *builder.CommandConfig) error {
	name := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgPccName))
	description := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgPccDescription))
	newUser := resources.PrivateCrossConnect{
		PrivateCrossConnect: ionoscloud.PrivateCrossConnect{
			Properties: &ionoscloud.PrivateCrossConnectProperties{
				Name:        &name,
				Description: &description,
			},
		},
	}
	u, resp, err := c.Pccs().Create(newUser)
	if err != nil {
		return err
	}
	return c.Printer.Print(getPccPrint(resp, c, getPcc(u)))
}

func RunPccUpdate(c *builder.CommandConfig) error {
	oldPcc, resp, err := c.Pccs().Get(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgPccId)))
	if err != nil {
		return err
	}
	newProperties := getPccInfo(oldPcc, c)
	pccUpd, resp, err := c.Pccs().Update(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgPccId)), *newProperties)
	if err != nil {
		return err
	}
	return c.Printer.Print(getPccPrint(resp, c, getPcc(pccUpd)))
}

func RunPccDelete(c *builder.CommandConfig) error {
	err := utils.AskForConfirm(c.Stdin, c.Printer, "delete private cross-connect")
	if err != nil {
		return err
	}
	resp, err := c.Pccs().Delete(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgPccId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getPccPrint(resp, c, nil))
}

func getPccInfo(oldUser *resources.PrivateCrossConnect, c *builder.CommandConfig) *resources.PrivateCrossConnectProperties {
	var namePcc, description string
	if properties, ok := oldUser.GetPropertiesOk(); ok && properties != nil {
		if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgPccName)) {
			namePcc = viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgPccName))
		} else {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				namePcc = *name
			}
		}
		if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgPccDescription)) {
			description = viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgPccDescription))
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

// Output Printing

var defaultPccCols = []string{"PccId", "Name", "Description"}

type PccPrint struct {
	PccId       string `json:"PccId,omitempty"`
	Name        string `json:"Name,omitempty"`
	Description string `json:"Description,omitempty"`
}

func getPccPrint(resp *resources.Response, c *builder.CommandConfig, pccs []resources.PrivateCrossConnect) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.ParentName
			r.Verb = c.Name
			r.WaitFlag = viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait))
		}
		if pccs != nil {
			r.OutputJSON = pccs
			r.KeyValue = getPccsKVMaps(pccs)
			r.Columns = getPccCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), os.Stderr)
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

func getPccPeerPrint(c *builder.CommandConfig, pccs []resources.Peer) printer.Result {
	r := printer.Result{}
	if c != nil {
		if pccs != nil {
			r.OutputJSON = pccs
			r.KeyValue = getPccPeersKVMaps(pccs)
			r.Columns = defaultPccPeersCols
		}
	}
	return r
}

func getPccCols(flagName string, outErr io.Writer) []string {
	if viper.IsSet(flagName) {
		var pccCols []string
		columnsMap := map[string]string{
			"PccId":       "PccId",
			"Name":        "Name",
			"Description": "Description",
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
		viper.GetString(config.ArgServerUrl),
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
