package share

import (
	"fmt"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/helpers"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
)

func PreRunGroupResourceIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgGroupId, cloudapiv6.ArgResourceId)
}

func PreRunGroupResourceDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgGroupId, cloudapiv6.ArgResourceId},
		[]string{cloudapiv6.ArgGroupId, cloudapiv6.ArgAll},
	)
}

func RunShareListAll(c *core.CommandConfig) error {
	// Don't apply listQueryParams to parent resource, as it would have unexpected side effects on the results
	groups, _, err := c.CloudApiV6Services.Groups().List()
	if err != nil {
		return err
	}

	var allShares = make([]ionoscloud.GroupShares, 0)

	totalTime := time.Duration(0)

	for _, group := range helpers.GetGroups(groups) {
		shares, resp, err := c.CloudApiV6Services.Groups().ListShares(*group.GetId())
		if err != nil {
			return err
		}

		allShares = append(allShares, shares.GroupShares)
		totalTime += resp.RequestTime
	}

	if totalTime != time.Duration(0) {
		c.Verbose(constants.MessageRequestTime, totalTime)
	}

	return c.Printer(allGroupShareCols).Prefix("*.items").Print(allShares)
}

func RunShareList(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		return RunShareListAll(c)
	}

	shares, resp, err := c.CloudApiV6Services.Groups().ListShares(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId)))
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allGroupShareCols).Prefix("items").Print(shares.GroupShares)
}

func PreRunShareList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgGroupId},
		[]string{cloudapiv6.ArgAll},
	); err != nil {
		return err
	}

	return nil
}

func RunShareGet(c *core.CommandConfig) error {
	c.Verbose("Getting Share with Resource ID: %v from Group with ID: %v...",
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId)))

	s, resp, err := c.CloudApiV6Services.Groups().GetShare(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceId)),
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allGroupShareCols).Print(s.GroupShare)
}

func RunShareCreate(c *core.CommandConfig) error {
	editPrivilege := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgEditPrivilege))
	sharePrivilege := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgSharePrivilege))

	input := resources.GroupShare{
		GroupShare: ionoscloud.GroupShare{
			Properties: &ionoscloud.GroupShareProperties{
				EditPrivilege:  &editPrivilege,
				SharePrivilege: &sharePrivilege,
			},
		},
	}

	c.Verbose("Properties set for creating the Share: EditPrivilege: %v, SharePrivilege: %v", editPrivilege, sharePrivilege)
	c.Verbose("Adding Share for Resource ID: %v from Group with ID: %v...",
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId)))

	shareAdded, resp, err := c.CloudApiV6Services.Groups().AddShare(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceId)),
		input,
	)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allGroupShareCols).Print(shareAdded.GroupShare)
}

func RunShareUpdate(c *core.CommandConfig) error {
	s, _, err := c.CloudApiV6Services.Groups().GetShare(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceId)))
	if err != nil {
		return err
	}

	properties := getShareUpdateInfo(s, c)
	newShare := resources.GroupShare{
		GroupShare: ionoscloud.GroupShare{
			Properties: &properties.GroupShareProperties,
		},
	}

	c.Verbose("Updating Share for Resource ID: %v from Group with ID: %v...",
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId)))

	shareUpdated, resp, err := c.CloudApiV6Services.Groups().UpdateShare(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceId)),
		newShare,
	)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allGroupShareCols).Print(shareUpdated.GroupShare)
}

func RunShareDelete(c *core.CommandConfig) error {
	shareId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceId))
	groupId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId))

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllShares(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete share from group", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	c.Verbose("Starting deleting Share with Resource ID: %v from Group with ID: %v...", shareId, groupId)

	resp, err := c.CloudApiV6Services.Groups().RemoveShare(groupId, shareId)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("Group Share successfully deleted")
	return nil
}

func getShareUpdateInfo(oldShare *resources.GroupShare, c *core.CommandConfig) *resources.GroupShareProperties {
	var sharePrivilege, editPrivilege bool

	if properties, ok := oldShare.GetPropertiesOk(); ok && properties != nil {
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgEditPrivilege)) {
			editPrivilege = viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgEditPrivilege))

			c.Verbose("Property EditPrivilege set: %v", editPrivilege)
		} else {
			if e, ok := properties.GetEditPrivilegeOk(); ok && e != nil {
				editPrivilege = *e
			}
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgSharePrivilege)) {
			sharePrivilege = viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgSharePrivilege))

			c.Verbose("Property SharePrivilege set: %v", sharePrivilege)
		} else {
			if e, ok := properties.GetSharePrivilegeOk(); ok && e != nil {
				sharePrivilege = *e
			}
		}
	}

	return &resources.GroupShareProperties{
		GroupShareProperties: ionoscloud.GroupShareProperties{
			EditPrivilege:  &editPrivilege,
			SharePrivilege: &sharePrivilege,
		},
	}
}

func DeleteAllShares(c *core.CommandConfig) error {
	groupId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId))

	c.Verbose("Group ID: %v", groupId)

	return core.DeleteAll(c, core.DeleteAllOptions[ionoscloud.GroupShare]{
		Resource: "GroupShare",
		List: func() ([]ionoscloud.GroupShare, error) {
			groupShares, _, err := c.CloudApiV6Services.Groups().ListShares(groupId)
			if err != nil {
				return nil, err
			}

			items, ok := groupShares.GetItemsOk()
			if !ok || items == nil {
				return nil, fmt.Errorf("could not get items of Group Shares")
			}

			return *items, nil
		},
		Summary: func(share ionoscloud.GroupShare) string {
			id := ""
			if v, ok := share.GetIdOk(); ok && v != nil {
				id = *v
			}
			return fmt.Sprintf("id: %s", id)
		},
		ID: func(share ionoscloud.GroupShare) string {
			if id, ok := share.GetIdOk(); ok && id != nil {
				return *id
			}
			return ""
		},
		Delete: func(share ionoscloud.GroupShare) error {
			resp, err := c.CloudApiV6Services.Groups().RemoveShare(groupId, *share.GetId())
			if resp != nil && request.GetId(resp) != "" {
				c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
			}
			return err
		},
	})
}
