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

func user() *core.Command {
	ctx := context.TODO()
	userCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "user",
			Short:            "User Operations",
			Long:             `The sub-command of ` + "`" + `ionosctl user` + "`" + ` allows you to list, get, create, update, delete Users under your account. To add Users to a Group, check the ` + "`" + `ionosctl group` + "`" + ` commands. To add S3Keys to a User, check the ` + "`" + `ionosctl s3key` + "`" + ` commands.`,
			TraverseChildren: true,
		},
	}
	globalFlags := userCmd.GlobalFlags()
	globalFlags.StringSlice(config.ArgCols, defaultUserCols, "Columns to be printed in the standard output")
	_ = viper.BindPFlag(core.GetGlobalFlagName(userCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = userCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultUserCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	core.NewCommand(ctx, userCmd, core.CommandBuilder{
		Namespace:  "user",
		Resource:   "user",
		Verb:       "list",
		ShortDesc:  "List Users",
		LongDesc:   "Use this command to get a list of existing Users available on your account.",
		Example:    listUserExample,
		PreCmdRun:  noPreRun,
		CmdRun:     RunUserList,
		InitClient: true,
	})

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, userCmd, core.CommandBuilder{
		Namespace:  "user",
		Resource:   "user",
		Verb:       "get",
		ShortDesc:  "Get a User",
		LongDesc:   "Use this command to retrieve details about a specific User.\n\nRequired values to run command:\n\n* User Id",
		Example:    getUserExample,
		PreCmdRun:  PreRunUserId,
		CmdRun:     RunUserGet,
		InitClient: true,
	})
	get.AddStringFlag(config.ArgUserId, "", "", config.RequiredFlagUserId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getUsersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, userCmd, core.CommandBuilder{
		Namespace: "user",
		Resource:  "user",
		Verb:      "create",
		ShortDesc: "Create a User under a particular contract",
		LongDesc: `Use this command to create a User under a particular contract. You need to specify the firstname, lastname, email and password for the new User.

Note: The password set here cannot be updated through the API currently. It is recommended that a new User log into the DCD and change their password.

Required values to run a command:

* User First Name
* User Last Name
* User Email
* User Password`,
		Example:    createUserExample,
		PreCmdRun:  PreRunUserNameEmailPwd,
		CmdRun:     RunUserCreate,
		InitClient: true,
	})
	create.AddStringFlag(config.ArgUserFirstName, "", "", "The firstname for the User "+config.RequiredFlag)
	create.AddStringFlag(config.ArgUserLastName, "", "", "The lastname for the User "+config.RequiredFlag)
	create.AddStringFlag(config.ArgUserEmail, "", "", "The email for the User "+config.RequiredFlag)
	create.AddStringFlag(config.ArgUserPassword, "", "", "The password for the User (must be at least 5 characters long) "+config.RequiredFlag)
	create.AddBoolFlag(config.ArgUserAdministrator, "", false, "Assigns the User to have administrative rights")
	create.AddBoolFlag(config.ArgUserForceSecAuth, "", false, "Indicates if secure (two-factor) authentication should be forced for the User")

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, userCmd, core.CommandBuilder{
		Namespace: "user",
		Resource:  "user",
		Verb:      "update",
		ShortDesc: "Update a User",
		LongDesc: `Use this command to update details about a specific User including their privileges.

Note: The password attribute is immutable. It is not allowed in update requests. It is recommended that the new User log into the DCD and change their password.

Required values to run command:

* User Id`,
		Example:    updateUserExample,
		PreCmdRun:  noPreRun,
		CmdRun:     RunUserUpdate,
		InitClient: true,
	})
	update.AddStringFlag(config.ArgUserFirstName, "", "", "The firstname for the User")
	update.AddStringFlag(config.ArgUserLastName, "", "", "The lastname for the User")
	update.AddStringFlag(config.ArgUserEmail, "", "", "The email for the User")
	update.AddBoolFlag(config.ArgUserAdministrator, "", false, "Assigns the User to have administrative rights")
	update.AddBoolFlag(config.ArgUserForceSecAuth, "", false, "Indicates if secure (two-factor) authentication should be forced for the User")
	update.AddStringFlag(config.ArgUserId, "", "", config.RequiredFlagUserId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getUsersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, userCmd, core.CommandBuilder{
		Namespace: "user",
		Resource:  "user",
		Verb:      "delete",
		ShortDesc: "Blacklists the User, disabling them",
		LongDesc: `This command blacklists the User, disabling them. The User is not completely purged, therefore if you anticipate needing to create a User with the same name in the future, we suggest renaming the User before you delete it.

Required values to run command:

* User Id`,
		Example:    deleteUserExample,
		PreCmdRun:  PreRunUserId,
		CmdRun:     RunUserDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(config.ArgUserId, "", "", config.RequiredFlagUserId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getUsersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	userCmd.AddCommand(userS3key())

	return userCmd
}

func PreRunUserId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.NS, config.ArgUserId)
}

