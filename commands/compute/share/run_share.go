package share

import (
	"errors"
	"fmt"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/helpers"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/waiter"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
	"github.com/ionos-cloud/ionosctl/v6/internal/waitfor"
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
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, totalTime))
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput(
		"*.items", jsonpaths.Share, allShares, tabheaders.GetHeaders(allGroupShareCols, defaultGroupShareCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}

func RunShareList(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		return RunShareListAll(c)
	}

	shares, resp, err := c.CloudApiV6Services.Groups().ListShares(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId)))
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("items", jsonpaths.Share, shares.GroupShares,
		tabheaders.GetHeadersAllDefault(defaultGroupShareCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
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
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Getting Share with Resource ID: %v from Group with ID: %v...",
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId))))

	s, resp, err := c.CloudApiV6Services.Groups().GetShare(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceId)),
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.Share, s.GroupShare,
		tabheaders.GetHeadersAllDefault(defaultGroupShareCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
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

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Properties set for creating the Share: EditPrivilege: %v, SharePrivilege: %v", editPrivilege, sharePrivilege))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Adding Share for Resource ID: %v from Group with ID: %v...",
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId))))

	shareAdded, resp, err := c.CloudApiV6Services.Groups().AddShare(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceId)),
		input,
	)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.Share, shareAdded.GroupShare,
		tabheaders.GetHeadersAllDefault(defaultGroupShareCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
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

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Updating Share for Resource ID: %v from Group with ID: %v...",
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId))))

	shareUpdated, resp, err := c.CloudApiV6Services.Groups().UpdateShare(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceId)),
		newShare,
	)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.Share, shareUpdated.GroupShare,
		tabheaders.GetHeadersAllDefault(defaultGroupShareCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
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

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Starting deleting Share with Resource ID: %v from Group with ID: %v...", shareId, groupId))

	resp, err := c.CloudApiV6Services.Groups().RemoveShare(groupId, shareId)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", jsontabwriter.GenerateLogOutput("Group Share successfully deleted"))
	return nil
}

func getShareUpdateInfo(oldShare *resources.GroupShare, c *core.CommandConfig) *resources.GroupShareProperties {
	var sharePrivilege, editPrivilege bool

	if properties, ok := oldShare.GetPropertiesOk(); ok && properties != nil {
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgEditPrivilege)) {
			editPrivilege = viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgEditPrivilege))

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property EditPrivilege set: %v", editPrivilege))
		} else {
			if e, ok := properties.GetEditPrivilegeOk(); ok && e != nil {
				editPrivilege = *e
			}
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgSharePrivilege)) {
			sharePrivilege = viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgSharePrivilege))

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property SharePrivilege set: %v", sharePrivilege))
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

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Group ID: %v", groupId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Getting Group Shares..."))

	groupShares, resp, err := c.CloudApiV6Services.Groups().ListShares(groupId)
	if err != nil {
		return err
	}

	groupSharesItems, ok := groupShares.GetItemsOk()
	if !ok || groupSharesItems == nil {
		return fmt.Errorf("could not get items of Group Shares")
	}

	if len(*groupSharesItems) <= 0 {
		return fmt.Errorf("no Group Shares found")
	}

	var multiErr error
	for _, share := range *groupSharesItems {
		id := share.GetId()
		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Delete the GroupShare with Id: %s", *id), viper.GetBool(constants.ArgForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		resp, err = c.CloudApiV6Services.Groups().RemoveShare(groupId, *id)
		if resp != nil && request.GetId(resp) != "" {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
		}
		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
			continue
		}

		if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrWaitDeleteAll, c.Resource, *id, err))

		}
	}

	if multiErr != nil {
		return multiErr
	}

	return nil
}
