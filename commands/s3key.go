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

func s3key() *builder.Command {
	ctx := context.TODO()
	s3keyCmd := &builder.Command{
		Command: &cobra.Command{
			Use:              "s3key",
			Short:            "S3Key Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl s3key` + "`" + ` allow you to see information, to list, get, create, update, delete Users S3Keys. To view details about Users, check the ` + "`" + `ionosctl user` + "`" + ` commands.`,
			TraverseChildren: true,
		},
	}
	globalFlags := s3keyCmd.GlobalFlags()
	globalFlags.StringSlice(config.ArgCols, defaultS3KeyCols, "Columns to be printed in the standard output. You can also print SecretKey, using `--cols=\"S3KeyId,Active,SecretKey\"`")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(s3keyCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	globalFlags.String(config.ArgUserId, "", config.RequiredFlagUserId)
	_ = viper.BindPFlag(builder.GetGlobalFlagName(s3keyCmd.Name(), config.ArgUserId), globalFlags.Lookup(config.ArgUserId))
	_ = s3keyCmd.Command.RegisterFlagCompletionFunc(config.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getUsersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	builder.NewCommand(ctx, s3keyCmd, PreRunGlobalUserId, RunS3KeyList, "list", "List User S3Keys",
		"Use this command to get a list of S3Keys of a specified User.\n\nRequired values to run command:\n\n* User Id", listS3KeysExample, true)

	/*
		Get Command
	*/
	get := builder.NewCommand(ctx, s3keyCmd, PreRunGlobalUserIdKeyId, RunS3KeyGet, "get", "Get a User S3Key",
		"Use this command to get information about a specified S3Key from a specified User.\n\nRequired values to run command:\n\n* User Id\n* S3Key Id",
		getS3KeyExample, true)
	get.AddStringFlag(config.ArgS3KeyId, "", "", config.RequiredFlagS3KeyId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgS3KeyId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getS3KeyIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(s3keyCmd.Name(), config.ArgUserId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := builder.NewCommand(ctx, s3keyCmd, PreRunGlobalUserId, RunS3KeyCreate, "create", "Create a S3Key for a User",
		`Use this command to create a S3Key for a particular User. 

Note: A maximum of five S3 keys may be created for any given user.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option.

Required values to run command:
* User Id`, createS3KeyExample, true)
	create.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for S3Key to be created")
	create.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for a S3Key to be created [seconds]")

	/*
		Update Command
	*/
	update := builder.NewCommand(ctx, s3keyCmd, PreRunGlobalUserIdKeyIdActive, RunS3KeyUpdate, "update", "Update a S3Key",
		`Use this command to update a specified S3Key from a particular User. This operation allows you to enable or disable a specific S3Key.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option.

Required values to run command:
* User Id
* S3Key Id
* S3Key Active`, updateS3KeyExample, true)
	update.AddBoolFlag(config.ArgS3KeyActive, "", false, "Enable or disable a S3Key")
	update.AddStringFlag(config.ArgS3KeyId, "", "", config.RequiredFlagS3KeyId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgS3KeyId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getS3KeyIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(s3keyCmd.Name(), config.ArgUserId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for S3Key to be updated")
	update.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for a S3Key to be updated [seconds]")

	/*
		Delete Command
	*/
	deleteCmd := builder.NewCommand(ctx, s3keyCmd, PreRunGlobalUserIdKeyId, RunS3KeyDelete, "delete", "Delete a S3Key",
		"Use this command to delete a specific S3Key of an User.\n\nRequired values to run command:\n\n* User Id\n* S3Key Id",
		deleteS3KeyExample, true)
	deleteCmd.AddStringFlag(config.ArgS3KeyId, "", "", config.RequiredFlagS3KeyId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgS3KeyId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getS3KeyIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(s3keyCmd.Name(), config.ArgUserId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for S3Key to be deleted")
	deleteCmd.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for a S3Key to be deleted [seconds]")

	return s3keyCmd
}

func PreRunGlobalUserId(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredGlobalFlags(c.ParentName, config.ArgUserId)
}

func PreRunGlobalUserIdKeyId(c *builder.PreCommandConfig) error {
	var result *multierror.Error
	if err := builder.CheckRequiredGlobalFlags(c.ParentName, config.ArgUserId); err != nil {
		result = multierror.Append(result, err)
	}
	if err := builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgS3KeyId); err != nil {
		result = multierror.Append(result, err)
	}
	if result != nil {
		return result
	}
	return nil
}

func PreRunGlobalUserIdKeyIdActive(c *builder.PreCommandConfig) error {
	var result *multierror.Error
	if err := builder.CheckRequiredGlobalFlags(c.ParentName, config.ArgUserId); err != nil {
		result = multierror.Append(result, err)
	}
	if err := builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgS3KeyId, config.ArgS3KeyActive); err != nil {
		result = multierror.Append(result, err)
	}
	if result != nil {
		return result
	}
	return nil
}

func RunS3KeyList(c *builder.CommandConfig) error {
	ss, _, err := c.S3Keys().List(viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgUserId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getS3KeyPrint(nil, c, getS3Keys(ss)))
}

func RunS3KeyGet(c *builder.CommandConfig) error {
	s, _, err := c.S3Keys().Get(viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgUserId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgS3KeyId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getS3KeyPrint(nil, c, getS3Key(s)))
}

func RunS3KeyCreate(c *builder.CommandConfig) error {
	s, resp, err := c.S3Keys().Create(viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgUserId)))
	if err != nil {
		return err
	}

	if err = waitForAction(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getS3KeyPrint(resp, c, getS3Key(s)))
}

func RunS3KeyUpdate(c *builder.CommandConfig) error {
	active := viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgS3KeyActive))
	newKey := resources.S3Key{
		S3Key: ionoscloud.S3Key{
			Properties: &ionoscloud.S3KeyProperties{
				Active: &active,
			},
		},
	}
	s, resp, err := c.S3Keys().Update(viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgUserId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgS3KeyId)),
		newKey,
	)
	if err != nil {
		return err
	}

	if err = waitForAction(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getS3KeyPrint(resp, c, getS3Key(s)))
}

func RunS3KeyDelete(c *builder.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete s3key"); err != nil {
		return err
	}
	resp, err := c.S3Keys().Delete(viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgUserId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgS3KeyId)),
	)
	if err != nil {
		return err
	}

	if err = waitForAction(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getS3KeyPrint(resp, c, nil))
}

