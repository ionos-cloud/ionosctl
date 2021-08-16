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

func user() *core.Command {
	ctx := context.TODO()
	userCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "user",
			Aliases:          []string{"u"},
			Short:            "User Operations",
			Long:             "The sub-commands of `ionosctl user` allow you to list, get, create, update, delete Users under your account. To add Users to a Group, check the `ionosctl group user` commands. To add S3Keys to a User, check the `ionosctl user s3key` commands.",
			TraverseChildren: true,
		},
	}
	globalFlags := userCmd.GlobalFlags()
	globalFlags.StringSliceP(config.ArgCols, "", defaultUserCols, utils.ColsMessage(defaultUserCols))
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
		Aliases:    []string{"l", "ls"},
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
		Aliases:    []string{"g"},
		ShortDesc:  "Get a User",
		LongDesc:   "Use this command to retrieve details about a specific User.\n\nRequired values to run command:\n\n* User Id",
		Example:    getUserExample,
		PreCmdRun:  PreRunUserId,
		CmdRun:     RunUserGet,
		InitClient: true,
	})
	get.AddStringFlag(config.ArgUserId, config.ArgIdShort, "", config.RequiredFlagUserId)
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
		Aliases:   []string{"c"},
		ShortDesc: "Create a User under a particular contract",
		LongDesc: `Use this command to create a User under a particular contract. You need to specify the firstname, lastname, email and password for the new User.

Note: The password set here cannot be updated through the API currently. It is recommended that a new User log into the DCD and change their password.

Required values to run a command:

* First Name
* Last Name
* Email
* Password`,
		Example:    createUserExample,
		PreCmdRun:  PreRunUserNameEmailPwd,
		CmdRun:     RunUserCreate,
		InitClient: true,
	})
	create.AddStringFlag(config.ArgFirstName, "", "", "The first name for the User "+config.RequiredFlag)
	create.AddStringFlag(config.ArgLastName, "", "", "The last name for the User "+config.RequiredFlag)
	create.AddStringFlag(config.ArgEmail, config.ArgEmailShort, "", "The email for the User "+config.RequiredFlag)
	create.AddStringFlag(config.ArgPassword, config.ArgPasswordShort, "", "The password for the User (must be at least 5 characters long) "+config.RequiredFlag)
	create.AddBoolFlag(config.ArgAdmin, "", false, "Assigns the User to have administrative rights")
	create.AddBoolFlag(config.ArgForceSecAuth, "", false, "Indicates if secure (two-factor) authentication should be forced for the User")

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, userCmd, core.CommandBuilder{
		Namespace: "user",
		Resource:  "user",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
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
	update.AddStringFlag(config.ArgFirstName, "", "", "The first name for the User")
	update.AddStringFlag(config.ArgLastName, "", "", "The last name for the User")
	update.AddStringFlag(config.ArgEmail, config.ArgEmailShort, "", "The email for the User")
	update.AddBoolFlag(config.ArgAdmin, "", false, "Assigns the User to have administrative rights")
	update.AddBoolFlag(config.ArgForceSecAuth, "", false, "Indicates if secure (two-factor) authentication should be forced for the User")
	update.AddStringFlag(config.ArgUserId, config.ArgIdShort, "", config.RequiredFlagUserId)
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
		Aliases:   []string{"d"},
		ShortDesc: "Blacklists the User, disabling them",
		LongDesc: `This command blacklists the User, disabling them. The User is not completely purged, therefore if you anticipate needing to create a User with the same name in the future, we suggest renaming the User before you delete it.

Required values to run command:

* User Id`,
		Example:    deleteUserExample,
		PreCmdRun:  PreRunUserId,
		CmdRun:     RunUserDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(config.ArgUserId, config.ArgIdShort, "", config.RequiredFlagUserId)
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
	return core.CheckRequiredFlags(c.NS, config.ArgFirstName, config.ArgLastName, config.ArgEmail, config.ArgPassword)
}

func RunUserList(c *core.CommandConfig) error {
	users, resp, err := c.Users().List()
	if resp != nil {
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getUserPrint(nil, c, getUsers(users)))
}

