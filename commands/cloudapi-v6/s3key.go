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

func UserS3keyCmd() *core.Command {
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
	globalFlags.StringSliceP(config.ArgCols, "", defaultS3KeyCols, printer.ColsMessage(defaultS3KeyCols))
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
	list.AddStringFlag(cloudapiv6.ArgUserId, "", "", cloudapiv6.UserId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.UsersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddBoolFlag(config.ArgNoHeaders, "", false, "When using text output, don't print headers")

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
	get.AddStringFlag(cloudapiv6.ArgUserId, "", "", cloudapiv6.UserId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.UsersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapiv6.ArgS3KeyId, cloudapiv6.ArgIdShort, "", cloudapiv6.S3KeyId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgS3KeyId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.S3KeyIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgUserId))), cobra.ShellCompDirectiveNoFileComp
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
	create.AddStringFlag(cloudapiv6.ArgUserId, "", "", cloudapiv6.UserId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.UsersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
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
	update.AddStringFlag(cloudapiv6.ArgUserId, "", "", cloudapiv6.UserId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.UsersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(cloudapiv6.ArgS3KeyActive, "", false, "Enable or disable an User S3Key. E.g.: --s3key-active=true, --s3key-active=false")
	update.AddStringFlag(cloudapiv6.ArgS3KeyId, cloudapiv6.ArgIdShort, "", cloudapiv6.S3KeyId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgS3KeyId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.S3KeyIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgUserId))), cobra.ShellCompDirectiveNoFileComp
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
		PreCmdRun:  PreRunUserKeyDelete,
		CmdRun:     RunUserS3KeyDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(cloudapiv6.ArgUserId, "", "", cloudapiv6.UserId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.UsersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(cloudapiv6.ArgS3KeyId, cloudapiv6.ArgIdShort, "", cloudapiv6.S3KeyId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgS3KeyId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.S3KeyIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.ArgUserId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for Request for User S3Key deletion to be executed")
	deleteCmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete all the S3Keys of an User.")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for User S3Key deletion [seconds]")

	return s3keyCmd
}

func PreRunUserKeyIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgUserId, cloudapiv6.ArgS3KeyId)
}

func PreRunUserKeyDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgUserId, cloudapiv6.ArgS3KeyId},
		[]string{cloudapiv6.ArgUserId, cloudapiv6.ArgAll},
	)
}

func RunUserS3KeyList(c *core.CommandConfig) error {
	ss, resp, err := c.CloudApiV6Services.S3Keys().List(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId)))
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getS3KeyPrint(nil, c, getS3Keys(ss)))
}

func RunUserS3KeyGet(c *core.CommandConfig) error {
	c.Printer.Verbose("S3 keys with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgS3KeyId)))
	s, resp, err := c.CloudApiV6Services.S3Keys().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgS3KeyId)),
	)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getS3KeyPrint(nil, c, getS3Key(s)))
}

func RunUserS3KeyCreate(c *core.CommandConfig) error {
	userId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId))
	c.Printer.Verbose("Creating S3 Key for User with ID: %v", userId)
	s, resp, err := c.CloudApiV6Services.S3Keys().Create(userId)
	if resp != nil && printer.GetId(resp) != "" {
		c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getS3KeyPrint(resp, c, getS3Key(s)))
}

func RunUserS3KeyUpdate(c *core.CommandConfig) error {
	active := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgS3KeyActive))
	c.Printer.Verbose("Property Active set: %v", active)
	newKey := resources.S3Key{
		S3Key: ionoscloud.S3Key{
			Properties: &ionoscloud.S3KeyProperties{
				Active: &active,
			},
		},
	}
	userId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId))
	keyId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgS3KeyId))
	c.Printer.Verbose("Creating S3 Key with ID: %v for User with ID: %v", keyId, userId)
	s, resp, err := c.CloudApiV6Services.S3Keys().Update(userId, keyId, newKey)
	if resp != nil && printer.GetId(resp) != "" {
		c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getS3KeyPrint(resp, c, getS3Key(s)))
}

func RunUserS3KeyDelete(c *core.CommandConfig) error {
	userId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId))
	s3KeyId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgS3KeyId))
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllS3keys(c); err != nil {
			return err
		}
		return c.Printer.Print(printer.Result{Resource: c.Resource, Verb: c.Verb})
	} else {
		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete s3key"); err != nil {
			return err
		}
		c.Printer.Verbose("User ID: %v", userId)
		c.Printer.Verbose("Starting deleting S3 Key with ID: %v...", s3KeyId)
		resp, err := c.CloudApiV6Services.S3Keys().Delete(userId, s3KeyId)
		if resp != nil && printer.GetId(resp) != "" {
			c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
		}
		if err != nil {
			return err
		}
		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
			return err
		}
		return c.Printer.Print(getS3KeyPrint(resp, c, nil))
	}
}

func DeleteAllS3keys(c *core.CommandConfig) error {
	userId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId))
	c.Printer.Verbose("User ID: %v", userId)
	c.Printer.Verbose("Getting S3 Keys...")
	s3Keys, resp, err := c.CloudApiV6Services.S3Keys().List(userId)
	if err != nil {
		return err
	}
	if s3KeysItems, ok := s3Keys.GetItemsOk(); ok && s3KeysItems != nil {
		if len(*s3KeysItems) > 0 {
			_ = c.Printer.Print("S3 keys to be deleted:")
			for _, s3Key := range *s3KeysItems {
				if id, ok := s3Key.GetIdOk(); ok && id != nil {
					_ = c.Printer.Print("S3 key Id: " + *id)
				}
			}
			if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete all the S3Keys"); err != nil {
				return err
			}
			c.Printer.Verbose("Deleting all the S3Keys...")
			var multiErr error
			for _, s3Key := range *s3KeysItems {
				if id, ok := s3Key.GetIdOk(); ok && id != nil {
					c.Printer.Verbose("Staring deleting S3 keys with id: %v...", *id)
					resp, err = c.CloudApiV6Services.S3Keys().Delete(userId, *id)
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
			return errors.New("no S3 Keys found")
		}
	} else {
		return errors.New("could not get items of S3 Keys")
	}
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
