package cloudapi_v5

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/multierr"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v5/completer"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v5/query"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v5/waiter"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/internal/utils"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	cloudapiv5 "github.com/ionos-cloud/ionosctl/services/cloudapi-v5"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v5/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const backupUnitNote = "NOTE: To login with backup agent use: https://backup.ionos.com, with CONTRACT_NUMBER-BACKUP_UNIT_NAME and BACKUP_UNIT_PASSWORD!"

func BackupunitCmd() *core.Command {
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
	globalFlags.StringSliceP(config.ArgCols, "", defaultBackupUnitCols, printer.ColsMessage(defaultBackupUnitCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(backupUnitCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = backupUnitCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultBackupUnitCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, backupUnitCmd, core.CommandBuilder{
		Namespace:  "backupunit",
		Resource:   "backupunit",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List BackupUnits",
		LongDesc:   "Use this command to get a list of existing BackupUnits available on your account.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.BackupUnitsFiltersUsage(),
		Example:    listBackupUnitsExample,
		PreCmdRun:  PreRunBackupUnitList,
		CmdRun:     RunBackupUnitList,
		InitClient: true,
	})
	list.AddIntFlag(cloudapiv5.ArgMaxResults, cloudapiv5.ArgMaxResultsShort, 0, "The maximum number of elements to return")
	list.AddStringFlag(cloudapiv5.ArgOrderBy, "", "", "Limits results to those containing a matching value for a specific property")
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupUnitsFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv5.ArgFilters, cloudapiv5.ArgFiltersShort, []string{""}, "Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2")
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupUnitsFilters(), cobra.ShellCompDirectiveNoFileComp
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
	get.AddStringFlag(cloudapiv5.ArgBackupUnitId, cloudapiv5.ArgIdShort, "", cloudapiv5.BackupUnitId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgBackupUnitId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupUnitsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
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
	getsso.AddStringFlag(cloudapiv5.ArgBackupUnitId, cloudapiv5.ArgIdShort, "", cloudapiv5.BackupUnitId, core.RequiredFlagOption())
	_ = getsso.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgBackupUnitId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupUnitsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
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
	create.AddStringFlag(cloudapiv5.ArgName, cloudapiv5.ArgNameShort, "", "Alphanumeric name you want to assign to the BackupUnit", core.RequiredFlagOption())
	create.AddStringFlag(cloudapiv5.ArgEmail, cloudapiv5.ArgEmailShort, "", "The e-mail address you want to assign to the BackupUnit", core.RequiredFlagOption())
	create.AddStringFlag(cloudapiv5.ArgPassword, cloudapiv5.ArgPasswordShort, "", "Alphanumeric password you want to assign to the BackupUnit", core.RequiredFlagOption())
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
	update.AddStringFlag(cloudapiv5.ArgPassword, cloudapiv5.ArgPasswordShort, "", "Alphanumeric password you want to update for the BackupUnit")
	update.AddStringFlag(cloudapiv5.ArgEmail, cloudapiv5.ArgEmailShort, "", "The e-mail address you want to update for the BackupUnit")
	update.AddStringFlag(cloudapiv5.ArgBackupUnitId, cloudapiv5.ArgIdShort, "", cloudapiv5.BackupUnitId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgBackupUnitId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupUnitsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
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
		PreCmdRun:  PreRunBackupUnitDelete,
		CmdRun:     RunBackupUnitDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(cloudapiv5.ArgBackupUnitId, cloudapiv5.ArgIdShort, "", cloudapiv5.BackupUnitId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgBackupUnitId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupUnitsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for BackupUnit deletion to be executed")
	deleteCmd.AddBoolFlag(cloudapiv5.ArgAll, cloudapiv5.ArgAllShort, false, "Delete all BackupUnits.")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for BackupUnit deletion [seconds]")

	return backupUnitCmd
}

func PreRunBackupUnitId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv5.ArgBackupUnitId)
}

func PreRunBackupUnitDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv5.ArgBackupUnitId},
		[]string{cloudapiv5.ArgAll},
	)
}

func PreRunBackupUnitList(c *core.PreCommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgFilters)) {
		return query.ValidateFilters(c, completer.BackupUnitsFilters(), completer.BackupUnitsFiltersUsage())
	}
	return nil
}

func PreRunBackupUnitNameEmailPwd(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv5.ArgName, cloudapiv5.ArgEmail, cloudapiv5.ArgPassword)
}

func RunBackupUnitList(c *core.CommandConfig) error {
	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	if !structs.IsZero(listQueryParams) {
		c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(listQueryParams))
	}
	backupUnits, resp, err := c.CloudApiV5Services.BackupUnit().List(listQueryParams)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getBackupUnitPrint(nil, c, getBackupUnits(backupUnits)))
}

func RunBackupUnitGet(c *core.CommandConfig) error {
	c.Printer.Verbose("Backup unit with id: %v is getting... ", viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgBackupUnitId)))
	u, resp, err := c.CloudApiV5Services.BackupUnit().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgBackupUnitId)))
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getBackupUnitPrint(nil, c, getBackupUnit(u)))
}

func RunBackupUnitGetSsoUrl(c *core.CommandConfig) error {
	c.Printer.Verbose("Backup unit with id: %v is getting... ", viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgBackupUnitId)))
	u, resp, err := c.CloudApiV5Services.BackupUnit().GetSsoUrl(viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgBackupUnitId)))
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getBackupUnitSSOPrint(c, u))
}

