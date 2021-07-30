package commands

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func userS3key() *core.Command {
	ctx := context.TODO()
	s3keyCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "s3key",
			Aliases:          []string{"k", "s3k"},
			Short:            "User S3Key Operations",
			Long:             "The sub-commands of `ionosctl user s3key` allow you to see information, to list, get, create, update, delete Users S3Keys.",
			TraverseChildren: true,
		},
	}
	globalFlags := s3keyCmd.GlobalFlags()
	globalFlags.StringSliceP(config.ArgCols, "", defaultS3KeyCols, utils.ColsMessage(defaultS3KeyCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(s3keyCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = s3keyCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultS3KeyCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, s3keyCmd, core.CommandBuilder{
		Namespace:  "user",
		Resource:   "s3key",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List User S3Keys",
		LongDesc:   "Use this command to get a list of S3Keys of a specified User.\n\nRequired values to run command:\n\n* User Id",
		Example:    listS3KeysExample,
		PreCmdRun:  PreRunUserId,
		CmdRun:     RunUserS3KeyList,
		InitClient: true,
	})
	list.AddStringFlag(config.ArgUserId, "", "", config.RequiredFlagUserId)
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getUsersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, s3keyCmd, core.CommandBuilder{
		Namespace:  "user",
		Resource:   "s3key",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a User S3Key",
		LongDesc:   "Use this command to get information about a specified S3Key from a specified User.\n\nRequired values to run command:\n\n* User Id\n* S3Key Id",
		Example:    getS3KeyExample,
		PreCmdRun:  PreRunUserKeyIds,
		CmdRun:     RunUserS3KeyGet,
		InitClient: true,
	})
	get.AddStringFlag(config.ArgUserId, "", "", config.RequiredFlagUserId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getUsersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(config.ArgS3KeyId, config.ArgIdShort, "", config.RequiredFlagS3KeyId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgS3KeyId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getS3KeyIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, config.ArgUserId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, s3keyCmd, core.CommandBuilder{
		Namespace: "user",
		Resource:  "s3key",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a S3Key for a User",
		LongDesc: `Use this command to create a S3Key for a particular User.

Note: A maximum of five S3 keys may be created for any given user.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* User Id`,
		Example:    createS3KeyExample,
		PreCmdRun:  PreRunUserId,
		CmdRun:     RunUserS3KeyCreate,
		InitClient: true,
	})
	create.AddStringFlag(config.ArgUserId, "", "", config.RequiredFlagUserId)
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getUsersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for User S3Key creation to be executed")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for User S3Key creation [seconds]")

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, s3keyCmd, core.CommandBuilder{
		Namespace: "user",
		Resource:  "s3key",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a S3Key",
		LongDesc: `Use this command to update a specified S3Key from a particular User. This operation allows you to enable or disable a specific S3Key.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* User Id
* S3Key Id
* S3Key Active`,
		Example:    updateS3KeyExample,
		PreCmdRun:  PreRunUserKeyIds,
		CmdRun:     RunUserS3KeyUpdate,
		InitClient: true,
	})
	update.AddStringFlag(config.ArgUserId, "", "", config.RequiredFlagUserId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getUsersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(config.ArgS3KeyActive, "", false, "Enable or disable an User S3Key")
	update.AddStringFlag(config.ArgS3KeyId, config.ArgIdShort, "", config.RequiredFlagS3KeyId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgS3KeyId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getS3KeyIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, config.ArgUserId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for User S3Key update to be executed")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for User S3Key update [seconds]")

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, s3keyCmd, core.CommandBuilder{
		Namespace:  "user",
		Resource:   "s3key",
		Verb:       "delete",
		Aliases:    []string{"d"},
		ShortDesc:  "Delete a S3Key",
		LongDesc:   "Use this command to delete a specific S3Key of an User.\n\nRequired values to run command:\n\n* User Id\n* S3Key Id",
		Example:    deleteS3KeyExample,
		PreCmdRun:  PreRunUserKeyIds,
		CmdRun:     RunUserS3KeyDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(config.ArgUserId, "", "", config.RequiredFlagUserId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getUsersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(config.ArgS3KeyId, config.ArgIdShort, "", config.RequiredFlagS3KeyId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgS3KeyId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getS3KeyIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, config.ArgUserId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for Request for User S3Key deletion to be executed")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for User S3Key deletion [seconds]")

	return s3keyCmd
}

func PreRunUserKeyIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.NS, config.ArgUserId, config.ArgS3KeyId)
}

func RunUserS3KeyList(c *core.CommandConfig) error {
	ss, _, err := c.S3Keys().List(viper.GetString(core.GetFlagName(c.NS, config.ArgUserId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getS3KeyPrint(nil, c, getS3Keys(ss)))
}

func RunUserS3KeyGet(c *core.CommandConfig) error {
	c.Printer.Infof("S3 keys with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, config.ArgS3KeyId)))
	s, _, err := c.S3Keys().Get(viper.GetString(core.GetFlagName(c.NS, config.ArgUserId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgS3KeyId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getS3KeyPrint(nil, c, getS3Key(s)))
}

func RunUserS3KeyCreate(c *core.CommandConfig) error {
	userId := viper.GetString(core.GetFlagName(c.NS, config.ArgUserId))
	c.Printer.Infof("Properties set for creating the S3key: UserId: %v", userId)
	s, resp, err := c.S3Keys().Create(userId)
	if resp != nil {
		c.Printer.Infof("Request href: %v ", resp.Header.Get("location"))
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getS3KeyPrint(resp, c, getS3Key(s)))
}

func RunUserS3KeyUpdate(c *core.CommandConfig) error {
	active := viper.GetBool(core.GetFlagName(c.NS, config.ArgS3KeyActive))
	c.Printer.Infof("Property Active set: %v", active)
	newKey := resources.S3Key{
		S3Key: ionoscloud.S3Key{
			Properties: &ionoscloud.S3KeyProperties{
				Active: &active,
			},
		},
	}
	s, resp, err := c.S3Keys().Update(viper.GetString(core.GetFlagName(c.NS, config.ArgUserId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgS3KeyId)),
		newKey,
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getS3KeyPrint(resp, c, getS3Key(s)))
}

func RunUserS3KeyDelete(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete s3key"); err != nil {
		return err
	}
	c.Printer.Infof("S3 keys with id: %v is deleting...", viper.GetString(core.GetFlagName(c.NS, config.ArgS3KeyId)))
	resp, err := c.S3Keys().Delete(viper.GetString(core.GetFlagName(c.NS, config.ArgUserId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgS3KeyId)),
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
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

func getS3KeyPrint(resp *resources.Response, c *core.CommandConfig, s []resources.S3Key) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForRequest = viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForRequest))
		}
		if s != nil {
			r.OutputJSON = s
			r.KeyValue = getS3KeysKVMaps(s)
			r.Columns = getS3KeyCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
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
	err := config.Load()
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
