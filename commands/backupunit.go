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

const backupUnitNote = "NOTE: To login with backup agent use: https://backup.ionos.com, with CONTRACT_NUMBER-BACKUP_UNIT_NAME and BACKUP_UNIT_PASSWORD!"

func backupunit() *core.Command {
	ctx := context.TODO()
	backupUnitCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "backupunit",
			Aliases:          []string{"b", "backup"},
			Short:            "BackupUnit Operations",
			Long:             "The sub-commands of `ionosctl backupunit` allow you to list, get, create, update, delete BackupUnits under your account.",
			TraverseChildren: true,
		},
	}
	globalFlags := backupUnitCmd.GlobalFlags()
	globalFlags.StringSliceP(config.ArgCols, "", defaultBackupUnitCols, utils.ColsMessage(defaultBackupUnitCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(backupUnitCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = backupUnitCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultBackupUnitCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	core.NewCommand(ctx, backupUnitCmd, core.CommandBuilder{
		Namespace:  "backupunit",
		Resource:   "backupunit",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List BackupUnits",
		LongDesc:   "Use this command to get a list of existing BackupUnits available on your account.",
		Example:    listBackupUnitsExample,
		PreCmdRun:  noPreRun,
		CmdRun:     RunBackupUnitList,
		InitClient: true,
	})

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, backupUnitCmd, core.CommandBuilder{
		Namespace:  "backupunit",
		Resource:   "backupunit",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a BackupUnit",
		LongDesc:   "Use this command to retrieve details about a specific BackupUnit.\n\nRequired values to run command:\n\n* BackupUnit Id",
		Example:    getBackupUnitExample,
		PreCmdRun:  PreRunBackupUnitId,
		CmdRun:     RunBackupUnitGet,
		InitClient: true,
	})
	get.AddStringFlag(config.ArgBackupUnitId, config.ArgIdShort, "", config.RequiredFlagBackupUnitId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgBackupUnitId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getBackupUnitsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get SSO URL Command
	*/
	getsso := core.NewCommand(ctx, backupUnitCmd, core.CommandBuilder{
		Namespace:  "backupunit",
		Resource:   "backupunit",
		Verb:       "get-sso-url",
		ShortDesc:  "Get BackupUnit SSO URL",
		LongDesc:   "Use this command to access the GUI with a Single Sign On URL that can be retrieved from the Cloud API using this request. If you copy the entire value returned and paste it into a browser, you will be logged into the BackupUnit GUI.\n\nRequired values to run command:\n\n* BackupUnit Id",
		Example:    getBackupUnitSSOExample,
		PreCmdRun:  PreRunBackupUnitId,
		CmdRun:     RunBackupUnitGetSsoUrl,
		InitClient: true,
	})
	getsso.AddStringFlag(config.ArgBackupUnitId, config.ArgIdShort, "", config.RequiredFlagBackupUnitId)
	_ = getsso.Command.RegisterFlagCompletionFunc(config.ArgBackupUnitId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getBackupUnitsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, backupUnitCmd, core.CommandBuilder{
		Namespace: "backupunit",
		Resource:  "backupunit",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a BackupUnit",
		LongDesc: `Use this command to create a BackupUnit under a particular contract. You need to specify the name, email and password for the new BackupUnit.

Notes:

* The name assigned to the BackupUnit will be concatenated with the contract number to create the login name for the backup system. The name may NOT be changed after creation.
* The password set here is used along with the login name described above to register the backup agent with the backup system. When setting the password, please make a note of it, as the value cannot be retrieved using the Cloud API.
* The e-mail address supplied here does NOT have to be the same as your Cloud API username. This e-mail address will receive service reports from the backup system.
* To login to backup agent, please use [https://dcd.ionos.com/latest/](https://dcd.ionos.com/latest/) and access BackupUnit Console or use [https://backup.ionos.com](https://backup.ionos.com)

Required values to run a command:

* Name
* Email
* Password`,
		Example:    createBackupUnitExample,
		PreCmdRun:  PreRunBackupUnitNameEmailPwd,
		CmdRun:     RunBackupUnitCreate,
		InitClient: true,
	})
	create.AddStringFlag(config.ArgName, config.ArgNameShort, "", "Alphanumeric name you want to assign to the BackupUnit "+config.RequiredFlag)
	create.AddStringFlag(config.ArgEmail, config.ArgEmailShort, "", "The e-mail address you want to assign to the BackupUnit "+config.RequiredFlag)
	create.AddStringFlag(config.ArgPassword, config.ArgPasswordShort, "", "Alphanumeric password you want to assign to the BackupUnit "+config.RequiredFlag)
	create.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for BackupUnit creation to be executed")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for BackupUnit creation [seconds]")

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, backupUnitCmd, core.CommandBuilder{
		Namespace: "backupunit",
		Resource:  "backupunit",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a BackupUnit",
		LongDesc: `Use this command to update details about a specific BackupUnit. The password and the email may be updated.

Required values to run command:

* BackupUnit Id`,
		Example:    updateBackupUnitExample,
		PreCmdRun:  PreRunBackupUnitId,
		CmdRun:     RunBackupUnitUpdate,
		InitClient: true,
	})
	update.AddStringFlag(config.ArgPassword, config.ArgPasswordShort, "", "Alphanumeric password you want to update for the BackupUnit")
	update.AddStringFlag(config.ArgEmail, config.ArgEmailShort, "", "The e-mail address you want to update for the BackupUnit")
	update.AddStringFlag(config.ArgBackupUnitId, config.ArgIdShort, "", config.RequiredFlagBackupUnitId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgBackupUnitId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getBackupUnitsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for BackupUnit update to be executed")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for BackupUnit update [seconds]")

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, backupUnitCmd, core.CommandBuilder{
		Namespace: "backupunit",
		Resource:  "backupunit",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a BackupUnit",
		LongDesc: `Use this command to delete a BackupUnit. Deleting a BackupUnit is a dangerous operation. A successful DELETE will remove the backup plans inside a BackupUnit, ALL backups associated with the BackupUnit, the backup user and finally the BackupUnit itself.

Required values to run command:

* BackupUnit Id`,
		Example:    deleteBackupUnitExample,
		PreCmdRun:  PreRunBackupUnitId,
		CmdRun:     RunBackupUnitDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(config.ArgBackupUnitId, config.ArgIdShort, "", config.RequiredFlagBackupUnitId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgBackupUnitId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getBackupUnitsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for BackupUnit deletion to be executed")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for BackupUnit deletion [seconds]")

	return backupUnitCmd
}

func PreRunBackupUnitId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.NS, config.ArgBackupUnitId)
}

