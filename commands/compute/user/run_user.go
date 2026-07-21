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
	c.Verbose("User with id: %v is getting...", c.Flags().String(cloudapiv6.ArgUserId))

	u, resp, err := c.CloudApiV6Services.Users().Get(c.Flags().String(cloudapiv6.ArgUserId))
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allUserCols).Print(u.User)
}

func RunUserCreate(c *core.CommandConfig) error {
	firstname := c.Flags().String(cloudapiv6.ArgFirstName)
	lastname := c.Flags().String(cloudapiv6.ArgLastName)
	email := c.Flags().String(cloudapiv6.ArgEmail)
	pwd := c.Flags().String(cloudapiv6.ArgPassword)
	secureAuth := c.Flags().Bool(cloudapiv6.ArgForceSecAuth)
	admin := c.Flags().Bool(cloudapiv6.ArgAdmin)

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
	oldUser, resp, err := c.CloudApiV6Services.Users().Get(c.Flags().String(cloudapiv6.ArgUserId))
	if err != nil {
		return err
	}

	newUser := getUserInfo(oldUser, c)

	c.Verbose("Updating User with ID: %v...", c.Flags().String(cloudapiv6.ArgUserId))

	userUpd, resp, err := c.CloudApiV6Services.Users().Update(c.Flags().String(cloudapiv6.ArgUserId), *newUser)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allUserCols).Print(userUpd.User)
}

func RunUserDelete(c *core.CommandConfig) error {
	userId := c.Flags().String(cloudapiv6.ArgUserId)

	if c.Flags().Bool(cloudapiv6.ArgAll) {
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
		if c.Flags().Changed(cloudapiv6.ArgFirstName) {
			firstName := c.Flags().String(cloudapiv6.ArgFirstName)

			c.Verbose("Property FirstName set: %v", firstName)

			userPropertiesPut.SetFirstname(firstName)
		} else {
			if firstnameOk, ok := properties.GetFirstnameOk(); ok && firstnameOk != nil {
				userPropertiesPut.SetFirstname(*firstnameOk)
			}
		}

		if c.Flags().Changed(cloudapiv6.ArgLastName) {
			lastName := c.Flags().String(cloudapiv6.ArgLastName)

			c.Verbose("Property LastName set: %v", lastName)

			userPropertiesPut.SetLastname(lastName)
		} else {
			if lastnameOk, ok := properties.GetLastnameOk(); ok && lastnameOk != nil {
				userPropertiesPut.SetLastname(*lastnameOk)
			}
		}

		if c.Flags().Changed(cloudapiv6.ArgEmail) {
			email := c.Flags().String(cloudapiv6.ArgEmail)

			c.Verbose("Property Email set: %v", email)

			userPropertiesPut.SetEmail(email)
		} else {
			if emailOk, ok := properties.GetEmailOk(); ok && emailOk != nil {
				userPropertiesPut.SetEmail(*emailOk)
			}
		}

		if c.Flags().Changed(cloudapiv6.ArgPassword) {
			password := c.Flags().String(cloudapiv6.ArgPassword)

			c.Verbose("Property Password set: %v", password)

			userPropertiesPut.SetPassword(password)
		}

		if c.Flags().Changed(cloudapiv6.ArgForceSecAuth) {
			forceSecureAuth := c.Flags().Bool(cloudapiv6.ArgForceSecAuth)

			c.Verbose("Property ForceSecAuth set: %v", forceSecureAuth)

			userPropertiesPut.SetForceSecAuth(forceSecureAuth)
		} else {
			if secAuthOk, ok := properties.GetForceSecAuthOk(); ok && secAuthOk != nil {
				userPropertiesPut.SetForceSecAuth(*secAuthOk)
			}
		}

		if c.Flags().Changed(cloudapiv6.ArgAdmin) {
			admin := c.Flags().Bool(cloudapiv6.ArgAdmin)

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
	c.Verbose("Getting Users...")

	users, _, err := c.CloudApiV6Services.Users().List()
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

	var multiErr error
	for _, user := range *usersItems {
		id := user.GetId()
		lastname := user.GetProperties().Lastname
		firstname := user.GetProperties().Firstname

		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Delete the User with Id: %s, LastName: %s, FirstName: %s", *id, *lastname, *firstname), viper.GetBool(constants.ArgForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		_, err = c.CloudApiV6Services.Users().Delete(*id)
		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
			continue
		}

		c.Msg(constants.MessageDeletingAll, c.Resource, *id)

	}

	if multiErr != nil {
		return multiErr
	}

	return nil
}

func RunGroupUserList(c *core.CommandConfig) error {
	users, resp, err := c.CloudApiV6Services.Groups().ListUsers(c.Flags().String(cloudapiv6.ArgGroupId))
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
	id := c.Flags().String(cloudapiv6.ArgUserId)
	groupId := c.Flags().String(cloudapiv6.ArgGroupId)

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
	if c.Flags().Bool(cloudapiv6.ArgAll) {
		if err := RemoveAllUsers(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "remove user from group", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	userId := c.Flags().String(cloudapiv6.ArgUserId)
	groupId := c.Flags().String(cloudapiv6.ArgGroupId)

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
	groupId := c.Flags().String(cloudapiv6.ArgGroupId)

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
