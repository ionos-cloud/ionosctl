package volume

import (
	"errors"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/viper"
)

func PreRunServerVolumeList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId); err != nil {
		return err
	}
	return nil
}

func PreRunDcServerVolumeIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId, cloudapiv6.ArgVolumeId)
}

func PreRunDcServerVolumeDetach(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId, cloudapiv6.ArgVolumeId},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId, cloudapiv6.ArgAll},
	)
}

func RunServerVolumeAttach(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))
	volumeId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId))

	c.Verbose(constants.DatacenterId, dcId)
	c.Verbose("Attaching Volume with ID: %v to Server with ID: %v...", volumeId, serverId)

	attachedVol, resp, err := c.CloudApiV6Services.Servers().AttachVolume(dcId, serverId, volumeId)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allVolumeCols).Print(attachedVol.Volume)
}

func RunServerVolumesList(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))

	c.Verbose(constants.DatacenterId, dcId)
	c.Verbose("Listing attached Volumes from Server with ID: %v...", serverId)

	attachedVols, resp, err := c.CloudApiV6Services.Servers().ListVolumes(dcId, serverId)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allVolumeCols).Prefix("items").Print(attachedVols.AttachedVolumes)
}

func RunServerVolumeGet(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))
	volumeId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId))

	c.Verbose(constants.DatacenterId, dcId)
	c.Verbose("Getting attached Volume with ID: %v from Server with ID: %v...", volumeId, serverId)

	attachedVol, resp, err := c.CloudApiV6Services.Servers().GetVolume(dcId, serverId, volumeId)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allVolumeCols).Print(attachedVol.Volume)
}

func RunServerVolumeDetach(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DetachAllServerVolumes(c); err != nil {
			return err
		}

		return nil
	}
	if !confirm.FAsk(c.Command.Command.InOrStdin(), "detach volume from server", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))
	volumeId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId))

	c.Verbose(constants.DatacenterId, dcId)
	c.Verbose("Detaching Volume with ID: %v from Server with ID: %v...", volumeId, serverId)

	resp, err := c.CloudApiV6Services.Servers().DetachVolume(dcId, serverId, volumeId)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("Volume successfully detached")
	return nil
}

func DetachAllServerVolumes(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))

	c.Verbose(constants.DatacenterId, dcId)
	c.Verbose("Server ID: %v", serverId)
	c.Verbose("Getting Volumes...")

	volumes, resp, err := c.CloudApiV6Services.Servers().ListVolumes(dcId, serverId)
	if err != nil {
		return err
	}

	volumesItems, ok := volumes.GetItemsOk()
	if !ok || volumesItems == nil {
		return fmt.Errorf("could not get items of Volumes")
	}

	if len(*volumesItems) <= 0 {
		return fmt.Errorf("no Volumes found")
	}

	c.Msg("Volumes to be detached:")

	var multiErr error
	for _, volume := range *volumesItems {
		id := volume.GetId()
		name := volume.GetProperties().GetName()

		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Detach the Volume with Id: %s, Name: %s", *id, *name), viper.GetBool(constants.ArgForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		resp, err = c.CloudApiV6Services.Servers().DetachVolume(dcId, serverId, *id)
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
