package commands

import (
	"context"
	"errors"
	"fmt"
	"os"

	"go.uber.org/multierr"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/query"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/waiter"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
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
	globalFlags.StringSliceP(constants.ArgCols, "", defaultUserCols, printer.ColsMessage(defaultUserCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(userCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = userCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultUserCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, userCmd, core.CommandBuilder{
		Namespace:  "user",
		Resource:   "user",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Users",
		LongDesc:   "Use this command to get a list of existing Users available on your account.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.UsersFiltersUsage(),
		Example:    listUserExample,
		PreCmdRun:  PreRunUserList,
		CmdRun:     RunUserList,
		InitClient: true,
	})
	list.AddInt32Flag(cloudapiv6.ArgMaxResults, cloudapiv6.ArgMaxResultsShort, cloudapiv6.DefaultMaxResults, cloudapiv6.ArgMaxResultsDescription)
	list.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultListDepth, cloudapiv6.ArgDepthDescription)
	list.AddStringFlag(cloudapiv6.ArgOrderBy, "", "", cloudapiv6.ArgOrderByDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.UsersFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv6.ArgFilters, cloudapiv6.ArgFiltersShort, []string{""}, cloudapiv6.ArgFiltersDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.UsersFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)

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
	get.AddUUIDFlag(cloudapiv6.ArgUserId, cloudapiv6.ArgIdShort, "", cloudapiv6.UserId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.UsersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)
	get.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.ArgDepthDescription)

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
	create.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultCreateDepth, cloudapiv6.ArgDepthDescription)

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
	update.AddStringFlag(cloudapiv6.ArgPassword, cloudapiv6.ArgPasswordShort, "", "The password for the User (must be at least 5 characters long)")
	update.AddBoolFlag(cloudapiv6.ArgAdmin, "", false, "Assigns the User to have administrative rights. E.g.: --admin=true, --admin=false")
	update.AddBoolFlag(cloudapiv6.ArgForceSecAuth, "", false, "Indicates if secure (two-factor) authentication should be forced for the User. E.g.: --force-secure-auth=true, --force-secure-auth=false")
	update.AddUUIDFlag(cloudapiv6.ArgUserId, cloudapiv6.ArgIdShort, "", cloudapiv6.UserId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.UsersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultUpdateDepth, cloudapiv6.ArgDepthDescription)

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
		PreCmdRun:  PreRunUserDelete,
		CmdRun:     RunUserDelete,
		InitClient: true,
	})
	deleteCmd.AddUUIDFlag(cloudapiv6.ArgUserId, cloudapiv6.ArgIdShort, "", cloudapiv6.UserId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.UsersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete all the Users.")
	deleteCmd.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultDeleteDepth, cloudapiv6.ArgDepthDescription)

	userCmd.AddCommand(UserS3keyCmd())

	return userCmd
}

func PreRunUserList(c *core.PreCommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgFilters)) {
		return query.ValidateFilters(c, completer.UsersFilters(), completer.UsersFiltersUsage())
	}
	return nil
}

func PreRunUserId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgUserId)
}

func PreRunUserDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgUserId},
		[]string{cloudapiv6.ArgAll},
	)
}

func PreRunUserNameEmailPwd(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgFirstName, cloudapiv6.ArgLastName, cloudapiv6.ArgEmail, cloudapiv6.ArgPassword)
}

func RunUserList(c *core.CommandConfig) error {
	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	users, resp, err := c.CloudApiV6Services.Users().List(listQueryParams)
	if resp != nil {
		c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getUserPrint(nil, c, getUsers(users)))
}

func RunUserGet(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	c.Printer.Verbose("User with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId)))
	u, resp, err := c.CloudApiV6Services.Users().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId)), queryParams)
	if resp != nil {
		c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getUserPrint(nil, c, getUser(u)))
}

func RunUserCreate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
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
	u, resp, err := c.CloudApiV6Services.Users().Create(newUser, queryParams)
	if resp != nil && printer.GetId(resp) != "" {
		c.Printer.Verbose(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getUserPrint(resp, c, getUser(u)))
}