func PreRunUserNameEmailPwd(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.NS, config.ArgUserFirstName, config.ArgUserLastName, config.ArgUserEmail, config.ArgUserPassword)
}

func RunUserList(c *core.CommandConfig) error {
	users, _, err := c.Users().List()
	if err != nil {
		return err
	}
	return c.Printer.Print(getUserPrint(nil, c, getUsers(users)))
}

func RunUserGet(c *core.CommandConfig) error {
	u, _, err := c.Users().Get(viper.GetString(core.GetFlagName(c.NS, config.ArgUserId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getUserPrint(nil, c, getUser(u)))
}

func RunUserCreate(c *core.CommandConfig) error {
	firstname := viper.GetString(core.GetFlagName(c.NS, config.ArgUserFirstName))
	lastname := viper.GetString(core.GetFlagName(c.NS, config.ArgUserLastName))
	email := viper.GetString(core.GetFlagName(c.NS, config.ArgUserEmail))
	pwd := viper.GetString(core.GetFlagName(c.NS, config.ArgUserPassword))
	secureAuth := viper.GetBool(core.GetFlagName(c.NS, config.ArgUserForceSecAuth))
	admin := viper.GetBool(core.GetFlagName(c.NS, config.ArgUserAdministrator))
	newUser := resources.User{
		User: ionoscloud.User{
			Properties: &ionoscloud.UserProperties{
				Firstname:     &firstname,
				Lastname:      &lastname,
				Email:         &email,
				Administrator: &admin,
				ForceSecAuth:  &secureAuth,
				Password:      &pwd,
			},
		},
	}
	u, resp, err := c.Users().Create(newUser)
	if err != nil {
		return err
	}
	return c.Printer.Print(getUserPrint(resp, c, getUser(u)))
}

func RunUserUpdate(c *core.CommandConfig) error {
	oldUser, resp, err := c.Users().Get(viper.GetString(core.GetFlagName(c.NS, config.ArgUserId)))
	if err != nil {
		return err
	}
	newProperties := getUserInfo(oldUser, c)
	newUser := resources.User{
		User: ionoscloud.User{
			Properties: &newProperties.UserProperties,
		},
	}
	userUpd, resp, err := c.Users().Update(viper.GetString(core.GetFlagName(c.NS, config.ArgUserId)), newUser)
	if err != nil {
		return err
	}
	return c.Printer.Print(getUserPrint(resp, c, getUser(userUpd)))
}

func RunUserDelete(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete user"); err != nil {
		return err
	}
	resp, err := c.Users().Delete(viper.GetString(core.GetFlagName(c.NS, config.ArgUserId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getUserPrint(resp, c, nil))
}

func getUserInfo(oldUser *resources.User, c *core.CommandConfig) *resources.UserProperties {
	var (
		firstName, lastName, email string
		forceSecureAuth, admin     bool
	)
	if properties, ok := oldUser.GetPropertiesOk(); ok && properties != nil {
		if viper.IsSet(core.GetFlagName(c.NS, config.ArgUserFirstName)) {
			firstName = viper.GetString(core.GetFlagName(c.NS, config.ArgUserFirstName))
		} else {
			if name, ok := properties.GetFirstnameOk(); ok && name != nil {
				firstName = *name
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, config.ArgUserLastName)) {
			lastName = viper.GetString(core.GetFlagName(c.NS, config.ArgUserLastName))
		} else {
			if name, ok := properties.GetLastnameOk(); ok && name != nil {
				lastName = *name
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, config.ArgUserEmail)) {
			email = viper.GetString(core.GetFlagName(c.NS, config.ArgUserEmail))
		} else {
			if e, ok := properties.GetEmailOk(); ok && e != nil {
				email = *e
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, config.ArgUserForceSecAuth)) {
			forceSecureAuth = viper.GetBool(core.GetFlagName(c.NS, config.ArgUserForceSecAuth))
		} else {
			if secAuth, ok := properties.GetForceSecAuthOk(); ok && secAuth != nil {
				forceSecureAuth = *secAuth
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, config.ArgUserAdministrator)) {
			admin = viper.GetBool(core.GetFlagName(c.NS, config.ArgUserAdministrator))
		} else {
			if administrator, ok := properties.GetAdministratorOk(); ok && administrator != nil {
				admin = *administrator
			}
		}
	}
	return &resources.UserProperties{
		UserProperties: ionoscloud.UserProperties{
			Firstname:     &firstName,
			Lastname:      &lastName,
			Email:         &email,
			Administrator: &admin,
			ForceSecAuth:  &forceSecureAuth,
		},
	}
}

func groupUser() *core.Command {
	ctx := context.TODO()
	groupUserCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "user",
			Short:            "Group User Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl group user` + "`" + ` allow you to list, add, remove Users from a Group.`,
			TraverseChildren: true,
		},
	}
	globalFlags := groupUserCmd.GlobalFlags()
	globalFlags.StringSlice(config.ArgCols, defaultGroupCols, "Columns to be printed in the standard output")
	_ = viper.BindPFlag(core.GetGlobalFlagName(groupUserCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = groupUserCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultGroupCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Users Command
	*/
	listUsers := core.NewCommand(ctx, groupUserCmd, core.CommandBuilder{
		Namespace:  "group",
		Resource:   "user",
		Verb:       "list",
		ShortDesc:  "List Users from a Group",
		LongDesc:   "Use this command to get a list of Users from a specific Group.\n\nRequired values to run command:\n\n* Group Id",
		Example:    listGroupUsersExample,
		PreCmdRun:  PreRunGroupId,
		CmdRun:     RunGroupUserList,
		InitClient: true,
	})
	listUsers.AddStringFlag(config.ArgGroupId, "", "", config.RequiredFlagGroupId)
	_ = listUsers.Command.RegisterFlagCompletionFunc(config.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getGroupsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Add User Command
	*/
	addUser := core.NewCommand(ctx, groupUserCmd, core.CommandBuilder{
		Namespace:  "group",
		Resource:   "user",
		Verb:       "add",
		ShortDesc:  "Add User to a Group",
		LongDesc:   "Use this command to add an existing User to a specific Group.\n\nRequired values to run command:\n\n* Group Id\n* User Id",
		Example:    addGroupUserExample,
		PreCmdRun:  PreRunGroupUserIds,
		CmdRun:     RunGroupUserAdd,
		InitClient: true,
	})
	addUser.AddStringFlag(config.ArgGroupId, "", "", config.RequiredFlagGroupId)
	_ = addUser.Command.RegisterFlagCompletionFunc(config.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getGroupsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	addUser.AddStringFlag(config.ArgUserId, "", "", config.RequiredFlagUserId)
	_ = addUser.Command.RegisterFlagCompletionFunc(config.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getUsersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Remove User Command
	*/
	removeUser := core.NewCommand(ctx, groupUserCmd, core.CommandBuilder{
		Namespace:  "group",
		Resource:   "user",
		Verb:       "remove",
		ShortDesc:  "Remove User from a Group",
		LongDesc:   "Use this command to remove a User from a Group.\n\nRequired values to run command:\n\n* Group Id\n* User Id",
		Example:    removeGroupUserExample,
		PreCmdRun:  PreRunGroupUserIds,
		CmdRun:     RunGroupUserRemove,
		InitClient: true,
	})
	removeUser.AddStringFlag(config.ArgGroupId, "", "", config.RequiredFlagGroupId)
	_ = removeUser.Command.RegisterFlagCompletionFunc(config.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getGroupsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	removeUser.AddStringFlag(config.ArgUserId, "", "", config.RequiredFlagUserId)
	_ = removeUser.Command.RegisterFlagCompletionFunc(config.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getGroupUsersIds(os.Stderr, viper.GetString(core.GetFlagName(removeUser.NS, config.ArgGroupId))), cobra.ShellCompDirectiveNoFileComp
	})

	return groupUserCmd
}

func RunGroupUserList(c *core.CommandConfig) error {
	users, _, err := c.Groups().ListUsers(viper.GetString(core.GetFlagName(c.NS, config.ArgGroupId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getUserPrint(nil, c, getGroupUsers(users)))
}

func RunGroupUserAdd(c *core.CommandConfig) error {
	id := viper.GetString(core.GetFlagName(c.NS, config.ArgUserId))
	u := resources.User{
		User: ionoscloud.User{
			Id: &id,
		},
	}
	userAdded, resp, err := c.Groups().AddUser(viper.GetString(core.GetFlagName(c.NS, config.ArgGroupId)), u)
	if err != nil {
		return err
	}
	return c.Printer.Print(getUserPrint(resp, c, getUser(userAdded)))
}

func RunGroupUserRemove(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "remove user from group"); err != nil {
		return err
	}
	resp, err := c.Groups().RemoveUser(
		viper.GetString(core.GetFlagName(c.NS, config.ArgGroupId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgUserId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getGroupPrint(resp, c, nil))
}

// Output Printing

var defaultUserCols = []string{"UserId", "Firstname", "Lastname", "Email", "S3CanonicalUserId", "Administrator", "ForceSecAuth", "SecAuthActive", "Active"}

type UserPrint struct {
	UserId            string `json:"UserId,omitempty"`
	Firstname         string `json:"Firstname,omitempty"`
	Lastname          string `json:"Lastname,omitempty"`
	Email             string `json:"Email,omitempty"`
	Administrator     bool   `json:"Administrator,omitempty"`
	ForceSecAuth      bool   `json:"ForceSecAuth,omitempty"`
	SecAuthActive     bool   `json:"SecAuthActive,omitempty"`
	S3CanonicalUserId string `json:"S3CanonicalUserId,omitempty"`
	Active            bool   `json:"Active,omitempty"`
}

func getUserPrint(resp *resources.Response, c *core.CommandConfig, users []resources.User) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForRequest = viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForRequest))
		}
		if users != nil {
			r.OutputJSON = users
			r.KeyValue = getUsersKVMaps(users)
			r.Columns = getUserCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getUserCols(flagName string, outErr io.Writer) []string {
	if viper.IsSet(flagName) {
		var userCols []string
		columnsMap := map[string]string{
			"UserId":            "UserId",
			"Firstname":         "Firstname",
			"Lastname":          "Lastname",
			"Email":             "Email",
			"Administrator":     "Administrator",
			"ForceSecAuth":      "ForceSecAuth",
			"SecAuthActive":     "SecAuthActive",
			"S3CanonicalUserId": "S3CanonicalUserId",
			"Active":            "Active",
		}
		for _, k := range viper.GetStringSlice(flagName) {
			col := columnsMap[k]
			if col != "" {
				userCols = append(userCols, col)
			} else {
				clierror.CheckError(errors.New("unknown column "+k), outErr)
			}
		}
		return userCols
	} else {
		return defaultUserCols
	}
}

func getUsers(users resources.Users) []resources.User {
	u := make([]resources.User, 0)
	if items, ok := users.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			u = append(u, resources.User{User: item})
		}
	}
	return u
}

func getGroupUsers(users resources.GroupMembers) []resources.User {
	u := make([]resources.User, 0)
	if items, ok := users.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			u = append(u, resources.User{User: item})
		}
	}
	return u
}

func getUser(u *resources.User) []resources.User {
	users := make([]resources.User, 0)
	if u != nil {
		users = append(users, resources.User{User: u.User})
	}
	return users
}

func getUsersKVMaps(us []resources.User) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(us))
	for _, u := range us {
		var uPrint UserPrint
		if id, ok := u.GetIdOk(); ok && id != nil {
			uPrint.UserId = *id
		}
		if properties, ok := u.GetPropertiesOk(); ok && properties != nil {
			if firstname, ok := properties.GetFirstnameOk(); ok && firstname != nil {
				uPrint.Firstname = *firstname
			}
			if lastname, ok := properties.GetLastnameOk(); ok && lastname != nil {
				uPrint.Lastname = *lastname
			}
			if email, ok := properties.GetEmailOk(); ok && email != nil {
				uPrint.Email = *email
			}
			if administrator, ok := properties.GetAdministratorOk(); ok && administrator != nil {
				uPrint.Administrator = *administrator
			}
			if forceSecAuth, ok := properties.GetForceSecAuthOk(); ok && forceSecAuth != nil {
				uPrint.ForceSecAuth = *forceSecAuth
			}
			if authActive, ok := properties.GetSecAuthActiveOk(); ok && authActive != nil {
				uPrint.SecAuthActive = *authActive
			}
			if canonicalUserId, ok := properties.GetS3CanonicalUserIdOk(); ok && canonicalUserId != nil {
				uPrint.S3CanonicalUserId = *canonicalUserId
			}
			if active, ok := properties.GetActiveOk(); ok && active != nil {
				uPrint.Active = *active
			}
		}
		o := structs.Map(uPrint)
		out = append(out, o)
	}
	return out
}

func getUsersIds(outErr io.Writer) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)
	userSvc := resources.NewUserService(clientSvc.Get(), context.TODO())
	users, _, err := userSvc.List()
	clierror.CheckError(err, outErr)
	usersIds := make([]string, 0)
	if items, ok := users.Users.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				usersIds = append(usersIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return usersIds
}

func getGroupUsersIds(outErr io.Writer, groupId string) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)
	groupSvc := resources.NewGroupService(clientSvc.Get(), context.TODO())
	users, _, err := groupSvc.ListUsers(groupId)
	clierror.CheckError(err, outErr)
	usersIds := make([]string, 0)
	if items, ok := users.GroupMembers.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				usersIds = append(usersIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return usersIds
}