// Output Printing

var defaultS3KeyCols = []string{"S3KeyId", "Active"}

type S3KeyPrint struct {
	S3KeyId   string `json:"S3KeyId,omitempty"`
	Active    bool   `json:"Active,omitempty"`
	SecretKey string `json:"SecretKey,omitempty"`
}

func getS3KeyPrint(resp *resources.Response, c *builder.CommandConfig, s []resources.S3Key) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.ParentName
			r.Verb = c.Name
			r.WaitFlag = viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait))
		}
		if s != nil {
			r.OutputJSON = s
			r.KeyValue = getS3KeysKVMaps(s)
			r.Columns = getS3KeyCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getS3KeyCols(flagName string, outErr io.Writer) []string {
	if viper.IsSet(flagName) {
		var keyCols []string
		columnsMap := map[string]string{
			"S3KeyId":   "S3KeyId",
			"Active":    "Active",
			"SecretKey": "SecretKey",
		}
		for _, k := range viper.GetStringSlice(flagName) {
			col := columnsMap[k]
			if col != "" {
				keyCols = append(keyCols, col)
			} else {
				clierror.CheckError(errors.New("unknown column "+k), outErr)
			}
		}
		return keyCols
	} else {
		return defaultS3KeyCols
	}
}

func getS3Keys(S3Keys resources.S3Keys) []resources.S3Key {
	ss := make([]resources.S3Key, 0)
	if items, ok := S3Keys.GetItemsOk(); ok && items != nil {
		for _, s := range *items {
			ss = append(ss, resources.S3Key{S3Key: s})
		}
	}
	return ss
}

func getS3Key(s *resources.S3Key) []resources.S3Key {
	ss := make([]resources.S3Key, 0)
	if s != nil {
		ss = append(ss, resources.S3Key{S3Key: s.S3Key})
	}
	return ss
}

func getS3KeysKVMaps(ss []resources.S3Key) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(ss))
	for _, s := range ss {
		o := getS3KeyKVMap(s)
		out = append(out, o)
	}
	return out
}

func getS3KeyKVMap(s resources.S3Key) map[string]interface{} {
	var ssPrint S3KeyPrint
	if ssId, ok := s.GetIdOk(); ok && ssId != nil {
		ssPrint.S3KeyId = *ssId
	}
	if properties, ok := s.GetPropertiesOk(); ok && properties != nil {
		if active, ok := properties.GetActiveOk(); ok && active != nil {
			ssPrint.Active = *active
		}
		if secretKey, ok := properties.GetSecretKeyOk(); ok && secretKey != nil {
			ssPrint.SecretKey = *secretKey
		}
	}
	return structs.Map(ssPrint)
}

func getS3KeyIds(outErr io.Writer, userId string) []string {
	err := config.LoadFile()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)
	S3KeySvc := resources.NewS3KeyService(clientSvc.Get(), context.TODO())
	S3Keys, _, err := S3KeySvc.List(userId)
	clierror.CheckError(err, outErr)
	ssIds := make([]string, 0)
	if items, ok := S3Keys.S3Keys.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				ssIds = append(ssIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return ssIds
}
