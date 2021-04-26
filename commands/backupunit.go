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

const backupUnitNote = "NOTE: To login with backup agent use: https://backup.ionos.com, with CONTRACT_NUMBER-BACKUP_UNIT_NAME and BACKUP_UNIT_PASSWORD!"

func backupunit() *builder.Command {
	ctx := context.TODO()
	backupUnitCmd := &builder.Command{
		Command: &cobra.Command{
			Use:              "backupunit",
			Short:            "BackupUnit Operations",
			Long:             `The sub-command of ` + "`" + `ionosctl backupunit` + "`" + ` allows you to list, get, create, update, delete BackupUnits under your account.`,
			TraverseChildren: true,
		},
	}
	globalFlags := backupUnitCmd.GlobalFlags()
	globalFlags.StringSlice(config.ArgCols, defaultBackupUnitCols, "Columns to be printed in the standard output")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(backupUnitCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))

	builder.NewCommand(ctx, backupUnitCmd, noPreRun, RunBackupUnitList, "list", "List BackupUnits",
		"Use this command to get a list of existing BackupUnits available on your account.", listBackupUnitsExample, true)

	get := builder.NewCommand(ctx, backupUnitCmd, PreRunBackupUnitIdValidate, RunBackupUnitGet, "get", "Get a BackupUnit",
		"Use this command to retrieve details about a specific BackupUnit.\n\nRequired values to run command:\n\n* BackupUnit Id", getBackupUnitExample, true)
	get.AddStringFlag(config.ArgBackupUnitId, "", "", config.RequiredFlagBackupUnitId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgBackupUnitId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getBackupUnitsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	getsso := builder.NewCommand(ctx, backupUnitCmd, PreRunBackupUnitIdValidate, RunBackupUnitGetSsoUrl, "get-sso-url", "Get BackupUnit SSO URL",
		"Use this command to access the GUI with a Single Sign On (SSO) URL that can be retrieved from the Cloud API using this request. If you copy the entire value returned and paste it into a browser, you will be logged into the BackupUnit GUI.\n\nRequired values to run command:\n\n* BackupUnit Id",
		getBackupUnitSSOExample, true)
	getsso.AddStringFlag(config.ArgBackupUnitId, "", "", config.RequiredFlagBackupUnitId)
	_ = getsso.Command.RegisterFlagCompletionFunc(config.ArgBackupUnitId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getBackupUnitsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	create := builder.NewCommand(ctx, backupUnitCmd, PreRunBackupUnitNameEmailPwdValidate, RunBackupUnitCreate, "create", "Create a BackupUnit",
		`Use this command to create a BackupUnit under a particular contract. You need to specify the name, email and password for the new BackupUnit.

Notes:

* The name assigned to the BackupUnit will be concatenated with the contract number to create the login name for the backup system. The name may NOT be changed after creation.
* The password set here is used along with the login name described above to register the backup agent with the backup system. When setting the password, please make a note of it, as the value cannot be retrieved using the Cloud API.
* The e-mail address supplied here does NOT have to be the same as your Cloud API username. This e-mail address will receive service reports from the backup system.
* To login to backup agent, please use https://dcd.ionos.com/latest/ and access BackupUnit Console or use https://backup.ionos.com

Required values to run a command:

* BackupUnit Name
* BackupUnit Email
* BackupUnit Password`, createBackupUnitExample, true)
	create.AddStringFlag(config.ArgBackupUnitName, "", "", "Alphanumeric name you want to assign to the BackupUnit "+config.RequiredFlag)
	create.AddStringFlag(config.ArgBackupUnitEmail, "", "", "The e-mail address you want to assign to the BackupUnit "+config.RequiredFlag)
	create.AddStringFlag(config.ArgBackupUnitPassword, "", "", "Alphanumeric password you want to assign to the BackupUnit "+config.RequiredFlag)
	create.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for BackupUnit to be created")
	create.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for BackupUnit to be created [seconds]")

	update := builder.NewCommand(ctx, backupUnitCmd, PreRunBackupUnitIdValidate, RunBackupUnitUpdate, "update", "Update a BackupUnit",
		`Use this command to update details about a specific BackupUnit. The password and the email may be updated.

Required values to run command:

* BackupUnit Id`, updateBackupUnitExample, true)
	update.AddStringFlag(config.ArgBackupUnitPassword, "", "", "Alphanumeric password you want to update for the BackupUnit")
	update.AddStringFlag(config.ArgBackupUnitEmail, "", "", "The e-mail address you want to update for the BackupUnit")
	update.AddStringFlag(config.ArgBackupUnitId, "", "", config.RequiredFlagBackupUnitId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgBackupUnitId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getBackupUnitsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for BackupUnit to be updated")
	update.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for BackupUnit to be updated [seconds]")

	deleteCmd := builder.NewCommand(ctx, backupUnitCmd, PreRunBackupUnitIdValidate, RunBackupUnitDelete, "delete", "Delete a BackupUnit",
		`Use this command to delete a BackupUnit. Deleting a BackupUnit is a dangerous operation. A successful DELETE will remove the backup plans inside a BackupUnit, ALL backups associated with the BackupUnit, the backup user and finally the BackupUnit itself.

Required values to run command:

* BackupUnit Id`, deleteBackupUnitExample, true)
	deleteCmd.AddStringFlag(config.ArgBackupUnitId, "", "", config.RequiredFlagBackupUnitId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgBackupUnitId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getBackupUnitsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for BackupUnit to be deleted")
	deleteCmd.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for BackupUnit to be deleted [seconds]")

	return backupUnitCmd
}

func PreRunBackupUnitIdValidate(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgBackupUnitId)
}

func PreRunBackupUnitNameEmailPwdValidate(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgBackupUnitName, config.ArgBackupUnitEmail, config.ArgBackupUnitPassword)
}

func RunBackupUnitList(c *builder.CommandConfig) error {
	backupUnits, _, err := c.BackupUnit().List()
	if err != nil {
		return err
	}
	return c.Printer.Print(getBackupUnitPrint(nil, c, getBackupUnits(backupUnits)))
}

func RunBackupUnitGet(c *builder.CommandConfig) error {
	u, _, err := c.BackupUnit().Get(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgBackupUnitId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getBackupUnitPrint(nil, c, getBackupUnit(u)))
}

func RunBackupUnitGetSsoUrl(c *builder.CommandConfig) error {
	u, _, err := c.BackupUnit().GetSsoUrl(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgBackupUnitId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getBackupUnitSSOPrint(c, u))
}

func RunBackupUnitCreate(c *builder.CommandConfig) error {
	name := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgBackupUnitName))
	email := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgBackupUnitEmail))
	pwd := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgBackupUnitPassword))
	newBackupUnit := resources.BackupUnit{
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
	if err = waitForAction(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	c.Printer.Print(backupUnitNote)
	return c.Printer.Print(getBackupUnitPrint(resp, c, getBackupUnit(u)))
}

func RunBackupUnitUpdate(c *builder.CommandConfig) error {
	newProperties := getBackupUnitInfo(c)
	backupUnitUpd, resp, err := c.BackupUnit().Update(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgBackupUnitId)), *newProperties)
	if err != nil {
		return err
	}
	if err = waitForAction(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getBackupUnitPrint(resp, c, getBackupUnit(backupUnitUpd)))
}