func RunUserUpdate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	oldUser, resp, err := c.CloudApiV6Services.Users().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId)), queryParams)
	if err != nil {
		return err
	}
	newUser := getUserInfo(oldUser, c)
	c.Printer.Verbose("Updating User with ID: %v...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId)))
	userUpd, resp, err := c.CloudApiV6Services.Users().Update(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId)), *newUser, queryParams)
	if resp != nil && printer.GetId(resp) != "" {
		c.Printer.Verbose(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getUserPrint(resp, c, getUser(userUpd)))
}

func RunUserDelete(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	userId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId))
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllUsers(c); err != nil {
			return err
		}
		return c.Printer.Print(printer.Result{Resource: c.Resource, Verb: c.Verb})
	} else {
		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete user"); err != nil {
			return err
		}
		c.Printer.Verbose("Starting deleting User with id: %v...", userId)
		resp, err := c.CloudApiV6Services.Users().Delete(userId, queryParams)
		if resp != nil && printer.GetId(resp) != "" {
			c.Printer.Verbose(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime)
		}
		if err != nil {
			return err
		}
		return c.Printer.Print(getUserPrint(resp, c, nil))
	}
}

func getUserInfo(oldUser *resources.User, c *core.CommandConfig) *resources.UserPut {
	userPropertiesPut := ionoscloud.UserPropertiesPut{}
	if properties, ok := oldUser.GetPropertiesOk(); ok && properties != nil {
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgFirstName)) {
			firstName := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFirstName))
			c.Printer.Verbose("Property FirstName set: %v", firstName)
			userPropertiesPut.SetFirstname(firstName)
		} else {
			if firstnameOk, ok := properties.GetFirstnameOk(); ok && firstnameOk != nil {
				userPropertiesPut.SetFirstname(*firstnameOk)
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgLastName)) {
			lastName := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLastName))
			c.Printer.Verbose("Property LastName set: %v", lastName)
			userPropertiesPut.SetLastname(lastName)
		} else {
			if lastnameOk, ok := properties.GetLastnameOk(); ok && lastnameOk != nil {
				userPropertiesPut.SetLastname(*lastnameOk)
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgEmail)) {
			email := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgEmail))
			c.Printer.Verbose("Property Email set: %v", email)
			userPropertiesPut.SetEmail(email)
		} else {
			if emailOk, ok := properties.GetEmailOk(); ok && emailOk != nil {
				userPropertiesPut.SetEmail(*emailOk)
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPassword)) {
			password := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPassword))
			c.Printer.Verbose("Property Password set: %v", password)
			userPropertiesPut.SetPassword(password)
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgForceSecAuth)) {
			forceSecureAuth := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgForceSecAuth))
			c.Printer.Verbose("Property ForceSecAuth set: %v", forceSecureAuth)
			userPropertiesPut.SetForceSecAuth(forceSecureAuth)
		} else {
			if secAuthOk, ok := properties.GetForceSecAuthOk(); ok && secAuthOk != nil {
				userPropertiesPut.SetForceSecAuth(*secAuthOk)
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgAdmin)) {
			admin := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAdmin))
			c.Printer.Verbose("Property Administrator set: %v", admin)
			userPropertiesPut.SetAdministrator(admin)
		} else {
			if administratorOk, ok := properties.GetAdministratorOk(); ok && administratorOk != nil {
				userPropertiesPut.SetAdministrator(*administratorOk)
			}
		}
	}
	return &resources.UserPut{
		UserPut: ionoscloud.UserPut{
			Properties: &userPropertiesPut,
		},
	}
}

