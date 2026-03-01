package label

import (
	"context"
	"errors"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/viper"
)

// Returns []core.FlagNameSetWithPredicate to be used as params to send to core.CheckRequiredFlagsSets funcs.
// If --resource-type datacenter, --datacenter-id is also required
// If --resource-type server,	  --datacenter-id and --server-id are also required
func generateFlagSets(c *core.PreCommandConfig, extraFlags ...string) []core.FlagNameSetWithPredicate {
	funcResourceTypeSetAndMatches := func(resource interface{}) bool {
		argResourceType := core.GetFlagName(c.NS, cloudapiv6.ArgResourceType)
		return !viper.IsSet(argResourceType) || viper.GetString(argResourceType) == resource
	}

	return []core.FlagNameSetWithPredicate{
		{
			FlagNameSet:    append([]string{cloudapiv6.ArgResourceType, cloudapiv6.ArgDataCenterId}, extraFlags...),
			Predicate:      funcResourceTypeSetAndMatches,
			PredicateParam: cloudapiv6.DatacenterResource,
		}, {
			FlagNameSet:    append([]string{cloudapiv6.ArgResourceType, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgVolumeId}, extraFlags...),
			Predicate:      funcResourceTypeSetAndMatches,
			PredicateParam: cloudapiv6.VolumeResource,
		}, {
			FlagNameSet:    append([]string{cloudapiv6.ArgResourceType, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId}, extraFlags...),
			Predicate:      funcResourceTypeSetAndMatches,
			PredicateParam: cloudapiv6.ServerResource,
		}, {
			FlagNameSet:    append([]string{cloudapiv6.ArgResourceType, cloudapiv6.ArgSnapshotId}, extraFlags...),
			Predicate:      funcResourceTypeSetAndMatches,
			PredicateParam: cloudapiv6.SnapshotResource,
		}, {
			FlagNameSet:    append([]string{cloudapiv6.ArgResourceType, cloudapiv6.ArgIpBlockId}, extraFlags...),
			Predicate:      funcResourceTypeSetAndMatches,
			PredicateParam: cloudapiv6.IpBlockResource,
		}, {
			FlagNameSet:    append([]string{cloudapiv6.ArgResourceType, cloudapiv6.ArgImageId}, extraFlags...),
			Predicate:      funcResourceTypeSetAndMatches,
			PredicateParam: cloudapiv6.ImageResource,
		},
	}
}

func PreRunResourceTypeLabelKey(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSetsIfPredicate(c.Command, c.NS, generateFlagSets(c, cloudapiv6.ArgLabelKey)...)
}

func PreRunResourceTypeLabelKeyRemove(c *core.PreCommandConfig) error {
	if all := viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll)); all {
		return nil
	}
	return core.CheckRequiredFlagsSetsIfPredicate(c.Command, c.NS,
		append(
			generateFlagSets(c, cloudapiv6.ArgLabelKey),
			generateFlagSets(c, cloudapiv6.ArgAll)...,
		)...,
	)
}

func PreRunResourceTypeLabelKeyValue(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSetsIfPredicate(c.Command, c.NS, generateFlagSets(c, cloudapiv6.ArgLabelKey, cloudapiv6.ArgLabelValue)...)
}

func PreRunLabelUrn(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgLabelUrn)
}

func PreRunLabelList(c *core.PreCommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgResourceType)) {
		return core.CheckRequiredFlagsSetsIfPredicate(c.Command, c.NS, generateFlagSets(c)...)
	}
	return core.NoPreRun(c)
}

func RunLabelList(c *core.CommandConfig) error {
	var out string
	switch viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceType)) {
	case cloudapiv6.DatacenterResource:
		return RunDataCenterLabelsList(c)
	case cloudapiv6.ServerResource:
		return RunServerLabelsList(c)
	case cloudapiv6.VolumeResource:
		return RunVolumeLabelsList(c)
	case cloudapiv6.IpBlockResource:
		return RunIpBlockLabelsList(c)
	case cloudapiv6.SnapshotResource:
		return RunSnapshotLabelsList(c)
	case cloudapiv6.ImageResource:
		return RunImageLabelsList(c)
	default:
		labelDcs, _, err := c.CloudApiV6Services.Labels().List()
		if err != nil {
			return err
		}

		cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

		out, err = jsontabwriter.GenerateOutput("items", jsonpaths.Label, labelDcs.Labels,
			tabheaders.GetHeadersAllDefault(defaultLabelCols, cols))
		if err != nil {
			return err
		}

		fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

		return nil
	}
}

func RunLabelGet(c *core.CommandConfig) error {
	resourceType := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceType))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey))
	labelValue := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelValue))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Getting label with label key: %v and label value: %v for %v...", labelKey, labelValue, resourceType))

	switch resourceType {
	case cloudapiv6.DatacenterResource:
		return RunDataCenterLabelGet(c)
	case cloudapiv6.ServerResource:
		return RunServerLabelGet(c)
	case cloudapiv6.VolumeResource:
		return RunVolumeLabelGet(c)
	case cloudapiv6.IpBlockResource:
		return RunIpBlockLabelGet(c)
	case cloudapiv6.SnapshotResource:
		return RunSnapshotLabelGet(c)
	case cloudapiv6.ImageResource:
		return RunImageLabelGet(c)
	default:
		fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", jsontabwriter.GenerateLogOutput(labelResourceWarning))

		return nil
	}
}