func RunBackupUnitDelete(c *builder.CommandConfig) error {
	err := utils.AskForConfirm(c.Stdin, c.Printer, "delete backup unit")
	if err != nil {
		return err
	}
	resp, err := c.BackupUnit().Delete(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgBackupUnitId)))
	if err != nil {
		return err
	}
	if err = waitForAction(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getBackupUnitPrint(resp, c, nil))
}

func getBackupUnitInfo(c *builder.CommandConfig) *resources.BackupUnitProperties {
	var properties resources.BackupUnitProperties
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgBackupUnitPassword)) {
		pwd := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgBackupUnitPassword))
		properties.SetPassword(pwd)
	}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgBackupUnitEmail)) {
		email := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgBackupUnitEmail))
		properties.SetEmail(email)
	}
	return &properties
}

// Output Printing

var defaultBackupUnitCols = []string{"BackupUnitId", "Name", "Email"}

type BackupUnitPrint struct {
	BackupUnitId     string `json:"BackupUnitId,omitempty"`
	Name             string `json:"Name,omitempty"`
	Email            string `json:"Email,omitempty"`
	BackupUnitSsoUrl string `json:"BackupUnitSsoUrl,omitempty"`
}

func getBackupUnitPrint(resp *resources.Response, c *builder.CommandConfig, backupUnits []resources.BackupUnit) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.ParentName
			r.Verb = c.Name
			r.WaitFlag = viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait))
		}
		if backupUnits != nil {
			r.OutputJSON = backupUnits
			r.KeyValue = getBackupUnitsKVMaps(backupUnits)
			r.Columns = getBackupUnitCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getBackupUnitSSOPrint(c *builder.CommandConfig, backupUnit *resources.BackupUnitSSO) printer.Result {
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

func getBackupUnitsIds(outErr io.Writer) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)
	backupUnitSvc := resources.NewBackupUnitService(clientSvc.Get(), context.TODO())
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