func DeleteAllUsers(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	c.Printer.Verbose("Getting Users...")
	users, resp, err := c.CloudApiV6Services.Users().List(cloudapiv6.ParentResourceListQueryParams)
	if err != nil {
		return err
	}
	if usersItems, ok := users.GetItemsOk(); ok && usersItems != nil {
		if len(*usersItems) > 0 {
			_ = c.Printer.Warn("Users to be deleted:")
			for _, user := range *usersItems {
				delIdAndName := ""
				if id, ok := user.GetIdOk(); ok && id != nil {
					delIdAndName += "User Id: " + *id
				}
				if properties, ok := user.GetPropertiesOk(); ok && properties != nil {
					if firstName, ok := properties.GetFirstnameOk(); ok && firstName != nil {
						delIdAndName += " User First Name: " + *firstName
					}
					if lastName, ok := properties.GetLastnameOk(); ok && lastName != nil {
						delIdAndName += " User Last Name: " + *lastName
					}
				}
				_ = c.Printer.Warn(delIdAndName)
			}
			if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete all the Users"); err != nil {
				return err
			}
			c.Printer.Verbose("Deleting all the Users...")
			var multiErr error
			for _, user := range *usersItems {
				if id, ok := user.GetIdOk(); ok && id != nil {
					c.Printer.Verbose("Starting deleting User with id: %v...", *id)
					resp, err = c.CloudApiV6Services.Users().Delete(*id, queryParams)
					if err != nil {
						multiErr = multierr.Append(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
						continue
					} else {
						_ = c.Printer.Warn(fmt.Sprintf(constants.MessageDeletingAll, c.Resource, *id))
					}
					if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
						multiErr = multierr.Append(multiErr, fmt.Errorf(constants.ErrWaitDeleteAll, c.Resource, *id, err))
						continue
					}
				}
			}
			if multiErr != nil {
				return multiErr
			}
			return nil
		} else {
			return errors.New("no Users found")
		}
	} else {
		return errors.New("could not get items of Users")
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
	listUsers.AddStringSliceFlag(constants.ArgCols, "", defaultUserCols, printer.ColsMessage(defaultUserCols))
	_ = listUsers.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultUserCols, cobra.ShellCompDirectiveNoFileComp
	})
	listUsers.AddUUIDFlag(cloudapiv6.ArgGroupId, "", "", cloudapiv6.GroupId, core.RequiredFlagOption())
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
	addUser.AddStringSliceFlag(constants.ArgCols, "", defaultUserCols, printer.ColsMessage(defaultUserCols))
	_ = addUser.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultUserCols, cobra.ShellCompDirectiveNoFileComp
	})
	addUser.AddUUIDFlag(cloudapiv6.ArgGroupId, "", "", cloudapiv6.GroupId, core.RequiredFlagOption())
	_ = addUser.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GroupsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	addUser.AddUUIDFlag(cloudapiv6.ArgUserId, cloudapiv6.ArgIdShort, "", cloudapiv6.UserId, core.RequiredFlagOption())
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
		PreCmdRun:  PreRunGroupUserRemove,
		CmdRun:     RunGroupUserRemove,
		InitClient: true,
	})
	removeUser.AddStringSliceFlag(constants.ArgCols, "", defaultUserCols, printer.ColsMessage(defaultUserCols))
	_ = removeUser.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultUserCols, cobra.ShellCompDirectiveNoFileComp
	})
	removeUser.AddUUIDFlag(cloudapiv6.ArgGroupId, "", "", cloudapiv6.GroupId, core.RequiredFlagOption())
	_ = removeUser.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GroupsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	removeUser.AddUUIDFlag(cloudapiv6.ArgUserId, cloudapiv6.ArgIdShort, "", cloudapiv6.UserId, core.RequiredFlagOption())
	_ = removeUser.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GroupUsersIds(os.Stderr, viper.GetString(core.GetFlagName(removeUser.NS, cloudapiv6.ArgGroupId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeUser.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Remove all Users fro a group.")

	return groupUserCmd
}

func RunGroupUserList(c *core.CommandConfig) error {
	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	users, resp, err := c.CloudApiV6Services.Groups().ListUsers(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId)), listQueryParams)
	if resp != nil {
		c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getUserPrint(nil, c, getGroupUsers(users)))
}

func PreRunGroupUserRemove(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgGroupId, cloudapiv6.ArgUserId},
		[]string{cloudapiv6.ArgGroupId, cloudapiv6.ArgAll},
	)
}