func RunUserGet(c *core.CommandConfig) error {
	c.Printer.Verbose("User with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, config.ArgUserId)))
	u, resp, err := c.Users().Get(viper.GetString(core.GetFlagName(c.NS, config.ArgUserId)))
	if resp != nil {
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getUserPrint(nil, c, getUser(u)))
}

func RunUserCreate(c *core.CommandConfig) error {
	firstname := viper.GetString(core.GetFlagName(c.NS, config.ArgFirstName))
	lastname := viper.GetString(core.GetFlagName(c.NS, config.ArgLastName))
	email := viper.GetString(core.GetFlagName(c.NS, config.ArgEmail))
	pwd := viper.GetString(core.GetFlagName(c.NS, config.ArgPassword))
	secureAuth := viper.GetBool(core.GetFlagName(c.NS, config.ArgForceSecAuth))
	admin := viper.GetBool(core.GetFlagName(c.NS, config.ArgAdmin))
	newUser := v5.UserPost{
		UserPost: ionoscloud.UserPost{
			Properties: &ionoscloud.UserPropertiesPost{
				Firstname:     &firstname,
				Lastname:      &lastname,
				Email:         &email,
				Administrator: &admin,
				ForceSecAuth:  &secureAuth,
				Password:      &pwd,
			},
		},
	}
	c.Printer.Verbose("Properties set for creating the user: Firstname: %v, Lastname: %v, Email: %v, ForceSecAuth: %v, Administrator: %v",
		firstname, lastname, email, secureAuth, admin)
	u, resp, err := c.Users().Create(newUser)
	if resp != nil {
		c.Printer.Verbose("Request href: %v ", resp.Header.Get("location"))
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getUserPrint(resp, c, getUser(u)))
}

func RunUserUpdate(c *core.CommandConfig) error {
	oldUser, resp, err := c.Users().Get(viper.GetString(core.GetFlagName(c.NS, config.ArgUserId)))
	if resp != nil {
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
	}
	if err != nil {
		return err
	}
	newUser := getUserInfo(oldUser, c)
	userUpd, resp, err := c.Users().Update(viper.GetString(core.GetFlagName(c.NS, config.ArgUserId)), *newUser)
	if err != nil {
		return err
	}
	return c.Printer.Print(getUserPrint(resp, c, getUser(userUpd)))
}

func RunUserDelete(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete user"); err != nil {
		return err
	}
	c.Printer.Verbose("User with id: %v is deleting...", viper.GetString(core.GetFlagName(c.NS, config.ArgUserId)))
	resp, err := c.Users().Delete(viper.GetString(core.GetFlagName(c.NS, config.ArgUserId)))
	if resp != nil {
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getUserPrint(resp, c, nil))
}

func getUserInfo(oldUser *v5.User, c *core.CommandConfig) *v5.UserPut {
	var (
		firstName, lastName, email string
		forceSecureAuth, admin     bool
	)
	if properties, ok := oldUser.GetPropertiesOk(); ok && properties != nil {
		if viper.IsSet(core.GetFlagName(c.NS, config.ArgFirstName)) {
			firstName = viper.GetString(core.GetFlagName(c.NS, config.ArgFirstName))
			c.Printer.Verbose("Property FirstName set: %v", firstName)
		} else {
			if name, ok := properties.GetFirstnameOk(); ok && name != nil {
				firstName = *name
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, config.ArgLastName)) {
			lastName = viper.GetString(core.GetFlagName(c.NS, config.ArgLastName))
			c.Printer.Verbose("Property LastName set: %v", lastName)
		} else {
			if name, ok := properties.GetLastnameOk(); ok && name != nil {
				lastName = *name
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, config.ArgEmail)) {
			email = viper.GetString(core.GetFlagName(c.NS, config.ArgEmail))
			c.Printer.Verbose("Property Email set: %v", email)
		} else {
			if e, ok := properties.GetEmailOk(); ok && e != nil {
				email = *e
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, config.ArgForceSecAuth)) {
			forceSecureAuth = viper.GetBool(core.GetFlagName(c.NS, config.ArgForceSecAuth))
			c.Printer.Verbose("Property ForceSecAuth set: %v", forceSecureAuth)
		} else {
			if secAuth, ok := properties.GetForceSecAuthOk(); ok && secAuth != nil {
				forceSecureAuth = *secAuth
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, config.ArgAdmin)) {
			admin = viper.GetBool(core.GetFlagName(c.NS, config.ArgAdmin))
			c.Printer.Verbose("Property Administrator set: %v", admin)
		} else {
			if administrator, ok := properties.GetAdministratorOk(); ok && administrator != nil {
				admin = *administrator
			}
		}
	}
	return &v5.UserPut{
		UserPut: ionoscloud.UserPut{
			Properties: &ionoscloud.UserPropertiesPut{
				Firstname:     &firstName,
				Lastname:      &lastName,
				Email:         &email,
				Administrator: &admin,
				ForceSecAuth:  &forceSecureAuth,
			},
		},
	}
}

