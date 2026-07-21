package group

import (
	"errors"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
)

func RunGroupResourceList(c *core.CommandConfig) error {
	c.Verbose("Listing Resources from Group with ID: %v...", c.Flags().String(cloudapiv6.ArgGroupId))

	resourcesListed, resp, err := c.CloudApiV6Services.Groups().ListResources(c.Flags().String(cloudapiv6.ArgGroupId))
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allResourceCols).Prefix("items").Print(resourcesListed.ResourceGroups)
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
		if err := removeAllUsersFromGroup(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "remove user from group", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	userId := c.Flags().String(cloudapiv6.ArgUserId)
	groupId := c.Flags().String(cloudapiv6.ArgGroupId)

	c.Verbose("User with id: %v is being removed from group with id: %v...", userId, groupId)

	resp, err := c.CloudApiV6Services.Groups().RemoveUser(groupId, userId)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("User successfully removed")

	return nil
}

func removeAllUsersFromGroup(c *core.CommandConfig) error {
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
