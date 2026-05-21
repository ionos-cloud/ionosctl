package backupunit

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

func PreRunBackupUnitId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgBackupUnitId)
}

func PreRunBackupUnitDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgBackupUnitId},
		[]string{cloudapiv6.ArgAll},
	)
}

func PreRunBackupUnitNameEmailPwd(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgName, cloudapiv6.ArgEmail, cloudapiv6.ArgPassword)
}

func RunBackupUnitList(c *core.CommandConfig) error {
	backupUnits, resp, err := c.CloudApiV6Services.BackupUnit().List()
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allCols).Prefix("items").Print(backupUnits.BackupUnits)
}

func RunBackupUnitGet(c *core.CommandConfig) error {
	c.Verbose("Backup unit with id: %v is getting... ", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgBackupUnitId)))

	u, resp, err := c.CloudApiV6Services.BackupUnit().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgBackupUnitId)))
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(u.BackupUnit)
}

func RunBackupUnitGetSsoUrl(c *core.CommandConfig) error {
	c.Verbose("Backup unit with id: %v is getting... ", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgBackupUnitId)))

	u, resp, err := c.CloudApiV6Services.BackupUnit().GetSsoUrl(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgBackupUnitId)))
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allSSOUrlCols).Print(u.BackupUnitSSO)
}

func RunBackupUnitCreate(c *core.CommandConfig) error {
	name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
	email := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgEmail))
	pwd := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPassword))

	newBackupUnit := resources.BackupUnit{
		BackupUnit: ionoscloud.BackupUnit{
			Properties: &ionoscloud.BackupUnitProperties{
				Name:     &name,
				Email:    &email,
				Password: &pwd,
			},
		},
	}

	c.Verbose("Properties set for creating the Backup Unit: Name: %v , Email: %v", name, email)

	u, resp, err := c.CloudApiV6Services.BackupUnit().Create(newBackupUnit)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg(backupUnitNote)

	return c.Printer(allCols).Print(u.BackupUnit)
}

func RunBackupUnitUpdate(c *core.CommandConfig) error {
	newProperties := getBackupUnitInfo(c)

	backupUnitUpd, resp, err := c.CloudApiV6Services.BackupUnit().Update(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgBackupUnitId)), *newProperties)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(backupUnitUpd.BackupUnit)
}

func RunBackupUnitDelete(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllBackupUnits(c); err != nil {
			return err
		}

		return nil
	}

	backupunitId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgBackupUnitId))
	backupunitDetails, _, err := c.CloudApiV6Services.BackupUnit().Get(backupunitId)
	if err != nil {
		return err
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("deleting Backup unit with id: %v, name: %s", backupunitId, *backupunitDetails.Properties.GetName()), viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	resp, err := c.CloudApiV6Services.BackupUnit().Delete(backupunitId)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("Backup Unit successfully deleted")

	return nil

}

func getBackupUnitInfo(c *core.CommandConfig) *resources.BackupUnitProperties {
	var properties resources.BackupUnitProperties
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPassword)) {
		pwd := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPassword))
		properties.SetPassword(pwd)

		c.Verbose("Property Password set")
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgEmail)) {
		email := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgEmail))
		properties.SetEmail(email)

		c.Verbose("Property Email set: %v", email)
	}

	return &properties
}

func DeleteAllBackupUnits(c *core.CommandConfig) error {
	c.Verbose("Getting Backup Units...")

	backupUnits, resp, err := c.CloudApiV6Services.BackupUnit().List()
	if err != nil {
		return err
	}

	backupUnitsItems, ok := backupUnits.GetItemsOk()
	if !ok || backupUnitsItems == nil {
		return fmt.Errorf("could not get Backup Unit items")
	}

	if len(*backupUnitsItems) <= 0 {
		return fmt.Errorf("no Backup Units found")
	}

	c.Msg("Backup Units to be deleted:")

	var multiErr error
	for _, backupUnit := range *backupUnitsItems {
		id := backupUnit.GetId()
		name := backupUnit.GetProperties().Name

		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Delete BackupUnit Id: %s , Name: %s ", *id, *name), viper.GetBool(constants.ArgForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		resp, err = c.CloudApiV6Services.BackupUnit().Delete(*id)
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
