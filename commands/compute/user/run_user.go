package user

import (
	"errors"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
)

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
	users, resp, err := c.CloudApiV6Services.Users().List()
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allUserCols).Prefix("items").Print(users.Users)
}

func RunUserGet(c *core.CommandConfig) error {
	c.Verbose("User with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId)))

	u, resp, err := c.CloudApiV6Services.Users().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId)))
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allUserCols).Print(u.User)
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

	c.Verbose("Properties set for creating the user: Firstname: %v, Lastname: %v, Email: %v, ForceSecAuth: %v, Administrator: %v",
		firstname, lastname, email, secureAuth, admin)
	c.Verbose("Creating User...")

	u, resp, err := c.CloudApiV6Services.Users().Create(newUser)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allUserCols).Print(u.User)
}

func RunUserUpdate(c *core.CommandConfig) error {
	oldUser, resp, err := c.CloudApiV6Services.Users().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId)))
	if err != nil {
		return err
	}

	newUser := getUserInfo(oldUser, c)

	c.Verbose("Updating User with ID: %v...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId)))

	userUpd, resp, err := c.CloudApiV6Services.Users().Update(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId)), *newUser)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allUserCols).Print(userUpd.User)
}

func RunUserDelete(c *core.CommandConfig) error {
	userId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId))

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllUsers(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete user", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	c.Verbose("Starting deleting User with id: %v...", userId)

	resp, err := c.CloudApiV6Services.Users().Delete(userId)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("User successfully deleted")

	return nil

}

func getUserInfo(oldUser *resources.User, c *core.CommandConfig) *resources.UserPut {
	userPropertiesPut := ionoscloud.UserPropertiesPut{}

	if properties, ok := oldUser.GetPropertiesOk(); ok && properties != nil {
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgFirstName)) {
			firstName := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFirstName))

			c.Verbose("Property FirstName set: %v", firstName)

			userPropertiesPut.SetFirstname(firstName)
		} else {
			if firstnameOk, ok := properties.GetFirstnameOk(); ok && firstnameOk != nil {
				userPropertiesPut.SetFirstname(*firstnameOk)
			}
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgLastName)) {
			lastName := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLastName))

			c.Verbose("Property LastName set: %v", lastName)

			userPropertiesPut.SetLastname(lastName)
		} else {
			if lastnameOk, ok := properties.GetLastnameOk(); ok && lastnameOk != nil {
				userPropertiesPut.SetLastname(*lastnameOk)
			}
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgEmail)) {
			email := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgEmail))

			c.Verbose("Property Email set: %v", email)

			userPropertiesPut.SetEmail(email)
		} else {
			if emailOk, ok := properties.GetEmailOk(); ok && emailOk != nil {
				userPropertiesPut.SetEmail(*emailOk)
			}
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPassword)) {
			password := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPassword))

			c.Verbose("Property Password set: %v", password)

			userPropertiesPut.SetPassword(password)
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgForceSecAuth)) {
			forceSecureAuth := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgForceSecAuth))

			c.Verbose("Property ForceSecAuth set: %v", forceSecureAuth)

			userPropertiesPut.SetForceSecAuth(forceSecureAuth)
		} else {
			if secAuthOk, ok := properties.GetForceSecAuthOk(); ok && secAuthOk != nil {
				userPropertiesPut.SetForceSecAuth(*secAuthOk)
			}
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgAdmin)) {
			admin := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAdmin))

			c.Verbose("Property Administrator set: %v", admin)

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
	return core.DeleteAll(c, core.DeleteAllOptions[ionoscloud.User]{
		Resource: "User",
		List: func() ([]ionoscloud.User, error) {
			users, _, err := c.CloudApiV6Services.Users().List()
			if err != nil {
				return nil, err
			}

			items, ok := users.GetItemsOk()
			if !ok || items == nil {
				return nil, fmt.Errorf("could not get items of Users")
			}

			return *items, nil
		},
		Summary: func(user ionoscloud.User) string {
			summary := fmt.Sprintf("id: %s", *user.GetId())
			if props, ok := user.GetPropertiesOk(); ok && props != nil {
				if firstname, ok := props.GetFirstnameOk(); ok && firstname != nil && *firstname != "" {
					summary = fmt.Sprintf("%s, firstname: %s", summary, *firstname)
				}
				if lastname, ok := props.GetLastnameOk(); ok && lastname != nil && *lastname != "" {
					summary = fmt.Sprintf("%s, lastname: %s", summary, *lastname)
				}
				if email, ok := props.GetEmailOk(); ok && email != nil && *email != "" {
					summary = fmt.Sprintf("%s, email: %s", summary, *email)
				}
			}
			return summary
		},
		ID: func(user ionoscloud.User) string {
			return *user.GetId()
		},
		Delete: func(user ionoscloud.User) error {
			_, err := c.CloudApiV6Services.Users().Delete(*user.GetId())
			return err
		},
	})
}

func RunGroupUserList(c *core.CommandConfig) error {
	users, resp, err := c.CloudApiV6Services.Groups().ListUsers(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId)))
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allUserCols).Prefix("items").Print(users.GroupMembers)
}

func PreRunGroupUserRemove(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgGroupId, cloudapiv6.ArgUserId},
		[]string{cloudapiv6.ArgGroupId, cloudapiv6.ArgAll},
	)
}

func RunGroupUserAdd(c *core.CommandConfig) error {
	id := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId))
	groupId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId))

	c.Verbose("User with id: %v is adding to group with id: %v...", id, groupId)

	u := ionoscloud.UserGroupPost{
		Id: &id,
	}

	userAdded, resp, err := c.CloudApiV6Services.Groups().AddUser(groupId, u)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allUserCols).Print(userAdded.User)
}

func RunGroupUserRemove(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := RemoveAllUsers(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "remove user from group", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	userId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId))
	groupId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId))

	c.Verbose("User with id: %v is adding to group with id: %v...", userId, groupId)

	resp, err := c.CloudApiV6Services.Groups().RemoveUser(groupId, userId)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("User successfully deleted")

	return nil

}

func RemoveAllUsers(c *core.CommandConfig) error {
	groupId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId))

	c.Verbose("Group ID: %v", groupId)
	c.Verbose("Getting Users...")

	users, resp, err := c.CloudApiV6Services.Groups().ListUsers(groupId)
	if err != nil {
		return err
	}

	usersItems, ok := users.GetItemsOk()
	if !ok || usersItems == nil {
		return fmt.Errorf("could not get items of Users")
	}

	if len(*usersItems) <= 0 {
		return fmt.Errorf("no Users found")
	}

	c.Msg("Users to be removed:")

	var multiErr error
	for _, user := range *usersItems {
		id := user.GetId()
		firstname := user.GetProperties().GetFirstname()
		lastname := user.GetProperties().GetLastname()

		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Remove the User with Id: %s, LastName: %s, FirstName: %s", *id, *lastname, *firstname), viper.GetBool(constants.ArgForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		resp, err = c.CloudApiV6Services.Groups().RemoveUser(groupId, *id)
		if resp != nil && request.GetId(resp) != "" {
			c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
		}
		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
			continue
		}

	}

	if multiErr != nil {
		return multiErr
	}

	return nil
}