func PreRunBackupUnitNameEmailPwd(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.NS, config.ArgName, config.ArgEmail, config.ArgPassword)
}

func RunBackupUnitList(c *core.CommandConfig) error {
	backupUnits, _, err := c.BackupUnit().List()
	if err != nil {
		return err
	}
	return c.Printer.Print(getBackupUnitPrint(nil, c, getBackupUnits(backupUnits)))
}

func RunBackupUnitGet(c *core.CommandConfig) error {
	u, _, err := c.BackupUnit().Get(viper.GetString(core.GetFlagName(c.NS, config.ArgBackupUnitId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getBackupUnitPrint(nil, c, getBackupUnit(u)))
}

func RunBackupUnitGetSsoUrl(c *core.CommandConfig) error {
	u, _, err := c.BackupUnit().GetSsoUrl(viper.GetString(core.GetFlagName(c.NS, config.ArgBackupUnitId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getBackupUnitSSOPrint(c, u))
}

func RunBackupUnitCreate(c *core.CommandConfig) error {
	name := viper.GetString(core.GetFlagName(c.NS, config.ArgName))
	email := viper.GetString(core.GetFlagName(c.NS, config.ArgEmail))
	pwd := viper.GetString(core.GetFlagName(c.NS, config.ArgPassword))
	newBackupUnit := v5.BackupUnit{
		BackupUnit: ionoscloud.BackupUnit{
			Properties: &ionoscloud.BackupUnitProperties{
				Name:     &name,
				Email:    &email,
				Password: &pwd,
			},
		},
	}
	u, resp, err := c.BackupUnit().Create(newBackupUnit)
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	c.Printer.Print(backupUnitNote)
	return c.Printer.Print(getBackupUnitPrint(resp, c, getBackupUnit(u)))
}

func RunBackupUnitUpdate(c *core.CommandConfig) error {
	newProperties := getBackupUnitInfo(c)
	backupUnitUpd, resp, err := c.BackupUnit().Update(viper.GetString(core.GetFlagName(c.NS, config.ArgBackupUnitId)), *newProperties)
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getBackupUnitPrint(resp, c, getBackupUnit(backupUnitUpd)))
}

func RunBackupUnitDelete(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete backup unit"); err != nil {
		return err
	}
	resp, err := c.BackupUnit().Delete(viper.GetString(core.GetFlagName(c.NS, config.ArgBackupUnitId)))
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getBackupUnitPrint(resp, c, nil))
}

func getBackupUnitInfo(c *core.CommandConfig) *v5.BackupUnitProperties {
	var properties v5.BackupUnitProperties
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgPassword)) {
		pwd := viper.GetString(core.GetFlagName(c.NS, config.ArgPassword))
		properties.SetPassword(pwd)
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgEmail)) {
		email := viper.GetString(core.GetFlagName(c.NS, config.ArgEmail))
		properties.SetEmail(email)
	}
	return &properties
}

