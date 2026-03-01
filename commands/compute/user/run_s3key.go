package user

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

func PreRunUserKeyIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgUserId, cloudapiv6.ArgS3KeyId)
}

func PreRunUserKeyDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgUserId, cloudapiv6.ArgS3KeyId},
		[]string{cloudapiv6.ArgUserId, cloudapiv6.ArgAll},
	)
}

func RunUserS3KeyList(c *core.CommandConfig) error {

	ss, resp, err := c.CloudApiV6Services.S3Keys().List(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId)))
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("items", jsonpaths.S3Key, ss.S3Keys,
		tabheaders.GetHeadersAllDefault(defaultS3KeyCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func RunUserS3KeyGet(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"S3 keys with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgS3KeyId))))

	s, resp, err := c.CloudApiV6Services.S3Keys().Get(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgS3KeyId)),
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.S3Key, s.S3Key,
		tabheaders.GetHeadersAllDefault(defaultS3KeyCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func RunUserS3KeyCreate(c *core.CommandConfig) error {
	userId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Creating S3 Key for User with ID: %v", userId))

	s, resp, err := c.CloudApiV6Services.S3Keys().Create(userId)
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

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.S3Key, s.S3Key,
		tabheaders.GetHeadersAllDefault(defaultS3KeyCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func RunUserS3KeyUpdate(c *core.CommandConfig) error {
	active := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgS3KeyActive))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property Active set: %v", active))

	newKey := resources.S3Key{
		S3Key: ionoscloud.S3Key{
			Properties: &ionoscloud.S3KeyProperties{
				Active: &active,
			},
		},
	}

	userId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId))
	keyId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgS3KeyId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Creating S3 Key with ID: %v for User with ID: %v", keyId, userId))

	s, resp, err := c.CloudApiV6Services.S3Keys().Update(userId, keyId, newKey)
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

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.S3Key, s.S3Key,
		tabheaders.GetHeadersAllDefault(defaultS3KeyCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func RunUserS3KeyDelete(c *core.CommandConfig) error {
	userId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId))
	s3KeyId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgS3KeyId))

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllS3keys(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete s3key", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("User ID: %v", userId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Starting deleting S3 Key with ID: %v...", s3KeyId))

	resp, err := c.CloudApiV6Services.S3Keys().Delete(userId, s3KeyId)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", jsontabwriter.GenerateLogOutput("S3 Key successfully deleted"))
	return nil
}

func DeleteAllS3keys(c *core.CommandConfig) error {
	userId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("User ID: %v", userId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Getting S3 Keys..."))

	s3Keys, resp, err := c.CloudApiV6Services.S3Keys().List(userId)
	if err != nil {
		return err
	}

	s3KeysItems, ok := s3Keys.GetItemsOk()
	if !ok || s3KeysItems == nil {
		return fmt.Errorf("could not get items of S3 Keys")
	}

	if len(*s3KeysItems) <= 0 {
		return fmt.Errorf("no S3 Keys found")
	}

	var multiErr error
	for _, s3Key := range *s3KeysItems {
		id := s3Key.GetId()
		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Delete the S3Key with Id: %s", *id), viper.GetBool(constants.ArgForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		resp, err = c.CloudApiV6Services.S3Keys().Delete(userId, *id)
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
