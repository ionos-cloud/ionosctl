package commands

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/completer"
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

func UserCmd() *core.Command {
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
	globalFlags.StringSliceP(config.ArgCols, "", defaultUserCols, printer.ColsMessage(defaultUserCols))
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
		PreCmdRun:  core.NoPreRun,
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
	get.AddStringFlag(cloudapiv6.ArgUserId, cloudapiv6.ArgIdShort, "", cloudapiv6.UserId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.UsersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
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
	create.AddStringFlag(cloudapiv6.ArgFirstName, "", "", "The first name for the User", core.RequiredFlagOption())
	create.AddStringFlag(cloudapiv6.ArgLastName, "", "", "The last name for the User", core.RequiredFlagOption())
	create.AddStringFlag(cloudapiv6.ArgEmail, cloudapiv6.ArgEmailShort, "", "The email for the User", core.RequiredFlagOption())
	create.AddStringFlag(cloudapiv6.ArgPassword, cloudapiv6.ArgPasswordShort, "", "The password for the User (must be at least 5 characters long)", core.RequiredFlagOption())
	create.AddBoolFlag(cloudapiv6.ArgAdmin, "", false, "Assigns the User to have administrative rights")
	create.AddBoolFlag(cloudapiv6.ArgForceSecAuth, "", false, "Indicates if secure (two-factor) authentication should be forced for the User")

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
		PreCmdRun:  PreRunUserId,
		CmdRun:     RunUserUpdate,
		InitClient: true,
	})
	update.AddStringFlag(cloudapiv6.ArgFirstName, "", "", "The first name for the User")
	update.AddStringFlag(cloudapiv6.ArgLastName, "", "", "The last name for the User")
	update.AddStringFlag(cloudapiv6.ArgEmail, cloudapiv6.ArgEmailShort, "", "The email for the User")
	update.AddBoolFlag(cloudapiv6.ArgAdmin, "", false, "Assigns the User to have administrative rights")
	update.AddBoolFlag(cloudapiv6.ArgForceSecAuth, "", false, "Indicates if secure (two-factor) authentication should be forced for the User")
	update.AddStringFlag(cloudapiv6.ArgUserId, cloudapiv6.ArgIdShort, "", cloudapiv6.UserId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.UsersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
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
	deleteCmd.AddStringFlag(cloudapiv6.ArgUserId, cloudapiv6.ArgIdShort, "", cloudapiv6.UserId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.UsersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	userCmd.AddCommand(UserS3keyCmd())

	return userCmd
}

func PreRunUserId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgUserId)
}

func PreRunUserNameEmailPwd(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgFirstName, cloudapiv6.ArgLastName, cloudapiv6.ArgEmail, cloudapiv6.ArgPassword)
}

func RunUserList(c *core.CommandConfig) error {
	users, resp, err := c.CloudApiV6Services.Users().List()
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getUserPrint(nil, c, getUsers(users)))
}

func RunUserGet(c *core.CommandConfig) error {
	c.Printer.Verbose("User with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId)))
	u, resp, err := c.CloudApiV6Services.Users().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId)))
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getUserPrint(nil, c, getUser(u)))
}

func RunUserCreate(c *core.CommandConfig) error {
	firstname := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFirstName))
	lastname := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLastName))
	email := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgEmail))
	pwd := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPassword))
	secureAuth := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgForceSecAuth))
	admin := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAdmin))
	newUser := resources.UserPost{
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
	c.Printer.Verbose("Creating User...")
	u, resp, err := c.CloudApiV6Services.Users().Create(newUser)
	if resp != nil {
		c.Printer.Verbose("Request href: %v ", resp.Header.Get("location"))
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getUserPrint(resp, c, getUser(u)))
}

func RunUserUpdate(c *core.CommandConfig) error {
	oldUser, resp, err := c.CloudApiV6Services.Users().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId)))
	if err != nil {
		return err
	}
	newUser := getUserInfo(oldUser, c)
	c.Printer.Verbose("Updating User with ID: %v...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId)))
	userUpd, resp, err := c.CloudApiV6Services.Users().Update(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId)), *newUser)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getUserPrint(resp, c, getUser(userUpd)))
}

func RunUserDelete(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete user"); err != nil {
		return err
	}
	c.Printer.Verbose("User with id: %v is deleting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId)))
	resp, err := c.CloudApiV6Services.Users().Delete(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId)))
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getUserPrint(resp, c, nil))
}

