package targetgroup

import (
	"errors"
	"fmt"

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

func PreRunTargetGroupId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgTargetGroupId)
}

func PreRunTargetGroupDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgTargetGroupId},
		[]string{cloudapiv6.ArgAll},
	)
}

func RunTargetGroupList(c *core.CommandConfig) error {

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Getting TargetGroups"))
	ss, resp, err := c.CloudApiV6Services.TargetGroups().List()
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("items", jsonpaths.TargetGroup, ss.TargetGroups,
		tabheaders.GetHeaders(allTargetGroupCols, defaultTargetGroupCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func RunTargetGroupGet(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		constants.TargetGroupId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId))))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Getting TargetGroup"))

	s, resp, err := c.CloudApiV6Services.TargetGroups().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)))
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.TargetGroup, s.TargetGroup,
		tabheaders.GetHeaders(allTargetGroupCols, defaultTargetGroupCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func RunTargetGroupCreate(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Creating TargetGroup"))

	s, resp, err := c.CloudApiV6Services.TargetGroups().Create(getTargetGroupNew(c))
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.TargetGroup, s.TargetGroup,
		tabheaders.GetHeaders(allTargetGroupCols, defaultTargetGroupCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func RunTargetGroupUpdate(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		constants.TargetGroupId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId))))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Updating TargetGroup"))

	s, resp, err := c.CloudApiV6Services.TargetGroups().Update(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)), getTargetGroupPropertiesSet(c))
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.TargetGroup, s.TargetGroup,
		tabheaders.GetHeaders(allTargetGroupCols, defaultTargetGroupCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func RunTargetGroupDelete(c *core.CommandConfig) error {
	var resp *resources.Response

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
			constants.TargetGroupId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId))))
		err := DeleteAllTargetGroup(c)
		if err != nil {
			return err
		}

		return nil
	}
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		constants.TargetGroupId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId))))

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete target group", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Deleting TargetGroup"))

	resp, err := c.CloudApiV6Services.TargetGroups().Delete(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)))
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", jsontabwriter.GenerateLogOutput("Target Group successfully deleted"))
	return nil
}

func DeleteAllTargetGroup(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateLogOutput("Getting Target Groups..."))

	targetGroups, resp, err := c.CloudApiV6Services.TargetGroups().List()
	if err != nil {
		return err
	}

	targetGroupItems, ok := targetGroups.GetItemsOk()
	if !ok || targetGroupItems == nil {
		return fmt.Errorf("could not get items of Target Groups")
	}

	if len(*targetGroupItems) <= 0 {
		return fmt.Errorf("no Target Groups found")
	}

	var multiErr error
	for _, tg := range *targetGroupItems {
		id := tg.GetId()
		name := tg.GetProperties().Name

		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Delete the Target Group with Id: %s, Name: %s", *id, *name), viper.GetBool(constants.ArgForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		resp, err = c.CloudApiV6Services.TargetGroups().Delete(*id)
		if resp != nil && request.GetId(resp) != "" {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
		}
		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
			continue
		}

		if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
		}
	}

	if multiErr != nil {
		return multiErr
	}

	return nil
}