func groupUser() *core.Command {
	ctx := context.TODO()
	groupUserCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "user",
			Aliases:          []string{"u"},
			Short:            "Group User Operations",
			Long:             "The sub-commands of `ionosctl group user` allow you to list, add, remove Users from a Group.",
			TraverseChildren: true,
		},
	}

	/*
		List Users Command
	*/
	listUsers := core.NewCommand(ctx, groupUserCmd, core.CommandBuilder{
		Namespace:  "group",
		Resource:   "user",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Users from a Group",
		LongDesc:   "Use this command to get a list of Users from a specific Group.\n\nRequired values to run command:\n\n* Group Id",
		Example:    listGroupUsersExample,
		PreCmdRun:  PreRunGroupId,
		CmdRun:     RunGroupUserList,
		InitClient: true,
	})
	listUsers.AddStringSliceFlag(config.ArgCols, "", defaultUserCols, utils.ColsMessage(defaultUserCols))
	_ = listUsers.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultUserCols, cobra.ShellCompDirectiveNoFileComp
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
		Aliases:    []string{"a"},
		ShortDesc:  "Add User to a Group",
		LongDesc:   "Use this command to add an existing User to a specific Group.\n\nRequired values to run command:\n\n* Group Id\n* User Id",
		Example:    addGroupUserExample,
		PreCmdRun:  PreRunGroupUserIds,
		CmdRun:     RunGroupUserAdd,
		InitClient: true,
	})
	addUser.AddStringSliceFlag(config.ArgCols, "", defaultUserCols, utils.ColsMessage(defaultUserCols))
	_ = addUser.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultUserCols, cobra.ShellCompDirectiveNoFileComp
	})
	addUser.AddStringFlag(config.ArgGroupId, "", "", config.RequiredFlagGroupId)
	_ = addUser.Command.RegisterFlagCompletionFunc(config.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getGroupsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	addUser.AddStringFlag(config.ArgUserId, config.ArgIdShort, "", config.RequiredFlagUserId)
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
		Aliases:    []string{"r"},
		ShortDesc:  "Remove User from a Group",
		LongDesc:   "Use this command to remove a User from a Group.\n\nRequired values to run command:\n\n* Group Id\n* User Id",
		Example:    removeGroupUserExample,
		PreCmdRun:  PreRunGroupUserIds,
		CmdRun:     RunGroupUserRemove,
		InitClient: true,
	})
	removeUser.AddStringSliceFlag(config.ArgCols, "", defaultUserCols, utils.ColsMessage(defaultUserCols))
	_ = removeUser.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultUserCols, cobra.ShellCompDirectiveNoFileComp
	})
	removeUser.AddStringFlag(config.ArgGroupId, "", "", config.RequiredFlagGroupId)
	_ = removeUser.Command.RegisterFlagCompletionFunc(config.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getGroupsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	removeUser.AddStringFlag(config.ArgUserId, config.ArgIdShort, "", config.RequiredFlagUserId)
	_ = removeUser.Command.RegisterFlagCompletionFunc(config.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getGroupUsersIds(os.Stderr, viper.GetString(core.GetFlagName(removeUser.NS, config.ArgGroupId))), cobra.ShellCompDirectiveNoFileComp
	})

	return groupUserCmd
}

func RunGroupUserList(c *core.CommandConfig) error {
	users, resp, err := c.Groups().ListUsers(viper.GetString(core.GetFlagName(c.NS, config.ArgGroupId)))
	if resp != nil {
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getUserPrint(nil, c, getGroupUsers(users)))
}

func RunGroupUserAdd(c *core.CommandConfig) error {
	id := viper.GetString(core.GetFlagName(c.NS, config.ArgUserId))
	u := v5.User{
		User: ionoscloud.User{
			Id: &id,
		},
	}
	groupId := viper.GetString(core.GetFlagName(c.NS, config.ArgGroupId))
	c.Printer.Verbose("Adding User with id: %v to group with id: %v...", id, groupId)
	userAdded, resp, err := c.Groups().AddUser(groupId, u)
	if resp != nil {
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getUserPrint(resp, c, getUser(userAdded)))
}

func RunGroupUserRemove(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "remove user from group"); err != nil {
		return err
	}
	id := viper.GetString(core.GetFlagName(c.NS, config.ArgUserId))
	groupId := viper.GetString(core.GetFlagName(c.NS, config.ArgGroupId))
	c.Printer.Verbose("Adding User with id: %v to group with id: %v...", id, groupId)
	resp, err := c.Groups().RemoveUser(groupId, id)
	if resp != nil {
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
	}
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

func getUserPrint(resp *v5.Response, c *core.CommandConfig, users []v5.User) printer.Result {
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
			if c.Resource != c.Namespace {
				r.Columns = getUserCols(core.GetFlagName(c.NS, config.ArgCols), c.Printer.GetStderr())
			} else {
				r.Columns = getUserCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
			}
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

func getUsers(users v5.Users) []v5.User {
	u := make([]v5.User, 0)
	if items, ok := users.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			u = append(u, v5.User{User: item})
		}
	}
	return u
}

func getGroupUsers(users v5.GroupMembers) []v5.User {
	u := make([]v5.User, 0)
	if items, ok := users.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			u = append(u, v5.User{User: item})
		}
	}
	return u
}

func getUser(u *v5.User) []v5.User {
	users := make([]v5.User, 0)
	if u != nil {
		users = append(users, v5.User{User: u.User})
	}
	return users
}

func getUsersKVMaps(us []v5.User) []map[string]interface{} {
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
	clientSvc, err := v5.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		config.GetServerUrl(),
	)
	clierror.CheckError(err, outErr)
	userSvc := v5.NewUserService(clientSvc.Get(), context.TODO())
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
	clientSvc, err := v5.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		config.GetServerUrl(),
	)
	clierror.CheckError(err, outErr)
	groupSvc := v5.NewGroupService(clientSvc.Get(), context.TODO())
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