func getUserInfo(oldUser *resources.User, c *core.CommandConfig) *resources.UserPut {
	var (
		firstName, lastName, email string
		forceSecureAuth, admin     bool
	)
	if properties, ok := oldUser.GetPropertiesOk(); ok && properties != nil {
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgFirstName)) {
			firstName = viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFirstName))
			c.Printer.Verbose("Property FirstName set: %v", firstName)
		} else {
			if name, ok := properties.GetFirstnameOk(); ok && name != nil {
				firstName = *name
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgLastName)) {
			lastName = viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLastName))
			c.Printer.Verbose("Property LastName set: %v", lastName)
		} else {
			if name, ok := properties.GetLastnameOk(); ok && name != nil {
				lastName = *name
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgEmail)) {
			email = viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgEmail))
			c.Printer.Verbose("Property Email set: %v", email)
		} else {
			if e, ok := properties.GetEmailOk(); ok && e != nil {
				email = *e
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgForceSecAuth)) {
			forceSecureAuth = viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgForceSecAuth))
			c.Printer.Verbose("Property ForceSecAuth set: %v", forceSecureAuth)
		} else {
			if secAuth, ok := properties.GetForceSecAuthOk(); ok && secAuth != nil {
				forceSecureAuth = *secAuth
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgAdmin)) {
			admin = viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAdmin))
			c.Printer.Verbose("Property Administrator set: %v", admin)
		} else {
			if administrator, ok := properties.GetAdministratorOk(); ok && administrator != nil {
				admin = *administrator
			}
		}
	}
	return &resources.UserPut{
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

func GroupUserCmd() *core.Command {
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
	listUsers.AddStringSliceFlag(config.ArgCols, "", defaultUserCols, printer.ColsMessage(defaultUserCols))
	_ = listUsers.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultUserCols, cobra.ShellCompDirectiveNoFileComp
	})
	listUsers.AddStringFlag(cloudapiv6.ArgGroupId, "", "", cloudapiv6.GroupId, core.RequiredFlagOption())
	_ = listUsers.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GroupsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
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
	addUser.AddStringSliceFlag(config.ArgCols, "", defaultUserCols, printer.ColsMessage(defaultUserCols))
	_ = addUser.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultUserCols, cobra.ShellCompDirectiveNoFileComp
	})
	addUser.AddStringFlag(cloudapiv6.ArgGroupId, "", "", cloudapiv6.GroupId, core.RequiredFlagOption())
	_ = addUser.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GroupsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	addUser.AddStringFlag(cloudapiv6.ArgUserId, cloudapiv6.ArgIdShort, "", cloudapiv6.UserId, core.RequiredFlagOption())
	_ = addUser.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.UsersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
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
	removeUser.AddStringSliceFlag(config.ArgCols, "", defaultUserCols, printer.ColsMessage(defaultUserCols))
	_ = removeUser.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultUserCols, cobra.ShellCompDirectiveNoFileComp
	})
	removeUser.AddStringFlag(cloudapiv6.ArgGroupId, "", "", cloudapiv6.GroupId, core.RequiredFlagOption())
	_ = removeUser.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GroupsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	removeUser.AddStringFlag(cloudapiv6.ArgUserId, cloudapiv6.ArgIdShort, "", cloudapiv6.UserId, core.RequiredFlagOption())
	_ = removeUser.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GroupUsersIds(os.Stderr, viper.GetString(core.GetFlagName(removeUser.NS, cloudapiv6.ArgGroupId))), cobra.ShellCompDirectiveNoFileComp
	})

	return groupUserCmd
}

func RunGroupUserList(c *core.CommandConfig) error {
	users, resp, err := c.CloudApiV6Services.Groups().ListUsers(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId)))
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getUserPrint(nil, c, getGroupUsers(users)))
}

func RunGroupUserAdd(c *core.CommandConfig) error {
	id := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId))
	groupId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId))
	c.Printer.Verbose("User with id: %v is adding to group with id: %v...", id, groupId)
	u := resources.User{
		User: ionoscloud.User{
			Id: &id,
		},
	}
	userAdded, resp, err := c.CloudApiV6Services.Groups().AddUser(groupId, u)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
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
	userId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId))
	groupId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId))
	c.Printer.Verbose("User with id: %v is adding to group with id: %v...", userId, groupId)
	resp, err := c.CloudApiV6Services.Groups().RemoveUser(
		groupId,
		userId,
	)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
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