// Output Printing

var defaultBackupUnitCols = []string{"BackupUnitId", "Name", "Email", "State"}

type BackupUnitPrint struct {
	BackupUnitId     string `json:"BackupUnitId,omitempty"`
	Name             string `json:"Name,omitempty"`
	Email            string `json:"Email,omitempty"`
	BackupUnitSsoUrl string `json:"BackupUnitSsoUrl,omitempty"`
	State            string `json:"State,omitempty"`
}

func getBackupUnitPrint(resp *v5.Response, c *core.CommandConfig, backupUnits []v5.BackupUnit) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForRequest = viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForRequest))
		}
		if backupUnits != nil {
			r.OutputJSON = backupUnits
			r.KeyValue = getBackupUnitsKVMaps(backupUnits)
			r.Columns = getBackupUnitCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getBackupUnitSSOPrint(c *core.CommandConfig, backupUnit *v5.BackupUnitSSO) printer.Result {
	r := printer.Result{}
	if c != nil {
		if backupUnit != nil {
			r.OutputJSON = backupUnit
			r.KeyValue = getBackupUnitsSSOKVMaps(backupUnit)
			r.Columns = []string{"BackupUnitSsoUrl"}
		}
	}
	return r
}

func getBackupUnitCols(flagName string, outErr io.Writer) []string {
	if viper.IsSet(flagName) {
		var backupUnitCols []string
		columnsMap := map[string]string{
			"BackupUnitId": "BackupUnitId",
			"Name":         "Name",
			"Password":     "Password",
			"Email":        "Email",
			"State":        "State",
		}
		for _, k := range viper.GetStringSlice(flagName) {
			col := columnsMap[k]
			if col != "" {
				backupUnitCols = append(backupUnitCols, col)
			} else {
				clierror.CheckError(errors.New("unknown column "+k), outErr)
			}
		}
		return backupUnitCols
	} else {
		return defaultBackupUnitCols
	}
}

func getBackupUnits(backupUnits v5.BackupUnits) []v5.BackupUnit {
	u := make([]v5.BackupUnit, 0)
	if items, ok := backupUnits.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			u = append(u, v5.BackupUnit{BackupUnit: item})
		}
	}
	return u
}

func getBackupUnit(u *v5.BackupUnit) []v5.BackupUnit {
	backupUnits := make([]v5.BackupUnit, 0)
	if u != nil {
		backupUnits = append(backupUnits, v5.BackupUnit{BackupUnit: u.BackupUnit})
	}
	return backupUnits
}

func getBackupUnitsKVMaps(us []v5.BackupUnit) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(us))
	for _, u := range us {
		var uPrint BackupUnitPrint
		if id, ok := u.GetIdOk(); ok && id != nil {
			uPrint.BackupUnitId = *id
		}
		if properties, ok := u.GetPropertiesOk(); ok && properties != nil {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				uPrint.Name = *name
			}
			if email, ok := properties.GetEmailOk(); ok && email != nil {
				uPrint.Email = *email
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

func getBackupUnitsSSOKVMaps(u *v5.BackupUnitSSO) []map[string]interface{} {
	out := make([]map[string]interface{}, 0)
	var uPrint BackupUnitPrint
	if url, ok := u.GetSsoUrlOk(); ok && url != nil {
		uPrint.BackupUnitSsoUrl = *url
	}
	o := structs.Map(uPrint)
	out = append(out, o)
	return out
}

func getBackupUnitsIds(outErr io.Writer) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := v5.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		config.GetServerUrl(),
	)
	clierror.CheckError(err, outErr)
	backupUnitSvc := v5.NewBackupUnitService(clientSvc.Get(), context.TODO())
	backupUnits, _, err := backupUnitSvc.List()
	clierror.CheckError(err, outErr)
	backupUnitsIds := make([]string, 0)
	if items, ok := backupUnits.BackupUnits.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				backupUnitsIds = append(backupUnitsIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return backupUnitsIds
}