func RunGroupUserAdd(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	id := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId))
	groupId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId))
	c.Printer.Verbose("User with id: %v is adding to group with id: %v...", id, groupId)
	u := resources.User{
		User: ionoscloud.User{
			Id: &id,
		},
	}
	userAdded, resp, err := c.CloudApiV6Services.Groups().AddUser(groupId, u, queryParams)
	if resp != nil && printer.GetId(resp) != "" {
		c.Printer.Verbose(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getUserPrint(resp, c, getUser(userAdded)))
}

func RunGroupUserRemove(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := RemoveAllUsers(c); err != nil {
			return err
		}
		return c.Printer.Print(printer.Result{Resource: c.Resource, Verb: c.Verb})
	} else {
		if err := utils.AskForConfirm(c.Stdin, c.Printer, "remove user from group"); err != nil {
			return err
		}
		userId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId))
		groupId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId))
		c.Printer.Verbose("User with id: %v is adding to group with id: %v...", userId, groupId)
		resp, err := c.CloudApiV6Services.Groups().RemoveUser(groupId, userId, queryParams)
		if resp != nil && printer.GetId(resp) != "" {
			c.Printer.Verbose(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime)
		}
		if err != nil {
			return err
		}
		return c.Printer.Print(getGroupPrint(resp, c, nil))
	}
}

func RemoveAllUsers(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	groupId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId))
	c.Printer.Verbose("Group ID: %v", groupId)
	c.Printer.Verbose("Getting Users...")
	users, resp, err := c.CloudApiV6Services.Groups().ListUsers(groupId, cloudapiv6.ParentResourceListQueryParams)
	if err != nil {
		return err
	}
	if usersItems, ok := users.GetItemsOk(); ok && usersItems != nil {
		if len(*usersItems) > 0 {
			_ = c.Printer.Warn("Users to be removed:")
			for _, user := range *usersItems {
				delIdAndName := ""
				if id, ok := user.GetIdOk(); ok && id != nil {
					delIdAndName += "User Id: " + *id
				}
				if properties, ok := user.GetPropertiesOk(); ok && properties != nil {
					if firstName, ok := properties.GetFirstnameOk(); ok && firstName != nil {
						delIdAndName += " User First Name: " + *firstName
					}
					if lastName, ok := properties.GetLastnameOk(); ok && lastName != nil {
						delIdAndName += " User Last Name: " + *lastName
					}
				}
				_ = c.Printer.Warn(delIdAndName)
			}
			if err := utils.AskForConfirm(c.Stdin, c.Printer, "removing all the Users"); err != nil {
				return err
			}
			c.Printer.Verbose("Removing all the Users...")
			var multiErr error
			for _, user := range *usersItems {
				if id, ok := user.GetIdOk(); ok && id != nil {
					c.Printer.Verbose("Starting removing User with id: %v...", *id)
					resp, err = c.CloudApiV6Services.Groups().RemoveUser(groupId, *id, queryParams)
					if resp != nil && printer.GetId(resp) != "" {
						c.Printer.Verbose(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime)
					}
					if err != nil {
						multiErr = multierr.Append(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
						continue
					} else {
						_ = c.Printer.Warn(fmt.Sprintf(constants.MessageDeletingAll, c.Resource, *id))
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
			return errors.New("no Users found")
		}
	} else {
		return errors.New("could not get items of Users")
	}
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
			r.WaitForRequest = viper.GetBool(core.GetFlagName(c.NS, constants.ArgWaitForRequest))
		}
		if users != nil {
			r.OutputJSON = users
			r.KeyValue = getUsersKVMaps(users)
			if c.Resource != c.Namespace {
				r.Columns = printer.GetHeadersAllDefault(defaultUserCols, viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols)))
			} else {
				r.Columns = printer.GetHeadersAllDefault(defaultUserCols, viper.GetStringSlice(core.GetGlobalFlagName(c.NS, constants.ArgCols)))
			}
		}
	}
	return r
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
