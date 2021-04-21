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

func user() *builder.Command {
	ctx := context.TODO()
	userCmd := &builder.Command{
		Command: &cobra.Command{
			Use:              "user",
			Short:            "User Operations",
			Long:             `The sub-command of ` + "`" + `ionosctl user` + "`" + ` allows you to list, get, create, update, delete Users under your account. To add Users to a Group, check the ` + "`" + `ionosctl group` + "`" + ` commands.`,
			TraverseChildren: true,
		},
	}
	globalFlags := userCmd.Command.PersistentFlags()
	globalFlags.StringSlice(config.ArgCols, defaultUserCols, "Columns to be printed in the standard output")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(userCmd.Command.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))

	/*
		List Command
	*/
	builder.NewCommand(ctx, userCmd, noPreRun, RunUserList, "list", "List Users",
		"Use this command to get a list of existing Users available on your account.", listUserExample, true)

	/*
		Get Command
	*/
	get := builder.NewCommand(ctx, userCmd, PreRunUserIdValidate, RunUserGet, "get", "Get a User",
		"Use this command to retrieve details about a specific User.\n\nRequired values to run command:\n\n* User Id", getUserExample, true)
	get.AddStringFlag(config.ArgUserId, "", "", config.RequiredFlagUserId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getUsersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := builder.NewCommand(ctx, userCmd, PreRunUserNameEmailPwdValidate, RunUserCreate, "create", "Create a User under a particular contract",
		`Use this command to create a User under a particular contract. You need to specify the firstname, lastname, email and password for the new User.

Note: The password set here cannot be updated through the API currently. It is recommended that a new User log into the DCD and change their password.

Required values to run a command:

* User First Name
* User Last Name
* User Email
* User Password`, createUserExample, true)
	create.AddStringFlag(config.ArgUserFirstName, "", "", "The firstname for the User "+config.RequiredFlag)
	create.AddStringFlag(config.ArgUserLastName, "", "", "The lastname for the User "+config.RequiredFlag)
	create.AddStringFlag(config.ArgUserEmail, "", "", "The email for the User "+config.RequiredFlag)
	create.AddStringFlag(config.ArgUserPassword, "", "", "The password for the User (must be at least 5 characters long) "+config.RequiredFlag)
	create.AddBoolFlag(config.ArgUserAdministrator, "", false, "Assigns the User to have administrative rights")
	create.AddBoolFlag(config.ArgUserForceSecAuth, "", false, "Indicates if secure (two-factor) authentication should be forced for the User")

	/*
		Update Command
	*/
	update := builder.NewCommand(ctx, userCmd, noPreRun, RunUserUpdate, "update", "Update a User",
		`Use this command to update details about a specific User including their privileges.

Note: The password attribute is immutable. It is not allowed in update requests. It is recommended that the new User log into the DCD and change their password.

Required values to run command:

* User Id`, updateUserExample, true)
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
	deleteCmd := builder.NewCommand(ctx, userCmd, PreRunUserIdValidate, RunUserDelete, "delete", "Blacklists the User, disabling them",
		`This command blacklists the User, disabling them. The User is not completely purged, therefore if you anticipate needing to create a User with the same name in the future, we suggest renaming the User before you delete it.

Required values to run command:

* User Id`, deleteUserExample, true)
	deleteCmd.AddStringFlag(config.ArgUserId, "", "", config.RequiredFlagUserId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getUsersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	return userCmd
}

func PreRunUserIdValidate(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgUserId)
}

func PreRunUserNameEmailPwdValidate(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgUserFirstName, config.ArgUserLastName, config.ArgUserEmail, config.ArgUserPassword)
}

func RunUserList(c *builder.CommandConfig) error {
	users, _, err := c.Users().List()
	if err != nil {
		return err
	}
	return c.Printer.Print(getUserPrint(nil, c, getUsers(users)))
}

func RunUserGet(c *builder.CommandConfig) error {
	u, _, err := c.Users().Get(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgUserId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getUserPrint(nil, c, getUser(u)))
}

func RunUserCreate(c *builder.CommandConfig) error {
	firstname := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgUserFirstName))
	lastname := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgUserLastName))
	email := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgUserEmail))
	pwd := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgUserPassword))
	secureAuth := viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgUserForceSecAuth))
	admin := viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgUserAdministrator))
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

func RunUserUpdate(c *builder.CommandConfig) error {
	oldUser, resp, err := c.Users().Get(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgUserId)))
	if err != nil {
		return err
	}
	newProperties := getUserInfo(oldUser, c)
	newUser := resources.User{
		User: ionoscloud.User{
			Properties: &newProperties.UserProperties,
		},
	}
	userUpd, resp, err := c.Users().Update(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgUserId)), newUser)
	if err != nil {
		return err
	}
	return c.Printer.Print(getUserPrint(resp, c, getUser(userUpd)))
}

func RunUserDelete(c *builder.CommandConfig) error {
	err := utils.AskForConfirm(c.Stdin, c.Printer, "delete user")
	if err != nil {
		return err
	}
	resp, err := c.Users().Delete(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgUserId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getUserPrint(resp, c, nil))
}

func getUserInfo(oldUser *resources.User, c *builder.CommandConfig) *resources.UserProperties {
	var (
		firstName, lastName, email string
		forceSecureAuth, admin     bool
	)
	if properties, ok := oldUser.GetPropertiesOk(); ok && properties != nil {
		if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgUserFirstName)) {
			firstName = viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgUserFirstName))
		} else {
			if name, ok := properties.GetFirstnameOk(); ok && name != nil {
				firstName = *name
			}
		}
		if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgUserLastName)) {
			lastName = viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgUserLastName))
		} else {
			if name, ok := properties.GetLastnameOk(); ok && name != nil {
				lastName = *name
			}
		}
		if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgUserEmail)) {
			email = viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgUserEmail))
		} else {
			if e, ok := properties.GetEmailOk(); ok && e != nil {
				email = *e
			}
		}
		if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgUserForceSecAuth)) {
			forceSecureAuth = viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgUserForceSecAuth))
		} else {
			if secAuth, ok := properties.GetForceSecAuthOk(); ok && secAuth != nil {
				forceSecureAuth = *secAuth
			}
		}
		if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgUserAdministrator)) {
			admin = viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgUserAdministrator))
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

func getUserPrint(resp *resources.Response, c *builder.CommandConfig, users []resources.User) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.ParentName
			r.Verb = c.Name
			r.WaitFlag = viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait))
		}
		if users != nil {
			r.OutputJSON = users
			r.KeyValue = getUsersKVMaps(users)
			r.Columns = getUserCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr())
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