func getTargetGroupNew(c *core.CommandConfig) resources.TargetGroup {
	input := resources.TargetGroupProperties{}
	// Set Required Properties
	input.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Property Name set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))))

	input.SetAlgorithm(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgAlgorithm)))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Property Algorithm set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgAlgorithm))))

	input.SetProtocol(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol)))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Property Protocol set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol))))

	inputHealthCheck := resources.TargetGroupHealthCheck{}

	// Set Properties for Health Check for Target Group
	inputHealthCheck.SetCheckTimeout(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgCheckTimeout)))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Property CheckTimeout for HealthCheck set: %v", viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgCheckTimeout))))

	inputHealthCheck.SetCheckInterval(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgCheckInterval)))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Property CheckInterval for HealthCheck set: %v", viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgCheckInterval))))

	inputHealthCheck.SetRetries(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgRetries)))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Property Retries for HealthCheck set: %v", viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgRetries))))

	// Set Health Check for Target Group
	input.SetHealthCheck(inputHealthCheck.TargetGroupHealthCheck)
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Setting HealthCheck"))

	inputHttpHealthCheck := resources.TargetGroupHttpHealthCheck{}
	// Set Properties for Http Health Check for Target Group
	inputHttpHealthCheck.SetPath(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPath)))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Property Path for HttpHealthCheck set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPath))))

	inputHttpHealthCheck.SetMethod(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgMethod)))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Property Method for HttpHealthCheck set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgMethod))))

	inputHttpHealthCheck.SetMatchType(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgMatchType)))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Property MatchType for HttpHealthCheck set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgMatchType))))

	inputHttpHealthCheck.SetResponse(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResponse)))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Property Response for HttpHealthCheck set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResponse))))

	inputHttpHealthCheck.SetRegex(viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgRegex)))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Property Regex for HttpHealthCheck set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRegex))))

	inputHttpHealthCheck.SetNegate(viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgNegate)))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Property Negate for HttpHealthCheck set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNegate))))

	// Set Http Health Check for Target Group
	input.SetHttpHealthCheck(inputHttpHealthCheck.TargetGroupHttpHealthCheck)
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Setting HttpHealthCheck"))

	return resources.TargetGroup{
		TargetGroup: ionoscloud.TargetGroup{
			Properties: &input.TargetGroupProperties,
		},
	}
}

func getTargetGroupPropertiesSet(c *core.CommandConfig) *resources.TargetGroupProperties {
	input := resources.TargetGroupProperties{}
	// Set new values for Required Properties
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
		input.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
			"Property Name set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgAlgorithm)) {
		input.SetAlgorithm(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgAlgorithm)))

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
			"Property Algorithm set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgAlgorithm))))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol)) {
		input.SetProtocol(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol)))

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
			"Property Protocol set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol))))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgCheckTimeout)) ||
		viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgCheckInterval)) ||
		viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgRetries)) {
		inputHealthCheck := resources.TargetGroupHealthCheck{}

		// Set new values for Health Check Properties
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgCheckTimeout)) {
			inputHealthCheck.SetCheckTimeout(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgCheckTimeout)))

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
				"Property CheckTimeout for HealthCheck set: %v", viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgCheckTimeout))))
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgCheckInterval)) {
			inputHealthCheck.SetCheckInterval(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgCheckInterval)))
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
				"Property CheckInterval for HealthCheck set: %v", viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgCheckInterval))))
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgRetries)) {
			inputHealthCheck.SetRetries(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgRetries)))
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
				"Property Retries for HealthCheck set: %v", viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgRetries))))
		}

		input.SetHealthCheck(inputHealthCheck.TargetGroupHealthCheck)
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Updating HealthCheck"))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPath)) || viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgMethod)) ||
		viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgResponse)) || viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgRegex)) ||
		viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgNegate)) || viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgMatchType)) {
		inputHttpHealthCheck := resources.TargetGroupHttpHealthCheck{}

		// Set new values for Health Check Properties
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPath)) {
			inputHttpHealthCheck.SetPath(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPath)))
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
				"Property Path for HttpHealthCheck set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPath))))
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgMethod)) {
			inputHttpHealthCheck.SetMethod(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgMethod)))
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
				"Property Method for HttpHealthCheck set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgMethod))))
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgResponse)) {
			inputHttpHealthCheck.SetResponse(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResponse)))
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
				"Property Response for HttpHealthCheck set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResponse))))
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgMatchType)) {
			inputHttpHealthCheck.SetMatchType(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgMatchType)))
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
				"Property MatchType for HttpHealthCheck set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgMatchType))))
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgRegex)) {
			inputHttpHealthCheck.SetRegex(viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgRegex)))
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
				"Property Regex for HttpHealthCheck set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRegex))))
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgNegate)) {
			inputHttpHealthCheck.SetNegate(viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgNegate)))
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
				"Property Negate for HttpHealthCheck set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNegate))))
		}

		input.SetHttpHealthCheck(inputHttpHealthCheck.TargetGroupHttpHealthCheck)
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Updating HttpHealthCheck"))
	}

	return &input
}