func RunBackupUnitCreate(c *core.CommandConfig) error {
	name := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgName))
	email := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgEmail))
	pwd := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgPassword))
	newBackupUnit := resources.BackupUnit{
		BackupUnit: ionoscloud.BackupUnit{
			Properties: &ionoscloud.BackupUnitProperties{
				Name:     &name,
				Email:    &email,
				Password: &pwd,
			},
		},
	}
	c.Printer.Verbose("Properties set for creating the Backup Unit: Name: %v , Email: %v", name, email)
	u, resp, err := c.CloudApiV5Services.BackupUnit().Create(newBackupUnit)
	if resp != nil {
		c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	err = c.Printer.Print(backupUnitNote)
	if err != nil {
		return err
	}
	return c.Printer.Print(getBackupUnitPrint(resp, c, getBackupUnit(u)))
}

func RunBackupUnitUpdate(c *core.CommandConfig) error {
	newProperties := getBackupUnitInfo(c)
	backupUnitUpd, resp, err := c.CloudApiV5Services.BackupUnit().Update(viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgBackupUnitId)), *newProperties)
	if resp != nil {
		c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getBackupUnitPrint(resp, c, getBackupUnit(backupUnitUpd)))
}

func RunBackupUnitDelete(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv5.ArgAll)) {
		if err := DeleteAllBackupUnits(c); err != nil {
			return err
		}
		return c.Printer.Print(printer.Result{Resource: c.Resource, Verb: c.Verb})
	} else {
		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete backup unit"); err != nil {
			return err
		}
		c.Printer.Verbose("Starting deleting Backup unit with id: %v...", viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgBackupUnitId)))
		resp, err := c.CloudApiV5Services.BackupUnit().Delete(viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgBackupUnitId)))
		if resp != nil {
			c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
		}
		if err != nil {
			return err
		}
		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
			return err
		}
		return c.Printer.Print(getBackupUnitPrint(resp, c, nil))
	}
}

func getBackupUnitInfo(c *core.CommandConfig) *resources.BackupUnitProperties {
	var properties resources.BackupUnitProperties
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgPassword)) {
		pwd := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgPassword))
		properties.SetPassword(pwd)
		c.Printer.Verbose("Property Password set")
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgEmail)) {
		email := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgEmail))
		properties.SetEmail(email)
		c.Printer.Verbose("Property Email set: %v", email)
	}
	return &properties
}

func DeleteAllBackupUnits(c *core.CommandConfig) error {
	_ = c.Printer.Print("Backup Units to be deleted:")
	backupUnits, _, err := c.CloudApiV5Services.BackupUnit().List(resources.ListQueryParams{})
	if err != nil {
		return err
	}
	if backupUnitsItems, ok := backupUnits.GetItemsOk(); ok && backupUnitsItems != nil {
		for _, backupUnit := range *backupUnitsItems {
			var messageLog string
			if id, ok := backupUnit.GetIdOk(); ok && id != nil {
				messageLog = fmt.Sprintf("Backup Unit Id: %v", *id)
			}
			if properties, ok := backupUnit.GetPropertiesOk(); ok && properties != nil {
				if name, ok := properties.GetNameOk(); ok && name != nil {
					messageLog = fmt.Sprintf("%v Backup Unit Name: %v", messageLog, *name)
				}
			}
			_ = c.Printer.Print(messageLog)
		}
		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete all the Backup Units"); err != nil {
			return err
		}
		c.Printer.Verbose("Deleting all the BackupUnits...")
		var multiErr error
		for _, backupUnit := range *backupUnitsItems {
			if id, ok := backupUnit.GetIdOk(); ok && id != nil {
				c.Printer.Verbose("Starting deleting Backup Unit with id: %v...", *id)
				resp, err := c.CloudApiV5Services.BackupUnit().Delete(*id)
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
					return err
				}
			}
		}
		if multiErr != nil {
			return multiErr
		}
		return nil
	} else {
		return errors.New("could not get items of Backup Units")
	}
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

func getBackupUnitPrint(resp *resources.Response, c *core.CommandConfig, backupUnits []resources.BackupUnit) printer.Result {
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

func getBackupUnitSSOPrint(c *core.CommandConfig, backupUnit *resources.BackupUnitSSO) printer.Result {
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

func getBackupUnits(backupUnits resources.BackupUnits) []resources.BackupUnit {
	u := make([]resources.BackupUnit, 0)
	if items, ok := backupUnits.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			u = append(u, resources.BackupUnit{BackupUnit: item})
		}
	}
	return u
}

func getBackupUnit(u *resources.BackupUnit) []resources.BackupUnit {
	backupUnits := make([]resources.BackupUnit, 0)
	if u != nil {
		backupUnits = append(backupUnits, resources.BackupUnit{BackupUnit: u.BackupUnit})
	}
	return backupUnits
}

func getBackupUnitsKVMaps(us []resources.BackupUnit) []map[string]interface{} {
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

func getBackupUnitsSSOKVMaps(u *resources.BackupUnitSSO) []map[string]interface{} {
	out := make([]map[string]interface{}, 0)
	var uPrint BackupUnitPrint
	if url, ok := u.GetSsoUrlOk(); ok && url != nil {
		uPrint.BackupUnitSsoUrl = *url
	}
	o := structs.Map(uPrint)
	out = append(out, o)
	return out
}
