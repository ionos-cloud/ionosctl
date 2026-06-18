package user

import (
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
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allS3KeyCols).Prefix("items").Print(ss.S3Keys)
}

func RunUserS3KeyGet(c *core.CommandConfig) error {
	c.Verbose("S3 keys with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgS3KeyId)))

	s, resp, err := c.CloudApiV6Services.S3Keys().Get(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgS3KeyId)),
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allS3KeyCols).Print(s.S3Key)
}

func RunUserS3KeyCreate(c *core.CommandConfig) error {
	userId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId))

	c.Verbose("Creating S3 Key for User with ID: %v", userId)

	s, resp, err := c.CloudApiV6Services.S3Keys().Create(userId)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allS3KeyCols).Print(s.S3Key)
}

func RunUserS3KeyUpdate(c *core.CommandConfig) error {
	active := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgS3KeyActive))

	c.Verbose("Property Active set: %v", active)

	newKey := resources.S3Key{
		S3Key: ionoscloud.S3Key{
			Properties: &ionoscloud.S3KeyProperties{
				Active: &active,
			},
		},
	}

	userId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId))
	keyId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgS3KeyId))

	c.Verbose("Creating S3 Key with ID: %v for User with ID: %v", keyId, userId)

	s, resp, err := c.CloudApiV6Services.S3Keys().Update(userId, keyId, newKey)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allS3KeyCols).Print(s.S3Key)
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

	c.Verbose("User ID: %v", userId)
	c.Verbose("Starting deleting S3 Key with ID: %v...", s3KeyId)

	resp, err := c.CloudApiV6Services.S3Keys().Delete(userId, s3KeyId)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("S3 Key successfully deleted")
	return nil
}

func DeleteAllS3keys(c *core.CommandConfig) error {
	userId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserId))

	c.Verbose("User ID: %v", userId)

	return core.DeleteAll(c, core.DeleteAllOptions[ionoscloud.S3Key]{
		Resource: "S3Key",
		List: func() ([]ionoscloud.S3Key, error) {
			s3Keys, _, err := c.CloudApiV6Services.S3Keys().List(userId)
			if err != nil {
				return nil, err
			}

			items, ok := s3Keys.GetItemsOk()
			if !ok || items == nil {
				return nil, fmt.Errorf("could not get items of S3 Keys")
			}

			return *items, nil
		},
		Summary: func(s3Key ionoscloud.S3Key) string {
			return fmt.Sprintf("id: %s", *s3Key.GetId())
		},
		ID: func(s3Key ionoscloud.S3Key) string {
			return *s3Key.GetId()
		},
		Delete: func(s3Key ionoscloud.S3Key) error {
			resp, err := c.CloudApiV6Services.S3Keys().Delete(userId, *s3Key.GetId())
			if resp != nil && request.GetId(resp) != "" {
				c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
			}
			return err
		},
	})
}