func RunLabelGetByUrn(c *core.CommandConfig) error {
	urn := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelUrn))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Getting label with urn: %v", urn))

	labelDc, _, err := c.CloudApiV6Services.Labels().GetByUrn(urn)
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.Label, labelDc.Label,
		tabheaders.GetHeadersAllDefault(defaultLabelCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}

func RunLabelAdd(c *core.CommandConfig) error {
	switch viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceType)) {
	case cloudapiv6.DatacenterResource:
		return RunDataCenterLabelAdd(c)
	case cloudapiv6.ServerResource:
		return RunServerLabelAdd(c)
	case cloudapiv6.VolumeResource:
		return RunVolumeLabelAdd(c)
	case cloudapiv6.IpBlockResource:
		return RunIpBlockLabelAdd(c)
	case cloudapiv6.SnapshotResource:
		return RunSnapshotLabelAdd(c)
	case cloudapiv6.ImageResource:
		return RunImageLabelAdd(c)
	default:
		fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", jsontabwriter.GenerateLogOutput(labelResourceWarning))

		return nil
	}
}

func RunLabelRemove(c *core.CommandConfig) error {
	resourceType := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceType))

	if all := viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll)) && resourceType == ""; all {
		return RunLabelRemoveAll(c)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("removing label from %v...", resourceType))

	switch resourceType {
	case cloudapiv6.DatacenterResource:
		return RunDataCenterLabelRemove(c)
	case cloudapiv6.ServerResource:
		return RunServerLabelRemove(c)
	case cloudapiv6.VolumeResource:
		return RunVolumeLabelRemove(c)
	case cloudapiv6.IpBlockResource:
		return RunIpBlockLabelRemove(c)
	case cloudapiv6.SnapshotResource:
		return RunSnapshotLabelRemove(c)
	case cloudapiv6.ImageResource:
		return RunImageLabelRemove(c)
	default:
		fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", jsontabwriter.GenerateLogOutput(labelResourceWarning))

		return nil
	}
}

func RunLabelRemoveAll(c *core.CommandConfig) error {
	labels, _, err := client.Must().CloudClient.LabelsApi.LabelsGet(context.Background()).Execute()

	var multiErr error
	for _, label := range *labels.GetItems() {
		key := *label.GetProperties().GetKey()
		value := *label.GetProperties().GetValue()
		resourceId := *label.GetProperties().GetResourceId()

		t := label.GetProperties().GetResourceType()

		if !confirm.FAsk(c.Command.Command.InOrStdin(),
			fmt.Sprintf("Delete Label with Id: %s, Key: '%s', Value: '%s'", resourceId, key, value),
			viper.GetBool(constants.ArgForce)) {
			continue
		}

		switch *t {
		case "datacenter":
			_, err = client.Must().CloudClient.LabelsApi.DatacentersLabelsDelete(context.Background(),
				resourceId, key).Execute()
			if err != nil {
				multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, resourceId, err))
				continue
			}
		case "volume":
			datacenter, _, err := client.Must().CloudClient.DataCentersApi.DatacentersGet(context.Background()).Execute()
			if err != nil {
				multiErr = errors.Join(multiErr, fmt.Errorf("error occurred getting datacenter with label ID: %v. error: %w", resourceId, err))
				continue
			}
			_, err = client.Must().CloudClient.LabelsApi.DatacentersVolumesLabelsDelete(context.Background(), *datacenter.Id,
				resourceId, key).Execute()
			if err != nil {
				multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, resourceId, err))
				continue
			}
		case "server":
			datacenter, _, err := client.Must().CloudClient.DataCentersApi.DatacentersGet(context.Background()).Execute()
			if err != nil {
				multiErr = errors.Join(multiErr, fmt.Errorf("error occurred getting datacenter with label ID: %v. error: %w", resourceId, err))
				continue
			}
			_, err = client.Must().CloudClient.LabelsApi.DatacentersServersLabelsDelete(context.Background(), *datacenter.Id,
				resourceId, key).Execute()
			if err != nil {
				multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, resourceId, err))
				continue
			}
		case "ipblock":
			_, err = client.Must().CloudClient.LabelsApi.IpblocksLabelsDelete(context.Background(),
				resourceId, key).Execute()
			if err != nil {
				multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, resourceId, err))
				continue
			}
		case "image":
			_, err = client.Must().CloudClient.LabelsApi.ImagesLabelsDelete(context.Background(),
				resourceId, key).Execute()
			if err != nil {
				multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, resourceId, err))
				continue
			}
		case "snapshot":
			_, err = client.Must().CloudClient.LabelsApi.SnapshotsLabelsDelete(context.Background(),
				resourceId, key).Execute()
			if err != nil {
				multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, resourceId, err))
				continue
			}
		}
		if multiErr != nil {
			return multiErr
		}
	}
	return nil
}
